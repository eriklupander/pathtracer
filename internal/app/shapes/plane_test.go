package shapes

import (
	"github.com/eriklupander/pt/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPlane(t *testing.T) {
	pl := NewPlane()
	n1 := pl.NormalAtLocal(geom.NewPoint(0, 0, 0), nil)
	n2 := pl.NormalAtLocal(geom.NewPoint(10, 0, -10), nil)
	n3 := pl.NormalAtLocal(geom.NewPoint(-6, 0, 150), nil)
	assert.True(t, geom.TupleEquals(n1, geom.NewVector(0, 1, 0)))
	assert.True(t, geom.TupleEquals(n2, geom.NewVector(0, 1, 0)))
	assert.True(t, geom.TupleEquals(n3, geom.NewVector(0, 1, 0)))

}

func TestPlane_IntersectLocalParallellMisses(t *testing.T) {
	pl := NewPlane()
	r := geom.NewRay(geom.NewPoint(0, 10, 0), geom.NewVector(0, 0, 1))
	xs := Intersections{}
	pl.IntersectLocal(r, &xs)
	assert.Len(t, xs, 0)
}
func TestPlane_IntersectLocalCoplanarMisses(t *testing.T) {
	pl := NewPlane()
	r := geom.NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
	xs := Intersections{}
	pl.IntersectLocal(r, &xs)
	assert.Len(t, xs, 0)
}
func TestPlane_IntersectLocalFromAbove(t *testing.T) {
	pl := NewPlane()
	r := geom.NewRay(geom.NewPoint(0, 1, 0), geom.NewVector(0, -1, 0))
	xs := Intersections{}
	pl.IntersectLocal(r, &xs)
	assert.Len(t, xs, 1)
	assert.Equal(t, 1.0, xs[0].T)
	assert.Equal(t, pl.ID(), xs[0].S.ID())
}
func TestPlane_IntersectLocalFromBelow(t *testing.T) {
	pl := NewPlane()
	r := geom.NewRay(geom.NewPoint(0, -1, 0), geom.NewVector(0, 1, 0))
	xs := Intersections{}
	pl.IntersectLocal(r, &xs)
	assert.Len(t, xs, 1)
	assert.Equal(t, 1.0, xs[0].T)
	assert.Equal(t, pl.ID(), xs[0].S.ID())
}
