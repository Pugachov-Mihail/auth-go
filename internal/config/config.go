package configapp

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

type Config struct {
	Env         string        `yaml:"env" env-required:"dev"`
	StoragePath ConfigDB      `yaml:"storage_path" env-required:"true"`
	GRPC        GRPCConf      `yaml:"grpc" env-required:"true"`
	Secret      string        `yaml:"secret"`
	TokenTTL    time.Duration `yaml:"tokenTTL"`
}

type GRPCConf struct {
	Port    int           `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}

type ConfigDB struct {
	Host   string `yaml:"POSTGRES_HOST"`
	UserDb string `yaml:"POSTGRES_USER"`
	DbName string `yaml:"POSTGRES_DB"`
	PassDb string `yaml:"POSTGRES_PASSWORD"`
	PortDb string `yaml:"POSTGRES_PORT"`
}

func MustLoad() *Config {
	fetchPath := fetchConfigPath()
	if fetchPath == "" {
		panic("Пустой файл конфигурации")
	}

	return MustLoadByPath(fetchPath)
}

func MustLoadByPath(fetchPath string) *Config {
	if _, err := os.Stat(fetchPath); os.IsNotExist(err) {
		panic("Не существует файл конфигурации: " + fetchPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(fetchPath, &cfg); err != nil {
		panic("Ошибка чтения конфига" + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func SetupLoger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
