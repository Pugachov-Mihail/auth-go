package kafka_user

import (
	configapp "auth/internal/config"
	"context"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"net"
	"strconv"
)

type Conf struct {
	kf *kafka.Conn
}

type RegisterUser struct {
	Id      string `json:"id"`
	SteamID int64  `json:"steam_id"`
}

func New(cfg *configapp.Config) (*Conf, error) {

	conn, err := kafka.Dial("tcp", cfg.Kafka.Broker)
	if err != nil {
		panic(err.Error())
	}

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}

	var controllerConn *kafka.Conn

	controllerConn, err = kafka.Dial("tcp",
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	return &Conf{kf: conn}, nil
}

// UserRegisterKafka данный метод будет отправлять прокси сервисам о том, что новый пользователь зарегался
// Также будет отправка данных для выгрузки информации о матчах пользователя
func (c *Conf) UserRegisterKafka(logs *slog.Logger, userId int64, steamId int64) error {
	flog := logs.With(slog.String("kafka write", strconv.FormatInt(userId, 10)))

	w := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "create-user",
		Balancer: &kafka.LeastBytes{},
	}

	flog.Info("Коннект с топиком create-user")

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(strconv.FormatInt(userId, 10)),
			Value: []byte(strconv.FormatInt(steamId, 10)),
		},
	)

	if err != nil {
		flog.Warn("ошибка записи в кафку: ", err)
	}

	if err := w.Close(); err != nil {
		flog.Warn("ошибка закрытия соединения записи в кафку: ", err)
	}
	flog.Info("Данные о пользователе отправленны")

	return nil
}
