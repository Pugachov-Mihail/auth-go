package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := configapp.MustLoad()
	log := configapp.SetupLoger(cfg.Env)
	log.Debug("Init app")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.GRPC.TimeOut, cfg)

	go application.GRPC.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stop app", slog.String("signal", sign.String()))

	application.GRPC.Stop()
	log.Info("app stop")
}
