package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	WriteServer struct {
		Host string
		Port string
	}
	StreamServer struct {
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

		MaxOpenConns int
		MaxIdleConns int
	}
	Redis struct {
		Host    string
		Port    string
		Timeout time.Duration
	}
}

func (cfg Config) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
		cfg.Postgres.Timeout/time.Second,
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

	cfg.WriteServer.Host = getEnv("LCOM_WRITE_SERVER_HOST", "localhost")
	cfg.WriteServer.Port = getEnv("LCOM_WRITE_SERVER_PORT", "8081")

	cfg.StreamServer.Host = getEnv("LCOM_STREAM_SERVER_HOST", "localhost")
	cfg.StreamServer.Port = getEnv("LCOM_STREAM_SERVER_PORT", "8082")

	cfg.Postgres.Username = mustEnv("LCOM_PG_APP_USER")
	cfg.Postgres.Password = mustEnv("LCOM_PG_APP_PASSWORD")
	cfg.Postgres.Host = mustEnv("LCOM_PG_HOST")
	cfg.Postgres.Port = mustEnv("LCOM_PG_PORT")
	cfg.Postgres.DB = mustEnv("LCOM_PG_DB")
	cfg.Postgres.Timeout = 5 * time.Second

	cfg.Postgres.MaxOpenConns = 25
	cfg.Postgres.MaxIdleConns = 25

	cfg.Redis.Host = getEnv("LCOM_REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("LCOM_REDIS_PORT", "6379")
	cfg.Redis.Timeout = 5 * time.Second

	return cfg
}
