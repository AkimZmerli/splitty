package splitty

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Manager is the main Bubble Tea model for split pane terminal multiplexing.
// Create one with New() and pass it to tea.NewProgram().
type Manager struct {
	// Layout tree
	root      node
	focusedID string

	// Terminal dimensions
	width  int
	height int

	// Configuration
	shell      string
	env        []string
	theme      Theme
	keyMap     KeyMap
	minWidth   int
	minHeight  int
	statusBar  bool
	mouse      bool
	presetName string

	// State
	zoomed       bool
	zoomedID     string
	preZoomRoot  node
	broadcasting bool
	ready        bool

	// Runtime
	log *logger
}

// Ensure Manager implements tea.Model.
var _ tea.Model = (*Manager)(nil)

// New creates a new split pane Manager with the given options.
func New(opts ...Option) *Manager {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	m := &Manager{
		shell:     shell,
		theme:     DefaultTheme,
		keyMap:    DefaultKeyMap(),
		minWidth:  10,
		minHeight: 3,
		statusBar: true,
		mouse:     true,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Init initializes the Manager, creating the initial pane layout.
func (m *Manager) Init() tea.Cmd {
	return nil
}

// Update handles messages from Bubble Tea.
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		if m.mouse {
			return m.handleMouse(msg)
		}

	case PaneOutputMsg:
		return m, m.readPaneCmd(msg.PaneID)

	case PaneClosedMsg:
		if msg.Err != nil {
			m.closePane(msg.PaneID)
		}
	}

	return m, nil
}

// View renders the current state of all panes.
func (m *Manager) View() string {
	if m.root == nil {
		return "Initializing..."
	}

	var content string
	if m.zoomed {
		content = m.renderZoomed()
	} else {
		content = m.renderNode(m.root, m.width, m.statusBarHeight())
	}

	if m.statusBar {
		bar := m.renderStatusBar()
		return lipgloss.JoinVertical(lipgloss.Left, content, bar)
	}

	return content
}

func (m *Manager) statusBarHeight() int {
	if m.statusBar {
		return m.height - 1
	}
	return m.height
}

func (m *Manager) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	if !m.ready {
		m.ready = true
		return m.initLayout()
	}

	// Recalculate layout
	m.layoutAll()
	return m, nil
}

func (m *Manager) initLayout() (tea.Model, tea.Cmd) {
	h := m.statusBarHeight()

	if m.presetName != "" {
		if builder, ok := presetRegistry[m.presetName]; ok {
			m.root = builder(m.shell, m.env, m.width, h)
		}
	}

	if m.root == nil {
		pane := newPane(m.shell, m.env, m.width, h)
		m.root = &leafNode{pane: pane}
	}

	leaves := m.root.leaves()
	if len(leaves) > 0 {
		m.focusedID = leaves[0].pane.ID
	}

	m.layoutAll()

	// Start all panes and begin reading
	var cmds []tea.Cmd
	for _, leaf := range leaves {
		if err := leaf.pane.start(m.shell, m.env); err != nil {
			m.log.error("failed to start pane", "id", leaf.pane.ID, "err", err)
			continue
		}
		cmds = append(cmds, m.readPaneCmd(leaf.pane.ID))
	}
	return m, tea.Batch(cmds...)
}

func (m *Manager) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.SplitVertical):
		return m.Split(Vertical)
	case key.Matches(msg, m.keyMap.SplitHorizontal):
		return m.Split(Horizontal)
	case key.Matches(msg, m.keyMap.Close):
		return m.Close()
	case key.Matches(msg, m.keyMap.FocusLeft):
		m.focusDirection(Vertical, false)
	case key.Matches(msg, m.keyMap.FocusRight):
		m.focusDirection(Vertical, true)
	case key.Matches(msg, m.keyMap.FocusUp):
		m.focusDirection(Horizontal, false)
	case key.Matches(msg, m.keyMap.FocusDown):
		m.focusDirection(Horizontal, true)
	case key.Matches(msg, m.keyMap.FocusCycle):
		m.cycleFocus(true)
	case key.Matches(msg, m.keyMap.FocusCycleBack):
		m.cycleFocus(false)
	case key.Matches(msg, m.keyMap.Zoom):
		return m.toggleZoom()
	case key.Matches(msg, m.keyMap.Swap):
		m.Swap()
	case key.Matches(msg, m.keyMap.Broadcast):
		m.broadcasting = !m.broadcasting
	case key.Matches(msg, m.keyMap.ResizeLeft):
		m.Resize(Vertical, -0.05)
	case key.Matches(msg, m.keyMap.ResizeRight):
		m.Resize(Vertical, 0.05)
	case key.Matches(msg, m.keyMap.ResizeUp):
		m.Resize(Horizontal, -0.05)
	case key.Matches(msg, m.keyMap.ResizeDown):
		m.Resize(Horizontal, 0.05)
	default:
		// Forward keystrokes to pane(s)
		data := keyToBytes(msg)
		if len(data) > 0 {
			m.SendInput(data)
		}
	}
	return m, nil
}

func (m *Manager) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Click to focus: find which pane was clicked
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
		leaves := m.root.leaves()
		for _, leaf := range leaves {
			p := leaf.pane
			if msg.X >= p.X && msg.X < p.X+p.Width &&
				msg.Y >= p.Y && msg.Y < p.Y+p.Height {
				m.focusedID = p.ID
				break
			}
		}
	}
	return m, nil
}

