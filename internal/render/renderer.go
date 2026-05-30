// Package render defines the framework-independent rendering seam. A Renderer
// turns a RenderContext into a string and never imports bubbletea, so output
// can be golden-tested without starting a program. ui is the only package that
// imports bubbletea.
package render

import (
	"github.com/michi-1221/tty-clock/internal/caps"
	"github.com/michi-1221/tty-clock/internal/clock"
	"github.com/michi-1221/tty-clock/internal/config"
	"github.com/michi-1221/tty-clock/internal/theme"
)

// RenderContext is everything a Renderer needs. Same input → same output.
type RenderContext struct {
	Now        clock.TimeSnapshot
	Theme      theme.Theme
	Format     config.FormatOptions
	Gran       clock.Granularity
	Width      int // available area after ui subtracts the footer
	Height     int
	Scale      int     // integer font magnification chosen by ui; 0/unset means ×1
	CellAspect float64 // resolved cell height/width for the analog dial; always > 0
	Caps       caps.Capabilities
}

// Renderer draws one clock mode (digital/analog).
type Renderer interface {
	Render(ctx RenderContext) string      // pure; must not panic on tiny sizes
	MinSize(ctx RenderContext) (w, h int) // smallest drawable cell box
	Name() string                         // "digital" / "analog"
}
