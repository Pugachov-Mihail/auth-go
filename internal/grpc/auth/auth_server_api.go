package auth_server

import (
	"auth/internal/domain/models"
	"auth/internal/validator/auth_validate"
	"auth/internal/validator/base_validate"
	authServer "auth/protos/gen/dota_traker.auth.v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.3 --all

type Auth interface {
	RegisterUser(
		ctx context.Context,
		login string,
		password string,
		email string,
		steamId int64) (int64, error)
	LoginUser(ctx context.Context, login string, password string) (string, error)
	RolesUser(ctx context.Context, uid int64) (models.Roles, error)
	AccessPermission(ctx context.Context, token string) (string, error)
}

type AuthServerApi struct {
	authServer.UnimplementedAuthServerServer
	auth Auth
}

func RegisterAuthServerApi(grpc *grpc.Server, auth Auth) {
	authServer.RegisterAuthServerServer(grpc, &AuthServerApi{auth: auth})
}

func (a *AuthServerApi) AuthLogin(
	ctx context.Context, req *authServer.AuthLoginRequest) (*authServer.AuthLoginResponse, error) {
	if !base_validate.ValidateLoginRequest(req) {
		return nil, status.Error(codes.InvalidArgument, "Пустые данные")
	}

	token, err := a.auth.LoginUser(ctx, req.GetLogin(), req.GetPassword())

	if err != nil {
		return nil, status.Error(codes.Internal, "Ошибка авторизации")
	}

	return &authServer.AuthLoginResponse{Token: token}, nil
}

func (a *AuthServerApi) AuthRegistration(
	ctx context.Context, req *authServer.AuthRegistrationRequest) (*authServer.AuthRegistrationResponse, error) {

	if ok, _ := auth_validate.ValidateEmail(req.GetEmail()); !ok {
		return nil, status.Error(codes.InvalidArgument, "Некорректная почта")
	}

	if !base_validate.ValidatePassword(req) {
		return nil, status.Error(codes.InvalidArgument, "Пароли не совпадают")
	}

	id, err := a.auth.RegisterUser(ctx, req.GetLogin(), req.GetPassword(), req.GetEmail(), req.GetSteamId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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

func (a *AuthServerApi) AuthAccessPermission(
	ctx context.Context,
	req *authServer.AuthLoginResponse) (*authServer.AccessPermissionResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует токен")
	}

	permission, err := a.auth.AccessPermission(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "Отсутствует доступ")
	}

	return &authServer.AccessPermissionResponse{AccessPermission: permission}, nil
}
