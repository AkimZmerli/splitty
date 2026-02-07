package splitty

// CopyMode represents the state of vim-style copy mode navigation.
type CopyMode struct {
	Active bool
	PaneID string

	// Cursor position (terminal-relative coordinates)
	CursorRow int
	CursorCol int

	// Visual selection
	Visual     bool
	VisualLine bool
	VisualRow  int
	VisualCol  int

	// Search
	Searching bool
	SearchBuf string
	SearchFwd bool // true = forward (/), false = backward (?)
	Matches   [][2]int
	MatchIdx  int
}

// handleCopyModeKey processes a key press while in copy mode.
// Returns true if the key was consumed.
func (m *Manager) handleCopyModeKey(key string) bool {
	p := m.findPane(m.copyMode.PaneID)
	if p == nil {
		m.exitCopyMode()
		return true
	}

	// Search input mode
	if m.copyMode.Searching {
		return m.handleSearchInput(key, p)
	}

	switch key {
	// Exit
	case "q", "esc":
		m.exitCopyMode()

	// Movement
	case "h":
		if m.copyMode.CursorCol > 0 {
			m.copyMode.CursorCol--
		}
	case "j":
		if m.copyMode.CursorRow < p.Height-1 {
			m.copyMode.CursorRow++
		} else {
			p.scrollDown(1)
		}
	case "k":
		if m.copyMode.CursorRow > 0 {
			m.copyMode.CursorRow--
		} else {
			p.scrollUp(1)
		}
	case "l":
		if m.copyMode.CursorCol < p.Width-1 {
			m.copyMode.CursorCol++
		}

	// Line movement
	case "0":
		m.copyMode.CursorCol = 0
	case "$":
		m.copyMode.CursorCol = p.Width - 1
	case "^":
		m.copyMode.CursorCol = m.firstNonBlank(p, m.copyMode.CursorRow)

	// Word movement
	case "w":
		m.copyMode.CursorCol = m.nextWord(p, m.copyMode.CursorRow, m.copyMode.CursorCol)
	case "b":
		m.copyMode.CursorCol = m.prevWord(p, m.copyMode.CursorRow, m.copyMode.CursorCol)
	case "e":
		m.copyMode.CursorCol = m.endWord(p, m.copyMode.CursorRow, m.copyMode.CursorCol)

	// Page movement
	case "ctrl+u":
		lines := p.Height / 2
		p.scrollUp(lines)
		m.copyMode.CursorRow -= lines
		if m.copyMode.CursorRow < 0 {
			m.copyMode.CursorRow = 0
		}
	case "ctrl+d":
		lines := p.Height / 2
		p.scrollDown(lines)
		m.copyMode.CursorRow += lines
		if m.copyMode.CursorRow >= p.Height {
			m.copyMode.CursorRow = p.Height - 1
		}

	// Buffer top/bottom
	case "g":
		p.scrollUp(999999)
		m.copyMode.CursorRow = 0
		m.copyMode.CursorCol = 0
	case "G":
		p.resetScroll()
		m.copyMode.CursorRow = p.Height - 1
		m.copyMode.CursorCol = 0

	// Visual mode
	case "v":
		if m.copyMode.Visual && !m.copyMode.VisualLine {
			m.copyMode.Visual = false
		} else {
			m.copyMode.Visual = true
			m.copyMode.VisualLine = false
			m.copyMode.VisualRow = m.copyMode.CursorRow
			m.copyMode.VisualCol = m.copyMode.CursorCol
		}
	case "V":
		if m.copyMode.VisualLine {
			m.copyMode.Visual = false
			m.copyMode.VisualLine = false
		} else {
			m.copyMode.Visual = true
			m.copyMode.VisualLine = true
			m.copyMode.VisualRow = m.copyMode.CursorRow
			m.copyMode.VisualCol = 0
		}

	// Yank
	case "y":
		if m.copyMode.Visual {
			m.yankVisualSelection(p)
			m.exitCopyMode()
		}

	// Search
	case "/":
		m.copyMode.Searching = true
		m.copyMode.SearchFwd = true
		m.copyMode.SearchBuf = ""
	case "?":
		m.copyMode.Searching = true
		m.copyMode.SearchFwd = false
		m.copyMode.SearchBuf = ""
	case "n":
		m.searchNext(p)
	case "N":
		m.searchPrev(p)

	default:
		return false
	}
	return true
}

