package suite

import (
	configapp "auth/internal/config"
	auth "auth/protos/gen/dota_traker.auth.v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

const Host = "localhost"

type Suite struct {
	*testing.T
	Cfg        *configapp.Config
	AuthClient auth.AuthServerClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	//TODO переделать файдл конфигов для тестов
	cfg := configapp.MustLoadByPath("../config/dev.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.TimeOut)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(),
		getAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc connect failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: auth.NewAuthServerClient(cc),
	}
}

func getAddress(cfg *configapp.Config) string {
	return net.JoinHostPort(Host, "8002")
}
