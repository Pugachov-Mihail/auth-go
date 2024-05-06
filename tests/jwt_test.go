package tests

import (
	configapp "auth/internal/config"
	"auth/internal/domain/models"
	"auth/internal/service/lib/jwt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	type args struct {
		user     models.User
		secret   string
		duration time.Duration
		st       configapp.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1",
			args: args{
				user: models.User{
					Email:    gofakeit.Email(),
					Id:       gofakeit.Int64(),
					PassHash: []byte(gofakeit.Password(true, true, true, true, false, 10)),
					SteamId:  gofakeit.Int64(),
				},
				secret:   "app",
				duration: time.Microsecond * 1,
				st: configapp.Config{
					Secret: "app",
				},
			},
		},
		{
			name: "Case 2",
			args: args{
				user: models.User{
					Email:    gofakeit.Email(),
					Id:       gofakeit.Int64(),
					PassHash: []byte(gofakeit.Password(true, true, true, true, false, 10)),
					SteamId:  gofakeit.Int64(),
				},
				secret:   "app",
				duration: time.Second * 60,
				st: configapp.Config{
					Secret: "app",
				},
			},
		},
		{
			name: "Case 3",
			args: args{
				user: models.User{
					Email:    gofakeit.Email(),
					Id:       gofakeit.Int64(),
					PassHash: []byte(gofakeit.Password(true, true, true, true, false, 10)),
					SteamId:  gofakeit.Int64(),
				},
				secret:   "app",
				duration: time.Second * 60,
				st: configapp.Config{
					Secret: "app",
				},
			},
		},
		{
			name: "Case 4",
			args: args{
				user: models.User{
					Email:    gofakeit.Email(),
					Id:       gofakeit.Int64(),
					PassHash: []byte(gofakeit.Password(true, true, true, true, false, 10)),
					SteamId:  gofakeit.Int64(),
				},
				secret:   "app",
				duration: time.Second * 60,
				st: configapp.Config{
					Secret: "app",
				},
			},
		},
		{
			name: "Case 5",
			args: args{
				user: models.User{
					Email:    gofakeit.Email(),
					Id:       gofakeit.Int64(),
					PassHash: []byte(gofakeit.Password(true, true, true, true, false, 10)),
					SteamId:  gofakeit.Int64(),
				},
				secret:   "app",
				duration: time.Second * 60,
				st: configapp.Config{
					Secret: "app",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jwt.NewToken(tt.args.user, tt.args.secret, tt.args.duration)
			require.NoError(t, err)
			assert.NotEmpty(t, got)

			token := jwt.ValidateToken(got, tt.args.st)
			require.NoError(t, err)
			assert.False(t, token)
		})
	}
}
