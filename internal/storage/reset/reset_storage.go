package reset

import (
	configapp "auth/internal/config"
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
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
	query := `UPDATE users_my SET email=$1 WHERE id=$2;`
	fu := `SELECT id FROM users_my WHERE email=$1;`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var user uint64

	err := s.db.QueryRowContext(dbCtx, query, email, uid)
	if err != nil {
		return 0, fmt.Errorf("ошибка смены почты: %w ", err)
	}
	if err := s.db.QueryRowContext(dbCtx, fu, email).Scan(&user); err != nil {
		return 0, fmt.Errorf("ошибка получения пользователя по новой почте: %w", err)
	}

	return int64(user), nil
}
func (s *StorageReset) Password(ctx context.Context, password string, userId int64) (int64, error) {
	query := `UPDATE users_my SET password=$1 WHERE id=$2;`
	fu := `SELECT id FROM users_my WHERE id=$1;`

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("ошибка генерации хеша пароля: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var user uint64

	if err := s.db.QueryRowContext(dbCtx, query, passHash, userId); err != nil {
		return 0, fmt.Errorf("ошибка смены пароля: %v", err)
	}

	if err := s.db.QueryRowContext(dbCtx, fu, userId).Scan(&user); err != nil {
		return 0, fmt.Errorf("ошибка получения пользователя при смене пароля: %w", err)
	}

	return int64(user), nil
}
func (s *StorageReset) IdSteam(ctx context.Context, steamId int64, userId int64) (int64, error) {
	query := `UPDATE users_my SET steam_id=$1 WHERE id=$2;`
	fu := `SELECT id FROM users_my WHERE id=$1;`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var user uint64

	err := s.db.QueryRowContext(dbCtx, query, steamId, userId)
	if err != nil {
		return 0, fmt.Errorf("ошибка смены Steam Id: %v ", err)
	}
	if err := s.db.QueryRowContext(dbCtx, fu, userId).Scan(&user); err != nil {
		return 0, fmt.Errorf("ошибка получения Steam Id: %w", err)
	}

	return int64(user), nil
}
