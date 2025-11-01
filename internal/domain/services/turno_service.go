package services

import (
	"context"
	"errors"
	"time"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
)

type TurnoService interface {
	CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest, clientID int64) (*entities.Turno, error)
	GetTurnoByID(ctx context.Context, id int64, clientID int64) (*entities.TurnoConDetalles, error)
	GetTurnos(ctx context.Context, filters *entities.TurnoFilters, clientID int64) ([]*entities.TurnoConDetalles, int, error)
	UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest, clientID int64) (*entities.Turno, error)
	DeleteTurno(ctx context.Context, id int64, clientID int64) error
	CambiarEstado(ctx context.Context, id int64, req *entities.CambiarEstadoRequest, clientID int64) error

	// Vistas específicas
	GetTurnosByDia(ctx context.Context, fecha time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error)
	GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error)
	GetTurnosByProfesional(ctx context.Context, profesionalID int64, fechaDesde, fechaHasta *time.Time, clientID int64) ([]*entities.TurnoConDetalles, error)

	// Validaciones
	CheckDisponibilidad(ctx context.Context, req *entities.DisponibilidadRequest, clientID int64) (*entities.DisponibilidadResponse, error)
}

type turnoService struct {
	turnoRepo repositories.TurnoRepository
}

func NewTurnoService(turnoRepo repositories.TurnoRepository) TurnoService {
	return &turnoService{
		turnoRepo: turnoRepo,
	}
}

func (s *turnoService) CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest, clientID int64) (*entities.Turno, error) {
	// Validar que la fecha no sea en el pasado
	if req.FechaHora.Before(time.Now()) {
		return nil, errors.New("no se puede crear un turno en el pasado")
	}

	// Verificar disponibilidad del profesional
	dispReq := &entities.DisponibilidadRequest{
		ProfesionalUserID: req.ProfesionalUserID,
		FechaHora:         req.FechaHora,
		DuracionMinutos:   req.DuracionMinutos,
	}

	dispResp, err := s.turnoRepo.CheckDisponibilidad(ctx, dispReq, clientID)
	if err != nil {
		return nil, err
	}
	if !dispResp.Disponible {
		return nil, errors.New(dispResp.Mensaje)
	}

	// Calcular hora_fin
	horaFin := req.FechaHora.Add(time.Duration(req.DuracionMinutos) * time.Minute)

	turno := &entities.Turno{
		ClientID:          clientID,
		PacienteID:        req.PacienteID,
		ProfesionalUserID: req.ProfesionalUserID,
		TipoServicio:      req.TipoServicio,
		FechaHora:         req.FechaHora,
		DuracionMinutos:   req.DuracionMinutos,
		HoraFin:           horaFin,
		Estado:            entities.EstadoPendiente,
		Observaciones:     req.Observaciones,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return s.turnoRepo.Create(ctx, turno)
}

func (s *turnoService) GetTurnoByID(ctx context.Context, id int64, clientID int64) (*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetByID(ctx, id, clientID)
}

func (s *turnoService) GetTurnos(ctx context.Context, filters *entities.TurnoFilters, clientID int64) ([]*entities.TurnoConDetalles, int, error) {
	// Validar paginación
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.PageSize > 100 {
		filters.PageSize = 100
	}

	return s.turnoRepo.GetAll(ctx, filters, clientID)
}

func (s *turnoService) UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest, clientID int64) (*entities.Turno, error) {
	// Obtener turno existente
	existingTurno, err := s.turnoRepo.GetByID(ctx, id, clientID)
	if err != nil {
		return nil, err
	}

	// Validar que no se pueda modificar un turno ya completado o cancelado
	if existingTurno.Estado == entities.EstadoCompletado || existingTurno.Estado == entities.EstadoCancelado {
		return nil, errors.New("no se puede modificar un turno completado o cancelado")
	}

	// Construir turno actualizado (empezar con valores existentes)
	turno := &entities.Turno{
		ID:                id,
		ClientID:          clientID,
		PacienteID:        existingTurno.PacienteID,
		ProfesionalUserID: existingTurno.ProfesionalUserID,
		TipoServicio:      existingTurno.TipoServicio,
		FechaHora:         existingTurno.FechaHora,
		DuracionMinutos:   existingTurno.DuracionMinutos,
		HoraFin:           existingTurno.HoraFin,
		Estado:            existingTurno.Estado,
		Observaciones:     existingTurno.Observaciones,
		UpdatedAt:         time.Now(),
	}

	// Aplicar cambios
	if req.ProfesionalUserID != nil {
		turno.ProfesionalUserID = *req.ProfesionalUserID
	}
	if req.TipoServicio != nil {
		turno.TipoServicio = *req.TipoServicio
	}
	if req.DuracionMinutos != nil {
		turno.DuracionMinutos = *req.DuracionMinutos
	}
	if req.Observaciones != nil {
		turno.Observaciones = *req.Observaciones
	}

	// Si se cambia la fecha/hora, validar disponibilidad
	if req.FechaHora != nil {
		if req.FechaHora.Before(time.Now()) {
			return nil, errors.New("no se puede programar un turno en el pasado")
		}

		// Verificar disponibilidad
		dispReq := &entities.DisponibilidadRequest{
			ProfesionalUserID: turno.ProfesionalUserID,
			FechaHora:         *req.FechaHora,
			DuracionMinutos:   turno.DuracionMinutos,
			TurnoID:           &id, // Excluir el turno actual
		}

		dispResp, err := s.turnoRepo.CheckDisponibilidad(ctx, dispReq, clientID)
		if err != nil {
			return nil, err
		}
		if !dispResp.Disponible {
			return nil, errors.New(dispResp.Mensaje)
		}

		turno.FechaHora = *req.FechaHora
	}

	// Recalcular hora_fin
	turno.HoraFin = turno.FechaHora.Add(time.Duration(turno.DuracionMinutos) * time.Minute)

	return s.turnoRepo.Update(ctx, id, turno, clientID)
}

