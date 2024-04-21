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

func New(log *slog.Logger, port int, storagePath configapp.ConfigDB, tokenTT time.Duration) *App {
	storage, err := auth_storage.New(storagePath)
	if err != nil {
		panic(err)
	}
	auth_service := auth.New(log, storage, storage, tokenTT)
	grpcAuth := grpcapp.New(log, port, auth_service)
	return &App{GRPC: grpcAuth}
}
