package tracer

import (
	"github.com/eriklupander/pathtracer/cmd"
	"github.com/eriklupander/pathtracer/internal/app/canvas"
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/scenes"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestPathTracer_Render(t *testing.T) {
	cmd.FromConfig()
	cmd.Cfg.Width = 320
	cmd.Cfg.Height = 240
	jobs := make(chan *job)
	wg := sync.WaitGroup{}
	testee := NewCtx(1, scenes.OCLScene()(), canvas.NewCanvas(320, 240), jobs, &wg)
	var cameraRay = geom.NewEmptyRay()
	testee.rayForPixelPathTracer(210, 173, &cameraRay)

	color := testee.trace(cameraRay)
	assert.NotNil(t, color)
}
