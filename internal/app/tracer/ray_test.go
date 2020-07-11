package tracer

import (
	"fmt"
	"github.com/eriklupander/pt/internal/app/geom"
	"github.com/eriklupander/pt/internal/app/shapes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRay(t *testing.T) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(4, 5, 6))
	assert.True(t, geom.TupleEquals(r.Origin, geom.NewPoint(1, 2, 3)))
	assert.True(t, geom.TupleEquals(r.Direction, geom.NewVector(4, 5, 6)))
}

func TestDistanceFromPoint(t *testing.T) {
	r := geom.NewRay(geom.NewPoint(2, 3, 4), geom.NewVector(1, 0, 0))
	p1 := Position(r, 0)
	assert.Equal(t, geom.NewPoint(2, 3, 4), p1)
}

//
//func TestIntersectSphereShape(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	interects := IntersectRayWithShape(s, r)
//	assert.Len(t, interects, 2)
//	assert.Equal(t, 4.0, interects[0].T)
//	assert.Equal(t, 6.0, interects[1].T)
//}
//
//func TestIntersectSphereAtTangent(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 1, -5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	interects := IntersectRayWithShape(s, r)
//	assert.Len(t, interects, 2)
//	assert.Equal(t, 5.0, interects[0].T)
//	assert.Equal(t, 5.0, interects[1].T)
//}
//
//func TestIntersectMissSphere(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 2, -5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	interects := IntersectRayWithShape(s, r)
//	assert.Len(t, interects, 0)
//}
//
//func TestIntersectSphereWhenOriginatingFromCenterOfSphere(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	interects := IntersectRayWithShape(s, r)
//	assert.Len(t, interects, 2)
//	assert.Equal(t, -1.0, interects[0].T)
//	assert.Equal(t, 1.0, interects[1].T)
//}
//
//func TestIntersectSphereBehindRay(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 0, 5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	interects := IntersectRayWithShape(s, r)
//	assert.Len(t, interects, 2)
//	assert.Equal(t, -6.0, interects[0].T)
//	assert.Equal(t, -4.0, interects[1].T)
//}

func TestHitWhenAllIntersectsHavePositiveT(t *testing.T) {
	s := shapes.NewSphere()
	i1 := shapes.NewIntersection(1.0, s)
	i2 := shapes.NewIntersection(2.0, s)
	xs := []shapes.Intersection{i1, i2}
	i, found := Hit(xs)
	assert.True(t, found)
	assert.True(t, shapes.IntersectionEqual(i, i1))
}

func TestHitWhenSomeIntersectsHaveNegativeT(t *testing.T) {
	s := shapes.NewSphere()
	i1 := shapes.NewIntersection(-1.0, s)
	i2 := shapes.NewIntersection(1.0, s)
	xs := []shapes.Intersection{i1, i2}
	i, _ := Hit(xs)
	assert.True(t, shapes.IntersectionEqual(i, i2))
}

func TestHitWhenAllIntersectsHaveNegativeT(t *testing.T) {
	s := shapes.NewSphere()
	i1 := shapes.NewIntersection(-2.0, s)
	i2 := shapes.NewIntersection(-1.0, s)
	xs := []shapes.Intersection{i1, i2}
	_, found := Hit(xs)
	assert.False(t, found)
}

// NOTE! This test has been commented out since the list of intersections
// passed to Hit() always has been sorted. This is an optimization.
//func TestHitIsLowestNonNegativeT(t *testing.T) {
//	s := shapes.NewSphere()
//	i1 := NewIntersection(5.0, s)
//	i2 := NewIntersection(7.0, s)
//	i3 := NewIntersection(-3, s)
//	i4 := NewIntersection(2.0, s)
//	intersections := []Intersection{i1, i2, i3, i4}
//	i, _ := Hit(intersections)
//	assert.True(t, IntersectionEqual(i, i4))
//}

func TestTranslateRay(t *testing.T) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(0, 1, 0))
	m1 := geom.Translate(3, 4, 5)
	r2 := TransformRay(r, m1)
	assert.True(t, geom.TupleEquals(r2.Origin, geom.NewPoint(4, 6, 8)))
	assert.True(t, geom.TupleEquals(r2.Direction, geom.NewVector(0, 1, 0)))
}
func BenchmarkTransformRay(b *testing.B) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(0, 1, 0))
	m1 := geom.Translate(3, 4, 5)
	var r2 geom.Ray
	for i := 0; i < b.N; i++ {
		r2 = TransformRay(r, m1)
	}
	fmt.Printf("%v\n", r2)
}
func BenchmarkTransformRayPtr(b *testing.B) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(0, 1, 0))
	m1 := geom.Translate(3, 4, 5)
	var r2 geom.Ray
	for i := 0; i < b.N; i++ {
		TransformRayPtr(r, m1, &r2)
	}
	fmt.Printf("%v\n", r2)
}

//
//func TestScaleRay(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(0, 1, 0))
//	m1 := Scale(2, 3, 4)
//	r2 := TransformRay(r, m1)
//	assert.True(t, TupleEquals(r2.Origin, NewPoint(2, 6, 12)))
//	assert.True(t, TupleEquals(r2.Direction, geom.NewVector(0, 3, 0)))
//}

// Replaced in chapter 9
//func TestIntersectScaledSphereWithRay(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	s.SetTransform(Scale(2, 2, 2))
//	intersections := IntersectRayWithShape(s, r)
//	assert.Len(t, intersections, 2)
//	assert.Equal(t, 3.0, intersections[0].T)
//	assert.Equal(t, 7.0, intersections[1].T)
//}

//func TestIntersectScaledSphereWithRay2(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	s.SetTransform(Scale(2, 2, 2))
//	intersections := IntersectRayWithShape(s, r)
//	assert.Len(t, intersections, 2)
//	assert.Equal(t, s.GetLocalRay().Origin, geom.NewPoint(0, 0, -2.5))
//	assert.Equal(t, s.GetLocalRay().Direction, geom.NewVector(0, 0, 0.5))
//}
//
//func TestIntersectTranslatedSphereWithRay(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	s.SetTransform(geom.Translate(5, 0, 0))
//	intersections := IntersectRayWithShape(s, r)
//	assert.Len(t, intersections, 0)
//}
//
//func TestIntersectTranslatedSphereWithRay2(t *testing.T) {
//	r := geom.NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
//	s := shapes.NewSphere()
//	s.SetTransform(geom.Translate(5, 0, 0))
//	_ = IntersectRayWithShape(s, r)
//	assert.Equal(t, s.GetLocalRay().Origin, geom.NewPoint(-5, 0, -5))
//	assert.Equal(t, s.GetLocalRay().Direction, geom.NewVector(0, 0, 1))
//}
