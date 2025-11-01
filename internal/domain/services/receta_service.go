package services

import (
	"context"
	"fmt"
	"math"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
)

type RecetaService interface {
	CreateReceta(ctx context.Context, req *entities.CreateRecetaRequest) (*entities.Receta, error)
	GetRecetaByID(ctx context.Context, id int64) (*entities.Receta, error)
	GetRecetas(ctx context.Context, limit, offset int) ([]*entities.Receta, int, error)
	GetRecetasByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) ([]*entities.Receta, int, error)
	GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error)
	CheckDioptriasChange(ctx context.Context, pacienteID int64) (*entities.DioptriasChange, error)
	UpdateReceta(ctx context.Context, id int64, req *entities.UpdateRecetaRequest) (*entities.Receta, error)
	DeleteReceta(ctx context.Context, id int64) error
}

type recetaService struct {
	recetaRepo repositories.RecetaRepository
}

func NewRecetaService(recetaRepo repositories.RecetaRepository) RecetaService {
	return &recetaService{
		recetaRepo: recetaRepo,
	}
}

func (s *recetaService) CreateReceta(ctx context.Context, req *entities.CreateRecetaRequest) (*entities.Receta, error) {
	if err := s.validateCreateRecetaRequest(req); err != nil {
		return nil, err
	}

	receta := &entities.Receta{
		PacienteID:    req.PacienteID,
		Fecha:         req.Fecha,
		ODEsfera:      req.ODEsfera,
		ODCilindro:    req.ODCilindro,
		ODEje:         req.ODEje,
		OIEsfera:      req.OIEsfera,
		OICilindro:    req.OICilindro,
		OIEje:         req.OIEje,
		TipoLente:     req.TipoLente,
		Observaciones: req.Observaciones,
	}

	return s.recetaRepo.Create(ctx, receta)
}

func (s *recetaService) GetRecetaByID(ctx context.Context, id int64) (*entities.Receta, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid receta ID")
	}

	return s.recetaRepo.GetByID(ctx, id)
}

func (s *recetaService) GetRecetas(ctx context.Context, limit, offset int) ([]*entities.Receta, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	recetas, err := s.recetaRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.recetaRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return recetas, total, nil
}

func (s *recetaService) GetRecetasByPacienteID(ctx context.Context, pacienteID int64, limit, offset int) ([]*entities.Receta, int, error) {
	if pacienteID <= 0 {
		return nil, 0, fmt.Errorf("invalid paciente ID")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	recetas, err := s.recetaRepo.GetByPacienteID(ctx, pacienteID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.recetaRepo.CountByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, 0, err
	}

	return recetas, total, nil
}

func (s *recetaService) GetHistorialByPacienteID(ctx context.Context, pacienteID int64) ([]*entities.Receta, error) {
	if pacienteID <= 0 {
		return nil, fmt.Errorf("invalid paciente ID")
	}

	return s.recetaRepo.GetHistorialByPacienteID(ctx, pacienteID)
}

func (s *recetaService) CheckDioptriasChange(ctx context.Context, pacienteID int64) (*entities.DioptriasChange, error) {
	if pacienteID <= 0 {
		return nil, fmt.Errorf("invalid paciente ID")
	}

	recetas, err := s.recetaRepo.GetLastTwoByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, err
	}

	if len(recetas) < 2 {
		return &entities.DioptriasChange{
			ODChange: 0,
			OIChange: 0,
			HasAlert: false,
			Message:  "No hay suficientes recetas para comparar",
		}, nil
	}

	// recetas[0] es la más reciente, recetas[1] es la anterior
	latest := recetas[0]
	previous := recetas[1]

	odChange := math.Abs(latest.ODEsfera - previous.ODEsfera)
	oiChange := math.Abs(latest.OIEsfera - previous.OIEsfera)

	hasAlert := odChange > 0.25 || oiChange > 0.25

	message := "Sin cambios significativos"
	if hasAlert {
		message = fmt.Sprintf("⚠️ Cambio significativo detectado: OD %.2fD, OI %.2fD", odChange, oiChange)
	}

	return &entities.DioptriasChange{
		ODChange: odChange,
		OIChange: oiChange,
		HasAlert: hasAlert,
		Message:  message,
	}, nil
}

func (s *recetaService) UpdateReceta(ctx context.Context, id int64, req *entities.UpdateRecetaRequest) (*entities.Receta, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid receta ID")
	}

	if err := s.validateUpdateRecetaRequest(req); err != nil {
		return nil, err
	}

	existingReceta, err := s.recetaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Fecha != nil {
		existingReceta.Fecha = *req.Fecha
	}
	if req.ODEsfera != nil {
		existingReceta.ODEsfera = *req.ODEsfera
	}
	if req.ODCilindro != nil {
		existingReceta.ODCilindro = *req.ODCilindro
	}
	if req.ODEje != nil {
		existingReceta.ODEje = *req.ODEje
	}
	if req.OIEsfera != nil {
		existingReceta.OIEsfera = *req.OIEsfera
	}
	if req.OICilindro != nil {
		existingReceta.OICilindro = *req.OICilindro
	}
	if req.OIEje != nil {
		existingReceta.OIEje = *req.OIEje
	}
	if req.TipoLente != nil {
		existingReceta.TipoLente = *req.TipoLente
	}
	if req.Observaciones != nil {
		existingReceta.Observaciones = *req.Observaciones
	}

	return s.recetaRepo.Update(ctx, id, existingReceta)
}

func (s *recetaService) DeleteReceta(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid receta ID")
	}

	return s.recetaRepo.Delete(ctx, id)
}

func (s *recetaService) validateCreateRecetaRequest(req *entities.CreateRecetaRequest) error {
	if req.PacienteID <= 0 {
		return fmt.Errorf("paciente_id is required")
	}
	if req.Fecha.IsZero() {
		return fmt.Errorf("fecha is required")
	}
	if req.ODEsfera < -20 || req.ODEsfera > 20 {
		return fmt.Errorf("od_esfera must be between -20 and 20")
	}
	if req.OIEsfera < -20 || req.OIEsfera > 20 {
		return fmt.Errorf("oi_esfera must be between -20 and 20")
	}
	if req.TipoLente == "" {
		return fmt.Errorf("tipo_lente is required")
	}
	validTipos := map[string]bool{
		"monofocal":  true,
		"bifocal":    true,
		"multifocal": true,
		"progresivo": true,
	}
	if !validTipos[req.TipoLente] {
		return fmt.Errorf("tipo_lente must be one of: monofocal, bifocal, multifocal, progresivo")
	}
	return nil
}

func (s *recetaService) validateUpdateRecetaRequest(req *entities.UpdateRecetaRequest) error {
	if req.ODEsfera != nil && (*req.ODEsfera < -20 || *req.ODEsfera > 20) {
		return fmt.Errorf("od_esfera must be between -20 and 20")
	}
	if req.OIEsfera != nil && (*req.OIEsfera < -20 || *req.OIEsfera > 20) {
		return fmt.Errorf("oi_esfera must be between -20 and 20")
	}
	if req.TipoLente != nil {
		validTipos := map[string]bool{
			"monofocal":  true,
			"bifocal":    true,
			"multifocal": true,
			"progresivo": true,
		}
		if !validTipos[*req.TipoLente] {
			return fmt.Errorf("tipo_lente must be one of: monofocal, bifocal, multifocal, progresivo")
		}
	}
	return nil
}
