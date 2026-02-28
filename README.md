<div align="center">

```
     ___       ___ __  __
    / __|_ __ | (_) |_| |_ _  _
    \__ \ '_ \| | |  _|  _| || |
    |___/ .__/|_|_|\__|\__|\_, |
        |_|                |__/
```

### The terminal multiplexer built for agentic coding

[![Go Reference](https://pkg.go.dev/badge/github.com/AkimZmerli/splitty.svg)](https://pkg.go.dev/github.com/AkimZmerli/splitty)
[![Go Report Card](https://goreportcard.com/badge/github.com/AkimZmerli/splitty)](https://goreportcard.com/report/github.com/AkimZmerli/splitty)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/AkimZmerli/splitty)](https://golang.org)

</div>

---

**Splitty** is a Go terminal multiplexer built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) with one guiding principle: **your left hand never leaves home row.**

Navigate panes with `Ctrl+W/A/S/D` -- the same muscle memory you use in games. When you're running 3 AI agents in parallel and need to check on each one, movement should feel instant and familiar. Not `Ctrl+B` then arrow keys. Not clicking tabs. Just WASD.

Built for the era of agentic coding, where developers orchestrate multiple AI sessions simultaneously and **observability across panes is everything.**

## Why Splitty?

Agentic coding is changing how developers work. You're no longer running one terminal session -- you're watching a code agent, a test runner, and a review agent all working at the same time. You need to:

- **See everything at once** -- split panes with real terminal emulation
- **Navigate instantly** -- WASD feels like a game, not a chore
- **Stay in flow** -- left hand navigates, right hand stays on the mouse or keys
- **Monitor activity** -- know which agent is working and which is waiting

Existing multiplexers (tmux, screen) were built for shell management. Splitty is built for **agent orchestration**.

## Quick Start

```bash
go get github.com/AkimZmerli/splitty
```

Ten lines to a working multi-pane terminal:

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

Hit `Ctrl+V` to split vertically. Hit `Ctrl+H` to split horizontally. Navigate with `Ctrl+W/A/S/D`. You're multiplexing.

## Keybindings

The keybindings are designed around left-hand ergonomics. Your hand rests on WASD for navigation, with splitting and management keys within natural reach.

### Navigation (WASD)

| Key | Action |
|-----|--------|
| `Ctrl+W` | Focus up |
| `Ctrl+A` | Focus left |
| `Ctrl+S` | Focus down |
| `Ctrl+D` | Focus right |
| `Tab` | Next pane |
| `Shift+Tab` | Previous pane |

### Pane Management

| Key | Action |
|-----|--------|
| `Ctrl+V` | Split vertical |
| `Ctrl+H` | Split horizontal |
| `Ctrl+C` | Close pane |
| `Ctrl+Z` | Toggle zoom |
| `Ctrl+X` | Swap panes |
| `Ctrl+B` | Toggle broadcast |
| `Ctrl+Shift+Arrow` | Resize pane |

### Scrollback

| Key | Action |
|-----|--------|
| `Ctrl+K` | Scroll up |
| `Ctrl+J` | Scroll down |
| `Ctrl+U` | Page up |
| `Ctrl+N` | Page down |
| `Ctrl+Home` | Scroll to top |
| `Ctrl+End` | Scroll to bottom |

All keybindings are fully customizable:

```go
km := splitty.DefaultKeyMap()
km.FocusUp = key.NewBinding(key.WithKeys("ctrl+w"))
splitty.New(splitty.WithKeyMap(km))
```

## Features

### WASD Navigation

The core design decision. Navigate between panes with `Ctrl+W/A/S/D` using the same spatial awareness you use in games. Your left hand stays planted, your right hand stays free. When you're watching 4 agent sessions, switching between them should be reflexive, not a sequence of prefix keys.

### Splitting

Every pane can be split vertically or horizontally. Under the hood, Splitty manages a binary tree of `splitNode` and `leafNode` elements. The split ratio defaults to 50/50 and can be resized dynamically.

```go
m.Split(splitty.Vertical)
m.Split(splitty.Horizontal)
```

### Themes

Five built-in themes. Each one styles pane borders, the status bar, dividers, and zoom/broadcast indicators.

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

Define your own `splitty.Theme` struct for full control over every style.

### Presets

Skip the manual splitting and jump straight into a layout.

```go
splitty.New(splitty.WithPreset(splitty.PresetDev))
```

| Preset | Layout |
|--------|--------|
| `PresetSingle` | One full-screen pane |
| `PresetDev` | 60/40 vertical split, right side split horizontally |
| `PresetTriple` | Three equal columns |
| `PresetQuad` | 2x2 grid |

Register your own:

```go
splitty.RegisterPreset("agentPair", func(shell string, env []string, w, h int) splitty.Node { ... })
```

### Scrollback Buffer

Configurable scrollback history per pane. Border color changes to indicate when you're viewing history. Auto-scroll resumes when you return to the bottom.

```go
splitty.New(splitty.WithScrollbackLines(2000))
```

### Zoom

Temporarily maximize the focused pane to the full terminal size. The rest of your layout is preserved in the background.

```go
m.Zoom()
m.Unzoom()
```

### Broadcast

Broadcast mode sends every keystroke to all panes simultaneously. Run the same command across multiple agent sessions at once.

```go
m.SetBroadcast(true)
m.SendInput([]byte("continue\n"))
```

### Persistence

Save your layout to JSON and reload it later. Splitty persists the tree structure, split ratios, working directories, and shell paths.

```go
m.SaveLayout("~/.config/splitty/layout.json")
m.LoadLayout("~/.config/splitty/layout.json")
```

### Mouse Support

Click to focus panes. Scroll wheel for scrollback navigation (3 lines per notch). Because sometimes a trackpad just wins.

### Embedding

Splitty's `Manager` implements `tea.Model`, so you can embed it inside any Bubble Tea application.

```go
type App struct {
    splits *splitty.Manager
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    updated, cmd := a.splits.Update(msg)
    a.splits = updated.(*splitty.Manager)
    return a, cmd
}
```

## Configuration

Functional options pattern. Mix and match.

```go
m := splitty.New(
    splitty.WithShell("/bin/zsh"),
    splitty.WithTheme(splitty.Nord),
    splitty.WithKeyMap(myKeyMap),
    splitty.WithPreset(splitty.PresetQuad),
    splitty.WithScrollbackLines(1000),
    splitty.WithMinSize(12, 4),
    splitty.WithStatusBar(true),
    splitty.WithMouse(true),
    splitty.WithEnv([]string{"TERM=xterm-256color"}),
    splitty.WithLogger(myLogger),
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

- **Binary tree layout** -- Every split creates a `splitNode` with two children. Leaves hold `Pane` instances with their own PTY and virtual screen.
- **PTY management** -- Each pane spawns a real shell process via `creack/pty`, with full ANSI terminal emulation in the `terminal` package.
- **Bubble Tea integration** -- `Manager` implements `Init`, `Update`, and `View`. All pane I/O flows through the Bubble Tea message loop.

For a deeper dive, see [`ARCHITECTURE.md`](./ARCHITECTURE.md).

## Roadmap

See [`ROADMAP.md`](./ROADMAP.md) for planned features including agent observability, activity indicators, output search, and framework integrations.

## Contributing

Contributions are welcome. Whether it is a bug report, a new theme, a feature idea, or a pull request -- all of it helps.

1. Fork the repo
2. Create a feature branch (`git checkout -b my-feature`)
3. Commit your changes
4. Open a pull request

Please keep PRs focused and include tests where possible.

## License

Splitty is released under the [MIT License](./LICENSE).

---

<div align="center">

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for developers who run agents like they play games.

</div>
