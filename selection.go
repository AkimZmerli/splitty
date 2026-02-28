package splitty

import "unicode"

// SelectionMode determines the granularity of text selection.
type SelectionMode int

const (
	SelectNone SelectionMode = iota
	SelectChar
	SelectWord
	SelectLine
)

// Selection represents a text selection within a pane.
type Selection struct {
	Active   bool
	Mode     SelectionMode
	PaneID   string
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}

// Normalize returns the selection with start <= end.
func (s *Selection) Normalize() (startRow, startCol, endRow, endCol int) {
	if s.StartRow < s.EndRow || (s.StartRow == s.EndRow && s.StartCol <= s.EndCol) {
		return s.StartRow, s.StartCol, s.EndRow, s.EndCol
	}
	return s.EndRow, s.EndCol, s.StartRow, s.StartCol
}

// Contains returns true if the given row/col is within the selection.
func (s *Selection) Contains(row, col int) bool {
	if !s.Active {
		return false
	}
	sr, sc, er, ec := s.Normalize()
	if row < sr || row > er {
		return false
	}
	if s.Mode == SelectLine {
		return true
	}
	if sr == er {
		return col >= sc && col <= ec
	}
	if row == sr {
		return col >= sc
	}
	if row == er {
		return col <= ec
	}
	return true
}

// isWordChar returns true if the rune is part of a word (not whitespace/punctuation).
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
