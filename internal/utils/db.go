package utils

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewDynamoDBClient() *dynamodb.DynamoDB {
	config := &aws.Config{
		Region: aws.String(getEnv("AWS_REGION", "us-east-1")),
	}

	// Configure for local development
	if getEnv("IS_OFFLINE", "false") == "true" {
		config.Endpoint = aws.String("http://localhost:8000")
		config.Credentials = credentials.NewStaticCredentials("dummy", "dummy", "")
	}

	sess := session.Must(session.NewSession(config))
	return dynamodb.New(sess)
}

func GetTableName() string {
	return getEnv("PATIENTS_TABLE", "patients")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
