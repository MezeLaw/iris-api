package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"iris-api/internal/handlers"
	"iris-api/internal/repository"
	"iris-api/internal/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db := utils.NewDynamoDBClient()
	tableName := utils.GetTableName()
	repo := repository.NewPatientRepository(db, tableName)

	handler := handlers.NewDeletePatientHandler(repo)
	return handler.Handle(ctx, request)
}

func main() {
	lambda.Start(handler)
}
