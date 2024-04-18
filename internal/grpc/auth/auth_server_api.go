package auth

import (
	authServer "auth/protos/gen/dota_traker.auth.v1"
	"context"
	"google.golang.org/grpc"
)

type AuthServerApi struct {
	authServer.UnimplementedAuthServerServer
}

func RegisterAuthServerApi(grpc *grpc.Server) {
	authServer.RegisterAuthServerServer(grpc, &AuthServerApi{})
}

func (a *AuthServerApi) AuthLogin(
	ctx context.Context, req *authServer.AuthLoginRequest) (*authServer.AuthRegistrationResponse, error) {
	panic(req)
	return &authServer.AuthRegistrationResponse{
		Token: "dasdas",
	}, nil
}

func (a *AuthServerApi) AuthRegistration(
	ctx context.Context, req *authServer.AuthRegistrationRequest) (*authServer.AuthRegistrationResponse, error) {
	panic(req)
	return &authServer.AuthRegistrationResponse{Token: "Adssad"}, nil
}

func (a *AuthServerApi) AuthRoles(
	ctx context.Context, req *authServer.AuthRolesRequest) (*authServer.AuthRolesResponse, error) {
	panic(req)
	return &authServer.AuthRolesResponse{RolesFlag: false}, nil
}
