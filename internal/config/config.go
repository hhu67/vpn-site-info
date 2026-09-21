package config

import "os"

type Config struct {
	PSQL      string
	JWTSecret string
	Port      string
}

func Load() *Config {
	return &Config{
		PSQL:      os.Getenv("PSQL"),
		JWTSecret: getEnvOrDefault("JWT_SECRET", "default-secret-change-me"),
		Port:      getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
