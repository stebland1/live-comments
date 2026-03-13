package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Postgres struct {
		Host     string
		Port     string
		Username string
		Password string
		DB       string
		Timeout  time.Duration
	}
}

func (cfg Config) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
		cfg.Postgres.Timeout,
	)
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("%s env var must be set", key))
	}

	return v
}

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func Load() Config {
	cfg := Config{}

	cfg.Server.Host = getEnv("LCOM_SERVER_HOST", "localhost")
	cfg.Server.Port = getEnv("LCOM_SERVER_PORT", "8081")

	cfg.Postgres.Username = mustEnv("LCOM_PG_APP_USER")
	cfg.Postgres.Password = mustEnv("LCOM_PG_APP_PASSWORD")
	cfg.Postgres.Host = mustEnv("LCOM_PG_HOST")
	cfg.Postgres.Port = mustEnv("LCOM_PG_PORT")
	cfg.Postgres.DB = mustEnv("LCOM_PG_DB")
	cfg.Postgres.Timeout = 5 * time.Second

	return cfg
}
