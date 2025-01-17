package render_test

import (
	"context"
	"testing"

	"github.com/nholuongutworks/common/test"
	"github.com/nholuongut/scope/probe/kubernetes"
	"github.com/nholuongut/scope/render"
	"github.com/nholuongut/scope/render/expected"
	"github.com/nholuongut/scope/test/fixture"
	"github.com/nholuongut/scope/test/reflect"
	"github.com/nholuongut/scope/test/utils"
)

func TestPodRenderer(t *testing.T) {
	have := utils.Prune(render.PodRenderer.Render(context.Background(), fixture.Report).Nodes)
	want := utils.Prune(expected.RenderedPods)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

var filterNonKubeSystem = render.Transformers([]render.Transformer{
	render.Complement(render.IsNamespace("kube-system")),
	render.FilterUnconnectedPseudo,
})

func TestPodFilterRenderer(t *testing.T) {
	// tag on containers or pod namespace in the topology and ensure
	// it is filtered out correctly.
	input := fixture.Report.Copy()
	input.Pod.Nodes[fixture.ClientPodNodeID] = input.Pod.Nodes[fixture.ClientPodNodeID].WithLatests(map[string]string{
		kubernetes.Namespace: "kube-system",
	})

	have := utils.Prune(render.Render(context.Background(), input, render.PodRenderer, filterNonKubeSystem).Nodes)
	want := utils.Prune(expected.RenderedPods.Copy())
	delete(want, fixture.ClientPodNodeID)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestPodServiceRenderer(t *testing.T) {
	have := utils.Prune(render.PodServiceRenderer.Render(context.Background(), fixture.Report).Nodes)
	want := utils.Prune(expected.RenderedPodServices)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestPodServiceFilterRenderer(t *testing.T) {
	// tag on containers or pod namespace in the topology and ensure
	// it is filtered out correctly.
	input := fixture.Report.Copy()
	input.Service.Nodes[fixture.ServiceNodeID] = input.Service.Nodes[fixture.ServiceNodeID].WithLatests(map[string]string{
		kubernetes.Namespace: "kube-system",
	})

	have := utils.Prune(render.Render(context.Background(), input, render.PodServiceRenderer, filterNonKubeSystem).Nodes)
	want := utils.Prune(expected.RenderedPodServices.Copy())
	delete(want, fixture.ServiceNodeID)
	delete(want, render.IncomingInternetID)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}
