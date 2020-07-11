package scenes

import (
	"github.com/eriklupander/pt/cmd"
	"github.com/eriklupander/pt/internal/app/camera"
	"github.com/eriklupander/pt/internal/app/geom"
	"github.com/eriklupander/pt/internal/app/material"
	"github.com/eriklupander/pt/internal/app/shapes"
	"math"
)

func ReferenceScene() func() *Scene {
	return func() *Scene {

		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(-2, 2.0, -4), geom.NewPoint(0, 0.5, 0))
		//lightBulb := shapes.NewSphere()
		//lightBulb.SetTransform(geom.Translate(0, 5, -1))
		//
		//lightBulb.Material = material.NewLightBulb()
		//
		//s1 := shapes.NewSphere()
		//s1.Material = material.NewDiffuse(1,.7, .7)
		//
		//floor := shapes.NewPlane()
		//floor.SetTransform(geom.Translate(0, -1, 0))
		//floor.Material = material.NewDiffuse(.9,.8, .95)

		lightBulb := shapes.NewSphere()
		lightBulb.SetMaterial(material.NewLightBulb())
		lightBulb.SetTransform(geom.Translate(-5, 4.5, -3))
		lightBulb.SetTransform(geom.Scale(1.5, 0.5, 1.5))

		floor := shapes.NewPlane()
		floor.SetTransform(geom.Translate(0, 0.01, 0))
		floor.SetMaterial(material.NewDiffuse(1, 0.5, 0.5))

		ceil := shapes.NewPlane()
		ceil.SetTransform(geom.Translate(0, 5, 0))
		ceilMat := material.NewDefaultMaterial()
		ceil.SetMaterial(ceilMat)

		wall := shapes.NewPlane()
		wall.SetMaterial(material.NewDiffuse(0.9, 0.9, 0.9))
		wall.SetTransform(geom.Translate(0, 0, 8))
		wall.SetTransform(geom.RotateX(math.Pi / 2))

		// transparent sphere
		middle := shapes.NewSphere()
		middle.SetTransform(geom.Translate(-0.5, 0.75, 0.5))
		middle.SetTransform(geom.Scale(0.75, 0.75, 0.75))
		glassMtrl := material.NewDiffuse(0.8, 0.8, 0.9)
		middle.SetMaterial(glassMtrl)

		s1 := shapes.NewSphere()
		s1.SetTransform(geom.Multiply(geom.Translate(-2, 0.25, -1), geom.Scale(0.25, 0.25, 0.25)))
		mat1 := material.NewDiffuse(1, 0.1, 0.1)
		s1.SetMaterial(mat1)

		s2 := shapes.NewSphere()
		//s2.CastShadow = false
		s2.SetTransform(geom.Multiply(geom.Translate(-1, 0.25, -1), geom.Scale(0.25, 0.25, 0.25)))
		mat2 := material.NewDiffuse(0.1, 1.0, 0.1)
		s2.SetMaterial(mat2)

		s3 := shapes.NewSphere()
		s3.SetTransform(geom.Multiply(geom.Translate(0, 0.25, -1), geom.Scale(0.25, 0.25, 0.25)))
		mat3 := material.NewDiffuse(0.1, 0.1, 1)
		s3.SetMaterial(mat3)

		return &Scene{
			Camera: cam,
			Objects: []shapes.Shape{
				lightBulb,
				floor,
				ceil,
				middle,
				wall,
				s1,
				s2,
				s3,
			},
		}
	}
}
