package main

import "github.com/charmbracelet/bubbles/key"

type CustomKeyMap struct {
	SwitchPanel    key.Binding
	JumpToSelected key.Binding
	GoToParentDir  key.Binding
	LaunchDefault  key.Binding
	WriteResult    key.Binding
	Refresh        key.Binding
	Quit           key.Binding
}

func newKeyMap() *CustomKeyMap {
	return &CustomKeyMap{
		SwitchPanel: key.NewBinding(
			key.WithKeys("shift+down", "shift+up", "alt+J", "alt+K"),
			key.WithHelp("shift + ↑/↓", "Switch panel"),
		),
		JumpToSelected: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Jump to selected item"),
		),
		GoToParentDir: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "Jump to parent directory"),
		),
		LaunchDefault: key.NewBinding(
			key.WithKeys("."),
			key.WithHelp(".", "Launch default application"),
		),
		WriteResult: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "Write result to disk"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("shift + r", "Refresh catelogue"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
	}
}
