package terminal

import (
	"strings"
	"testing"
)

func TestNewScreen(t *testing.T) {
	s := NewScreen(80, 24)
	if s.Width() != 80 {
		t.Errorf("expected width 80, got %d", s.Width())
	}
	if s.Height() != 24 {
		t.Errorf("expected height 24, got %d", s.Height())
	}
	if s.Title() != "" {
		t.Errorf("expected empty title, got %q", s.Title())
	}
}

func TestScreenWriteASCII(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("ABC"))

	// Characters should be at row 0, cols 0-2
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("cell(0,0): expected A, got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][1].Rune != 'B' {
		t.Errorf("cell(0,1): expected B, got %c", s.cells[0][1].Rune)
	}
	if s.cells[0][2].Rune != 'C' {
		t.Errorf("cell(0,2): expected C, got %c", s.cells[0][2].Rune)
	}
}

func TestScreenWriteNewline(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("A\nB"))

	if s.cells[0][0].Rune != 'A' {
		t.Errorf("expected A at row 0")
	}
	// After LF, cursor moves down but col stays (no CR)
	// Actually LF only moves down. B is at (1, 1) since col didn't reset
	// Wait, after writing A, col=1. LF moves down. B at (1,1).
	if s.cells[1][1].Rune != 'B' {
		t.Errorf("expected B at row 1, col 1, got %c", s.cells[1][1].Rune)
	}
}

func TestScreenWriteCRLF(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("AB\r\nCD"))

	if s.cells[0][0].Rune != 'A' {
		t.Error("expected A at row 0 col 0")
	}
	if s.cells[0][1].Rune != 'B' {
		t.Error("expected B at row 0 col 1")
	}
	if s.cells[1][0].Rune != 'C' {
		t.Error("expected C at row 1 col 0")
	}
	if s.cells[1][1].Rune != 'D' {
		t.Error("expected D at row 1 col 1")
	}
}

func TestScreenResize(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("Hello"))

	s.Resize(20, 10)
	if s.Width() != 20 {
		t.Errorf("expected width 20, got %d", s.Width())
	}
	if s.Height() != 10 {
		t.Errorf("expected height 10, got %d", s.Height())
	}

	// Original content should be preserved
	if s.cells[0][0].Rune != 'H' {
		t.Errorf("expected H preserved, got %c", s.cells[0][0].Rune)
	}
}

func TestScreenResizeShrink(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("Hello World"))

	s.Resize(5, 3)
	if s.Width() != 5 {
		t.Errorf("expected width 5, got %d", s.Width())
	}
	if s.Height() != 3 {
		t.Errorf("expected height 3, got %d", s.Height())
	}

	// First 5 chars should be preserved
	expected := "Hello"
	for i, ch := range expected {
		if s.cells[0][i].Rune != ch {
			t.Errorf("cell(0,%d): expected %c, got %c", i, ch, s.cells[0][i].Rune)
		}
	}
}

