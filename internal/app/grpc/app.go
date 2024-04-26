package grpcapp

import (
	serverAdmin "auth/internal/grpc/admin"
	serverAuth "auth/internal/grpc/auth"
	server_reset "auth/internal/grpc/reset"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, authServese serverAuth.Auth) *App {
	grpcServer := grpc.NewServer()
	serverAdmin.RegisterServerApi(grpcServer)
	serverAuth.RegisterAuthServerApi(grpcServer, authServese)
	server_reset.RegisterResetServerApi(grpcServer)
	reflection.Register(grpcServer)

	return &App{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *App) Run() error {
	log := a.log.With("Auth service Running", a.port)

	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", a.port))

	if err != nil {
		log.Error("Error Run Auth service")
		return fmt.Errorf("%s: %w", "Auth Run", err)
	}

	log.Info("Auth Run", slog.String("addr", l.Addr().String()))

	if err := a.grpcServer.Serve(l); err != nil {
		log.Error("Error Run GRPC Auth service")
		return fmt.Errorf("%s: %w", "Auth Run", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.With("Auth service Spoping", a.port).
		Info("Auth GRPC stop")

	a.grpcServer.GracefulStop()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
