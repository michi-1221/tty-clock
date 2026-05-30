package theme

import "testing"

func TestPresetsComplete(t *testing.T) {
	for _, name := range Names() {
		th, err := Resolve(name, nil)
		if err != nil {
			t.Errorf("Resolve(%q) error: %v", name, err)
			continue
		}
		if th.Primary == "" || th.Accent == "" || th.Secondary == "" || th.Muted == "" {
			t.Errorf("preset %q has an empty role: %+v", name, th)
		}
	}
}

func TestResolveUnknownName(t *testing.T) {
	if _, err := Resolve("does-not-exist", nil); err == nil {
		t.Error("unknown name without override should error")
	}
	full := &Palette{Primary: "#111111", Accent: "#222222", Secondary: "#333333", Muted: "#444444"}
	th, err := Resolve("does-not-exist", full)
	if err != nil {
		t.Fatalf("unknown name with full override should resolve: %v", err)
	}
	if string(th.Primary) != "#111111" {
		t.Errorf("Primary = %q, want #111111", th.Primary)
	}
}

func TestResolvePartialOverride(t *testing.T) {
	th, err := Resolve("tokyo-night", &Palette{Accent: "#ff0000"})
	if err != nil {
		t.Fatal(err)
	}
	if string(th.Accent) != "#ff0000" {
		t.Errorf("Accent override = %q, want #ff0000", th.Accent)
	}
	if string(th.Primary) != "#c0caf5" {
		t.Errorf("Primary = %q, want preset #c0caf5", th.Primary)
	}
}

func TestResolveBadHex(t *testing.T) {
	if _, err := Resolve("tokyo-night", &Palette{Accent: "red"}); err == nil {
		t.Error("invalid hex override should error")
	}
}

func TestNextCycle(t *testing.T) {
	if got := Next("tokyo-night"); got != "dracula" {
		t.Errorf("Next(tokyo-night) = %q, want dracula", got)
	}
	if got := Next("monochrome"); got != "tokyo-night" {
		t.Errorf("Next(monochrome) wrap = %q, want tokyo-night", got)
	}
	if got := Next("custom-unknown"); got != "tokyo-night" {
		t.Errorf("Next(unknown) = %q, want tokyo-night", got)
	}
	// Walking the whole ring returns to the start.
	name := Names()[0]
	for range Names() {
		name = Next(name)
	}
	if name != Names()[0] {
		t.Errorf("full cycle ended at %q, want %q", name, Names()[0])
	}
}
