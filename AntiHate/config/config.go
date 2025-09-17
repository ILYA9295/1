package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string
	AppPort          string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	return &Config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		PostgresUser:     getEnv("POSTGRES_USER", "bot_user"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "bot_pass"),
		PostgresDB:       getEnv("POSTGRES_DB", "bot_db"),
		PostgresHost:     getEnv("POSTGRES_HOST", "db"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		AppPort:          getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
