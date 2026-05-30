package theme

import "testing"

func TestResolveGradientPreset(t *testing.T) {
	th, err := Resolve("sunset", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !th.HasGradient() {
		t.Fatal("sunset should have a gradient")
	}
	if th.GradDir != Vertical {
		t.Errorf("sunset direction = %v, want vertical", th.GradDir)
	}
	// Endpoints match the first and last stop; the midpoint differs from both.
	if got := string(th.ColorAt(0)); got != "#ffe27a" {
		t.Errorf("ColorAt(0) = %q, want #ffe27a", got)
	}
	if got := string(th.ColorAt(1)); got != "#a23bc0" {
		t.Errorf("ColorAt(1) = %q, want #a23bc0", got)
	}
	mid := string(th.ColorAt(0.5))
	if mid == "#ffe27a" || mid == "#a23bc0" {
		t.Errorf("ColorAt(0.5) = %q, want an interpolated color", mid)
	}
}

func TestRainbowIsHorizontal(t *testing.T) {
	th, err := Resolve("rainbow", nil)
	if err != nil {
		t.Fatal(err)
	}
	if th.GradDir != Horizontal {
		t.Errorf("rainbow direction = %v, want horizontal", th.GradDir)
	}
}

func TestColorAtMidpointInterpolates(t *testing.T) {
	th, err := Resolve("x", &Palette{
		Primary: "#808080", Accent: "#808080", Secondary: "#808080", Muted: "#808080",
		Gradient: &GradientSpec{Stops: []string{"#000000", "#ffffff"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := string(th.ColorAt(0.5)); got != "#808080" {
		t.Errorf("midpoint of black→white = %q, want #808080", got)
	}
}

func TestFlatThemeHasNoGradient(t *testing.T) {
	th, err := Resolve("tokyo-night", nil)
	if err != nil {
		t.Fatal(err)
	}
	if th.HasGradient() {
		t.Error("tokyo-night should be flat")
	}
	if got := th.ColorAt(0.5); got != th.Primary {
		t.Errorf("flat ColorAt = %q, want Primary %q", got, th.Primary)
	}
}

func TestResolveBadGradientStop(t *testing.T) {
	_, err := Resolve("tokyo-night", &Palette{
		Gradient: &GradientSpec{Stops: []string{"#000000", "nope"}},
	})
	if err == nil {
		t.Error("invalid gradient stop should error")
	}
}

func TestResolveBadGradientDirection(t *testing.T) {
	_, err := Resolve("tokyo-night", &Palette{
		Gradient: &GradientSpec{Stops: []string{"#000000", "#ffffff"}, Direction: "sideways"},
	})
	if err == nil {
		t.Error("invalid gradient direction should error")
	}
}

func TestShorthandHexGradient(t *testing.T) {
	th, err := Resolve("x", &Palette{
		Primary: "#fff", Accent: "#fff", Secondary: "#fff", Muted: "#fff",
		Gradient: &GradientSpec{Stops: []string{"#f00", "#00f"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := string(th.ColorAt(0)); got != "#ff0000" {
		t.Errorf("shorthand #f00 → %q, want #ff0000", got)
	}
}
