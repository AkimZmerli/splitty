package terminal

import "fmt"

// Mouse button values for VT encoding.
const (
	MouseButtonLeft     = 0
	MouseButtonMiddle   = 1
	MouseButtonRight    = 2
	MouseButtonRelease  = 3
	MouseButtonWheelUp  = 64
	MouseButtonWheelDown = 65
)

// Mouse action types.
const (
	MouseActionPress   = 0
	MouseActionRelease = 1
	MouseActionMotion  = 2
)

// EncodeMouseEvent encodes a mouse event for transmission to a terminal application.
// button: VT button value (0=left, 1=middle, 2=right, 64=wheelUp, 65=wheelDown)
// action: 0=press, 1=release, 2=motion
// x, y: 0-based terminal coordinates
// encoding: 0=X10/Normal, 1006=SGR
func EncodeMouseEvent(button, action, x, y, encoding int) []byte {
	if encoding == 1006 {
		return encodeSGR(button, action, x, y)
	}
	return encodeX10(button, action, x, y)
}

// encodeSGR produces SGR-encoded mouse events: ESC [ < Cb ; Cx ; Cy M/m
// Coords are 1-based. M=press/motion, m=release.
func encodeSGR(button, action, x, y int) []byte {
	cb := button
	if action == MouseActionMotion {
		cb += 32
	}

	final := byte('M')
	if action == MouseActionRelease {
		final = 'm'
	}

	return []byte(fmt.Sprintf("\x1b[<%d;%d;%d%c", cb, x+1, y+1, final))
}

// encodeX10 produces X10/Normal-encoded mouse events: ESC [ M Cb Cx Cy
// Cb = button + 32, Cx/Cy = coord + 33. Coords capped at 222 (255-33).
func encodeX10(button, action, x, y int) []byte {
	cb := button + 32
	if action == MouseActionMotion {
		cb += 32
	}
	if action == MouseActionRelease {
		cb = 3 + 32 // release is encoded as button 3
	}

	// Cap coordinates for X10 encoding (max value 255, offset 33)
	cx := x + 33
	cy := y + 33
	if cx > 255 {
		cx = 255
	}
	if cy > 255 {
		cy = 255
	}

	return []byte{0x1b, '[', 'M', byte(cb), byte(cx), byte(cy)}
}
