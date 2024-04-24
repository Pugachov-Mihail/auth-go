package reset_service

import (
	"context"
	"log/slog"
)

type Reset struct {
	log          *slog.Logger
	resetStorage ResetStorage
}

type ResetStorage interface {
	Email(ctx context.Context, email string, uid int64) (uint64, error)
}

func (r *Reset) ResetEmail(ctx context.Context, email string, userId int64) (int64, error) {
	panic(email)
}
func (r *Reset) ResetPassword(ctx context.Context, password string, userId int64) (int64, error) {
	panic(password)
}
func (r *Reset) ResetIdSteam(ctx context.Context, steamId int64, userId int64) (int64, error) {
	panic(steamId)
}
