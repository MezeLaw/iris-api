package entities

import (
	"database/sql"
	"time"
)

// AntecedentesVisuales representa el historial visual de un paciente
type AntecedentesVisuales struct {
	ID         int64 `json:"id" db:"id"`
	PacienteID int64 `json:"paciente_id" db:"paciente_id"`

	// Uso de lentes de contacto
	UsaLentesContacto      bool           `json:"usa_lentes_contacto" db:"usa_lentes_contacto"`
	TipoLenteActual        sql.NullString `json:"tipo_lente_actual,omitempty" db:"tipo_lente_actual"`
	TiempoUsoDiario        sql.NullString `json:"tiempo_uso_diario,omitempty" db:"tiempo_uso_diario"`
	MarcaModelo            sql.NullString `json:"marca_modelo,omitempty" db:"marca_modelo"`
	FechaUltimaAdaptacion  sql.NullTime   `json:"fecha_ultima_adaptacion,omitempty" db:"fecha_ultima_adaptacion"`
	MolestiaComplicaciones sql.NullString `json:"molestias_complicaciones,omitempty" db:"molestias_complicaciones"`

	// Auditoría
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAntecedentesVisualesRequest representa los datos para crear antecedentes visuales
type CreateAntecedentesVisualesRequest struct {
	UsaLentesContacto      bool    `json:"usa_lentes_contacto"`
	TipoLenteActual        *string `json:"tipo_lente_actual,omitempty"`
	TiempoUsoDiario        *string `json:"tiempo_uso_diario,omitempty"`
	MarcaModelo            *string `json:"marca_modelo,omitempty"`
	FechaUltimaAdaptacion  *string `json:"fecha_ultima_adaptacion,omitempty"` // formato: YYYY-MM-DD
	MolestiaComplicaciones *string `json:"molestias_complicaciones,omitempty"`
}

// UpdateAntecedentesVisualesRequest representa los datos para actualizar antecedentes visuales
type UpdateAntecedentesVisualesRequest struct {
	UsaLentesContacto      *bool   `json:"usa_lentes_contacto,omitempty"`
	TipoLenteActual        *string `json:"tipo_lente_actual,omitempty"`
	TiempoUsoDiario        *string `json:"tiempo_uso_diario,omitempty"`
	MarcaModelo            *string `json:"marca_modelo,omitempty"`
	FechaUltimaAdaptacion  *string `json:"fecha_ultima_adaptacion,omitempty"` // formato: YYYY-MM-DD
	MolestiaComplicaciones *string `json:"molestias_complicaciones,omitempty"`
}
