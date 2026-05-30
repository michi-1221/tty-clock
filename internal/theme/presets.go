package theme

// DefaultPreset is the theme used when config omits one (採用: Tokyo Night).
const DefaultPreset = "tokyo-night"

type preset struct {
	name string
	pal  Palette
}

// presets defines the bundled themes in cycle order; the first is the default.
// Flat themes set the five color roles; gradient themes add a Gradient (the
// flat roles stay as a fallback for non-gradient surfaces, e.g. the date).
// monochrome is kept last so the 't' cycle wraps cleanly back to the default.
var presets = []preset{
	{"tokyo-night", Palette{Primary: "#c0caf5", Accent: "#7aa2f7", Secondary: "#bb9af7", Muted: "#565f89", Background: "#1a1b26"}},
	{"dracula", Palette{Primary: "#f8f8f2", Accent: "#ff79c6", Secondary: "#bd93f9", Muted: "#6272a4", Background: "#282a36"}},
	{"nord", Palette{Primary: "#d8dee9", Accent: "#88c0d0", Secondary: "#81a1c1", Muted: "#4c566a", Background: "#2e3440"}},
	{"gruvbox", Palette{Primary: "#ebdbb2", Accent: "#fe8019", Secondary: "#fabd2f", Muted: "#928374", Background: "#282828"}},
	{"catppuccin-mocha", Palette{Primary: "#cdd6f4", Accent: "#cba6f7", Secondary: "#fab387", Muted: "#6c7086", Background: "#1e1e2e"}},
	{"solarized-dark", Palette{Primary: "#839496", Accent: "#b58900", Secondary: "#2aa198", Muted: "#586e75", Background: "#002b36"}},
	{"sunset", Palette{
		Primary: "#ff7e5f", Accent: "#ffd26f", Secondary: "#ffb199", Muted: "#7a5c8e", Background: "#1a1326",
		Gradient: &GradientSpec{Stops: []string{"#ffe27a", "#ff9e5e", "#ff5e62", "#a23bc0"}},
	}},
	{"aurora", Palette{
		Primary: "#00d9b8", Accent: "#5efce8", Secondary: "#a0e9ff", Muted: "#3a5a6a", Background: "#0a1622",
		Gradient: &GradientSpec{Stops: []string{"#00f5a0", "#00d9f5", "#8e54e9"}},
	}},
	{"rainbow", Palette{
		Primary: "#ffca3a", Accent: "#ff595e", Secondary: "#8ac926", Muted: "#555555", Background: "#101014",
		Gradient: &GradientSpec{Stops: []string{"#ff595e", "#ff924c", "#ffca3a", "#8ac926", "#1982c4", "#6a4c93"}, Direction: "horizontal"},
	}},
	{"monochrome", Palette{Primary: "#ffffff", Accent: "#ffffff", Secondary: "#bbbbbb", Muted: "#666666", Background: "#000000"}},
}

// presetByName returns the raw palette for a preset name.
func presetByName(name string) (Palette, bool) {
	for _, p := range presets {
		if p.name == name {
			return p.pal, true
		}
	}
	return Palette{}, false
}

// Names returns every preset name in cycle order.
func Names() []string {
	out := make([]string, len(presets))
	for i, p := range presets {
		out[i] = p.name
	}
	return out
}

// Next returns the preset name after the given one, wrapping around. Names not
// in the preset list (e.g. a custom-only theme) start the cycle at the first
// preset — so the 't' key always walks clean presets (採用 #7).
func Next(name string) string {
	for i, p := range presets {
		if p.name == name {
			return presets[(i+1)%len(presets)].name
		}
	}
	return presets[0].name
}
