package tracer

import (
	"github.com/eriklupander/pt/cmd"
	"github.com/eriklupander/pt/internal/app/canvas"
	"github.com/eriklupander/pt/internal/app/geom"
	"github.com/eriklupander/pt/internal/app/scenes"
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
	testee := NewCtx(1, scenes.ReferenceScene()(), canvas.NewCanvas(320, 240), jobs, &wg)
	var cameraRay = geom.NewEmptyRay()
	testee.rayForPixelPathTracer(120, 169, &cameraRay)

	color := testee.trace(cameraRay)
	assert.NotNil(t, color)
}
