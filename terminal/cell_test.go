package terminal

import (
	"strings"
	"testing"
)

func TestDefaultStyle(t *testing.T) {
	s := DefaultStyle()
	if s.Bold || s.Dim || s.Italic || s.Underline || s.Blink ||
		s.Reverse || s.Hidden || s.Strike {
		t.Error("default style should have no attributes set")
	}
	if s.FG.Type != ColorDefault {
		t.Error("default FG should be ColorDefault")
	}
	if s.BG.Type != ColorDefault {
		t.Error("default BG should be ColorDefault")
	}
}

func TestEmptyCell(t *testing.T) {
	c := EmptyCell()
	if c.Rune != ' ' {
		t.Errorf("expected space, got %c", c.Rune)
	}
	if c.Style != DefaultStyle() {
		t.Error("empty cell should have default style")
	}
}

func TestStyleEqual(t *testing.T) {
	s1 := DefaultStyle()
	s2 := DefaultStyle()
	if !s1.Equal(s2) {
		t.Error("two default styles should be equal")
	}

	s2.Bold = true
	if s1.Equal(s2) {
		t.Error("different styles should not be equal")
	}
}

func TestStyleToANSIDefault(t *testing.T) {
	s := DefaultStyle()
	ansi := s.ToANSI()
	if ansi != "\x1b[0m" {
		t.Errorf("default style ANSI: expected \\x1b[0m, got %q", ansi)
	}
}

func TestStyleToANSIAttributes(t *testing.T) {
	tests := []struct {
		name     string
		style    Style
		contains string
	}{
		{"bold", Style{Bold: true}, ";1"},
		{"dim", Style{Dim: true}, ";2"},
		{"italic", Style{Italic: true}, ";3"},
		{"underline", Style{Underline: true}, ";4"},
		{"blink", Style{Blink: true}, ";5"},
		{"reverse", Style{Reverse: true}, ";7"},
		{"hidden", Style{Hidden: true}, ";8"},
		{"strike", Style{Strike: true}, ";9"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ansi := tt.style.ToANSI()
			if !strings.Contains(ansi, tt.contains) {
				t.Errorf("expected ANSI to contain %q, got %q", tt.contains, ansi)
			}
		})
	}
}

func TestStyleToANSIMultipleAttributes(t *testing.T) {
	s := Style{Bold: true, Italic: true, Underline: true}
	ansi := s.ToANSI()
	if !strings.Contains(ansi, ";1") {
		t.Error("missing bold")
	}
	if !strings.Contains(ansi, ";3") {
		t.Error("missing italic")
	}
	if !strings.Contains(ansi, ";4") {
		t.Error("missing underline")
	}
}

func TestWriteColor16(t *testing.T) {
	t.Run("foreground standard", func(t *testing.T) {
		s := Style{FG: Color{Type: Color16, Value: 1}} // red
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";31") {
			t.Errorf("expected ;31, got %q", ansi)
		}
	})

	t.Run("foreground bright", func(t *testing.T) {
		s := Style{FG: Color{Type: Color16, Value: 9}} // bright red
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";91") {
			t.Errorf("expected ;91, got %q", ansi)
		}
	})

	t.Run("background standard", func(t *testing.T) {
		s := Style{BG: Color{Type: Color16, Value: 2}} // green
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";42") {
			t.Errorf("expected ;42, got %q", ansi)
		}
	})

	t.Run("background bright", func(t *testing.T) {
		s := Style{BG: Color{Type: Color16, Value: 10}} // bright green
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";102") {
			t.Errorf("expected ;102, got %q", ansi)
		}
	})
}

func TestWriteColor256(t *testing.T) {
	t.Run("foreground", func(t *testing.T) {
		s := Style{FG: Color{Type: Color256, Value: 196}}
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";38;5;196") {
			t.Errorf("expected ;38;5;196, got %q", ansi)
		}
	})

	t.Run("background", func(t *testing.T) {
		s := Style{BG: Color{Type: Color256, Value: 33}}
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";48;5;33") {
			t.Errorf("expected ;48;5;33, got %q", ansi)
		}
	})
}

func TestWriteColorRGB(t *testing.T) {
	t.Run("foreground", func(t *testing.T) {
		// 0xFF0000 = red
		s := Style{FG: Color{Type: ColorRGB, Value: 0xFF0000}}
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";38;2;255;0;0") {
			t.Errorf("expected ;38;2;255;0;0, got %q", ansi)
		}
	})

	t.Run("background", func(t *testing.T) {
		// 0x00FF00 = green
		s := Style{BG: Color{Type: ColorRGB, Value: 0x00FF00}}
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";48;2;0;255;0") {
			t.Errorf("expected ;48;2;0;255;0, got %q", ansi)
		}
	})

	t.Run("mixed RGB", func(t *testing.T) {
		// 0x1A2B3C
		s := Style{FG: Color{Type: ColorRGB, Value: 0x1A2B3C}}
		ansi := s.ToANSI()
		if !strings.Contains(ansi, ";38;2;26;43;60") {
			t.Errorf("expected ;38;2;26;43;60, got %q", ansi)
		}
	})
}

func TestColorConstants(t *testing.T) {
	if ColorDefault != 0 {
		t.Error("ColorDefault should be 0")
	}
	if Color16 == ColorDefault {
		t.Error("Color16 should differ from ColorDefault")
	}
	if Color256 == Color16 {
		t.Error("Color256 should differ from Color16")
	}
	if ColorRGB == Color256 {
		t.Error("ColorRGB should differ from Color256")
	}
}
