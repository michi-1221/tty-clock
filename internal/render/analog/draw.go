package analog

import "math"

// line draws an integer Bresenham segment, inclusive of both endpoints.
func line(c *Canvas, x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := sign(x1 - x0)
	sy := sign(y1 - y0)
	err := dx + dy
	for {
		c.Set(x0, y0)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

// ellipsePoint returns the dot at clock angle deg (0 = 12 o'clock, increasing
// clockwise) on the ellipse with radii (rx,ry) about (cx,cy). Screen y grows
// downward, so the vertical term is negated.
func ellipsePoint(cx, cy, rx, ry int, deg float64) (int, int) {
	rad := deg * math.Pi / 180
	// The explicit float64() around each product rounds it before the add,
	// defeating FMA contraction (the arm64 backend fuses a*b+c, amd64 does not).
	// Without this, iround flips at .5 boundaries and the dial renders a few
	// dots differently per architecture (golden tests then fail off-arm64).
	x := float64(cx) + float64(float64(rx)*math.Sin(rad))
	y := float64(cy) - float64(float64(ry)*math.Cos(rad))
	return iround(x), iround(y)
}

// ellipse draws an ellipse outline by sweeping the angle finely enough to stay
// gapless (deterministic; no floating-point midpoint state to drift).
func ellipse(c *Canvas, cx, cy, rx, ry int) {
	steps := 4 * (rx + ry)
	if steps < 8 {
		steps = 8
	}
	for i := 0; i < steps; i++ {
		x, y := ellipsePoint(cx, cy, rx, ry, 360*float64(i)/float64(steps))
		c.Set(x, y)
	}
}

// hand draws a hand from the center out to frac of the dial radii at deg.
func hand(c *Canvas, cx, cy, rx, ry int, frac, deg float64) {
	x, y := ellipsePoint(cx, cy, iround(float64(rx)*frac), iround(float64(ry)*frac), deg)
	line(c, cx, cy, x, y)
}

// ticks draws the 12 hour marks just inside the rim.
func ticks(c *Canvas, cx, cy, rx, ry int) {
	for k := 0; k < 12; k++ {
		deg := float64(k) * 30
		ox, oy := ellipsePoint(cx, cy, rx, ry, deg)
		ix, iy := ellipsePoint(cx, cy, iround(float64(rx)*0.88), iround(float64(ry)*0.88), deg)
		line(c, ix, iy, ox, oy)
	}
}

func iround(f float64) int { return int(math.Round(f)) }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	default:
		return 0
	}
}
