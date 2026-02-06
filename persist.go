package splitty

import (
	"encoding/json"
	"os"
)

// layoutData is the JSON-serializable representation of the layout tree.
type layoutData struct {
	Version int       `json:"version"`
	Tree    *nodeData `json:"tree"`
}

type nodeData struct {
	Type      string    `json:"type"`               // "leaf" or "split"
	Direction string    `json:"direction,omitempty"` // "vertical" or "horizontal"
	Ratio     float64   `json:"ratio,omitempty"`
	Title     string    `json:"title,omitempty"`
	CWD       string    `json:"cwd,omitempty"`
	Shell     string    `json:"shell,omitempty"`
	First     *nodeData `json:"first,omitempty"`
	Second    *nodeData `json:"second,omitempty"`
}

// SaveLayout serializes the current layout tree to a JSON file.
// Only the tree structure and pane metadata are saved (not terminal content).
func (m *Manager) SaveLayout(path string) error {
	if m.root == nil {
		return nil
	}

	data := layoutData{
		Version: 1,
		Tree:    serializeNode(m.root, m.shell),
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0o644)
}

// LoadLayout restores the layout from a JSON file.
// New PTY sessions are spawned for each pane with the saved working directories.
func (m *Manager) LoadLayout(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var data layoutData
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if data.Tree == nil {
		return nil
	}

	newRoot := deserializeNode(data.Tree, m.shell, m.env, m.scrollbackLines)
	if newRoot == nil {
		return nil
	}

	// Close all existing panes
	if m.root != nil {
		for _, leaf := range m.root.leaves() {
			leaf.pane.close()
		}
	}

	m.root = newRoot
	m.zoomed = false
	m.zoomedID = ""

	// Focus first pane
	leaves := m.root.leaves()
	if len(leaves) > 0 {
		m.focusedID = leaves[0].pane.ID
	}

	m.layoutAll()

	m.log.info("layout loaded", "path", path, "panes", len(leaves))
	return nil
}

func serializeNode(n node, shell string) *nodeData {
	switch n := n.(type) {
	case *leafNode:
		return &nodeData{
			Type:  "leaf",
			Title: n.pane.Title,
			CWD:   n.pane.CWD,
			Shell: shell,
		}
	case *splitNode:
		return &nodeData{
			Type:      "split",
			Direction: n.dir.String(),
			Ratio:     n.ratio,
			First:     serializeNode(n.first, shell),
			Second:    serializeNode(n.second, shell),
		}
	}
	return nil
}

func deserializeNode(data *nodeData, shell string, env []string, scrollbackSize int) node {
	if data == nil {
		return nil
	}

	switch data.Type {
	case "leaf":
		s := shell
		if data.Shell != "" {
			s = data.Shell
		}
		p := newPane(s, env, 80, 24, scrollbackSize) // dimensions will be recalculated by layoutAll
		if data.CWD != "" {
			p.CWD = data.CWD
		}
		if data.Title != "" {
			p.Title = data.Title
		}
		return &leafNode{pane: p}

	case "split":
		dir := Vertical
		if data.Direction == "horizontal" {
			dir = Horizontal
		}
		first := deserializeNode(data.First, shell, env, scrollbackSize)
		second := deserializeNode(data.Second, shell, env, scrollbackSize)
		if first == nil || second == nil {
			return first // or second, whichever is non-nil
		}
		return &splitNode{
			dir:    dir,
			ratio:  data.Ratio,
			first:  first,
			second: second,
		}
	}
	return nil
}
