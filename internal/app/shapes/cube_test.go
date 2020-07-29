package shapes

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCubeRayHits(t *testing.T) {
	c := NewCube()

	tc := []cubetest{
		{name: "+x", p: geom.NewPoint(5, 0.5, 0), v: geom.NewVector(-1, 0, 0), t1: 4.0, t2: 6.0},
		{name: "-x", p: geom.NewPoint(-5, 0.5, 0), v: geom.NewVector(1, 0, 0), t1: 4.0, t2: 6.0},
		{name: "+y", p: geom.NewPoint(0.5, 5, 0), v: geom.NewVector(0, -1, 0), t1: 4.0, t2: 6.0},
		{name: "-y", p: geom.NewPoint(0.5, -5, 0), v: geom.NewVector(0, 1, 0), t1: 4.0, t2: 6.0},
		{name: "+z", p: geom.NewPoint(0.5, 0, 5), v: geom.NewVector(0, 0, -1), t1: 4.0, t2: 6.0},
		{name: "-z", p: geom.NewPoint(0.5, 0, -5), v: geom.NewVector(0, 0, 1), t1: 4.0, t2: 6.0},
	}
	for _, test := range tc {
		xs := Intersections{}
		r := geom.NewRay(test.p, test.v)
		c.IntersectLocal(r, &xs)
		assert.Equal(t, test.t1, xs[0].T)
		assert.Equal(t, test.t2, xs[1].T)
	}
}

func TestCubeRayMisses(t *testing.T) {
	c := NewCube()

	tc := []cubetest{
		{p: geom.NewPoint(-2, 0, 0), v: geom.NewVector(0.2673, 0.5345, 0.8018)},
		{p: geom.NewPoint(0, -2, 0), v: geom.NewVector(0.8018, 0.2673, 0.5345)},
		{p: geom.NewPoint(0, 0, -2), v: geom.NewVector(0.5345, 0.8018, 0.2673)},
		{p: geom.NewPoint(2, 0, 2), v: geom.NewVector(0, 0, -1)},
		{p: geom.NewPoint(0, 2, 2), v: geom.NewVector(0, -1, 0)},
		{p: geom.NewPoint(2, 2, 0), v: geom.NewVector(-1, 0, 0)},
	}
	for _, test := range tc {
		xs := Intersections{}
		r := geom.NewRay(test.p, test.v)
		c.IntersectLocal(r, &xs)
		assert.Len(t, xs, 0)
	}
}

func TestCubeNormal(t *testing.T) {
	c := NewCube()

	tc := []cubetest{
		{p: geom.NewPoint(1, 0.5, -0.8), v: geom.NewVector(1, 0, 0)},
		{p: geom.NewPoint(-1, -0.2, 0.9), v: geom.NewVector(-1, 0, 0)},
		{p: geom.NewPoint(-0.4, 1, -0.1), v: geom.NewVector(0, 1, 0)},
		{p: geom.NewPoint(0.3, -1, -0.7), v: geom.NewVector(0, -1, 0)},
		{p: geom.NewPoint(-0.6, 0.3, 1), v: geom.NewVector(0, 0, 1)},
		{p: geom.NewPoint(0.4, 0.4, -1), v: geom.NewVector(0, 0, -1)},
		{p: geom.NewPoint(1, 1, 1), v: geom.NewVector(1, 0, 0)},
		{p: geom.NewPoint(-1, -1, -1), v: geom.NewVector(-1, 0, 0)},
	}
	for _, test := range tc {
		n := c.NormalAtLocal(test.p, nil)
		assert.Equal(t, test.v, n)
	}
}

type cubetest struct {
	name string
	p    geom.Tuple4
	v    geom.Tuple4
	t1   float64
	t2   float64
}
