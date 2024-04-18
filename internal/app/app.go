package app

import (
	grpcapp "auth/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPC *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	grpcAuth := grpcapp.New(log, port)
	return &App{GRPC: grpcAuth}
}
