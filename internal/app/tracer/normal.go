package tracer

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/shapes"
)

func NormalAtPtr(s shapes.Shape, worldPoint geom.Tuple4, intersection *shapes.Intersection, localPoint *geom.Tuple4) geom.Tuple4 {

	// transform point from world to object space, including recursively traversing any parent object
	// transforms.
	WorldToObjectPtr(s, worldPoint, localPoint)

	// normal in local space given the shape's implementation
	objectNormal := s.NormalAtLocal(*localPoint, intersection)

	// convert normal from object space back into world space, again recursively applying any
	// parent transforms.
	NormalToWorldPtr(s, &objectNormal)
	return objectNormal
}

// in - normal * 2 * dot(in, normal)
func Reflect(vec geom.Tuple4, normal geom.Tuple4) geom.Tuple4 {
	dotScalar := geom.Dot(vec, normal)
	norm := geom.MultiplyByScalar(geom.MultiplyByScalar(normal, 2.0), dotScalar)
	return geom.Sub(vec, norm)
}

// in - normal * 2 * dot(in, normal)
func ReflectPtr(vec geom.Tuple4, normal geom.Tuple4, out *geom.Tuple4) {
	dotScalar := geom.Dot(vec, normal)
	norm := geom.MultiplyByScalar(geom.MultiplyByScalar(normal, 2.0), dotScalar)
	geom.SubPtr(vec, norm, out)
}

func NormalAt(s shapes.Shape, worldPoint geom.Tuple4, intersection *shapes.Intersection) geom.Tuple4 {

	// transform point from world to object space, including recursively traversing any parent object
	// transforms.
	localPoint := WorldToObject(s, worldPoint)

	// normal in local space given the shape's implementation
	objectNormal := s.NormalAtLocal(localPoint, intersection)

	// convert normal from object space back into world space, again recursively applying any
	// parent transforms.
	return NormalToWorld(s, objectNormal)
}

func WorldToObject(shape shapes.Shape, point geom.Tuple4) geom.Tuple4 {
	if shape.GetParent() != nil {
		point = WorldToObject(shape.GetParent(), point)
	}
	return geom.MultiplyByTuple(shape.GetInverse(), point)
}

func WorldToObjectPtr(shape shapes.Shape, point geom.Tuple4, out *geom.Tuple4) {
	if shape.GetParent() != nil {
		WorldToObjectPtr(shape.GetParent(), point, &point)
	}
	i := shape.GetInverse()
	geom.MultiplyByTuplePtr(&i, &point, out)
}

func NormalToWorld(shape shapes.Shape, normal geom.Tuple4) geom.Tuple4 {
	normal = geom.MultiplyByTuple(shape.GetInverseTranspose(), normal)
	normal[3] = 0.0 // set w to 0
	normal = geom.Normalize(normal)

	if shape.GetParent() != nil {
		normal = NormalToWorld(shape.GetParent(), normal)
	}
	return normal
}

func NormalToWorldPtr(shape shapes.Shape, normal *geom.Tuple4) {
	it := shape.GetInverseTranspose()
	geom.MultiplyByTuplePtr(&it, normal, normal)
	normal[3] = 0.0 // set w to 0
	geom.NormalizePtr(normal, normal)

	if shape.GetParent() != nil {
		NormalToWorldPtr(shape.GetParent(), normal)
	}
}
