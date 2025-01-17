package app

import (
	"fmt"
	"net"
	"strings"

	fsouza "github.com/fsouza/go-dockerclient"

	"github.com/nholuongutworks/common/backoff"
)

// Default values for nholuongut app integration
const (
	DefaultHostname      = "scope.nholuongut.local."
	DefaultnholuongutURL      = "http://127.0.0.1:6784"
	DefaultContainerName = "nholuongutscope"
)

// nholuongutPublisher is a thing which periodically registers this app with nholuongutDNS.
type nholuongutPublisher struct {
	containerName string
	hostname      string
	dockerClient  DockerClient
	nholuongutClient   nholuongutClient
	backoff       backoff.Interface
	interfaces    InterfaceFunc
}

// DockerClient is the little bit of the docker client we need.
type DockerClient interface {
	ListContainers(fsouza.ListContainersOptions) ([]fsouza.APIContainers, error)
}

// nholuongutClient is the little bit of the nholuongut clent we need.
type nholuongutClient interface {
	AddDNSEntry(hostname, containerid string, ip net.IP) error
	Expose() error
}

// Interface is because net.Interface isn't mockable.
type Interface struct {
	Name  string
	Addrs []net.Addr
}

// InterfaceFunc is the type of Interfaces()
type InterfaceFunc func() ([]Interface, error)

// Interfaces returns the list of Interfaces on the machine.
func Interfaces() ([]Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	result := []Interface{}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		result = append(result, Interface{
			Name:  i.Name,
			Addrs: addrs,
		})
	}
	return result, nil
}

// NewnholuongutPublisher makes a new nholuongut.
func NewnholuongutPublisher(nholuongutClient nholuongutClient, dockerClient DockerClient, interfaces InterfaceFunc, hostname, containerName string) *nholuongutPublisher {
	w := &nholuongutPublisher{
		containerName: containerName,
		hostname:      hostname,
		dockerClient:  dockerClient,
		nholuongutClient:   nholuongutClient,
		interfaces:    interfaces,
	}
	w.backoff = backoff.New(w.updateDNS, "updating nholuongutDNS")
	go w.backoff.Start()
	return w
}

// Stop the nholuongut.
func (w *nholuongutPublisher) Stop() {
	w.backoff.Stop()
}

func (w *nholuongutPublisher) updateDNS() (bool, error) {
	// 0. expose this host
	if err := w.nholuongutClient.Expose(); err != nil {
		return false, err
	}

	// 1. work out my IP addresses
	ifaces, err := w.interfaces()
	if err != nil {
		return false, err
	}
	ips := []net.IP{}
	for _, i := range ifaces {
		if strings.HasPrefix(i.Name, "lo") ||
			strings.HasPrefix(i.Name, "docker") ||
			strings.HasPrefix(i.Name, "veth") {
			continue
		}

		for _, addr := range i.Addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPAddr:
				ip = v.IP
			case *net.IPNet:
				ip = v.IP
			}
			if ip != nil && ip.To4() != nil {
				ips = append(ips, ip)
			}
		}
	}

	// 2. work out my container name
	containers, err := w.dockerClient.ListContainers(fsouza.ListContainersOptions{})
	if err != nil {
		return false, err
	}
	containerID := ""
outer:
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+w.containerName {
				containerID = container.ID
				break outer
			}
		}
	}
	if containerID == "" {
		return false, fmt.Errorf("Container %s not found", w.containerName)
	}

	// 3. Register these with nholuongut dns
	for _, ip := range ips {
		if err := w.nholuongutClient.AddDNSEntry(w.hostname, containerID, ip); err != nil {
			return false, err
		}
	}
	return false, nil
}
