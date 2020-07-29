package shapes

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/material"
	"math"
	"math/rand"
)

func NewCube() *Cube {
	m1 := geom.New4x4()  //NewMat4x4(make([]float64, 16))
	inv := geom.New4x4() //NewMat4x4(make([]float64, 16))
	invTranspose := geom.New4x4()
	savedXs := make([]Intersection, 2)
	for i := 0; i < 2; i++ {
		savedXs[i] = NewIntersection(0.0, nil)
	}

	return &Cube{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        m1,
			Inverse:          inv,
			InverseTranspose: invTranspose,
			Material:         material.NewDefaultMaterial(),
		},
		savedXs:    savedXs,
		CastShadow: true,
	}
}

type Cube struct {
	Basic
	Label    string
	parent   Shape
	savedRay geom.Ray
	savedXs  []Intersection

	CastShadow bool
}

func (c *Cube) Name() string {
	return c.Label
}

func (c *Cube) CastsShadow() bool {
	return c.CastShadow
}

func (c *Cube) ID() int64 {
	return c.Id
}

func (c *Cube) GetTransform() geom.Mat4x4 {
	return c.Transform
}
func (c *Cube) GetInverse() geom.Mat4x4 {
	return c.Inverse
}
func (c *Cube) GetInverseTranspose() geom.Mat4x4 {
	return c.InverseTranspose
}

func (c *Cube) SetTransform(transform geom.Mat4x4) {
	c.Transform = geom.Multiply(c.Transform, transform)
	c.Inverse = geom.Inverse(c.Transform)
	c.InverseTranspose = geom.Transpose(c.Inverse)
}

func (c *Cube) GetMaterial() material.Material {
	return c.Material
}

func (c *Cube) SetMaterial(material material.Material) {
	c.Material = material
}

func (c *Cube) IntersectLocal(ray geom.Ray, xs *Intersections) {
	// There is supposed  to be a way to optimize this for fewer checks by looking at early values.
	xtmin, xtmax := checkAxis(ray.Origin.Get(0), ray.Direction.Get(0))
	ytmin, ytmax := checkAxis(ray.Origin.Get(1), ray.Direction.Get(1))
	ztmin, ztmax := checkAxis(ray.Origin.Get(2), ray.Direction.Get(2))

	// Om det största av min-värdena är större än det minsta max-värdet.
	tmin := max(xtmin, ytmin, ztmin)
	tmax := min(xtmax, ytmax, ztmax)
	if tmin > tmax {
		return
	}

	// use allocated slice and structs
	c.savedXs[0].T = tmin
	c.savedXs[0].S = c
	c.savedXs[1].T = tmax
	c.savedXs[1].S = c

	*xs = append(*xs, c.savedXs...)
}

// NormalAtLocal uses the fact that given a unit cube, the point of the surface axis X,Y or Z is always either
// 1.0 for positive XYZ and -1.0 for negative XYZ. I.e - if the point is 0.4, 1, -0.5, we know that the
// point is on the top Y surface and we can return a 0,1,0 normal
func (c *Cube) NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4 {
	maxc := max(math.Abs(point.Get(0)), math.Abs(point.Get(1)), math.Abs(point.Get(2)))
	if maxc == math.Abs(point.Get(0)) {
		return geom.NewVector(point.Get(0), 0, 0)
	} else if maxc == math.Abs(point.Get(1)) {
		return geom.NewVector(0, point.Get(1), 0)
	}
	return geom.NewVector(0, 0, point.Get(2))
}

func (c *Cube) GetLocalRay() geom.Ray {
	return c.savedRay
}
func (c *Cube) GetParent() Shape {
	return c.parent
}
func (c *Cube) SetParent(shape Shape) {
	c.parent = shape
}
func (c *Cube) Init() {
	c.savedXs = make([]Intersection, 2)
}
func checkAxis(origin float64, direction float64) (min float64, max float64) {
	tminNumerator := -1 - origin
	tmaxNumerator := 1 - origin
	var tmin, tmax float64
	if math.Abs(direction) >= geom.Epsilon {
		tmin = tminNumerator / direction
		tmax = tmaxNumerator / direction
	} else {
		tmin = tminNumerator * math.Inf(1)
		tmax = tmaxNumerator * math.Inf(1)
	}
	if tmin > tmax {
		// swap
		temp := tmin
		tmin = tmax
		tmax = temp
	}
	return tmin, tmax
}

func max(values ...float64) float64 {
	c := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > c {
			c = values[i]
		}
	}
	return c
}

func min(values ...float64) float64 {
	c := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < c {
			c = values[i]
		}
	}
	return c
}
