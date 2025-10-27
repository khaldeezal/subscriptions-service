package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port:        getenv("APP_PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/subscriptions?sslmode=disable"),
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
