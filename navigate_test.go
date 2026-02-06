package splitty

import "testing"

func TestFocusPane(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	t.Run("focus existing pane", func(t *testing.T) {
		m.FocusPane("b")
		if m.focusedID != "b" {
			t.Errorf("expected focused b, got %s", m.focusedID)
		}
	})

	t.Run("focus non-existing pane does nothing", func(t *testing.T) {
		m.focusedID = "a"
		m.FocusPane("z")
		if m.focusedID != "a" {
			t.Errorf("focus should not change, got %s", m.focusedID)
		}
	})

	t.Run("nil root does nothing", func(t *testing.T) {
		m2 := testManager(nil, "")
		m2.FocusPane("a") // should not panic
	})
}

func TestFocusedPane(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	t.Run("returns focused pane", func(t *testing.T) {
		p := m.FocusedPane()
		if p == nil {
			t.Fatal("expected non-nil pane")
		}
		if p.ID != "a" {
			t.Errorf("expected pane a, got %s", p.ID)
		}
	})

	t.Run("returns nil for invalid id", func(t *testing.T) {
		m.focusedID = "z"
		p := m.FocusedPane()
		if p != nil {
			t.Error("expected nil for invalid focus")
		}
	})
}

func TestPanes(t *testing.T) {
	root, _, _, _ := buildThreePane()
	m := testManager(root, "a")

	panes := m.Panes()
	if len(panes) != 3 {
		t.Fatalf("expected 3 panes, got %d", len(panes))
	}
	expected := []string{"a", "b", "c"}
	for i, p := range panes {
		if p.ID != expected[i] {
			t.Errorf("pane %d: expected %s, got %s", i, expected[i], p.ID)
		}
	}

	t.Run("nil root returns nil", func(t *testing.T) {
		m2 := testManager(nil, "")
		if m2.Panes() != nil {
			t.Error("expected nil panes for nil root")
		}
	})
}

func TestCycleFocus(t *testing.T) {
	root, _, _, _ := buildThreePane()
	m := testManager(root, "a")

	t.Run("cycle forward", func(t *testing.T) {
		m.focusedID = "a"
		m.cycleFocus(true)
		if m.focusedID != "b" {
			t.Errorf("expected b, got %s", m.focusedID)
		}
		m.cycleFocus(true)
		if m.focusedID != "c" {
			t.Errorf("expected c, got %s", m.focusedID)
		}
		// Wraps around
		m.cycleFocus(true)
		if m.focusedID != "a" {
			t.Errorf("expected wrap to a, got %s", m.focusedID)
		}
	})

	t.Run("cycle backward", func(t *testing.T) {
		m.focusedID = "a"
		m.cycleFocus(false)
		if m.focusedID != "c" {
			t.Errorf("expected wrap to c, got %s", m.focusedID)
		}
		m.cycleFocus(false)
		if m.focusedID != "b" {
			t.Errorf("expected b, got %s", m.focusedID)
		}
	})

	t.Run("single pane does nothing", func(t *testing.T) {
		single := testLeaf("x")
		m2 := testManager(single, "x")
		m2.cycleFocus(true)
		if m2.focusedID != "x" {
			t.Errorf("expected x, got %s", m2.focusedID)
		}
	})

	t.Run("nil root does nothing", func(t *testing.T) {
		m2 := testManager(nil, "")
		m2.cycleFocus(true) // should not panic
	})

	t.Run("invalid focus does nothing", func(t *testing.T) {
		root2, _, _ := buildTwoPane()
		m2 := testManager(root2, "z")
		m2.cycleFocus(true)
		if m2.focusedID != "z" {
			t.Errorf("expected z (unchanged), got %s", m2.focusedID)
		}
	})
}

func TestFocusDirection(t *testing.T) {
	// Use a two-pane vertical split: a | b
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	t.Run("move right from a to b", func(t *testing.T) {
		m.focusedID = "a"
		m.focusDirection(Vertical, true) // right
		if m.focusedID != "b" {
			t.Errorf("expected b, got %s", m.focusedID)
		}
	})

	t.Run("move left from b to a", func(t *testing.T) {
		m.focusedID = "b"
		m.focusDirection(Vertical, false) // left
		if m.focusedID != "a" {
			t.Errorf("expected a, got %s", m.focusedID)
		}
	})

	t.Run("move right from b stays at b", func(t *testing.T) {
		m.focusedID = "b"
		m.focusDirection(Vertical, true) // no sibling to the right
		if m.focusedID != "b" {
			t.Errorf("expected b (no change), got %s", m.focusedID)
		}
	})

	t.Run("wrong axis does nothing", func(t *testing.T) {
		m.focusedID = "a"
		m.focusDirection(Horizontal, true) // no horizontal split
		if m.focusedID != "a" {
			t.Errorf("expected a (no change), got %s", m.focusedID)
		}
	})

	t.Run("nil root does nothing", func(t *testing.T) {
		m2 := testManager(nil, "")
		m2.focusDirection(Vertical, true) // should not panic
	})
}