func (s *turnoService) DeleteTurno(ctx context.Context, id int64, clientID int64) error {
	// Verificar que el turno existe
	_, err := s.turnoRepo.GetByID(ctx, id, clientID)
	if err != nil {
		return err
	}

	return s.turnoRepo.Delete(ctx, id, clientID)
}

func (s *turnoService) CambiarEstado(ctx context.Context, id int64, req *entities.CambiarEstadoRequest, clientID int64) error {
	// Validar que el estado sea válido
	if !req.Estado.IsValid() {
		return errors.New("estado inválido")
	}

	// Obtener turno existente
	existingTurno, err := s.turnoRepo.GetByID(ctx, id, clientID)
	if err != nil {
		return err
	}

	// Validaciones de estado
	if existingTurno.Estado == entities.EstadoCancelado && req.Estado != entities.EstadoCancelado {
		return errors.New("no se puede cambiar el estado de un turno cancelado")
	}

	if existingTurno.Estado == entities.EstadoCompletado && req.Estado != entities.EstadoCompletado {
		return errors.New("no se puede cambiar el estado de un turno completado")
	}

	return s.turnoRepo.CambiarEstado(ctx, id, req.Estado, clientID)
}

// Vistas específicas
func (s *turnoService) GetTurnosByDia(ctx context.Context, fecha time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetByDia(ctx, fecha, profesionalID, clientID)
}

func (s *turnoService) GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, profesionalID *int64, clientID int64) ([]*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetBySemana(ctx, fechaInicio, profesionalID, clientID)
}

func (s *turnoService) GetTurnosByProfesional(ctx context.Context, profesionalID int64, fechaDesde, fechaHasta *time.Time, clientID int64) ([]*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetByProfesional(ctx, profesionalID, fechaDesde, fechaHasta, clientID)
}

// Validaciones
func (s *turnoService) CheckDisponibilidad(ctx context.Context, req *entities.DisponibilidadRequest, clientID int64) (*entities.DisponibilidadResponse, error) {
	return s.turnoRepo.CheckDisponibilidad(ctx, req, clientID)
}
