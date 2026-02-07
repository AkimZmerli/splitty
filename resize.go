package splitty

// clampRatio clamps a split ratio to the valid range [0.1, 0.9].
func clampRatio(r float64) float64 {
	if r < 0.1 {
		return 0.1
	}
	if r > 0.9 {
		return 0.9
	}
	return r
}

// Resize adjusts the split ratio of the nearest ancestor split node
// along the given axis by delta (positive = grow first child).
func (m *Manager) Resize(dir Direction, delta float64) {
	if m.root == nil || m.zoomed {
		return
	}

	_, path := m.root.findLeaf(m.focusedID)
	if path == nil {
		return
	}

	// Walk up to find the nearest split matching the direction
	for i := len(path) - 1; i >= 0; i-- {
		sn, ok := path[i].(*splitNode)
		if !ok {
			continue
		}
		if sn.dir != dir {
			continue
		}

		sn.ratio += delta
		if sn.ratio < 0.1 {
			sn.ratio = 0.1
		}
		if sn.ratio > 0.9 {
			sn.ratio = 0.9
		}

		m.layoutAll()
		return
	}
}

// findBorder checks if the given screen coordinates are on a split border.
// Returns the split node and its direction if found, nil otherwise.
func findBorder(root node, mouseX, mouseY, totalWidth, totalHeight int) *splitNode {
	return findBorderAt(root, 0, 0, totalWidth, totalHeight, mouseX, mouseY)
}

func findBorderAt(n node, x, y, width, height, mx, my int) *splitNode {
	sn, ok := n.(*splitNode)
	if !ok {
		return nil
	}

	if sn.dir == Vertical {
		leftW := int(float64(width) * sn.ratio)
		borderX := x + leftW
		// Hit if mouse is within 1 cell of the border column
		if mx >= borderX-1 && mx <= borderX && my >= y && my < y+height {
			return sn
		}
		// Recurse into children
		rightW := width - leftW
		if found := findBorderAt(sn.first, x, y, leftW, height, mx, my); found != nil {
			return found
		}
		return findBorderAt(sn.second, x+leftW, y, rightW, height, mx, my)
	}

	// Horizontal
	topH := int(float64(height) * sn.ratio)
	borderY := y + topH
	if my >= borderY-1 && my <= borderY && mx >= x && mx < x+width {
		return sn
	}
	bottomH := height - topH
	if found := findBorderAt(sn.first, x, y, width, topH, mx, my); found != nil {
		return found
	}
	return findBorderAt(sn.second, x, y+topH, width, bottomH, mx, my)
}

// layoutAll recalculates all pane positions and sizes from the tree.
func (m *Manager) layoutAll() {
	if m.root == nil {
		return
	}
	h := m.statusBarHeight()
	layoutNode(m.root, 0, 0, m.width, h)
}

// layoutNode recursively assigns positions and sizes to all panes.
func layoutNode(n node, x, y, width, height int) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	switch n := n.(type) {
	case *leafNode:
		if n.pane != nil {
			n.pane.X = x
			n.pane.Y = y
			// Account for border (2 chars each side)
			innerW := width - 2
			innerH := height - 2
			if innerW < 1 {
				innerW = 1
			}
			if innerH < 1 {
				innerH = 1
			}
			n.pane.resize(innerW, innerH)
		}

	case *splitNode:
		if n.dir == Vertical {
			leftW := int(float64(width) * n.ratio)
			rightW := width - leftW
			if leftW < 1 {
				leftW = 1
			}
			if rightW < 1 {
				rightW = 1
			}
			layoutNode(n.first, x, y, leftW, height)
			layoutNode(n.second, x+leftW, y, rightW, height)
		} else {
			topH := int(float64(height) * n.ratio)
			bottomH := height - topH
			if topH < 1 {
				topH = 1
			}
			if bottomH < 1 {
				bottomH = 1
			}
			layoutNode(n.first, x, y, width, topH)
			layoutNode(n.second, x, y+topH, width, bottomH)
		}
	}
}
