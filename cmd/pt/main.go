package main

import (
	"flag"
	"fmt"
	"github.com/eriklupander/pathtracer/internal/app/scenes"
	"github.com/eriklupander/pathtracer/internal/app/tracer"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	var workers, width, height, samples int
	var rednerScene string
	var profile bool

	flag.IntVar(&workers, "workers", runtime.NumCPU(), "Number of workers")
	flag.IntVar(&width, "width", 640, "Width of screen")
	flag.IntVar(&height, "height", 480, "Height of screen")
	flag.IntVar(&samples, "samples", 256, "Number of samples to generate")
	flag.StringVar(&rednerScene, "scene", "reference", "Scene to use")
	flag.BoolVar(&profile, "profile", false, "Enable profiling")

	flag.Parse()

	fmt.Printf("Running with %d workers\n", workers)
	fmt.Printf("-- samples %d\n", samples)
	fmt.Printf("-- scene %s\n", rednerScene)

	if profile {
		f, err := os.Create("default.pgo")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	//runtime.SetBlockProfileRate(1)
	//runtime.SetMutexProfileFraction(1)
	// we need a webserver to get the pprof going
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var scene func() *scenes.Scene
	switch rednerScene {
	case "reference":
		scene = scenes.OCLScene(width, height)
	case "cornell":
		scene = scenes.OCLScene(width, height)
	default:
		scene = scenes.OCLScene(width, height)
	}

	t := tracer.NewPathTracer(width, height, workers)
	t.Render(scene)
}
