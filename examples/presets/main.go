package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/AkimZmerli/splitty"
)

func main() {
	// Start with a developer layout: 60/40 vertical split
	// with the right side split horizontally
	m := splitty.New(
		splitty.WithPreset(splitty.PresetDev),
		splitty.WithTheme(splitty.Catppuccin),
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
