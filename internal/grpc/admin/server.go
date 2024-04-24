package admin

import (
	"auth/internal/validator/base_validate"
	adminServer "auth/protos/gen/dota_traker.admin.v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
)

type serverAdminApi struct {
	adminServer.UnimplementedAdminPanelServer
}

func RegisterServerApi(grpc *grpc.Server) {
	adminServer.RegisterAdminPanelServer(grpc, &serverAdminApi{})
}

func (s *serverAdminApi) AdminPermission(
	ctx context.Context,
	req *adminServer.AdminPermissionRequest) (*adminServer.AdminPermissionResponse, error) {
	var log *slog.Logger
	log.With("Admin Permission")

	if !base_validate.ValidatorAdminPermission(req) {
		log.Warn("Пустые данные")
		return nil, fmt.Errorf("Пустые данные")
	}

	return &adminServer.AdminPermissionResponse{Success: false}, nil
}

func (s *serverAdminApi) AdminSettingsPanel(
	ctx context.Context,
	req *adminServer.AdminSettiongsPanelRequest) (*adminServer.AdminSettingsPanelResponse, error) {
	var log *slog.Logger
	log.With("Admin Permission")

	if !base_validate.ValidatorAdminSetting(req) {
		log.Warn("Пустые данные")
		return nil, fmt.Errorf("Пустые данные")
	}
	return &adminServer.AdminSettingsPanelResponse{Service: 12}, nil
}

func (s *serverAdminApi) AdminListInformation(
	ctx context.Context,
	req *adminServer.AdminListInformationRequest) (*adminServer.AdminListInformationsResponse, error) {
	panic("not implemented")
}
