package reset

import (
	configapp "auth/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type StorageReset struct {
	db *sql.DB
}

func New(storagePath configapp.ConfigDB) (*StorageReset, error) {
	connDb := "postgres://" + storagePath.UserDb + ":" + storagePath.PassDb + "@" + storagePath.Host +
		":" + storagePath.PortDb + "/" + storagePath.DbName + "?sslmode=disable"

	db, err := sql.Open("postgres", connDb)

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка базы ", err)
		return nil, fmt.Errorf("%s: %v", "postgre", err)
	}

	return &StorageReset{db: db}, nil
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
