package docker_test

import (
	"io"
	"reflect"
	"testing"
	"time"

	commonTest "github.com/nholuongutworks/common/test"
	"github.com/nholuongut/scope/common/xfer"
	"github.com/nholuongut/scope/probe/controls"
	"github.com/nholuongut/scope/probe/docker"
	"github.com/nholuongut/scope/report"
	"github.com/nholuongut/scope/test"
)

func TestControls(t *testing.T) {
	mdc := newMockClient()
	setupStubs(mdc, func() {
		hr := controls.NewDefaultHandlerRegistry()
		registry, _ := docker.NewRegistry(docker.RegistryOptions{
			Interval:        10 * time.Second,
			HandlerRegistry: hr,
		})
		defer registry.Stop()

		for _, tc := range []struct{ command, result string }{
			{docker.StopContainer, "stopped"},
			{docker.StartContainer, "started"},
			{docker.RestartContainer, "restarted"},
			{docker.PauseContainer, "paused"},
			{docker.UnpauseContainer, "unpaused"},
		} {
			result := hr.HandleControlRequest(xfer.Request{
				Control: tc.command,
				NodeID:  report.MakeContainerNodeID("a1b2c3d4e5"),
			})
			if !reflect.DeepEqual(result, xfer.Response{
				Error: tc.result,
			}) {
				t.Error(result)
			}
		}
	})
}

type mockPipe struct{}

func (mockPipe) Ends() (io.ReadWriter, io.ReadWriter)                { return nil, nil }
func (mockPipe) CopyToWebsocket(io.ReadWriter, xfer.Websocket) error { return nil }
func (mockPipe) Close() error                                        { return nil }
func (mockPipe) Closed() bool                                        { return false }
func (mockPipe) OnClose(func())                                      {}

func TestPipes(t *testing.T) {
	oldNewPipe := controls.NewPipe
	defer func() { controls.NewPipe = oldNewPipe }()
	controls.NewPipe = func(_ controls.PipeClient, _ string) (string, xfer.Pipe, error) {
		return "pipeid", mockPipe{}, nil
	}

	mdc := newMockClient()
	setupStubs(mdc, func() {
		hr := controls.NewDefaultHandlerRegistry()
		registry, _ := docker.NewRegistry(docker.RegistryOptions{
			Interval:        10 * time.Second,
			HandlerRegistry: hr,
		})
		defer registry.Stop()

		test.Poll(t, 100*time.Millisecond, true, func() interface{} {
			_, ok := registry.GetContainer("ping")
			return ok
		})

		for _, want := range []struct {
			control  string
			response xfer.Response
		}{
			{
				control: docker.AttachContainer,
				response: xfer.Response{
					Pipe:   "pipeid",
					RawTTY: true,
				},
			},

			{
				control: docker.ExecContainer,
				response: xfer.Response{
					Pipe:             "pipeid",
					RawTTY:           true,
					ResizeTTYControl: docker.ResizeExecTTY,
				},
			},
		} {
			result := hr.HandleControlRequest(xfer.Request{
				Control: want.control,
				NodeID:  report.MakeContainerNodeID("ping"),
			})
			if !reflect.DeepEqual(result, want.response) {
				t.Errorf("diff %s: %s", want.control, commonTest.Diff(want, result))
			}
		}
	})
}

func TestDockerImageName(t *testing.T) {
	for _, input := range []struct{ in, name string }{
		{"foo/bar", "foo/bar"},
		{"foo/bar:baz", "foo/bar"},
		{"reg:123/foo/bar:baz", "foo/bar"},
		{"docker-registry.domain.name:5000/repo/image1:ver", "repo/image1"},
		{"foo", "foo"},
	} {
		name := docker.ImageNameWithoutTag(input.in)
		if name != input.name {
			t.Fatalf("%s: %s != %s", input.in, name, input.name)
		}
	}
}
