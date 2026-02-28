package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	stepForm     = 0
	stepReview   = 1
	stepConfirm  = 2
	maxFields    = 5
)

type Field struct {
	label       string
	placeholder string
	required    bool
}

type Model struct {
	step          int
	focusedField  int
	inputs        []textinput.Model
	fields        []Field
	formData      map[string]string
	errorMsg      string
	successMsg    string
	reviewport    viewport.Model
	confirmChoice int
	width         int
	height        int
}

func NewModel() Model {
	fields := []Field{
		{label: "Project Name", placeholder: "My Project", required: true},
		{label: "Description", placeholder: "What is this project about?", required: true},
		{label: "Author", placeholder: "Your Name", required: true},
		{label: "Email", placeholder: "your@email.com", required: true},
		{label: "Repository", placeholder: "https://github.com/user/project", required: false},
	}

	inputs := make([]textinput.Model, len(fields))
	for i, field := range fields {
		inputs[i] = textinput.New()
		inputs[i].Placeholder = field.placeholder
		inputs[i].CharLimit = 100
		inputs[i].Width = 40
	}
	inputs[0].Focus()

	vp := viewport.New(80, 15)
	vp.Style = lipgloss.NewStyle().Padding(1)

	return Model{
		step:         stepForm,
		focusedField: 0,
		inputs:       inputs,
		fields:       fields,
		formData:     make(map[string]string),
		width:        80,
		height:       24,
		reviewport:   vp,
		confirmChoice: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		if m.step == stepForm {
			return m.updateForm(msg)
		} else if m.step == stepReview {
			return m.updateReview(msg)
		} else if m.step == stepConfirm {
			return m.updateConfirm(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.reviewport.Width = msg.Width - 4
		m.reviewport.Height = msg.Height - 10
		return m, nil
	}

	return m, nil
}

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down", "j":
			m.focusedField = (m.focusedField + 1) % len(m.inputs)
			m.updateInputFocus()
			return m, nil

		case "up", "k":
			m.focusedField = (m.focusedField - 1 + len(m.inputs)) % len(m.inputs)
			m.updateInputFocus()
			return m, nil

		case "enter":
			m.errorMsg = ""
			// Validate required fields
			for i, field := range m.fields {
				if field.required && m.inputs[i].Value() == "" {
					m.errorMsg = fmt.Sprintf("'%s' is required", field.label)
					m.focusedField = i
					m.updateInputFocus()
					return m, nil
				}
			}
			// Save form data
			for i, field := range m.fields {
				m.formData[field.label] = m.inputs[i].Value()
			}
			m.step = stepReview
			m.updateReviewContent()
			return m, nil
		}
	}

	// Update inputs
	for i := range m.inputs {
		m.inputs[i], _ = m.inputs[i].Update(msg)
	}

	return m, nil
}

func (m Model) updateReview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			m.reviewport.LineUp(1)
			return m, nil
		case "down", "j":
			m.reviewport.LineDown(1)
			return m, nil
		case "enter":
			m.step = stepConfirm
			m.confirmChoice = 0
			return m, nil
		case "left", "h":
			m.step = stepForm
			return m, nil
		}
	}

	return m, nil
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "h":
			m.confirmChoice = 0
			return m, nil
		case "right", "l":
			m.confirmChoice = 1
			return m, nil
		case "enter":
			if m.confirmChoice == 0 {
				m.step = stepReview
			} else {
				m.successMsg = "✅ Configuration saved successfully!"
				return m, tea.Cmd(func() tea.Msg {
					return nil
				})
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) updateInputFocus() {
	for i := range m.inputs {
		if i == m.focusedField {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m *Model) updateReviewContent() {
	var content strings.Builder
	content.WriteString("\n📋 Review Your Information\n\n")
	for i, field := range m.fields {
		value := m.inputs[i].Value()
		if value == "" {
			value = "—"
		}
		content.WriteString(fmt.Sprintf("  %s:\n    %s\n\n", field.label, value))
	}
	content.WriteString("\nNaviage with arrow keys (↑/↓), Enter to confirm or ← to edit\n")
	m.reviewport.SetContent(content.String())
}

func (m Model) View() tea.View {
	var content string
	switch m.step {
	case stepForm:
		content = m.viewForm()
	case stepReview:
		content = m.viewReview()
	case stepConfirm:
		content = m.viewConfirm()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m Model) viewForm() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF69B4")).
		Render("📝 Configuration Form")

	var fields strings.Builder
	for i, field := range m.fields {
		label := field.label
		if field.required {
			label += " *"
		}

		labelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7AA2F7"))

		if i == m.focusedField {
			labelStyle = labelStyle.Bold(true)
		}

		fields.WriteString(labelStyle.Render(label) + "\n")
		fields.WriteString(m.inputs[i].View() + "\n\n")
	}

	var helpText string
	if m.errorMsg != "" {
		helpText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true).
			Render("❌ " + m.errorMsg)
	} else {
		helpText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5FFF87")).
			Render("✓ Fields marked with * are required")
	}

	content := fmt.Sprintf("%s\n\n%s\n\nNavigate: ↑/↓ | Submit: Enter | Quit: Ctrl+C\n%s",
		title, fields.String(), helpText)

	return lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(content)
}

func (m Model) viewReview() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7AA2F7")).
		Render("✓ Review Information")

	viewport := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(m.reviewport.View())

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9B1D6")).
		Render("↑/↓ Scroll | Enter Confirm | ← Edit")

	return lipgloss.NewStyle().
		Padding(2, 4).
		Render(title + "\n\n" + viewport + "\n\n" + footer)
}

func (m Model) viewConfirm() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#CBA6F7")).
		Render("⚡ Confirm Action")

	message := "Are you sure you want to save this configuration?"

	buttonNo := m.renderButton("No, Edit More", m.confirmChoice == 0)
	buttonYes := m.renderButton("Yes, Save", m.confirmChoice == 1)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		buttonNo,
		"    ",
		buttonYes,
	)

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9B1D6")).
		Render("← → Navigate | Enter Select")

	if m.successMsg != "" {
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5FFF87")).
			Bold(true).
			Padding(1).
			Render(m.successMsg)
		return lipgloss.NewStyle().
			Padding(2, 4).
			Render(title + "\n\n" + successStyle + "\n\nPress Q to quit")
	}

	content := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
		title, message, buttons, footer)

	return lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(content)
}

func (m Model) renderButton(label string, focused bool) string {
	style := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder())

	if focused {
		style = style.
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)
	} else {
		style = style.
			BorderForeground(lipgloss.Color("240")).
			Foreground(lipgloss.Color("240"))
	}

	return style.Render(label)
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