func TestCSICursorMovement(t *testing.T) {
	s := NewScreen(20, 10)

	t.Run("CUP - move to position", func(t *testing.T) {
		s.Write([]byte("\x1b[3;5H")) // row 3, col 5 (1-based)
		if s.cursor.Row != 2 || s.cursor.Col != 4 {
			t.Errorf("expected (2,4), got (%d,%d)", s.cursor.Row, s.cursor.Col)
		}
	})

	t.Run("CUU - cursor up", func(t *testing.T) {
		s.Write([]byte("\x1b[5;5H")) // start at (4,4)
		s.Write([]byte("\x1b[2A"))   // up 2
		if s.cursor.Row != 2 {
			t.Errorf("expected row 2, got %d", s.cursor.Row)
		}
	})

	t.Run("CUD - cursor down", func(t *testing.T) {
		s.Write([]byte("\x1b[1;1H")) // start at (0,0)
		s.Write([]byte("\x1b[3B"))   // down 3
		if s.cursor.Row != 3 {
			t.Errorf("expected row 3, got %d", s.cursor.Row)
		}
	})

	t.Run("CUF - cursor forward", func(t *testing.T) {
		s.Write([]byte("\x1b[1;1H")) // start at (0,0)
		s.Write([]byte("\x1b[5C"))   // right 5
		if s.cursor.Col != 5 {
			t.Errorf("expected col 5, got %d", s.cursor.Col)
		}
	})

	t.Run("CUB - cursor back", func(t *testing.T) {
		s.Write([]byte("\x1b[1;10H")) // start at (0,9)
		s.Write([]byte("\x1b[3D"))    // left 3
		if s.cursor.Col != 6 {
			t.Errorf("expected col 6, got %d", s.cursor.Col)
		}
	})

	t.Run("CHA - cursor horizontal absolute", func(t *testing.T) {
		s.Write([]byte("\x1b[1;1H")) // reset
		s.Write([]byte("\x1b[8G"))   // col 8 (1-based)
		if s.cursor.Col != 7 {
			t.Errorf("expected col 7, got %d", s.cursor.Col)
		}
	})

	t.Run("VPA - line position absolute", func(t *testing.T) {
		s.Write([]byte("\x1b[5d")) // row 5 (1-based)
		if s.cursor.Row != 4 {
			t.Errorf("expected row 4, got %d", s.cursor.Row)
		}
	})

	t.Run("CNL - cursor next line", func(t *testing.T) {
		s.Write([]byte("\x1b[3;5H")) // start at (2,4)
		s.Write([]byte("\x1b[2E"))   // next line x2
		if s.cursor.Row != 4 || s.cursor.Col != 0 {
			t.Errorf("expected (4,0), got (%d,%d)", s.cursor.Row, s.cursor.Col)
		}
	})

	t.Run("CPL - cursor previous line", func(t *testing.T) {
		s.Write([]byte("\x1b[5;5H")) // start at (4,4)
		s.Write([]byte("\x1b[2F"))   // prev line x2
		if s.cursor.Row != 2 || s.cursor.Col != 0 {
			t.Errorf("expected (2,0), got (%d,%d)", s.cursor.Row, s.cursor.Col)
		}
	})
}

func TestCSICursorClamp(t *testing.T) {
	s := NewScreen(10, 5)

	t.Run("cursor up clamped at 0", func(t *testing.T) {
		s.Write([]byte("\x1b[1;1H"))  // go to (0,0)
		s.Write([]byte("\x1b[100A")) // up 100
		if s.cursor.Row != 0 {
			t.Errorf("expected row 0, got %d", s.cursor.Row)
		}
	})

	t.Run("cursor down clamped at height-1", func(t *testing.T) {
		s.Write([]byte("\x1b[1;1H"))  // go to (0,0)
		s.Write([]byte("\x1b[100B")) // down 100
		if s.cursor.Row != 4 {
			t.Errorf("expected row 4, got %d", s.cursor.Row)
		}
	})

	t.Run("cursor right clamped at width-1", func(t *testing.T) {
		s.Write([]byte("\x1b[1;1H"))  // go to (0,0)
		s.Write([]byte("\x1b[100C")) // right 100
		if s.cursor.Col != 9 {
			t.Errorf("expected col 9, got %d", s.cursor.Col)
		}
	})
}

