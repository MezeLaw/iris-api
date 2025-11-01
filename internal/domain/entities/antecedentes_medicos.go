package entities

import (
	"database/sql"
	"time"
)

// AntecedentesMedicos representa el historial médico de un paciente
type AntecedentesMedicos struct {
	ID         int64 `json:"id" db:"id"`
	PacienteID int64 `json:"paciente_id" db:"paciente_id"`

	// Enfermedades sistémicas
	TieneDiabetes     bool           `json:"tiene_diabetes" db:"tiene_diabetes"`
	TieneHipertension bool           `json:"tiene_hipertension" db:"tiene_hipertension"`
	TieneAlergias     bool           `json:"tiene_alergias" db:"tiene_alergias"`
	DetalleAlergias   sql.NullString `json:"detalle_alergias,omitempty" db:"detalle_alergias"`
	OtrasEnfermedades sql.NullString `json:"otras_enfermedades,omitempty" db:"otras_enfermedades"`

	// Medicación
	MedicacionHabitual sql.NullString `json:"medicacion_habitual,omitempty" db:"medicacion_habitual"`

	// Cirugías
	CirugiasPrevias  sql.NullString `json:"cirugias_previas,omitempty" db:"cirugias_previas"`
	CirugiasOculares sql.NullString `json:"cirugias_oculares,omitempty" db:"cirugias_oculares"`

	// Auditoría
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAntecedentesMedicosRequest representa los datos para crear antecedentes médicos
type CreateAntecedentesMedicosRequest struct {
	TieneDiabetes      bool    `json:"tiene_diabetes"`
	TieneHipertension  bool    `json:"tiene_hipertension"`
	TieneAlergias      bool    `json:"tiene_alergias"`
	DetalleAlergias    *string `json:"detalle_alergias,omitempty"`
	OtrasEnfermedades  *string `json:"otras_enfermedades,omitempty"`
	MedicacionHabitual *string `json:"medicacion_habitual,omitempty"`
	CirugiasPrevias    *string `json:"cirugias_previas,omitempty"`
	CirugiasOculares   *string `json:"cirugias_oculares,omitempty"`
}

// UpdateAntecedentesMedicosRequest representa los datos para actualizar antecedentes médicos
type UpdateAntecedentesMedicosRequest struct {
	TieneDiabetes      *bool   `json:"tiene_diabetes,omitempty"`
	TieneHipertension  *bool   `json:"tiene_hipertension,omitempty"`
	TieneAlergias      *bool   `json:"tiene_alergias,omitempty"`
	DetalleAlergias    *string `json:"detalle_alergias,omitempty"`
	OtrasEnfermedades  *string `json:"otras_enfermedades,omitempty"`
	MedicacionHabitual *string `json:"medicacion_habitual,omitempty"`
	CirugiasPrevias    *string `json:"cirugias_previas,omitempty"`
	CirugiasOculares   *string `json:"cirugias_oculares,omitempty"`
}
