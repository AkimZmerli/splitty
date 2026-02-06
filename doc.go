// Package splitty provides split pane terminal multiplexing
// for Bubble Tea applications.
//
// Splitty manages multiple terminal panes in a configurable layout using
// a binary tree structure. Each pane runs an independent shell session via
// PTY, with full ANSI terminal emulation.
//
// Basic usage:
//
//	m := splitty.New()
//	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
//	if _, err := p.Run(); err != nil {
//	    log.Fatal(err)
//	}
//
// Custom configuration:
//
//	m := splitty.New(
//	    splitty.WithShell("/bin/zsh"),
//	    splitty.WithTheme(splitty.TokyoNight),
//	    splitty.WithPreset(splitty.PresetDev),
//	)
package splitty