func TestCSIEraseDisplay(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("ABCDE"))
	s.Write([]byte("\r\nFGHIJ"))
	s.Write([]byte("\r\nKLMNO"))

	t.Run("erase below", func(t *testing.T) {
		s2 := NewScreen(5, 3)
		s2.Write([]byte("ABCDE\r\nFGHIJ\r\nKLMNO"))
		s2.Write([]byte("\x1b[2;3H")) // row 1, col 2
		s2.Write([]byte("\x1b[0J"))   // erase below
		// Row 1 from col 2 and all of row 2 should be erased
		if s2.cells[1][2].Rune != ' ' {
			t.Error("cell(1,2) should be erased")
		}
		if s2.cells[2][0].Rune != ' ' {
			t.Error("cell(2,0) should be erased")
		}
		// Row 1 before col 2 should be preserved
		if s2.cells[1][0].Rune != 'F' {
			t.Errorf("cell(1,0) should be F, got %c", s2.cells[1][0].Rune)
		}
	})

	t.Run("erase all", func(t *testing.T) {
		s2 := NewScreen(5, 3)
		s2.Write([]byte("ABCDE\r\nFGHIJ\r\nKLMNO"))
		s2.Write([]byte("\x1b[2J")) // erase entire display
		for r := 0; r < 3; r++ {
			for c := 0; c < 5; c++ {
				if s2.cells[r][c].Rune != ' ' {
					t.Errorf("cell(%d,%d) should be space, got %c", r, c, s2.cells[r][c].Rune)
				}
			}
		}
	})
}

func TestCSIEraseLine(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("ABCDE"))
	s.Write([]byte("\x1b[1;3H")) // row 0, col 2
	s.Write([]byte("\x1b[0K"))   // erase right

	// Cols 2-4 should be erased
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("cell(0,0) should be A, got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][1].Rune != 'B' {
		t.Errorf("cell(0,1) should be B, got %c", s.cells[0][1].Rune)
	}
	if s.cells[0][2].Rune != ' ' {
		t.Errorf("cell(0,2) should be erased, got %c", s.cells[0][2].Rune)
	}
	if s.cells[0][4].Rune != ' ' {
		t.Errorf("cell(0,4) should be erased, got %c", s.cells[0][4].Rune)
	}
}

func TestCSIEraseLineLeft(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("ABCDE"))
	s.Write([]byte("\x1b[1;3H")) // row 0, col 2
	s.Write([]byte("\x1b[1K"))   // erase left

	// Cols 0-2 should be erased
	if s.cells[0][0].Rune != ' ' {
		t.Errorf("cell(0,0) should be erased, got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][2].Rune != ' ' {
		t.Errorf("cell(0,2) should be erased, got %c", s.cells[0][2].Rune)
	}
	// Cols 3-4 preserved
	if s.cells[0][3].Rune != 'D' {
		t.Errorf("cell(0,3) should be D, got %c", s.cells[0][3].Rune)
	}
}

func TestCSIEraseEntireLine(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("ABCDE"))
	s.Write([]byte("\x1b[1;3H")) // row 0, col 2
	s.Write([]byte("\x1b[2K"))   // erase entire line

	for c := 0; c < 5; c++ {
		if s.cells[0][c].Rune != ' ' {
			t.Errorf("cell(0,%d) should be erased, got %c", c, s.cells[0][c].Rune)
		}
	}
}

