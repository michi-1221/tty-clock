package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/michi-1221/tty-clock/internal/theme"
)

// Colorize fills a multi-line art string with the theme's color. With no
// gradient it applies the single flat color. With a gradient it interpolates
// the stops along the theme's axis: vertical colors each row by its top→bottom
// position; horizontal colors each cell by its left→right position. Only a
// foreground color is added, so the visible width is unchanged and callers can
// still center/join the result.
func Colorize(r *lipgloss.Renderer, art string, th theme.Theme, flat lipgloss.Color) string {
	if !th.HasGradient() {
		return r.NewStyle().Foreground(flat).Render(art)
	}
	lines := strings.Split(art, "\n")
	if th.GradDir == theme.Horizontal {
		for i, ln := range lines {
			lines[i] = paintRow(r, ln, th)
		}
		return strings.Join(lines, "\n")
	}
	n := len(lines)
	for i, ln := range lines {
		lines[i] = r.NewStyle().Foreground(th.ColorAt(frac(i, n))).Render(ln)
	}
	return strings.Join(lines, "\n")
}

// paintRow colors each cell of a single line by its horizontal position.
func paintRow(r *lipgloss.Renderer, line string, th theme.Theme) string {
	runes := []rune(line)
	var b strings.Builder
	for j, ch := range runes {
		b.WriteString(r.NewStyle().Foreground(th.ColorAt(frac(j, len(runes)))).Render(string(ch)))
	}
	return b.String()
}

// frac maps index i of n items to [0,1]; a single item sits at 0.
func frac(i, n int) float64 {
	if n <= 1 {
		return 0
	}
	return float64(i) / float64(n-1)
}
