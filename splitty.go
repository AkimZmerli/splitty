package splitty

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AkimZmerli/splitty/terminal"
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
	shell           string
	env             []string
	theme           Theme
	keyMap          KeyMap
	minWidth        int
	minHeight       int
	statusBar       bool
	mouse           bool
	presetName      string
	scrollbackLines int
	scrollSpeed     int

	// State
	zoomed       bool
	zoomedID     string
	preZoomRoot  node
	broadcasting bool
	ready        bool
	themeIndex   int

	// Copy mode state
	copyMode CopyMode

	// Text selection state
	selection  Selection
	lastClickX int
	lastClickY int
	clickCount int
	lastClickT int64 // unix millis

	// Drag-to-resize state
	dragging       bool
	dragSplit      *splitNode
	dragDir        Direction
	dragOriginX    int
	dragOriginY    int
	dragStartRatio float64

	// Context menu
	menu contextMenu

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
		shell:           shell,
		theme:           TokyoNight,
		keyMap:          DefaultKeyMap(),
		minWidth:        10,
		minHeight:       3,
		statusBar:       true,
		mouse:           true,
		scrollbackLines: 1000,
		scrollSpeed:     3,
		menu:            newContextMenu(),
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
		cmd := m.readPaneCmd(msg.PaneID)
		return m, cmd

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
		content = lipgloss.JoinVertical(lipgloss.Left, content, bar)
	}

	if m.menu.visible {
		content = m.menu.overlayOnView(content, m.width, m.height, m.theme)
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
	m.menu.hide()

	if !m.ready {
		m.ready = true
		return m.initLayout()
	}

	// Recalculate layout
	m.layoutAll()
	return m, nil
}

// executeMenuAction dispatches the currently selected menu action.
func (m *Manager) executeMenuAction() (tea.Model, tea.Cmd) {
	action := m.menu.selectedAction()
	targetID := m.menu.paneID
	m.menu.hide()

	// Focus the target pane before acting on it
	m.focusedID = targetID

	switch action {
	case actionSplitVertical:
		return m.Split(Vertical)
	case actionSplitHorizontal:
		return m.Split(Horizontal)
	case actionClosePane:
		return m.ClosePane(targetID)
	}
	return m, nil
}

