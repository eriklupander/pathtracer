package tracer

import (
	"github.com/eriklupander/pt/internal/app/geom"
	"github.com/eriklupander/pt/internal/app/shapes"
	"math"
)

// Position multiplies direction of ray with the passed distance and adds the result onto the origin.
// Used for finding the position along a ray.
func Position(r geom.Ray, distance float64) geom.Tuple4 {
	add := geom.MultiplyByScalar(r.Direction, distance)
	pos := geom.Add(r.Origin, add)
	return pos
}

func PositionPtr(r geom.Ray, distance float64, out *geom.Tuple4) {
	add := geom.MultiplyByScalar(r.Direction, distance)
	geom.AddPtr(r.Origin, add, out)
}

// TODO only used by unit tests, fix so tests use IntersectRayWithShapePtr and remove
//func IntersectRayWithShape(s shapes.Shape, r2 geom.Ray) []shapes.Intersection {
//
//	// transform ray with inverse of shape transformation matrix to be able to intersect a translated/rotated/skewed shape
//	r := TransformRay(r2, s.GetInverse())
//
//	// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
//	return s.IntersectLocal(r)
//}

//func IntersectRayWithShapePtr(s shapes.Shape, r2 geom.Ray, in *geom.Ray) []shapes.Intersection {
//	//calcstats.Incr()
//	// transform ray with inverse of shape transformation matrix to be able to intersect a translated/rotated/skewed shape
//	TransformRayPtr(r2, s.GetInverse(), in)
//
//	// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
//	return s.IntersectLocal(*in)
//}

// Hit finds the first intersection with a positive T (the passed intersections are assumed to have been sorted already)
func Hit(intersections []shapes.Intersection) (shapes.Intersection, bool) {

	// Filter out all negatives
	for i := 0; i < len(intersections); i++ {
		if intersections[i].T > 0.0 {
			return intersections[i], true
		}
	}

	return shapes.Intersection{}, false
}

func TransformRay(r geom.Ray, m1 geom.Mat4x4) geom.Ray {
	origin := geom.MultiplyByTuple(m1, r.Origin)
	direction := geom.MultiplyByTuple(m1, r.Direction)
	return geom.NewRay(origin, direction)
}

func TransformRayPtr(r geom.Ray, m1 geom.Mat4x4, out *geom.Ray) {
	geom.MultiplyByTuplePtr(&m1, &r.Origin, &out.Origin)
	geom.MultiplyByTuplePtr(&m1, &r.Direction, &out.Direction)
}

func Schlick(comps Computation) float64 {
	// find the cosine of the angle between the eye and normal vectors using Dot
	cos := geom.Dot(comps.EyeVec, comps.NormalVec)
	// total internal reflection can only occur if n1 > n2
	if comps.N1 > comps.N2 {
		n := comps.N1 / comps.N2
		sin2Theta := (n * n) * (1.0 - (cos * cos))
		if sin2Theta > 1.0 {
			return 1.0
		}
		// compute cosine of theta_t using trig identity
		cosTheta := math.Sqrt(1.0 - sin2Theta)

		// when n1 > n2, use cos(theta_t) instead
		cos = cosTheta
	}
	temp := (comps.N1 - comps.N2) / (comps.N1 + comps.N2)
	r0 := temp * temp
	return r0 + (1-r0)*math.Pow(1-cos, 5)
}

//
//func IntersectWithWorldPtr(scene *scenes.Scene, r geom.Ray, intersections shapes.Intersections, inRay *geom.Ray) []shapes.Intersection {
//	for idx := range scene.Objects {
//		intersections := IntersectRayWithShapePtr(scene.Objects[idx], r, inRay)
//
//		for innerIdx := range intersections {
//			intersections = append(intersections, intersections[innerIdx])
//		}
//	}
//	if len(intersections) > 1 {
//		sort.Sort(intersections)
//	}
//	return intersections
//}
