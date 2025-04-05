package config

import (
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseURL string        `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@db:5432/postgres"`
	TTL         time.Duration `env:"TTL" env-default:"5m"`
	Level       string        `env:"Level" env-default:"debug"`
	PortMetrics string        `env:"PortMetrics" env-default:"9090"`
	LogLevel    slog.Level
}

var AppConfig Config

func Init() error {
	if err := cleanenv.ReadEnv(&AppConfig); err != nil {
		slog.Error("Error", "error", err)
		return err
	}
	return nil
}
