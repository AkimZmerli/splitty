<div align="center">

```
     ___       ___ __  __
    / __|_ __ | (_) |_| |_ _  _
    \__ \ '_ \| | |  _|  _| || |
    |___/ .__/|_|_|\__|\__|\_, |
        |_|                |__/
```

### Split pane terminal multiplexing for your Bubble Tea apps

[![Go Reference](https://pkg.go.dev/badge/github.com/AkimZmerli/splitty.svg)](https://pkg.go.dev/github.com/AkimZmerli/splitty)
[![Go Report Card](https://goreportcard.com/badge/github.com/AkimZmerli/splitty)](https://goreportcard.com/report/github.com/AkimZmerli/splitty)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/AkimZmerli/splitty)](https://golang.org)

</div>

---

**Splitty** is a standalone Go library that brings full terminal multiplexing to
[Bubble Tea](https://github.com/charmbracelet/bubbletea) applications. Split,
navigate, zoom, broadcast, and persist terminal panes -- all without leaving
your TUI.

Think of it as tmux that lives *inside* your Go program.

## Highlights

- **Binary tree layout engine** -- split any pane vertically or horizontally, infinitely deep
- **5 built-in themes** -- Default, Tokyo Night, Dracula, Nord, and Catppuccin
- **Vim-style navigation** -- `Ctrl+h/j/k/l` to jump between panes, `Tab` to cycle
- **Zoom** -- blow up any pane to fullscreen and back with `Ctrl+z`
- **Broadcast mode** -- type into *every* pane at once (hello, cluster ops)
- **Layout persistence** -- save and restore your pane arrangements as JSON
- **Layout presets** -- spin up `single`, `dev`, `triple`, or `quad` layouts in one call
- **Mouse support** -- click to focus, because sometimes a trackpad just wins
- **Fully embeddable** -- Splitty is a `tea.Model`; drop it into any Bubble Tea app

## Installation

```bash
go get github.com/AkimZmerli/splitty
```

## Quick Start

Ten lines. That's all it takes to get a working split-pane terminal.

```go
package main

import (
    "log"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/AkimZmerli/splitty"
)

func main() {
    m := splitty.New(
        splitty.WithTheme(splitty.TokyoNight),
    )
    p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Hit `Ctrl+\` to split vertically. Hit `Ctrl+-` to split horizontally. You are
now a multiplexer.

## Features

### Splitting

Every pane can be split vertically or horizontally. Under the hood, Splitty
manages a binary tree of `splitNode` and `leafNode` elements. The split ratio
defaults to 50/50 and can be resized dynamically.

```go
// Programmatic splitting
m.Split(splitty.Vertical)
m.Split(splitty.Horizontal)
```

### Navigation

Move between panes with directional focus (vim-style `h/j/k/l`) or cycle
through them sequentially with `Tab` and `Shift+Tab`. Click a pane with the
mouse to focus it instantly.

```go
m.Focus(splitty.Vertical)    // move focus along the vertical axis
m.FocusPane("pane-3")        // jump to a specific pane by ID
```

### Themes

Splitty ships with five carefully crafted themes. Each one styles pane borders,
the status bar, dividers, and zoom/broadcast indicators.

```go
splitty.New(splitty.WithTheme(splitty.Dracula))
```

| Theme | Vibe |
|-------|------|
| `DefaultTheme` | Clean ANSI 256 -- works everywhere |
| `TokyoNight` | Cool blue and muted purple |
| `Dracula` | The classic dark palette |
| `Nord` | Arctic, icy, Scandinavian calm |
| `Catppuccin` | Warm pastels on deep mocha |

You can also define your own `splitty.Theme` struct for full control over every
style.

### Presets

Skip the manual splitting and jump straight into a layout.

```go
splitty.New(splitty.WithPreset(splitty.PresetDev))
```

| Preset | Layout |
|--------|--------|
| `PresetSingle` | One full-screen pane |
| `PresetDev` | 60/40 vertical split, right side split horizontally (editor + terminal + logs) |
| `PresetTriple` | Three equal columns |
| `PresetQuad` | 2x2 grid |

Want something custom? Register your own:

```go
splitty.RegisterPreset("myLayout", func(shell string, env []string, w, h int) splitty.Node { ... })
```

### Zoom

Temporarily maximize the focused pane to the full terminal size. The rest of
your layout is preserved in the background -- unzoom to snap right back.

```go
m.Zoom()       // maximize focused pane
m.Unzoom()     // restore previous layout
m.IsZoomed()   // check zoom state
```

### Broadcast

Flip on broadcast mode and every keystroke goes to *all* panes simultaneously.
Perfect for running the same command across multiple sessions.

```go
m.SetBroadcast(true)
m.SendInput([]byte("uptime\n"))    // sent to every pane
m.IsBroadcasting()                 // check broadcast state
```

### Persistence

Save your carefully arranged layout to a JSON file and reload it later. Splitty
persists the tree structure, split ratios, working directories, and shell paths.
New PTY sessions are spawned on load.

```go
m.SaveLayout("~/.config/splitty/layout.json")
m.LoadLayout("~/.config/splitty/layout.json")
```

### Embedding

Splitty's `Manager` implements `tea.Model`, so you can embed it inside a larger
Bubble Tea application. Nest it alongside other models, wrap it in a container,
or use it as one view in a multi-page TUI.

```go
type App struct {
    splits *splitty.Manager
    // ... your other models
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    updated, cmd := a.splits.Update(msg)
    a.splits = updated.(*splitty.Manager)
    return a, cmd
}
```

## Default Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+\` | Split vertical |
| `Ctrl+-` | Split horizontal |
| `Ctrl+w` | Close pane |
| `Ctrl+h` | Focus left |
| `Ctrl+j` | Focus down |
| `Ctrl+k` | Focus up |
| `Ctrl+l` | Focus right |
| `Tab` | Next pane |
| `Shift+Tab` | Previous pane |
| `Ctrl+z` | Toggle zoom |
| `Ctrl+x` | Swap panes |
| `Ctrl+b` | Toggle broadcast |
| `Ctrl+Shift+Arrow` | Resize pane |

All keybindings are fully customizable via `WithKeyMap`:

```go
km := splitty.DefaultKeyMap()
km.SplitVertical = key.NewBinding(key.WithKeys("ctrl+d"))
splitty.New(splitty.WithKeyMap(km))
```

## Configuration

Splitty uses the functional options pattern. Mix and match to taste.

```go
m := splitty.New(
    splitty.WithShell("/bin/zsh"),              // shell for new panes
    splitty.WithTheme(splitty.Nord),            // visual theme
    splitty.WithKeyMap(myKeyMap),               // custom keybindings
    splitty.WithPreset(splitty.PresetQuad),     // initial layout
    splitty.WithMinSize(12, 4),                 // minimum pane dimensions
    splitty.WithStatusBar(true),                // bottom status bar
    splitty.WithMouse(true),                    // mouse click-to-focus
    splitty.WithEnv([]string{"TERM=xterm-256color"}), // extra env vars
    splitty.WithLogger(myLogger),               // charmbracelet/log logger
)
```

## API at a Glance

```
splitty.New(opts...)          -> *Manager (implements tea.Model)

Manager methods:
  Split(dir)                  Split the focused pane
  Close()                     Close the focused pane
  ClosePane(id)               Close a specific pane
  Focus(dir)                  Move focus directionally
  FocusPane(id)               Focus a pane by ID
  FocusedPane()               Get the focused pane
  Panes()                     List all panes
  Zoom() / Unzoom()           Maximize / restore a pane
  Swap()                      Swap focused pane with its sibling
  Resize(dir, delta)          Adjust split ratio
  SaveLayout(path)            Persist layout to JSON
  LoadLayout(path)            Restore layout from JSON
  SetBroadcast(bool)          Toggle broadcast mode
  SendInput([]byte)           Send raw input to pane(s)

Messages:
  PaneSplitMsg                Emitted when a pane is split
  PaneClosedMsg               Emitted when a pane is closed
  PaneFocusedMsg              Emitted when focus changes
  PaneOutputMsg               Emitted when a pane produces output
  PaneResizedMsg              Emitted when a pane is resized
  LayoutLoadedMsg             Emitted when a layout is loaded
```

## Examples

The [`examples/`](./examples/) directory contains runnable demos:

| Example | What it shows |
|---------|---------------|
| [`basic`](./examples/basic) | Minimal setup -- just `New()` and go |
| [`presets`](./examples/presets) | Starting with a preset layout |
| [`custom-theme`](./examples/custom-theme) | Building your own theme |
| [`custom-keys`](./examples/custom-keys) | Remapping keybindings |
| [`embedded`](./examples/embedded) | Embedding Splitty in a larger app |

## Architecture

Splitty is organized around a few core concepts:

- **Binary tree layout** -- Every split creates a `splitNode` with two children.
  Leaves hold `Pane` instances with their own PTY and virtual screen.
- **PTY management** -- Each pane spawns a real shell process via
  `creack/pty`, with full ANSI terminal emulation in the `terminal` package.
- **Bubble Tea integration** -- `Manager` implements `Init`, `Update`, and
  `View`. All pane I/O flows through the Bubble Tea message loop.

For a deeper dive, see [`ARCHITECTURE.md`](./ARCHITECTURE.md).

## Contributing

Contributions are welcome and appreciated! Whether it is a bug report, a new
theme, a feature idea, or a pull request -- all of it helps.

1. Fork the repo
2. Create a feature branch (`git checkout -b my-feature`)
3. Commit your changes
4. Open a pull request

Please keep PRs focused and include tests where possible.

## License

Splitty is released under the [MIT License](./LICENSE).

---

<div align="center">

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and an
unreasonable fondness for split panes.

</div>
