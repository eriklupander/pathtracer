package shapes

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/material"
	"math"
	"math/rand"
)

func NewCylinder() *Cylinder {
	m1 := geom.New4x4()
	inv := geom.New4x4()
	invTranspose := geom.New4x4()

	return &Cylinder{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        m1,
			Inverse:          inv,
			InverseTranspose: invTranspose,
			Material:         material.NewDefaultMaterial(),
		},
		MinY: math.Inf(-1),
		MaxY: math.Inf(1),
		//savedXs:          savedXs,
		CastShadow: true,
	}
}

func NewCylinderMM(min, max float64) *Cylinder {
	c := NewCylinder()
	c.MinY = min
	c.MaxY = max
	return c
}

func NewCylinderMMC(min, max float64, closed bool) *Cylinder {
	c := NewCylinder()
	c.MinY = min
	c.MaxY = max
	c.closed = closed
	return c
}

type Cylinder struct {
	Basic
	parent     Shape
	savedRay   geom.Ray
	MinY       float64
	MaxY       float64
	closed     bool
	CastShadow bool

	//savedXs []Intersection
}

func (c *Cylinder) CastsShadow() bool {
	return c.CastShadow
}

func (c *Cylinder) ID() int64 {
	return c.Id
}

func (c *Cylinder) GetTransform() geom.Mat4x4 {
	return c.Transform
}

func (c *Cylinder) GetInverse() geom.Mat4x4 {
	return c.Inverse
}
func (c *Cylinder) GetInverseTranspose() geom.Mat4x4 {
	return c.InverseTranspose
}

func (c *Cylinder) SetTransform(transform geom.Mat4x4) {
	c.Transform = geom.Multiply(c.Transform, transform)
	c.Inverse = geom.Inverse(c.Transform)
	c.InverseTranspose = geom.Transpose(c.Inverse)
}

func (c *Cylinder) GetMaterial() material.Material {
	return c.Material
}

func (c *Cylinder) SetMaterial(material material.Material) {
	c.Material = material
}

func (c *Cylinder) IntersectLocal(ray geom.Ray, xs *Intersections) {
	rdx2 := ray.Direction.Get(0) * ray.Direction.Get(0)
	rdz2 := ray.Direction.Get(2) * ray.Direction.Get(2)

	a := rdx2 + rdz2
	if math.Abs(a) < geom.Epsilon {
		c.intercectCaps(ray, xs)
		return
	}

	b := 2*ray.Origin.Get(0)*ray.Direction.Get(0) +
		2*ray.Origin.Get(2)*ray.Direction.Get(2)

	rox2 := ray.Origin.Get(0) * ray.Origin.Get(0)
	roz2 := ray.Origin.Get(2) * ray.Origin.Get(2)

	c1 := rox2 + roz2 - 1

	disc := b*b - 4*a*c1

	// ray does not intersect the cylinder
	if disc < 0 {
		return
	}

	t0 := (-b - math.Sqrt(disc)) / (2 * a)
	t1 := (-b + math.Sqrt(disc)) / (2 * a)

	y0 := ray.Origin.Get(1) + t0*ray.Direction.Get(1)
	if y0 > c.MinY && y0 < c.MaxY {
		*xs = append(*xs, NewIntersection(t0, c))
	}

	y1 := ray.Origin.Get(1) + t1*ray.Direction.Get(1)
	if y1 > c.MinY && y1 < c.MaxY {
		*xs = append(*xs, NewIntersection(t1, c))
	}

	c.intercectCaps(ray, xs)
}

func (c *Cylinder) NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4 {

	// compute the square of the distance from the y axis
	dist := math.Pow(point.Get(0), 2) + math.Pow(point.Get(2), 2)
	if dist < 1 && point.Get(1) >= c.MaxY-geom.Epsilon {
		return geom.NewVector(0, 1, 0)
	} else if dist < 1 && point.Get(1) <= c.MinY+geom.Epsilon {
		return geom.NewVector(0, -1, 0)
	} else {
		return geom.NewVector(point.Get(0), 0, point.Get(2))
	}
}

func (c *Cylinder) GetLocalRay() geom.Ray {
	return c.savedRay
}
func (c *Cylinder) GetParent() Shape {
	return c.parent
}
func (c *Cylinder) SetParent(shape Shape) {
	c.parent = shape
}

func (c *Cylinder) Init() {}

func checkCap(ray geom.Ray, t float64) bool {
	x := ray.Origin.Get(0) + t*ray.Direction.Get(0)
	z := ray.Origin.Get(2) + t*ray.Direction.Get(2)
	return math.Pow(x, 2)+math.Pow(z, 2) <= 1.0
}

func (c *Cylinder) intercectCaps(ray geom.Ray, xs *Intersections) {
	if !c.closed || math.Abs(ray.Direction.Get(1)) < geom.Epsilon {
		return
	}

	// check for an intersection with the lower end cap by intersecting
	// the ray with the plane at y=cyl.minimum
	t := (c.MinY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if checkCap(ray, t) {
		*xs = append(*xs, NewIntersection(t, c))
	}

	// check for an intersection with the upper end cap by intersecting
	// the ray with the plane at y=cyl.maximum
	t = (c.MaxY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if checkCap(ray, t) {
		*xs = append(*xs, NewIntersection(t, c))
	}
}
