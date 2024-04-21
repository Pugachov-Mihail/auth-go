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
)

var (
	ErrorUserNotFound = errors.New("пользователь не найден")
	ErrorUserExists   = errors.New("пользователь существует")
)

type Storage struct {
	db *sql.DB
}

func New(storagePath configapp.ConfigDB) (*Storage, error) {
	connDb := storagePath.UserDb + ":" + storagePath.PassDb + "@" + storagePath.Host + ":" + storagePath.PortDb + "/" + storagePath.DbName

	db, err := sql.Open("postgres", connDb)

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Ошибка закрытия базы ", err)
		}
	}()

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка базы ", err)
		return nil, fmt.Errorf("%s: %w", "postgre", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	log *slog.Logger,
	email string,
	passHash []byte,
	login string,
	steamId int64) (uid int64, err error) {
	stmt, err := s.db.Prepare("INSERT INTO user(email, pass_hash, steam_id, login) VALUES(?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("Save user: %w", err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash, steamId, login)
	if err != nil {
		return 0, fmt.Errorf("Save user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error("Ошибка сохранения пользователя: ", email, err)
		return 0, fmt.Errorf("Save user: %w", err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, login string) (models.User, error) {
	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM user WHERE id = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("Gets user: %w", err)
	}
	row := stmt.QueryRowContext(ctx, login)

	var user models.User

	err = row.Scan(&user.Id, &user.PassHash, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("User: %w", err)
		}
		return models.User{}, fmt.Errorf("User: %w", err)
	}

	return user, nil
}

func (s *Storage) RolesUser(ctx context.Context, userId int64) (models.Roles, error) {
	stmt, err := s.db.Prepare("SELECT roles_name, roles_flag FROM roles where user_id = ?")
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
		return models.Roles{}, fmt.Errorf("Roles: %w", err)
	}
	return roles, nil
}
