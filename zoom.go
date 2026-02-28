package splitty

import tea "charm.land/bubbletea/v2"

// Zoom maximizes the focused pane to fill the entire screen.
func (m *Manager) Zoom() tea.Cmd {
	if m.root == nil || m.zoomed {
		return nil
	}
	m.zoomed = true
	m.zoomedID = m.focusedID

	// Resize the zoomed pane to full size
	p := m.findPane(m.zoomedID)
	if p != nil {
		h := m.statusBarHeight()
		p.resize(m.width-2, h-2)
	}

	m.log.debug("zoomed pane", "id", m.zoomedID)
	return nil
}

// Unzoom restores the layout to its pre-zoom state.
func (m *Manager) Unzoom() tea.Cmd {
	if !m.zoomed {
		return nil
	}
	m.zoomed = false
	m.zoomedID = ""

	// Recalculate all pane sizes
	m.layoutAll()

	m.log.debug("unzoomed")
	return nil
}

// IsZoomed returns true if a pane is currently maximized.
func (m *Manager) IsZoomed() bool {
	return m.zoomed
}

// toggleZoom toggles zoom on the focused pane.
func (m *Manager) toggleZoom() (tea.Model, tea.Cmd) {
	if m.zoomed {
		cmd := m.Unzoom()
		return m, cmd
	}
	cmd := m.Zoom()
	return m, cmd
}

// Swap exchanges the focused pane with its sibling in the tree.
func (m *Manager) Swap() {
	if m.root == nil || m.zoomed {
		return
	}

	parent, _ := findParent(m.root, m.focusedID)
	if parent == nil {
		return
	}

	// Swap first and second children
	parent.first, parent.second = parent.second, parent.first

	m.layoutAll()
	m.log.debug("swapped panes")
}
