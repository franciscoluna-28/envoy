package config

import (
	"os"
)

type Config struct {
	port string
	dbURL string
}

func New() *Config {
	return &Config{
		port: getEnv("PORT", "8080"),
		dbURL: getEnv("DB_URL", "app.db"),
	}
}

func (c *Config) Port() string {
	return c.port
}

func (c *Config) DatabaseURL() string {
	return c.dbURL
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
