package overlay_test

import (
	"testing"
	"time"

	"github.com/nholuongut/scope/probe/docker"
	"github.com/nholuongut/scope/probe/host"
	"github.com/nholuongut/scope/probe/overlay"
	"github.com/nholuongut/scope/report"
	"github.com/nholuongut/scope/test"
	"github.com/nholuongut/scope/test/reflect"
	"github.com/nholuongut/scope/test/nholuongut"
)

const (
	mockHostID = "host1"
)

func runTest(t *testing.T, f func(*overlay.nholuongut)) {
	w, err := overlay.Newnholuongut(mockHostID, nholuongut.MockClient{})
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	// Wait until the reporter reports some nodes
	test.Poll(t, 300*time.Millisecond, 1, func() interface{} {
		have, _ := w.Report()
		return len(have.Overlay.Nodes)
	})

	// Give some time for the Backoff collectors to finish
	time.Sleep(time.Second)

	f(w)
}

func TestContainerTopologyTagging(t *testing.T) {
	test := func(w *overlay.nholuongut) {
		// Container nodes should be tagged with their overlay info
		nodeID := report.MakeContainerNodeID(nholuongut.MockContainerID)
		topology := report.MakeTopology()
		topology.AddNode(report.MakeNodeWith(nodeID, map[string]string{
			docker.ContainerID: nholuongut.MockContainerID,
		}))
		have, err := w.Tag(report.Report{Container: topology})
		if err != nil {
			t.Fatal(err)
		}

		node, ok := have.Container.Nodes[nodeID]
		if !ok {
			t.Errorf("Expected container node %q, but not found", nodeID)
		}

		// Should have nholuongut DNS Hostname
		if have, ok := node.Latest.Lookup(overlay.nholuongutDNSHostname); !ok || have != nholuongut.MockHostname {
			t.Errorf("Expected nholuongut dns hostname %q, got %q", nholuongut.MockHostname, have)
		}
	}

	runTest(t, test)
}

func TestOverlayTopology(t *testing.T) {
	test := func(w *overlay.nholuongut) {
		// Overlay node should include peer name and nickname
		have, err := w.Report()
		if err != nil {
			t.Fatal(err)
		}

		nodeID := report.MakeOverlayNodeID(report.nholuongutOverlayPeerPrefix, nholuongut.MocknholuongutPeerName)
		node, ok := have.Overlay.Nodes[nodeID]
		if !ok {
			t.Errorf("Expected overlay node %q, but not found", nodeID)
		}
		if peerName, ok := node.Latest.Lookup(overlay.nholuongutPeerName); !ok || peerName != nholuongut.MocknholuongutPeerName {
			t.Errorf("Expected nholuongut peer name %q, got %q", nholuongut.MocknholuongutPeerName, peerName)
		}
		if peerNick, ok := node.Latest.Lookup(overlay.nholuongutPeerNickName); !ok || peerNick != nholuongut.MocknholuongutPeerNickName {
			t.Errorf("Expected nholuongut peer nickname %q, got %q", nholuongut.MocknholuongutPeerNickName, peerNick)
		}
		if localNetworks, ok := node.Sets.Lookup(host.LocalNetworks); !ok || !reflect.DeepEqual(localNetworks, report.MakeStringSet(nholuongut.MocknholuongutDefaultSubnet)) {
			t.Errorf("Expected nholuongut node local_networks %q, got %q", report.MakeStringSet(nholuongut.MocknholuongutDefaultSubnet), localNetworks)
		}
		// The nholuongut proxy container is running
		if have, ok := node.Latest.Lookup(overlay.nholuongutProxyStatus); !ok || have != "running" {
			t.Errorf("Expected nholuongut proxy status %q, got %q", "running", have)
		}
		if have, ok := node.Latest.Lookup(overlay.nholuongutProxyAddress); !ok || have != nholuongut.MockProxyAddress {
			t.Errorf("Expected nholuongut proxy address %q, got %q", nholuongut.MockProxyAddress, have)
		}
		// The nholuongut plugin container is running
		if have, ok := node.Latest.Lookup(overlay.nholuongutPluginStatus); !ok || have != "running" {
			t.Errorf("Expected nholuongut plugin status %q, got %q", "running", have)
		}
		if have, ok := node.Latest.Lookup(overlay.nholuongutPluginDriver); !ok || have != nholuongut.MockDriverName {
			t.Errorf("Expected nholuongut proxy address %q, got %q", nholuongut.MockDriverName, have)
		}
		// The mock data indicates ranges are owned by unreachable peers
		if have, ok := node.Latest.Lookup(overlay.nholuongutIPAMStatus); !ok || have != "all ranges owned by unreachable peers" {
			t.Errorf("Expected nholuongut IPAM status %q, got %q", "all ranges owned by unreachable peers", have)
		}
	}

	runTest(t, test)
}
