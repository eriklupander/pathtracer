package shapes

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/material"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewSphere() *Sphere {

	return &Sphere{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        geom.New4x4(),
			Inverse:          geom.New4x4(),
			InverseTranspose: geom.New4x4(),
			Material: material.Material{
				Color:           geom.Tuple4{1, .5, .5},
				Emission:        geom.Tuple4{0, 0, 0},
				RefractiveIndex: 1,
			},
		},
		savedVec:    geom.NewVector(0, 0, 0),
		savedNormal: geom.NewVector(0, 0, 0),
		savedRay:    geom.NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 0)),
		xsCache:     make([]Intersection, 2),
		//xsEmpty:     make([]Intersection, 0),
		originPoint: geom.NewPoint(0, 0, 0),
		CastShadow:  true,
	}
}

type Sphere struct {
	Basic
	parent   Shape
	savedRay geom.Ray

	// cached stuff
	originPoint geom.Tuple4
	savedVec    geom.Tuple4
	xsCache     []Intersection
	xsEmpty     []Intersection

	savedNormal geom.Tuple4

	CastShadow bool
}

func (s *Sphere) ID() int64 {
	return s.Id
}

func (s *Sphere) CastsShadow() bool {
	return s.CastShadow
}

func (s *Sphere) GetParent() Shape {
	return s.parent
}

func (s *Sphere) NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4 {
	geom.SubPtr(point, s.originPoint, &s.savedNormal)
	return s.savedNormal
}

func (s *Sphere) GetLocalRay() geom.Ray {
	return s.savedRay
}

// IntersectLocal implements Sphere-ray intersection
func (s *Sphere) IntersectLocal(ray geom.Ray, xs *Intersections) {
	s.savedRay = ray
	//s.XsCache = s.XsCache[:0]
	// this is a vector from the origin of the ray to the center of the sphere at 0,0,0
	//SubPtr(r.Origin, s.originPoint, &s.savedVec)

	// Note that doing the Subtraction inlined was much faster than letting SubPtr do it.
	// Shouldn't the SubPtr be inlined by the compiler? Need to figure out what's going on here...
	for i := 0; i < 4; i++ {
		s.savedVec[i] = ray.Origin[i] - s.originPoint[i]
	}

	// This dot product is
	a := geom.Dot(ray.Direction, ray.Direction)

	// Take the dot of the direction and the vector from ray origin to sphere center times 2
	b := 2.0 * geom.Dot(ray.Direction, s.savedVec)

	// Take the dot of the two sphereToRay vectors and decrease by 1 (is that because the sphere is unit length 1?
	c := geom.Dot(s.savedVec, s.savedVec) - 1.0

	// calculate the discriminant
	discriminant := (b * b) - 4*a*c
	if discriminant < 0.0 {
		return // s.xsEmpty
	}

	// finally, find the intersection distances on our ray. Some values:
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	s.xsCache[0].T = t1
	s.xsCache[1].T = t2
	s.xsCache[0].S = s
	s.xsCache[1].S = s
	*xs = append(*xs, s.xsCache...)
}

func (s *Sphere) GetTransform() geom.Mat4x4 {
	return s.Transform
}
func (s *Sphere) GetInverse() geom.Mat4x4 {
	return s.Inverse
}
func (s *Sphere) GetInverseTranspose() geom.Mat4x4 {
	return s.InverseTranspose
}
func (s *Sphere) GetMaterial() material.Material {
	return s.Material
}

// SetTransform passes a pointer to the Sphere on which to apply the translation matrix
func (s *Sphere) SetTransform(translation geom.Mat4x4) {
	s.Transform = geom.Multiply(s.Transform, translation)
	s.Inverse = geom.Inverse(s.Transform)
	s.InverseTranspose = geom.Transpose(s.Inverse)
}

// SetMaterial passes a pointer to the Sphere on which to set the material
func (s *Sphere) SetMaterial(m material.Material) {
	s.Material = m
}

func (s *Sphere) SetParent(shape Shape) {
	s.parent = shape
}
