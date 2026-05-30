package analog

import (
	"strconv"
	"strings"

	"github.com/michi-1221/tty-clock/internal/clock"
	"github.com/michi-1221/tty-clock/internal/render"
)

// AnalogRenderer draws an analog dial (face, 12 ticks, hour numbers, and
// hour/minute/second hands) with braille dots. v1 is single-color (採用 #2);
// per-role colors come later.
type AnalogRenderer struct{}

// New returns a stateless analog renderer.
func New() AnalogRenderer { return AnalogRenderer{} }

// Name implements render.Renderer.
func (AnalogRenderer) Name() string { return "analog" }

// minCols/minRows is the smallest legible dial, in cells. Constant and
// independent of ctx.Scale, so the ui fallback decision is stable.
const (
	minCols = 20
	minRows = 10
)

// MinSize implements render.Renderer.
func (AnalogRenderer) MinSize(render.RenderContext) (w, h int) { return minCols, minRows }

// dialRadii returns the largest visually-round ellipse radii (in dots) that fit
// a wpx×hpx canvas for the given cell aspect (height/width). A round dial needs
// rx/ry = aspect/2 in dot units (braille dots are square only when aspect == 2).
func dialRadii(wpx, hpx int, aspect float64) (rx, ry int) {
	if aspect <= 0 {
		aspect = 2.0
	}
	maxX, maxY := wpx/2-1, hpx/2-1
	ry = maxY
	rx = iround(float64(ry) * aspect / 2)
	if rx > maxX { // width-limited: shrink to fit, then recompute ry to stay round
		rx = maxX
		ry = iround(float64(rx) * 2 / aspect)
	}
	if rx < 1 {
		rx = 1
	}
	if ry < 1 {
		ry = 1
	}
	return rx, ry
}

// Render implements render.Renderer. The dial fills the available cell area
// (ctx.Scale is ignored) and is shaped round via ctx.CellAspect.
func (AnalogRenderer) Render(ctx render.RenderContext) string {
	cols, rows := max(ctx.Width, minCols), max(ctx.Height, minRows)
	wpx, hpx := cols*2, rows*4
	cx, cy := wpx/2, hpx/2
	rx, ry := dialRadii(wpx, hpx, ctx.CellAspect)

	c := NewCanvas(wpx, hpx)
	ellipse(c, cx, cy, rx, ry)
	ticks(c, cx, cy, rx, ry)
	if ctx.Format.ShowNumbers {
		drawNumbers(c, cx, cy, rx, ry) // hour numbers as braille dots, matching the dial
	}
	hand(c, cx, cy, rx, ry, 0.55, hourAngle(ctx.Now)) // hour (shortest)
	hand(c, cx, cy, rx, ry, 0.80, minAngle(ctx.Now))  // minute
	if ctx.Format.ShowSeconds && ctx.Gran == clock.GranSeconds {
		hand(c, cx, cy, rx, ry, 0.92, secAngle(ctx.Now)) // second hand; hidden by 's'
	}
	return ctx.Caps.Renderer.NewStyle().Foreground(ctx.Theme.Accent).Render(c.String())
}

// digitFont is a 3×5 dot matrix for 0–9, plotted into the braille canvas so the
// hour numbers are pixelated like the rest of the dial. '#' marks a lit dot.
var digitFont = map[rune][]string{
	'0': {"###", "# #", "# #", "# #", "###"},
	'1': {" # ", " # ", " # ", " # ", " # "},
	'2': {"###", "  #", "###", "#  ", "###"},
	'3': {"###", "  #", "###", "  #", "###"},
	'4': {"# #", "# #", "###", "  #", "  #"},
	'5': {"###", "#  ", "###", "  #", "###"},
	'6': {"###", "#  ", "###", "# #", "###"},
	'7': {"###", "  #", "  #", "  #", "  #"},
	'8': {"###", "# #", "###", "# #", "###"},
	'9': {"###", "# #", "###", "  #", "###"},
}

const (
	digitW   = 3
	digitH   = 5
	digitGap = 1
)

// drawNumbers plots the hour numbers 1–12 as braille dots just inside the ticks.
func drawNumbers(c *Canvas, cx, cy, rx, ry int) {
	lx, ly := iround(float64(rx)*0.72), iround(float64(ry)*0.72)
	for h := 1; h <= 12; h++ {
		x, y := ellipsePoint(cx, cy, lx, ly, float64(h%12)*30)
		drawDigits(c, x, y, strconv.Itoa(h))
	}
}

// drawDigits plots s centered at (cx,cy) using the dot-matrix digit font.
func drawDigits(c *Canvas, cx, cy int, s string) {
	total := len(s)*digitW + (len(s)-1)*digitGap
	x0 := cx - total/2
	y0 := cy - digitH/2
	for _, ch := range s {
		if glyph, ok := digitFont[ch]; ok {
			for row := 0; row < digitH; row++ {
				for col := 0; col < digitW && col < len(glyph[row]); col++ {
					if glyph[row][col] == '#' {
						c.Set(x0+col, y0+row)
					}
				}
			}
		}
		x0 += digitW + digitGap
	}
}

func gridString(g [][]rune) string {
	out := make([]string, len(g))
	for i, row := range g {
		out[i] = string(row)
	}
	return strings.Join(out, "\n")
}

// Clock-face angles in degrees: 0 = 12 o'clock, increasing clockwise.
func secAngle(n clock.TimeSnapshot) float64  { return 6 * float64(n.Second) }
func minAngle(n clock.TimeSnapshot) float64  { return 6*float64(n.Minute) + 0.1*float64(n.Second) }
func hourAngle(n clock.TimeSnapshot) float64 { return 30*float64(n.Hour12%12) + 0.5*float64(n.Minute) }
