package tracer

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/shapes"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNormalOnSphereAtX(t *testing.T) {
	s := shapes.NewSphere()
	normalVector := NormalAt(s, geom.NewPoint(1, 0, 0), nil)
	assert.True(t, geom.TupleEquals(normalVector, geom.NewVector(1, 0, 0)))
}
func TestNormalOnSphereAtY(t *testing.T) {
	s := shapes.NewSphere()
	normalVector := NormalAt(s, geom.NewPoint(0, 1, 0), nil)
	assert.True(t, geom.TupleEquals(normalVector, geom.NewVector(0, 1, 0)))
}
func TestNormalAtPointOnSphereAtZ(t *testing.T) {
	s := shapes.NewSphere()
	normalVector := NormalAt(s, geom.NewPoint(0, 0, 1), nil)
	assert.True(t, geom.TupleEquals(normalVector, geom.NewVector(0, 0, 1)))
}
func TestNormalOnSphereAtNonAxial(t *testing.T) {
	s := shapes.NewSphere()
	nonAxial := math.Sqrt(3.0) / 3.0
	normalVector := NormalAt(s, geom.NewPoint(nonAxial, nonAxial, nonAxial), nil)
	assert.InEpsilon(t, nonAxial, normalVector.Get(0), geom.Epsilon)
	assert.InEpsilon(t, nonAxial, normalVector.Get(1), geom.Epsilon)
	assert.InEpsilon(t, nonAxial, normalVector.Get(2), geom.Epsilon)
}
func TestNormalIsNormalized(t *testing.T) {
	s := shapes.NewSphere()
	nonAxial := math.Sqrt(3.0) / 3.0
	normalVector := NormalAt(s, geom.NewPoint(nonAxial, nonAxial, nonAxial), nil)
	normalizedNormalVector := geom.Normalize(normalVector)
	assert.True(t, geom.TupleEquals(normalVector, normalizedNormalVector))
}
func TestComputeNormalOnTranslatedSphere(t *testing.T) {

	s := shapes.NewSphere()
	s.SetTransform(geom.Translate(0, 1, 0))
	normalVector := NormalAt(s, geom.NewPoint(0, 1.70711, -0.70711), nil)
	assert.Equal(t, 0.0, normalVector.Get(0))
	assert.InEpsilon(t, 0.70711, normalVector.Get(1), geom.Epsilon)
	assert.InEpsilon(t, -0.70711, normalVector.Get(2), geom.Epsilon)
}

func TestComputeNormalOnTransformedSphere(t *testing.T) {

	s := shapes.NewSphere()
	m1 := geom.Multiply(geom.Scale(1, 0.5, 1), geom.RotateZ(math.Pi/5.0))
	s.SetTransform(m1)
	normalVector := NormalAt(s, geom.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2), nil)
	assert.Equal(t, 0.0, normalVector.Get(0))
	assert.InEpsilon(t, 0.97014, normalVector.Get(1), geom.Epsilon)
	assert.InEpsilon(t, -0.24254, normalVector.Get(2), geom.Epsilon)
}

// Reflecting a vector approaching at 45°
func TestReflectRay(t *testing.T) {
	v := geom.NewVector(1, -1, 0)
	normal := geom.NewVector(0, 1, 0) // straight up
	reflectV := Reflect(v, normal)
	assert.True(t, geom.TupleEquals(geom.NewVector(1, 1, 0), reflectV))
}
func TestReflectRaySlanted(t *testing.T) {
	v := geom.NewVector(0, -1, 0) // Pointing straight down
	fortyFive := math.Sqrt(2) / 2.0
	normal := geom.NewVector(fortyFive, fortyFive, 0) // straight up
	reflectV := Reflect(v, normal)
	assert.True(t, geom.TupleEquals(geom.NewVector(1, 0, 0), reflectV))
}

//
//func TestPrecomputingReflectionVector(t *testing.T) {
//	pl := shapes.NewPlane()
//	ray := geom.NewRay(geom.NewPoint(0, 1, -1), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
//	intersections := NewIntersection(math.Sqrt(2), pl)
//	comps := PrepareComputationForIntersection(intersections, ray)
//	assert.Equal(t, 0.0, comps.ReflectVec.Get(0))
//	assert.InEpsilon(t, math.Sqrt(2)/2, comps.ReflectVec.Get(1), geom.Epsilon)
//	assert.InEpsilon(t, math.Sqrt(2)/2, comps.ReflectVec.Get(2), geom.Epsilon)
//}
