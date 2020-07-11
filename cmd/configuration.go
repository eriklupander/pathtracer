package cmd

import "github.com/spf13/viper"

type Config struct {
	Width   int
	Height  int
	Workers int
	Samples int
}

var Cfg *Config

func FromConfig() {
	Cfg = &Config{
		Width:   viper.GetInt("width"),
		Height:  viper.GetInt("height"),
		Workers: viper.GetInt("workers"),
		Samples: viper.GetInt("samples"),
	}
}
