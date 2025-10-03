package services

import (
	"context"
	"errors"
	"time"

	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
)

type TurnoService struct {
	turnoRepo repositories.TurnoRepository
}

func NewTurnoService(turnoRepo repositories.TurnoRepository) *TurnoService {
	return &TurnoService{
		turnoRepo: turnoRepo,
	}
}

func (s *TurnoService) CreateTurno(ctx context.Context, req *entities.CreateTurnoRequest) (*entities.Turno, error) {
	// Validar que la fecha no sea en el pasado
	if req.FechaHora.Before(time.Now()) {
		return nil, errors.New("no se puede crear un turno en el pasado")
	}

	// Verificar disponibilidad del profesional
	disponible, err := s.turnoRepo.CheckDisponibilidad(ctx, req.ContactologoID, req.FechaHora, req.DuracionMinutos, nil)
	if err != nil {
		return nil, err
	}
	if !disponible {
		return nil, errors.New("el profesional no está disponible en ese horario")
	}

	turno := &entities.Turno{
		PacienteID:          req.PacienteID,
		ContactologoID:      req.ContactologoID,
		FechaHora:           req.FechaHora,
		DuracionMinutos:     req.DuracionMinutos,
		TipoServicio:        req.TipoServicio,
		Estado:              entities.EstadoPendiente,
		Motivo:              req.Motivo,
		Observaciones:       req.Observaciones,
		RecordatorioEnviado: false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return s.turnoRepo.Create(ctx, turno)
}

func (s *TurnoService) GetTurnoByID(ctx context.Context, id int64) (*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetByID(ctx, id)
}

func (s *TurnoService) GetTurnos(ctx context.Context, filter *entities.TurnoFilter) ([]*entities.TurnoConDetalles, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return s.turnoRepo.GetAll(ctx, filter)
}

func (s *TurnoService) UpdateTurno(ctx context.Context, id int64, req *entities.UpdateTurnoRequest) (*entities.Turno, error) {
	// Obtener turno existente
	existingTurno, err := s.turnoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validar que no se pueda modificar un turno ya completado o cancelado
	if existingTurno.Estado == entities.EstadoCompletado || existingTurno.Estado == entities.EstadoCancelado {
		return nil, errors.New("no se puede modificar un turno completado o cancelado")
	}

	// Construir turno actualizado
	turno := &entities.Turno{
		ID:                  id,
		PacienteID:          existingTurno.PacienteID,
		ContactologoID:      existingTurno.ContactologoID,
		FechaHora:           existingTurno.FechaHora,
		DuracionMinutos:     existingTurno.DuracionMinutos,
		TipoServicio:        existingTurno.TipoServicio,
		Estado:              existingTurno.Estado,
		Motivo:              existingTurno.Motivo,
		Observaciones:       existingTurno.Observaciones,
		RecordatorioEnviado: existingTurno.RecordatorioEnviado,
		UpdatedAt:           time.Now(),
	}

	// Aplicar cambios
	if req.FechaHora != nil {
		if req.FechaHora.Before(time.Now()) {
			return nil, errors.New("no se puede programar un turno en el pasado")
		}

		duracion := existingTurno.DuracionMinutos
		if req.DuracionMinutos != nil {
			duracion = *req.DuracionMinutos
		}

		// Verificar disponibilidad si se cambia la fecha/hora
		disponible, err := s.turnoRepo.CheckDisponibilidad(ctx, existingTurno.ContactologoID, *req.FechaHora, duracion, &id)
		if err != nil {
			return nil, err
		}
		if !disponible {
			return nil, errors.New("el profesional no está disponible en ese horario")
		}

		turno.FechaHora = *req.FechaHora
	}

	if req.DuracionMinutos != nil {
		turno.DuracionMinutos = *req.DuracionMinutos
	}
	if req.TipoServicio != nil {
		turno.TipoServicio = *req.TipoServicio
	}
	if req.Estado != nil {
		turno.Estado = *req.Estado
	}
	if req.Motivo != nil {
		turno.Motivo = *req.Motivo
	}
	if req.Observaciones != nil {
		turno.Observaciones = *req.Observaciones
	}

	return s.turnoRepo.Update(ctx, id, turno)
}

func (s *TurnoService) CancelTurno(ctx context.Context, id int64, motivo string) error {
	// Obtener turno existente
	existingTurno, err := s.turnoRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Validar que no esté ya cancelado o completado
	if existingTurno.Estado == entities.EstadoCancelado {
		return errors.New("el turno ya está cancelado")
	}
	if existingTurno.Estado == entities.EstadoCompletado {
		return errors.New("no se puede cancelar un turno completado")
	}

	return s.turnoRepo.CancelTurno(ctx, id, motivo)
}

func (s *TurnoService) DeleteTurno(ctx context.Context, id int64) error {
	return s.turnoRepo.Delete(ctx, id)
}

// Vistas específicas
func (s *TurnoService) GetTurnosByDia(ctx context.Context, fecha time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetByDia(ctx, fecha, contactologoID)
}

func (s *TurnoService) GetTurnosBySemana(ctx context.Context, fechaInicio time.Time, contactologoID *int64) ([]*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetBySemana(ctx, fechaInicio, contactologoID)
}

func (s *TurnoService) GetTurnosByProfesional(ctx context.Context, contactologoID int64, fechaDesde, fechaHasta time.Time) ([]*entities.TurnoConDetalles, error) {
	return s.turnoRepo.GetByProfesional(ctx, contactologoID, fechaDesde, fechaHasta)
}

// Alertas
func (s *TurnoService) GetProximosTurnosAlert(ctx context.Context) (*entities.ProximosTurnosAlert, error) {
	// Obtener turnos de las próximas 24 horas
	turnos24h, err := s.turnoRepo.GetProximosTurnos(ctx, 24)
	if err != nil {
		return nil, err
	}

	// Obtener turnos sin confirmar
	turnosSinConfirmar, err := s.turnoRepo.GetTurnosSinConfirmar(ctx)
	if err != nil {
		return nil, err
	}

	return &entities.ProximosTurnosAlert{
		TurnosProximas24h:  turnos24h,
		TurnosSinConfirmar: turnosSinConfirmar,
		Count24h:           len(turnos24h),
		CountSinConfirmar:  len(turnosSinConfirmar),
	}, nil
}

func (s *TurnoService) CountTurnos(ctx context.Context, filter *entities.TurnoFilter) (int, error) {
	return s.turnoRepo.Count(ctx, filter)
}
