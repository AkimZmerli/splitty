package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/AkimZmerli/splitty"
)

func main() {
	m := splitty.New()

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
