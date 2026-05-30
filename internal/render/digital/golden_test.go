package digital

import (
	"flag"
	"io"
	"os"
	"path/filepath"
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

// asciiCaps returns capabilities with a forced Ascii profile and ascii font, so
// Render output is deterministic plain text (★3) — no global renderer, no color.
func asciiCaps() caps.Capabilities {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.Ascii)
	return caps.New(r, false) // unicode=false → ascii font
}

func mustTheme(t *testing.T) theme.Theme {
	t.Helper()
	th, err := theme.Resolve("tokyo-night", nil)
	if err != nil {
		t.Fatal(err)
	}
	return th
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
		t.Fatalf("read golden %s: %v (run: go test ./... -update)", path, err)
	}
	if got != string(want) {
		t.Errorf("golden %s mismatch:\n--- got ---\n%q\n--- want ---\n%q", name, got, string(want))
	}
}

func TestDigitalGolden(t *testing.T) {
	fixed := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC) // 15:04:05, odd second
	th := mustTheme(t)
	d := New()

	cases := []struct {
		name string
		fmt  config.FormatOptions
		gran clock.Granularity
	}{
		{
			name: "digital_24h_seconds",
			fmt:  config.FormatOptions{Hour24: true, ShowSeconds: true, ShowDate: true, DateFormat: "Mon 2006-01-02", Font: "block"},
			gran: clock.GranSeconds,
		},
		{
			name: "digital_12h_ampm_nosecs",
			fmt:  config.FormatOptions{Hour24: false, ShowSeconds: false, ShowDate: false, ShowAMPM: true, Font: "block"},
			gran: clock.GranSeconds,
		},
		{
			name: "digital_blink_off",
			fmt:  config.FormatOptions{Hour24: true, ShowSeconds: true, ShowDate: false, BlinkColon: true, Font: "block"},
			gran: clock.GranSeconds,
		},
		{
			name: "digital_minutes_hides_seconds",
			fmt:  config.FormatOptions{Hour24: true, ShowSeconds: true, ShowDate: false, Font: "block"},
			gran: clock.GranMinutes,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := render.RenderContext{
				Now:    clock.NewSnapshot(fixed),
				Theme:  th,
				Format: tc.fmt,
				Gran:   tc.gran,
				Width:  80,
				Height: 24,
				Caps:   asciiCaps(),
			}
			checkGolden(t, tc.name, d.Render(ctx))
		})
	}
}

func TestDigitalGoldenScaled(t *testing.T) {
	ctx := render.RenderContext{
		Now:    clock.NewSnapshot(time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)),
		Theme:  mustTheme(t),
		Format: config.FormatOptions{Hour24: true, ShowSeconds: true, ShowDate: true, DateFormat: "Mon 2006-01-02", Font: "block"},
		Gran:   clock.GranSeconds,
		Width:  120, Height: 30, Scale: 2,
		Caps: asciiCaps(),
	}
	checkGolden(t, "digital_scale2", New().Render(ctx))
}

func TestScaleClampedToFit(t *testing.T) {
	ctx := render.RenderContext{
		Now:    clock.NewSnapshot(time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)),
		Theme:  mustTheme(t),
		Format: config.FormatOptions{Hour24: true, ShowSeconds: true, Font: "block"},
		Gran:   clock.GranSeconds,
		Width:  60, Height: 14, Scale: 9, // oversized request must be clamped to fit
		Caps: asciiCaps(),
	}
	out := New().Render(ctx)
	if w := lipgloss.Width(out); w > ctx.Width {
		t.Errorf("scaled width %d exceeds available %d", w, ctx.Width)
	}
	if h := lipgloss.Height(out); h > ctx.Height {
		t.Errorf("scaled height %d exceeds available %d", h, ctx.Height)
	}
}

func TestDigitalRendersColorInTrueColor(t *testing.T) {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.TrueColor)
	ctx := render.RenderContext{
		Now:    clock.NewSnapshot(time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)),
		Theme:  mustTheme(t),
		Format: config.FormatOptions{Hour24: true, ShowSeconds: true, Font: "block"},
		Gran:   clock.GranSeconds,
		Width:  80, Height: 24,
		Caps: caps.New(r, true),
	}
	out := New().Render(ctx)
	if !containsESC(out) {
		t.Error("TrueColor render should contain ANSI color escapes")
	}
}

func containsESC(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b {
			return true
		}
	}
	return false
}
