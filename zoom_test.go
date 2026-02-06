package splitty

import "testing"

func TestZoom(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	t.Run("zoom sets state", func(t *testing.T) {
		m.Zoom()
		if !m.zoomed {
			t.Error("zoomed should be true")
		}
		if m.zoomedID != "a" {
			t.Errorf("zoomedID: expected a, got %s", m.zoomedID)
		}
	})

	t.Run("zoom when already zoomed does nothing", func(t *testing.T) {
		m.zoomed = true
		m.zoomedID = "a"
		cmd := m.Zoom()
		if cmd != nil {
			t.Error("expected nil cmd when already zoomed")
		}
	})

	t.Run("zoom with nil root does nothing", func(t *testing.T) {
		m2 := testManager(nil, "")
		cmd := m2.Zoom()
		if cmd != nil {
			t.Error("expected nil cmd with nil root")
		}
		if m2.zoomed {
			t.Error("should not be zoomed")
		}
	})
}

func TestUnzoom(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")
	m.zoomed = true
	m.zoomedID = "a"

	t.Run("unzoom clears state", func(t *testing.T) {
		m.Unzoom()
		if m.zoomed {
			t.Error("zoomed should be false")
		}
		if m.zoomedID != "" {
			t.Errorf("zoomedID should be empty, got %s", m.zoomedID)
		}
	})

	t.Run("unzoom when not zoomed does nothing", func(t *testing.T) {
		m.zoomed = false
		cmd := m.Unzoom()
		if cmd != nil {
			t.Error("expected nil cmd when not zoomed")
		}
	})
}

func TestIsZoomed(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	if m.IsZoomed() {
		t.Error("should not be zoomed initially")
	}
	m.zoomed = true
	if !m.IsZoomed() {
		t.Error("should be zoomed")
	}
}

func TestToggleZoom(t *testing.T) {
	root, _, _ := buildTwoPane()
	m := testManager(root, "a")

	// Toggle on
	m.toggleZoom()
	if !m.zoomed {
		t.Error("should be zoomed after first toggle")
	}

	// Toggle off
	m.toggleZoom()
	if m.zoomed {
		t.Error("should not be zoomed after second toggle")
	}
}

func TestSwap(t *testing.T) {
	root, _, _ := buildTwoPane()
	sn := root.(*splitNode)
	m := testManager(root, "a")

	origFirst := sn.first
	origSecond := sn.second

	m.Swap()

	if sn.first != origSecond {
		t.Error("first should now be the original second")
	}
	if sn.second != origFirst {
		t.Error("second should now be the original first")
	}
}

func TestSwapNilRoot(t *testing.T) {
	m := testManager(nil, "")
	m.Swap() // should not panic
}

func TestSwapZoomed(t *testing.T) {
	root, _, _ := buildTwoPane()
	sn := root.(*splitNode)
	m := testManager(root, "a")
	m.zoomed = true

	origFirst := sn.first
	m.Swap()
	if sn.first != origFirst {
		t.Error("swap should be no-op when zoomed")
	}
}

func TestSwapSinglePane(t *testing.T) {
	leaf := testLeaf("x")
	m := testManager(leaf, "x")
	m.Swap() // should not panic, no parent to swap
}
