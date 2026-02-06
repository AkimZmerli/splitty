package splitty

import "github.com/AkimZmerli/splitty/terminal"

// testPane creates a Pane for testing without spawning a real PTY.
func testPane(id string) *Pane {
	return &Pane{
		ID:     id,
		Title:  "test",
		CWD:    "/tmp",
		Width:  80,
		Height: 24,
		screen: terminal.NewScreen(80, 24),
	}
}

// testLeaf creates a leafNode wrapping a test pane.
func testLeaf(id string) *leafNode {
	return &leafNode{pane: testPane(id)}
}

// testManager creates a Manager with a pre-built tree for testing.
// No PTYs are started; this is purely for testing tree logic.
func testManager(root node, focusedID string) *Manager {
	m := &Manager{
		shell:     "/bin/sh",
		theme:     DefaultTheme,
		keyMap:    DefaultKeyMap(),
		minWidth:  10,
		minHeight: 3,
		statusBar: true,
		width:     160,
		height:    48,
		root:      root,
		focusedID: focusedID,
		log:       newLogger(nil),
		menu:      newContextMenu(),
	}
	return m
}

// buildTwoPane creates a vertical split with panes "a" and "b".
//
//	split(V)
//	├── leaf "a"
//	└── leaf "b"
func buildTwoPane() (node, *leafNode, *leafNode) {
	a := testLeaf("a")
	b := testLeaf("b")
	root := &splitNode{
		dir:    Vertical,
		ratio:  0.5,
		first:  a,
		second: b,
	}
	return root, a, b
}

// buildThreePane creates a vertical split where the second child is a
// horizontal split, yielding panes "a", "b", "c".
//
//	split(V)
//	├── leaf "a"
//	└── split(H)
//	    ├── leaf "b"
//	    └── leaf "c"
func buildThreePane() (node, *leafNode, *leafNode, *leafNode) {
	a := testLeaf("a")
	b := testLeaf("b")
	c := testLeaf("c")
	inner := &splitNode{
		dir:    Horizontal,
		ratio:  0.5,
		first:  b,
		second: c,
	}
	root := &splitNode{
		dir:    Vertical,
		ratio:  0.5,
		first:  a,
		second: inner,
	}
	return root, a, b, c
}

// buildQuadPane creates a 2x2 grid with panes "a", "b", "c", "d".
//
//	split(V)
//	├── split(H)
//	│   ├── leaf "a"
//	│   └── leaf "b"
//	└── split(H)
//	    ├── leaf "c"
//	    └── leaf "d"
func buildQuadPane() (node, *leafNode, *leafNode, *leafNode, *leafNode) {
	a := testLeaf("a")
	b := testLeaf("b")
	c := testLeaf("c")
	d := testLeaf("d")
	left := &splitNode{
		dir:    Horizontal,
		ratio:  0.5,
		first:  a,
		second: b,
	}
	right := &splitNode{
		dir:    Horizontal,
		ratio:  0.5,
		first:  c,
		second: d,
	}
	root := &splitNode{
		dir:    Vertical,
		ratio:  0.5,
		first:  left,
		second: right,
	}
	return root, a, b, c, d
}
