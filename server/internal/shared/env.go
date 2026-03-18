package shared

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port          string `env:"PORT" envDefault:"8080"`
	DatabaseURL   string `env:"DATABASE_URL,required"`
	JWTSecret     string `env:"JWT_SECRET,required"`
	Env           string `env:"APP_ENV" envDefault:"development"`
	EncryptionKey string `env:"ENCRYPTION_KEY,required"`
	ChecksumKey   string `env:"CHECKSUM_KEY,required"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Missing environment variables: %v", err)
	}

	return cfg
}
