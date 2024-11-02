package overlay

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/nholuongutworks/common/backoff"
	"github.com/nholuongut/scope/common/nholuongut"
	"github.com/nholuongut/scope/probe/docker"
	"github.com/nholuongut/scope/probe/host"
	"github.com/nholuongut/scope/report"
)

// Keys for use in Node
const (
	nholuongutPeerName                          = "nholuongut_peer_name"
	nholuongutPeerNickName                      = "nholuongut_peer_nick_name"
	nholuongutDNSHostname                       = "nholuongut_dns_hostname"
	nholuongutMACAddress                        = "nholuongut_mac_address"
	nholuongutVersion                           = "nholuongut_version"
	nholuongutEncryption                        = "nholuongut_encryption"
	nholuongutProtocol                          = "nholuongut_protocol"
	nholuongutPeerDiscovery                     = "nholuongut_peer_discovery"
	nholuongutTargetCount                       = "nholuongut_target_count"
	nholuongutConnectionCount                   = "nholuongut_connection_count"
	nholuongutPeerCount                         = "nholuongut_peer_count"
	nholuongutTrustedSubnets                    = "nholuongut_trusted_subnet_count"
	nholuongutIPAMTableID                       = "nholuongut_ipam_table"
	nholuongutIPAMStatus                        = "nholuongut_ipam_status"
	nholuongutIPAMRange                         = "nholuongut_ipam_range"
	nholuongutIPAMDefaultSubnet                 = "nholuongut_ipam_default_subnet"
	nholuongutDNSTableID                        = "nholuongut_dns_table"
	nholuongutDNSDomain                         = "nholuongut_dns_domain"
	nholuongutDNSUpstream                       = "nholuongut_dns_upstream"
	nholuongutDNSTTL                            = "nholuongut_dns_ttl"
	nholuongutDNSEntryCount                     = "nholuongut_dns_entry_count"
	nholuongutProxyTableID                      = "nholuongut_proxy_table"
	nholuongutProxyStatus                       = "nholuongut_proxy_status"
	nholuongutProxyAddress                      = "nholuongut_proxy_address"
	nholuongutPluginTableID                     = "nholuongut_plugin_table"
	nholuongutPluginStatus                      = "nholuongut_plugin_status"
	nholuongutPluginDriver                      = "nholuongut_plugin_driver"
	nholuongutConnectionsConnection             = "nholuongut_connection_connection"
	nholuongutConnectionsState                  = "nholuongut_connection_state"
	nholuongutConnectionsInfo                   = "nholuongut_connection_info"
	nholuongutConnectionsTablePrefix            = "nholuongut_connections_table_"
	nholuongutConnectionsMulticolumnTablePrefix = "nholuongut_connections_multicolumn_table_"
)

