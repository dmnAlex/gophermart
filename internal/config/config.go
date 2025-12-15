package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type Config struct {
	LogLevel             string `env:"LOG_LEVEL"`
	RunAddress           string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	MigrationsPath       string `env:"MIGRATIONS_PATH"`
	JWTSecret            string `env:"JWT_SECRET"`
}

func New() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "service run address and port")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "localhost:8081", "accrual run address and port")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "database uri")
	flag.StringVar(&cfg.MigrationsPath, "m", "./migrations", "migrations path")
	flag.StringVar(&cfg.JWTSecret, "j", "defaultsecret", "JWT secret")

	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "parse env")
	}

	return cfg, nil
}
