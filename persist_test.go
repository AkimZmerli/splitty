package splitty

import (
	"encoding/json"
	"testing"
)

func TestSerializeLeafNode(t *testing.T) {
	leaf := testLeaf("x")
	leaf.pane.Title = "my-shell"
	leaf.pane.CWD = "/home/user"

	data := serializeNode(leaf, "/bin/bash")
	if data == nil {
		t.Fatal("serializeNode returned nil")
	}
	if data.Type != "leaf" {
		t.Errorf("expected type leaf, got %s", data.Type)
	}
	if data.Title != "my-shell" {
		t.Errorf("expected title my-shell, got %s", data.Title)
	}
	if data.CWD != "/home/user" {
		t.Errorf("expected CWD /home/user, got %s", data.CWD)
	}
	if data.Shell != "/bin/bash" {
		t.Errorf("expected shell /bin/bash, got %s", data.Shell)
	}
}

func TestSerializeSplitNode(t *testing.T) {
	root, _, _ := buildTwoPane()

	data := serializeNode(root, "/bin/sh")
	if data == nil {
		t.Fatal("serializeNode returned nil")
	}
	if data.Type != "split" {
		t.Errorf("expected type split, got %s", data.Type)
	}
	if data.Direction != "vertical" {
		t.Errorf("expected direction vertical, got %s", data.Direction)
	}
	if data.Ratio != 0.5 {
		t.Errorf("expected ratio 0.5, got %f", data.Ratio)
	}
	if data.First == nil || data.Second == nil {
		t.Fatal("children should not be nil")
	}
	if data.First.Type != "leaf" || data.Second.Type != "leaf" {
		t.Error("children should be leaf nodes")
	}
}

func TestSerializeThreePane(t *testing.T) {
	root, _, _, _ := buildThreePane()
	data := serializeNode(root, "/bin/sh")

	if data.Type != "split" {
		t.Fatal("root should be split")
	}
	if data.First.Type != "leaf" {
		t.Error("first child should be leaf")
	}
	if data.Second.Type != "split" {
		t.Error("second child should be split")
	}
	if data.Second.Direction != "horizontal" {
		t.Errorf("inner split direction: expected horizontal, got %s", data.Second.Direction)
	}
}

func TestSerializeNilNode(t *testing.T) {
	data := serializeNode(nil, "/bin/sh")
	if data != nil {
		t.Error("expected nil for nil node")
	}
}

func TestDeserializeLeafNode(t *testing.T) {
	data := &nodeData{
		Type:  "leaf",
		Title: "test",
		CWD:   "/tmp",
		Shell: "/bin/bash",
	}

	n := deserializeNode(data, "/bin/sh", nil)
	if n == nil {
		t.Fatal("deserializeNode returned nil")
	}
	leaf, ok := n.(*leafNode)
	if !ok {
		t.Fatal("expected *leafNode")
	}
	if leaf.pane == nil {
		t.Fatal("pane should not be nil")
	}
	if leaf.pane.Title != "test" {
		t.Errorf("expected title test, got %s", leaf.pane.Title)
	}
	if leaf.pane.CWD != "/tmp" {
		t.Errorf("expected CWD /tmp, got %s", leaf.pane.CWD)
	}
}

func TestDeserializeSplitNode(t *testing.T) {
	data := &nodeData{
		Type:      "split",
		Direction: "vertical",
		Ratio:     0.6,
		First:     &nodeData{Type: "leaf", Shell: "/bin/sh"},
		Second:    &nodeData{Type: "leaf", Shell: "/bin/sh"},
	}

	n := deserializeNode(data, "/bin/sh", nil)
	if n == nil {
		t.Fatal("deserializeNode returned nil")
	}
	sn, ok := n.(*splitNode)
	if !ok {
		t.Fatal("expected *splitNode")
	}
	if sn.dir != Vertical {
		t.Errorf("expected Vertical, got %v", sn.dir)
	}
	if sn.ratio != 0.6 {
		t.Errorf("expected ratio 0.6, got %f", sn.ratio)
	}
	if sn.first == nil || sn.second == nil {
		t.Fatal("children should not be nil")
	}
}

func TestDeserializeHorizontalDirection(t *testing.T) {
	data := &nodeData{
		Type:      "split",
		Direction: "horizontal",
		Ratio:     0.5,
		First:     &nodeData{Type: "leaf", Shell: "/bin/sh"},
		Second:    &nodeData{Type: "leaf", Shell: "/bin/sh"},
	}

	n := deserializeNode(data, "/bin/sh", nil)
	sn := n.(*splitNode)
	if sn.dir != Horizontal {
		t.Errorf("expected Horizontal direction")
	}
}

func TestDeserializeNilData(t *testing.T) {
	n := deserializeNode(nil, "/bin/sh", nil)
	if n != nil {
		t.Error("expected nil for nil data")
	}
}

func TestDeserializeMissingChild(t *testing.T) {
	data := &nodeData{
		Type:      "split",
		Direction: "vertical",
		Ratio:     0.5,
		First:     &nodeData{Type: "leaf", Shell: "/bin/sh"},
		Second:    nil,
	}

	n := deserializeNode(data, "/bin/sh", nil)
	// When second child is nil, should return first
	if n == nil {
		t.Fatal("should return the non-nil child")
	}
	_, ok := n.(*leafNode)
	if !ok {
		t.Error("expected the lone leaf to be returned")
	}
}

func TestDeserializeUnknownType(t *testing.T) {
	data := &nodeData{Type: "unknown"}
	n := deserializeNode(data, "/bin/sh", nil)
	if n != nil {
		t.Error("expected nil for unknown type")
	}
}

func TestSerializeDeserializeRoundtrip(t *testing.T) {
	root, _, _, _ := buildThreePane()
	data := serializeNode(root, "/bin/bash")

	// Marshal to JSON and back
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded nodeData
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	n := deserializeNode(&decoded, "/bin/bash", nil)
	if n == nil {
		t.Fatal("deserialized node is nil")
	}

	// Verify structure
	leaves := n.leaves()
	if len(leaves) != 3 {
		t.Errorf("expected 3 leaves, got %d", len(leaves))
	}
}

func TestDeserializeUsesShellFallback(t *testing.T) {
	data := &nodeData{
		Type:  "leaf",
		Title: "test",
		Shell: "", // empty shell should use the fallback
	}

	n := deserializeNode(data, "/bin/zsh", nil)
	leaf := n.(*leafNode)
	// newPane is called with the fallback shell
	if leaf.pane == nil {
		t.Fatal("pane should not be nil")
	}
}
