package canvas

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCanvas(t *testing.T) {
	canvas := NewCanvas(10, 20)
	assert.Equal(t, 10, canvas.W)
	assert.Equal(t, 20, canvas.H)
	for _, px := range canvas.Pixels {
		assert.True(t, px.Get(0) == 0.0)
		assert.True(t, px.Get(1) == 0.0)
		assert.True(t, px.Get(2) == 0.0)
	}
}

func TestCanvas_WritePixel(t *testing.T) {
	canvas := NewCanvas(10, 20)
	canvas.WritePixel(2, 3, geom.NewColor(1, 0, 0))
	px := canvas.ColorAt(2, 3)
	assert.True(t, px.Get(0) == 1.0)
	assert.True(t, px.Get(1) == 0.0)
	assert.True(t, px.Get(2) == 0.0)
}
