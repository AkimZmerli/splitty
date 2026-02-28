package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/AkimZmerli/splitty"
)

// app demonstrates embedding splitty.Manager inside a larger Bubble Tea application.
type app struct {
	splits      *splitty.Manager
	showSidebar bool
	width       int
	height      int
}

func newApp() *app {
	return &app{
		splits: splitty.New(
			splitty.WithTheme(splitty.Glacier),
			splitty.WithStatusBar(false), // we render our own chrome
		),
	}
}

func (a *app) Init() tea.Cmd {
	return a.splits.Init()
}

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+t" {
			a.showSidebar = !a.showSidebar
			return a, nil
		}
		if msg.String() == "ctrl+q" {
			return a, tea.Quit
		}
	}

	updated, cmd := a.splits.Update(msg)
	a.splits = updated.(*splitty.Manager)
	return a, cmd
}

func (a *app) View() tea.View {
	sidebarWidth := 0
	if a.showSidebar {
		sidebarWidth = 25
	}

	splitsContent := a.splits.View().Content

	var content string
	if sidebarWidth > 0 {
		sidebar := lipgloss.NewStyle().
			Width(sidebarWidth - 1).
			Height(a.height).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#88C0D0")).
			Padding(1).
			Render("Sidebar\n\nCtrl+T toggle\nCtrl+Q quit")

		content = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, splitsContent)
	} else {
		content = splitsContent
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func main() {
	p := tea.NewProgram(newApp())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
