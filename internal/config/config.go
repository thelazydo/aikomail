package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL  string `envconfig:"DATABASE_URL" default:"postgres://user:password@db:5432/email_verifier?sslmode=disable"`
	Port         string `envconfig:"PORT"         default:"8080"`
	SMTPDomain   string `envconfig:"SMTP_DOMAIN"  default:"localhost"`
	SMTPAddress  string `envconfig:"SMTP_FROM"      default:"verify@example.com"`
	RateLimit    int    `envconfig:"RATE_LIMIT"   default:"5"`
	WorkersCount int    `envconfig:"WORKERS_COUNT"   default:"5"`
}

var CONFIG Config

func LoadConfig() (*Config, error) {
	url := os.Getenv("DATABASE_URL")
	slog.Info("db url existing", "url", url)
	if err := godotenv.Load(); err != nil {
		// ignore if no .env file
		slog.Info("no .env file found, relying on environment variables")
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	CONFIG = cfg

	return &cfg, nil
}