func (m *Manager) initLayout() (tea.Model, tea.Cmd) {
	h := m.statusBarHeight()

	if m.presetName != "" {
		if builder, ok := presetRegistry[m.presetName]; ok {
			m.root = builder(m.shell, m.env, m.width, h, m.scrollbackLines)
		}
	}

	if m.root == nil {
		pane := newPane(m.shell, m.env, m.width, h, m.scrollbackLines)
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
	// Intercept keys when context menu is visible
	if m.menu.visible {
		switch msg.Type {
		case tea.KeyEscape:
			m.menu.hide()
			return m, nil
		case tea.KeyUp:
			m.menu.moveUp()
			return m, nil
		case tea.KeyDown:
			m.menu.moveDown()
			return m, nil
		case tea.KeyEnter:
			return m.executeMenuAction()
		default:
			m.menu.hide()
			return m, nil
		}
	}

	// Copy mode intercepts all keys
	if m.copyMode.Active {
		keyStr := msg.String()
		m.handleCopyModeKey(keyStr)
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keyMap.EnterCopyMode):
		m.enterCopyMode()
		return m, nil
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
	case key.Matches(msg, m.keyMap.ScrollUp):
		if p := m.findPane(m.focusedID); p != nil {
			p.scrollUp(1)
		}
	case key.Matches(msg, m.keyMap.ScrollDown):
		if p := m.findPane(m.focusedID); p != nil {
			p.scrollDown(1)
		}
	case key.Matches(msg, m.keyMap.ScrollPageUp):
		if p := m.findPane(m.focusedID); p != nil {
			p.scrollUp(p.Height / 2)
		}
	case key.Matches(msg, m.keyMap.ScrollPageDown):
		if p := m.findPane(m.focusedID); p != nil {
			p.scrollDown(p.Height / 2)
		}
	case key.Matches(msg, m.keyMap.ScrollToTop):
		if p := m.findPane(m.focusedID); p != nil {
			p.scrollUp(999999)
		}
	case key.Matches(msg, m.keyMap.ScrollToBottom):
		if p := m.findPane(m.focusedID); p != nil {
			p.resetScroll()
		}
	case key.Matches(msg, m.keyMap.CycleTheme):
		m.themeIndex = (m.themeIndex + 1) % len(themeList)
		m.theme = themeList[m.themeIndex].Theme
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
	// Context menu interactions
	if m.menu.visible {
		contentH := m.statusBarHeight()
		switch {
		case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft:
			if m.menu.hitTest(msg.X, msg.Y, m.width, contentH) {
				idx := m.menu.hitTestItem(msg.X, msg.Y, m.width, contentH)
				if idx >= 0 {
					m.menu.selected = idx
					return m.executeMenuAction()
				}
				return m, nil
			}
			// Click outside — dismiss menu
			m.menu.hide()
			return m, nil

		case msg.Action == tea.MouseActionMotion:
			idx := m.menu.hitTestItem(msg.X, msg.Y, m.width, contentH)
			if idx >= 0 {
				m.menu.selected = idx
			}
			return m, nil
		}
		return m, nil
	}

	// Drag-to-resize handling (takes priority over everything else)
	if m.dragging {
		switch msg.Action {
		case tea.MouseActionMotion:
			h := m.statusBarHeight()
			if m.dragDir == Vertical {
				newRatio := float64(msg.X) / float64(m.width)
				m.dragSplit.ratio = clampRatio(newRatio)
			} else {
				newRatio := float64(msg.Y) / float64(h)
				m.dragSplit.ratio = clampRatio(newRatio)
			}
			m.layoutAll()
		case tea.MouseActionRelease:
			m.dragging = false
			m.dragSplit = nil
		}
		return m, nil
	}

	// Check for border click to start drag (not while zoomed)
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft && !m.zoomed {
		h := m.statusBarHeight()
		if sn := findBorder(m.root, msg.X, msg.Y, m.width, h); sn != nil {
			m.dragging = true
			m.dragSplit = sn
			m.dragDir = sn.dir
			m.dragOriginX = msg.X
			m.dragOriginY = msg.Y
			m.dragStartRatio = sn.ratio
			return m, nil
		}
	}

	// Find the pane under the mouse cursor
	hitPane := m.findPaneAt(msg.X, msg.Y)

	// If the pane has mouse mode and Shift is not held, forward events to PTY
	if hitPane != nil && hitPane.hasMouseMode() && !msg.Shift {
		button, action, ok := m.translateMouseEvent(msg)
		if !ok {
			return m, nil
		}
		// Translate to pane-relative coordinates (inside the border)
		termX := msg.X - hitPane.X - 1
		termY := msg.Y - hitPane.Y - 1
		// Clamp to pane dimensions
		if termX < 0 {
			termX = 0
		}
		if termY < 0 {
			termY = 0
		}
		if termX >= hitPane.Width {
			termX = hitPane.Width - 1
		}
		if termY >= hitPane.Height {
			termY = hitPane.Height - 1
		}
		hitPane.forwardMouse(button, action, termX, termY)
		// Also focus the pane on press
		if msg.Action == tea.MouseActionPress {
			m.focusedID = hitPane.ID
		}
		return m, nil
	}

	// --- Splitty-level mouse handling (no mouse mode or Shift bypass) ---

	// Mouse wheel scrolling
	if msg.Action == tea.MouseActionPress && (msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown) {
		if hitPane != nil {
			if msg.Button == tea.MouseButtonWheelUp {
				hitPane.scrollUp(m.scrollSpeed)
			} else {
				hitPane.scrollDown(m.scrollSpeed)
			}
		}
		return m, nil
	}

	// Right-click: show context menu (not while zoomed)
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonRight && !m.zoomed {
		if hitPane != nil {
			m.focusedID = hitPane.ID
			m.menu.show(msg.X, msg.Y, hitPane.ID)
		}
		return m, nil
	}

	// Left-click: focus + selection
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
		if hitPane != nil {
			m.focusedID = hitPane.ID
			termX := msg.X - hitPane.X - 1
			termY := msg.Y - hitPane.Y - 1
			if termX < 0 {
				termX = 0
			}
			if termY < 0 {
				termY = 0
			}

			// Multi-click detection
			now := time.Now().UnixMilli()
			if msg.X == m.lastClickX && msg.Y == m.lastClickY && now-m.lastClickT < 500 {
				m.clickCount++
				if m.clickCount > 3 {
					m.clickCount = 1
				}
			} else {
				m.clickCount = 1
			}
			m.lastClickX = msg.X
			m.lastClickY = msg.Y
			m.lastClickT = now

			switch m.clickCount {
			case 1:
				// Single click: clear selection, start char select
				m.selection = Selection{
					Active:   true,
					Mode:     SelectChar,
					PaneID:   hitPane.ID,
					StartRow: termY,
					StartCol: termX,
					EndRow:   termY,
					EndCol:   termX,
				}
			case 2:
				// Double click: select word
				m.selectWord(hitPane, termX, termY)
			case 3:
				// Triple click: select line
				m.selection = Selection{
					Active:   true,
					Mode:     SelectLine,
					PaneID:   hitPane.ID,
					StartRow: termY,
					StartCol: 0,
					EndRow:   termY,
					EndCol:   hitPane.Width - 1,
				}
			}
		} else {
			m.selection.Active = false
		}
		return m, nil
	}

	// Mouse drag: extend char selection
	if msg.Action == tea.MouseActionMotion && msg.Button == tea.MouseButtonLeft {
		if m.selection.Active && m.selection.Mode == SelectChar {
			if p := m.findPane(m.selection.PaneID); p != nil {
				termX := msg.X - p.X - 1
				termY := msg.Y - p.Y - 1
				if termX < 0 {
					termX = 0
				}
				if termY < 0 {
					termY = 0
				}
				if termX >= p.Width {
					termX = p.Width - 1
				}
				if termY >= p.Height {
					termY = p.Height - 1
				}
				m.selection.EndRow = termY
				m.selection.EndCol = termX
			}
		}
		return m, nil
	}

	// Mouse release: copy selection to clipboard
	if msg.Action == tea.MouseActionRelease {
		if m.selection.Active {
			m.copySelectionToClipboard()
		}
		return m, nil
	}

	return m, nil
}

