package digital

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Font is a fixed-height bitmap font: every glyph has exactly Rows rows, and
// all rows of a single glyph share the same display width.
type Font struct {
	Name   string
	Rows   int
	Glyphs map[rune][]string
}

var fonts = map[string]Font{
	"block": {Name: "block", Rows: 5, Glyphs: Block5},
	"ascii": {Name: "ascii", Rows: 5, Glyphs: ASCII5},
}

// Select picks the font to use. A "block" (or not-yet-implemented "seg7", or
// empty) request degrades to "ascii" when the terminal can't render block
// glyphs. An explicit "ascii" is always honored (manual escape hatch, ★7).
func Select(name string, unicode bool) Font {
	if name == "ascii" {
		return fonts["ascii"]
	}
	if !unicode {
		return fonts["ascii"]
	}
	return fonts["block"]
}

// Compose renders a string of digit/colon runes into a multi-line block by
// joining glyphs horizontally, separated by a uniform 1-cell gap so adjacent
// strokes don't merge into one bar. When blinkOff is true, ':' is swapped for a
// same-width blank so the colon disappears without shifting any digit (★5).
// Runes with no glyph are skipped defensively (callers pass digits and ':').
func (f Font) Compose(s string, blinkOff bool) string {
	cols := make([]string, 0, len(s)*2)
	first := true
	for _, r := range s {
		g, ok := f.Glyphs[r]
		if !ok {
			continue
		}
		if r == ':' && blinkOff {
			g = blankLike(g)
		}
		if !first {
			cols = append(cols, f.spacer())
		}
		cols = append(cols, strings.Join(g, "\n"))
		first = false
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}

// spacer is a 1-cell-wide, Rows-tall column of blanks placed between glyphs.
func (f Font) spacer() string {
	rows := make([]string, f.Rows)
	for i := range rows {
		rows[i] = " "
	}
	return strings.Join(rows, "\n")
}

// scaleArt magnifies a composed multi-line block by an integer factor: each
// rune is repeated f times horizontally and each row f times vertically. All
// glyph runes are width-1, so cell == rune and the aspect ratio is preserved.
func scaleArt(art string, f int) string {
	if f <= 1 {
		return art
	}
	var out []string
	for _, line := range strings.Split(art, "\n") {
		var b strings.Builder
		b.Grow(len(line) * f)
		for _, r := range line {
			for i := 0; i < f; i++ {
				b.WriteRune(r)
			}
		}
		wide := b.String()
		for i := 0; i < f; i++ {
			out = append(out, wide)
		}
	}
	return strings.Join(out, "\n")
}

// blankLike returns a glyph of the same shape filled with spaces, preserving
// each row's display width.
func blankLike(g []string) []string {
	out := make([]string, len(g))
	for i, row := range g {
		out[i] = strings.Repeat(" ", lipgloss.Width(row))
	}
	return out
}
