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
	ScrollUp        key.Binding
	ScrollDown      key.Binding
	ScrollPageUp    key.Binding
	ScrollPageDown  key.Binding
	ScrollToTop     key.Binding
	ScrollToBottom  key.Binding
	CycleTheme      key.Binding
	EnterCopyMode   key.Binding
}

// DefaultKeyMap returns the default keybindings for split pane operations.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		SplitVertical: key.NewBinding(
			key.WithKeys("ctrl+v"),
			key.WithHelp("ctrl+v", "split vertical"),
		),
		SplitHorizontal: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "split horizontal"),
		),
		Close: key.NewBinding(
			key.WithKeys("ctrl+q"),
			key.WithHelp("ctrl+q", "close pane"),
		),
		FocusLeft: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "focus left"),
		),
		FocusDown: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "focus down"),
		),
		FocusUp: key.NewBinding(
			key.WithKeys("ctrl+w"),
			key.WithHelp("ctrl+w", "focus up"),
		),
		FocusRight: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "focus right"),
		),
		FocusCycle: key.NewBinding(
			key.WithKeys("ctrl+tab"),
			key.WithHelp("ctrl+tab", "next pane"),
		),
		FocusCycleBack: key.NewBinding(
			key.WithKeys("ctrl+shift+tab"),
			key.WithHelp("ctrl+shift+tab", "prev pane"),
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
		ScrollUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "scroll down"),
		),
		ScrollPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "page up"),
		),
		ScrollPageDown: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "page down"),
		),
		ScrollToTop: key.NewBinding(
			key.WithKeys("ctrl+home"),
			key.WithHelp("ctrl+home", "scroll to top"),
		),
		ScrollToBottom: key.NewBinding(
			key.WithKeys("ctrl+end"),
			key.WithHelp("ctrl+end", "scroll to bottom"),
		),
		CycleTheme: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "cycle theme"),
		),
		EnterCopyMode: key.NewBinding(
			key.WithKeys("ctrl+["),
			key.WithHelp("ctrl+[", "copy mode"),
		),
	}
}
