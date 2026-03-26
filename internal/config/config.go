package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN     string
	JWT_SECRET string
	HTTP_PORT  string
	GRPC_PORT  string
}

func LoadConfig() *Config {
	// Load .env file jika ada
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return &Config{
		DB_DSN:     getEnv("DB_DSN", "root:rootpassword@tcp(127.0.0.1:3306)/warehouse_auth?charset=utf8mb4&parseTime=True&loc=Local"),
		JWT_SECRET: getEnv("JWT_SECRET", "rahasia-super-aman"),
		HTTP_PORT:  getEnv("HTTP_PORT", "8081"),
		GRPC_PORT:  getEnv("GRPC_PORT", "50051"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
