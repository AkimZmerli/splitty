package terminal

import (
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

type parseState int

const (
	stateGround  parseState = iota
	stateEscape             // Received ESC
	stateCSI                // Received ESC [
	stateOSC                // Received ESC ]
	stateCharset            // Received ESC (
)

// Screen is a virtual terminal that parses ANSI escape sequences
// and maintains a cell grid for rendering.
type Screen struct {
	mu           sync.RWMutex
	width        int
	height       int
	cells        [][]Cell
	cursor       Cursor
	savedCursor  Cursor
	scrollTop    int
	scrollBottom int
	title        string

	// Parser state machine
	state    parseState
	paramBuf []byte
	oscBuf   []byte
	private  bool

	// UTF-8 multi-byte accumulator
	utfBuf []byte

	// Pending wrap: cursor hit right margin but hasn't wrapped yet.
	// Next printable char will wrap; CR/LF/cursor-move will clear it.
	wrapPending bool
}

// NewScreen creates a new screen with the given dimensions.
func NewScreen(width, height int) *Screen {
	s := &Screen{
		width:        width,
		height:       height,
		cursor:       DefaultCursor(),
		scrollTop:    0,
		scrollBottom: height - 1,
		state:        stateGround,
		paramBuf:     make([]byte, 0, 64),
		oscBuf:       make([]byte, 0, 256),
		utfBuf:       make([]byte, 0, 4),
	}
	s.cells = s.newGrid(width, height)
	return s
}

func (s *Screen) newGrid(width, height int) [][]Cell {
	grid := make([][]Cell, height)
	for i := range grid {
		grid[i] = make([]Cell, width)
		for j := range grid[i] {
			grid[i][j] = EmptyCell()
		}
	}
	return grid
}

// Resize changes the screen dimensions, preserving existing content.
func (s *Screen) Resize(width, height int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resizeInternal(width, height)
}

func (s *Screen) resizeInternal(width, height int) {
	newGrid := s.newGrid(width, height)
	minH := s.height
	if height < minH {
		minH = height
	}
	minW := s.width
	if width < minW {
		minW = width
	}
	for row := 0; row < minH; row++ {
		for col := 0; col < minW; col++ {
			newGrid[row][col] = s.cells[row][col]
		}
	}
	s.cells = newGrid
	s.width = width
	s.height = height
	s.scrollBottom = height - 1
	if s.scrollTop >= height {
		s.scrollTop = 0
	}
	if s.cursor.Row >= height {
		s.cursor.Row = height - 1
	}
	if s.cursor.Col >= width {
		s.cursor.Col = width - 1
	}
}

// Width returns the screen width.
func (s *Screen) Width() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.width
}

// Height returns the screen height.
func (s *Screen) Height() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.height
}

// Title returns the terminal title set via OSC sequences.
func (s *Screen) Title() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.title
}

// Write processes raw terminal output bytes, parsing ANSI escape sequences
// and updating the cell grid. Implements io.Writer.
func (s *Screen) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range data {
		s.processByte(b)
	}
	return len(data), nil
}

func (s *Screen) processByte(b byte) {
	switch s.state {
	case stateGround:
		s.processGround(b)
	case stateEscape:
		s.processEscape(b)
	case stateCSI:
		s.processCSI(b)
	case stateOSC:
		s.processOSC(b)
	case stateCharset:
		s.state = stateGround // consume one byte and return
	}
}

func (s *Screen) flushUtfBuf() {
	for range s.utfBuf {
		s.writeRune(unicode.ReplacementChar)
	}
	s.utfBuf = s.utfBuf[:0]
}

