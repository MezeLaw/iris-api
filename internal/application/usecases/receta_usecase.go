package usecases

import (
	"context"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/services"
)

type RecetaUseCase interface {
	CreateReceta(ctx context.Context, req *entities.CreateRecetaRequest) (*entities.Receta, error)
	GetRecetaByID(ctx context.Context, id int64) (*entities.Receta, error)
	GetRecetas(ctx context.Context, limit, offset int) (*GetRecetasResponse, error)
	GetRecetasByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) (*GetRecetasResponse, error)
	GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error)
	CheckDioptriasChange(ctx context.Context, pacienteID int64) (*entities.DioptriasChange, error)
	UpdateReceta(ctx context.Context, id int64, req *entities.UpdateRecetaRequest) (*entities.Receta, error)
	DeleteReceta(ctx context.Context, id int64) error
}

type GetRecetasResponse struct {
	Recetas []*entities.Receta `json:"recetas"`
	Total   int                `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
	HasMore bool               `json:"has_more"`
}

type recetaUseCase struct {
	recetaService services.RecetaService
}

func NewRecetaUseCase(recetaService services.RecetaService) RecetaUseCase {
	return &recetaUseCase{
		recetaService: recetaService,
	}
}

func (uc *recetaUseCase) CreateReceta(ctx context.Context, req *entities.CreateRecetaRequest) (*entities.Receta, error) {
	return uc.recetaService.CreateReceta(ctx, req)
}

func (uc *recetaUseCase) GetRecetaByID(ctx context.Context, id int64) (*entities.Receta, error) {
	return uc.recetaService.GetRecetaByID(ctx, id)
}

func (uc *recetaUseCase) GetRecetas(ctx context.Context, limit, offset int) (*GetRecetasResponse, error) {
	recetas, total, err := uc.recetaService.GetRecetas(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	hasMore := (offset + limit) < total

	return &GetRecetasResponse{
		Recetas: recetas,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	}, nil
}

func (uc *recetaUseCase) GetRecetasByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) (*GetRecetasResponse, error) {
	recetas, total, err := uc.recetaService.GetRecetasByPacienteID(ctx, pacienteID, limit, offset)
	if err != nil {
		return nil, err
	}

	hasMore := (offset + limit) < total

	return &GetRecetasResponse{
		Recetas: recetas,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	}, nil
}

func (uc *recetaUseCase) GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	return uc.recetaService.GetHistorialByPacienteID(ctx, pacienteID)
}

func (uc *recetaUseCase) CheckDioptriasChange(ctx context.Context, pacienteID int64) (*entities.DioptriasChange, error) {
	return uc.recetaService.CheckDioptriasChange(ctx, pacienteID)
}

func (uc *recetaUseCase) UpdateReceta(ctx context.Context, id int64, req *entities.UpdateRecetaRequest) (*entities.Receta, error) {
	return uc.recetaService.UpdateReceta(ctx, id, req)
}

func (uc *recetaUseCase) DeleteReceta(ctx context.Context, id int64) error {
	return uc.recetaService.DeleteReceta(ctx, id)
}