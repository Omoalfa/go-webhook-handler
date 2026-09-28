package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Environment   string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	WebhookSecret string
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func LoadConfig() (*Config, error) {
	loadEnv()

	config := Config{}

	port, exist := os.LookupEnv("PORT")
	if !exist {
		port = "8080"
	}
	config.Port = port

	env, exist := os.LookupEnv("ENVIRONMENT")
	if !exist {
		env = "development"
	}
	config.Environment = env

	databaseURL, exist := os.LookupEnv("DATABASE_URL")
	if !exist || databaseURL == "" {
		return &config, errors.New("DATABASE_URL env is required")
	}
	config.DatabaseURL = databaseURL

	redisURL, exist := os.LookupEnv("REDIS_URL")
	if !exist || redisURL == "" {
		return &config, errors.New("REDIS_URL env is required")
	}
	config.RedisURL = redisURL

	jwtSecret, exist := os.LookupEnv("JWT_SECRET")
	if !exist || jwtSecret == "" {
		return &config, errors.New("JWT_SECRET is required")
	}
	config.JWTSecret = jwtSecret

	webhookSecret, exist := os.LookupEnv("WEBHOOK_SECRET")
	if !exist || webhookSecret == "" {
		return &config, errors.New("WEBHOOK_SECRET is required")
	}
	config.WebhookSecret = webhookSecret

	return &config, nil
}
