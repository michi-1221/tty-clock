// Package caps detects terminal capabilities: the lipgloss color profile (via
// an explicit renderer) and a best-effort Unicode/block-glyph guess. It owns
// an explicit *lipgloss.Renderer so styling never depends on the global
// default renderer — that keeps golden tests deterministic (★3).
package caps

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Capabilities carries the rendering decisions derived from the environment.
type Capabilities struct {
	Renderer   *lipgloss.Renderer
	Profile    termenv.Profile
	Unicode    bool    // best-effort: can the terminal draw block/braille glyphs?
	CellAspect float64 // terminal cell height/width; 0 if unknown (callers default to 2.0)
}

// Detect builds capabilities from the program's output writer.
func Detect(w io.Writer) Capabilities {
	r := lipgloss.NewRenderer(w)
	return Capabilities{
		Renderer:   r,
		Profile:    r.ColorProfile(),
		Unicode:    detectUnicode(),
		CellAspect: detectCellAspect(w),
	}
}

// detectCellAspect measures the terminal cell height/width ratio from its pixel
// size (TIOCGWINSZ) when the writer is a TTY that reports it; 0 otherwise. Used
// to draw a visually round analog dial regardless of the font's cell shape.
func detectCellAspect(w io.Writer) float64 {
	f, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return 0
	}
	return cellAspectFromFd(f.Fd())
}

// New builds capabilities around an explicit renderer (used in tests, where
// the profile is forced via r.SetColorProfile).
func New(r *lipgloss.Renderer, unicode bool) Capabilities {
	return Capabilities{Renderer: r, Profile: r.ColorProfile(), Unicode: unicode}
}

// detectUnicode is a heuristic (★7): there is no reliable way to know whether a
// terminal can draw a given glyph, so we trust TERM and the locale.
func detectUnicode() bool {
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	for _, k := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if v := os.Getenv(k); v != "" {
			u := strings.ToUpper(v)
			return strings.Contains(u, "UTF-8") || strings.Contains(u, "UTF8")
		}
	}
	return true // no locale info: assume a modern UTF-8 terminal
}
