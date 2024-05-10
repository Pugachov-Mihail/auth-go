package app

import (
	grpcapp "auth/internal/app/grpc"
	configapp "auth/internal/config"
	kafka_user "auth/internal/kafka"
	"auth/internal/service/auth"
	resetservice "auth/internal/service/reset"
	authstorage "auth/internal/storage/auth"
	"auth/internal/storage/reset"
	"log/slog"
	"time"
)

type App struct {
	GRPC *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath configapp.ConfigDB, tokenTT time.Duration, cfg *configapp.Config) *App {
	storage, err := authstorage.New(storagePath)
	resetStorage, err := reset.New(storagePath)
	if err != nil {
		panic(err)
	}

	conf, err := kafka_user.New(cfg)
	if err != nil {
		log.Warn("Ошибка коннекта с кафкой")
		return nil
	}
	authService := auth.New(log, storage, storage, storage, tokenTT, cfg, conf)
	resetService := resetservice.New(log, resetStorage)
	grpcAuth := grpcapp.New(log, port, authService, resetService)

	return &App{GRPC: grpcAuth}
}
