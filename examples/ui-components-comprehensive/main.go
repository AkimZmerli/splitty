package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap defines custom key bindings
type KeyMap struct {
	Next    key.Binding
	Prev    key.Binding
	Help    key.Binding
	Quit    key.Binding
	Submit  key.Binding
	Tab     key.Binding
	ShiftTab key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Submit},
		{k.Tab, k.ShiftTab, k.Help, k.Quit},
	}
}

var keys = KeyMap{
	Next: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev field"),
	),
}

const (
	stateSpinner = iota
	stateTextInput
	stateTextArea
	stateProgress
	statePaginator
)

type Model struct {
	state         int
	activeField   int
	help          help.Model
	showHelp      bool
	width         int
	height        int
	windowHeight  int

	// Spinner state
	spinner spinner.Model
	loading int

	// Text input state
	inputs   []textinput.Model
	focused  int
	formData map[string]string

	// Text area state
	textarea textarea.Model
	textdata string

	// Progress state
	progress progress.Model
	progVal  float64

	// Paginator state
	paginator paginator.Model
	pages     []string
	pageIdx   int
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Initialize text inputs for form
	inputs := make([]textinput.Model, 3)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Name"
	inputs[0].CharLimit = 30
	inputs[0].Width = 30

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Email"
	inputs[1].CharLimit = 50
	inputs[1].Width = 30

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Message"
	inputs[2].CharLimit = 100
	inputs[2].Width = 30

	// Initialize text area
	ta := textarea.New()
	ta.Placeholder = "Enter some text..."
	ta.CharLimit = 500
	ta.SetWidth(50)
	ta.SetHeight(8)

	// Initialize progress bar
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(30),
		progress.WithoutPercentage(),
	)

	// Initialize paginator
	pages := []string{
		"Page 1: Welcome to the comprehensive Bubbletea component test!",
		"Page 2: This demonstrates pagination navigation.",
		"Page 3: You can move through pages with arrow keys.",
		"Page 4: Each page can contain different content.",
		"Page 5: This is the final page of the demo.",
	}
	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = 1
	p.SetTotalPages(len(pages))

	m := Model{
		spinner:      s,
		inputs:       inputs,
		focused:      0,
		formData:     make(map[string]string),
		textarea:     ta,
		progress:     prog,
		progVal:      0.0,
		paginator:    p,
		pages:        pages,
		pageIdx:      0,
		help:         help.New(),
		windowHeight: 20,
		width:        80,
		height:       24,
	}

	m.inputs[0].Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
			return m, nil
		case key.Matches(msg, keys.Next):
			if m.state < statePaginator {
				m.state++
				m.focused = 0
				if m.state == stateTextInput {
					m.inputs[0].Focus()
				} else if m.state == stateTextArea {
					m.textarea.Focus()
				}
			}
			return m, nil
		case key.Matches(msg, keys.Prev):
			if m.state > stateSpinner {
				m.state--
				m.focused = 0
				if m.state == stateTextInput {
					m.inputs[0].Focus()
				}
			}
			return m, nil
		}

		// Handle state-specific input
		switch m.state {
		case stateSpinner:
			return m, m.spinner.Tick

		case stateTextInput:
			if key.Matches(msg, keys.Tab) {
				m.focused = (m.focused + 1) % len(m.inputs)
				for i := range m.inputs {
					m.inputs[i].Blur()
				}
				m.inputs[m.focused].Focus()
				return m, nil
			}
			if key.Matches(msg, keys.ShiftTab) {
				m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
				for i := range m.inputs {
					m.inputs[i].Blur()
				}
				m.inputs[m.focused].Focus()
				return m, nil
			}
			var cmds []tea.Cmd
			for i := range m.inputs {
				m.inputs[i], _ = m.inputs[i].Update(msg)
			}
			m.formData["name"] = m.inputs[0].Value()
			m.formData["email"] = m.inputs[1].Value()
			m.formData["message"] = m.inputs[2].Value()
			return m, tea.Batch(cmds...)

		case stateTextArea:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			m.textdata = m.textarea.Value()
			return m, cmd

		case stateProgress:
			if key.Matches(msg, keys.Submit) {
				m.progVal = 0
			}
			m.progVal += 0.1
			if m.progVal > 1.0 {
				m.progVal = 1.0
			}
			return m, nil

		case statePaginator:
			if key.Matches(msg, keys.Next) {
				m.paginator.NextPage()
				m.pageIdx = m.paginator.Page
			} else if key.Matches(msg, keys.Prev) {
				m.paginator.PrevPage()
				m.pageIdx = m.paginator.Page
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.windowHeight = msg.Height - 8
		m.textarea.SetWidth(msg.Width - 4)
		m.textarea.SetHeight(6)
		return m, nil

	case spinner.TickMsg:
		if m.state == stateSpinner {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case progress.FrameMsg:
		if m.state == stateProgress {
			var cmd tea.Cmd
			model, cmd := m.progress.Update(msg)
			m.progress = model.(progress.Model)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.showHelp {
		return m.helpView()
	}

	content := ""
	switch m.state {
	case stateSpinner:
		content = m.spinnerView()
	case stateTextInput:
		content = m.textInputView()
	case stateTextArea:
		content = m.textAreaView()
	case stateProgress:
		content = m.progressView()
	case statePaginator:
		content = m.paginatorView()
	}

	footer := m.footerView()
	return content + "\n" + footer
}

func (m Model) spinnerView() string {
	title := titleStyle().Render("⏳ Spinner Component")
	spinner := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(fmt.Sprintf("%s Loading...", m.spinner.View()))
	info := infoStyle().Render("Press Right/L to continue, Q to quit")
	return containerStyle().Render(title + "\n\n" + spinner + "\n\n" + info)
}

func (m Model) textInputView() string {
	title := titleStyle().Render("📝 Text Input Component")
	var b strings.Builder
	for i, input := range m.inputs {
		b.WriteString(input.View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n\n")
		}
	}
	info := infoStyle().Render("Tab/Shift+Tab to navigate fields, Right/L to continue")
	return containerStyle().Render(title + "\n\n" + b.String() + "\n\n" + info)
}

func (m Model) textAreaView() string {
	title := titleStyle().Render("📄 Text Area Component")
	info := infoStyle().Render("Type to edit, Right/L to continue, Q to quit")
	preview := previewStyle().Render(fmt.Sprintf("Length: %d / 500", len(m.textdata)))
	return containerStyle().Render(title + "\n\n" + m.textarea.View() + "\n\n" + preview + "\n" + info)
}

func (m Model) progressView() string {
	title := titleStyle().Render("📊 Progress Bar Component")
	prog := m.progress.ViewAs(m.progVal)
	percent := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Render(fmt.Sprintf("%.0f%%", m.progVal*100))
	info := infoStyle().Render("Press Enter to reset progress, Right/L to continue")
	return containerStyle().Render(title + "\n\n" + prog + "\n" + percent + "\n\n" + info)
}

func (m Model) paginatorView() string {
	title := titleStyle().Render("📖 Paginator Component")
	page := lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(m.pages[m.pageIdx])
	paginatorView := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(m.paginator.View())
	info := infoStyle().Render("Left/H and Right/L to navigate pages")
	return containerStyle().Render(title + "\n\n" + page + "\n\n" + paginatorView + "\n" + info)
}

func (m Model) helpView() string {
	return m.help.View(keys)
}

func (m Model) footerView() string {
	stateNames := []string{"Spinner", "Text Input", "Text Area", "Progress", "Paginator"}
	stateIndicators := make([]string, len(stateNames))
	for i, name := range stateNames {
		if i == m.state {
			stateIndicators[i] = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("63")).
				Render("● " + name)
		} else {
			stateIndicators[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("○ " + name)
		}
	}
	return lipgloss.NewStyle().
		Padding(0, 1).
		Render(strings.Join(stateIndicators, " | ") + " | Press ? for help")
}

// Style functions
func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF69B4")).
		MarginBottom(1)
}

func infoStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7AA2F7")).
		Italic(true)
}

func previewStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9B1D6"))
}

func containerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
