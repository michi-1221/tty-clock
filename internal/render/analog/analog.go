package analog

import (
	"math"
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
	return render.Colorize(ctx.Caps.Renderer, c.String(), ctx.Theme, ctx.Theme.Accent)
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
	digitW    = 3
	digitH    = 5
	labelFrac = 0.72 // hour numbers sit at this fraction of the dial radius
)

// numberSize returns the target glyph size (in dots) for the hour numbers. The
// height grows gently with the dial — one extra dot per ~16 dots of radius above
// the legible base digitH — then is capped so adjacent two-digit labels (10–12)
// never collide: their combined width must fit the chord between adjacent hours
// at the label radius, 2·labelFrac·r·sin(15°). Width tracks height at the font's
// 3:5 aspect. Because the glyph is resized (not block-scaled), the size can land
// between the old ×1 and ×2 steps.
func numberSize(rx, ry int) (tw, th int) {
	r := rx
	if ry < r {
		r = ry
	}
	chord := 2 * labelFrac * float64(r) * math.Sin(15*math.Pi/180)
	th = digitH + r/16
	for {
		tw = iround(float64(th) * digitW / digitH)
		if tw < digitW {
			tw = digitW
		}
		if th <= digitH || float64(2*tw+glyphGap(tw)) <= chord {
			break
		}
		th-- // too wide for the chord: shrink a step and re-fit
	}
	return tw, th
}

// glyphGap is the inter-digit spacing for a glyph tw dots wide, keeping the
// font's original 1:3 gap-to-width proportion (at least one dot).
func glyphGap(tw int) int {
	g := iround(float64(tw) / 3)
	if g < 1 {
		g = 1
	}
	return g
}

// srcIndex maps a destination pixel to its source pixel using centered
// nearest-neighbor sampling — sampling at the pixel center keeps symmetric
// glyphs (0, 8, 1) symmetric after resizing, unlike floor-based sampling.
func srcIndex(dst, dstN, srcN int) int {
	s := iround((float64(dst)+0.5)*float64(srcN)/float64(dstN) - 0.5)
	if s < 0 {
		s = 0
	}
	if s >= srcN {
		s = srcN - 1
	}
	return s
}

// drawNumbers plots the hour numbers 1–12 as braille dots just inside the ticks,
// sized by numberSize so they grow with the dial.
func drawNumbers(c *Canvas, cx, cy, rx, ry int) {
	tw, th := numberSize(rx, ry)
	lx, ly := iround(float64(rx)*labelFrac), iround(float64(ry)*labelFrac)
	for h := 1; h <= 12; h++ {
		x, y := ellipsePoint(cx, cy, lx, ly, float64(h%12)*30)
		drawDigits(c, x, y, strconv.Itoa(h), tw, th)
	}
}

// drawDigits plots s centered at (cx,cy), each glyph resized to tw×th dots via
// centered nearest-neighbor sampling of the 3×5 dot-matrix font.
func drawDigits(c *Canvas, cx, cy int, s string, tw, th int) {
	if tw < 1 {
		tw = digitW
	}
	if th < 1 {
		th = digitH
	}
	gap := glyphGap(tw)
	total := len(s)*tw + (len(s)-1)*gap
	x0 := cx - total/2
	y0 := cy - th/2
	for _, ch := range s {
		if glyph, ok := digitFont[ch]; ok {
			for dy := 0; dy < th; dy++ {
				row := glyph[srcIndex(dy, th, digitH)]
				for dx := 0; dx < tw; dx++ {
					if sc := srcIndex(dx, tw, digitW); sc < len(row) && row[sc] == '#' {
						c.Set(x0+dx, y0+dy)
					}
				}
			}
		}
		x0 += tw + gap
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
