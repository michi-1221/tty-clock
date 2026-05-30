// Package digital renders the clock as giant block-font digits. It is a pure
// render.Renderer: no time.Now(), no bubbletea, deterministic for golden tests.
package digital

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/michi-1221/tui-clock/internal/clock"
	"github.com/michi-1221/tui-clock/internal/render"
)

// DigitalRenderer draws HH:MM(:SS) as stacked glyph rows, with an optional
// AM/PM label and date line beneath (both small text, never block glyphs ★4).
type DigitalRenderer struct{}

// New returns a stateless digital renderer.
func New() DigitalRenderer { return DigitalRenderer{} }

// Name implements render.Renderer.
func (DigitalRenderer) Name() string { return "digital" }

// timeDigits builds the digit/colon string fed to the font. Seconds are shown
// only when both requested and the granularity is per-second (採用 #4).
func timeDigits(ctx render.RenderContext) string {
	hour := ctx.Now.Hour24
	if !ctx.Format.Hour24 {
		hour = ctx.Now.Hour12 // always zero-padded below → no digit jitter (採用 #8)
	}
	if showSeconds(ctx) {
		return fmt.Sprintf("%02d:%02d:%02d", hour, ctx.Now.Minute, ctx.Now.Second)
	}
	return fmt.Sprintf("%02d:%02d", hour, ctx.Now.Minute)
}

func showSeconds(ctx render.RenderContext) bool {
	return ctx.Format.ShowSeconds && ctx.Gran == clock.GranSeconds
}

// blinkOff reports whether the colon should be hidden this frame. Blinking is
// only active at second granularity (採用 #4).
func blinkOff(ctx render.RenderContext) bool {
	return ctx.Format.BlinkColon && ctx.Gran == clock.GranSeconds && !ctx.Now.BlinkOn
}

// Render implements render.Renderer.
func (DigitalRenderer) Render(ctx render.RenderContext) string {
	font := Select(ctx.Format.Font, ctx.Caps.Unicode)
	r := ctx.Caps.Renderer

	base := font.Compose(timeDigits(ctx), blinkOff(ctx))
	art := scaleArt(base, effectiveScale(ctx, font, base))
	parts := []string{r.NewStyle().Foreground(ctx.Theme.Primary).Render(art)}

	if !ctx.Format.Hour24 && ctx.Format.ShowAMPM {
		ampm := "AM"
		if ctx.Now.IsPM {
			ampm = "PM"
		}
		parts = append(parts, r.NewStyle().Foreground(ctx.Theme.Accent).Render(ampm))
	}

	if ctx.Format.ShowDate {
		date := ctx.Now.T.Format(ctx.Format.DateFormat)
		parts = append(parts, r.NewStyle().Foreground(ctx.Theme.Secondary).Render(date))
	}

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

// extraRows counts the unscaled text rows (AM/PM, date) stacked under the
// giant digits.
func extraRows(ctx render.RenderContext) int {
	n := 0
	if !ctx.Format.Hour24 && ctx.Format.ShowAMPM {
		n++
	}
	if ctx.Format.ShowDate {
		n++
	}
	return n
}

// effectiveScale clamps ui's requested Scale (Scale<=0 means ×1) to what fits
// the available area for the current content. Render runs only after ui's
// MinSize check confirms ×1 fits, so the result is always ≥ 1.
func effectiveScale(ctx render.RenderContext, font Font, base string) int {
	s := ctx.Scale
	if s < 1 {
		s = 1
	}
	baseW := lipgloss.Width(base)
	if baseW < 1 {
		baseW = 1
	}
	if fitW := ctx.Width / baseW; fitW >= 1 && fitW < s {
		s = fitW
	}
	if fitH := (ctx.Height - extraRows(ctx)) / font.Rows; fitH >= 1 && fitH < s {
		s = fitH
	}
	if s < 1 {
		s = 1
	}
	return s
}

// MinSize implements render.Renderer: the cell box the current layout needs.
func (DigitalRenderer) MinSize(ctx render.RenderContext) (w, h int) {
	font := Select(ctx.Format.Font, ctx.Caps.Unicode)
	art := font.Compose(timeDigits(ctx), false)
	w = lipgloss.Width(art)
	h = font.Rows
	if !ctx.Format.Hour24 && ctx.Format.ShowAMPM {
		h++
	}
	if ctx.Format.ShowDate {
		if dw := lipgloss.Width(ctx.Now.T.Format(ctx.Format.DateFormat)); dw > w {
			w = dw
		}
		h++
	}
	return w, h
}
