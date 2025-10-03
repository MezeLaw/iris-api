package entities

import "time"

type PacienteActivo struct {
	ID               int64     `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Email            string    `json:"email" db:"email"`
	UltimoTurno      time.Time `json:"ultimo_turno" db:"ultimo_turno"`
	TotalTurnos      int       `json:"total_turnos" db:"total_turnos"`
	TurnosPendientes int       `json:"turnos_pendientes" db:"turnos_pendientes"`
}

type PacienteInactivo struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	UltimoTurno time.Time `json:"ultimo_turno" db:"ultimo_turno"`
	DiasInactivo int      `json:"dias_inactivo" db:"dias_inactivo"`
}

type ReportePacientesActivos struct {
	Pacientes []*PacienteActivo `json:"pacientes"`
	Total     int               `json:"total"`
}

type ReportePacientesInactivos struct {
	Pacientes []*PacienteInactivo `json:"pacientes"`
	Total     int                 `json:"total"`
}
