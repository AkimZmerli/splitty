package splitty

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// menuAction identifies a context menu action.
type menuAction int

const (
	actionSplitVertical menuAction = iota
	actionSplitHorizontal
	actionClosePane
)

// menuItem represents a single entry in the context menu.
type menuItem struct {
	label  string
	action menuAction
}

// contextMenu manages the right-click popup menu state.
type contextMenu struct {
	visible  bool
	x        int
	y        int
	items    []menuItem
	selected int
	paneID   string
}

func newContextMenu() contextMenu {
	return contextMenu{
		items: []menuItem{
			{label: "Split Vertical", action: actionSplitVertical},
			{label: "Split Horizontal", action: actionSplitHorizontal},
			{label: "Close Pane", action: actionClosePane},
		},
	}
}

func (c *contextMenu) show(x, y int, paneID string) {
	c.visible = true
	c.x = x
	c.y = y
	c.selected = 0
	c.paneID = paneID
}

func (c *contextMenu) hide() {
	c.visible = false
}

func (c *contextMenu) moveUp() {
	if c.selected > 0 {
		c.selected--
	}
}

func (c *contextMenu) moveDown() {
	if c.selected < len(c.items)-1 {
		c.selected++
	}
}

func (c *contextMenu) selectedAction() menuAction {
	return c.items[c.selected].action
}

// menuWidth returns the rendered width of the menu (label + padding + border).
func (c *contextMenu) menuWidth() int {
	maxLen := 0
	for _, item := range c.items {
		if len(item.label) > maxLen {
			maxLen = len(item.label)
		}
	}
	return maxLen + 4 // 2 chars padding on each side
}

// menuHeight returns the rendered height of the menu (items + border).
func (c *contextMenu) menuHeight() int {
	return len(c.items) + 2 // top and bottom border
}

// clampPosition adjusts the menu position to stay within screen bounds.
func (c *contextMenu) clampPosition(screenW, screenH int) (int, int) {
	x := c.x
	y := c.y
	mw := c.menuWidth()
	mh := c.menuHeight()

	if x+mw > screenW {
		x = screenW - mw
	}
	if x < 0 {
		x = 0
	}
	if y+mh > screenH {
		y = screenH - mh
	}
	if y < 0 {
		y = 0
	}
	return x, y
}

// hitTest returns true if the given coordinates are inside the menu.
func (c *contextMenu) hitTest(x, y, screenW, screenH int) bool {
	mx, my := c.clampPosition(screenW, screenH)
	mw := c.menuWidth()
	mh := c.menuHeight()
	return x >= mx && x < mx+mw && y >= my && y < my+mh
}

// hitTestItem returns the item index at the given y coordinate, or -1.
func (c *contextMenu) hitTestItem(x, y, screenW, screenH int) int {
	mx, my := c.clampPosition(screenW, screenH)
	mw := c.menuWidth()
	// Items start 1 row below menu top (border row)
	itemY := y - my - 1
	if x >= mx && x < mx+mw && itemY >= 0 && itemY < len(c.items) {
		return itemY
	}
	return -1
}

// render produces the menu box as a styled string with line breaks.
func (c *contextMenu) render(theme Theme) string {
	maxLen := 0
	for _, item := range c.items {
		if len(item.label) > maxLen {
			maxLen = len(item.label)
		}
	}

	normal := lipgloss.NewStyle().
		Padding(0, 1)
	highlight := lipgloss.NewStyle().
		Padding(0, 1).
		Background(theme.StatusBar.GetBackground()).
		Foreground(theme.StatusBar.GetForeground())

	var rows []string
	for i, item := range c.items {
		padded := item.label + strings.Repeat(" ", maxLen-len(item.label))
		if i == c.selected {
			rows = append(rows, highlight.Render(padded))
		} else {
			rows = append(rows, normal.Render(padded))
		}
	}

	content := strings.Join(rows, "\n")
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderActive.GetBorderTopForeground()).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252"))
	return box.Render(content)
}

// overlayOnView composites the menu box onto the background view at the
// clamped (x, y) position. Uses ANSI-aware string operations so that
// escape sequences in the background are not corrupted.
func (c *contextMenu) overlayOnView(bg string, screenW, screenH int, theme Theme) string {
	if !c.visible {
		return bg
	}

	menuStr := c.render(theme)
	mx, my := c.clampPosition(screenW, screenH)
	menuLines := strings.Split(menuStr, "\n")

	bgLines := strings.Split(bg, "\n")

	// Ensure bgLines has enough rows
	for len(bgLines) < screenH {
		bgLines = append(bgLines, strings.Repeat(" ", screenW))
	}

	for i, mLine := range menuLines {
		row := my + i
		if row < 0 || row >= len(bgLines) {
			continue
		}

		bgLine := bgLines[row]
		mLineW := ansi.StringWidth(mLine)
		bgW := ansi.StringWidth(bgLine)

		// Build: leftBg + menuLine + rightBg
		var result strings.Builder

		if mx > 0 {
			if bgW >= mx {
				result.WriteString(ansi.Truncate(bgLine, mx, ""))
			} else {
				result.WriteString(bgLine)
				result.WriteString(strings.Repeat(" ", mx-bgW))
			}
		}

		result.WriteString(mLine)

		rightStart := mx + mLineW
		if rightStart < bgW {
			result.WriteString(ansi.TruncateLeft(bgLine, rightStart, ""))
		}

		bgLines[row] = result.String()
	}

	return strings.Join(bgLines, "\n")
}
