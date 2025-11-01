package entities

import (
	"database/sql"
	"time"
)

type EstadoTurno string

const (
	EstadoPendiente  EstadoTurno = "pendiente"
	EstadoConfirmado EstadoTurno = "confirmado"
	EstadoCancelado  EstadoTurno = "cancelado"
	EstadoCompletado EstadoTurno = "completado"
	EstadoNoAsistio  EstadoTurno = "no_asistio"
)

// IsValid checks if the estado is valid
func (e EstadoTurno) IsValid() bool {
	switch e {
	case EstadoPendiente, EstadoConfirmado, EstadoCancelado, EstadoCompletado, EstadoNoAsistio:
		return true
	default:
		return false
	}
}

type Turno struct {
	ID                int64        `json:"id" db:"id"`
	ClientID          int64        `json:"client_id" db:"client_id"`
	PacienteID        int64        `json:"paciente_id" db:"paciente_id"`
	ProfesionalUserID int64        `json:"profesional_user_id" db:"profesional_user_id"`
	TipoServicio      string       `json:"tipo_servicio" db:"tipo_servicio"`
	FechaHora         time.Time    `json:"fecha_hora" db:"fecha_hora"`
	DuracionMinutos   int          `json:"duracion_minutos" db:"duracion_minutos"`
	HoraFin           time.Time    `json:"hora_fin" db:"hora_fin"`
	Estado            EstadoTurno  `json:"estado" db:"estado"`
	Observaciones     string       `json:"observaciones,omitempty" db:"observaciones"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CreateTurnoRequest struct {
	PacienteID        int64     `json:"paciente_id" validate:"required"`
	ProfesionalUserID int64     `json:"profesional_user_id" validate:"required"`
	TipoServicio      string    `json:"tipo_servicio" validate:"required,max=100"`
	FechaHora         time.Time `json:"fecha_hora" validate:"required"`
	DuracionMinutos   int       `json:"duracion_minutos" validate:"required,min=15,max=240"`
	Observaciones     string    `json:"observaciones,omitempty" validate:"max=1000"`
}

type UpdateTurnoRequest struct {
	ProfesionalUserID *int64     `json:"profesional_user_id,omitempty"`
	TipoServicio      *string    `json:"tipo_servicio,omitempty" validate:"omitempty,max=100"`
	FechaHora         *time.Time `json:"fecha_hora,omitempty"`
	DuracionMinutos   *int       `json:"duracion_minutos,omitempty" validate:"omitempty,min=15,max=240"`
	Observaciones     *string    `json:"observaciones,omitempty" validate:"omitempty,max=1000"`
}

type CambiarEstadoRequest struct {
	Estado EstadoTurno `json:"estado" validate:"required,oneof=pendiente confirmado cancelado completado no_asistio"`
}

type TurnoConDetalles struct {
	Turno
	PacienteNombre      string `json:"paciente_nombre" db:"paciente_nombre"`
	PacienteApellido    string `json:"paciente_apellido" db:"paciente_apellido"`
	PacienteEmail       string `json:"paciente_email" db:"paciente_email"`
	ProfesionalNombre   string `json:"profesional_nombre" db:"profesional_nombre"`
	ProfesionalApellido string `json:"profesional_apellido" db:"profesional_apellido"`
	ProfesionalEmail    string `json:"profesional_email" db:"profesional_email"`
}

type TurnoFilters struct {
	FechaDesde        *time.Time   `form:"fecha_desde"`
	FechaHasta        *time.Time   `form:"fecha_hasta"`
	ProfesionalUserID *int64       `form:"profesional_user_id"`
	PacienteID        *int64       `form:"paciente_id"`
	Estado            *EstadoTurno `form:"estado"`
	Page              int          `form:"page"`
	PageSize          int          `form:"page_size"`
}

type DisponibilidadRequest struct {
	ProfesionalUserID int64     `json:"profesional_user_id" validate:"required"`
	FechaHora         time.Time `json:"fecha_hora" validate:"required"`
	DuracionMinutos   int       `json:"duracion_minutos" validate:"required,min=15"`
	TurnoID           *int64    `json:"turno_id,omitempty"` // Para excluir el turno al editar
}

type DisponibilidadResponse struct {
	Disponible      bool    `json:"disponible"`
	Mensaje         string  `json:"mensaje,omitempty"`
	TurnosConflicto []Turno `json:"turnos_conflicto,omitempty"`
}
