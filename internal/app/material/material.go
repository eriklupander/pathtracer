package material

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
)

type Material struct {
	Color           geom.Tuple4
	Emission        geom.Tuple4
	RefractiveIndex float64
	Reflectivity    float64
}

func NewDefaultMaterial() Material {
	return Material{
		Color:           geom.Tuple4{1, 1, 1},
		Emission:        geom.Tuple4{0, 0, 0},
		RefractiveIndex: 0,
	}
}

func NewDiffuse(r, g, b float64) Material {
	return Material{
		Color:           geom.Tuple4{r, g, b},
		Emission:        geom.Tuple4{0, 0, 0},
		RefractiveIndex: 0,
	}
}
func NewMirror() Material {
	return Material{
		Color:           geom.Tuple4{1, 1, 1},
		Emission:        geom.Tuple4{0, 0, 0},
		RefractiveIndex: 0,
		Reflectivity:    1.0,
	}
}
func NewLightBulb() Material {
	return Material{
		Color:           geom.Tuple4{0, 0, 0},
		Emission:        geom.Tuple4{8, 8, 8},
		RefractiveIndex: 0,
	}
}
