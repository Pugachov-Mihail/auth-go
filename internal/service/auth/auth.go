package auth

import (
	configapp "auth/internal/config"
	"auth/internal/domain/models"
	"auth/internal/service/lib/jwt"
	auth_storage "auth/internal/storage/auth"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.3 --all

const (
	Register   = "Register"
	Login      = "Login"
	Roles      = "Roles"
	ErrInvalid = "invalid credentials"
)

type Auth struct {
	Log             *slog.Logger
	TokenTTL        time.Duration
	UsrProvider     UserProvider
	UsrSaver        UserSaver
	registerNewUser RegisterNewUser
	Secret          string
	Cfg             configapp.Config
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		log *slog.Logger,
		email string,
		passHash []byte,
		login string,
		steamId int64) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (models.User, error)
	RolesUser(ctx context.Context, uid int64) (models.Roles, error)
}

type RegisterNewUser interface {
	UserRegisterKafka(ctx context.Context, log *slog.Logger, userId int64, steamId int64) (bool, error)
}

// New конструктор сервисного слоя Auth
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTl time.Duration,
	cfg *configapp.Config,
) *Auth {
	return &Auth{
		UsrSaver:    userSaver,
		UsrProvider: userProvider,
		TokenTTL:    tokenTTl,
		Log:         log,
		Cfg:         *cfg,
	}
}

func (a *Auth) LoginUser(ctx context.Context, login string, password string, secret string) (string, error) {

	log := a.Log.With(
		slog.String("Auth ", Login),
		slog.String("login", login))

	log.Info("Пользователь " + login + "залогинился")

	user, err := a.UsrProvider.User(ctx, login)

	if err != nil {
		if errors.Is(err, auth_storage.ErrorUserNotFound) {
			a.Log.Warn("Пользователь не найден;", err)

			return "", fmt.Errorf("%s: %s", Login, ErrInvalid)
		}
		a.Log.Error("Ошибка получения пользователя;", err)
		return "", fmt.Errorf("%s: %w", Login, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.Log.Warn(ErrInvalid, err)
		return "", fmt.Errorf("%s: %s", Login, ErrInvalid)
	}

	token, err := jwt.NewToken(user, a.Cfg.Secret, a.TokenTTL)
	if err != nil {
		a.Log.With("Ошибка генерации токена;", err)
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

	log := a.Log.With(
		slog.String("Auth ", Register),
		slog.String("login", login))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Ошибка генерации хеша пароля;", err)
		return 0, fmt.Errorf("%s: %w", Register, err)
	}

	id, err := a.UsrSaver.SaveUser(ctx, log, login, passHash, email, steamId)
	if err != nil {
		if errors.Is(err, auth_storage.ErrorUserExists) {
			log.Warn("Пользователь существует;", err)
			return 0, fmt.Errorf("%s: %w", Register, err)
		}
		log.Error("Ошибка сохранения пользователя;", err)
		return 0, fmt.Errorf("%s: %w", Register, err)
	}

	log.Info("Пользователь ", login, " зарегистрировался")

	////TODO Сделать вызов кафки для передачи остальным мс о том что пользователь создан
	//_, err = a.registerNewuser.UserRegisterKafka(ctx, log, id, steamId)
	//if err != nil {
	//	log.Error("Ошибка передачи информации о регистрации пользователя" + strconv.FormatInt(id, 10) + login)
	//}

	return id, nil
}

func (a *Auth) RolesUser(ctx context.Context, uid int64) (models.Roles, error) {
	log := a.Log.With(
		slog.String("Auth ", Roles),
		slog.Int64("", uid))
	log.Info("Find roles user " + strconv.FormatInt(uid, 10))

	roles, err := a.UsrProvider.RolesUser(ctx, uid)

	if err != nil {
		return models.Roles{}, fmt.Errorf("%s: %w", Roles, err)
	}
	log.Info("Проверка ролей пользователя " + strconv.FormatInt(uid, 10))

	return roles, err
}
