package utils

import (
	"github.com/nholuongut/scope/report"
)

// Prune returns a copy of the Nodes with all information not strictly
// necessary for rendering nodes and edges in the UI cut away.
func Prune(nodes report.Nodes) report.Nodes {
	result := report.Nodes{}
	for id, node := range nodes {
		result[id] = PruneNode(node)
	}
	return result
}

// PruneNode returns a copy of the Node with all information not strictly
// necessary for rendering nodes and edges stripped away. Specifically, that
// means cutting out parts of the Node.
func PruneNode(node report.Node) report.Node {
	prunedChildren := report.MakeNodeSet()
	node.Children.ForEach(func(child report.Node) {
		prunedChildren = prunedChildren.Add(PruneNode(child))
	})
	return report.MakeNode(
		node.ID).
		WithTopology(node.Topology).
		WithAdjacent(node.Adjacency...).
		WithChildren(prunedChildren)
}
