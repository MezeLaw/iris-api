package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"iris-api/internal/models"
	"iris-api/internal/repository"
	"iris-api/pkg/response"
)

type UpdatePatientHandler struct {
	repo repository.PatientRepositoryInterface
}

func NewUpdatePatientHandler(repo repository.PatientRepositoryInterface) *UpdatePatientHandler {
	return &UpdatePatientHandler{
		repo: repo,
	}
}

func (h *UpdatePatientHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	patientID := request.PathParameters["id"]
	if patientID == "" {
		return response.BadRequest("Patient ID is required"), nil
	}

	var req models.UpdatePatientRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.BadRequest("Invalid request body"), nil
	}

	existingPatient, err := h.repo.GetPatient(patientID)
	if err != nil {
		return response.InternalServerError("Failed to check patient existence"), nil
	}

	if existingPatient == nil {
		return response.NotFound("Patient not found"), nil
	}

	updatedPatient, err := h.repo.UpdatePatient(patientID, req)
	if err != nil {
		return response.InternalServerError("Failed to update patient"), nil
	}

	return response.Success(updatedPatient, "Patient updated successfully"), nil
}