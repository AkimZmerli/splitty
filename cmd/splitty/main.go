package main

import (
	"fmt"
	"os"

	"github.com/AkimZmerli/splitty"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := splitty.New(
		splitty.WithTheme(splitty.TokyoNight),
		splitty.WithScrollbackLines(1000),
		splitty.WithMouse(true),
		splitty.WithStatusBar(true),
	)

	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
