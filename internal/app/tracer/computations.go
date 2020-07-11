package tracer

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/shapes"
)

func NewComputation() Computation {
	containers := make([]shapes.Shape, 8)
	containers = containers[:0]
	return Computation{
		T:          0,
		Object:     nil,
		Point:      geom.NewPoint(0, 0, 0),
		EyeVec:     geom.NewVector(0, 0, 0),
		NormalVec:  geom.NewVector(0, 0, 0),
		Inside:     false,
		OverPoint:  geom.NewPoint(0, 0, 0),
		UnderPoint: geom.NewPoint(0, 0, 0),
		ReflectVec: geom.NewVector(0, 0, 0),
		N1:         0,
		N2:         0,

		localPoint:   geom.NewPoint(0, 0, 0),
		containers:   containers,
		cachedOffset: geom.NewVector(0, 0, 0),
	}
}

func PrepareComputationForIntersectionPtr(i shapes.Intersection, r geom.Ray, comps *Computation, xs ...shapes.Intersection) {
	comps.T = i.T
	comps.Object = i.S
	PositionPtr(r, i.T, &comps.Point)
	geom.NegatePtr(r.Direction, &comps.EyeVec)
	comps.NormalVec = NormalAt(i.S, comps.Point, &i) //  fix
	//comps.NormalVec = NormalAtPtr(i.S, comps.Point, &i, &comps.localPoint) //  fix
	ReflectPtr(r.Direction, comps.NormalVec, &comps.ReflectVec)

	comps.Inside = false
	if geom.Dot(comps.EyeVec, comps.NormalVec) < 0 {
		comps.Inside = true
		geom.NegatePtr(comps.NormalVec, &comps.NormalVec) // fix
	}
	geom.MultiplyByScalarPtr(comps.NormalVec, geom.Epsilon, &comps.cachedOffset)
	geom.AddPtr(comps.Point, comps.cachedOffset, &comps.OverPoint)
	geom.SubPtr(comps.Point, comps.cachedOffset, &comps.UnderPoint)

	comps.N1 = 1.0
	comps.N2 = 1.0

	comps.containers = comps.containers[:0] // make([]Shape, 0)
	for idx := range xs {
		if xs[idx].S.ID() == i.S.ID() && i.T == xs[idx].T {
			if len(comps.containers) == 0 {
				comps.N1 = 1.0
			} else {
				comps.N1 = comps.containers[len(comps.containers)-1].GetMaterial().RefractiveIndex
			}
		}

		index := indexOf(xs[idx].S, comps.containers)
		if index > -1 {
			copy(comps.containers[index:], comps.containers[index+1:])    // Shift a[i+1:] left one indexs[idx].
			comps.containers[len(comps.containers)-1] = nil               // Erase last element (write zero value).
			comps.containers = comps.containers[:len(comps.containers)-1] // Truncate slice.
		} else {
			comps.containers = append(comps.containers, xs[idx].S)
		}

		if xs[idx].S.ID() == i.S.ID() && xs[idx].T == i.T {
			if len(comps.containers) == 0 {
				comps.N2 = 1.0
			} else {
				comps.N2 = comps.containers[len(comps.containers)-1].GetMaterial().RefractiveIndex
			}
			break
		}
	}
}

func indexOf(s shapes.Shape, list []shapes.Shape) int {
	for idx := range list {
		if list[idx].ID() == s.ID() {
			return idx
		}
	}
	return -1
}

type Computation struct {
	T          float64
	Object     shapes.Shape
	Point      geom.Tuple4
	EyeVec     geom.Tuple4
	NormalVec  geom.Tuple4
	Inside     bool
	OverPoint  geom.Tuple4
	UnderPoint geom.Tuple4
	ReflectVec geom.Tuple4
	N1         float64
	N2         float64

	// cached stuff
	localPoint   geom.Tuple4
	containers   []shapes.Shape
	cachedOffset geom.Tuple4
}

func NewLightData() LightData {
	return LightData{
		//Ambient:        rgb.NewColor(0, 0, 0),
		//Diffuse:        rgb.NewColor(0, 0, 0),
		//Specular:       geom.NewColor(0, 0, 0),
		//EffectiveColor: geom.NewColor(0, 0, 0),
		LightVec:   geom.NewVector(0, 0, 0),
		ReflectVec: geom.NewVector(0, 0, 0),
	}
}

// LightData is used for pre-allocated memory for lighting computations.
type LightData struct {
	//Ambient        Tuple4
	//Diffuse        Tuple4
	//Specular       Tuple4
	//EffectiveColor Tuple4
	LightVec   geom.Tuple4
	ReflectVec geom.Tuple4
}
