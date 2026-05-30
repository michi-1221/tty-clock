package theme

// DefaultPreset is the theme used when config omits one (採用: Tokyo Night).
const DefaultPreset = "tokyo-night"

type preset struct {
	name string
	pal  Palette
}

// presets defines the bundled themes in cycle order; the first is the default.
// Fields are positional: {Primary, Accent, Secondary, Muted, Background}.
var presets = []preset{
	{"tokyo-night", Palette{"#c0caf5", "#7aa2f7", "#bb9af7", "#565f89", "#1a1b26"}},
	{"dracula", Palette{"#f8f8f2", "#ff79c6", "#bd93f9", "#6272a4", "#282a36"}},
	{"nord", Palette{"#d8dee9", "#88c0d0", "#81a1c1", "#4c566a", "#2e3440"}},
	{"gruvbox", Palette{"#ebdbb2", "#fe8019", "#fabd2f", "#928374", "#282828"}},
	{"catppuccin-mocha", Palette{"#cdd6f4", "#cba6f7", "#fab387", "#6c7086", "#1e1e2e"}},
	{"solarized-dark", Palette{"#839496", "#b58900", "#2aa198", "#586e75", "#002b36"}},
	{"monochrome", Palette{"#ffffff", "#ffffff", "#bbbbbb", "#666666", "#000000"}},
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
