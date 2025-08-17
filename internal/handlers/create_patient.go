package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"iris-api/internal/models"
	"iris-api/internal/repository"
	"iris-api/pkg/response"
)

type CreatePatientHandler struct {
	repo repository.PatientRepositoryInterface
}

func NewCreatePatientHandler(repo repository.PatientRepositoryInterface) *CreatePatientHandler {
	return &CreatePatientHandler{
		repo: repo,
	}
}

func (h *CreatePatientHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.CreatePatientRequest

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.BadRequest("Invalid request body"), nil
	}

	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Phone == "" || req.DateOfBirth == "" {
		return response.BadRequest("Missing required fields: first_name, last_name, email, phone, date_of_birth"), nil
	}

	patient, err := h.repo.CreatePatient(req)
	if err != nil {
		return response.InternalServerError("Failed to create patient"), nil
	}

	return response.Success(patient, "Patient created successfully"), nil
}