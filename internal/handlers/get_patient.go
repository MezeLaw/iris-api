package handlers

import (
	"context"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"iris-api/internal/repository"
	"iris-api/pkg/response"
)

type GetPatientHandler struct {
	repo repository.PatientRepositoryInterface
}

func NewGetPatientHandler(repo repository.PatientRepositoryInterface) *GetPatientHandler {
	return &GetPatientHandler{
		repo: repo,
	}
}

func (h *GetPatientHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	patientID := request.PathParameters["id"]
	if patientID != "" {
		patient, err := h.repo.GetPatient(patientID)
		if err != nil {
			return response.InternalServerError("Failed to get patient"), nil
		}

		if patient == nil {
			return response.NotFound("Patient not found"), nil
		}

		return response.Success(patient, "Patient retrieved successfully"), nil
	}

	limit := 20
	if limitStr := request.QueryStringParameters["limit"]; limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	lastKey := request.QueryStringParameters["last_key"]

	patients, nextKey, err := h.repo.ListPatients(limit, lastKey)
	if err != nil {
		return response.InternalServerError("Failed to list patients"), nil
	}

	responseData := map[string]interface{}{
		"patients": patients,
		"next_key": nextKey,
		"count":    len(patients),
	}

	return response.Success(responseData, "Patients retrieved successfully"), nil
}