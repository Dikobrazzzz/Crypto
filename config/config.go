package config

import (
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseURL string        `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@db:5432/postgres"`
	TTL         time.Duration `env:"TTL" env-default:"5m"`
	Level       slog.Leveler  `env:"Level" env-default:"slog.LevelDebug"`
}

var AppConfig Config

func Init() {
	if err := cleanenv.ReadEnv(&AppConfig); err != nil {
		slog.Error("Error", "error", err)
	}
}
