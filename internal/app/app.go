package app

import (
	grpcapp "auth/internal/app/grpc"
	configapp "auth/internal/config"
	"auth/internal/service/auth"
	auth_storage "auth/internal/storage/auth"
	"log/slog"
	"time"
)

type App struct {
	GRPC *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath configapp.ConfigDB, tokenTT time.Duration, cfg *configapp.Config) *App {
	storage, err := auth_storage.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, tokenTT, cfg)
	grpcAuth := grpcapp.New(log, port, authService)

	return &App{GRPC: grpcAuth}
}
