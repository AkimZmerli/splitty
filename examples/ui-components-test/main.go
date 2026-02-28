package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Model represents the state of the UI test application
type Model struct {
	step       int          // 0: spinner, 1: textinput, 2: summary
	spinner    spinner.Model
	textinput  textinput.Model
	inputValue string
	err        error
}

// InitialModel creates and initializes the model
func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Enter your name..."
	ti.CharLimit = 50
	ti.SetWidth(30)

	return Model{
		step:      0,
		spinner:   s,
		textinput: ti,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	if m.step == 0 {
		return m.spinner.Tick
	}
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.step == 0 {
				// Move from spinner to text input
				m.step = 1
				m.textinput.Focus()
				return m, textinput.Blink
			} else if m.step == 1 {
				// Move from text input to summary
				m.inputValue = m.textinput.Value()
				m.step = 2
				return m, nil
			} else if m.step == 2 {
				// Restart
				return InitialModel(), m.spinner.Tick
			}
		}

		if m.step == 0 {
			// Spinner step: no input handling needed
			return m, m.spinner.Tick
		}

		if m.step == 1 {
			// Text input step
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

	case spinner.TickMsg:
		if m.step == 0 {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() tea.View {
	s := lipgloss.NewStyle()
	titleStyle := s.Copy().
		Bold(true).
		Foreground(lipgloss.Color("#FF69B4")).
		MarginBottom(1)

	instructionStyle := s.Copy().
		Foreground(lipgloss.Color("#7AA2F7")).
		Italic(true).
		MarginBottom(2)

	containerStyle := s.Copy().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	var rendered string
	switch m.step {
	case 0:
		// Spinner step
		content := fmt.Sprintf("%s Loading UI Components...\n\n", m.spinner.View())
		content += instructionStyle.Render("Press Enter to continue")
		rendered = containerStyle.Render(
			titleStyle.Render("🎨 Bubbletea Components Test") + "\n" + content,
		)

	case 1:
		// Text input step
		content := fmt.Sprintf("Enter your name:\n%s\n\n", m.textinput.View())
		content += instructionStyle.Render("Press Enter to submit")
		rendered = containerStyle.Render(
			titleStyle.Render("✏️  Text Input Component") + "\n" + content,
		)

	case 2:
		// Summary step
		successStyle := s.Copy().
			Foreground(lipgloss.Color("#5FFF87")).
			Bold(true)

		content := fmt.Sprintf("Hello, %s! 👋\n\n", successStyle.Render(m.inputValue))
		content += "You've successfully tested:\n"
		content += "  ✓ Spinner component\n"
		content += "  ✓ Text input component\n"
		content += "  ✓ Style composition with Lipgloss\n\n"
		content += instructionStyle.Render("Press Enter to restart, or Q to quit")

		rendered = containerStyle.Render(
			titleStyle.Render("🎉 Summary") + "\n" + content,
		)
	}

	v := tea.NewView(rendered)
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