// readPaneCmd returns a tea.Cmd that reads one chunk from a pane's PTY.
// After reading, it returns a PaneOutputMsg which triggers another read.
func (m *Manager) readPaneCmd(paneID string) tea.Cmd {
	return func() tea.Msg {
		leaf, _ := m.root.findLeaf(paneID)
		if leaf == nil || leaf.pane.closed || leaf.pane.pty == nil {
			return nil
		}
		buf := make([]byte, 4096)
		n, err := leaf.pane.pty.Read(buf)
		if err != nil {
			return PaneClosedMsg{PaneID: paneID, Err: err}
		}
		leaf.pane.screen.Write(buf[:n])
		return PaneOutputMsg{PaneID: paneID, Data: buf[:n]}
	}
}

func (m *Manager) findPane(id string) *Pane {
	if m.root == nil {
		return nil
	}
	leaf, _ := m.root.findLeaf(id)
	if leaf == nil {
		return nil
	}
	return leaf.pane
}

// renderNode recursively renders the layout tree.
func (m *Manager) renderNode(n node, width, height int) string {
	switch n := n.(type) {
	case *leafNode:
		return m.renderPane(n.pane, width, height)
	case *splitNode:
		return m.renderSplit(n, width, height)
	}
	return ""
}

func (m *Manager) renderPane(p *Pane, width, height int) string {
	content := p.render()
	style := m.theme.BorderInactive
	if p.ID == m.focusedID {
		style = m.theme.BorderActive
	}
	return style.
		Width(width - 2).  // account for border
		Height(height - 2). // account for border
		Render(content)
}

func (m *Manager) renderSplit(sn *splitNode, width, height int) string {
	if sn.dir == Vertical {
		leftW := int(float64(width) * sn.ratio)
		rightW := width - leftW
		if leftW < 1 {
			leftW = 1
		}
		if rightW < 1 {
			rightW = 1
		}
		left := m.renderNode(sn.first, leftW, height)
		right := m.renderNode(sn.second, rightW, height)
		return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	}

	// Horizontal
	topH := int(float64(height) * sn.ratio)
	bottomH := height - topH
	if topH < 1 {
		topH = 1
	}
	if bottomH < 1 {
		bottomH = 1
	}
	top := m.renderNode(sn.first, width, topH)
	bottom := m.renderNode(sn.second, width, bottomH)
	return lipgloss.JoinVertical(lipgloss.Left, top, bottom)
}

func (m *Manager) renderZoomed() string {
	p := m.findPane(m.zoomedID)
	if p == nil {
		return ""
	}
	return m.renderPane(p, m.width, m.statusBarHeight())
}

func (m *Manager) renderStatusBar() string {
	leaves := m.root.leaves()
	paneCount := len(leaves)

	// Find focused pane index
	idx := 0
	for i, l := range leaves {
		if l.pane.ID == m.focusedID {
			idx = i + 1
			break
		}
	}

	var parts []string
	parts = append(parts, "Ctrl+\\ split")
	parts = append(parts, "Ctrl+- hsplit")
	parts = append(parts, "Ctrl+w close")
	parts = append(parts, "Ctrl+hjkl nav")

	if m.zoomed {
		parts = append(parts, m.theme.ZoomIndicator)
	}
	if m.broadcasting {
		parts = append(parts, m.theme.BroadcastIndicator)
	}

	left := strings.Join(parts, "  ")
	right := fmt.Sprintf("Pane %d/%d", idx, paneCount)

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 0 {
		gap = 0
	}
	bar := left + strings.Repeat(" ", gap) + right
	return m.theme.StatusBar.Width(m.width).Render(bar)
}

// keyToBytes converts a tea.KeyMsg to raw bytes for PTY input.
func keyToBytes(msg tea.KeyMsg) []byte {
	// Handle special keys
	switch msg.Type {
	case tea.KeyEnter:
		return []byte{'\r'}
	case tea.KeyTab:
		return []byte{'\t'}
	case tea.KeyBackspace:
		return []byte{0x7f}
	case tea.KeyEscape:
		return []byte{0x1b}
	case tea.KeyUp:
		return []byte("\x1b[A")
	case tea.KeyDown:
		return []byte("\x1b[B")
	case tea.KeyRight:
		return []byte("\x1b[C")
	case tea.KeyLeft:
		return []byte("\x1b[D")
	case tea.KeyHome:
		return []byte("\x1b[H")
	case tea.KeyEnd:
		return []byte("\x1b[F")
	case tea.KeyPgUp:
		return []byte("\x1b[5~")
	case tea.KeyPgDown:
		return []byte("\x1b[6~")
	case tea.KeyDelete:
		return []byte("\x1b[3~")
	case tea.KeySpace:
		return []byte{' '}
	case tea.KeyRunes:
		return []byte(string(msg.Runes))
	}

	// Control keys (ctrl+a = 0x01, etc.)
	if msg.Type >= tea.KeyCtrlA && msg.Type <= tea.KeyCtrlZ {
		return []byte{byte(msg.Type - tea.KeyCtrlA + 1)}
	}

	return nil
}
