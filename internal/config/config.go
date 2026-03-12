package config

import "os"

type Config struct {
	Server struct {
		Host string
		Port string
	}
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

	return cfg
}
