package main

import (
	"github.com/eriklupander/pathtracer/cmd"
	"github.com/eriklupander/pathtracer/internal/app/scenes"
	"github.com/eriklupander/pathtracer/internal/app/tracer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
)

func main() {

	var configFlags = pflag.NewFlagSet("config", pflag.ExitOnError)
	configFlags.Int("workers", runtime.NumCPU(), "number of workers")
	configFlags.Int("width", 640, "Image width")
	configFlags.Int("height", 480, "Image height")
	configFlags.Int("samples", 1, "Number of samples per pixel")
	configFlags.String("scene", "reference", "scene from /scenes")

	if err := configFlags.Parse(os.Args[1:]); err != nil {
		panic(err.Error())
	}
	if err := viper.BindPFlags(configFlags); err != nil {
		panic(err.Error())
	}
	viper.AutomaticEnv()

	cmd.FromConfig()
	logrus.Printf("Running with %d CPUs\n", viper.GetInt("workers"))

	//runtime.SetBlockProfileRate(1)
	//runtime.SetMutexProfileFraction(1)
	// we need a webserver to get the pprof going
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var scene func() *scenes.Scene
	switch viper.GetString("scene") {
	case "reference":
		scene = scenes.OCLScene()
	case "cornell":
		scene = scenes.OCLScene()
	default:
		scene = scenes.OCLScene()
	}

	t := tracer.NewPathTracer()
	t.Render(scene)
}
