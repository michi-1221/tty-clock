package ui

import "github.com/charmbracelet/bubbles/key"

// keyMap is the single source of truth for keys and their help labels.
type keyMap struct {
	Quit       key.Binding
	ToggleMode key.Binding // framed for phase 2 (analog); a no-op in phase 1 (★8)
	CycleTheme key.Binding
	ToggleSecs key.Binding
	Reload     key.Binding // framed for phase 2 (live reload)
	Help       key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"), key.WithHelp("q", "quit")),
		ToggleMode: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "mode")),
		CycleTheme: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "theme")),
		ToggleSecs: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "seconds")),
		Reload:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

// ShortHelp / FullHelp implement help.KeyMap. Phase 1 advertises only the keys
// that actually do something here; m/r are wired up in phase 2.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleSecs, k.CycleTheme, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleSecs, k.CycleTheme},
		{k.Help, k.Quit},
	}
}
