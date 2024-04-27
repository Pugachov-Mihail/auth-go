package auth_storage

import (
	"context"
	"database/sql"
)

type StorageReset struct {
	db *sql.DB
}

func (s *StorageReset) Email(ctx context.Context, email string, uid int64) (int64, error) {
	panic(ctx)
}
func (s *StorageReset) Password(ctx context.Context, password string, userId int64) (int64, error) {
	panic(ctx)
}
func (s *StorageReset) IdSteam(ctx context.Context, steamId int64, userId int64) (int64, error) {
	panic(ctx)
}
