package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"fmt"
)

func main() {
	cfg := configapp.MustLoad()
	log := configapp.SetupLoger(cfg.Env)

	application := app.New(log, cfg.GRPC.Port)

	application.GRPC.MustRun()

	log.Debug("Init app")
	fmt.Println(cfg)
}
