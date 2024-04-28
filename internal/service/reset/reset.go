package reset_service

import (
	"context"
	"fmt"
	"log/slog"
)

type Reset struct {
	log          *slog.Logger
	resetStorage ResetStorage
}

type ResetStorage interface {
	Email(ctx context.Context, email string, uid int64) (int64, error)
	Password(ctx context.Context, password string, userId int64) (int64, error)
	IdSteam(ctx context.Context, steamId int64, userId int64) (int64, error)
}

func New(
	log *slog.Logger,
	resetStorage ResetStorage,
) *Reset {
	return &Reset{
		log:          log,
		resetStorage: resetStorage,
	}
}

func (r *Reset) ResetEmailStore(ctx context.Context, email string, userId int64) (int64, error) {
	idUser, err := r.resetStorage.Email(ctx, email, userId)
	if err != nil {
		r.log.Error("ошибка изменения e-mail: ", err)
		return 0, fmt.Errorf("ошибка изменения e-mail: %w", err)
	}
	return idUser, nil
}

func (r *Reset) ResetPassword(ctx context.Context, password string, userId int64) (int64, error) {
	idUser, err := r.resetStorage.Password(ctx, password, userId)

	if err != nil {
		return 0, fmt.Errorf("ошибка изменения пароля: %w", err)
	}

	return idUser, nil
}

func (r *Reset) ResetIdSteam(ctx context.Context, steamId int64, userId int64) (int64, error) {
	_, err := r.resetStorage.IdSteam(ctx, steamId, userId)
	if err != nil {
		return 0, fmt.Errorf("ошибка изменения steam id: %w ", err)
	}
	//TODO навесить кафку что бы все мс связанные со стим айди были вкурсе

	return 0, err
}