func (m *Manager) handleSearchInput(key string, p *Pane) bool {
	switch key {
	case "esc":
		m.copyMode.Searching = false
		m.copyMode.SearchBuf = ""
		m.copyMode.Matches = nil
	case "enter":
		m.copyMode.Searching = false
		m.executeSearch(p)
	case "backspace":
		if len(m.copyMode.SearchBuf) > 0 {
			m.copyMode.SearchBuf = m.copyMode.SearchBuf[:len(m.copyMode.SearchBuf)-1]
		}
	default:
		if len(key) == 1 {
			m.copyMode.SearchBuf += key
		}
	}
	return true
}

func (m *Manager) executeSearch(p *Pane) {
	if m.copyMode.SearchBuf == "" {
		return
	}
	matches := p.screen.SearchText(m.copyMode.SearchBuf)
	m.copyMode.Matches = matches
	m.copyMode.MatchIdx = 0
	if len(matches) > 0 {
		m.jumpToMatch(p)
	}
}

func (m *Manager) searchNext(p *Pane) {
	if len(m.copyMode.Matches) == 0 {
		return
	}
	m.copyMode.MatchIdx = (m.copyMode.MatchIdx + 1) % len(m.copyMode.Matches)
	m.jumpToMatch(p)
}

func (m *Manager) searchPrev(p *Pane) {
	if len(m.copyMode.Matches) == 0 {
		return
	}
	m.copyMode.MatchIdx--
	if m.copyMode.MatchIdx < 0 {
		m.copyMode.MatchIdx = len(m.copyMode.Matches) - 1
	}
	m.jumpToMatch(p)
}

func (m *Manager) jumpToMatch(p *Pane) {
	if m.copyMode.MatchIdx >= len(m.copyMode.Matches) {
		return
	}
	match := m.copyMode.Matches[m.copyMode.MatchIdx]
	m.copyMode.CursorRow = match[0]
	m.copyMode.CursorCol = match[1]
	// Clamp
	if m.copyMode.CursorRow >= p.Height {
		m.copyMode.CursorRow = p.Height - 1
	}
	if m.copyMode.CursorCol >= p.Width {
		m.copyMode.CursorCol = p.Width - 1
	}
}

func (m *Manager) enterCopyMode() {
	p := m.findPane(m.focusedID)
	if p == nil {
		return
	}
	m.copyMode = CopyMode{
		Active:    true,
		PaneID:    p.ID,
		CursorRow: p.Height / 2,
		CursorCol: 0,
	}
}

func (m *Manager) exitCopyMode() {
	if m.copyMode.Active {
		if p := m.findPane(m.copyMode.PaneID); p != nil {
			p.resetScroll()
		}
	}
	m.copyMode = CopyMode{}
}

func (m *Manager) yankVisualSelection(p *Pane) {
	var sr, sc, er, ec int
	if m.copyMode.VisualLine {
		sr = m.copyMode.VisualRow
		sc = 0
		er = m.copyMode.CursorRow
		ec = p.Width - 1
	} else {
		sr = m.copyMode.VisualRow
		sc = m.copyMode.VisualCol
		er = m.copyMode.CursorRow
		ec = m.copyMode.CursorCol
	}
	if sr > er || (sr == er && sc > ec) {
		sr, sc, er, ec = er, ec, sr, sc
	}

	sel := Selection{
		Active:   true,
		Mode:     SelectChar,
		PaneID:   p.ID,
		StartRow: sr,
		StartCol: sc,
		EndRow:   er,
		EndCol:   ec,
	}
	if m.copyMode.VisualLine {
		sel.Mode = SelectLine
	}
	m.selection = sel
	m.copySelectionToClipboard()
	m.selection.Active = false
}

