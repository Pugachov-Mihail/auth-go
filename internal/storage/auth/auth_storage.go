package auth_storage

import (
	configapp "auth/internal/config"
	"auth/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"time"
)

var (
	ErrorUserNotFound = errors.New("пользователь не найден")
	ErrorUserExists   = errors.New("пользователь существует")
)

type Storage struct {
	db *sql.DB
}

func New(storagePath configapp.ConfigDB) (*Storage, error) {
	connDb := "postgres://" + storagePath.UserDb + ":" + storagePath.PassDb + "@" + storagePath.Host +
		":" + storagePath.PortDb + "/" + storagePath.DbName + "?sslmode=disable"

	db, err := sql.Open("postgres", connDb)

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка базы ", err)
		return nil, fmt.Errorf("%s: %v", "postgre", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	log *slog.Logger,
	email string,
	passHash []byte,
	login string,
	steamId int64) (int64, error) {

	exists, err := s.UserExists(ctx, email)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("пользователь существует")
	}

	query := `INSERT INTO users_my (email, pass_hash, steam_id, login) VALUES ($1,$2,$3,$4) RETURNING id;`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var pk int64

	err = s.db.QueryRowContext(dbCtx, query, email, passHash, steamId, login).Scan(&pk)

	if err != nil {
		return 0, fmt.Errorf("save user: %w", err)
	}

	return pk, nil
}

func (s *Storage) User(ctx context.Context, login string) (models.User, error) {
	query := `SELECT id, email, pass_hash FROM users_my WHERE login = $1;`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var user models.User

	err := s.db.QueryRowContext(dbCtx, query, login).Scan(&user.Id, &user.Email, &user.PassHash)

	if err != nil {
		return models.User{}, fmt.Errorf("Gets user: %w", err)
	}

	return user, nil
}

func (s *Storage) RolesUser(ctx context.Context, userId int64) (models.Roles, error) {
	stmt, err := s.db.Prepare("SELECT roles_name, roles_flag FROM roles where user_id = ?;")
	if err != nil {
		return models.Roles{}, fmt.Errorf("Gets user: %w", err)
	}
	row := stmt.QueryRowContext(ctx, userId)

	var roles models.Roles

	err = row.Scan(&roles.RolesName, &roles.RolesFlag)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Roles{}, fmt.Errorf("Roles: %w", err)
		}
		return models.Roles{}, fmt.Errorf("roles: %w", err)
	}
	return roles, nil
}

func (s *Storage) UserExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT email FROM users_my WHERE email = $1;`

	var user models.User

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := s.db.QueryRowContext(dbCtx, query, email).Scan(&user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists user error: %w", err)
	}

	if email == user.Email {
		return true, ErrorUserExists
	}
	return false, nil
}

func (s *Storage) PermissionAccess(ctx context.Context, token string) (models.User, error) {
	query := `SELECT id, email FROM users_my WHERE id=$1;`
	userId := `SELECT user_id FROM access_token WHERE token=$1;`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var user models.User

	if err := s.db.QueryRowContext(dbCtx, userId, token).Scan(&user.Id); err != nil {
		return models.User{}, fmt.Errorf("ошибка получения пользователя")
	}

	if err := s.db.QueryRowContext(dbCtx, query, user.Id).Scan(&user.Id, &user.Email); err != nil {
		return models.User{}, fmt.Errorf("ошибка получения пользователя")
	}

	return user, nil
}

func (s *Storage) SaveToken(ctx context.Context, token string, id int64) error {
	query := `INSERT INTO access_token(user_id, token) VALUES ($1, $2);`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.db.QueryRowContext(dbCtx, query, id, token)
	if err != nil {
		return fmt.Errorf("ошибка сохранения токена")
	}
	return nil
}

func (s *Storage) RefreshToken(ctx context.Context, tokenNew string, tokenOld string) error {
	query := `UPDATE access_token SET token=$1 WHERE token = $2;`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := s.db.QueryRowContext(dbCtx, query, tokenNew, tokenOld)
	if err != nil {
		return fmt.Errorf("ошибка обновления токена")
	}
	return nil

}
