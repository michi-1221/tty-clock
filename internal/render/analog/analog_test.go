package analog

import (
	"strings"
	"testing"
)

func TestDialRadii(t *testing.T) {
	// aspect 2 (cells exactly 1:2) → braille dots square → rx == ry.
	if rx, ry := dialRadii(160, 80, 2.0); rx != ry {
		t.Errorf("aspect 2 should give rx==ry, got %d,%d", rx, ry)
	}
	// taller cells (aspect > 2) → ellipse stretched horizontally: rx > ry.
	if rx, ry := dialRadii(160, 80, 3.0); rx <= ry {
		t.Errorf("aspect 3 should give rx>ry, got %d,%d", rx, ry)
	}
	// must always fit the canvas (and stay ≥ 1).
	rx, ry := dialRadii(40, 80, 3.0)
	if rx < 1 || ry < 1 || rx > 40/2-1 || ry > 80/2-1 {
		t.Errorf("radii must fit canvas and be ≥1, got %d,%d", rx, ry)
	}
}

func TestDrawDigitsBaseSizeMatchesFont(t *testing.T) {
	// At the native 3×5 size, every lit dot maps 1:1 to a glyph '#'.
	base := 0
	for _, row := range digitFont['8'] {
		base += strings.Count(row, "#")
	}
	c := NewCanvas(60, 60)
	drawDigits(c, 30, 30, "8", digitW, digitH)
	got := 0
	for x := 0; x < 60; x++ {
		for y := 0; y < 60; y++ {
			if c.Test(x, y) {
				got++
			}
		}
	}
	if got != base {
		t.Errorf("drawDigits(\"8\") at base size lit %d dots, want %d", got, base)
	}
}

func TestDrawDigitsStaysSymmetric(t *testing.T) {
	// Centered nearest-neighbor resize must keep left/right-symmetric glyphs
	// symmetric — otherwise enlarged numbers look lopsided.
	for _, ch := range []rune{'0', '8', '1'} {
		const tw, th = 5, 8
		c := NewCanvas(40, 40)
		drawDigits(c, 20, 20, string(ch), tw, th)
		x0 := 20 - tw/2
		y0 := 20 - th/2
		for dy := 0; dy < th; dy++ {
			for dx := 0; dx < tw; dx++ {
				l := c.Test(x0+dx, y0+dy)
				r := c.Test(x0+tw-1-dx, y0+dy)
				if l != r {
					t.Errorf("glyph %q row %d not horizontally symmetric at col %d", ch, dy, dx)
				}
			}
		}
	}
}

func TestNumberSizeGrowsWithDial(t *testing.T) {
	// Tiny dials stay at the base size; the glyph grows with the dial but stays
	// below the old ×2 (height 10), and never shrinks below the legible base.
	if tw, th := numberSize(18, 18); tw != digitW || th != digitH {
		t.Errorf("small dial numberSize = %d×%d, want %d×%d", tw, th, digitW, digitH)
	}
	_, th40 := numberSize(39, 39) // matches the analog_40x20 golden
	if th40 <= digitH || th40 >= 2*digitH {
		t.Errorf("40x20 dial height = %d, want between %d and %d", th40, digitH, 2*digitH)
	}
	if _, thBig := numberSize(70, 70); thBig < th40 {
		t.Error("numberSize height should grow with the dial radius")
	}
}

func TestDrawNumbersAreDots(t *testing.T) {
	// Numbers must be braille dots in the canvas, not overlaid text — so the
	// rendered runes are all in the braille block.
	c := NewCanvas(160, 160)
	drawNumbers(c, 80, 80, 70, 70)
	for _, r := range c.String() {
		if r != '\n' && (r < 0x2800 || r > 0x28FF) {
			t.Fatalf("number rendering produced a non-braille rune %#x", r)
		}
	}
}
