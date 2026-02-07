package terminal

import (
	"testing"
)

func TestEncodeMouseEventSGR(t *testing.T) {
	tests := []struct {
		name     string
		button   int
		action   int
		x, y     int
		expected string
	}{
		{
			name:     "left press at 0,0",
			button:   MouseButtonLeft,
			action:   MouseActionPress,
			x:        0, y: 0,
			expected: "\x1b[<0;1;1M",
		},
		{
			name:     "right press at 10,20",
			button:   MouseButtonRight,
			action:   MouseActionPress,
			x:        10, y: 20,
			expected: "\x1b[<2;11;21M",
		},
		{
			name:     "left release at 5,5",
			button:   MouseButtonLeft,
			action:   MouseActionRelease,
			x:        5, y: 5,
			expected: "\x1b[<0;6;6m",
		},
		{
			name:     "motion with left button at 3,4",
			button:   MouseButtonLeft,
			action:   MouseActionMotion,
			x:        3, y: 4,
			expected: "\x1b[<32;4;5M",
		},
		{
			name:     "wheel up at 0,0",
			button:   MouseButtonWheelUp,
			action:   MouseActionPress,
			x:        0, y: 0,
			expected: "\x1b[<64;1;1M",
		},
		{
			name:     "wheel down at 0,0",
			button:   MouseButtonWheelDown,
			action:   MouseActionPress,
			x:        0, y: 0,
			expected: "\x1b[<65;1;1M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeMouseEvent(tt.button, tt.action, tt.x, tt.y, 1006)
			if string(got) != tt.expected {
				t.Errorf("got %q, want %q", string(got), tt.expected)
			}
		})
	}
}

func TestEncodeMouseEventX10(t *testing.T) {
	tests := []struct {
		name   string
		button int
		action int
		x, y   int
	}{
		{"left press", MouseButtonLeft, MouseActionPress, 0, 0},
		{"right press", MouseButtonRight, MouseActionPress, 5, 5},
		{"release", MouseButtonLeft, MouseActionRelease, 0, 0},
		{"motion", MouseButtonLeft, MouseActionMotion, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeMouseEvent(tt.button, tt.action, tt.x, tt.y, 0)
			if len(got) != 6 {
				t.Fatalf("expected 6 bytes, got %d", len(got))
			}
			if got[0] != 0x1b || got[1] != '[' || got[2] != 'M' {
				t.Error("expected ESC [ M prefix")
			}

			// Verify coordinate encoding
			expectedCx := byte(tt.x + 33)
			expectedCy := byte(tt.y + 33)
			if got[4] != expectedCx {
				t.Errorf("x: got %d, want %d", got[4], expectedCx)
			}
			if got[5] != expectedCy {
				t.Errorf("y: got %d, want %d", got[5], expectedCy)
			}

			// Verify button encoding
			if tt.action == MouseActionRelease {
				if got[3] != byte(3+32) {
					t.Errorf("release: got cb=%d, want %d", got[3], 3+32)
				}
			} else if tt.action == MouseActionMotion {
				if got[3] != byte(tt.button+32+32) {
					t.Errorf("motion: got cb=%d, want %d", got[3], tt.button+32+32)
				}
			} else {
				if got[3] != byte(tt.button+32) {
					t.Errorf("press: got cb=%d, want %d", got[3], tt.button+32)
				}
			}
		})
	}
}

func TestEncodeMouseEventX10CoordCap(t *testing.T) {
	// Coordinates beyond 222 should be capped at 255
	got := EncodeMouseEvent(MouseButtonLeft, MouseActionPress, 300, 300, 0)
	if got[4] != 255 || got[5] != 255 {
		t.Errorf("coords should be capped at 255, got cx=%d cy=%d", got[4], got[5])
	}
}
