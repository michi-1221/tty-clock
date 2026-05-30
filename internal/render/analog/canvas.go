// Package analog renders the clock as an analog dial drawn with Unicode braille
// dots. It is a pure render.Renderer: no time.Now(), no bubbletea, deterministic
// for golden tests.
package analog

const brailleBase = 0x2800

// pixelMap maps a sub-cell dot (row y%4, column x%2) to its braille bit, per the
// Unicode dot numbering: dots 1,2,3,7 in the left column; 4,5,6,8 in the right.
var pixelMap = [4][2]uint8{
	{0x01, 0x08}, // dots 1, 4
	{0x02, 0x10}, // dots 2, 5
	{0x04, 0x20}, // dots 3, 6
	{0x40, 0x80}, // dots 7, 8
}

// Canvas is a braille dot grid: wpx×hpx dots packed into cols×rows terminal
// cells (2 dots wide, 4 dots tall per cell). Braille dots are ~square, so a
// circle drawn with equal radii looks round without aspect correction.
type Canvas struct {
	wpx, hpx   int
	cols, rows int
	cells      [][]uint8
}

// NewCanvas allocates a canvas for wpx×hpx dots (ceil division to cells).
func NewCanvas(wpx, hpx int) *Canvas {
	if wpx < 1 {
		wpx = 1
	}
	if hpx < 1 {
		hpx = 1
	}
	cols := (wpx + 1) / 2
	rows := (hpx + 3) / 4
	cells := make([][]uint8, rows)
	for r := range cells {
		cells[r] = make([]uint8, cols)
	}
	return &Canvas{wpx: wpx, hpx: hpx, cols: cols, rows: rows, cells: cells}
}

// Set turns on the dot at (x,y). Out-of-range coordinates are a no-op, so every
// drawing routine can funnel through Set without bounds checks and never panics.
func (c *Canvas) Set(x, y int) {
	if x < 0 || y < 0 || x >= c.wpx || y >= c.hpx {
		return
	}
	c.cells[y/4][x/2] |= pixelMap[y%4][x%2]
}

// Test reports whether the dot at (x,y) is on.
func (c *Canvas) Test(x, y int) bool {
	if x < 0 || y < 0 || x >= c.wpx || y >= c.hpx {
		return false
	}
	return c.cells[y/4][x/2]&pixelMap[y%4][x%2] != 0
}

// Clear turns off the dot at (x,y).
func (c *Canvas) Clear(x, y int) {
	if x < 0 || y < 0 || x >= c.wpx || y >= c.hpx {
		return
	}
	c.cells[y/4][x/2] &^= pixelMap[y%4][x%2]
}

// Grid returns the canvas as a rows×cols grid of braille runes, so callers can
// overlay other glyphs (e.g. hour numbers) onto specific cells before styling.
func (c *Canvas) Grid() [][]rune {
	g := make([][]rune, c.rows)
	for r := 0; r < c.rows; r++ {
		g[r] = make([]rune, c.cols)
		for col := 0; col < c.cols; col++ {
			g[r][col] = rune(brailleBase | int(c.cells[r][col]))
		}
	}
	return g
}

// String renders the canvas as rows of braille runes — every row exactly cols
// runes wide (empty cells are U+2800 so the width stays constant).
func (c *Canvas) String() string { return gridString(c.Grid()) }
