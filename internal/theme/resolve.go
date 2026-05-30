package theme

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var hexRe = regexp.MustCompile(`^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$`)

// Resolve builds a Theme from a preset name plus an optional inline override.
// Only non-empty override fields replace preset values (partial override is
// allowed). An unknown name is legal only when the override fully specifies
// every required role; otherwise it is an error.
func Resolve(name string, override *Palette) (Theme, error) {
	base, known := presetByName(name)
	if !known {
		if override == nil {
			return Theme{}, fmt.Errorf("unknown theme %q (known: %s)", name, strings.Join(Names(), ", "))
		}
		base = Palette{} // override must fill everything
	}
	merged := merge(base, override)

	for _, c := range []struct{ role, val string }{
		{"primary", merged.Primary}, {"accent", merged.Accent},
		{"secondary", merged.Secondary}, {"muted", merged.Muted},
		{"background", merged.Background},
	} {
		if c.val != "" && !hexRe.MatchString(c.val) {
			return Theme{}, fmt.Errorf("theme %q: invalid %s color %q (want #rgb or #rrggbb)", name, c.role, c.val)
		}
	}
	if err := requireComplete(merged); err != nil {
		return Theme{}, fmt.Errorf("theme %q: %w", name, err)
	}

	return Theme{
		Name:          name,
		Primary:       lipgloss.Color(merged.Primary),
		Accent:        lipgloss.Color(merged.Accent),
		Secondary:     lipgloss.Color(merged.Secondary),
		Muted:         lipgloss.Color(merged.Muted),
		Background:    lipgloss.Color(merged.Background),
		HasBackground: false,
	}, nil
}

func merge(base Palette, o *Palette) Palette {
	if o == nil {
		return base
	}
	if o.Primary != "" {
		base.Primary = o.Primary
	}
	if o.Accent != "" {
		base.Accent = o.Accent
	}
	if o.Secondary != "" {
		base.Secondary = o.Secondary
	}
	if o.Muted != "" {
		base.Muted = o.Muted
	}
	if o.Background != "" {
		base.Background = o.Background
	}
	return base
}

// requireComplete checks the foreground roles are present. Background is
// optional (an empty background means "use the terminal background").
func requireComplete(p Palette) error {
	var missing []string
	if p.Primary == "" {
		missing = append(missing, "primary")
	}
	if p.Accent == "" {
		missing = append(missing, "accent")
	}
	if p.Secondary == "" {
		missing = append(missing, "secondary")
	}
	if p.Muted == "" {
		missing = append(missing, "muted")
	}
	if len(missing) > 0 {
		return fmt.Errorf("incomplete palette, missing: %s", strings.Join(missing, ", "))
	}
	return nil
}
