# Theme Enhancements

## 1. Save Preferred Theme as Primary

### Goal
Allow users to persist their preferred theme choice across sessions.

### Current Behavior
- Themes cycle with `ctrl+t`
- Always starts with Tokyo Night on launch
- No persistence of user preference

### Proposed Solution
Add config file persistence to save theme preference.

**Config Location**: `~/.config/splitty/config.json`

**Structure**:
```json
{
  "theme": "Nightshade",
  "scrollSpeed": 3,
  "scrollbackLines": 1000
}
```

### Implementation Options

#### Option A: Auto-save (Recommended)
- Every theme cycle auto-saves to config
- Simple, immediate feedback
- User's last choice is always remembered

#### Option B: Explicit Save
- Add `ctrl+shift+t` to "lock in" current theme
- More control, but requires extra step
- Good for "trying out" themes without saving

#### Option C: Menu Option
- Add to help menu or startup menu
- Most discoverable but slowest workflow

### Implementation Tasks
- [ ] Create config module (`config.go`)
- [ ] Add config file read/write functions
- [ ] Load saved theme in `New()` constructor
- [ ] Save theme on cycle (or explicit command)
- [ ] Add `WithConfigPath()` option for custom config location
- [ ] Update documentation

---

## 2. Animated Theme Transitions

### Goal
Smooth, gradual color transitions when switching themes instead of instant changes.

### Visual Effect
```
Tokyo Night (blue #7AA2F7)
    ↓ (smooth fade over 1 second)
Nightshade (purple #BD93F9)
    ↓
Glacier (arctic blue #88C0D0)
```

### Technical Approach

**Color Interpolation**:
- Parse hex colors → RGB tuples
- Linear interpolation (lerp) between RGB values
- Generate intermediate colors at each frame

**Animation System**:
- Use `tea.Tick` for ~60fps updates
- Transition duration: 500ms-1000ms (configurable)
- Animate all theme properties simultaneously

**Code Structure**:
```go
type ThemeTransition struct {
    from       Theme
    to         Theme
    progress   float64  // 0.0 to 1.0
    duration   time.Duration
    startTime  time.Time
}

func (t *ThemeTransition) Interpolate() Theme {
    // Lerp all colors based on progress
}
```

### Design Questions

**Q1: When to animate?**
- Option A: Only when user presses `ctrl+t` (subtle, on-demand)
- Option B: Auto-loop demo mode with separate keybinding `ctrl+shift+a`
- Option C: Both - `ctrl+t` animates, new key for auto-loop demo

**Q2: Transition Duration**
- 500ms: Snappy, noticeable
- 1000ms: Smooth, cinematic
- 2000ms: Very slow, demo mode
- Configurable via `WithThemeTransitionDuration()`?

**Q3: Auto-loop Demo Mode?**
- Continuously cycle through all themes with transitions
- Good for showcasing splitty
- Separate keybinding: `ctrl+shift+a` (auto-cycle)
- Disable with any key press or same keybinding

### Implementation Tasks
- [ ] Create `transition.go` module
- [ ] Implement RGB color parsing and lerping
- [ ] Add `ThemeTransition` state to Manager
- [ ] Implement `tea.Tick` animation loop
- [ ] Add `WithAnimateThemes()` option (enable/disable)
- [ ] Add `WithThemeTransitionDuration()` option
- [ ] Add auto-loop demo mode (optional)
- [ ] Update all theme color applications to use interpolated theme
- [ ] Add tests for color interpolation
- [ ] Update documentation and README

### Color Interpolation Math
```go
func lerpColor(from, to lipgloss.Color, t float64) lipgloss.Color {
    r1, g1, b1 := parseHex(from)
    r2, g2, b2 := parseHex(to)

    r := lerp(r1, r2, t)
    g := lerp(g1, g2, t)
    b := lerp(b1, b2, t)

    return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

func lerp(a, b float64, t float64) uint8 {
    return uint8(a + (b-a)*t)
}
```

---

## Priority & Sequencing

### Recommended Order:
1. **Phase 1**: Save Preferred Theme (simpler, immediate value)
2. **Phase 2**: Animated Transitions (polish, wow factor)

### Decisions Needed
- [ ] Config save strategy (auto-save vs explicit)
- [ ] Animation trigger (ctrl+t only vs auto-loop mode)
- [ ] Transition duration default
- [ ] Whether to include auto-loop demo mode

---

## Theme System Architecture

### Current Themes
1. Tokyo Night (default on load)
2. Nightshade
3. Glacier
4. Sorbet

### Theme Properties (all need interpolation)
- `BorderActive`
- `BorderInactive`
- `BorderScrollback`
- `BorderScrollbackFocused`
- `BorderResize`
- `BorderCopyMode`
- `BorderCopyModeFocused`
- `Divider`
- `StatusBar` (background + foreground)
- `StatusText`

### Current Keybindings
- `ctrl+t`: Cycle theme (increment themeIndex)
- Proposed: `ctrl+shift+t`: Save current theme
- Proposed: `ctrl+shift+a`: Auto-loop demo mode
