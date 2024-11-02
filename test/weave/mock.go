package nholuongut

import (
	"net"

	"github.com/nholuongut/scope/common/nholuongut"
)

// Constants used for testing
const (
	MocknholuongutPeerName      = "winnebago"
	MocknholuongutPeerNickName  = "winny"
	MocknholuongutDefaultSubnet = "10.32.0.1/12"
	MockContainerID        = "83183a667c01"
	MockHostname           = "hostname.nholuongut.local"
	MockProxyAddress       = "unix:///foo/bar/nholuongut.sock"
	MockDriverName         = "nholuongut_mock"
)

// MockClient is a mock version of nholuongut.Client
type MockClient struct{}

// Status implements nholuongut.Client
func (MockClient) Status() (nholuongut.Status, error) {
	return nholuongut.Status{
		Router: nholuongut.Router{
			Name: MocknholuongutPeerName,
			Peers: []nholuongut.Peer{
				{
					Name:     MocknholuongutPeerName,
					NickName: MocknholuongutPeerNickName,
				},
			},
		},
		DNS: &nholuongut.DNS{
			Entries: []struct {
				Hostname    string
				ContainerID string
				Tombstone   int64
			}{
				{
					Hostname:    MockHostname + ".",
					ContainerID: MockContainerID,
					Tombstone:   0,
				},
			},
		},
		IPAM: &nholuongut.IPAM{
			DefaultSubnet: MocknholuongutDefaultSubnet,
			Entries: []struct {
				Size        uint32
				IsKnownPeer bool
			}{
				{Size: 0, IsKnownPeer: false},
			},
		},
		Proxy: &nholuongut.Proxy{
			Addresses: []string{MockProxyAddress},
		},
		Plugin: &nholuongut.Plugin{
			DriverName: MockDriverName,
		},
	}, nil
}

// AddDNSEntry implements nholuongut.Client
func (MockClient) AddDNSEntry(fqdn, containerid string, ip net.IP) error {
	return nil
}

// Expose implements nholuongut.Client
func (MockClient) Expose() error {
	return nil
}
