// Package config defines the JSON configuration schema and its defaults.
package config

import "github.com/michi-1221/tty-clock/internal/theme"

// Config is the on-disk schema. Strings (not enums) so JSON stays friendly;
// they are parsed/validated via the clock package.
type Config struct {
	Mode        string         `json:"mode"`        // "digital" | "analog"
	Theme       string         `json:"theme"`       // preset name
	CustomTheme *theme.Palette `json:"customTheme"` // optional inline override (nil ok)
	Granularity string         `json:"granularity"` // "seconds" | "minutes"
	CellAspect  float64        `json:"cellAspect"`  // analog dial: cell height/width; 0 = auto-detect (fallback 2.0)
	ShowHelp    bool           `json:"showHelp"`    // help line visible at startup; default true
	Format      FormatOptions  `json:"format"`
}

// FormatOptions controls what the digital clock shows.
type FormatOptions struct {
	Hour24      bool   `json:"hour24"`      // default true
	ShowSeconds bool   `json:"showSeconds"` // default true
	ShowDate    bool   `json:"showDate"`    // default true
	DateFormat  string `json:"dateFormat"`  // Go layout; default "Mon 2006-01-02"
	BlinkColon  bool   `json:"blinkColon"`  // default false
	Font        string `json:"font"`        // "block" | "ascii" | "seg7"; default "block"
	ShowNumbers bool   `json:"showNumbers"` // analog: draw 1–12 on the face; default true
}

// DefaultConfig returns the baseline a loaded file is overlaid onto, so any
// omitted field keeps its default.
func DefaultConfig() Config {
	return Config{
		Mode:        "digital",
		Theme:       theme.DefaultPreset,
		Granularity: "seconds",
		ShowHelp:    true,
		Format: FormatOptions{
			Hour24:      true,
			ShowSeconds: true,
			ShowDate:    true,
			DateFormat:  "Mon 2006-01-02",
			BlinkColon:  false,
			Font:        "block",
			ShowNumbers: true,
		},
	}
}
