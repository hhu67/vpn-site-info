package config

import (
	"log"
	"os"
)

type Config struct {
	PSQL      string
	JWTSecret string
	Port      string
}

func Load() *Config {
	psql := os.Getenv("PSQL")
	jwtSecret := getEnvOrDefault("JWT_SECRET", "default-secret-change-me")
	port := getEnvOrDefault("PORT", "8080")

	log.Printf("[CONFIG] PSQL: %s", psql)
	log.Printf("[CONFIG] JWT_SECRET: %d chars", len(jwtSecret))
	log.Printf("[CONFIG] PORT: %s", port)

	return &Config{
		PSQL:      psql,
		JWTSecret: jwtSecret,
		Port:      port,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		log.Printf("[CONFIG] %s loaded from environment", key)
		return value
	}
	log.Printf("[CONFIG] %s using default value", key)
	return defaultValue
}
