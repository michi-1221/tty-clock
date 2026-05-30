package theme

import (
	"fmt"
	"math"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// rgb is a parsed 8-bit-per-channel color used for gradient interpolation.
type rgb struct{ r, g, b uint8 }

// color renders the rgb back to a #rrggbb lipgloss color.
func (c rgb) color() lipgloss.Color {
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", c.r, c.g, c.b))
}

// ColorAt returns the gradient color at frac ∈ [0,1] (clamped), interpolating
// linearly between the surrounding stops. With no gradient it returns Primary;
// the render painter calls this per row (vertical) or per cell (horizontal).
func (t Theme) ColorAt(frac float64) lipgloss.Color {
	n := len(t.gradient)
	switch {
	case n == 0:
		return t.Primary
	case n == 1 || frac <= 0:
		return t.gradient[0].color()
	case frac >= 1:
		return t.gradient[n-1].color()
	}
	pos := frac * float64(n-1)
	i := int(pos)
	return lerp(t.gradient[i], t.gradient[i+1], pos-float64(i)).color()
}

func lerp(a, b rgb, t float64) rgb {
	return rgb{lerp8(a.r, b.r, t), lerp8(a.g, b.g, t), lerp8(a.b, b.b, t)}
}

func lerp8(a, b uint8, t float64) uint8 {
	return uint8(math.Round(float64(a) + (float64(b)-float64(a))*t))
}

// parseHex parses a #rgb or #rrggbb string (validated by hexRe) into an rgb.
func parseHex(s string) (rgb, error) {
	if !hexRe.MatchString(s) {
		return rgb{}, fmt.Errorf("invalid color %q (want #rgb or #rrggbb)", s)
	}
	h := s[1:]
	if len(h) == 3 { // expand shorthand #rgb → #rrggbb
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return rgb{}, err
	}
	return rgb{uint8(v >> 16), uint8(v >> 8), uint8(v)}, nil
}

// parseDirection maps the JSON direction string to a Direction (default vertical).
func parseDirection(s string) (Direction, error) {
	switch s {
	case "", "vertical":
		return Vertical, nil
	case "horizontal":
		return Horizontal, nil
	default:
		return Vertical, fmt.Errorf("invalid direction %q (want vertical or horizontal)", s)
	}
}
