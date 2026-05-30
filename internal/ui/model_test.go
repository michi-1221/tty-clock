package ui

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/michi-1221/tty-clock/internal/caps"
	"github.com/michi-1221/tty-clock/internal/clock"
	"github.com/michi-1221/tty-clock/internal/config"
)

func testCaps() caps.Capabilities {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.Ascii)
	return caps.New(r, false) // unicode=false → ascii/digital
}

func unicodeCaps() caps.Capabilities {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.Ascii)
	return caps.New(r, true) // unicode=true → analog allowed
}

func newModel(t *testing.T) Model {
	t.Helper()
	return New(config.DefaultConfig(), "", testCaps())
}

func runeKey(r rune) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}} }

func TestNewSeedsNow(t *testing.T) {
	m := newModel(t)
	if m.now.IsZero() {
		t.Error("New() must seed now (★1) so the first frame is correct")
	}
}

func TestTickReArms(t *testing.T) {
	m := newModel(t)
	at := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	nm, cmd := m.Update(tickMsg{t: at, gen: m.tickGen})
	got := nm.(Model)
	if !got.now.Equal(at) {
		t.Errorf("now = %v, want %v", got.now, at)
	}
	if cmd == nil {
		t.Error("a current-generation tick MUST return a re-arm command (★2)")
	}
}

func TestStaleTickDropped(t *testing.T) {
	m := newModel(t)
	orig := m.now
	nm, cmd := m.Update(tickMsg{t: time.Now().Add(time.Hour), gen: m.tickGen + 99})
	got := nm.(Model)
	if !got.now.Equal(orig) {
		t.Error("stale-generation tick must not change now")
	}
	if cmd != nil {
		t.Error("stale-generation tick must not re-arm (would multiply loops)")
	}
}

func TestToggleSeconds(t *testing.T) {
	m := newModel(t)
	before := m.fmtOpts.ShowSeconds
	nm, _ := m.Update(runeKey('s'))
	if nm.(Model).fmtOpts.ShowSeconds == before {
		t.Error("'s' should toggle ShowSeconds")
	}
}

func TestCycleTheme(t *testing.T) {
	m := newModel(t)
	before := m.theme.Name
	nm, _ := m.Update(runeKey('t'))
	if nm.(Model).theme.Name == before {
		t.Errorf("'t' should change theme from %q", before)
	}
}

func TestToggleHelp(t *testing.T) {
	m := newModel(t)
	if !m.showHelp {
		t.Fatal("help line should be visible by default")
	}
	nm, _ := m.Update(runeKey('?'))
	if nm.(Model).showHelp {
		t.Error("'?' should hide the help line when visible")
	}
	nm2, _ := nm.(Model).Update(runeKey('?'))
	if !nm2.(Model).showHelp {
		t.Error("'?' again should show the help line")
	}
}

func TestHiddenHelpRemovesFooter(t *testing.T) {
	m := newModel(t)
	nm, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	withHelp := nm.(Model).View()
	hidden, _ := nm.(Model).Update(runeKey('?'))
	withoutHelp := hidden.(Model).View()
	if strings.Contains(withoutHelp, "seconds") || strings.Contains(withoutHelp, "theme") {
		t.Error("hiding help should remove the seconds/theme labels from the footer")
	}
	if lipgloss.Height(withoutHelp) < lipgloss.Height(withHelp) {
		// sanity: hiding the footer shouldn't add rows
		t.Logf("with help %d rows, without %d", lipgloss.Height(withHelp), lipgloss.Height(withoutHelp))
	}
}

func TestQuit(t *testing.T) {
	m := newModel(t)
	nm, cmd := m.Update(runeKey('q'))
	if !nm.(Model).quitting {
		t.Error("'q' should set quitting")
	}
	if cmd == nil {
		t.Error("'q' should return tea.Quit")
	}
	if msg := cmd(); msg == nil {
		t.Error("quit command should yield a message")
	}
}

func TestModeToggle(t *testing.T) {
	m := newModel(t)
	if m.mode != clock.ModeDigital {
		t.Fatal("default mode should be digital")
	}
	nm, _ := m.Update(runeKey('m'))
	if nm.(Model).mode != clock.ModeAnalog {
		t.Error("'m' should switch to analog")
	}
	nm2, _ := nm.(Model).Update(runeKey('m'))
	if nm2.(Model).mode != clock.ModeDigital {
		t.Error("'m' again should switch back to digital")
	}
}

