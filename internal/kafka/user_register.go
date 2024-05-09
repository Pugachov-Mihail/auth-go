package kafka

import (
	configapp "auth/internal/config"
	"context"
	"log/slog"
	"strconv"
)

type Conf struct {
	conf *configapp.KafkaConfig
}

type RegisterUser struct {
	Id      string `json:"id"`
	SteamID int64  `json:"steam_id"`
}

// UserRegisterKafka данный метод будет отправлять прокси сервисам о том, что новый пользователь зарегался
// Также будет отправка данных для выгрузки информации о матчах пользователя
func (k *Conf) UserRegisterKafka(ctx context.Context, logs *slog.Logger, userId int64, steamId int64) (bool, error) {
	log := logs.With(slog.String("Кафка", strconv.FormatInt(steamId, 10)))

	userInfo := struct {
		userID  string
		steamId string
	}{strconv.FormatInt(userId, 10), strconv.FormatInt(steamId, 10)}

	log.Info("Отправка данных о пользователе")
	log.Debug("", userInfo)
	log.Debug("", k.conf.Broker)

	//conn, err := kafka.DialLeader(context.Background(), "tcp", k.conf.Broker, "create-user", 1)
	//if err != nil {
	//	log.Error("Ошибка соединения с кафкой: ", err)
	//	return false, fmt.Errorf("ошибка кафки %w", err)
	//}
	//
	//if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
	//	log.Error("Ошибка соединения с кафкой: ", err)
	//	return false, fmt.Errorf("ошибка записи кафки %w", err)
	//}
	//_, err = conn.WriteMessages(kafka.Message{Value: []byte(userInfo.steamId)})
	//if err != nil {
	//	log.Error("ошибка записи данных", err)
	//	return false, fmt.Errorf("ошибка записи данных %w", err)
	//}
	//if err := conn.Close(); err != nil {
	//	log.Error("ошибка закрытия соединения", err)
	//	return false, fmt.Errorf("ошибка закрытия соединения %w", err)
	//}
	return true, nil
}