func TestFocusDirectionThreePane(t *testing.T) {
	// Three pane layout:
	//   split(V)
	//   ├── leaf "a"
	//   └── split(H)
	//       ├── leaf "b"
	//       └── leaf "c"
	root, _, _, _ := buildThreePane()
	m := testManager(root, "a")

	t.Run("right from a goes to b", func(t *testing.T) {
		m.focusedID = "a"
		m.focusDirection(Vertical, true)
		if m.focusedID != "b" {
			t.Errorf("expected b, got %s", m.focusedID)
		}
	})

	t.Run("down from b goes to c", func(t *testing.T) {
		m.focusedID = "b"
		m.focusDirection(Horizontal, true) // down
		if m.focusedID != "c" {
			t.Errorf("expected c, got %s", m.focusedID)
		}
	})

	t.Run("up from c goes to b", func(t *testing.T) {
		m.focusedID = "c"
		m.focusDirection(Horizontal, false) // up
		if m.focusedID != "b" {
			t.Errorf("expected b, got %s", m.focusedID)
		}
	})

	t.Run("left from b goes to a", func(t *testing.T) {
		m.focusedID = "b"
		m.focusDirection(Vertical, false) // left
		if m.focusedID != "a" {
			t.Errorf("expected a, got %s", m.focusedID)
		}
	})

	t.Run("left from c goes to a", func(t *testing.T) {
		m.focusedID = "c"
		m.focusDirection(Vertical, false)
		if m.focusedID != "a" {
			t.Errorf("expected a, got %s", m.focusedID)
		}
	})
}

func TestFirstLeaf(t *testing.T) {
	root, a, _, _ := buildThreePane()

	t.Run("split node", func(t *testing.T) {
		got := firstLeaf(root)
		if got != a {
			t.Errorf("expected leaf a, got %v", got)
		}
	})

	t.Run("leaf node returns itself", func(t *testing.T) {
		leaf := testLeaf("x")
		got := firstLeaf(leaf)
		if got != leaf {
			t.Error("expected leaf to return itself")
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		got := firstLeaf(nil)
		if got != nil {
			t.Error("expected nil")
		}
	})
}

func TestLastLeaf(t *testing.T) {
	root, _, _, c := buildThreePane()

	t.Run("split node", func(t *testing.T) {
		got := lastLeaf(root)
		if got != c {
			t.Errorf("expected leaf c, got %v", got)
		}
	})

	t.Run("leaf node returns itself", func(t *testing.T) {
		leaf := testLeaf("x")
		got := lastLeaf(leaf)
		if got != leaf {
			t.Error("expected leaf to return itself")
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		got := lastLeaf(nil)
		if got != nil {
			t.Error("expected nil")
		}
	})
}

func TestContainsNode(t *testing.T) {
	root, a, b, c := buildThreePane()
	sn := root.(*splitNode)
	inner := sn.second

	t.Run("root contains itself", func(t *testing.T) {
		if !containsNode(root, root) {
			t.Error("tree should contain itself")
		}
	})

	t.Run("root contains leaf a", func(t *testing.T) {
		if !containsNode(root, a) {
			t.Error("root should contain leaf a")
		}
	})

	t.Run("root contains leaf c", func(t *testing.T) {
		if !containsNode(root, c) {
			t.Error("root should contain leaf c")
		}
	})

	t.Run("inner contains b", func(t *testing.T) {
		if !containsNode(inner, b) {
			t.Error("inner should contain leaf b")
		}
	})

	t.Run("inner does not contain a", func(t *testing.T) {
		if containsNode(inner, a) {
			t.Error("inner should not contain leaf a")
		}
	})

	t.Run("leaf does not contain different leaf", func(t *testing.T) {
		if containsNode(a, b) {
			t.Error("leaf a should not contain leaf b")
		}
	})
}