func (s *Screen) processGround(b byte) {
	// Handle control characters — flush any pending UTF-8 first
	switch {
	case b == 0x1b: // ESC
		s.flushUtfBuf()
		s.wrapPending = false
		s.state = stateEscape
		return
	case b == '\n': // LF
		s.flushUtfBuf()
		s.wrapPending = false
		s.lineFeed()
		return
	case b == '\r': // CR
		s.flushUtfBuf()
		s.wrapPending = false
		s.cursor.Col = 0
		return
	case b == '\b': // BS
		s.flushUtfBuf()
		s.wrapPending = false
		if s.cursor.Col > 0 {
			s.cursor.Col--
		}
		return
	case b == '\t': // TAB
		s.flushUtfBuf()
		s.wrapPending = false
		s.cursor.Col = (s.cursor.Col/8 + 1) * 8
		if s.cursor.Col >= s.width {
			s.cursor.Col = s.width - 1
		}
		return
	case b == 0x07: // BEL
		s.flushUtfBuf()
		return
	case b < 0x20: // other C0 controls
		s.flushUtfBuf()
		return
	}

	// Printable bytes — handle UTF-8 decoding
	switch {
	case b < 0x80:
		// ASCII fast path
		s.flushUtfBuf()
		s.writeRune(rune(b))

	case b >= 0xC0 && b <= 0xF7:
		// Leading byte of a multi-byte sequence
		s.flushUtfBuf()
		s.utfBuf = append(s.utfBuf, b)

	case b >= 0x80 && b < 0xC0:
		// Continuation byte
		if len(s.utfBuf) == 0 {
			// Orphan continuation byte
			s.writeRune(unicode.ReplacementChar)
			return
		}
		s.utfBuf = append(s.utfBuf, b)
		if utf8.FullRune(s.utfBuf) {
			r, size := utf8.DecodeRune(s.utfBuf)
			if r == utf8.RuneError && size <= 1 {
				r = unicode.ReplacementChar
			}
			s.writeRune(r)
			s.utfBuf = s.utfBuf[:0]
		} else if len(s.utfBuf) >= 4 {
			// Too many bytes without completing — flush as errors
			s.flushUtfBuf()
		}

	default:
		// Invalid byte (0xF8-0xFF)
		s.flushUtfBuf()
		s.writeRune(unicode.ReplacementChar)
	}
}

func (s *Screen) processEscape(b byte) {
	switch b {
	case '[':
		s.state = stateCSI
		s.paramBuf = s.paramBuf[:0]
		s.private = false
	case ']':
		s.state = stateOSC
		s.oscBuf = s.oscBuf[:0]
	case '(':
		s.state = stateCharset
	case '7': // DECSC
		s.savedCursor = s.cursor
		s.state = stateGround
	case '8': // DECRC
		s.cursor = s.savedCursor
		s.state = stateGround
	case 'D': // IND
		s.lineFeed()
		s.state = stateGround
	case 'M': // RI
		s.reverseIndex()
		s.state = stateGround
	case 'c': // RIS
		s.fullReset()
		s.state = stateGround
	default:
		s.state = stateGround
	}
}

func (s *Screen) processCSI(b byte) {
	switch {
	case b == '?':
		s.private = true
	case b >= '0' && b <= '9', b == ';':
		s.paramBuf = append(s.paramBuf, b)
	case b >= 0x40: // final byte
		s.dispatchCSI(b)
		s.state = stateGround
	}
}

func (s *Screen) processOSC(b byte) {
	switch {
	case b == 0x07: // BEL terminates
		s.handleOSC()
		s.state = stateGround
	case b == 0x1b: // ESC might start ST
		s.handleOSC()
		s.state = stateEscape
	default:
		s.oscBuf = append(s.oscBuf, b)
	}
}

func (s *Screen) handleOSC() {
	data := string(s.oscBuf)
	if strings.HasPrefix(data, "0;") || strings.HasPrefix(data, "2;") {
		s.title = data[2:]
	}
}

// parseParams splits the CSI parameter buffer into integers.
func (s *Screen) parseParams() []int {
	if len(s.paramBuf) == 0 {
		return nil
	}
	parts := strings.Split(string(s.paramBuf), ";")
	params := make([]int, len(parts))
	for i, p := range parts {
		if p == "" {
			params[i] = 0
		} else {
			n, _ := strconv.Atoi(p)
			params[i] = n
		}
	}
	return params
}

func (s *Screen) paramDefault(idx, def int, params []int) int {
	if idx < len(params) && params[idx] != 0 {
		return params[idx]
	}
	return def
}

