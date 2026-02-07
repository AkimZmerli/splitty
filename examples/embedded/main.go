package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	case tea.KeyMsg:
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

func (a *app) View() string {
	sidebarWidth := 0
	if a.showSidebar {
		sidebarWidth = 25
	}

	splitsView := a.splits.View()

	if sidebarWidth > 0 {
		sidebar := lipgloss.NewStyle().
			Width(sidebarWidth - 1).
			Height(a.height).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#88C0D0")).
			Padding(1).
			Render("Sidebar\n\nCtrl+T toggle\nCtrl+Q quit")

		return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, splitsView)
	}

	return splitsView
}

func main() {
	p := tea.NewProgram(newApp(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
