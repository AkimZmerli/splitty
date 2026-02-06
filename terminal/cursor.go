package terminal

// Cursor represents the terminal cursor state.
type Cursor struct {
	Row     int
	Col     int
	Visible bool
	Style   Style
}

// DefaultCursor returns a visible cursor at position (0,0) with default style.
func DefaultCursor() Cursor {
	return Cursor{
		Row:     0,
		Col:     0,
		Visible: true,
		Style:   DefaultStyle(),
	}
}
