package scenes

import (
	"github.com/eriklupander/pt/internal/app/camera"
	"github.com/eriklupander/pt/internal/app/shapes"
)

type Scene struct {
	Camera  camera.Camera
	Objects []shapes.Shape
}
