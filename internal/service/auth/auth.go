package auth

import (
	configapp "auth/internal/config"
	"auth/internal/domain/models"
	kafka_user "auth/internal/kafka"
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
	Permission = "Permission"
	ErrInvalid = "invalid credentials"
)

type Auth struct {
	Log             *slog.Logger
	TokenTTL        time.Duration
	UsrProvider     UserProvider
	UsrSaver        UserSaver
	TokenSaver      TokenSaver
	Cfg             configapp.Config
	registerNewUser RegisterNewUser
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
	Logout(ctx context.Context, token string) error
	RolesUser(ctx context.Context, uid int64) (models.Roles, error)
	PermissionAccess(ctx context.Context, token string) (models.User, error)
}

type TokenSaver interface {
	SaveToken(ctx context.Context, token string, id int64) (int64, error)
	RefreshToken(ctx context.Context, tokenNew string, tokenOld string) error
}

type RegisterNewUser interface {
	UserRegisterKafka(logs *slog.Logger, userId int64, steamId int64) error
}

// New конструктор сервисного слоя Auth
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	TokenSaver TokenSaver,
	tokenTTl time.Duration,
	cfg *configapp.Config,
	kf *kafka_user.Conf,
) *Auth {
	return &Auth{
		UsrSaver:        userSaver,
		UsrProvider:     userProvider,
		TokenSaver:      TokenSaver,
		TokenTTL:        tokenTTl,
		Log:             log,
		Cfg:             *cfg,
		registerNewUser: kf,
	}
}

func (a *Auth) LoginUser(ctx context.Context, login string, password string) (string, error) {

	log := a.Log.With(
		slog.String("Auth ", Login),
		slog.String("login", login))

	log.Info("Попытка пользователю: " + login + " залогиниться")

	user, err := a.UsrProvider.User(ctx, login)

	if err != nil {
		if errors.Is(err, auth_storage.ErrorUserNotFound) {
			log.Warn("Пользователь не найден;", err)

			return "", fmt.Errorf("%s: %s", Login, ErrInvalid)
		}
		log.Error("Ошибка получения пользователя;", err)
		return "", fmt.Errorf("%s: %w", Login, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn(ErrInvalid, err)
		return "", fmt.Errorf("%s: %s", Login, ErrInvalid)
	}

	token, err := jwt.NewToken(user, a.Cfg.Secret, a.TokenTTL)
	if err != nil {
		log.Warn("Ошибка генерации токена;", err)
		return "", fmt.Errorf("%s: %w", Login, err)
	}

	_, err = a.TokenSaver.SaveToken(ctx, token, user.Id)
	if err != nil {
		log.Warn("ошибка сохранения токена ", err)
		return "", fmt.Errorf("ошибка сохранения токена %w", err)
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

	id, err := a.UsrSaver.SaveUser(ctx, log, email, passHash, login, steamId)
	if err != nil {
		if errors.Is(err, auth_storage.ErrorUserExists) {
			log.Warn("Пользователь существует: ", err)
			return 0, fmt.Errorf("%s: %w", Register, err)
		}
		log.Error("Ошибка сохранения пользователя: ", err)
		return 0, fmt.Errorf("%s: %w", Register, err)
	}

	log.Info("Пользователь ", login, " зарегестрировался")

	if err := a.registerNewUser.UserRegisterKafka(a.Log, id, steamId); err != nil {
		log.Error("Ошибка передачи информации о регистрации пользователя" + strconv.FormatInt(id, 10) + login)
		return 0, fmt.Errorf("ошибка кафки: %w", err)
	}

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

// AccessPermission обновление токена
func (a *Auth) AccessPermission(ctx context.Context, token string) (string, error) {
	//TODO написать тесты и поправить ручку
	log := a.Log.With(slog.String("Auth", Permission))

	user, err := a.UsrProvider.PermissionAccess(ctx, token)
	if err != nil {
		log.Warn("не обработанный токен")
		return "", fmt.Errorf("ошибка обработки токена: %w", err)
	}

	if jwt.ValidateToken(token, a.Cfg) {
		log.Debug("Токен не сменен")
		return token, nil
	}

	tokenNew, err := jwt.NewToken(user, a.Cfg.Secret, a.TokenTTL)

	if err = a.TokenSaver.RefreshToken(ctx, tokenNew, token); err != nil {
		log.Warn("Ошибка обновления токена")
		return "", fmt.Errorf("ошибка обновления токена: %w", err)
	}

	log.Debug("Токен сменен")
	return tokenNew, nil
}

func (a *Auth) LogoutUser(ctx context.Context, token string) (bool, error) {
	log := a.Log.With(slog.String("Auth", "Logout"))

	if err := a.UsrProvider.Logout(ctx, token); err != nil {
		log.Warn("Ошибка выхода из профиля: ", err)
		return false, fmt.Errorf("ошибка выхода из профиля: %w", err)
	}

	return true, nil
}