func (s *Screen) dispatchCSI(final byte) {
	s.wrapPending = false
	params := s.parseParams()

	switch final {
	case 'A': // CUU - cursor up
		n := s.paramDefault(0, 1, params)
		s.cursor.Row = clamp(s.cursor.Row-n, 0, s.height-1)

	case 'B': // CUD - cursor down
		n := s.paramDefault(0, 1, params)
		s.cursor.Row = clamp(s.cursor.Row+n, 0, s.height-1)

	case 'C': // CUF - cursor forward
		n := s.paramDefault(0, 1, params)
		s.cursor.Col = clamp(s.cursor.Col+n, 0, s.width-1)

	case 'D': // CUB - cursor back
		n := s.paramDefault(0, 1, params)
		s.cursor.Col = clamp(s.cursor.Col-n, 0, s.width-1)

	case 'E': // CNL - cursor next line
		n := s.paramDefault(0, 1, params)
		s.cursor.Col = 0
		s.cursor.Row = clamp(s.cursor.Row+n, 0, s.height-1)

	case 'F': // CPL - cursor previous line
		n := s.paramDefault(0, 1, params)
		s.cursor.Col = 0
		s.cursor.Row = clamp(s.cursor.Row-n, 0, s.height-1)

	case 'G': // CHA - cursor horizontal absolute
		col := s.paramDefault(0, 1, params) - 1
		s.cursor.Col = clamp(col, 0, s.width-1)

	case 'H', 'f': // CUP - cursor position
		row := s.paramDefault(0, 1, params) - 1
		col := s.paramDefault(1, 1, params) - 1
		s.cursor.Row = clamp(row, 0, s.height-1)
		s.cursor.Col = clamp(col, 0, s.width-1)

	case 'J': // ED - erase display
		s.eraseDisplay(s.paramDefault(0, 0, params))

	case 'K': // EL - erase line
		s.eraseLine(s.paramDefault(0, 0, params))

	case 'L': // IL - insert lines
		s.insertLines(s.paramDefault(0, 1, params))

	case 'M': // DL - delete lines
		s.deleteLines(s.paramDefault(0, 1, params))

	case 'P': // DCH - delete characters
		s.deleteChars(s.paramDefault(0, 1, params))

	case '@': // ICH - insert characters
		s.insertChars(s.paramDefault(0, 1, params))

	case 'X': // ECH - erase characters
		s.eraseChars(s.paramDefault(0, 1, params))

	case 'S': // SU - scroll up
		s.scrollUp(s.paramDefault(0, 1, params))

	case 'T': // SD - scroll down
		s.scrollDown(s.paramDefault(0, 1, params))

	case 'd': // VPA - line position absolute
		row := s.paramDefault(0, 1, params) - 1
		s.cursor.Row = clamp(row, 0, s.height-1)

	case 'r': // DECSTBM - set scrolling region
		top := s.paramDefault(0, 1, params) - 1
		bottom := s.paramDefault(1, s.height, params) - 1
		s.scrollTop = clamp(top, 0, s.height-1)
		s.scrollBottom = clamp(bottom, 0, s.height-1)
		s.cursor.Row = 0
		s.cursor.Col = 0

	case 'm': // SGR
		s.handleSGR(params)

	case 'h': // SM - set mode
		if s.private {
			s.handleDECSet(params)
		}

	case 'l': // RM - reset mode
		if s.private {
			s.handleDECReset(params)
		}

	case 's': // SCP - save cursor
		s.savedCursor = s.cursor

	case 'u': // RCP - restore cursor
		s.cursor = s.savedCursor

	case 'n', 'c': // DSR, DA - device reports (ignore)
	}
}

func (s *Screen) handleSGR(params []int) {
	if len(params) == 0 {
		params = []int{0}
	}
	for i := 0; i < len(params); i++ {
		p := params[i]
		switch {
		case p == 0:
			s.cursor.Style = DefaultStyle()
		case p == 1:
			s.cursor.Style.Bold = true
		case p == 2:
			s.cursor.Style.Dim = true
		case p == 3:
			s.cursor.Style.Italic = true
		case p == 4:
			s.cursor.Style.Underline = true
		case p == 5:
			s.cursor.Style.Blink = true
		case p == 7:
			s.cursor.Style.Reverse = true
		case p == 8:
			s.cursor.Style.Hidden = true
		case p == 9:
			s.cursor.Style.Strike = true
		case p == 22:
			s.cursor.Style.Bold = false
			s.cursor.Style.Dim = false
		case p == 23:
			s.cursor.Style.Italic = false
		case p == 24:
			s.cursor.Style.Underline = false
		case p == 25:
			s.cursor.Style.Blink = false
		case p == 27:
			s.cursor.Style.Reverse = false
		case p == 28:
			s.cursor.Style.Hidden = false
		case p == 29:
			s.cursor.Style.Strike = false
		case p >= 30 && p <= 37:
			s.cursor.Style.FG = Color{Type: Color16, Value: uint32(p - 30)}
		case p == 38: // extended foreground
			i = s.parseExtendedColor(params, i, true)
		case p == 39:
			s.cursor.Style.FG = Color{}
		case p >= 40 && p <= 47:
			s.cursor.Style.BG = Color{Type: Color16, Value: uint32(p - 40)}
		case p == 48: // extended background
			i = s.parseExtendedColor(params, i, false)
		case p == 49:
			s.cursor.Style.BG = Color{}
		case p >= 90 && p <= 97:
			s.cursor.Style.FG = Color{Type: Color16, Value: uint32(p - 90 + 8)}
		case p >= 100 && p <= 107:
			s.cursor.Style.BG = Color{Type: Color16, Value: uint32(p - 100 + 8)}
		}
	}
}

