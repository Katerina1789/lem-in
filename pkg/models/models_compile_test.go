package models

import "testing"

func TestCompilation(t *testing.T) {
	var _ FlowGraph
	var _ FlowNode
	var _ FlowEdge
	t.Log("Models compiled successfully")
}
