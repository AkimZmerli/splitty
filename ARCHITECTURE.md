# Splitty Architecture

This document describes the internal architecture of Splitty, a Go library that
provides split pane terminal multiplexing for
[Bubble Tea](https://github.com/charmbracelet/bubbletea) applications.

```
import "github.com/AkimZmerli/splitty"
```

---

## 1. Overview

Splitty is a single-package Go library (plus one internal sub-package) that turns
any Bubble Tea program into a split pane terminal multiplexer. It manages multiple
pseudo-terminal sessions arranged in a binary tree layout, with a built-in ANSI
terminal emulator that translates raw PTY output into renderable content for Bubble
Tea's `View()` cycle.

The library is designed around a single entry point: `splitty.New()` returns a
`*Manager` that implements `tea.Model`. The consumer passes it to
`tea.NewProgram()` and gets a working terminal multiplexer with no additional
wiring.

**Single-package rationale.** Splitty is consumed as a library, not a framework.
All public types -- `Manager`, `Theme`, `KeyMap`, `Option`, `Direction`, and the
various `tea.Msg` types -- live in the root `splitty` package. This keeps the
import path flat (`import "github.com/AkimZmerli/splitty"`) and avoids forcing
consumers to juggle multiple sub-package imports for common operations. The only
sub-package is `terminal/`, which encapsulates the virtual terminal emulator and
PTY wrapper as a bounded internal domain with no upward dependencies.

---

## 2. File Organization

Splitty uses a **vertical slice** architecture. Each file owns a complete feature
from public API through internal logic, rather than organizing by technical layer
(e.g., separating "models" from "controllers"). Shared infrastructure lives in a
small set of core files.

### Feature Slices

Each feature slice contains all the methods for one capability of the Manager:

| File            | Feature                     | Key Types / Functions                                      |
|-----------------|-----------------------------|------------------------------------------------------------|
| `split.go`      | Split and close panes       | `Split()`, `Close()`, `ClosePane()`, `replaceNode()`       |
| `navigate.go`   | Focus navigation            | `Focus()`, `FocusPane()`, `focusDirection()`, `cycleFocus()`|
| `resize.go`     | Resize and layout calc      | `Resize()`, `layoutAll()`, `layoutNode()`                  |
| `zoom.go`       | Zoom/maximize and swap      | `Zoom()`, `Unzoom()`, `toggleZoom()`, `Swap()`             |
| `broadcast.go`  | Broadcast input to all panes| `SetBroadcast()`, `SendInput()`                            |
| `persist.go`    | Save/load layouts as JSON   | `SaveLayout()`, `LoadLayout()`, `serializeNode()`          |
| `presets.go`    | Named layout presets        | `RegisterPreset()`, preset builders, `presetRegistry`      |

### Core Infrastructure

| File            | Role                                                                      |
|-----------------|---------------------------------------------------------------------------|
| `splitty.go`    | `Manager` struct definition, `New()`, `Init()`, `Update()`, `View()`      |
| `node.go`       | Binary tree types: `node` interface, `leafNode`, `splitNode`              |
| `pane.go`       | `Pane` type with PTY lifecycle (`start`, `resize`, `write`, `close`)      |
| `messages.go`   | Public `tea.Msg` types and the `Direction` enum                           |
| `keys.go`       | `KeyMap` struct and `DefaultKeyMap()`                                     |
| `theme.go`      | `Theme` struct and 5 built-in themes                                      |
| `options.go`    | `Option` type (functional options) and all `With*()` constructors         |
| `logger.go`     | Nil-safe `charmbracelet/log` wrapper                                      |
| `doc.go`        | Package-level godoc                                                       |

### Terminal Sub-Package

The `terminal/` package is a **bounded domain** that has no imports from the
parent `splitty` package. It provides two things: a virtual terminal emulator
(`Screen`) and a PTY wrapper (`PTY`).

| File              | Contents                                                      |
|-------------------|---------------------------------------------------------------|
| `terminal/screen.go` | `Screen` struct, ANSI state machine parser, cell grid, `Render()` |
| `terminal/cell.go`   | `Cell`, `Style`, `Color` types, `ToANSI()` conversion        |
| `terminal/cursor.go` | `Cursor` state (row, col, visibility, current style)          |
| `terminal/pty.go`    | `PTY` wrapper over `creack/pty`                               |

---

## 3. Binary Tree Model

All layout state is encoded as a binary tree rooted at `Manager.root`. The tree
has two node types, defined by the `node` interface:

- **`leafNode`** -- a terminal leaf containing a single `*Pane`.
- **`splitNode`** -- an internal container with a `Direction` (Vertical or
  Horizontal), a `ratio` (float64, 0.0 to 1.0), and two children (`first`,
  `second`).

```
              splitNode
              dir=Vertical
              ratio=0.5
             /          \
        leafNode      splitNode
        [pane-1]      dir=Horizontal
                      ratio=0.6
                     /          \
                leafNode      leafNode
                [pane-2]      [pane-3]
```

This tree represents a layout where the screen is split vertically into two
halves. The left half contains pane-1. The right half is split horizontally into
a top 60% (pane-2) and bottom 40% (pane-3):

```
+------------+------------+
|            |            |
|            |   pane-2   |
|   pane-1   |   (60%)    |
|            |------------|
|            |   pane-3   |
|            |   (40%)    |
+------------+------------+
     50%           50%
```

### Split Algorithm

`Split(dir)` in `split.go`:

1. Find the focused leaf by ID.
2. Create a new `Pane` and wrap it in a `leafNode`.
3. Create a `splitNode` with `ratio=0.5`, the old leaf as `first`, the new leaf
   as `second`.
4. Replace the focused leaf in the tree with the new `splitNode` (via
   `replaceNode()`, which walks up to find the parent and calls
   `replaceChild()`).
5. Call `layoutAll()` to recalculate all positions and sizes.
6. Start the new pane's PTY and return a `readPaneCmd` to begin reading output.

### Close Algorithm

`ClosePane(id)` in `split.go`:

1. Close the pane's PTY.
2. Find the parent `splitNode` via `findParent()`.
3. Get the sibling of the closed leaf.
4. Find the grandparent `splitNode` via `findParentSplit()`.
5. Replace the parent with the sibling in the grandparent -- this is the
   **collapse** step. The parent `splitNode` is removed from the tree and the
   sibling is promoted up.
6. If the parent was the root, the sibling becomes the new root.
7. Call `layoutAll()` to recalculate layout.

```
Before close(pane-3):          After:

    splitNode                  splitNode
   /         \                /         \
leafNode   splitNode    leafNode      leafNode
[pane-1]  /         \   [pane-1]      [pane-2]
     leafNode  leafNode
     [pane-2]  [pane-3]
```

### Tree Traversal

The `node` interface provides three traversal methods:

- `findLeaf(id) (*leafNode, []node)` -- returns the matching leaf and the full
  path from root to leaf (used for directional navigation and resize).
- `leaves() []*leafNode` -- returns all leaves in left-to-right, top-to-bottom
  order (in-order traversal of the binary tree).
- `clone() node` -- deep-copies the tree structure (shallow-copies pane
  pointers).

---

## 4. Terminal Emulation Pipeline

Each pane runs an independent shell process via a PTY. The output pipeline is:

```
Shell process
     |
     | (raw bytes: text + ANSI escape sequences)
     v
PTY.Read(buf)           terminal/pty.go
     |
     v
Screen.Write(buf)       terminal/screen.go
     |
     | (byte-by-byte state machine)
     v
ANSI parser             terminal/screen.go
     |
     | (updates cell grid + cursor)
     v
Cell grid               [][]Cell
     |
     | (on View() call)
     v
Screen.Render()         terminal/screen.go
     |
     | (diff-style ANSI output with SGR optimization)
     v
ANSI string -> lipgloss border -> View()
```

### ANSI Parser State Machine

The parser in `Screen.processByte()` implements a five-state machine:

| State          | Entered On  | Handles                                              |
|----------------|-------------|------------------------------------------------------|
| `stateGround`  | (default)   | Printable chars, CR, LF, BS, TAB, BEL, ESC           |
| `stateEscape`  | ESC (0x1B)  | `[` -> CSI, `]` -> OSC, `(` -> charset, `7`/`8`, `D`/`M`/`c` |
| `stateCSI`     | ESC `[`     | Params (`0-9`, `;`), `?` private prefix, final byte dispatch |
| `stateOSC`     | ESC `]`     | Accumulates until BEL (0x07) or ESC, sets window title|
| `stateCharset` | ESC `(`     | Consumes one byte and returns to ground               |

The CSI dispatcher (`dispatchCSI`) handles 20+ commands including cursor
movement (CUU/CUD/CUF/CUB/CUP), erase operations (ED/EL/ECH), line
manipulation (IL/DL), scrolling (SU/SD), scroll regions (DECSTBM), SGR styling,
and DEC private modes (cursor visibility, alternate screen buffer).

### SGR Handling

`handleSGR()` processes Select Graphic Rendition parameters, supporting:

- Attributes: bold, dim, italic, underline, blink, reverse, hidden, strikethrough
- 16-color foreground/background (SGR 30-37, 40-47, 90-97, 100-107)
- 256-color palette (SGR 38;5;N / 48;5;N)
- 24-bit true color (SGR 38;2;R;G;B / 48;2;R;G;B)

### Render Output

`Screen.Render()` iterates the cell grid row by row, column by column, emitting
ANSI SGR sequences only when the style changes from the previous cell. This
minimizes the output size. The result is a string containing ANSI-styled text
with newline separators between rows, terminated by a reset (`ESC[0m`).

---

## 5. Message Flow

Splitty follows the standard Bubble Tea `Init() -> Update(msg) -> View()` loop.
All state mutations happen in `Update()`.

### Message Dispatch

`Manager.Update()` handles four message types:

```
tea.WindowSizeMsg  ->  handleWindowSize()  ->  initLayout() or layoutAll()
tea.KeyMsg         ->  handleKey()         ->  split/close/navigate/resize/zoom/broadcast/forward
tea.MouseMsg       ->  handleMouse()       ->  click-to-focus
PaneOutputMsg      ->  readPaneCmd()       ->  schedule next read
PaneClosedMsg      ->  closePane()         ->  remove pane from tree
```

### The Self-Re-Invoking Read Pattern

PTY reading uses a continuation pattern common in Bubble Tea applications. This
is the critical data flow that keeps terminal output streaming:

```
1. initLayout() starts each pane:
      cmds = append(cmds, m.readPaneCmd(pane.ID))

2. readPaneCmd(id) returns a tea.Cmd (a closure):
      func() tea.Msg {
          n, err := leaf.pane.pty.Read(buf)     // blocks until data
          leaf.pane.screen.Write(buf[:n])        // parse into cell grid
          return PaneOutputMsg{PaneID: id, Data: buf[:n]}
      }

3. Bubble Tea runtime executes the Cmd in a goroutine.
   When it returns, Update() receives PaneOutputMsg.

4. Update() handles PaneOutputMsg:
      return m, m.readPaneCmd(msg.PaneID)       // schedule next read

5. Goto step 2.
```

This creates a perpetual read loop for each pane. The loop terminates when
`pty.Read()` returns an error (process exit), which produces a `PaneClosedMsg`
instead of `PaneOutputMsg`.

### Key Input Flow

Non-intercepted keystrokes flow from the terminal into the pane's PTY:

```
tea.KeyMsg
     |
     v
handleKey() -- check against KeyMap bindings
     |
     | (no match: default case)
     v
keyToBytes(msg) -- convert tea.KeyMsg to raw bytes
     |
     v
SendInput(data) -- write to focused pane's PTY (or all panes if broadcasting)
     |
     v
PTY.Write(data) -- shell process receives input
```

The `keyToBytes()` function handles translation of Bubble Tea's key
representation back to raw terminal bytes: special keys (Enter, Tab, arrow keys)
are mapped to their escape sequences, control keys to their byte values
(ctrl+a = 0x01), and printable runes to their UTF-8 encoding.

---

## 6. Concurrency Model

Splitty operates within Bubble Tea's concurrency model, which serializes all
state access through the `Update()` function. The only concurrent access point
is the virtual terminal screen.

### Bubble Tea's Goroutine Model

Bubble Tea runs `tea.Cmd` functions in goroutines but delivers the resulting
`tea.Msg` back to `Update()` sequentially. This means:

- `readPaneCmd` runs in a goroutine (it blocks on `pty.Read()`).
- The returned `PaneOutputMsg` is delivered to `Update()` on the main loop.
- Tree mutations (split, close, navigate, resize) always happen on the main loop.

### Screen Mutex

The one exception is `Screen.Write()` -- called inside `readPaneCmd` goroutines
-- and `Screen.Render()` -- called from `View()` on the main loop. These can
overlap, so `Screen` uses a `sync.RWMutex`:

| Method              | Lock Type  | Called From                |
|---------------------|------------|---------------------------|
| `Screen.Write()`    | `mu.Lock()`   | `readPaneCmd` goroutine   |
| `Screen.Resize()`   | `mu.Lock()`   | `layoutAll()` on main loop|
| `Screen.Render()`   | `mu.RLock()`  | `View()` on main loop     |
| `Screen.Width()`    | `mu.RLock()`  | Any                       |
| `Screen.Height()`   | `mu.RLock()`  | Any                       |
| `Screen.Title()`    | `mu.RLock()`  | Any                       |

This is safe because:

- Multiple `Render()` calls can proceed concurrently (read lock).
- A `Write()` call blocks until any concurrent `Render()` finishes, and vice
  versa.
- `Resize()` takes a write lock, ensuring exclusive access during dimension
  changes.

### Pane ID Generation

Pane IDs are generated via `sync/atomic.AddUint64` on a package-level counter,
making ID generation safe from any goroutine without additional locking.

---

## 7. Rendering Pipeline

`View()` is called by Bubble Tea on every frame. It produces the final string
for the terminal.

### Normal Mode

```
View()
  |
  +-- renderNode(root, width, height)     recursive tree walk
  |     |
  |     +-- leafNode: renderPane(pane, w, h)
  |     |     |
  |     |     +-- pane.render()            -> Screen.Render() -> ANSI string
  |     |     +-- theme.BorderActive or    -> lipgloss.Style wraps content
  |     |         theme.BorderInactive        with rounded border
  |     |
  |     +-- splitNode: renderSplit(sn, w, h)
  |           |
  |           +-- Vertical: compute leftW/rightW from ratio
  |           |     renderNode(first, leftW, h)
  |           |     renderNode(second, rightW, h)
  |           |     lipgloss.JoinHorizontal(Top, left, right)
  |           |
  |           +-- Horizontal: compute topH/bottomH from ratio
  |                 renderNode(first, w, topH)
  |                 renderNode(second, w, bottomH)
  |                 lipgloss.JoinVertical(Left, top, bottom)
  |
  +-- renderStatusBar()                   bottom bar with hints + pane count
  |
  +-- lipgloss.JoinVertical(Left, content, bar)
```

### Zoomed Mode

When a pane is zoomed, `View()` skips the tree walk entirely:

```
View()
  +-- renderZoomed()
  |     +-- renderPane(zoomedPane, fullWidth, fullHeight)
  +-- renderStatusBar()
```

The zoomed pane is resized to the full terminal dimensions (minus border and
status bar). All other panes remain in the tree but are not rendered.

### Layout Calculation

`layoutAll()` (in `resize.go`) is separate from rendering. It recursively walks
the tree via `layoutNode()`, assigning `X`, `Y`, `Width`, and `Height` to each
pane and calling `pane.resize()` to update both the `Screen` dimensions and the
PTY window size (which sends `SIGWINCH` to the shell process).

This separation means layout is recalculated only on structural changes (split,
close, resize, window resize), not on every frame.

---

## 8. Extension Points

### Custom Themes

Create a `Theme` struct with lipgloss styles for active/inactive borders, the
status bar, divider appearance, and indicator strings:

```go
myTheme := splitty.Theme{
    BorderActive:       lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("#FF0000")),
    BorderInactive:     lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#555555")),
    StatusBar:          lipgloss.NewStyle().Background(lipgloss.Color("#111111")).Foreground(lipgloss.Color("#FFFFFF")),
    ZoomIndicator:      "[ZOOM]",
    BroadcastIndicator: "[BROADCAST]",
    // ...
}
m := splitty.New(splitty.WithTheme(myTheme))
```

Five built-in themes are provided: `DefaultTheme`, `TokyoNight`, `Dracula`,
`Nord`, and `Catppuccin`.

### Custom Keybindings

Replace the entire `KeyMap` or modify individual bindings:

```go
km := splitty.DefaultKeyMap()
km.SplitVertical = key.NewBinding(key.WithKeys("ctrl+d"))
m := splitty.New(splitty.WithKeyMap(km))
```

### Custom Presets

Register a named preset that builds an arbitrary tree:

```go
splitty.RegisterPreset("monitoring", func(shell string, env []string, w, h int) splitty.node {
    // Build and return a custom node tree
})
m := splitty.New(splitty.WithPreset("monitoring"))
```

Built-in presets: `PresetSingle` (1 pane), `PresetDev` (60/40 vertical with
horizontal sub-split), `PresetTriple` (3 columns), `PresetQuad` (2x2 grid).

### Embedding in Larger Applications

Since `Manager` implements `tea.Model`, it can be embedded as a component in a
larger Bubble Tea application. The host application wraps `Manager` in its own
model and forwards messages:

```go
type App struct {
    splitty *splitty.Manager
    // ... other state
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle app-level keys first, then delegate:
    updated, cmd := a.splitty.Update(msg)
    a.splitty = updated.(*splitty.Manager)
    return a, cmd
}
```

The `Panes()`, `FocusedPane()`, `FocusPane()`, `IsZoomed()`, and
`IsBroadcasting()` methods provide read access to layout state for the host
application.

---

## 9. Design Decisions

### Why a Binary Tree (Not a List or Grid)

A binary tree naturally models the recursive split operation. When a user splits
a pane, the leaf becomes an internal node with two children -- this is a
one-step tree insertion. Closing a pane is a one-step collapse (replace parent
with sibling). A flat list would require tracking adjacency relationships
separately, and a grid would not support uneven splits or nested subdivisions.

The binary tree also maps directly to the rendering algorithm: each `splitNode`
is a `JoinHorizontal` or `JoinVertical` operation, making the render function a
simple recursive descent.

### Why a Self-Rolled ANSI Parser (Not a Library Like midterm)

Splitty includes its own ANSI terminal emulator in the `terminal/` package for
three reasons:

1. **Control over the cell grid.** The emulator needs to produce a `[][]Cell`
   grid that can be rendered into ANSI strings suitable for embedding inside
   lipgloss-styled borders. Third-party terminal emulators typically target a
   different output format (e.g., a full terminal screen or a widget), making
   integration with Bubble Tea's string-based `View()` awkward.

2. **Minimal scope.** Splitty does not need a full VT100/xterm emulator. It
   needs enough to handle shells and common CLI tools (vim, htop, etc.). The
   parser handles CSI sequences, SGR attributes (including true color), scroll
   regions, cursor save/restore, alternate screen buffer, and OSC title
   sequences. This covers the vast majority of real-world terminal output.

3. **Dependency control.** Keeping the emulator in-tree means no version
   conflicts or unexpected behavior changes from upstream libraries.

### Why Functional Options

`New()` uses the functional options pattern (`WithShell()`, `WithTheme()`, etc.)
for several reasons:

- All options have sensible defaults. The zero-config `splitty.New()` works out
  of the box.
- New options can be added without breaking existing callers.
- Options are composable: `splitty.New(splitty.WithTheme(t), splitty.WithPreset(p))`.
- The pattern is idiomatic in the Charm ecosystem (`tea.NewProgram` uses the
  same approach).

### Why a Single Package

Splitting into multiple packages (e.g., `splitty/layout`, `splitty/terminal`,
`splitty/theme`) would force consumers to import multiple paths for common
operations. Since Splitty is a library, not a framework, the simplest consumer
experience is a single import. The `terminal/` sub-package exists only because
it is a genuinely independent domain with no upward dependencies -- the ANSI
parser and PTY wrapper know nothing about Bubble Tea, pane management, or layout
trees.
