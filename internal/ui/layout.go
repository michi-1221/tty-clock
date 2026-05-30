package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/michi-1221/tty-clock/internal/clock"
	"github.com/michi-1221/tty-clock/internal/render"
)

// defaultWidth/Height are used until a WindowSizeMsg arrives (e.g. non-TTY
// output), so the clock renders instead of showing "terminal too small"
// forever (★6).
const (
	defaultWidth  = 80
	defaultHeight = 24
)

// scaleRefWidth/Height anchor the font scale: a window of this size renders the
// font at ×1 (the original size). Larger windows step up in integer multiples.
const (
	scaleRefWidth  = 54
	scaleRefHeight = 10
)

// scaleFor maps a full window size to an integer font magnification:
// floor(min(w/54, h/10)), at least 1. So 54×10→1, 108×20→2, 162×30→3.
func scaleFor(w, h int) int {
	if s := min(w/scaleRefWidth, h/scaleRefHeight); s > 1 {
		return s
	}
	return 1
}

// renderContext snapshots model state into the pure render input.
func (m Model) renderContext(w, h, scale int) render.RenderContext {
	aspect := m.caps.CellAspect // measured cell height/width
	if m.cfg.CellAspect > 0 {
		aspect = m.cfg.CellAspect // explicit config override wins
	}
	if aspect <= 0 {
		aspect = 2.0 // fallback when the terminal doesn't report pixel size
	}
	return render.RenderContext{
		Now:        clock.NewSnapshot(m.now),
		Theme:      m.theme,
		Format:     m.fmtOpts,
		Gran:       m.gran,
		Width:      w,
		Height:     h,
		Scale:      scale,
		CellAspect: aspect,
		Caps:       m.caps,
	}
}

// view assembles the full screen: a centered clock above a footer.
func (m Model) view() string {
	w, h := m.width, m.height
	if w <= 0 || h <= 0 {
		w, h = defaultWidth, defaultHeight
	}
	r := m.caps.Renderer

	footer := m.footer(w)
	footerH := 0
	if footer != "" {
		footerH = lipgloss.Height(footer)
	}

	availH := h - footerH
	if availH < 1 {
		availH = 1
	}

	// Scale is anchored to the full window size (not the footer-adjusted area).
	body := m.body(w, availH, scaleFor(w, h))
	centered := r.Place(w, availH, lipgloss.Center, lipgloss.Center, body)
	if footer == "" {
		return centered // help hidden and no error: the clock uses the whole screen
	}
	return lipgloss.JoinVertical(lipgloss.Left, centered, footer)
}

// footer is the single-line short help, hidden when the user toggles help off
// with '?'. A transient error line is shown above it (and even when help is
// hidden). Returns "" when there is nothing to show, so view() reserves no row.
func (m Model) footer(w int) string {
	var lines []string
	if m.err != nil {
		lines = append(lines, m.caps.Renderer.NewStyle().Foreground(m.theme.Accent).Render(m.err.Error()))
	}
	if m.showHelp {
		lines = append(lines, m.help.ShortHelpView(m.keys.ShortHelp()))
	}
	if len(lines) == 0 {
		return ""
	}
	return lipgloss.PlaceHorizontal(w, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, lines...))
}

// body renders the clock, degrading to a one-line time and finally to a notice
// when the area is too small — never panicking.
func (m Model) body(w, availH, scale int) string {
	ctx := m.renderContext(w, availH, scale)
	rdr := m.activeRenderer(ctx)
	mw, mh := rdr.MinSize(ctx)
	switch {
	case w >= mw && availH >= mh:
		return rdr.Render(ctx)
	case availH >= 1:
		line := m.compactTime()
		if lipgloss.Width(line) <= w {
			return line
		}
		return "terminal too small"
	default:
		return "terminal too small"
	}
}

// compactTime is the plain one-line fallback (e.g. very small terminals).
func (m Model) compactTime() string {
	sn := clock.NewSnapshot(m.now)
	hour := sn.Hour24
	if !m.fmtOpts.Hour24 {
		hour = sn.Hour12
	}
	var s string
	if m.fmtOpts.ShowSeconds && m.gran == clock.GranSeconds {
		s = fmt.Sprintf("%02d:%02d:%02d", hour, sn.Minute, sn.Second)
	} else {
		s = fmt.Sprintf("%02d:%02d", hour, sn.Minute)
	}
	if !m.fmtOpts.Hour24 && m.fmtOpts.ShowAMPM {
		if sn.IsPM {
			s += " PM"
		} else {
			s += " AM"
		}
	}
	return m.caps.Renderer.NewStyle().Foreground(m.theme.Primary).Render(s)
}
