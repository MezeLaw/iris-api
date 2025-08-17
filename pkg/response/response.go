package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type APIResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(data interface{}, message string) events.APIGatewayProxyResponse {
	response := SuccessResponse{
		Data:    data,
		Message: message,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}
}

func Error(statusCode int, error string, message string) events.APIGatewayProxyResponse {
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}
}

func BadRequest(message string) events.APIGatewayProxyResponse {
	return Error(400, "BAD_REQUEST", message)
}

func NotFound(message string) events.APIGatewayProxyResponse {
	return Error(404, "NOT_FOUND", message)
}

func InternalServerError(message string) events.APIGatewayProxyResponse {
	return Error(500, "INTERNAL_SERVER_ERROR", message)
}
