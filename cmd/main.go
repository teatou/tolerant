package main

import (
	"os"

	"github.com/teatou/tolerant/internal/config"
	"github.com/teatou/tolerant/pkg/mylogger"
)

const configEnv = "CONFIG"

func main() {
	val, ok := os.LookupEnv(configEnv)
	if !ok {
		val = "configs/dev.yaml"
	}

	cfg, err := config.LoadConfig(val)
	if err != nil {
		panic("uploading config error")
	}

	logger, err := mylogger.NewZapLogger(cfg.Logger.Level)
	if err != nil {
		panic("making mylogger error")
	}
	defer logger.Sync()

	// init storage
	// init cache

	// router

	// server
}
