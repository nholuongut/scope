package render_test

import (
	"context"
	"testing"

	"github.com/nholuongutworks/common/test"
	"github.com/nholuongut/scope/render"
	"github.com/nholuongut/scope/render/expected"
	"github.com/nholuongut/scope/test/fixture"
	"github.com/nholuongut/scope/test/reflect"
	"github.com/nholuongut/scope/test/utils"
)

func TestEndpointRenderer(t *testing.T) {
	have := utils.Prune(render.EndpointRenderer.Render(context.Background(), fixture.Report).Nodes)
	want := utils.Prune(expected.RenderedEndpoints)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestProcessRenderer(t *testing.T) {
	have := utils.Prune(render.ProcessRenderer.Render(context.Background(), fixture.Report).Nodes)
	want := utils.Prune(expected.RenderedProcesses)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestProcessNameRenderer(t *testing.T) {
	have := utils.Prune(render.ProcessNameRenderer.Render(context.Background(), fixture.Report).Nodes)
	want := utils.Prune(expected.RenderedProcessNames)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}
