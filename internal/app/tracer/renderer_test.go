package tracer

import (
	"fmt"
	"github.com/eriklupander/pathtracer/cmd"
	"github.com/eriklupander/pathtracer/internal/app/canvas"
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/scenes"
	"github.com/eriklupander/pathtracer/internal/app/shapes"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"sort"
	"sync"
	"testing"
)

func TestPathTracer_Render(t *testing.T) {
	cmd.FromConfig()
	cmd.Cfg.Width = 320
	cmd.Cfg.Height = 240
	jobs := make(chan *job)
	wg := sync.WaitGroup{}
	testee := NewCtx(1, scenes.OCLScene()(), canvas.NewCanvas(320, 240), jobs, &wg)
	var cameraRay = geom.NewEmptyRay()
	testee.rayForPixelPathTracer(210, 173, &cameraRay)

	color := testee.trace(cameraRay)
	assert.NotNil(t, color)
}

func TestMysort(t *testing.T) {
	xs := make([]shapes.Intersection, 4)
	xs[0].T = 4
	xs[1].T = -4
	xs[2].T = 2
	xs[3].T = 3

	xs = mysort(xs)
	assert.EqualValues(t, xs[0].T, -4)
	assert.EqualValues(t, xs[1].T, 2)
	assert.EqualValues(t, xs[2].T, 3)
	assert.EqualValues(t, xs[3].T, 4)
}

func TestMysort2(t *testing.T) {
	xs := make([]shapes.Intersection, 5)
	xs[0].T = 1
	xs[1].T = 1
	xs[2].T = -4
	xs[3].T = 2
	xs[4].T = -3

	xs = mysort(xs)
	assert.EqualValues(t, xs[0].T, -4)
	assert.EqualValues(t, xs[1].T, -3)
	assert.EqualValues(t, xs[2].T, 1)
	assert.EqualValues(t, xs[3].T, 1)
	assert.EqualValues(t, xs[4].T, 2)
}

func BenchmarkMySort(b *testing.B) {
	xs := make([]shapes.Intersection, 4)
	xs[0].T = 1
	xs[1].T = -4
	xs[2].T = 2
	xs[3].T = -3
	for i := 0; i < b.N; i++ {
		xs = mysort(xs)
	}
	fmt.Printf("%v\n", xs)
}

func BenchmarkSortSort(b *testing.B) {
	xs := make([]shapes.Intersection, 4)
	xs[0].T = 1
	xs[1].T = -4
	xs[2].T = 2
	xs[3].T = -3
	x := shapes.Intersections(xs)
	for i := 0; i < b.N; i++ {
		sort.Sort(x)
	}
	fmt.Printf("%v\n", x)
}

func BenchmarkBigMySort(b *testing.B) {
	xs := make([]shapes.Intersection, 7)
	for i := 0; i < 7; i++ {
		xs[i].T = rand.Float64() * 1000.0
	}
	for i := 0; i < b.N; i++ {
		xs = mysort(xs)
	}
	fmt.Printf("%v\n", xs)
}

func BenchmarkBigSortSort(b *testing.B) {
	xs := make([]shapes.Intersection, 7)
	for i := 0; i < 7; i++ {
		xs[i].T = rand.Float64() * 1000.0
	}
	x := shapes.Intersections(xs)
	for i := 0; i < b.N; i++ {
		sort.Sort(x)
	}
	fmt.Printf("%v\n", x)
}

func TestPositionOnCircle(t *testing.T) {

	step := (math.Pi * 2) / 8
	for i := 0; i < 8; i++ {
		fmt.Printf("degree: %v sin: %v cos: %v\n", RadToDeg*float64(float64(i)*step), math.Sin(float64(i)*(step)), math.Cos(float64(i)*(step)))
	}
	origin := geom.NewPoint(0, 0, 0)
	direction := geom.NewVector(1, 0, 0)
	newpos := positionOnCircle(origin, direction)
	assert.EqualValues(t, 0, newpos[2])

}

const (
	RadToDeg = 180 / math.Pi
	DegToRad = math.Pi / 180
)
