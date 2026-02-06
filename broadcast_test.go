package splitty

import "testing"

func TestSetBroadcast(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	m.SetBroadcast(true)
	if !m.broadcasting {
		t.Error("broadcasting should be true")
	}

	m.SetBroadcast(false)
	if m.broadcasting {
		t.Error("broadcasting should be false")
	}
}

func TestIsBroadcasting(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	if m.IsBroadcasting() {
		t.Error("should not be broadcasting initially")
	}

	m.broadcasting = true
	if !m.IsBroadcasting() {
		t.Error("should be broadcasting")
	}
}

func TestSendInput(t *testing.T) {
	// SendInput with nil PTYs just shouldn't panic
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	t.Run("send to focused pane", func(t *testing.T) {
		m.SendInput([]byte("hello")) // no PTY, should not panic
	})

	t.Run("broadcast mode", func(t *testing.T) {
		m.broadcasting = true
		m.SendInput([]byte("hello")) // broadcasts to all, no PTYs
		m.broadcasting = false
	})

	t.Run("nil root", func(t *testing.T) {
		m2 := testManager(nil, "")
		m2.SendInput([]byte("hello")) // should not panic
	})
}
