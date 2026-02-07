package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/AkimZmerli/splitty"
)

func main() {
	keys := splitty.DefaultKeyMap()
	keys.SplitVertical = key.NewBinding(
		key.WithKeys("ctrl+v"),
		key.WithHelp("ctrl+v", "split vertical"),
	)
	keys.SplitHorizontal = key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "split horizontal"),
	)

	m := splitty.New(
		splitty.WithKeyMap(keys),
		splitty.WithTheme(splitty.Nightshade),
	)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
