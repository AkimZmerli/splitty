package splitty

import "testing"

func TestLayoutNodeLeaf(t *testing.T) {
	leaf := testLeaf("x")
	layoutNode(leaf, 10, 5, 80, 24)

	p := leaf.pane
	if p.X != 10 {
		t.Errorf("expected X=10, got %d", p.X)
	}
	if p.Y != 5 {
		t.Errorf("expected Y=5, got %d", p.Y)
	}
	// Inner dimensions = width-2, height-2 (border)
	if p.Width != 78 {
		t.Errorf("expected Width=78, got %d", p.Width)
	}
	if p.Height != 22 {
		t.Errorf("expected Height=22, got %d", p.Height)
	}
}

func TestLayoutNodeVerticalSplit(t *testing.T) {
	a := testLeaf("a")
	b := testLeaf("b")
	root := &splitNode{
		dir:    Vertical,
		ratio:  0.5,
		first:  a,
		second: b,
	}

	layoutNode(root, 0, 0, 100, 40)

	// Left pane: x=0, width=50
	if a.pane.X != 0 {
		t.Errorf("a.X: expected 0, got %d", a.pane.X)
	}
	if a.pane.Width != 48 { // 50 - 2 (border)
		t.Errorf("a.Width: expected 48, got %d", a.pane.Width)
	}

	// Right pane: x=50, width=50
	if b.pane.X != 50 {
		t.Errorf("b.X: expected 50, got %d", b.pane.X)
	}
	if b.pane.Width != 48 { // 50 - 2
		t.Errorf("b.Width: expected 48, got %d", b.pane.Width)
	}

	// Both have same Y and height
	if a.pane.Y != 0 || b.pane.Y != 0 {
		t.Error("both panes should start at Y=0")
	}
	if a.pane.Height != 38 || b.pane.Height != 38 {
		t.Errorf("heights: a=%d b=%d, expected 38", a.pane.Height, b.pane.Height)
	}
}

func TestLayoutNodeHorizontalSplit(t *testing.T) {
	a := testLeaf("a")
	b := testLeaf("b")
	root := &splitNode{
		dir:    Horizontal,
		ratio:  0.5,
		first:  a,
		second: b,
	}

	layoutNode(root, 0, 0, 80, 40)

	// Top pane
	if a.pane.Y != 0 {
		t.Errorf("a.Y: expected 0, got %d", a.pane.Y)
	}
	if a.pane.Height != 18 { // 20 - 2
		t.Errorf("a.Height: expected 18, got %d", a.pane.Height)
	}

	// Bottom pane
	if b.pane.Y != 20 {
		t.Errorf("b.Y: expected 20, got %d", b.pane.Y)
	}
	if b.pane.Height != 18 { // 20 - 2
		t.Errorf("b.Height: expected 18, got %d", b.pane.Height)
	}

	// Both should have same X and width
	if a.pane.X != 0 || b.pane.X != 0 {
		t.Error("both panes should start at X=0")
	}
}

func TestLayoutNodeMinimumSize(t *testing.T) {
	leaf := testLeaf("x")

	// Width and height below 1 are clamped
	layoutNode(leaf, 0, 0, 0, 0)
	// layoutNode clamps to minimum 1x1, inner is max(1-2, 1) = 1
	if leaf.pane.Width < 1 {
		t.Errorf("Width should be at least 1, got %d", leaf.pane.Width)
	}
	if leaf.pane.Height < 1 {
		t.Errorf("Height should be at least 1, got %d", leaf.pane.Height)
	}
}

func TestLayoutNodeNilPane(t *testing.T) {
	leaf := &leafNode{pane: nil}
	// Should not panic
	layoutNode(leaf, 0, 0, 80, 24)
}

func TestLayoutNodeAsymmetricRatio(t *testing.T) {
	a := testLeaf("a")
	b := testLeaf("b")
	root := &splitNode{
		dir:    Vertical,
		ratio:  0.75,
		first:  a,
		second: b,
	}

	layoutNode(root, 0, 0, 100, 40)

	// a gets 75% = 75 cols, b gets 25 cols
	expectedLeftW := int(float64(100) * 0.75)
	if a.pane.X != 0 {
		t.Errorf("a.X: expected 0, got %d", a.pane.X)
	}
	if b.pane.X != expectedLeftW {
		t.Errorf("b.X: expected %d, got %d", expectedLeftW, b.pane.X)
	}
}

func TestManagerResize(t *testing.T) {
	root, _, _ := buildTwoPane()
	sn := root.(*splitNode)
	m := testManager(root, "a")

	origRatio := sn.ratio

	t.Run("resize increases ratio", func(t *testing.T) {
		m.Resize(Vertical, 0.1)
		if sn.ratio <= origRatio {
			t.Errorf("ratio should increase, got %f", sn.ratio)
		}
	})

	t.Run("resize decreases ratio", func(t *testing.T) {
		sn.ratio = 0.5
		m.Resize(Vertical, -0.1)
		if sn.ratio >= 0.5 {
			t.Errorf("ratio should decrease, got %f", sn.ratio)
		}
	})

	t.Run("ratio clamped at minimum 0.1", func(t *testing.T) {
		sn.ratio = 0.15
		m.Resize(Vertical, -0.5)
		if sn.ratio < 0.1 {
			t.Errorf("ratio should be clamped at 0.1, got %f", sn.ratio)
		}
	})

	t.Run("ratio clamped at maximum 0.9", func(t *testing.T) {
		sn.ratio = 0.85
		m.Resize(Vertical, 0.5)
		if sn.ratio > 0.9 {
			t.Errorf("ratio should be clamped at 0.9, got %f", sn.ratio)
		}
	})

	t.Run("wrong axis does nothing", func(t *testing.T) {
		sn.ratio = 0.5
		m.Resize(Horizontal, 0.1) // tree is vertical only
		if sn.ratio != 0.5 {
			t.Errorf("ratio should not change, got %f", sn.ratio)
		}
	})

	t.Run("nil root does nothing", func(t *testing.T) {
		m2 := testManager(nil, "")
		m2.Resize(Vertical, 0.1) // should not panic
	})

	t.Run("zoomed does nothing", func(t *testing.T) {
		sn.ratio = 0.5
		m.zoomed = true
		m.Resize(Vertical, 0.1)
		if sn.ratio != 0.5 {
			t.Errorf("ratio should not change when zoomed, got %f", sn.ratio)
		}
		m.zoomed = false
	})
}
