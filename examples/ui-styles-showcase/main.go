package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ColorScheme struct {
	name           string
	primaryColor   string
	secondaryColor string
	accentColor    string
	bgColor        string
}

var colorSchemes = []ColorScheme{
	{
		name:           "Tokyo Night",
		primaryColor:   "#7AA2F7",
		secondaryColor: "#414868",
		accentColor:    "#BB9AF7",
		bgColor:        "#1A1B26",
	},
	{
		name:           "Nightshade",
		primaryColor:   "#BD93F9",
		secondaryColor: "#44475A",
		accentColor:    "#FF79C6",
		bgColor:        "#282A36",
	},
	{
		name:           "Glacier",
		primaryColor:   "#88C0D0",
		secondaryColor: "#4C566A",
		accentColor:    "#81A1C1",
		bgColor:        "#2E3440",
	},
	{
		name:           "Sorbet",
		primaryColor:   "#CBA6F7",
		secondaryColor: "#585B70",
		accentColor:    "#F5C2E7",
		bgColor:        "#1E1E2E",
	},
	{
		name:           "Solarized",
		primaryColor:   "#268BD2",
		secondaryColor: "#586E75",
		accentColor:    "#2AA198",
		bgColor:        "#002B36",
	},
}

type Model struct {
	schemeIdx     int
	width         int
	height        int
	currentScheme ColorScheme
}

func NewModel() Model {
	return Model{
		schemeIdx:     0,
		width:         80,
		height:        24,
		currentScheme: colorSchemes[0],
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "h":
			m.schemeIdx = (m.schemeIdx - 1 + len(colorSchemes)) % len(colorSchemes)
			m.currentScheme = colorSchemes[m.schemeIdx]
			return m, nil
		case "right", "l":
			m.schemeIdx = (m.schemeIdx + 1) % len(colorSchemes)
			m.currentScheme = colorSchemes[m.schemeIdx]
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func (m Model) render() string {
	scheme := m.currentScheme

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(scheme.primaryColor)).
		Padding(1, 2).
		Render(fmt.Sprintf("🎨 Color Scheme: %s", scheme.name))

	// Color preview boxes
	primaryBox := lipgloss.NewStyle().
		Width(16).
		Height(3).
		Align(lipgloss.Center, lipgloss.Center).
		Background(lipgloss.Color(scheme.primaryColor)).
		Foreground(lipgloss.Color(scheme.bgColor)).
		Bold(true).
		Render("Primary")

	secondaryBox := lipgloss.NewStyle().
		Width(16).
		Height(3).
		Align(lipgloss.Center, lipgloss.Center).
		Background(lipgloss.Color(scheme.secondaryColor)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Render("Secondary")

	accentBox := lipgloss.NewStyle().
		Width(16).
		Height(3).
		Align(lipgloss.Center, lipgloss.Center).
		Background(lipgloss.Color(scheme.accentColor)).
		Foreground(lipgloss.Color(scheme.bgColor)).
		Bold(true).
		Render("Accent")

	colorRow := lipgloss.JoinHorizontal(lipgloss.Center, primaryBox, "  ", secondaryBox, "  ", accentBox)

	// Typography samples
	h1 := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(scheme.primaryColor)).
		Render("Heading 1 - Large Title")

	h2 := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(scheme.secondaryColor)).
		Render("Heading 2 - Subtitle")

	body := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render("Body text with standard styling for regular content display")

	emphasized := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(scheme.accentColor)).
		Render("Emphasized text")

	typographyBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(scheme.primaryColor)).
		Padding(1, 2).
		Render(fmt.Sprintf("%s\n%s\n%s %s", h1, h2, body, emphasized))

	// Layout demonstration
	leftPanel := lipgloss.NewStyle().
		Width(30).
		Height(6).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(scheme.primaryColor)).
		Padding(1).
		Render("Left Panel\n\n" +
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(scheme.secondaryColor)).
				Render("Sidebar content"))

	centerPanel := lipgloss.NewStyle().
		Width(30).
		Height(6).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(scheme.accentColor)).
		Padding(1).
		Render("Center Panel\n\n" +
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(scheme.primaryColor)).
				Render("Main content"))

	rightPanel := lipgloss.NewStyle().
		Width(30).
		Height(6).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(scheme.secondaryColor)).
		Padding(1).
		Render("Right Panel\n\n" +
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(scheme.accentColor)).
				Render("Info panel"))

	layoutRow := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, centerPanel, rightPanel)

	// Status bar
	schemeNav := lipgloss.NewStyle().
		Foreground(lipgloss.Color(scheme.primaryColor)).
		Render(fmt.Sprintf("Scheme %d/%d", m.schemeIdx+1, len(colorSchemes)))

	schemeIndicators := make([]string, len(colorSchemes))
	for i := range colorSchemes {
		if i == m.schemeIdx {
			schemeIndicators[i] = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(scheme.primaryColor)).
				Render("●")
		} else {
			schemeIndicators[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color(scheme.secondaryColor)).
				Render("○")
		}
	}

	statusBar := lipgloss.NewStyle().
		Padding(0, 2).
		Render(strings.Join(schemeIndicators, " ") + " | " + schemeNav + " | ← → to switch | Q to quit")

	// Combine all sections
	content := strings.Join([]string{
		header,
		"",
		colorRow,
		"",
		typographyBox,
		"",
		layoutRow,
	}, "\n")

	// Final container
	final := lipgloss.NewStyle().
		Padding(1).
		Render(content)

	// Add footer with navigation
	footer := lipgloss.NewStyle().
		BorderTop(true).
		BorderForeground(lipgloss.Color(scheme.secondaryColor)).
		Padding(1, 2).
		Render(statusBar)

	return final + "\n" + footer
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
