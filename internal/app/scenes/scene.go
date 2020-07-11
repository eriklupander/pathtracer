package scenes

import (
	"github.com/eriklupander/pathtracer/internal/app/camera"
	"github.com/eriklupander/pathtracer/internal/app/shapes"
)

type Scene struct {
	Camera  camera.Camera
	Objects []shapes.Shape
}
