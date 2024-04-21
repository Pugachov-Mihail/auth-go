package auth

import (
	"auth/internal/domain/models"
	"auth/internal/service/lib/jwt"
	auth_storage "auth/internal/storage/auth"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const (
	Register   = "Register"
	Login      = "Login"
	Roles      = "Roles"
	ErrInvalid = "invalid credentials"
)

type Auth struct {
	log         *slog.Logger
	tokenTTL    time.Duration
	usrProvider UserProvider
	usrSaver    UserSaver
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		password []byte,
		login string,
		steamId int64) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (models.User, error)
	RolesUser(ctx context.Context, uid int64) (models.Roles, error)
}

// New конструктор сервисного слоя Auth
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTl time.Duration) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		tokenTTL:    tokenTTl,
		log:         log,
	}
}

func (a *Auth) LoginUser(ctx context.Context, login string, password string, secret string) (string, error) {

	log := a.log.With(
		slog.String("Auth ", Login),
		slog.String("login", login))
	log.Info("Register user", login)

	user, err := a.usrProvider.User(ctx, login)

	if err != nil {
		if errors.Is(err, auth_storage.ErrorUserNotFound) {
			a.log.Warn("Пользователь не найден", err)

			return "", fmt.Errorf("%s: %w", Login, ErrInvalid)
		}
		a.log.Error("Ошибка получения пользователя", err)
		return "", fmt.Errorf("%s: %w", Login, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn(ErrInvalid, err)
		return "", fmt.Errorf("%s: %w", Login, ErrInvalid)
	}
	token, err := jwt.NewToken(user, secret, a.tokenTTL)
	if err != nil {
		a.log.With("Ошибка генерации токена", err)
		return "", fmt.Errorf("%s: %w", Login, err)
	}
	return token, nil
}

func (a *Auth) RegisterUser(
	ctx context.Context,
	login string,
	password string,
	email string,
	steamId int64) (int64, error) {

	log := a.log.With(
		slog.String("Auth ", Register),
		slog.String("login", login))

	log.Info("Registering user", login)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Ошибка генерации хеша пароля", err)
		return 0, fmt.Errorf("%s: %w", Register, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, login, passHash, email, steamId)
	if err != nil {
		if errors.Is(err, auth_storage.ErrorUserExists) {
			log.Warn("Пользователь существует", err)
			return 0, fmt.Errorf("%s: %w", Register, err)
		}
		log.Error("Ошибка сохранения пользователя", err)
		return 0, fmt.Errorf("%s: %w", Register, err)
	}

	log.Info("Пользователь ", login, " зарегистрировался")

	return id, nil
}

func (a *Auth) RolesUser(ctx context.Context, uid int64) (models.Roles, error) {
	log := a.log.With(
		slog.String("Auth ", Roles),
		slog.Int64("", uid))
	log.Info("Find roles user", uid)

	roles, err := a.usrProvider.RolesUser(ctx, uid)

	if err != nil {
		return models.Roles{}, fmt.Errorf("%s: %w", Roles, err)
	}
	log.Info("Проверка ролей пользователя", uid)

	return roles, err
}
