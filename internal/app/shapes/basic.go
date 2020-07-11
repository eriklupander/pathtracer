package shapes

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/material"
)

type Shape interface {
	ID() int64
	GetTransform() geom.Mat4x4
	GetInverse() geom.Mat4x4
	GetInverseTranspose() geom.Mat4x4
	SetTransform(transform geom.Mat4x4)
	GetMaterial() material.Material
	SetMaterial(material material.Material)
	IntersectLocal(ray geom.Ray, xs *Intersections)
	NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4
	GetLocalRay() geom.Ray
	GetParent() Shape
	SetParent(shape Shape)
	CastsShadow() bool
}

type Basic struct {
	Id               int64
	Transform        geom.Mat4x4
	Inverse          geom.Mat4x4
	InverseTranspose geom.Mat4x4
	Material         material.Material
}
