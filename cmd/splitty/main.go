package main

import (
	"fmt"
	"os"

	"github.com/AkimZmerli/splitty"
	tea "charm.land/bubbletea/v2"
)

func main() {
	m := splitty.New(
		splitty.WithTheme(splitty.TokyoNight),
		splitty.WithScrollbackLines(1000),
		splitty.WithMouse(true),
		splitty.WithStatusBar(true),
	)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
