package usecases

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type PacienteUseCase interface {
	// CRUD Pacientes
	CreatePaciente(ctx context.Context, req *entities.CreatePacienteRequest) (*entities.Paciente, error)
	GetPacienteByID(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error)
	GetPacienteComplete(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error)
	ListPacientes(ctx context.Context, clientID int64, page, pageSize int) (*entities.PacienteListResponse, error)
	SearchPacientes(ctx context.Context, clientID int64, query string, page, pageSize int) (*entities.PacienteListResponse, error)
	UpdatePaciente(ctx context.Context, id int64, clientID int64, req *entities.UpdatePacienteRequest) (*entities.Paciente, error)
	DeletePaciente(ctx context.Context, id int64, clientID int64) error

	// Antecedentes
	CreateOrUpdateAntecedentesMedicos(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesMedicosRequest) (*entities.AntecedentesMedicos, error)
	GetAntecedentesMedicos(ctx context.Context, pacienteID int64) (*entities.AntecedentesMedicos, error)
	CreateOrUpdateAntecedentesVisuales(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesVisualesRequest) (*entities.AntecedentesVisuales, error)
	GetAntecedentesVisuales(ctx context.Context, pacienteID int64) (*entities.AntecedentesVisuales, error)

	// Exámenes Visuales
	CreateExamenVisual(ctx context.Context, req *entities.CreateExamenVisualRequest) (*entities.ExamenVisual, error)
	GetExamenVisual(ctx context.Context, id int64) (*entities.ExamenVisual, error)
	GetExamenesVisualesByPaciente(ctx context.Context, pacienteID int64) ([]entities.ExamenVisual, error)
	UpdateExamenVisual(ctx context.Context, id int64, req *entities.UpdateExamenVisualRequest) (*entities.ExamenVisual, error)
	DeleteExamenVisual(ctx context.Context, id int64) error
	CompararExamenes(ctx context.Context, examenAnteriorID, examenActualID int64) (*entities.DiferenciaRefraccion, error)
}

type pacienteUseCase struct {
	pacienteService services.PacienteService
}

func NewPacienteUseCase(pacienteService services.PacienteService) PacienteUseCase {
	return &pacienteUseCase{
		pacienteService: pacienteService,
	}
}

// CRUD Pacientes

func (uc *pacienteUseCase) CreatePaciente(ctx context.Context, req *entities.CreatePacienteRequest) (*entities.Paciente, error) {
	return uc.pacienteService.CreatePaciente(ctx, req)
}

func (uc *pacienteUseCase) GetPacienteByID(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	return uc.pacienteService.GetPacienteByID(ctx, id, clientID)
}

func (uc *pacienteUseCase) GetPacienteComplete(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	return uc.pacienteService.GetPacienteComplete(ctx, id, clientID)
}

func (uc *pacienteUseCase) ListPacientes(ctx context.Context, clientID int64, page, pageSize int) (*entities.PacienteListResponse, error) {
	return uc.pacienteService.ListPacientes(ctx, clientID, page, pageSize)
}

func (uc *pacienteUseCase) SearchPacientes(ctx context.Context, clientID int64, query string, page, pageSize int) (*entities.PacienteListResponse, error) {
	return uc.pacienteService.SearchPacientes(ctx, clientID, query, page, pageSize)
}

func (uc *pacienteUseCase) UpdatePaciente(ctx context.Context, id int64, clientID int64, req *entities.UpdatePacienteRequest) (*entities.Paciente, error) {
	return uc.pacienteService.UpdatePaciente(ctx, id, clientID, req)
}

func (uc *pacienteUseCase) DeletePaciente(ctx context.Context, id int64, clientID int64) error {
	return uc.pacienteService.DeletePaciente(ctx, id, clientID)
}

// Antecedentes

func (uc *pacienteUseCase) CreateOrUpdateAntecedentesMedicos(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesMedicosRequest) (*entities.AntecedentesMedicos, error) {
	return uc.pacienteService.CreateOrUpdateAntecedentesMedicos(ctx, pacienteID, req)
}

func (uc *pacienteUseCase) GetAntecedentesMedicos(ctx context.Context, pacienteID int64) (*entities.AntecedentesMedicos, error) {
	return uc.pacienteService.GetAntecedentesMedicos(ctx, pacienteID)
}

func (uc *pacienteUseCase) CreateOrUpdateAntecedentesVisuales(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesVisualesRequest) (*entities.AntecedentesVisuales, error) {
	return uc.pacienteService.CreateOrUpdateAntecedentesVisuales(ctx, pacienteID, req)
}

func (uc *pacienteUseCase) GetAntecedentesVisuales(ctx context.Context, pacienteID int64) (*entities.AntecedentesVisuales, error) {
	return uc.pacienteService.GetAntecedentesVisuales(ctx, pacienteID)
}

// Exámenes Visuales

func (uc *pacienteUseCase) CreateExamenVisual(ctx context.Context, req *entities.CreateExamenVisualRequest) (*entities.ExamenVisual, error) {
	return uc.pacienteService.CreateExamenVisual(ctx, req)
}

func (uc *pacienteUseCase) GetExamenVisual(ctx context.Context, id int64) (*entities.ExamenVisual, error) {
	return uc.pacienteService.GetExamenVisual(ctx, id)
}

func (uc *pacienteUseCase) GetExamenesVisualesByPaciente(ctx context.Context, pacienteID int64) ([]entities.ExamenVisual, error) {
	return uc.pacienteService.GetExamenesVisualesByPaciente(ctx, pacienteID)
}

func (uc *pacienteUseCase) UpdateExamenVisual(ctx context.Context, id int64, req *entities.UpdateExamenVisualRequest) (*entities.ExamenVisual, error) {
	return uc.pacienteService.UpdateExamenVisual(ctx, id, req)
}

func (uc *pacienteUseCase) DeleteExamenVisual(ctx context.Context, id int64) error {
	return uc.pacienteService.DeleteExamenVisual(ctx, id)
}

func (uc *pacienteUseCase) CompararExamenes(ctx context.Context, examenAnteriorID, examenActualID int64) (*entities.DiferenciaRefraccion, error) {
	return uc.pacienteService.CompararExamenes(ctx, examenAnteriorID, examenActualID)
}
