package kafka

import (
	"context"
	"log/slog"
)

type Kafka struct {
	values string
}

// UserRegisterKafka данный метод будет отправлять прокси сервисам о том, что новый пользователь зарегался
// Также будет отправка данных для выгрузки информации о матчах пользователя
func (k *Kafka) UserRegisterKafka(ctx context.Context, log *slog.Logger, userId int64, steamId int64) (bool, error) {
	log.Info("Отправка данных о пользователе")
	return false, nil
}
