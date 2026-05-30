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

func TestDrawDigits(t *testing.T) {
	c := NewCanvas(60, 60)
	drawDigits(c, 30, 30, "8")
	want := 0
	for _, row := range digitFont['8'] {
		want += strings.Count(row, "#")
	}
	got := 0
	for x := 0; x < 60; x++ {
		for y := 0; y < 60; y++ {
			if c.Test(x, y) {
				got++
			}
		}
	}
	if got != want {
		t.Errorf("drawDigits(\"8\") lit %d dots, want %d (one per glyph '#')", got, want)
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
