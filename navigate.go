package splitty

// Focus moves focus to the pane in the given direction relative to current focus.
func (m *Manager) Focus(dir Direction) {
	m.focusDirection(dir, true)
}

// FocusPane sets focus to the pane with the given ID.
func (m *Manager) FocusPane(id string) {
	if m.root == nil {
		return
	}
	if leaf, _ := m.root.findLeaf(id); leaf != nil {
		m.focusedID = id
	}
}

// FocusedPane returns the currently focused pane, or nil.
func (m *Manager) FocusedPane() *Pane {
	return m.findPane(m.focusedID)
}

// Panes returns all panes in the layout in left-to-right, top-to-bottom order.
func (m *Manager) Panes() []*Pane {
	if m.root == nil {
		return nil
	}
	leaves := m.root.leaves()
	panes := make([]*Pane, len(leaves))
	for i, l := range leaves {
		panes[i] = l.pane
	}
	return panes
}

// focusDirection implements directional navigation using the two-phase algorithm:
// 1. Walk up from focused leaf to find a split node matching the axis
// 2. Walk down into the opposite subtree to the nearest edge leaf
func (m *Manager) focusDirection(axis Direction, forward bool) {
	if m.root == nil {
		return
	}

	leaf, path := m.root.findLeaf(m.focusedID)
	if leaf == nil || len(path) < 2 {
		return
	}

	// Phase 1: walk up to find a matching split
	for i := len(path) - 1; i >= 0; i-- {
		sn, ok := path[i].(*splitNode)
		if !ok {
			continue
		}
		if sn.dir != axis {
			continue
		}

		// Check if we're on the correct side for this direction
		var childInPath node
		if i+1 < len(path) {
			childInPath = path[i+1]
		} else {
			childInPath = leaf
		}

		isInFirst := sn.first == childInPath || containsNode(sn.first, childInPath)

		// forward=true means go to second child (right/down)
		// forward=false means go to first child (left/up)
		if forward && isInFirst {
			// Enter second subtree, pick the nearest edge (first leaf)
			target := firstLeaf(sn.second)
			if target != nil {
				m.focusedID = target.pane.ID
			}
			return
		}
		if !forward && !isInFirst {
			// Enter first subtree, pick the farthest edge (last leaf)
			target := lastLeaf(sn.first)
			if target != nil {
				m.focusedID = target.pane.ID
			}
			return
		}
	}
}

// cycleFocus moves focus to the next or previous pane in order.
func (m *Manager) cycleFocus(forward bool) {
	if m.root == nil {
		return
	}

	leaves := m.root.leaves()
	if len(leaves) <= 1 {
		return
	}

	// Find current index
	idx := -1
	for i, l := range leaves {
		if l.pane.ID == m.focusedID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return
	}

	if forward {
		idx = (idx + 1) % len(leaves)
	} else {
		idx = (idx - 1 + len(leaves)) % len(leaves)
	}
	m.focusedID = leaves[idx].pane.ID
}

// firstLeaf returns the leftmost/topmost leaf in a subtree.
func firstLeaf(n node) *leafNode {
	switch n := n.(type) {
	case *leafNode:
		return n
	case *splitNode:
		return firstLeaf(n.first)
	}
	return nil
}

// lastLeaf returns the rightmost/bottommost leaf in a subtree.
func lastLeaf(n node) *leafNode {
	switch n := n.(type) {
	case *leafNode:
		return n
	case *splitNode:
		return lastLeaf(n.second)
	}
	return nil
}

// containsNode checks if a tree contains a specific node.
func containsNode(tree, target node) bool {
	if tree == target {
		return true
	}
	sn, ok := tree.(*splitNode)
	if !ok {
		return false
	}
	return containsNode(sn.first, target) || containsNode(sn.second, target)
}