// selectWord selects the word under the cursor in the given pane.
func (m *Manager) selectWord(p *Pane, col, row int) {
	cell := p.screen.GetCell(row, col)
	r := cell.Rune
	if r == 0 {
		r = ' '
	}

	// Find word boundaries
	startCol := col
	endCol := col

	if isWordChar(r) {
		// Expand left
		for startCol > 0 {
			c := p.screen.GetCell(row, startCol-1)
			cr := c.Rune
			if cr == 0 {
				cr = ' '
			}
			if !isWordChar(cr) {
				break
			}
			startCol--
		}
		// Expand right
		for endCol < p.Width-1 {
			c := p.screen.GetCell(row, endCol+1)
			cr := c.Rune
			if cr == 0 {
				cr = ' '
			}
			if !isWordChar(cr) {
				break
			}
			endCol++
		}
	}

	m.selection = Selection{
		Active:   true,
		Mode:     SelectWord,
		PaneID:   p.ID,
		StartRow: row,
		StartCol: startCol,
		EndRow:   row,
		EndCol:   endCol,
	}
}

// copySelectionToClipboard copies the selected text via OSC 52.
func (m *Manager) copySelectionToClipboard() {
	if !m.selection.Active {
		return
	}
	p := m.findPane(m.selection.PaneID)
	if p == nil {
		return
	}
	sr, sc, er, ec := m.selection.Normalize()
	text := p.screen.GetText(sr, sc, er, ec)
	if text == "" {
		return
	}

	// OSC 52 clipboard: \x1b]52;c;<base64>\x07
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	osc52 := fmt.Sprintf("\x1b]52;c;%s\x07", encoded)
	// Write to stdout (the outer terminal) to set clipboard
	os.Stdout.WriteString(osc52)
}

