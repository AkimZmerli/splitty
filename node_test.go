package splitty

import "testing"

func TestLeafNodeFindLeaf(t *testing.T) {
	leaf := testLeaf("x")

	t.Run("found", func(t *testing.T) {
		found, path := leaf.findLeaf("x")
		if found == nil {
			t.Fatal("expected to find leaf")
		}
		if found != leaf {
			t.Error("returned wrong leaf")
		}
		if len(path) != 1 || path[0] != leaf {
			t.Error("path should contain only the leaf itself")
		}
	})

	t.Run("not found", func(t *testing.T) {
		found, path := leaf.findLeaf("y")
		if found != nil {
			t.Error("expected nil for non-matching id")
		}
		if path != nil {
			t.Error("expected nil path")
		}
	})

	t.Run("nil pane", func(t *testing.T) {
		empty := &leafNode{pane: nil}
		found, _ := empty.findLeaf("x")
		if found != nil {
			t.Error("expected nil when pane is nil")
		}
	})
}

func TestLeafNodeLeaves(t *testing.T) {
	leaf := testLeaf("x")
	leaves := leaf.leaves()
	if len(leaves) != 1 {
		t.Fatalf("expected 1 leaf, got %d", len(leaves))
	}
	if leaves[0] != leaf {
		t.Error("leaf should return itself")
	}
}

func TestLeafNodeClone(t *testing.T) {
	leaf := testLeaf("x")
	cloned := leaf.clone()
	cl, ok := cloned.(*leafNode)
	if !ok {
		t.Fatal("clone should return a *leafNode")
	}
	if cl == leaf {
		t.Error("clone should return a different pointer")
	}
	if cl.pane != leaf.pane {
		t.Error("clone should share the same pane pointer")
	}
}

func TestSplitNodeFindLeaf(t *testing.T) {
	root, a, _, c := buildThreePane()

	t.Run("find first child", func(t *testing.T) {
		found, path := root.findLeaf("a")
		if found == nil {
			t.Fatal("expected to find leaf a")
		}
		if found != a {
			t.Error("returned wrong leaf")
		}
		// Path should be: root splitNode -> leaf a
		if len(path) < 2 {
			t.Errorf("path too short: %d", len(path))
		}
	})

	t.Run("find nested child", func(t *testing.T) {
		found, path := root.findLeaf("c")
		if found == nil {
			t.Fatal("expected to find leaf c")
		}
		if found != c {
			t.Error("returned wrong leaf")
		}
		// Path: root -> inner split -> leaf c
		if len(path) < 3 {
			t.Errorf("path too short for nested leaf: %d", len(path))
		}
	})

	t.Run("not found", func(t *testing.T) {
		found, _ := root.findLeaf("z")
		if found != nil {
			t.Error("expected nil for non-existing id")
		}
	})
}

func TestSplitNodeLeaves(t *testing.T) {
	root, _, _, _ := buildThreePane()
	leaves := root.leaves()
	if len(leaves) != 3 {
		t.Fatalf("expected 3 leaves, got %d", len(leaves))
	}
	ids := make([]string, len(leaves))
	for i, l := range leaves {
		ids[i] = l.pane.ID
	}
	// In-order: a, b, c
	expected := []string{"a", "b", "c"}
	for i, e := range expected {
		if ids[i] != e {
			t.Errorf("leaf %d: expected %q, got %q", i, e, ids[i])
		}
	}
}

func TestSplitNodeClone(t *testing.T) {
	root, _, _, _ := buildThreePane()
	cloned := root.clone()
	sn, ok := cloned.(*splitNode)
	if !ok {
		t.Fatal("clone should return a *splitNode")
	}
	if sn == root {
		t.Error("clone should return a different pointer")
	}
	// Verify structure is preserved
	origLeaves := root.leaves()
	clonedLeaves := sn.leaves()
	if len(clonedLeaves) != len(origLeaves) {
		t.Fatalf("leaf count mismatch: %d vs %d", len(clonedLeaves), len(origLeaves))
	}
	for i := range origLeaves {
		if clonedLeaves[i] == origLeaves[i] {
			t.Error("cloned leaf should be a different pointer")
		}
		if clonedLeaves[i].pane.ID != origLeaves[i].pane.ID {
			t.Errorf("pane ID mismatch at %d", i)
		}
	}
}

func TestReplaceChild(t *testing.T) {
	a := testLeaf("a")
	b := testLeaf("b")
	sn := &splitNode{dir: Vertical, ratio: 0.5, first: a, second: b}

	replacement := testLeaf("x")

	t.Run("replace first", func(t *testing.T) {
		ok := sn.replaceChild(a, replacement)
		if !ok {
			t.Error("replaceChild should return true")
		}
		if sn.first != replacement {
			t.Error("first child should be replaced")
		}
	})

	t.Run("replace second", func(t *testing.T) {
		ok := sn.replaceChild(b, replacement)
		if !ok {
			t.Error("replaceChild should return true")
		}
		if sn.second != replacement {
			t.Error("second child should be replaced")
		}
	})

	t.Run("not a child", func(t *testing.T) {
		other := testLeaf("other")
		ok := sn.replaceChild(other, replacement)
		if ok {
			t.Error("replaceChild should return false for non-child")
		}
	})
}

func TestSibling(t *testing.T) {
	a := testLeaf("a")
	b := testLeaf("b")
	sn := &splitNode{dir: Vertical, ratio: 0.5, first: a, second: b}

	t.Run("sibling of first is second", func(t *testing.T) {
		sib := sn.sibling(a)
		if sib != b {
			t.Error("sibling of first should be second")
		}
	})

	t.Run("sibling of second is first", func(t *testing.T) {
		sib := sn.sibling(b)
		if sib != a {
			t.Error("sibling of second should be first")
		}
	})

	t.Run("not a child returns nil", func(t *testing.T) {
		other := testLeaf("other")
		sib := sn.sibling(other)
		if sib != nil {
			t.Error("sibling of non-child should be nil")
		}
	})
}

func TestFindParent(t *testing.T) {
	root, _, _, _ := buildThreePane()
	sn := root.(*splitNode)
	inner := sn.second.(*splitNode)

	t.Run("direct child of root", func(t *testing.T) {
		parent, child := findParent(root, "a")
		if parent != sn {
			t.Error("parent of a should be root split")
		}
		if child == nil {
			t.Error("child should not be nil")
		}
	})

	t.Run("nested child", func(t *testing.T) {
		parent, child := findParent(root, "b")
		if parent != inner {
			t.Error("parent of b should be inner split")
		}
		if child == nil {
			t.Error("child should not be nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		parent, _ := findParent(root, "z")
		if parent != nil {
			t.Error("expected nil parent for non-existing id")
		}
	})

	t.Run("leaf root", func(t *testing.T) {
		leaf := testLeaf("x")
		parent, _ := findParent(leaf, "x")
		if parent != nil {
			t.Error("leaf root has no parent")
		}
	})
}
