package render

import (
	"io"
	"regexp"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/michi-1221/tty-clock/internal/theme"
)

func trueColor() *lipgloss.Renderer {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.TrueColor)
	return r
}

var (
	ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	fgRe   = regexp.MustCompile(`38;2;\d+;\d+;\d+`)
)

func strip(s string) string   { return ansiRe.ReplaceAllString(s, "") }
func distinctFG(s string) int { return len(set(fgRe.FindAllString(s, -1))) }
func set(xs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		m[x] = struct{}{}
	}
	return m
}

func mustResolve(t *testing.T, name string) theme.Theme {
	t.Helper()
	th, err := theme.Resolve(name, nil)
	if err != nil {
		t.Fatal(err)
	}
	return th
}

func TestColorizeFlatIsTransparentUnderAscii(t *testing.T) {
	// A flat theme under the Ascii profile must not alter the text at all, so
	// the golden tests of the renderers stay byte-stable.
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.Ascii)
	art := "AB\nCD"
	if got := Colorize(r, art, mustResolve(t, "tokyo-night"), lipgloss.Color("#ffffff")); got != art {
		t.Errorf("flat Ascii Colorize changed text:\n got %q\nwant %q", got, art)
	}
}

func TestColorizeVerticalUsesDistinctColorsPerRow(t *testing.T) {
	art := "AA\nBB\nCC" // 3 rows → 3 gradient fractions
	out := Colorize(trueColor(), art, mustResolve(t, "sunset"), lipgloss.Color("#000000"))
	if strip(out) != art {
		t.Errorf("visible text changed: %q", strip(out))
	}
	if n := distinctFG(out); n < 2 {
		t.Errorf("vertical gradient used %d distinct colors, want ≥2", n)
	}
}

func TestColorizeHorizontalUsesDistinctColorsPerCell(t *testing.T) {
	art := "ABCDEF" // single row, 6 cells
	out := Colorize(trueColor(), art, mustResolve(t, "rainbow"), lipgloss.Color("#000000"))
	if strip(out) != art {
		t.Errorf("visible text changed: %q", strip(out))
	}
	if n := distinctFG(out); n < 2 {
		t.Errorf("horizontal gradient used %d distinct colors, want ≥2", n)
	}
}

func TestColorizePreservesWidth(t *testing.T) {
	art := "12:34:56"
	out := Colorize(trueColor(), art, mustResolve(t, "rainbow"), lipgloss.Color("#000000"))
	if w := lipgloss.Width(out); w != lipgloss.Width(art) {
		t.Errorf("colorized width = %d, want %d", w, lipgloss.Width(art))
	}
}
