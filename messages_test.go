package splitty

import "testing"

func TestDirectionString(t *testing.T) {
	tests := []struct {
		dir  Direction
		want string
	}{
		{Vertical, "vertical"},
		{Horizontal, "horizontal"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.dir.String()
			if got != tt.want {
				t.Errorf("Direction(%d).String() = %q, want %q", tt.dir, got, tt.want)
			}
		})
	}
}

func TestDirectionConstants(t *testing.T) {
	// Ensure constants are distinct
	if Vertical == Horizontal {
		t.Error("Vertical and Horizontal should be different values")
	}
}