// isPaneAdjacentToDrag returns true if the pane is a direct child of the drag split.
func (m *Manager) isPaneAdjacentToDrag(p *Pane) bool {
	if m.dragSplit == nil {
		return false
	}
	for _, leaf := range m.dragSplit.first.leaves() {
		if leaf.pane == p {
			return true
		}
	}
	for _, leaf := range m.dragSplit.second.leaves() {
		if leaf.pane == p {
			return true
		}
	}
	return false
}

// findPaneAt returns the pane at the given screen coordinates, or nil.
func (m *Manager) findPaneAt(x, y int) *Pane {
	if m.root == nil {
		return nil
	}
	for _, leaf := range m.root.leaves() {
		p := leaf.pane
		if x >= p.X && x < p.X+p.Width+2 && // +2 for border
			y >= p.Y && y < p.Y+p.Height+2 {
			return p
		}
	}
	return nil
}

// translateMouseEvent converts a tea.MouseMsg into VT button and action values.
// Returns (button, action, ok).
func (m *Manager) translateMouseEvent(msg tea.MouseMsg) (int, int, bool) {
	var button int
	switch msg.Button {
	case tea.MouseButtonLeft:
		button = terminal.MouseButtonLeft
	case tea.MouseButtonMiddle:
		button = terminal.MouseButtonMiddle
	case tea.MouseButtonRight:
		button = terminal.MouseButtonRight
	case tea.MouseButtonWheelUp:
		button = terminal.MouseButtonWheelUp
	case tea.MouseButtonWheelDown:
		button = terminal.MouseButtonWheelDown
	default:
		// For motion without a button, use left as default
		button = terminal.MouseButtonLeft
	}

	var action int
	switch msg.Action {
	case tea.MouseActionPress:
		action = terminal.MouseActionPress
	case tea.MouseActionRelease:
		action = terminal.MouseActionRelease
	case tea.MouseActionMotion:
		action = terminal.MouseActionMotion
	default:
		return 0, 0, false
	}

	return button, action, true
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
		_, _ = leaf.pane.screen.Write(buf[:n])
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
	var content string
	if m.copyMode.Active && m.copyMode.PaneID == p.ID {
		highlight := m.getCopyModeHighlight(p)
		content = p.screen.RenderWithHighlight(highlight)
	} else if m.selection.Active && m.selection.PaneID == p.ID {
		content = p.renderWithSelection(&m.selection)
	} else {
		content = p.render()
	}
	var style lipgloss.Style

	// Copy mode border
	if m.copyMode.Active && m.copyMode.PaneID == p.ID {
		if p.ID == m.focusedID {
			style = m.theme.BorderCopyModeFocused
		} else {
			style = m.theme.BorderCopyMode
		}
	} else if m.dragging && m.dragSplit != nil && m.isPaneAdjacentToDrag(p) {
		style = m.theme.BorderResize
	} else if p.isScrolledBack() {
		if p.ID == m.focusedID {
			style = m.theme.BorderScrollbackFocused
		} else {
			style = m.theme.BorderScrollback
		}
	} else {
		if p.ID == m.focusedID {
			style = m.theme.BorderActive
		} else {
			style = m.theme.BorderInactive
		}
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
	parts = append(parts, "Ctrl+v split")
	parts = append(parts, "Ctrl+h hsplit")
	parts = append(parts, "Ctrl+q close")
	parts = append(parts, "Ctrl+wasd nav")

	if m.zoomed {
		parts = append(parts, m.theme.ZoomIndicator)
	}
	if m.broadcasting {
		parts = append(parts, m.theme.BroadcastIndicator)
	}
	if m.copyMode.Active {
		indicator := "[COPY]"
		if m.copyMode.Searching {
			if m.copyMode.SearchFwd {
				indicator = fmt.Sprintf("[COPY /%-10s]", m.copyMode.SearchBuf)
			} else {
				indicator = fmt.Sprintf("[COPY ?%-10s]", m.copyMode.SearchBuf)
			}
		} else if m.copyMode.Visual {
			indicator = "[COPY VISUAL]"
		}
		parts = append(parts, indicator)
	}

	left := strings.Join(parts, "  ")
	themeName := themeList[m.themeIndex].Name
	right := fmt.Sprintf("%s  Pane %d/%d", themeName, idx, paneCount)

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
