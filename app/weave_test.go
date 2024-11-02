package app_test

import (
	"net"
	"sync"
	"testing"
	"time"

	fsouza "github.com/fsouza/go-dockerclient"

	"github.com/nholuongut/scope/app"
	"github.com/nholuongut/scope/test"
)

type mockDockerClient struct{}

func (mockDockerClient) ListContainers(fsouza.ListContainersOptions) ([]fsouza.APIContainers, error) {
	return []fsouza.APIContainers{
		{
			Names: []string{"/" + containerName},
			ID:    containerID,
		},
		{
			Names: []string{"/notme"},
			ID:    "1234abcd",
		},
	}, nil
}

type entry struct {
	containerid string
	ip          net.IP
}

type mocknholuongutClient struct {
	sync.Mutex
	published map[string]entry
}

func (m *mocknholuongutClient) AddDNSEntry(hostname, containerid string, ip net.IP) error {
	m.Lock()
	defer m.Unlock()
	m.published[hostname] = entry{containerid, ip}
	return nil
}

func (m *mocknholuongutClient) Expose() error {
	return nil
}

const (
	hostname      = "foo.nholuongut"
	containerName = "bar"
	containerID   = "a1b2c3d4"
)

var (
	ip = net.ParseIP("1.2.3.4")
)

func Testnholuongut(t *testing.T) {
	nholuongutClient := &mocknholuongutClient{
		published: map[string]entry{},
	}
	dockerClient := mockDockerClient{}
	interfaces := func() ([]app.Interface, error) {
		return []app.Interface{
			{
				Name: "eth0",
				Addrs: []net.Addr{
					&net.IPAddr{
						IP: ip,
					},
				},
			},
			{
				Name: "docker0",
				Addrs: []net.Addr{
					&net.IPAddr{
						IP: net.ParseIP("4.3.2.1"),
					},
				},
			},
		}, nil
	}
	publisher := app.NewnholuongutPublisher(
		nholuongutClient, dockerClient, interfaces,
		hostname, containerName)
	defer publisher.Stop()

	want := map[string]entry{
		hostname: {containerID, ip},
	}
	test.Poll(t, 100*time.Millisecond, want, func() interface{} {
		nholuongutClient.Lock()
		defer nholuongutClient.Unlock()
		result := map[string]entry{}
		for k, v := range nholuongutClient.published {
			result[k] = v
		}
		return result
	})
}
