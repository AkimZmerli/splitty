package splitty

import tea "charm.land/bubbletea/v2"

// Split divides the focused pane in the given direction.
// Returns the updated model and a command to start the new pane's PTY.
func (m *Manager) Split(dir Direction) (tea.Model, tea.Cmd) {
	if m.root == nil || m.zoomed {
		return m, nil
	}

	// Check minimum size constraints
	h := m.statusBarHeight()
	if dir == Vertical && m.width/2 < m.minWidth {
		return m, nil
	}
	if dir == Horizontal && h/2 < m.minHeight {
		return m, nil
	}

	focused := m.findPane(m.focusedID)
	if focused == nil {
		return m, nil
	}

	// Create new pane
	newP := newPane(m.shell, m.env, focused.Width/2, focused.Height/2, m.scrollbackLines)

	newLeaf := &leafNode{pane: newP}
	oldLeaf := &leafNode{pane: focused}

	// Create split node
	split := &splitNode{
		dir:    dir,
		ratio:  0.5,
		first:  oldLeaf,
		second: newLeaf,
	}

	// Replace the focused leaf with the new split
	if m.root == nil {
		m.root = split
	} else {
		m.replaceNode(focused.ID, split)
	}

	// Recalculate layout
	m.layoutAll()

	// Start new pane
	if err := newP.start(m.shell, m.env); err != nil {
		m.log.error("failed to start pane", "id", newP.ID, "err", err)
		return m, nil
	}

	m.log.debug("split pane", "parent", focused.ID, "new", newP.ID, "dir", dir)

	return m, tea.Batch(
		m.readPaneCmd(newP.ID),
		func() tea.Msg {
			return PaneSplitMsg{
				ParentID:  focused.ID,
				NewPaneID: newP.ID,
				Direction: dir,
			}
		},
	)
}

// Close closes the currently focused pane.
func (m *Manager) Close() (tea.Model, tea.Cmd) {
	if m.root == nil {
		return m, nil
	}

	// If only one pane, quit
	leaves := m.root.leaves()
	if len(leaves) <= 1 {
		leaves[0].pane.close()
		return m, tea.Quit
	}

	return m.ClosePane(m.focusedID)
}

// ClosePane closes a specific pane by ID.
func (m *Manager) ClosePane(id string) (tea.Model, tea.Cmd) {
	if m.root == nil {
		return m, nil
	}

	// Unzoom first if needed
	if m.zoomed {
		m.zoomed = false
		m.zoomedID = ""
		m.preZoomRoot = nil
	}

	pane := m.findPane(id)
	if pane == nil {
		return m, nil
	}
	pane.close()

	// Find parent and collapse
	parent, child := findParent(m.root, id)
	if parent == nil {
		// This is the root leaf
		return m, tea.Quit
	}

	sibling := parent.sibling(child)

	// Replace parent with sibling in grandparent
	grandparent, _ := m.findParentNode(parent)
	if grandparent == nil {
		// Parent is the root
		m.root = sibling
	} else {
		grandparent.replaceChild(parent, sibling)
	}

	// Move focus to sibling's first leaf
	siblingLeaves := sibling.leaves()
	if len(siblingLeaves) > 0 {
		m.focusedID = siblingLeaves[0].pane.ID
	}

	m.layoutAll()

	return m, func() tea.Msg {
		return PaneClosedMsg{PaneID: id}
	}
}

// closePane is an internal close used for PTY errors.
func (m *Manager) closePane(id string) {
	m.ClosePane(id)
}

// replaceNode replaces the leaf with the given pane ID with a new node.
func (m *Manager) replaceNode(paneID string, replacement node) {
	// Check if root is the target
	if leaf, ok := m.root.(*leafNode); ok && leaf.pane.ID == paneID {
		m.root = replacement
		return
	}

	parent, child := findParent(m.root, paneID)
	if parent != nil {
		parent.replaceChild(child, replacement)
	}
}

// findParentNode finds the parent splitNode of a given splitNode.
func (m *Manager) findParentNode(target *splitNode) (*splitNode, node) {
	return findParentSplit(m.root, target)
}

func findParentSplit(root node, target *splitNode) (*splitNode, node) {
	sn, ok := root.(*splitNode)
	if !ok {
		return nil, nil
	}
	if sn.first == target || sn.second == target {
		return sn, target
	}
	if parent, child := findParentSplit(sn.first, target); parent != nil {
		return parent, child
	}
	return findParentSplit(sn.second, target)
}