// Word navigation helpers

func (m *Manager) firstNonBlank(p *Pane, row int) int {
	for col := 0; col < p.Width; col++ {
		c := p.screen.GetCell(row, col)
		r := c.Rune
		if r == 0 {
			r = ' '
		}
		if r != ' ' {
			return col
		}
	}
	return 0
}

func (m *Manager) nextWord(p *Pane, row, col int) int {
	// Skip current word chars
	for col < p.Width-1 {
		c := p.screen.GetCell(row, col)
		r := c.Rune
		if r == 0 {
			r = ' '
		}
		if !isWordChar(r) {
			break
		}
		col++
	}
	// Skip non-word chars
	for col < p.Width-1 {
		c := p.screen.GetCell(row, col)
		r := c.Rune
		if r == 0 {
			r = ' '
		}
		if isWordChar(r) {
			break
		}
		col++
	}
	return col
}

func (m *Manager) prevWord(p *Pane, row, col int) int {
	if col > 0 {
		col--
	}
	// Skip non-word chars backward
	for col > 0 {
		c := p.screen.GetCell(row, col)
		r := c.Rune
		if r == 0 {
			r = ' '
		}
		if isWordChar(r) {
			break
		}
		col--
	}
	// Skip word chars backward to start of word
	for col > 0 {
		c := p.screen.GetCell(row, col-1)
		r := c.Rune
		if r == 0 {
			r = ' '
		}
		if !isWordChar(r) {
			break
		}
		col--
	}
	return col
}

func (m *Manager) endWord(p *Pane, row, col int) int {
	if col < p.Width-1 {
		col++
	}
	// Skip non-word chars forward
	for col < p.Width-1 {
		c := p.screen.GetCell(row, col)
		r := c.Rune
		if r == 0 {
			r = ' '
		}
		if isWordChar(r) {
			break
		}
		col++
	}
	// Skip word chars forward to end
	for col < p.Width-1 {
		c := p.screen.GetCell(row, col+1)
		r := c.Rune
		if r == 0 {
			r = ' '
		}
		if !isWordChar(r) {
			break
		}
		col++
	}
	return col
}

// getCopyModeHighlight builds a highlight grid for copy mode cursor and visual selection.
func (m *Manager) getCopyModeHighlight(p *Pane) [][]bool {
	highlight := make([][]bool, p.Height)
	for row := 0; row < p.Height; row++ {
		highlight[row] = make([]bool, p.Width)
	}

	// Cursor cell
	if m.copyMode.CursorRow >= 0 && m.copyMode.CursorRow < p.Height &&
		m.copyMode.CursorCol >= 0 && m.copyMode.CursorCol < p.Width {
		highlight[m.copyMode.CursorRow][m.copyMode.CursorCol] = true
	}

	// Visual selection
	if m.copyMode.Visual {
		sr, sc := m.copyMode.VisualRow, m.copyMode.VisualCol
		er, ec := m.copyMode.CursorRow, m.copyMode.CursorCol
		if sr > er || (sr == er && sc > ec) {
			sr, sc, er, ec = er, ec, sr, sc
		}
		for row := sr; row <= er && row < p.Height; row++ {
			if row < 0 {
				continue
			}
			if m.copyMode.VisualLine {
				for col := 0; col < p.Width; col++ {
					highlight[row][col] = true
				}
			} else {
				colStart := 0
				colEnd := p.Width - 1
				if row == sr {
					colStart = sc
				}
				if row == er {
					colEnd = ec
				}
				for col := colStart; col <= colEnd && col < p.Width; col++ {
					if col >= 0 {
						highlight[row][col] = true
					}
				}
			}
		}
	}

	// Search match highlights
	for _, match := range m.copyMode.Matches {
		row, col := match[0], match[1]
		if row >= 0 && row < p.Height && col >= 0 && col < p.Width {
			for i := 0; i < len(m.copyMode.SearchBuf) && col+i < p.Width; i++ {
				highlight[row][col+i] = true
			}
		}
	}

	return highlight
}
