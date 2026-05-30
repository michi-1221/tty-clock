package ui

import "github.com/charmbracelet/bubbles/key"

// keyMap is the single source of truth for keys and their help labels.
type keyMap struct {
	Quit       key.Binding
	ToggleMode key.Binding // m: digital ⇄ analog
	CycleTheme key.Binding
	ToggleSecs key.Binding
	Reload     key.Binding // r: re-read the config file
	Edit       key.Binding // e: open the config file in $EDITOR
	Help       key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"), key.WithHelp("q", "quit")),
		ToggleMode: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "mode")),
		CycleTheme: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "theme")),
		ToggleSecs: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "seconds")),
		Reload:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload")),
		Edit:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit config")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

// ShortHelp / FullHelp implement help.KeyMap — the footer is generated from
// these, so every working key must appear here.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleSecs, k.CycleTheme, k.ToggleMode, k.Reload, k.Edit, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleSecs, k.CycleTheme, k.ToggleMode},
		{k.Reload, k.Edit, k.Help, k.Quit},
	}
}
