package splitty

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for split pane operations.
type KeyMap struct {
	SplitVertical   key.Binding
	SplitHorizontal key.Binding
	Close           key.Binding
	FocusLeft       key.Binding
	FocusDown       key.Binding
	FocusUp         key.Binding
	FocusRight      key.Binding
	FocusCycle      key.Binding
	FocusCycleBack  key.Binding
	Zoom            key.Binding
	Swap            key.Binding
	ResizeLeft      key.Binding
	ResizeRight     key.Binding
	ResizeUp        key.Binding
	ResizeDown      key.Binding
	Broadcast       key.Binding
}

// DefaultKeyMap returns the default keybindings for split pane operations.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		SplitVertical: key.NewBinding(
			key.WithKeys("ctrl+\\"),
			key.WithHelp("ctrl+\\", "split vertical"),
		),
		SplitHorizontal: key.NewBinding(
			key.WithKeys("ctrl+-"),
			key.WithHelp("ctrl+-", "split horizontal"),
		),
		Close: key.NewBinding(
			key.WithKeys("ctrl+w"),
			key.WithHelp("ctrl+w", "close pane"),
		),
		FocusLeft: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "focus left"),
		),
		FocusDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "focus down"),
		),
		FocusUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "focus up"),
		),
		FocusRight: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "focus right"),
		),
		FocusCycle: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next pane"),
		),
		FocusCycleBack: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev pane"),
		),
		Zoom: key.NewBinding(
			key.WithKeys("ctrl+z"),
			key.WithHelp("ctrl+z", "toggle zoom"),
		),
		Swap: key.NewBinding(
			key.WithKeys("ctrl+x"),
			key.WithHelp("ctrl+x", "swap panes"),
		),
		ResizeLeft: key.NewBinding(
			key.WithKeys("ctrl+shift+left"),
			key.WithHelp("ctrl+shift+left", "resize left"),
		),
		ResizeRight: key.NewBinding(
			key.WithKeys("ctrl+shift+right"),
			key.WithHelp("ctrl+shift+right", "resize right"),
		),
		ResizeUp: key.NewBinding(
			key.WithKeys("ctrl+shift+up"),
			key.WithHelp("ctrl+shift+up", "resize up"),
		),
		ResizeDown: key.NewBinding(
			key.WithKeys("ctrl+shift+down"),
			key.WithHelp("ctrl+shift+down", "resize down"),
		),
		Broadcast: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "toggle broadcast"),
		),
	}
}
