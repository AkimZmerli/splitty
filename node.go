package splitty

// node is the binary tree interface for layout management.
// Each node is either a leafNode (terminal pane) or a splitNode (container).
type node interface {
	findLeaf(id string) (*leafNode, []node)
	leaves() []*leafNode
	clone() node
}

// leafNode is a terminal leaf in the layout tree.
type leafNode struct {
	pane *Pane
}

func (l *leafNode) findLeaf(id string) (*leafNode, []node) {
	if l.pane != nil && l.pane.ID == id {
		return l, []node{l}
	}
	return nil, nil
}

func (l *leafNode) leaves() []*leafNode {
	return []*leafNode{l}
}

func (l *leafNode) clone() node {
	return &leafNode{pane: l.pane}
}

// splitNode is an internal node that divides space between two children.
type splitNode struct {
	dir    Direction
	ratio  float64 // 0.0-1.0: first child gets ratio * available space
	first  node
	second node
}

func (s *splitNode) findLeaf(id string) (*leafNode, []node) {
	if leaf, path := s.first.findLeaf(id); leaf != nil {
		return leaf, append([]node{s}, path...)
	}
	if leaf, path := s.second.findLeaf(id); leaf != nil {
		return leaf, append([]node{s}, path...)
	}
	return nil, nil
}

func (s *splitNode) leaves() []*leafNode {
	return append(s.first.leaves(), s.second.leaves()...)
}

func (s *splitNode) clone() node {
	return &splitNode{
		dir:    s.dir,
		ratio:  s.ratio,
		first:  s.first.clone(),
		second: s.second.clone(),
	}
}

// replaceChild replaces a direct child node with a new node.
func (s *splitNode) replaceChild(old, replacement node) bool {
	if s.first == old {
		s.first = replacement
		return true
	}
	if s.second == old {
		s.second = replacement
		return true
	}
	return false
}

// sibling returns the other child when given one child.
func (s *splitNode) sibling(child node) node {
	if s.first == child {
		return s.second
	}
	if s.second == child {
		return s.first
	}
	return nil
}

// findParent finds the parent splitNode of a node with the given leaf id.
func findParent(root node, id string) (*splitNode, node) {
	sn, ok := root.(*splitNode)
	if !ok {
		return nil, nil
	}

	// Check direct children
	if leaf, ok := sn.first.(*leafNode); ok && leaf.pane != nil && leaf.pane.ID == id {
		return sn, sn.first
	}
	if leaf, ok := sn.second.(*leafNode); ok && leaf.pane != nil && leaf.pane.ID == id {
		return sn, sn.second
	}

	// Recurse
	if parent, child := findParent(sn.first, id); parent != nil {
		return parent, child
	}
	return findParent(sn.second, id)
}
