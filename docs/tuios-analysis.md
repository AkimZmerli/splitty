# Tuios vs Splitty: Architecture Analysis

Research comparing [tuios](https://github.com/Gaurav-Gosain/tuios) terminal multiplexer patterns with Splitty's current implementation, focused on scrolling, copying, mouse interactions, and rendering.

---

## 1. Rendering Pipeline

### Tuios Approach

Tuios uses a multi-level caching + compositing system built on bubbletea v2 + lipgloss:

- **Canvas/Layer compositing** — Each window is an independent `lipgloss.Layer` with X/Y/Z positioning. The compositor handles overlap and clipping (`render.go`).
- **3-tier dirty flags** — `Dirty`, `ContentDirty`, `PositionDirty` per window. Position changes reuse cached content strings. Content changes regenerate terminal text but reuse layer positioning.
- **Style batching** — Groups consecutive cells with identical attributes into a single styled string via `batchBuilder`/`flushBatch` (`render_terminal.go:174`). Avoids per-cell ANSI escape generation.
- **Style cache** — Hash-based `StyleCache` (`stylecache.go`) maps `(cell attrs, fg, bg, cursor state) → lipgloss.Style` with 1024-entry LRU. Uses `maphash` for fast hashing of cell attributes.
- **Object pooling** — `sync.Pool` for `strings.Builder`, `[]byte`, `lipgloss.Style`, `HighlightGrid`, layer slices (`pool.go`). Significant GC pressure reduction.
- **Tick-based frame rate** — 60 FPS normal, 30 FPS during interactions (`update.go:224-235`). Frame skipping when no content changed and no animations active.
- **Resize indicator** — Shows "Resizing... WxH" placeholder during drag instead of re-rendering terminal content.

### Splitty Current State

- Full re-render through tcell on every draw cycle
- No content caching per pane
- No style batching — each cell is styled individually
- No dirty tracking — everything re-renders
- No frame skipping
- O(width × height × panes) per frame

---

## 2. Scrollback Buffer

### Tuios (`internal/vt/scrollback.go`)

```go
type Scrollback struct {
    lines           []uv.Line  // Pre-allocated ring buffer
    maxLines        int        // Default: 10,000
    head, tail      int        // Ring buffer pointers
    full            bool
    lastWidthCaptured int      // For resize reflow
    softWrapped     []bool     // Track soft-wrapped vs hard-wrapped lines
}
```

Key differences from Splitty:
- **Soft-wrap tracking** — `softWrapped []bool` array tracks which lines are hard newlines vs soft wraps, enabling future reflow on resize
- **Default 10,000 lines** (vs Splitty's 1,000)
- `SetMaxLines()` dynamically resizes, keeping most recent lines
- `Reflow()` stub for terminal resize rewrapping (architected but not yet implemented)
- Uses `uv.Line` (slice of cells from charmbracelet/ultraviolet) directly

### Splitty (`terminal/screen.go`)

Ring buffer works well but lacks soft-wrap awareness, so resize breaks scrollback line alignment. Index math (`oldestIdx`) is recalculated per-cell during render instead of being cached once per render pass.

---

## 3. Copy Mode / Text Selection

### Tuios: Full Vim-Style Copy Mode

Files: `internal/input/copymode_handlers.go`, `copymode_motion.go`, `copymode_visual.go`, `copymode_search.go`, `copymode_char_search.go`

**State machine** with four states:
- `CopyModeNormal` — Navigation
- `CopyModeSearch` — Typing a search query
- `CopyModeVisualChar` — Character-wise visual selection
- `CopyModeVisualLine` — Line-wise visual selection

**Navigation motions:**
| Key | Action |
|-----|--------|
| `h/j/k/l` | Character movement |
| `w/W/b/B/e/E` | Word movement (vim word vs WORD) |
| `0/^/$` | Line start / first non-blank / line end |
| `gg/G` | Buffer top / bottom |
| `Ctrl+U/D` | Half-page up/down |
| `Ctrl+B/F` | Full page up/down |
| `H/M/L` | Screen top / middle / bottom |
| `{/}` | Paragraph up/down |
| `%` | Matching bracket |
| `f/F/t/T + char` | Find character on line |
| `;/,` | Repeat/reverse char search |
| `n/N` | Next/prev search match |

**Count prefixes:** `10j` moves down 10 lines, `5w` moves 5 words forward.

**Search:** `/` forward, `?` backward, with highlighted matches across scrollback + screen. UTF-8 aware with byte-to-column position conversion.

**Visual selection:** `v` for char, `V` for line. Mouse drag also enters visual mode. `y` or `c` yanks to system clipboard via `tea.SetClipboard()`.

**Auto-enter/exit:** Mouse wheel up auto-enters copy mode. Scrolling down to bottom auto-exits.

### Splitty Current State

Only basic scrollback navigation (ctrl+k/j/u/n). No selection, no copy, no search within scrollback, no vim motions.

---

## 4. Mouse Interactions

### Tuios (`internal/input/mouse.go`)

#### Event Forwarding to Inner Terminals
Detects if inner terminal has mouse mode enabled (`HasMouseMode()`) or is in alt screen. Translates absolute coordinates to terminal-relative coordinates and forwards via `EncodeMouseEvent()` → PTY:

```go
termX := screenX - window.X - 1  // Account for left border
termY := screenY - window.Y - 1  // Account for top border

adjustedMouse := uv.MouseClickEvent{
    X: termX, Y: termY,
    Button: uv.MouseButton(mouse.Button),
    Mod:    uv.KeyMod(mouse.Mod),
}
sendMouseClickToWindow(focusedWindow, adjustedMouse)
```

Supports both daemon mode (encode to ANSI escape sequences → PTY) and local mode (direct VT emulator call).

#### Hit Testing
Z-order aware `findClickedWindow()` iterates all windows, finds topmost (highest Z) containing the click point.

#### Double/Triple Click
Tracks click timing + position with 500ms threshold:
- Single click: character selection start
- Double click: word selection (via `isWordChar` boundary detection)
- Triple click: line selection

#### Drag Operations
- Left drag: window move (with edge snapping zones)
- Right drag: corner-based resize (detects quadrant)
- Copy mode drag: visual selection

#### Smart Resize During Drag
```go
// Visual resize only during drag (no PTY resize)
focusedWindow.ResizeVisual(newWidth, newHeight)
o.PendingResizes[focusedWindow.ID] = [2]int{newWidth, newHeight}

// On mouse release, apply all deferred PTY resizes
for i := range o.Windows {
    if dims, exists := o.PendingResizes[o.Windows[i].ID]; exists {
        o.Windows[i].Resize(dims[0], dims[1])
    }
}
```

150ms delay after resize release for shell prompt to redraw via SIGWINCH.

#### Mouse Wheel
- In terminal mode with mouse-aware app: forwards wheel events to inner terminal
- In terminal mode without mouse app: enters copy mode on wheel up, scrolls 3 lines per notch
- In copy mode: scrolls via `MoveUp`/`MoveDown`, exits when reaching bottom

### Splitty Current State

- Mouse wheel scrolls 3 lines (hard-coded)
- Left click focuses pane
- Right click shows context menu
- No mouse forwarding to inner terminals
- No text selection via mouse
- No drag resize

---

## 5. VT Emulator Comparison

### Tuios (`internal/vt/`)

Vendored fork of charmbracelet/ultraviolet with:
- **Callbacks** for AltScreen, CursorStyle, Title changes — app state stays in sync
- **`Touched()` method** returns modified lines (dirty-line tracking at buffer level)
- **Mouse mode detection** — `HasMouseMode()` checks all mouse tracking modes (X10, Normal, ButtonEvent, AnyEvent)
- **Mouse event encoding** — `EncodeMouseEvent()` produces X10 or SGR escape sequences
- **Kitty graphics + Sixel passthrough** support
- **Bracketed paste mode** detection and wrapping

### Splitty (`terminal/screen.go`)

Custom VT emulator — functional but simpler:
- No dirty-line tracking
- No mouse mode detection or encoding
- No callback system for mode changes
- No alt screen tracking at emulator level

---

## 6. Object Pooling & Memory (`internal/pool/pool.go`)

Tuios pools frequently allocated objects:

```go
var (
    stringBuilderPool = sync.Pool{New: func() any { return &strings.Builder{} }}
    layerPool         = sync.Pool{New: func() any { return &layers }}
    byteSlicePool     = sync.Pool{New: func() any { buf := make([]byte, 32*1024); return &buf }}
    stylePool         = sync.Pool{New: func() any { style := lipgloss.NewStyle(); return &style }}
    highlightGridPool = sync.Pool{New: func() any { return &HighlightGrid{} }}
)
```

`HighlightGrid` is a sparse boolean grid for marking search matches, visual selections, etc. Uses `[][]bool` with lazy row allocation — only allocates rows that actually have highlights.

---

## Recommendations for Splitty

### High Impact, Lower Effort

1. **Content caching per pane** — Cache the rendered string, only rebuild when content actually changes. Add `ContentDirty` vs `PositionDirty` dirty flags.
2. **Style batching** — Group consecutive cells with same attributes into one ANSI escape + text run. Dramatically cuts render output size.
3. **Frame skipping** — Don't re-render if nothing changed since last frame.
4. **Cache scrollback index math** — Pre-calculate `oldestIdx` once per render, not per cell.

### High Impact, Medium Effort

5. **Copy mode with vim motions** — Tuios's `CopyMode` state machine is well-structured and directly applicable to Splitty.
6. **Mouse event forwarding** — Check if running app has mouse mode, translate coords, encode and write to PTY.
7. **Visual-only resize during drag** with deferred PTY resize on release.
8. **Text selection** — Single/double/triple click modes with clipboard integration.

### Medium Impact, Higher Effort

9. **Object pooling** for string builders and cell arrays via `sync.Pool`.
10. **Style cache** with hash-based lookup to avoid rebuilding identical `tcell.Style` objects.
11. **Soft-wrap tracking** in scrollback for proper reflow on terminal resize.
12. **Auto-scroll improvements** — Re-enable auto-scroll when new output arrives while scrolled back.
