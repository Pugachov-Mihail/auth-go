package server_reset

import (
	resetService "auth/protos/gen/dota_traker.reset.v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.3 --all

type Reset struct {
	resetService.UnimplementedResetAuthDataServer
	reset ResetStorage
}

type ResetStorage interface {
	ResetPassword(ctx context.Context, password string, userId int64) (int64, error)
	ResetIdSteam(ctx context.Context, steamId int64, userId int64) (int64, error)
	ResetEmailStore(ctx context.Context, email string, userId int64) (int64, error)
}

func RegisterResetServerApi(grpc *grpc.Server, storage ResetStorage) {
	resetService.RegisterResetAuthDataServer(grpc, &Reset{reset: storage})
}

func (r *Reset) ResetSteamId(ctx context.Context, req *resetService.ResetSteamIdRequests) (*resetService.ResetResponse, error) {
	if req.GetSteamId() != 0 || req.GetIdUser() != 0 {
		return nil, status.Error(codes.InvalidArgument, "Нет данных")
	}

	uid, err := r.reset.ResetIdSteam(ctx, req.GetSteamId(), req.GetIdUser())

	if err != nil {
		return nil, status.Error(codes.Internal, "Ошибка изменения Steam id")
	}

	return &resetService.ResetResponse{UserId: uid}, nil
}

func (r *Reset) ResetEmail(ctx context.Context, req *resetService.ResetEmailRequests) (*resetService.ResetResponse, error) {
	if req.GetEmail() == "" || req.GetIdUser() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Нет данных")
	}

	uid, err := r.reset.ResetEmailStore(ctx, req.GetEmail(), req.GetIdUser())
	if err != nil {
		return nil, status.Error(codes.Internal, "Ошибка изменения почты")
	}

	return &resetService.ResetResponse{UserId: uid}, nil
}

func (r *Reset) ResetPassword(ctx context.Context, req *resetService.ResetPasswordRequests) (*resetService.ResetResponse, error) {
	if req.GetPassword() != "" || req.GetPassword2() != "" {
		return nil, status.Error(codes.InvalidArgument, "Нет данных")
	}

	if req.GetPassword() == req.GetPassword2() {
		return nil, status.Error(codes.InvalidArgument, "Пароли не совпадают")
	}

	uid, err := r.reset.ResetPassword(ctx, req.GetPassword(), req.GetIdUser())
	if err != nil {
		return nil, status.Error(codes.Internal, "Ошибка изменения пароля")
	}

	return &resetService.ResetResponse{UserId: uid}, nil
}