var (
	containerNotRunningRE = regexp.MustCompile(`Container .* is not running\n`)

	containerMetadata = report.MetadataTemplates{
		nholuongutMACAddress:  {ID: nholuongutMACAddress, Label: "nholuongut MAC", From: report.FromLatest, Priority: 17},
		nholuongutDNSHostname: {ID: nholuongutDNSHostname, Label: "nholuongut DNS Name", From: report.FromLatest, Priority: 18},
	}

	nholuongutMetadata = report.MetadataTemplates{
		nholuongutVersion:         {ID: nholuongutVersion, Label: "Version", From: report.FromLatest, Priority: 1},
		nholuongutProtocol:        {ID: nholuongutProtocol, Label: "Protocol", From: report.FromLatest, Priority: 2},
		nholuongutPeerName:        {ID: nholuongutPeerName, Label: "Name", From: report.FromLatest, Priority: 3},
		nholuongutEncryption:      {ID: nholuongutEncryption, Label: "Encryption", From: report.FromLatest, Priority: 4},
		nholuongutPeerDiscovery:   {ID: nholuongutPeerDiscovery, Label: "Peer discovery", From: report.FromLatest, Priority: 5},
		nholuongutTargetCount:     {ID: nholuongutTargetCount, Label: "Targets", From: report.FromLatest, Priority: 6},
		nholuongutConnectionCount: {ID: nholuongutConnectionCount, Label: "Connections", From: report.FromLatest, Priority: 8},
		nholuongutPeerCount:       {ID: nholuongutPeerCount, Label: "Peers", From: report.FromLatest, Priority: 7},
		nholuongutTrustedSubnets:  {ID: nholuongutTrustedSubnets, Label: "Trusted subnets", From: report.FromSets, Priority: 9},
	}

	nholuongutTableTemplates = report.TableTemplates{
		nholuongutIPAMTableID: {
			ID:    nholuongutIPAMTableID,
			Label: "IPAM",
			Type:  report.PropertyListType,
			FixedRows: map[string]string{
				nholuongutIPAMStatus:        "Status",
				nholuongutIPAMRange:         "Range",
				nholuongutIPAMDefaultSubnet: "Default subnet",
			},
		},
		nholuongutDNSTableID: {
			ID:    nholuongutDNSTableID,
			Label: "DNS",
			Type:  report.PropertyListType,
			FixedRows: map[string]string{
				nholuongutDNSDomain:     "Domain",
				nholuongutDNSUpstream:   "Upstream",
				nholuongutDNSTTL:        "TTL",
				nholuongutDNSEntryCount: "Entries",
			},
		},
		nholuongutProxyTableID: {
			ID:    nholuongutProxyTableID,
			Label: "Proxy",
			Type:  report.PropertyListType,
			FixedRows: map[string]string{
				nholuongutProxyStatus:  "Status",
				nholuongutProxyAddress: "Address",
			},
		},
		nholuongutPluginTableID: {
			ID:    nholuongutPluginTableID,
			Label: "Plugin",
			Type:  report.PropertyListType,
			FixedRows: map[string]string{
				nholuongutPluginStatus: "Status",
				nholuongutPluginDriver: "Driver name",
			},
		},
		nholuongutConnectionsMulticolumnTablePrefix: {
			ID:     nholuongutConnectionsMulticolumnTablePrefix,
			Type:   report.MulticolumnTableType,
			Prefix: nholuongutConnectionsMulticolumnTablePrefix,
			Columns: []report.Column{
				{
					ID:    nholuongutConnectionsConnection,
					Label: "Connections",
				},
				{
					ID:    nholuongutConnectionsState,
					Label: "State",
				},
				{
					ID:    nholuongutConnectionsInfo,
					Label: "Info",
				},
			},
		},
		// Kept for backward-compatibility.
		nholuongutConnectionsTablePrefix: {
			ID:     nholuongutConnectionsTablePrefix,
			Label:  "Connections",
			Type:   report.PropertyListType,
			Prefix: nholuongutConnectionsTablePrefix,
		},
	}
)

// nholuongut represents a single nholuongut router, presumably on the same host
// as the probe. It is both a Reporter and a Tagger: it produces an Overlay
// topology, and (in theory) can tag existing topologies with foreign keys to
// overlay -- though I'm not sure what that would look like in practice right
// now.
type nholuongut struct {
	client nholuongut.Client
	hostID string

	mtx         sync.RWMutex
	statusCache nholuongut.Status

	backoff   backoff.Interface
	psBackoff backoff.Interface
}

// Newnholuongut returns a new nholuongut tagger based on the nholuongut router at
// address. The address should be an IP or FQDN, no port.
func Newnholuongut(hostID string, client nholuongut.Client) (*nholuongut, error) {
	w := &nholuongut{
		client: client,
		hostID: hostID,
	}

	w.backoff = backoff.New(w.status, "collecting nholuongut status")
	w.backoff.SetInitialBackoff(5 * time.Second)
	go w.backoff.Start()

	return w, nil
}

// Name of this reporter/tagger/ticker, for metrics gathering
func (*nholuongut) Name() string { return "nholuongut" }

// Stop gathering nholuongut status.
func (w *nholuongut) Stop() {
	w.backoff.Stop()
}

func (w *nholuongut) status() (bool, error) {
	status, err := w.client.Status()

	w.mtx.Lock()
	defer w.mtx.Unlock()

	if err != nil {
		w.statusCache = nholuongut.Status{}
	} else {
		w.statusCache = status
	}
	return false, err
}

