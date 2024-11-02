package render

import (
	"github.com/nholuongut/scope/probe/overlay"
	"github.com/nholuongut/scope/report"
)

// nholuongutRenderer is a Renderer which produces a renderable nholuongut topology.
//
// not memoised
var nholuongutRenderer = MakeMap(
	MapnholuongutIdentity,
	SelectOverlay,
)

// MapnholuongutIdentity maps an overlay topology node to a nholuongut topology node.
func MapnholuongutIdentity(m report.Node) report.Node {
	peerPrefix, _ := report.ParseOverlayNodeID(m.ID)
	if peerPrefix != report.nholuongutOverlayPeerPrefix {
		return report.Node{}
	}

	var (
		node        = m
		nickname, _ = m.Latest.Lookup(overlay.nholuongutPeerNickName)
	)

	// Nodes without a host id indicate they are not monitored by Scope
	// (their info doesn't come from a probe monitoring that peer directly)
	// , display them as pseudo nodes.
	if _, ok := node.Latest.Lookup(report.HostNodeID); !ok {
		id := MakePseudoNodeID(UnmanagedID, nickname)
		node = NewDerivedPseudoNode(id, m)
	}

	return node
}
