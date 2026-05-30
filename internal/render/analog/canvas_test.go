package analog

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestCanvasRoundTrip(t *testing.T) {
	for x := 0; x < 2; x++ {
		for y := 0; y < 4; y++ {
			c := NewCanvas(2, 4) // exactly one cell
			c.Set(x, y)
			if !c.Test(x, y) {
				t.Errorf("Set/Test(%d,%d) not on", x, y)
			}
			if c.cells[0][0] != pixelMap[y][x] {
				t.Errorf("dot(%d,%d): cell=%#x want %#x", x, y, c.cells[0][0], pixelMap[y][x])
			}
			c.Clear(x, y)
			if c.Test(x, y) {
				t.Errorf("Clear(%d,%d) still on", x, y)
			}
		}
	}
}

func TestCanvasSinglePixelRune(t *testing.T) {
	c := NewCanvas(2, 4)
	c.Set(0, 0) // dot 1 → U+2801
	if got := []rune(c.String())[0]; got != rune(0x2801) {
		t.Errorf("Set(0,0) rune = %#x, want 0x2801", got)
	}
	c2 := NewCanvas(2, 4)
	c2.Set(1, 3) // dot 8 → U+2880
	if got := []rune(c2.String())[0]; got != rune(0x2880) {
		t.Errorf("Set(1,3) rune = %#x, want 0x2880", got)
	}
}

func TestCanvasOutOfRangeNoop(t *testing.T) {
	c := NewCanvas(4, 4)
	for _, p := range [][2]int{{-1, 0}, {0, -1}, {4, 0}, {0, 4}, {100, 100}} {
		c.Set(p[0], p[1]) // must not panic
		if c.Test(p[0], p[1]) {
			t.Errorf("out-of-range (%d,%d) reads on", p[0], p[1])
		}
	}
	for _, line := range strings.Split(c.String(), "\n") {
		for _, r := range line {
			if r != rune(brailleBase) {
				t.Errorf("canvas not blank after out-of-range sets: %#x", r)
			}
		}
	}
}

func TestCanvasStringDimensions(t *testing.T) {
	c := NewCanvas(41, 43) // odd dims → ceil: cols=21, rows=11
	if c.cols != 21 || c.rows != 11 {
		t.Fatalf("cols,rows = %d,%d want 21,11", c.cols, c.rows)
	}
	lines := strings.Split(c.String(), "\n")
	if len(lines) != 11 {
		t.Errorf("got %d rows, want 11", len(lines))
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 21 {
			t.Errorf("row %d width %d, want 21", i, w)
		}
	}
}
