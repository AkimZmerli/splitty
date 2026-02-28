package main

import (
	"fmt"
	"io"
	"os"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Item represents a list item
type Item struct {
	title       string
	description string
}

func (i Item) FilterValue() string { return i.title }

// ItemDelegate implements list.ItemDelegate
type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 2 }
func (d ItemDelegate) Spacing() int                            { return 1 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, _ list.Model, index int, item list.Item) {
	i := item.(Item)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF69B4")).
		Render(i.title)

	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7AA2F7")).
		Render(i.description)

	selected := " "
	if index%2 == 0 {
		selected = "›"
	}

	fmt.Fprintf(w, "%s %s\n  %s", selected, title, desc)
}

// Model represents the application state
type Model struct {
	list list.Model
}

// InitialModel creates and initializes the model
func InitialModel() Model {
	items := []list.Item{
		Item{title: "Spinner", description: "Great for loading states and animations"},
		Item{title: "Text Input", description: "Single and multi-line text input fields"},
		Item{title: "List", description: "Scrollable lists with selection support"},
		Item{title: "Paginator", description: "Navigate through pages of content"},
		Item{title: "Viewport", description: "Scrollable content area with full control"},
		Item{title: "Progress", description: "Visual progress indicators"},
		Item{title: "Status Bar", description: "Display status information at bottom"},
		Item{title: "Help", description: "Interactive help and key binding display"},
	}

	l := list.New(items, ItemDelegate{}, 50, 15)
	l.Title = "🎨 Bubbletea Components"
	l.SetStatusBarItemName("component", "components")

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF69B4"))

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7AA2F7"))

	return Model{list: l}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if _, ok := m.list.SelectedItem().(Item); ok {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m Model) View() tea.View {
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9B1D6")).
		Italic(true).
		Render("\nNavigate with arrow keys, select with Enter, quit with Q")

	v := tea.NewView(m.list.View() + instructions)
	v.AltScreen = true
	return v
}

func main() {
	m := InitialModel()

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
