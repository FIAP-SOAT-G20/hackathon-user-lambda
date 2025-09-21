package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/paramstore"
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

func Load(ctx context.Context) *Config {
	// Load .env if present (dev convenience)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: .env not found or failed to load: %v", err)
	}

	// Try to get JWT settings from Parameter Store first, fallback to env vars
	jwtSecretParam := getEnv("JWT_SECRET_PARAMETER_NAME", "")
	jwtExpirationParam := getEnv("JWT_EXPIRATION_PARAMETER_NAME", "")

	var jwtSecret, jwtExpStr string

	if jwtSecretParam != "" {
		jwtSecret = paramstore.GetParameterWithFallback(ctx, jwtSecretParam, "")
	}
	if jwtSecret == "" {
		// Fallback to environment variable
		var ok bool
		jwtSecret, ok = os.LookupEnv("JWT_SECRET")
		if !ok || jwtSecret == "" {
			log.Fatal("JWT_SECRET must be provided via Parameter Store or environment variable")
		}
	}

	if jwtExpirationParam != "" {
		jwtExpStr = paramstore.GetParameterWithFallback(ctx, jwtExpirationParam, "24h")
	} else {
		jwtExpStr = getEnv("JWT_EXPIRATION", "24h")
	}

	exp, err := time.ParseDuration(jwtExpStr)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRATION %q, defaulting to 24h", jwtExpStr)
		exp = 24 * time.Hour
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
