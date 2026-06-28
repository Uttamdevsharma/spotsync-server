package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration time.Duration
	BcryptCost    int
}

func LoadConfig() *Config {
	godotenv.Load()

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DB_URL", getEnv("DATABASE_URL", ""))
	jwtSecret := getEnv("JWT_SECRET", "spotsync-jwt-secret-key-change-me")

	jwtExpStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExp, err := time.ParseDuration(jwtExpStr)
	if err != nil {
		jwtExp = 24 * time.Hour
	}

	costStr := getEnv("BCRYPT_COST", "10")
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		cost = 10
	}

	return &Config{
		ServerPort:    port,
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		JWTExpiration: jwtExp,
		BcryptCost:    cost,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