func TestActiveRendererForcesDigitalWithoutUnicode(t *testing.T) {
	m := New(config.DefaultConfig(), "", testCaps()) // unicode=false
	m.mode = clock.ModeAnalog
	if got := m.activeRenderer(m.renderContext(80, 40, 1)).Name(); got != "digital" {
		t.Errorf("non-UTF-8 terminal must force digital even in analog mode, got %q", got)
	}
}

func TestActiveRendererAnalogAndFallback(t *testing.T) {
	m := New(config.DefaultConfig(), "", unicodeCaps())
	m.mode = clock.ModeAnalog
	if got := m.activeRenderer(m.renderContext(80, 40, 1)).Name(); got != "analog" {
		t.Errorf("large UTF-8 analog should select analog, got %q", got)
	}
	if got := m.activeRenderer(m.renderContext(10, 6, 1)).Name(); got != "digital" {
		t.Errorf("analog below MinSize should fall back to digital, got %q", got)
	}
	if m.mode != clock.ModeAnalog {
		t.Error("fallback must preserve mode for recovery")
	}
}

func TestViewWithoutSizeDoesNotPanic(t *testing.T) {
	m := newModel(t)
	// No WindowSizeMsg yet: width/height are 0; must fall back, not break (★6).
	out := m.View()
	if strings.TrimSpace(out) == "" {
		t.Error("View with no size should still render (default size fallback)")
	}
}

func TestViewTinyTerminal(t *testing.T) {
	m := newModel(t)
	nm, _ := m.Update(tea.WindowSizeMsg{Width: 3, Height: 2})
	out := nm.(Model).View()
	if strings.TrimSpace(out) == "" {
		t.Error("tiny terminal should still render something (no panic, no blank)")
	}
}

func TestScaleFor(t *testing.T) {
	cases := []struct{ w, h, want int }{
		{54, 10, 1},  // reference → ×1
		{108, 20, 2}, // both 2× → ×2
		{162, 30, 3}, // both 3× → ×3
		{110, 15, 1}, // height short → height-limited to ×1
		{40, 8, 1},   // below reference → floor at 1
		{53, 9, 1},
	}
	for _, c := range cases {
		if got := scaleFor(c.w, c.h); got != c.want {
			t.Errorf("scaleFor(%d,%d) = %d, want %d", c.w, c.h, got, c.want)
		}
	}
}

