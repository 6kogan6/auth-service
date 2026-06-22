package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL      string
	JWTSecret  string
	ServerAddr string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DBURL:      os.Getenv("DB_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerAddr: os.Getenv("SERVER_ADDR"),
	}

	if cfg.DBURL == "" {
		return Config{}, errors.New("DB_URL не задан")
	}

	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET не задан")
	}

	if cfg.ServerAddr == "" {
		return Config{}, errors.New("SERVER_ADDR не задан")
	}

	return cfg, nil
}
