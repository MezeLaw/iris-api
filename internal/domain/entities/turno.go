package entities

import (
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

type TipoServicio string

const (
	TipoConsulta         TipoServicio = "consulta"
	TipoControlVision    TipoServicio = "control_vision"
	TipoAdaptacionLentes TipoServicio = "adaptacion_lentes"
	TipoSeguimiento      TipoServicio = "seguimiento"
	TipoEmergencia       TipoServicio = "emergencia"
)

type Turno struct {
	ID                  int64        `json:"id" db:"id"`
	PacienteID          int64        `json:"paciente_id" db:"paciente_id"`
	ContactologoID      int64        `json:"contactologo_id" db:"contactologo_id"`
	FechaHora           time.Time    `json:"fecha_hora" db:"fecha_hora"`
	DuracionMinutos     int          `json:"duracion_minutos" db:"duracion_minutos"`
	TipoServicio        TipoServicio `json:"tipo_servicio" db:"tipo_servicio"`
	Estado              EstadoTurno  `json:"estado" db:"estado"`
	Motivo              string       `json:"motivo,omitempty" db:"motivo"`
	Observaciones       string       `json:"observaciones,omitempty" db:"observaciones"`
	RecordatorioEnviado bool         `json:"recordatorio_enviado" db:"recordatorio_enviado"`
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at" db:"updated_at"`
}

type CreateTurnoRequest struct {
	PacienteID      int64        `json:"paciente_id" validate:"required"`
	ContactologoID  int64        `json:"contactologo_id" validate:"required"`
	FechaHora       time.Time    `json:"fecha_hora" validate:"required"`
	DuracionMinutos int          `json:"duracion_minutos" validate:"required,min=15,max=240"`
	TipoServicio    TipoServicio `json:"tipo_servicio" validate:"required,oneof=consulta control_vision adaptacion_lentes seguimiento emergencia"`
	Motivo          string       `json:"motivo,omitempty" validate:"max=500"`
	Observaciones   string       `json:"observaciones,omitempty" validate:"max=1000"`
}

type UpdateTurnoRequest struct {
	FechaHora       *time.Time    `json:"fecha_hora,omitempty"`
	DuracionMinutos *int          `json:"duracion_minutos,omitempty" validate:"omitempty,min=15,max=240"`
	TipoServicio    *TipoServicio `json:"tipo_servicio,omitempty" validate:"omitempty,oneof=consulta control_vision adaptacion_lentes seguimiento emergencia"`
	Estado          *EstadoTurno  `json:"estado,omitempty" validate:"omitempty,oneof=pendiente confirmado cancelado completado no_asistio"`
	Motivo          *string       `json:"motivo,omitempty" validate:"omitempty,max=500"`
	Observaciones   *string       `json:"observaciones,omitempty" validate:"omitempty,max=1000"`
}

type TurnoConDetalles struct {
	Turno
	PacienteNombre     string `json:"paciente_nombre" db:"paciente_nombre"`
	PacienteEmail      string `json:"paciente_email" db:"paciente_email"`
	ContactologoNombre string `json:"contactologo_nombre" db:"contactologo_nombre"`
	ContactologoEmail  string `json:"contactologo_email" db:"contactologo_email"`
}

type TurnoFilter struct {
	PacienteID     *int64        `form:"paciente_id"`
	ContactologoID *int64        `form:"contactologo_id"`
	TipoServicio   *TipoServicio `form:"tipo_servicio"`
	Estado         *EstadoTurno  `form:"estado"`
	FechaDesde     *time.Time    `form:"fecha_desde"`
	FechaHasta     *time.Time    `form:"fecha_hasta"`
	Limit          int           `form:"limit"`
	Offset         int           `form:"offset"`
}

type ProximosTurnosAlert struct {
	TurnosProximas24h  []*TurnoConDetalles `json:"turnos_proximas_24h"`
	TurnosSinConfirmar []*TurnoConDetalles `json:"turnos_sin_confirmar"`
	Count24h           int                 `json:"count_24h"`
	CountSinConfirmar  int                 `json:"count_sin_confirmar"`
}
