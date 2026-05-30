package analog

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/michi-1221/tty-clock/internal/caps"
	"github.com/michi-1221/tty-clock/internal/clock"
	"github.com/michi-1221/tty-clock/internal/config"
	"github.com/michi-1221/tty-clock/internal/render"
	"github.com/michi-1221/tty-clock/internal/theme"
)

var update = flag.Bool("update", false, "update golden files")

func capsProfile(p termenv.Profile) caps.Capabilities {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(p)
	return caps.New(r, true)
}

func mustTheme(t *testing.T) theme.Theme {
	t.Helper()
	th, err := theme.Resolve("tokyo-night", nil)
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// ctxAt builds a render context at a fixed, hand-distinct time (10:10:30).
func ctxAt(t *testing.T, w, h int, prof termenv.Profile) render.RenderContext {
	return render.RenderContext{
		Now:    clock.NewSnapshot(time.Date(2026, 1, 2, 10, 10, 30, 0, time.UTC)),
		Theme:  mustTheme(t),
		Format: config.FormatOptions{ShowSeconds: true, ShowNumbers: true},
		Gran:   clock.GranSeconds,
		Width:  w, Height: h,
		CellAspect: 2.0, // square-dot baseline for deterministic goldens
		Caps:       capsProfile(prof),
	}
}

func checkGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", "golden", name+".txt")
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run: go test ./internal/render/analog -update)", path, err)
	}
	if got != string(want) {
		t.Errorf("golden %s mismatch:\n--- got ---\n%s\n--- want ---\n%s", name, got, string(want))
	}
}

func TestAnalogGolden(t *testing.T) {
	checkGolden(t, "analog_40x20", New().Render(ctxAt(t, 40, 20, termenv.Ascii)))
	checkGolden(t, "analog_30x14", New().Render(ctxAt(t, 30, 14, termenv.Ascii)))
}

func TestAnalogTrueColorHasESC(t *testing.T) {
	if out := New().Render(ctxAt(t, 40, 20, termenv.TrueColor)); !containsESC(out) {
		t.Error("TrueColor render should contain ANSI color escapes")
	}
}

func TestAnalogTinyNoPanic(t *testing.T) {
	if out := New().Render(ctxAt(t, 2, 1, termenv.Ascii)); out == "" { // below MinSize
		t.Error("tiny render should still produce output (clamped, no panic)")
	}
}

func TestSecondHandGated(t *testing.T) {
	on := ctxAt(t, 40, 20, termenv.Ascii) // ShowSeconds true
	off := on
	off.Format.ShowSeconds = false
	if litCells(New().Render(on)) <= litCells(New().Render(off)) {
		t.Error("'s' (ShowSeconds off) should drop the analog second hand → fewer lit cells")
	}
}

// litCells counts non-blank braille cells (a proxy for how much is drawn).
func litCells(s string) int {
	n := 0
	for _, r := range s {
		if r != '\n' && r != rune(0x2800) {
			n++
		}
	}
	return n
}

func containsESC(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b {
			return true
		}
	}
	return false
}

var fgRe = regexp.MustCompile(`38;2;\d+;\d+;\d+`)

func TestAnalogGradientPaintsMultipleColors(t *testing.T) {
	sunset, err := theme.Resolve("sunset", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := ctxAt(t, 40, 20, termenv.TrueColor)
	ctx.Theme = sunset
	seen := map[string]struct{}{}
	for _, m := range fgRe.FindAllString(New().Render(ctx), -1) {
		seen[m] = struct{}{}
	}
	if len(seen) < 2 {
		t.Errorf("gradient dial used %d distinct colors, want ≥2", len(seen))
	}
}
