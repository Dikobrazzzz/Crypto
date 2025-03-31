package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseURL string        `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@db:5432/postgres"`
	TTL         time.Duration `env:"TTL" env-default:"5m"`
	Level       string        `env:"Level" env-default:"debug"`
	PortMetrics string        `env:"PortMetrics" env-default:"9090"`
	Rlevel      slog.Level
}

var AppConfig Config

func Init() {
	if err := cleanenv.ReadEnv(&AppConfig); err != nil {
		slog.Error("Error", "error", err)
		os.Exit(1)
	}

	switch AppConfig.Level {
	case "debug":
		AppConfig.Rlevel = slog.LevelDebug
	case "info":
		AppConfig.Rlevel = slog.LevelInfo
	case "error":
		AppConfig.Rlevel = slog.LevelError
	default:
		slog.Error("Invalid log error")
		os.Exit(1)
	}
}
