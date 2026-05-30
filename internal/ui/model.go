// Package ui is the only package that imports bubbletea. It owns the state
// machine: Model + Init/Update/View, key toggles, and the tick loop.
package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/michi-1221/tty-clock/internal/caps"
	"github.com/michi-1221/tty-clock/internal/clock"
	"github.com/michi-1221/tty-clock/internal/config"
	"github.com/michi-1221/tty-clock/internal/render"
	"github.com/michi-1221/tty-clock/internal/render/analog"
	"github.com/michi-1221/tty-clock/internal/render/digital"
	"github.com/michi-1221/tty-clock/internal/theme"
)

// Model is the whole application state (one struct, one reducer).
type Model struct {
	cfg        config.Config
	theme      theme.Theme
	gran       clock.Granularity
	mode       clock.Mode
	fmtOpts    config.FormatOptions
	digital    render.Renderer // both renderers are held; activeRenderer picks one
	analog     render.Renderer
	now        time.Time // latest tick instant; View never calls time.Now()
	tickGen    uint64
	width      int
	height     int
	caps       caps.Capabilities
	keys       keyMap
	help       help.Model
	showHelp   bool
	configPath string // kept for phase-2 reload
	err        error  // transient, shown in the footer; cleared on next success
	quitting   bool
}

// New builds the initial model from a loaded config. It resolves the theme,
// mode, and granularity, and — critically — seeds now with time.Now() so the
// very first frame (drawn before the first tick) shows the correct time (★1).
func New(cfg config.Config, configPath string, c caps.Capabilities) Model {
	mode, _ := clock.ParseMode(cfg.Mode)
	gran, _ := clock.ParseGranularity(cfg.Granularity)

	th, err := theme.Resolve(cfg.Theme, cfg.CustomTheme)
	if err != nil {
		th, _ = theme.Resolve(theme.DefaultPreset, nil) // fall back, surface err in footer
	}

	return Model{
		cfg:        cfg,
		theme:      th,
		gran:       gran,
		mode:       mode,
		fmtOpts:    cfg.Format,
		digital:    digital.New(),
		analog:     analog.New(),
		now:        time.Now(),
		caps:       c,
		keys:       defaultKeyMap(),
		help:       help.New(),
		showHelp:   true, // help line visible by default; '?' hides it
		configPath: configPath,
		err:        err,
	}
}

// Init starts the tick loop. tea.Every aligns the first tick to the next
// wall-clock boundary; now is already seeded so the initial frame is correct.
func (m Model) Init() tea.Cmd {
	return tickCmd(m.gran, m.tickGen)
}

// Update is the reducer: messages in, (new model, command) out.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if msg.gen != m.tickGen {
			return m, nil // drop a stale generation's tick (★2)
		}
		m.now = msg.t
		return m, tickCmd(m.gran, m.tickGen) // always re-arm
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width // needed for help truncation (★8)
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		m.quitting = true
		return m, tea.Quit
	case key.Matches(msg, m.keys.ToggleSecs):
		m.fmtOpts.ShowSeconds = !m.fmtOpts.ShowSeconds
		m.err = nil
		return m, nil
	case key.Matches(msg, m.keys.CycleTheme):
		if th, err := theme.Resolve(theme.Next(m.theme.Name), nil); err == nil {
			m.theme = th
			m.err = nil
		}
		return m, nil
	case key.Matches(msg, m.keys.Help):
		m.showHelp = !m.showHelp
		return m, nil
	case key.Matches(msg, m.keys.ToggleMode):
		if m.mode == clock.ModeDigital {
			m.mode = clock.ModeAnalog
		} else {
			m.mode = clock.ModeDigital
		}
		m.err = nil
		return m, nil
	case key.Matches(msg, m.keys.Reload):
		nm, cmd := m.reload()
		return nm, cmd
	}
	return m, nil
}

// reload re-reads the config file and rebuilds ALL derived state from it, so the
// file is the single source of truth — ephemeral s/t/m toggles are discarded.
// It never crashes: any failure stays in m.err and leaves the running clock
// untouched. The tick is only re-armed when the granularity actually changed.
func (m Model) reload() (Model, tea.Cmd) {
	// Started with no file (pure defaults)? Pick up a file created since launch
	// and adopt it for this and future reloads.
	if m.configPath == "" {
		if p, err := config.Resolve(""); err == nil && p != "" {
			m.configPath = p
		}
	}

	cfg, err := config.Load(m.configPath) // Load("") returns DefaultConfig(), nil
	if err != nil {
		m.err = err // read/syntax/validate failure: keep everything, don't touch tick
		return m, nil
	}
	// Validate() does not check the theme, so a valid file can still name an
	// unknown theme — treat that as all-or-nothing and keep the running theme.
	th, terr := theme.Resolve(cfg.Theme, cfg.CustomTheme)
	if terr != nil {
		m.err = terr
		return m, nil
	}

	// mode/gran/font are known-valid here (Validate ran inside Load), so the
	// parse errors are impossible — ignore them for symmetry with New().
	newMode, _ := clock.ParseMode(cfg.Mode)
	newGran, _ := clock.ParseGranularity(cfg.Granularity)
	oldGran := m.gran

	m.cfg = cfg
	m.theme = th
	m.mode = newMode
	m.gran = newGran
	m.fmtOpts = cfg.Format
	m.err = nil

	// Re-arm only when granularity changed: bump the generation so the in-flight
	// old-gen tick is dropped (and does not re-arm), starting exactly one new
	// loop. Unchanged gran → return nil so the existing loop keeps running (never
	// spawn a parallel tea.Every).
	if newGran != oldGran {
		m.tickGen++
		return m, tickCmd(m.gran, m.tickGen)
	}
	return m, nil
}

// activeRenderer picks the renderer for the current mode, falling back to
// digital when analog can't be drawn: a non-UTF-8 terminal (braille has no
// ASCII form) or an area below analog's MinSize. mode stays in m.mode, so the
// analog face returns automatically once the terminal can show it again.
func (m Model) activeRenderer(ctx render.RenderContext) render.Renderer {
	if m.mode == clock.ModeAnalog && m.caps.Unicode {
		if aw, ah := m.analog.MinSize(ctx); ctx.Width >= aw && ctx.Height >= ah {
			return m.analog
		}
	}
	return m.digital
}

// View renders the current frame (pure function of the model).
func (m Model) View() string {
	if m.quitting {
		return "" // avoid a final flash inside the alt screen before it's torn down
	}
	return m.view()
}
