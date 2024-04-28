package app

import (
	grpcapp "auth/internal/app/grpc"
	configapp "auth/internal/config"
	"auth/internal/service/auth"
	reset_service "auth/internal/service/reset"
	auth_storage "auth/internal/storage/auth"
	"auth/internal/storage/reset"
	"log/slog"
	"time"
)

type App struct {
	GRPC *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath configapp.ConfigDB, tokenTT time.Duration, cfg *configapp.Config) *App {
	storage, err := auth_storage.New(storagePath)
	resetStorage, err := reset.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, tokenTT, cfg)
	resetService := reset_service.New(log, resetStorage)
	grpcAuth := grpcapp.New(log, port, authService, resetService)

	return &App{GRPC: grpcAuth}
}
