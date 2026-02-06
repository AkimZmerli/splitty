package terminal

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// PTY wraps a pseudo-terminal for a shell process.
type PTY struct {
	ptmx *os.File
	cmd  *exec.Cmd
}

// NewPTY starts a new shell process attached to a PTY with the given dimensions.
func NewPTY(shell, cwd string, env []string, cols, rows uint16) (*PTY, error) {
	cmd := exec.Command(shell)
	if cwd != "" {
		cmd.Dir = cwd
	}
	cmd.Env = append(os.Environ(), env...)

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
	if err != nil {
		return nil, err
	}

	return &PTY{
		ptmx: ptmx,
		cmd:  cmd,
	}, nil
}

// Read reads output from the PTY. Blocks until data is available.
func (p *PTY) Read(buf []byte) (int, error) {
	return p.ptmx.Read(buf)
}

// Write sends input to the PTY (forwarded to the shell process).
func (p *PTY) Write(data []byte) (int, error) {
	return p.ptmx.Write(data)
}

// Resize changes the PTY window size, sending SIGWINCH to the child process.
func (p *PTY) Resize(cols, rows uint16) error {
	return pty.Setsize(p.ptmx, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
}

// Close terminates the shell process and closes the PTY.
func (p *PTY) Close() error {
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		_ = p.cmd.Wait()
	}
	if p.ptmx != nil {
		return p.ptmx.Close()
	}
	return nil
}

// Fd returns the PTY master file descriptor.
func (p *PTY) Fd() uintptr {
	return p.ptmx.Fd()
}
