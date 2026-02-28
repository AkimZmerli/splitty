package splitty

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestKeyToBytes(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.KeyPressMsg
		expected []byte
	}{
		{"enter", tea.KeyPressMsg{Code: tea.KeyEnter}, []byte{'\r'}},
		{"tab", tea.KeyPressMsg{Code: tea.KeyTab}, []byte{'\t'}},
		{"backspace", tea.KeyPressMsg{Code: tea.KeyBackspace}, []byte{0x7f}},
		{"escape", tea.KeyPressMsg{Code: tea.KeyEscape}, []byte{0x1b}},
		{"up", tea.KeyPressMsg{Code: tea.KeyUp}, []byte("\x1b[A")},
		{"down", tea.KeyPressMsg{Code: tea.KeyDown}, []byte("\x1b[B")},
		{"right", tea.KeyPressMsg{Code: tea.KeyRight}, []byte("\x1b[C")},
		{"left", tea.KeyPressMsg{Code: tea.KeyLeft}, []byte("\x1b[D")},
		{"home", tea.KeyPressMsg{Code: tea.KeyHome}, []byte("\x1b[H")},
		{"end", tea.KeyPressMsg{Code: tea.KeyEnd}, []byte("\x1b[F")},
		{"pgup", tea.KeyPressMsg{Code: tea.KeyPgUp}, []byte("\x1b[5~")},
		{"pgdown", tea.KeyPressMsg{Code: tea.KeyPgDown}, []byte("\x1b[6~")},
		{"delete", tea.KeyPressMsg{Code: tea.KeyDelete}, []byte("\x1b[3~")},
		{"space", tea.KeyPressMsg{Code: tea.KeySpace}, []byte{' '}},
		{"runes", tea.KeyPressMsg{Text: "hi"}, []byte("hi")},
		{"ctrl+a", tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl}, []byte{0x01}},
		{"ctrl+c", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}, []byte{0x03}},
		{"ctrl+z", tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl}, []byte{0x1a}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keyToBytes(tt.msg)
			if len(got) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, expected %d", len(got), len(tt.expected))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("byte %d: got 0x%02x, expected 0x%02x", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestKeyToBytesUnknown(t *testing.T) {
	// F-keys or other unmapped types should return nil
	msg := tea.KeyPressMsg{Code: tea.KeyF1}
	got := keyToBytes(msg)
	if got != nil {
		t.Errorf("expected nil for unmapped key, got %v", got)
	}
}

func TestFindPane(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	t.Run("find existing pane", func(t *testing.T) {
		p := m.findPane("a")
		if p == nil {
			t.Fatal("expected to find pane a")
		}
		if p.ID != "a" {
			t.Errorf("expected a, got %s", p.ID)
		}
	})

	t.Run("find non-existing pane", func(t *testing.T) {
		p := m.findPane("z")
		if p != nil {
			t.Error("expected nil for non-existing pane")
		}
	})

	t.Run("nil root", func(t *testing.T) {
		m2 := testManager(nil, "")
		p := m2.findPane("a")
		if p != nil {
			t.Error("expected nil with nil root")
		}
	})
}

func TestManagerView(t *testing.T) {
	t.Run("nil root shows initializing", func(t *testing.T) {
		m := testManager(nil, "")
		view := m.View()
		if view.Content != "Initializing..." {
			t.Errorf("expected 'Initializing...', got %q", view.Content)
		}
	})
}

func TestStatusBarHeight(t *testing.T) {
	m := testManager(nil, "")
	m.height = 48

	t.Run("with status bar", func(t *testing.T) {
		m.statusBar = true
		h := m.statusBarHeight()
		if h != 47 {
			t.Errorf("expected 47, got %d", h)
		}
	})

	t.Run("without status bar", func(t *testing.T) {
		m.statusBar = false
		h := m.statusBarHeight()
		if h != 48 {
			t.Errorf("expected 48, got %d", h)
		}
	})
}

func TestReplaceNode(t *testing.T) {
	t.Run("replace root leaf", func(t *testing.T) {
		leaf := testLeaf("x")
		m := testManager(leaf, "x")
		replacement := testLeaf("y")
		m.replaceNode("x", replacement)
		if m.root != replacement {
			t.Error("root should be replaced")
		}
	})

	t.Run("replace child in tree", func(t *testing.T) {
		root, _, _ := buildTwoPane()
		sn := root.(*splitNode)
		m := testManager(root, "a")
		replacement := testLeaf("x")
		m.replaceNode("b", replacement)
		if sn.second != replacement {
			t.Error("second child should be replaced")
		}
	})
}

func TestFindParentSplit(t *testing.T) {
	root, _, _, _ := buildThreePane()
	sn := root.(*splitNode)
	inner := sn.second.(*splitNode)

	t.Run("find parent of inner split", func(t *testing.T) {
		parent, child := findParentSplit(root, inner)
		if parent != sn {
			t.Error("parent should be root split")
		}
		if child != inner {
			t.Error("child should be inner split")
		}
	})

	t.Run("root has no parent", func(t *testing.T) {
		parent, _ := findParentSplit(root, sn)
		// root splitNode's parent is not found since root is itself the search
		if parent != nil {
			t.Error("root should have no parent in itself")
		}
	})

	t.Run("leaf root returns nil", func(t *testing.T) {
		leaf := testLeaf("x")
		parent, _ := findParentSplit(leaf, sn)
		if parent != nil {
			t.Error("leaf root should return nil")
		}
	})
}