func TestSGRAttributes(t *testing.T) {
	s := NewScreen(10, 5)

	t.Run("bold", func(t *testing.T) {
		s.Write([]byte("\x1b[1mA"))
		if !s.cells[0][0].Style.Bold {
			t.Error("expected bold")
		}
	})

	t.Run("dim", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[2mB"))
		if !s.cells[0][0].Style.Dim {
			t.Error("expected dim")
		}
	})

	t.Run("italic", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[3mC"))
		if !s.cells[0][0].Style.Italic {
			t.Error("expected italic")
		}
	})

	t.Run("underline", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[4mD"))
		if !s.cells[0][0].Style.Underline {
			t.Error("expected underline")
		}
	})

	t.Run("reset clears all", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[1;3;4m\x1b[0mE"))
		cell := s.cells[0][0]
		if cell.Style.Bold || cell.Style.Italic || cell.Style.Underline {
			t.Error("reset should clear all attributes")
		}
	})
}

func TestSGRColors(t *testing.T) {
	s := NewScreen(10, 5)

	t.Run("16-color foreground", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[31mR")) // red
		if s.cells[0][0].Style.FG.Type != Color16 {
			t.Error("expected Color16 FG")
		}
		if s.cells[0][0].Style.FG.Value != 1 {
			t.Errorf("expected value 1 (red), got %d", s.cells[0][0].Style.FG.Value)
		}
	})

	t.Run("16-color background", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[44mB")) // blue bg
		if s.cells[0][0].Style.BG.Type != Color16 {
			t.Error("expected Color16 BG")
		}
		if s.cells[0][0].Style.BG.Value != 4 {
			t.Errorf("expected value 4 (blue), got %d", s.cells[0][0].Style.BG.Value)
		}
	})

	t.Run("bright foreground", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[91mR")) // bright red
		if s.cells[0][0].Style.FG.Value != 9 {
			t.Errorf("expected bright red (9), got %d", s.cells[0][0].Style.FG.Value)
		}
	})

	t.Run("256-color foreground", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[38;5;196mX"))
		if s.cells[0][0].Style.FG.Type != Color256 {
			t.Error("expected Color256")
		}
		if s.cells[0][0].Style.FG.Value != 196 {
			t.Errorf("expected 196, got %d", s.cells[0][0].Style.FG.Value)
		}
	})

	t.Run("RGB foreground", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[38;2;255;128;0mX"))
		if s.cells[0][0].Style.FG.Type != ColorRGB {
			t.Error("expected ColorRGB")
		}
		expected := uint32(255)<<16 | uint32(128)<<8 | uint32(0)
		if s.cells[0][0].Style.FG.Value != expected {
			t.Errorf("expected 0x%06X, got 0x%06X", expected, s.cells[0][0].Style.FG.Value)
		}
	})

	t.Run("default foreground", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[31m\x1b[39mX"))
		if s.cells[0][0].Style.FG.Type != ColorDefault {
			t.Error("expected default FG after SGR 39")
		}
	})

	t.Run("default background", func(t *testing.T) {
		s.Write([]byte("\x1b[0m\x1b[1;1H\x1b[41m\x1b[49mX"))
		if s.cells[0][0].Style.BG.Type != ColorDefault {
			t.Error("expected default BG after SGR 49")
		}
	})
}

func TestOSCTitle(t *testing.T) {
	s := NewScreen(10, 5)

	t.Run("set title with OSC 0", func(t *testing.T) {
		s.Write([]byte("\x1b]0;My Title\x07"))
		if s.Title() != "My Title" {
			t.Errorf("expected 'My Title', got %q", s.Title())
		}
	})

	t.Run("set title with OSC 2", func(t *testing.T) {
		s.Write([]byte("\x1b]2;Window Title\x07"))
		if s.Title() != "Window Title" {
			t.Errorf("expected 'Window Title', got %q", s.Title())
		}
	})
}

func TestScrollUp(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("AAAAA\r\nBBBBB\r\nCCCCC"))

	// Fill screen, cursor at (2,4) after wrapPending
	// Scroll up 1 line
	s.Write([]byte("\x1b[1S"))

	// Row 0 should now be what was row 1
	if s.cells[0][0].Rune != 'B' {
		t.Errorf("row 0 should be B after scroll, got %c", s.cells[0][0].Rune)
	}
	// Last row should be empty
	if s.cells[2][0].Rune != ' ' {
		t.Errorf("last row should be empty, got %c", s.cells[2][0].Rune)
	}
}

func TestScrollDown(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("AAAAA\r\nBBBBB\r\nCCCCC"))

	s.Write([]byte("\x1b[1T"))

	// First row should be empty
	if s.cells[0][0].Rune != ' ' {
		t.Errorf("first row should be empty, got %c", s.cells[0][0].Rune)
	}
	// Row 1 should be what was row 0
	if s.cells[1][0].Rune != 'A' {
		t.Errorf("row 1 should be A after scroll, got %c", s.cells[1][0].Rune)
	}
}

