// Package theme defines color palettes and resolved themes for the clock.
//
// Roles: Primary=digits/hands, Accent=seconds/colon/AM-PM, Secondary=date,
// Muted=clock face/ticks, Background=optional surface.
package theme

import "github.com/charmbracelet/lipgloss"

// Palette is the raw-hex form used in JSON config (customTheme override).
type Palette struct {
	Primary    string `json:"primary,omitempty"`
	Accent     string `json:"accent,omitempty"`
	Secondary  string `json:"secondary,omitempty"`
	Muted      string `json:"muted,omitempty"`
	Background string `json:"background,omitempty"`
}

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
	HasBackground bool // false by default: respect the terminal background
}
