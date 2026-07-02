package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.App.Environment = os.Getenv("APP_ENV")
	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.SSLMode = os.Getenv("DB_SSLMODE")
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.Logger.Level = os.Getenv("LOG_LEVEL")

	return cfg, nil
}
