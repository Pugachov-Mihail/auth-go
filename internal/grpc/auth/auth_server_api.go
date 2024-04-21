package auth_server

import (
	configapp "auth/internal/config"
	"auth/internal/domain/models"
	"auth/internal/validator"
	authServer "auth/protos/gen/dota_traker.auth.v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Auth interface {
	RegisterUser(
		ctx context.Context,
		login string,
		password string,
		email string,
		steamId int64) (int64, error)
	LoginUser(ctx context.Context, login string, password string, secret string) (string, error)
	RolesUser(ctx context.Context, uid int64) (models.Roles, error)
}

type AuthServerApi struct {
	authServer.UnimplementedAuthServerServer
	auth   Auth
	secret configapp.Config
}

func RegisterAuthServerApi(grpc *grpc.Server, auth Auth) {
	authServer.RegisterAuthServerServer(grpc, &AuthServerApi{auth: auth})
}

func (a *AuthServerApi) AuthLogin(
	ctx context.Context, req *authServer.AuthLoginRequest) (*authServer.AuthLoginResponse, error) {
	var log *slog.Logger
	log.With("Init Login logger")

	if !validator.ValidateLoginRequest(req) {
		log.Warn("Пустые данные")
		return nil, status.Error(codes.InvalidArgument, "Пустые данные")
	}

	token, err := a.auth.LoginUser(ctx, req.GetLogin(), req.GetPassword(), a.secret.Secret)

	if err != nil {
		return nil, status.Error(codes.Internal, "Ошибка авторизации")
	}

	return &authServer.AuthLoginResponse{Token: token}, nil
}

func (a *AuthServerApi) AuthRegistration(
	ctx context.Context, req *authServer.AuthRegistrationRequest) (*authServer.AuthRegistrationResponse, error) {
	var log *slog.Logger
	log.With("Init Registration")

	if !validator.ValidatePassword(req) {
		log.Warn("Пароли не совпадают")
		return nil, status.Error(codes.InvalidArgument, "Пароли не совпадают")
	}

	id, err := a.auth.RegisterUser(ctx, req.GetLogin(), req.GetPassword(), req.GetEmail(), req.GetSteamId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Ошибка создания пользователя")
	}

	return &authServer.AuthRegistrationResponse{UserId: id}, nil
}

func (a *AuthServerApi) AuthRoles(
	ctx context.Context, req *authServer.AuthRolesRequest) (*authServer.AuthRolesResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Пустой айди пользователя")
	}

	roles, err := a.auth.RolesUser(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Ошибка роли")
	}

	return &authServer.AuthRolesResponse{RolesFlag: roles.RolesFlag, RoleName: roles.RolesName}, nil
}
