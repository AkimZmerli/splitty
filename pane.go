package splitty

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/AkimZmerli/splitty/terminal"
)

var paneCounter uint64

func nextPaneID() string {
	id := atomic.AddUint64(&paneCounter, 1)
	return fmt.Sprintf("pane-%d", id)
}

// Pane represents a single terminal pane with its own PTY session.
type Pane struct {
	ID     string
	Title  string
	CWD    string
	X      int
	Y      int
	Width  int
	Height int

	screen     *terminal.Screen
	pty        *terminal.PTY
	closed     bool
	autoScroll bool // Auto-scroll to bottom on new output
}

func newPane(shell string, env []string, width, height, scrollbackSize int) *Pane {
	cwd, _ := os.Getwd()
	return &Pane{
		ID:         nextPaneID(),
		Title:      shell,
		CWD:        cwd,
		Width:      width,
		Height:     height,
		screen:     terminal.NewScreenWithScrollback(width, height, scrollbackSize),
		autoScroll: true,
	}
}

// start launches the shell process in this pane's PTY.
func (p *Pane) start(shell string, env []string) error {
	if p.pty != nil {
		return nil
	}
	pt, err := terminal.NewPTY(shell, p.CWD, env, uint16(p.Width), uint16(p.Height))
	if err != nil {
		return fmt.Errorf("start pty: %w", err)
	}
	p.pty = pt
	return nil
}

// resize updates the pane dimensions, resizing both the screen buffer and PTY.
func (p *Pane) resize(width, height int) {
	p.Width = width
	p.Height = height
	p.screen.Resize(width, height)
	if p.pty != nil {
		_ = p.pty.Resize(uint16(width), uint16(height))
	}
}

// write sends input to the pane's PTY.
func (p *Pane) write(data []byte) {
	if p.pty != nil && !p.closed {
		_, _ = p.pty.Write(data)
	}
}

// close terminates the pane's shell process and PTY.
func (p *Pane) close() {
	if p.closed {
		return
	}
	p.closed = true
	if p.pty != nil {
		_ = p.pty.Close()
	}
}

// render returns the pane's current screen content as an ANSI string.
func (p *Pane) render() string {
	return p.screen.Render()
}

// scrollUp scrolls the view backward into history.
func (p *Pane) scrollUp(lines int) {
	p.screen.ScrollViewUp(lines)
	p.autoScroll = false
}

// scrollDown scrolls the view forward toward live output.
func (p *Pane) scrollDown(lines int) {
	offset := p.screen.ScrollViewDown(lines)
	if offset == 0 {
		p.autoScroll = true
	}
}

// resetScroll jumps to the bottom and resumes auto-scroll.
func (p *Pane) resetScroll() {
	p.screen.ScrollViewDown(999999)
	p.autoScroll = true
}

// isScrolledBack returns true if the pane is viewing history.
func (p *Pane) isScrolledBack() bool {
	return p.screen.IsScrolledBack()
}