func TestTab(t *testing.T) {
	s := NewScreen(20, 5)
	s.Write([]byte("AB\tC"))

	// Tab stops at multiples of 8. "AB" puts cursor at col 2.
	// Tab goes to col 8. Then C at col 8.
	if s.cells[0][8].Rune != 'C' {
		t.Errorf("expected C at col 8, got %c", s.cells[0][8].Rune)
	}
}

func TestBackspace(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("AB\bC"))

	// Write A at 0, B at 1, backspace to 0, C overwrites at 1
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("expected A at col 0, got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][1].Rune != 'C' {
		t.Errorf("expected C at col 1 (backspace + overwrite), got %c", s.cells[0][1].Rune)
	}
}

func TestBackspaceAtStart(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("\bA"))
	// Backspace at col 0 should stay at 0
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("expected A at col 0, got %c", s.cells[0][0].Rune)
	}
}

func TestUTF8(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("é")) // 2-byte UTF-8: 0xC3 0xA9

	if s.cells[0][0].Rune != 'é' {
		t.Errorf("expected é, got %c", s.cells[0][0].Rune)
	}
}

func TestUTF8MultibyteJapanese(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("日")) // 3-byte UTF-8

	if s.cells[0][0].Rune != '日' {
		t.Errorf("expected 日, got %c (0x%X)", s.cells[0][0].Rune, s.cells[0][0].Rune)
	}
}

func TestUTF8Emoji(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("😀")) // 4-byte UTF-8

	if s.cells[0][0].Rune != '😀' {
		t.Errorf("expected 😀, got %c", s.cells[0][0].Rune)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		v, lo, hi, want int
	}{
		{5, 0, 10, 5},   // within range
		{-1, 0, 10, 0},  // below
		{15, 0, 10, 10}, // above
		{0, 0, 10, 0},   // at lower bound
		{10, 0, 10, 10}, // at upper bound
	}

	for _, tt := range tests {
		got := clamp(tt.v, tt.lo, tt.hi)
		if got != tt.want {
			t.Errorf("clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
}

func TestScreenRender(t *testing.T) {
	s := NewScreen(3, 2)
	s.Write([]byte("ABC\r\nDEF"))

	rendered := s.Render()

	// Should contain the characters
	if !strings.Contains(rendered, "A") || !strings.Contains(rendered, "F") {
		t.Errorf("rendered output should contain characters, got %q", rendered)
	}

	// Should end with reset
	if !strings.HasSuffix(rendered, "\x1b[0m") {
		t.Error("rendered output should end with ANSI reset")
	}
}

func TestDECSaveCursor(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("\x1b[3;5H")) // move to (2,4)
	s.Write([]byte("\x1b7"))     // save cursor (DECSC)
	s.Write([]byte("\x1b[1;1H")) // move to (0,0)
	s.Write([]byte("\x1b8"))     // restore cursor (DECRC)

	if s.cursor.Row != 2 || s.cursor.Col != 4 {
		t.Errorf("expected restored cursor at (2,4), got (%d,%d)", s.cursor.Row, s.cursor.Col)
	}
}

func TestCSISaveCursor(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("\x1b[3;5H")) // move to (2,4)
	s.Write([]byte("\x1b[s"))    // SCP - save cursor
	s.Write([]byte("\x1b[1;1H")) // move to (0,0)
	s.Write([]byte("\x1b[u"))    // RCP - restore cursor

	if s.cursor.Row != 2 || s.cursor.Col != 4 {
		t.Errorf("expected restored cursor at (2,4), got (%d,%d)", s.cursor.Row, s.cursor.Col)
	}
}

func TestDECSTBM(t *testing.T) {
	s := NewScreen(10, 10)
	s.Write([]byte("\x1b[3;8r")) // set scroll region rows 3-8

	if s.scrollTop != 2 {
		t.Errorf("expected scrollTop=2, got %d", s.scrollTop)
	}
	if s.scrollBottom != 7 {
		t.Errorf("expected scrollBottom=7, got %d", s.scrollBottom)
	}
	// CUP should reset to origin
	if s.cursor.Row != 0 || s.cursor.Col != 0 {
		t.Error("cursor should reset to (0,0) after DECSTBM")
	}
}

