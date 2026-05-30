// Package theme defines color palettes and resolved themes for the clock.
//
// Roles: Primary=digits/hands, Accent=seconds/colon/analog hands, Secondary=date,
// Muted=clock face/ticks, Background=optional surface.
package theme

import "github.com/charmbracelet/lipgloss"

// Palette is the raw-hex form used in JSON config (customTheme override).
type Palette struct {
	Primary    string        `json:"primary,omitempty"`
	Accent     string        `json:"accent,omitempty"`
	Secondary  string        `json:"secondary,omitempty"`
	Muted      string        `json:"muted,omitempty"`
	Background string        `json:"background,omitempty"`
	Gradient   *GradientSpec `json:"gradient,omitempty"`
}

// GradientSpec is an optional multi-stop gradient painted across the giant
// digits (and the analog dial). With fewer than two stops it is ignored and the
// theme renders flat with the Primary/Accent roles.
type GradientSpec struct {
	Stops     []string `json:"stops"`               // hex colors, first → last along the axis
	Direction string   `json:"direction,omitempty"` // "vertical" (default) | "horizontal"
}

// Direction is the axis a gradient is painted along.
type Direction int

const (
	Vertical   Direction = iota // top → bottom (default)
	Horizontal                  // left → right
)

// Theme is a resolved palette with lipgloss colors ready for styling. The
// color profile is applied later by the renderer that calls .Render(), so
// these stay profile-independent.
type Theme struct {
	Name          string
	Primary       lipgloss.Color
	Accent        lipgloss.Color
	Secondary     lipgloss.Color
	Muted         lipgloss.Color
	Background    lipgloss.Color
	HasBackground bool      // false by default: respect the terminal background
	GradDir       Direction // gradient axis when gradient is set
	gradient      []rgb     // resolved gradient stops; len < 2 → render flat
}

// HasGradient reports whether this theme paints a multi-stop gradient (≥2 stops)
// rather than a single flat color.
func (t Theme) HasGradient() bool { return len(t.gradient) >= 2 }
