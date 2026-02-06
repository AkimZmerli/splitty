package splitty

// SetBroadcast enables or disables broadcast mode.
// When broadcasting, keyboard input is sent to all panes simultaneously.
func (m *Manager) SetBroadcast(enabled bool) {
	m.broadcasting = enabled
	m.log.debug("broadcast mode", "enabled", enabled)
}

// IsBroadcasting returns true if broadcast mode is active.
func (m *Manager) IsBroadcasting() bool {
	return m.broadcasting
}

// SendInput sends raw bytes to the focused pane, or all panes if broadcasting.
func (m *Manager) SendInput(data []byte) {
	if m.root == nil {
		return
	}

	if m.broadcasting {
		for _, leaf := range m.root.leaves() {
			leaf.pane.write(data)
		}
		return
	}

	p := m.findPane(m.focusedID)
	if p != nil {
		p.write(data)
	}
}
