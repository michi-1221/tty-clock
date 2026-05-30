package digital

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func allFonts() []Font { return []Font{fonts["block"], fonts["ascii"]} }

// requiredRunes are every rune Compose may be asked to draw.
var requiredRunes = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':'}

func TestGlyphInvariants(t *testing.T) {
	for _, f := range allFonts() {
		for _, r := range requiredRunes {
			g, ok := f.Glyphs[r]
			if !ok {
				t.Errorf("font %q missing glyph %q", f.Name, r)
				continue
			}
			if len(g) != f.Rows {
				t.Errorf("font %q glyph %q has %d rows, want %d", f.Name, r, len(g), f.Rows)
			}
			w := lipgloss.Width(g[0])
			for i, row := range g {
				if got := lipgloss.Width(row); got != w {
					t.Errorf("font %q glyph %q row %d width %d, want %d (uniform)", f.Name, r, i, got, w)
				}
			}
		}
	}
}

func TestComposeNoDigitJitter(t *testing.T) {
	for _, f := range allFonts() {
		a := f.Compose("00:00:00", false)
		b := f.Compose("23:59:59", false)
		if lipgloss.Width(a) != lipgloss.Width(b) || lipgloss.Height(a) != lipgloss.Height(b) {
			t.Errorf("font %q: dims differ 00:00:00 (%dx%d) vs 23:59:59 (%dx%d)",
				f.Name, lipgloss.Width(a), lipgloss.Height(a), lipgloss.Width(b), lipgloss.Height(b))
		}
	}
}

func TestComposeBlinkSameDimensions(t *testing.T) {
	for _, f := range allFonts() {
		on := f.Compose("12:34:56", false)
		off := f.Compose("12:34:56", true)
		if lipgloss.Width(on) != lipgloss.Width(off) || lipgloss.Height(on) != lipgloss.Height(off) {
			t.Errorf("font %q: blink changes dims on=%dx%d off=%dx%d",
				f.Name, lipgloss.Width(on), lipgloss.Height(on), lipgloss.Width(off), lipgloss.Height(off))
		}
		if on == off {
			t.Errorf("font %q: blink-off should differ from blink-on", f.Name)
		}
	}
}

func TestSelectFallback(t *testing.T) {
	if Select("block", false).Name != "ascii" {
		t.Error("block + no-unicode should fall back to ascii")
	}
	if Select("block", true).Name != "block" {
		t.Error("block + unicode should stay block")
	}
	if Select("ascii", true).Name != "ascii" {
		t.Error("explicit ascii should always be ascii")
	}
}

func TestScaleArt(t *testing.T) {
	base := fonts["block"].Compose("12:34", false)
	bw, bh := lipgloss.Width(base), lipgloss.Height(base)
	for _, f := range []int{1, 2, 3} {
		got := scaleArt(base, f)
		if w, h := lipgloss.Width(got), lipgloss.Height(got); w != bw*f || h != bh*f {
			t.Errorf("scaleArt(f=%d) = %dx%d, want %dx%d", f, w, h, bw*f, bh*f)
		}
	}
	if scaleArt(base, 1) != base {
		t.Error("scaleArt(f=1) should be identity")
	}
	if scaleArt(base, 0) != base {
		t.Error("scaleArt(f<=0) should be identity")
	}
}