func (s *Screen) parseExtendedColor(params []int, i int, fg bool) int {
	if i+1 >= len(params) {
		return i
	}
	switch params[i+1] {
	case 5: // 256 color
		if i+2 < len(params) {
			c := Color{Type: Color256, Value: uint32(params[i+2])}
			if fg {
				s.cursor.Style.FG = c
			} else {
				s.cursor.Style.BG = c
			}
			return i + 2
		}
	case 2: // RGB
		if i+4 < len(params) {
			r, g, b := params[i+2], params[i+3], params[i+4]
			c := Color{Type: ColorRGB, Value: uint32(r)<<16 | uint32(g)<<8 | uint32(b)}
			if fg {
				s.cursor.Style.FG = c
			} else {
				s.cursor.Style.BG = c
			}
			return i + 4
		}
	}
	return i
}

func (s *Screen) handleDECSet(params []int) {
	for _, p := range params {
		switch p {
		case 25: // DECTCEM - show cursor
			s.cursor.Visible = true
		case 1049, 47, 1047: // alternate screen buffer
			s.clearScreen()
		}
	}
}

func (s *Screen) handleDECReset(params []int) {
	for _, p := range params {
		switch p {
		case 25:
			s.cursor.Visible = false
		case 1049, 47, 1047:
			s.clearScreen()
		}
	}
}

// --- Terminal operations ---

func (s *Screen) writeRune(r rune) {
	if s.wrapPending {
		s.cursor.Col = 0
		s.lineFeed()
		s.wrapPending = false
	}
	if s.cursor.Row >= 0 && s.cursor.Row < s.height &&
		s.cursor.Col >= 0 && s.cursor.Col < s.width {
		s.cells[s.cursor.Row][s.cursor.Col] = Cell{
			Rune:  r,
			Style: s.cursor.Style,
		}
	}
	s.cursor.Col++
	if s.cursor.Col >= s.width {
		s.cursor.Col = s.width - 1
		s.wrapPending = true
	}
}

func (s *Screen) lineFeed() {
	if s.cursor.Row == s.scrollBottom {
		s.scrollUp(1)
	} else if s.cursor.Row < s.height-1 {
		s.cursor.Row++
	}
}

func (s *Screen) reverseIndex() {
	if s.cursor.Row == s.scrollTop {
		s.scrollDown(1)
	} else if s.cursor.Row > 0 {
		s.cursor.Row--
	}
}

func (s *Screen) scrollUp(n int) {
	for i := 0; i < n; i++ {
		for row := s.scrollTop; row < s.scrollBottom; row++ {
			s.cells[row] = s.cells[row+1]
		}
		s.cells[s.scrollBottom] = make([]Cell, s.width)
		for j := range s.cells[s.scrollBottom] {
			s.cells[s.scrollBottom][j] = EmptyCell()
		}
	}
}

func (s *Screen) scrollDown(n int) {
	for i := 0; i < n; i++ {
		for row := s.scrollBottom; row > s.scrollTop; row-- {
			s.cells[row] = s.cells[row-1]
		}
		s.cells[s.scrollTop] = make([]Cell, s.width)
		for j := range s.cells[s.scrollTop] {
			s.cells[s.scrollTop][j] = EmptyCell()
		}
	}
}

func (s *Screen) eraseDisplay(mode int) {
	switch mode {
	case 0: // below
		s.eraseLine(0)
		for row := s.cursor.Row + 1; row < s.height; row++ {
			for col := 0; col < s.width; col++ {
				s.cells[row][col] = EmptyCell()
			}
		}
	case 1: // above
		s.eraseLine(1)
		for row := 0; row < s.cursor.Row; row++ {
			for col := 0; col < s.width; col++ {
				s.cells[row][col] = EmptyCell()
			}
		}
	case 2, 3: // all
		s.clearScreen()
	}
}

