package splitty

// Direction specifies split orientation.
type Direction int

const (
	// Vertical splits panes side by side (│).
	Vertical Direction = iota
	// Horizontal stacks panes top and bottom (─).
	Horizontal
)

// String returns the direction name.
func (d Direction) String() string {
	if d == Vertical {
		return "vertical"
	}
	return "horizontal"
}

// PaneSplitMsg is sent when a pane is split.
type PaneSplitMsg struct {
	ParentID  string
	NewPaneID string
	Direction Direction
}

// PaneClosedMsg is sent when a pane is closed.
type PaneClosedMsg struct {
	PaneID string
	Err    error
}

// PaneFocusedMsg is sent when focus changes to a new pane.
type PaneFocusedMsg struct {
	PaneID string
}

// PaneOutputMsg is sent when a pane's PTY produces output.
type PaneOutputMsg struct {
	PaneID string
	Data   []byte
}

// PaneResizedMsg is sent when a pane is resized.
type PaneResizedMsg struct {
	PaneID string
	Width  int
	Height int
}

// LayoutLoadedMsg is sent when a layout is loaded from file.
type LayoutLoadedMsg struct {
	Err error
}
