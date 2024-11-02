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

func TestHostRenderer(t *testing.T) {
	have := utils.Prune(render.HostRenderer.Render(context.Background(), fixture.Report).Nodes)
	want := utils.Prune(expected.RenderedHosts)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}