func TestDECSetShowCursor(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("\x1b[?25l")) // hide cursor
	if s.cursor.Visible {
		t.Error("cursor should be hidden")
	}
	s.Write([]byte("\x1b[?25h")) // show cursor
	if !s.cursor.Visible {
		t.Error("cursor should be visible")
	}
}

func TestFullReset(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("Hello"))
	s.Write([]byte("\x1b]0;Title\x07"))
	s.Write([]byte("\x1bc")) // RIS - full reset

	if s.cursor.Row != 0 || s.cursor.Col != 0 {
		t.Error("cursor should be at (0,0) after reset")
	}
	if s.Title() != "" {
		t.Error("title should be empty after reset")
	}
	// Screen should be cleared
	if s.cells[0][0].Rune != ' ' {
		t.Errorf("cells should be empty after reset, got %c", s.cells[0][0].Rune)
	}
}

func TestInsertLines(t *testing.T) {
	s := NewScreen(5, 4)
	s.Write([]byte("AAAAA\r\nBBBBB\r\nCCCCC\r\nDDDDD"))
	s.Write([]byte("\x1b[2;1H")) // row 1, col 0
	s.Write([]byte("\x1b[1L"))   // insert 1 line

	// Row 1 should be blank, B moves to row 2
	if s.cells[1][0].Rune != ' ' {
		t.Errorf("inserted row should be blank, got %c", s.cells[1][0].Rune)
	}
	if s.cells[2][0].Rune != 'B' {
		t.Errorf("row 2 should have B, got %c", s.cells[2][0].Rune)
	}
}

func TestDeleteLines(t *testing.T) {
	s := NewScreen(5, 4)
	s.Write([]byte("AAAAA\r\nBBBBB\r\nCCCCC\r\nDDDDD"))
	s.Write([]byte("\x1b[2;1H")) // row 1, col 0
	s.Write([]byte("\x1b[1M"))   // delete 1 line

	// Row 1 should now be C (was row 2)
	if s.cells[1][0].Rune != 'C' {
		t.Errorf("row 1 should be C, got %c", s.cells[1][0].Rune)
	}
	// Last row should be blank
	if s.cells[3][0].Rune != ' ' {
		t.Errorf("last row should be blank, got %c", s.cells[3][0].Rune)
	}
}

func TestWrapPending(t *testing.T) {
	s := NewScreen(3, 3)
	s.Write([]byte("ABCD"))

	// After writing 3 chars (A, B, C), cursor is at col 2 with wrapPending.
	// Writing D triggers wrap: col=0, lineFeed, then D at (1,0).
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("expected A at (0,0), got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][2].Rune != 'C' {
		t.Errorf("expected C at (0,2), got %c", s.cells[0][2].Rune)
	}
	if s.cells[1][0].Rune != 'D' {
		t.Errorf("expected D at (1,0) after wrap, got %c", s.cells[1][0].Rune)
	}
}

func TestDeleteChars(t *testing.T) {
	s := NewScreen(5, 2)
	s.Write([]byte("ABCDE"))
	s.Write([]byte("\x1b[1;2H")) // col 1
	s.Write([]byte("\x1b[2P"))   // delete 2 chars

	// "ABCDE" -> delete 2 at col 1 -> "ADEE " -> actually "ADE  "
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("expected A at col 0, got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][1].Rune != 'D' {
		t.Errorf("expected D at col 1, got %c", s.cells[0][1].Rune)
	}
	if s.cells[0][2].Rune != 'E' {
		t.Errorf("expected E at col 2, got %c", s.cells[0][2].Rune)
	}
}

