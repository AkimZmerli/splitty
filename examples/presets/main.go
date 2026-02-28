package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/AkimZmerli/splitty"
)

func main() {
	// Start with a developer layout: 60/40 vertical split
	// with the right side split horizontally
	m := splitty.New(
		splitty.WithPreset(splitty.PresetDev),
		splitty.WithTheme(splitty.Sorbet),
	)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
