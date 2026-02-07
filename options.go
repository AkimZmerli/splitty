package splitty

import "github.com/charmbracelet/log"

// Option configures the Manager.
type Option func(*Manager)

// WithShell sets the shell command for new panes.
// Defaults to $SHELL or /bin/sh.
func WithShell(shell string) Option {
	return func(m *Manager) {
		m.shell = shell
	}
}

// WithTheme sets the visual theme for pane borders and status bar.
func WithTheme(theme Theme) Option {
	return func(m *Manager) {
		m.theme = theme
	}
}

// WithKeyMap sets custom keybindings for split pane operations.
func WithKeyMap(km KeyMap) Option {
	return func(m *Manager) {
		m.keyMap = km
	}
}

// WithMinSize sets the minimum pane dimensions in columns and rows.
// Splits that would result in panes smaller than this are rejected.
// Default: 10 columns, 3 rows.
func WithMinSize(width, height int) Option {
	return func(m *Manager) {
		m.minWidth = width
		m.minHeight = height
	}
}

// WithStatusBar enables or disables the bottom status bar.
// Default: true.
func WithStatusBar(enabled bool) Option {
	return func(m *Manager) {
		m.statusBar = enabled
	}
}

// WithMouse enables or disables mouse support for focus and resize.
// Default: true.
func WithMouse(enabled bool) Option {
	return func(m *Manager) {
		m.mouse = enabled
	}
}

// WithPreset initializes the layout with a named preset.
// Available presets: PresetSingle, PresetDev, PresetTriple, PresetQuad.
func WithPreset(name string) Option {
	return func(m *Manager) {
		m.presetName = name
	}
}

// WithEnv sets additional environment variables for PTY sessions.
// These are appended to the current environment.
func WithEnv(env []string) Option {
	return func(m *Manager) {
		m.env = env
	}
}

// WithLogger sets a charmbracelet/log logger for debug output.
// By default, logging is disabled.
func WithLogger(l *log.Logger) Option {
	return func(m *Manager) {
		m.log = newLogger(l)
	}
}

// WithScrollbackLines sets the maximum scrollback history per pane.
// Default: 1000 lines. Set to 0 to disable scrollback.
func WithScrollbackLines(lines int) Option {
	return func(m *Manager) {
		m.scrollbackLines = lines
	}
}

// WithScrollSpeed sets the number of lines scrolled per mouse wheel notch.
// Default: 3 lines.
func WithScrollSpeed(lines int) Option {
	return func(m *Manager) {
		if lines < 1 {
			lines = 1
		}
		m.scrollSpeed = lines
	}
}
