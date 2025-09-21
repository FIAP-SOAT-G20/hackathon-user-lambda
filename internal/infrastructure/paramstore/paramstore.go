package paramstore

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var (
	client *ssm.Client
	once   sync.Once
)

func getClient(ctx context.Context) *ssm.Client {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("Failed to load AWS config: %v", err)
		}
		client = ssm.NewFromConfig(cfg)
	})
	return client
}

// GetParameter retrieves a parameter from AWS Systems Manager Parameter Store
func GetParameter(ctx context.Context, parameterName string) (string, error) {
	if parameterName == "" {
		return "", nil
	}

	ssmClient := getClient(ctx)

	input := &ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true),
	}

	result, err := ssmClient.GetParameter(ctx, input)
	if err != nil {
		log.Printf("Error getting parameter %s: %v", parameterName, err)
		return "", err
	}

	if result.Parameter == nil || result.Parameter.Value == nil {
		return "", nil
	}

	return *result.Parameter.Value, nil
}

// GetParameterWithFallback retrieves a parameter from Parameter Store with fallback to environment variable
func GetParameterWithFallback(ctx context.Context, parameterName, fallbackValue string) string {
	if parameterName == "" {
		return fallbackValue
	}

	value, err := GetParameter(ctx, parameterName)
	if err != nil || value == "" {
		log.Printf("Using fallback value for parameter %s", parameterName)
		return fallbackValue
	}

	return value
}
