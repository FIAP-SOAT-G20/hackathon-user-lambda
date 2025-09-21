package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string

	// DynamoDB
	AWSRegion      string
	UsersTableName string
	IdsTableName   string

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() *Config {
	// Load .env if present (dev convenience)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: .env not found or failed to load: %v", err)
	}

	jwtExpStr := getEnv("JWT_EXPIRATION", "24h")
	exp, err := time.ParseDuration(jwtExpStr)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRATION %q, defaulting to 24h", jwtExpStr)
		exp = 24 * time.Hour
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok || jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable must be set and non-empty")
	}

	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
		UsersTableName: getEnv("USERS_TABLE_NAME", "hackathon_users"),
		IdsTableName:   getEnv("IDS_TABLE_NAME", "hackathon_ids"),
		JWTSecret:      jwtSecret,
		JWTExpiration:  exp,
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