func TestLargeWindowScalesBody(t *testing.T) {
	m := newModel(t)
	small := m.body(60, 11, scaleFor(60, 12))
	big := m.body(130, 25, scaleFor(130, 26))
	if lipgloss.Height(big) <= lipgloss.Height(small) {
		t.Errorf("body height small=%d big=%d: a larger window should scale the clock up",
			lipgloss.Height(small), lipgloss.Height(big))
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReloadAppliesFileValues(t *testing.T) {
	path := writeConfig(t, `{"granularity":"minutes","theme":"dracula","format":{"hour24":false}}`)
	m := New(config.DefaultConfig(), path, testCaps())
	nm, _ := m.Update(runeKey('r'))
	got := nm.(Model)
	if got.gran != clock.GranMinutes {
		t.Errorf("gran = %v, want minutes", got.gran)
	}
	if got.theme.Name != "dracula" {
		t.Errorf("theme = %q, want dracula", got.theme.Name)
	}
	if got.fmtOpts.Hour24 {
		t.Error("hour24 should be false from file")
	}
	if got.err != nil {
		t.Errorf("unexpected err: %v", got.err)
	}
}

func TestReloadDiscardsEphemeralToggles(t *testing.T) {
	path := writeConfig(t, `{"theme":"tokyo-night","format":{"showSeconds":true}}`)
	m := New(config.DefaultConfig(), path, testCaps())
	m2, _ := m.Update(runeKey('s'))          // seconds → false (ephemeral)
	m3, _ := m2.(Model).Update(runeKey('t')) // theme → dracula (ephemeral)
	if m3.(Model).fmtOpts.ShowSeconds {
		t.Fatal("precondition: 's' should have turned seconds off")
	}
	nm, _ := m3.(Model).Update(runeKey('r'))
	got := nm.(Model)
	if !got.fmtOpts.ShowSeconds {
		t.Error("reload should restore showSeconds=true from file (toggle discarded)")
	}
	if got.theme.Name != "tokyo-night" {
		t.Errorf("reload should restore theme to tokyo-night, got %q", got.theme.Name)
	}
}

func TestReloadGranularityChangeReArms(t *testing.T) {
	path := writeConfig(t, `{"granularity":"minutes"}`)
	m := New(config.DefaultConfig(), path, testCaps()) // starts at seconds
	gen0 := m.tickGen
	nm, cmd := m.reload()
	if cmd == nil {
		t.Fatal("granularity change must return a re-arm command")
	}
	if nm.tickGen != gen0+1 {
		t.Errorf("tickGen = %d, want %d", nm.tickGen, gen0+1)
	}
	// An in-flight OLD-gen tick arriving after reload must be dropped.
	before := nm.now
	at := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	d, dcmd := nm.Update(tickMsg{t: at, gen: gen0})
	if dcmd != nil {
		t.Error("stale old-gen tick must not re-arm")
	}
	if !d.(Model).now.Equal(before) {
		t.Error("stale old-gen tick must not change now")
	}
	// The new-gen tick is accepted and keeps the single chain alive.
	n2, ncmd := nm.Update(tickMsg{t: at, gen: nm.tickGen})
	if ncmd == nil || !n2.(Model).now.Equal(at) {
		t.Error("new-gen tick should update now and re-arm")
	}
}

func TestReloadGranUnchangedDoesNotReArm(t *testing.T) {
	path := writeConfig(t, `{"granularity":"seconds","theme":"nord"}`)
	m := New(config.DefaultConfig(), path, testCaps()) // already seconds
	gen0 := m.tickGen
	nm, cmd := m.reload()
	if cmd != nil {
		t.Error("unchanged granularity must NOT re-arm (would spawn a parallel loop)")
	}
	if nm.tickGen != gen0 {
		t.Errorf("tickGen changed to %d, want %d", nm.tickGen, gen0)
	}
	if nm.theme.Name != "nord" {
		t.Error("reload should still apply non-granularity changes (theme=nord)")
	}
}

func TestReloadInvalidFileKeepsState(t *testing.T) {
	path := writeConfig(t, `{ broken`)
	m := New(config.DefaultConfig(), path, testCaps())
	gen0, gran0, theme0 := m.tickGen, m.gran, m.theme.Name
	nm, cmd := m.reload()
	if cmd != nil || nm.err == nil {
		t.Fatal("invalid file: want err set and nil cmd")
	}
	if nm.gran != gran0 || nm.theme.Name != theme0 || nm.tickGen != gen0 {
		t.Error("invalid reload must preserve all derived state")
	}
}

func TestReloadUnknownThemeKeepsState(t *testing.T) {
	// Valid JSON that passes Validate() (mode/gran/font ok) but names a theme
	// that theme.Resolve rejects — exercises the two-stage error path.
	path := writeConfig(t, `{"theme":"does-not-exist"}`)
	m := New(config.DefaultConfig(), path, testCaps())
	theme0 := m.theme.Name
	nm, cmd := m.reload()
	if cmd != nil || nm.err == nil {
		t.Fatal("unknown theme: want err set and nil cmd")
	}
	if nm.theme.Name != theme0 {
		t.Errorf("theme should be preserved (%q), got %q", theme0, nm.theme.Name)
	}
}

func TestReloadClearsPriorError(t *testing.T) {
	path := writeConfig(t, `{"theme":"gruvbox"}`)
	m := New(config.DefaultConfig(), path, testCaps())
	m.err = io.ErrUnexpectedEOF // simulate a prior transient error
	nm, _ := m.reload()
	if nm.err != nil {
		t.Errorf("successful reload should clear err, got %v", nm.err)
	}
}

func TestReloadKeyWired(t *testing.T) {
	path := writeConfig(t, `{"theme":"gruvbox"}`)
	m := New(config.DefaultConfig(), path, testCaps())
	nm, _ := m.Update(runeKey('r'))
	if nm.(Model).theme.Name != "gruvbox" {
		t.Error("'r' must trigger reload (theme should become gruvbox)")
	}
}

func TestQuittingViewEmpty(t *testing.T) {
	m := newModel(t)
	nm, _ := m.Update(runeKey('q'))
	if nm.(Model).View() != "" {
		t.Error("quitting View should be empty to avoid a final flash")
	}
}
