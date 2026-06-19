package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	RedisURL  string
	JWTSecret string
	APIURL    string
}

func Load() *Config {
	for _, path := range []string{"../.env", ".env", "../../.env"} {
		if err := godotenv.Load(path); err == nil {
			log.Printf("[config] .env cargado desde %s", path)
			break
		}
	}

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		RedisURL:  getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret: getEnv("JWT_SECRET", "cambiar_en_produccion"),
		APIURL:    getEnv("API_URL", "http://localhost:8000"),
	}

	log.Printf("[config] Puerto: %s | Redis: %s | API: %s",
		cfg.Port, cfg.RedisURL, cfg.APIURL)

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
