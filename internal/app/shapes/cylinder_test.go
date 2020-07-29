package shapes

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCylinderRayMisses(t *testing.T) {
	c := NewCylinder()

	tc := []cyltest{
		{p: geom.NewPoint(1, 0, 0), v: geom.NewVector(0, 1, 0)},
		{p: geom.NewPoint(0, 0, 0), v: geom.NewVector(0, 1, 0)},
		{p: geom.NewPoint(0, 0, -5), v: geom.NewVector(1, 1, 1)},
	}
	for _, test := range tc {
		xs := Intersections{}
		r := geom.NewRay(test.p, test.v)
		c.IntersectLocal(r, &xs)
		assert.Len(t, xs, 0)
	}
}
func TestCylinderRayHits(t *testing.T) {
	c := NewCylinder()

	tc := []cyltest{
		{p: geom.NewPoint(1, 0, -5), v: geom.Normalize(geom.NewVector(0, 0, 1)), t1: 5, t2: 5},
		{p: geom.NewPoint(0, 0, -5), v: geom.Normalize(geom.NewVector(0, 0, 1)), t1: 4, t2: 6},
		{p: geom.NewPoint(0.5, 0, -5), v: geom.Normalize(geom.NewVector(0.1, 1, 1)), t1: 6.80798191702732, t2: 7.088723439378861},
	}
	for _, test := range tc {
		xs := Intersections{}
		r := geom.NewRay(test.p, test.v)
		c.IntersectLocal(r, &xs)
		assert.Equal(t, 2, len(xs))
		assert.Equal(t, test.t1, xs[0].T)
		assert.Equal(t, test.t2, xs[1].T)
	}
}

func TestCylinderLocalNormal(t *testing.T) {
	c := NewCylinder()

	tc := []cyltest{
		{p: geom.NewPoint(1, 0, 0), v: geom.Normalize(geom.NewVector(1, 0, 0))},
		{p: geom.NewPoint(0, 5, -1), v: geom.Normalize(geom.NewVector(0, 0, -1))},
		{p: geom.NewPoint(0, -2, 1), v: geom.Normalize(geom.NewVector(0, 0, 1))},
		{p: geom.NewPoint(-1, 1, 0), v: geom.Normalize(geom.NewVector(-1, 0, 0))},
	}
	for _, test := range tc {
		n := c.NormalAtLocal(test.p, nil)
		assert.Equal(t, test.v, n)
	}
}

func TestIntersectCappedOpenCylinder(t *testing.T) {
	c := NewCylinderMM(1, 2)

	tc := []cyltest{
		{p: geom.NewPoint(0, 1.5, 0), v: geom.NewVector(0.1, 1, 0), t1: 0},
		{p: geom.NewPoint(0, 3, -5), v: geom.NewVector(0, 0, 1), t1: 0},
		{p: geom.NewPoint(0, 0, -5), v: geom.NewVector(0, 0, 1), t1: 0},
		{p: geom.NewPoint(0, 2, -5), v: geom.NewVector(0, 0, 1), t1: 0},
		{p: geom.NewPoint(0, 1, -5), v: geom.NewVector(0, 0, 1), t1: 0},
		{p: geom.NewPoint(0, 1.5, -2), v: geom.NewVector(0, 0, 1), t1: 2},
	}
	for _, test := range tc {
		xs := Intersections{}
		c.IntersectLocal(geom.NewRay(test.p, geom.Normalize(test.v)), &xs)
		assert.Len(t, xs, int(test.t1))
	}
}

func TestIntersectCappedClosedCylinder(t *testing.T) {
	c := NewCylinderMMC(1, 2, true)

	tc := []cyltest{
		{p: geom.NewPoint(0, 3, 0), v: geom.NewVector(0, -1, 0), t1: 2},
		{p: geom.NewPoint(0, 3, -2), v: geom.NewVector(0, -1, 2), t1: 2},
		{p: geom.NewPoint(0, 4, -2), v: geom.NewVector(0, -1, 1), t1: 2},
		{p: geom.NewPoint(0, 0, -2), v: geom.NewVector(0, 1, 2), t1: 2},
		{p: geom.NewPoint(0, -1, -2), v: geom.NewVector(0, 1, 1), t1: 2},
	}
	for _, test := range tc {
		xs := Intersections{}
		c.IntersectLocal(geom.NewRay(test.p, geom.Normalize(test.v)), &xs)
		assert.Equal(t, len(xs), int(test.t1))
	}
}

func TestCylinderNormalAtCap(t *testing.T) {
	cyl := NewCylinderMMC(1, 2, true)

	tc := []cyltest{
		{p: geom.NewPoint(0, 1, 0), v: geom.NewVector(0, -1, 0), t1: 2},
		{p: geom.NewPoint(0.5, 1, 0), v: geom.NewVector(0, -1, 0), t1: 2},
		{p: geom.NewPoint(0, 1, 0.5), v: geom.NewVector(0, -1, 0), t1: 2},
		{p: geom.NewPoint(0, 2, 0), v: geom.NewVector(0, 1, 0), t1: 2},
		{p: geom.NewPoint(0.5, 2, 0), v: geom.NewVector(0, 1, 0), t1: 2},
		{p: geom.NewPoint(0, 2, 0.5), v: geom.NewVector(0, 1, 0), t1: 2},
	}

	for _, test := range tc {
		n := cyl.NormalAtLocal(test.p, nil)
		assert.Equal(t, test.v, n)
	}
}

type cyltest struct {
	p  geom.Tuple4
	v  geom.Tuple4
	t1 float64
	t2 float64
}
