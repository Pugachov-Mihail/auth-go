package tests

import (
	"auth/internal/domain/models"
	"auth/internal/service/auth"
	"auth/internal/service/auth/mocks"
	"context"
	"log/slog"
	"reflect"
	"testing"
	"time"
)

func TestAuth_LoginUser(t *testing.T) {
	var ctx context.Context

	type fields struct {
		log         *slog.Logger
		tokenTTL    time.Duration
		usrProvider auth.UserProvider
		usrSaver    auth.UserSaver
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
		secret   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test Mock",
			args: args{
				ctx:      ctx,
				login:    "Adams",
				password: "12345",
				secret:   "dasdas",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userProvider := mocks.NewUserProvider(t)
			userSaver := mocks.NewUserSaver(t)

			userProvider.
				On("User", tt.args.ctx, tt.args.login).
				Return(nil, nil)

			a := &auth.Auth{
				Log:         tt.fields.log,
				TokenTTL:    tt.fields.tokenTTL,
				usrProvider: userProvider,
				usrSaver:    userSaver,
			}
			got, err := a.LoginUser(tt.args.ctx, tt.args.login, tt.args.password, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoginUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuth_RegisterUser(t *testing.T) {
	type fields struct {
		log         *slog.Logger
		tokenTTL    time.Duration
		usrProvider auth.UserProvider
		usrSaver    auth.UserSaver
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
		email    string
		steamId  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &auth.Auth{
				Log:         tt.fields.log,
				TokenTTL:    tt.fields.tokenTTL,
				usrProvider: tt.fields.usrProvider,
				usrSaver:    tt.fields.usrSaver,
			}
			got, err := a.RegisterUser(tt.args.ctx, tt.args.login, tt.args.password, tt.args.email, tt.args.steamId)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuth_RolesUser(t *testing.T) {
	type fields struct {
		log         *slog.Logger
		tokenTTL    time.Duration
		usrProvider auth.UserProvider
		usrSaver    auth.UserSaver
	}
	type args struct {
		ctx context.Context
		uid int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Roles
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &auth.Auth{
				Log:         tt.fields.log,
				TokenTTL:    tt.fields.tokenTTL,
				usrProvider: tt.fields.usrProvider,
				usrSaver:    tt.fields.usrSaver,
			}
			got, err := a.RolesUser(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("RolesUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RolesUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		log          *slog.Logger
		userSaver    auth.UserSaver
		userProvider auth.UserProvider
		tokenTTl     time.Duration
	}
	tests := []struct {
		name string
		args args
		want *auth.Auth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := auth.New(tt.args.log, tt.args.userSaver, tt.args.userProvider, tt.args.tokenTTl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
