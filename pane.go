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

	screen *terminal.Screen
	pty    *terminal.PTY
	closed bool
}

func newPane(shell string, env []string, width, height int) *Pane {
	cwd, _ := os.Getwd()
	return &Pane{
		ID:     nextPaneID(),
		Title:  shell,
		CWD:    cwd,
		Width:  width,
		Height: height,
		screen: terminal.NewScreen(width, height),
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