func TestInsertChars(t *testing.T) {
	s := NewScreen(5, 2)
	s.Write([]byte("ABCDE"))
	s.Write([]byte("\x1b[1;2H")) // col 1
	s.Write([]byte("\x1b[2@"))   // insert 2 chars

	// "ABCDE" -> insert 2 blanks at col 1 -> "A  BC"
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("expected A at col 0, got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][1].Rune != ' ' {
		t.Errorf("expected space at col 1, got %c", s.cells[0][1].Rune)
	}
	if s.cells[0][3].Rune != 'B' {
		t.Errorf("expected B at col 3, got %c", s.cells[0][3].Rune)
	}
}

func TestEraseChars(t *testing.T) {
	s := NewScreen(5, 2)
	s.Write([]byte("ABCDE"))
	s.Write([]byte("\x1b[1;2H")) // col 1
	s.Write([]byte("\x1b[2X"))   // erase 2 chars

	// "ABCDE" -> erase 2 at col 1 -> "A  DE"
	if s.cells[0][0].Rune != 'A' {
		t.Errorf("expected A at col 0, got %c", s.cells[0][0].Rune)
	}
	if s.cells[0][1].Rune != ' ' {
		t.Errorf("expected space at col 1, got %c", s.cells[0][1].Rune)
	}
	if s.cells[0][3].Rune != 'D' {
		t.Errorf("expected D at col 3, got %c", s.cells[0][3].Rune)
	}
}

func TestReverseIndex(t *testing.T) {
	s := NewScreen(10, 5)
	s.Write([]byte("\x1b[3;1H")) // row 2
	s.Write([]byte("\x1bM"))     // RI - reverse index

	if s.cursor.Row != 1 {
		t.Errorf("expected row 1, got %d", s.cursor.Row)
	}
}

func TestReverseIndexAtTop(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("AAAAA\r\nBBBBB\r\nCCCCC"))
	s.Write([]byte("\x1b[1;1H")) // row 0
	s.Write([]byte("\x1bM"))     // RI at scroll top -> scroll down

	// Row 0 should be blank (scrolled down)
	if s.cells[0][0].Rune != ' ' {
		t.Errorf("row 0 should be blank, got %c", s.cells[0][0].Rune)
	}
	// Row 1 should be A (was row 0)
	if s.cells[1][0].Rune != 'A' {
		t.Errorf("row 1 should be A, got %c", s.cells[1][0].Rune)
	}
}

func TestLineFeedAtBottom(t *testing.T) {
	s := NewScreen(5, 3)
	s.Write([]byte("AAAAA\r\nBBBBB\r\nCCCCC"))
	// Cursor is at row 2 (last row), col 4 with wrapPending
	s.Write([]byte("\n")) // LF at bottom -> scroll up

	// Row 0 should now be B, row 1 C, row 2 blank
	if s.cells[0][0].Rune != 'B' {
		t.Errorf("row 0 should be B, got %c", s.cells[0][0].Rune)
	}
	if s.cells[1][0].Rune != 'C' {
		t.Errorf("row 1 should be C, got %c", s.cells[1][0].Rune)
	}
	if s.cells[2][0].Rune != ' ' {
		t.Errorf("row 2 should be blank, got %c", s.cells[2][0].Rune)
	}
}

func TestResizeCursorClamp(t *testing.T) {
	s := NewScreen(10, 10)
	s.Write([]byte("\x1b[10;10H")) // cursor at (9,9)

	s.Resize(5, 5)
	if s.cursor.Row != 4 {
		t.Errorf("cursor row should be clamped to 4, got %d", s.cursor.Row)
	}
	if s.cursor.Col != 4 {
		t.Errorf("cursor col should be clamped to 4, got %d", s.cursor.Col)
	}
}

func TestScreenWriteImplementsIOWriter(t *testing.T) {
	s := NewScreen(10, 5)
	n, err := s.Write([]byte("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 bytes written, got %d", n)
	}
}
