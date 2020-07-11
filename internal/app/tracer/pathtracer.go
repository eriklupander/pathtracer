package tracer

import (
	"fmt"
	"github.com/eriklupander/pathtracer/cmd"
	canvas2 "github.com/eriklupander/pathtracer/internal/app/canvas"
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/scenes"
	"image"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

var originPoint = geom.NewPoint(0, 0, 0)
var black = geom.NewColor(0, 0, 0)
var white = geom.NewColor(1, 1, 1)

type PathTracer struct {

}

func (t *PathTracer) Render(sceneFactory func() *scenes.Scene) {

	st := time.Now()
	canvas := canvas2.NewCanvas(cmd.Cfg.Width, cmd.Cfg.Height)
	jobs := make(chan *job)

	wg := sync.WaitGroup{}
	wg.Add(canvas.H)

	// Create the render contexts, one per worker
	renderContexts := make([]*Ctx, cmd.Cfg.Workers)
	for i := 0; i < cmd.Cfg.Workers; i++ {
		renderContexts[i] = NewCtx(i, sceneFactory(), canvas, jobs, &wg)
	}

	// start workers
	for i := 0; i < cmd.Cfg.Workers; i++ {
		go renderContexts[i].workerFuncPerLine()
	}

	// start passing work to the workers, one line at a time
	for row := 0; row < cmd.Cfg.Height; row++ {
		jobs <- &job{row: row, col: 0}
		fmt.Printf("%d/%d\n", row, cmd.Cfg.Height)
	}

	wg.Wait()
	done := time.Now().Sub(st)
	fmt.Println("All done")
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	fmt.Printf("Finished in %v\n", time.Now().Sub(st))
	fmt.Printf("Samples taken: %v\n", sampleCount)
	fmt.Printf("Samples/s: %v\n", float64(sampleCount)/done.Seconds())
	writeImagePNG(canvas, "out.png")
}

func NewPathTracer() *PathTracer {
	return &PathTracer{}
}

func writeImagePNG(canvas *canvas2.Canvas, filename string) {
	fmt.Printf("writing output to file %v\n", filename)
	myImage := image.NewRGBA(image.Rect(0, 0, canvas.W, canvas.H))
	writeDataToPNG(canvas, myImage)
	outputFile, _ := os.Create(filename)
	defer outputFile.Close()
	_ = png.Encode(outputFile, myImage)
}

func writeDataToPNG(canvas *canvas2.Canvas, myImage *image.RGBA) {
	for i := 0; i < len(canvas.Pixels); i++ {
		myImage.Pix[i*4] = clamp(canvas.Pixels[i][0])
		myImage.Pix[i*4+1] = clamp(canvas.Pixels[i][1])
		myImage.Pix[i*4+2] = clamp(canvas.Pixels[i][2])
		myImage.Pix[i*4+3] = 255
	}
}

type job struct {
	row int
	col int
}

func clamp(clr float64) uint8 {
	c := clr * 255.0
	rounded := math.Round(c)
	if rounded > 255.0 {
		rounded = 255.0
	} else if rounded < 0.0 {
		rounded = 0.0
	}
	return uint8(rounded)
}
