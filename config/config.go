package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@db:5432/postgres"`
}

var AppConfig Config

func Init() {
	if err := cleanenv.ReadEnv(&AppConfig); err != nil {
		slog.Error("Error", "error", err)
	}
}