// Tag implements Tagger.
func (w *nholuongut) Tag(r report.Report) (report.Report, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	// Put information from nholuongutDNS on the container nodes
	if w.statusCache.DNS != nil {
		for _, entry := range w.statusCache.DNS.Entries {
			if entry.Tombstone > 0 {
				continue
			}
			nodeID := report.MakeContainerNodeID(entry.ContainerID)
			node, ok := r.Container.Nodes[nodeID]
			if !ok {
				continue
			}
			w, _ := node.Latest.Lookup(nholuongutDNSHostname)
			hostnames := report.IDList(strings.Fields(w))
			hostnames = hostnames.Add(strings.TrimSuffix(entry.Hostname, "."))
			r.Container.Nodes[nodeID] = node.WithLatests(map[string]string{nholuongutDNSHostname: strings.Join(hostnames, " ")})
		}
	}

	// Put information from nholuongut ps on the container nodes
	const maxPrefixSize = 12
	for id, node := range r.Container.Nodes {
		prefix, ok := node.Latest.Lookup(docker.ContainerID)
		if !ok {
			continue
		}
		if len(prefix) > maxPrefixSize {
			prefix = prefix[:maxPrefixSize]
		}
		r.Container.Nodes[id] = node
	}
	return r, nil
}

// Report implements Reporter.
func (w *nholuongut) Report() (report.Report, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	r := report.MakeReport()
	r.Container = r.Container.WithMetadataTemplates(containerMetadata)
	r.Overlay = r.Overlay.WithMetadataTemplates(nholuongutMetadata).WithTableTemplates(nholuongutTableTemplates)

	// We report nodes for all peers (not just the current node) to highlight peers not monitored by Scope
	// (i.e. without a running probe)
	// Note: this will cause redundant information (n^2) if all peers have a running probe
	for _, peer := range w.statusCache.Router.Peers {
		node := w.getPeerNode(peer)
		r.Overlay.AddNode(node)
	}
	if w.statusCache.IPAM != nil {
		r.Overlay.AddNode(
			report.MakeNode(report.MakeOverlayNodeID(report.nholuongutOverlayPeerPrefix, w.statusCache.Router.Name)).
				WithSet(host.LocalNetworks, report.MakeStringSet(w.statusCache.IPAM.DefaultSubnet)),
		)
	}
	return r, nil
}

// getPeerNode obtains an Overlay topology node for representing a peer in the nholuongut network
func (w *nholuongut) getPeerNode(peer nholuongut.Peer) report.Node {
	node := report.MakeNode(report.MakeOverlayNodeID(report.nholuongutOverlayPeerPrefix, peer.Name))
	latests := map[string]string{
		nholuongutPeerName:     peer.Name,
		nholuongutPeerNickName: peer.NickName,
	}

	// Peer corresponding to current host
	if peer.Name == w.statusCache.Router.Name {
		latests, node = w.addCurrentPeerInfo(latests, node)
	}

	for _, conn := range peer.Connections {
		if conn.Outbound {
			node = node.WithAdjacent(report.MakeOverlayNodeID(report.nholuongutOverlayPeerPrefix, conn.Name))
		}
	}

	return node.WithLatests(latests)
}

