package terminal

import (
	"fmt"
	"strings"
)

// Cell represents a single character cell in the terminal grid.
type Cell struct {
	Rune  rune
	Style Style
}

// Style represents the visual attributes of a terminal cell.
type Style struct {
	FG        Color
	BG        Color
	Bold      bool
	Dim       bool
	Italic    bool
	Underline bool
	Blink     bool
	Reverse   bool
	Hidden    bool
	Strike    bool
}

// Color represents a terminal color.
type Color struct {
	Type  ColorType
	Value uint32 // Color256: 0-255, ColorRGB: 0xRRGGBB
}

// ColorType identifies the color encoding.
type ColorType int

const (
	ColorDefault ColorType = iota
	Color16                // Standard 16 colors (0-15)
	Color256               // 256 palette colors (0-255)
	ColorRGB               // True color (24-bit RGB)
)

// DefaultStyle returns a style with all default values.
func DefaultStyle() Style {
	return Style{}
}

// EmptyCell returns a blank cell with default styling.
func EmptyCell() Cell {
	return Cell{Rune: ' ', Style: DefaultStyle()}
}

// ToANSI converts a Style to an ANSI SGR escape sequence string.
func (s Style) ToANSI() string {
	var b strings.Builder
	b.WriteString("\x1b[0")

	if s.Bold {
		b.WriteString(";1")
	}
	if s.Dim {
		b.WriteString(";2")
	}
	if s.Italic {
		b.WriteString(";3")
	}
	if s.Underline {
		b.WriteString(";4")
	}
	if s.Blink {
		b.WriteString(";5")
	}
	if s.Reverse {
		b.WriteString(";7")
	}
	if s.Hidden {
		b.WriteString(";8")
	}
	if s.Strike {
		b.WriteString(";9")
	}

	writeColor(&b, s.FG, true)
	writeColor(&b, s.BG, false)

	b.WriteByte('m')
	return b.String()
}

func writeColor(b *strings.Builder, c Color, fg bool) {
	switch c.Type {
	case Color16:
		base := 30
		if !fg {
			base = 40
		}
		if c.Value < 8 {
			fmt.Fprintf(b, ";%d", base+int(c.Value))
		} else {
			fmt.Fprintf(b, ";%d", base+60+int(c.Value)-8)
		}
	case Color256:
		if fg {
			fmt.Fprintf(b, ";38;5;%d", c.Value)
		} else {
			fmt.Fprintf(b, ";48;5;%d", c.Value)
		}
	case ColorRGB:
		r := (c.Value >> 16) & 0xFF
		g := (c.Value >> 8) & 0xFF
		bl := c.Value & 0xFF
		if fg {
			fmt.Fprintf(b, ";38;2;%d;%d;%d", r, g, bl)
		} else {
			fmt.Fprintf(b, ";48;2;%d;%d;%d", r, g, bl)
		}
	}
}

// Equal returns true if two styles are identical.
func (s Style) Equal(other Style) bool {
	return s == other
}
