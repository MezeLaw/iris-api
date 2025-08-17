package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"iris-api/internal/repository"
	"iris-api/pkg/response"
)

type DeletePatientHandler struct {
	repo repository.PatientRepositoryInterface
}

func NewDeletePatientHandler(repo repository.PatientRepositoryInterface) *DeletePatientHandler {
	return &DeletePatientHandler{
		repo: repo,
	}
}

func (h *DeletePatientHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	patientID := request.PathParameters["id"]
	if patientID == "" {
		return response.BadRequest("Patient ID is required"), nil
	}

	existingPatient, err := h.repo.GetPatient(patientID)
	if err != nil {
		return response.InternalServerError("Failed to check patient existence"), nil
	}

	if existingPatient == nil {
		return response.NotFound("Patient not found"), nil
	}

	if err := h.repo.DeletePatient(patientID); err != nil {
		return response.InternalServerError("Failed to delete patient"), nil
	}

	return response.Success(map[string]string{"id": patientID}, "Patient deleted successfully"), nil
}