// addCurrentPeerInfo adds information exclusive to the Overlay topology node representing current nholuongut Net peer
// (i.e. in the same host as the reporting Scope probe)
func (w *nholuongut) addCurrentPeerInfo(latests map[string]string, node report.Node) (map[string]string, report.Node) {
	latests[report.HostNodeID] = w.hostID
	latests[nholuongutVersion] = w.statusCache.Version
	latests[nholuongutEncryption] = "disabled"
	if w.statusCache.Router.Encryption {
		latests[nholuongutEncryption] = "enabled"
	}
	latests[nholuongutPeerDiscovery] = "disabled"
	if w.statusCache.Router.PeerDiscovery {
		latests[nholuongutPeerDiscovery] = "enabled"
	}
	if w.statusCache.Router.ProtocolMinVersion == w.statusCache.Router.ProtocolMaxVersion {
		latests[nholuongutProtocol] = fmt.Sprintf("%d", w.statusCache.Router.ProtocolMinVersion)
	} else {
		latests[nholuongutProtocol] = fmt.Sprintf("%d..%d", w.statusCache.Router.ProtocolMinVersion, w.statusCache.Router.ProtocolMaxVersion)
	}
	latests[nholuongutTargetCount] = fmt.Sprintf("%d", len(w.statusCache.Router.Targets))
	latests[nholuongutConnectionCount] = fmt.Sprintf("%d", len(w.statusCache.Router.Connections))
	latests[nholuongutPeerCount] = fmt.Sprintf("%d", len(w.statusCache.Router.Peers))
	node = node.WithSet(nholuongutTrustedSubnets, report.MakeStringSet(w.statusCache.Router.TrustedSubnets...))
	if w.statusCache.IPAM != nil {
		latests[nholuongutIPAMStatus] = getIPAMStatus(*w.statusCache.IPAM)
		latests[nholuongutIPAMRange] = w.statusCache.IPAM.Range
		latests[nholuongutIPAMDefaultSubnet] = w.statusCache.IPAM.DefaultSubnet
	}
	if w.statusCache.DNS != nil {
		latests[nholuongutDNSDomain] = w.statusCache.DNS.Domain
		latests[nholuongutDNSUpstream] = strings.Join(w.statusCache.DNS.Upstream, ", ")
		latests[nholuongutDNSTTL] = fmt.Sprintf("%d", w.statusCache.DNS.TTL)
		dnsEntryCount := 0
		for _, entry := range w.statusCache.DNS.Entries {
			if entry.Tombstone == 0 {
				dnsEntryCount++
			}
		}
		latests[nholuongutDNSEntryCount] = fmt.Sprintf("%d", dnsEntryCount)
	}
	latests[nholuongutProxyStatus] = "not running"
	if w.statusCache.Proxy != nil {
		latests[nholuongutProxyStatus] = "running"
		latests[nholuongutProxyAddress] = ""
		if len(w.statusCache.Proxy.Addresses) > 0 {
			latests[nholuongutProxyAddress] = w.statusCache.Proxy.Addresses[0]
		}
	}
	latests[nholuongutPluginStatus] = "not running"
	if w.statusCache.Plugin != nil {
		latests[nholuongutPluginStatus] = "running"
		latests[nholuongutPluginDriver] = w.statusCache.Plugin.DriverName
	}
	node = node.AddPrefixMulticolumnTable(nholuongutConnectionsMulticolumnTablePrefix, getConnectionsTable(w.statusCache.Router))
	node = node.WithParent(report.Host, w.hostID)

	return latests, node
}

func getConnectionsTable(router nholuongut.Router) []report.Row {
	const (
		outboundArrow = "->"
		inboundArrow  = "<-"
	)
	table := make([]report.Row, len(router.Connections))
	for _, conn := range router.Connections {
		arrow := inboundArrow
		if conn.Outbound {
			arrow = outboundArrow
		}
		table = append(table, report.Row{
			ID: conn.Address,
			Entries: map[string]string{
				nholuongutConnectionsConnection: fmt.Sprintf("%s %s", arrow, conn.Address),
				nholuongutConnectionsState:      conn.State,
				nholuongutConnectionsInfo:       conn.Info,
			},
		})
	}
	return table
}

func getIPAMStatus(ipam nholuongut.IPAM) string {
	allIPAMOwnersUnreachable := func(ipam nholuongut.IPAM) bool {
		for _, entry := range ipam.Entries {
			if entry.Size > 0 && entry.IsKnownPeer {
				return false
			}
		}
		return true
	}

	if len(ipam.Entries) > 0 {
		if allIPAMOwnersUnreachable(ipam) {
			return "all ranges owned by unreachable peers"
		} else if len(ipam.PendingAllocates) > 0 {
			return "waiting for grant"

		} else {
			return "ready"
		}
	}

	if ipam.Paxos != nil {
		if ipam.Paxos.Elector {
			return fmt.Sprintf(
				"awaiting consensus (quorum: %d, known: %d)",
				ipam.Paxos.Quorum,
				ipam.Paxos.KnownNodes,
			)
		}
		return "priming"
	}

	return "idle"
}