func (s *Screen) eraseLine(mode int) {
	if s.cursor.Row < 0 || s.cursor.Row >= s.height {
		return
	}
	switch mode {
	case 0: // right
		for col := s.cursor.Col; col < s.width; col++ {
			s.cells[s.cursor.Row][col] = EmptyCell()
		}
	case 1: // left
		end := s.cursor.Col
		if end >= s.width {
			end = s.width - 1
		}
		for col := 0; col <= end; col++ {
			s.cells[s.cursor.Row][col] = EmptyCell()
		}
	case 2: // entire line
		for col := 0; col < s.width; col++ {
			s.cells[s.cursor.Row][col] = EmptyCell()
		}
	}
}

func (s *Screen) insertLines(n int) {
	if s.cursor.Row < s.scrollTop || s.cursor.Row > s.scrollBottom {
		return
	}
	for i := 0; i < n; i++ {
		for row := s.scrollBottom; row > s.cursor.Row; row-- {
			s.cells[row] = s.cells[row-1]
		}
		s.cells[s.cursor.Row] = make([]Cell, s.width)
		for j := range s.cells[s.cursor.Row] {
			s.cells[s.cursor.Row][j] = EmptyCell()
		}
	}
}

func (s *Screen) deleteLines(n int) {
	if s.cursor.Row < s.scrollTop || s.cursor.Row > s.scrollBottom {
		return
	}
	for i := 0; i < n; i++ {
		for row := s.cursor.Row; row < s.scrollBottom; row++ {
			s.cells[row] = s.cells[row+1]
		}
		s.cells[s.scrollBottom] = make([]Cell, s.width)
		for j := range s.cells[s.scrollBottom] {
			s.cells[s.scrollBottom][j] = EmptyCell()
		}
	}
}

func (s *Screen) deleteChars(n int) {
	row := s.cursor.Row
	if row < 0 || row >= s.height {
		return
	}
	col := s.cursor.Col
	end := s.width - n
	if end < col {
		end = col
	}
	for i := col; i < end; i++ {
		s.cells[row][i] = s.cells[row][i+n]
	}
	for i := end; i < s.width; i++ {
		s.cells[row][i] = EmptyCell()
	}
}

func (s *Screen) insertChars(n int) {
	row := s.cursor.Row
	if row < 0 || row >= s.height {
		return
	}
	col := s.cursor.Col
	for i := s.width - 1; i >= col+n; i-- {
		s.cells[row][i] = s.cells[row][i-n]
	}
	end := col + n
	if end > s.width {
		end = s.width
	}
	for i := col; i < end; i++ {
		s.cells[row][i] = EmptyCell()
	}
}

func (s *Screen) eraseChars(n int) {
	row := s.cursor.Row
	if row < 0 || row >= s.height {
		return
	}
	end := s.cursor.Col + n
	if end > s.width {
		end = s.width
	}
	for i := s.cursor.Col; i < end; i++ {
		s.cells[row][i] = EmptyCell()
	}
}

func (s *Screen) clearScreen() {
	for row := 0; row < s.height; row++ {
		for col := 0; col < s.width; col++ {
			s.cells[row][col] = EmptyCell()
		}
	}
	s.cursor.Row = 0
	s.cursor.Col = 0
	s.wrapPending = false
}

func (s *Screen) fullReset() {
	s.clearScreen()
	s.cursor = DefaultCursor()
	s.scrollTop = 0
	s.scrollBottom = s.height - 1
	s.title = ""
	s.utfBuf = s.utfBuf[:0]
	s.wrapPending = false
}

// Render converts the cell grid to an ANSI string suitable for display
// inside a Bubble Tea View.
func (s *Screen) Render() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var buf strings.Builder
	buf.Grow(s.width * s.height * 2)

	var prevStyle Style
	first := true

	for row := 0; row < s.height; row++ {
		if row > 0 {
			buf.WriteByte('\n')
		}
		for col := 0; col < s.width; col++ {
			cell := s.cells[row][col]
			renderStyle := cell.Style
			if s.cursor.Visible && row == s.cursor.Row && col == s.cursor.Col {
				renderStyle.Reverse = !renderStyle.Reverse
			}
			if first || !renderStyle.Equal(prevStyle) {
				buf.WriteString(renderStyle.ToANSI())
				prevStyle = renderStyle
				first = false
			}
			if cell.Rune == 0 {
				buf.WriteByte(' ')
			} else {
				buf.WriteRune(cell.Rune)
			}
		}
	}
	buf.WriteString("\x1b[0m")
	return buf.String()
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
