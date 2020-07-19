package scenes

import (
	"github.com/eriklupander/pathtracer/cmd"
	"github.com/eriklupander/pathtracer/internal/app/camera"
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/material"
	"github.com/eriklupander/pathtracer/internal/app/shapes"
	"math"
)

func OCLScene() func() *Scene {
	return func() *Scene {

		cam := camera.NewCamera(cmd.Cfg.Width, cmd.Cfg.Height, math.Pi/3, geom.NewPoint(0, 0.1, -1.5), geom.NewPoint(0, 0.05, 0))

		//middleSphere := make([]shapes.Shape, 9)

		// left wall
		leftWall := shapes.NewPlane()
		leftWall.SetTransform(geom.Translate(-.6, 0, 0))
		leftWall.SetTransform(geom.RotateZ(math.Pi / 2))
		leftWall.SetMaterial(material.NewDiffuse(0.75, 0.25, 0.25))
		//
		//// right wall
		rightWall := shapes.NewPlane()
		rightWall.SetTransform(geom.Translate(.6, 0, 0))
		rightWall.SetTransform(geom.RotateZ(math.Pi / 2))
		rightWall.SetMaterial(material.NewDiffuse(0.25, 0.25, 0.75))

		// floor
		floor := shapes.NewPlane()
		floor.SetTransform(geom.Translate(0, -.4, 0))
		floor.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// ceiling
		ceil := shapes.NewPlane()
		ceil.SetTransform(geom.Translate(0, .4, 0))
		ceil.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// back wall
		backWall := shapes.NewPlane()
		backWall.SetTransform(geom.Translate(0, 0, .4))
		backWall.SetTransform(geom.RotateX(math.Pi / 2))
		backWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// front wall
		frontWall := shapes.NewPlane()
		frontWall.SetTransform(geom.Translate(0, 0, -2))
		frontWall.SetTransform(geom.RotateX(math.Pi / 2))
		frontWall.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))

		// left sphere
		leftSphere := shapes.NewSphere()
		leftSphere.SetTransform(geom.Translate(-0.25, -0.24, 0.1))
		leftSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		leftSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		//leftSphere.SetMaterial(material.NewMirror())

		// middle sphere
		middleSphere := shapes.NewSphere()
		middleSphere.SetTransform(geom.Translate(0, -0.24, -0.30))
		middleSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		//		middleSphere.SetMaterial(material.NewDiffuse(0.9, 0.8, 0.7))
		middleSphere.SetMaterial(material.NewGlass())

		// middle sphere
		rightSphere := shapes.NewSphere()
		rightSphere.SetTransform(geom.Translate(0.25, -0.24, 0.1))
		rightSphere.SetTransform(geom.Scale(0.16, 0.16, 0.16))
		rightSphere.SetMaterial(material.NewDiffuse(0.57, 0.86, 1))
		//middleSphere.SetMaterial(material.NewGlass())

		// lightsource
		lightsource := shapes.NewSphere()
		lightsource.SetTransform(geom.Translate(0, 1.36, 0))
		light := material.NewLightBulb()
		light.Emission = geom.NewColor(9, 8, 6)
		lightsource.SetMaterial(light)

		return &Scene{
			Camera:  cam,
			Objects: []shapes.Shape{floor, ceil, leftWall, rightWall, backWall, leftSphere, middleSphere, rightSphere, lightsource},
		}
	}

}
