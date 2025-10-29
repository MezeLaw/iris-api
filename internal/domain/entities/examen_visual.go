package entities

import (
	"database/sql"
	"time"
)

// ExamenVisual representa un examen visual/refracción de un paciente
type ExamenVisual struct {
	ID          int64     `json:"id" db:"id"`
	PacienteID  int64     `json:"paciente_id" db:"paciente_id"`
	FechaExamen time.Time `json:"fecha_examen" db:"fecha_examen"`

	// Agudeza visual sin corrección
	AvScOD sql.NullString `json:"av_sc_od,omitempty" db:"av_sc_od"` // Ojo derecho
	AvScOI sql.NullString `json:"av_sc_oi,omitempty" db:"av_sc_oi"` // Ojo izquierdo

	// Agudeza visual con corrección
	AvCcOD sql.NullString `json:"av_cc_od,omitempty" db:"av_cc_od"`
	AvCcOI sql.NullString `json:"av_cc_oi,omitempty" db:"av_cc_oi"`

	// Refracción Ojo Derecho (OD)
	ODEsfera   sql.NullFloat64 `json:"od_esfera,omitempty" db:"od_esfera"`
	ODCilindro sql.NullFloat64 `json:"od_cilindro,omitempty" db:"od_cilindro"`
	ODEje      sql.NullInt64   `json:"od_eje,omitempty" db:"od_eje"`
	ODAdd      sql.NullFloat64 `json:"od_add,omitempty" db:"od_add"`

	// Refracción Ojo Izquierdo (OI)
	OIEsfera   sql.NullFloat64 `json:"oi_esfera,omitempty" db:"oi_esfera"`
	OICilindro sql.NullFloat64 `json:"oi_cilindro,omitempty" db:"oi_cilindro"`
	OIEje      sql.NullInt64   `json:"oi_eje,omitempty" db:"oi_eje"`
	OIAdd      sql.NullFloat64 `json:"oi_add,omitempty" db:"oi_add"`

	// Tipo de lente indicado
	TipoLente     sql.NullString `json:"tipo_lente,omitempty" db:"tipo_lente"`
	TipoLenteOtro sql.NullString `json:"tipo_lente_otro,omitempty" db:"tipo_lente_otro"`

	// Observaciones y profesional
	Observaciones      sql.NullString `json:"observaciones,omitempty" db:"observaciones"`
	RealizadoPorUserID sql.NullInt64  `json:"realizado_por_user_id,omitempty" db:"realizado_por_user_id"`

	// Auditoría
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateExamenVisualRequest representa los datos para crear un examen visual
type CreateExamenVisualRequest struct {
	PacienteID  int64   `json:"paciente_id" binding:"required"`
	FechaExamen *string `json:"fecha_examen,omitempty"` // formato: YYYY-MM-DD, default: hoy

	// Agudeza visual sin corrección
	AvScOD *string `json:"av_sc_od,omitempty"`
	AvScOI *string `json:"av_sc_oi,omitempty"`

	// Agudeza visual con corrección
	AvCcOD *string `json:"av_cc_od,omitempty"`
	AvCcOI *string `json:"av_cc_oi,omitempty"`

	// Refracción Ojo Derecho (OD)
	ODEsfera   *float64 `json:"od_esfera,omitempty"`
	ODCilindro *float64 `json:"od_cilindro,omitempty"`
	ODEje      *int     `json:"od_eje,omitempty"`
	ODAdd      *float64 `json:"od_add,omitempty"`

	// Refracción Ojo Izquierdo (OI)
	OIEsfera   *float64 `json:"oi_esfera,omitempty"`
	OICilindro *float64 `json:"oi_cilindro,omitempty"`
	OIEje      *int     `json:"oi_eje,omitempty"`
	OIAdd      *float64 `json:"oi_add,omitempty"`

	// Tipo de lente indicado
	TipoLente     *string `json:"tipo_lente,omitempty"`
	TipoLenteOtro *string `json:"tipo_lente_otro,omitempty"`

	// Observaciones y profesional
	Observaciones      *string `json:"observaciones,omitempty"`
	RealizadoPorUserID *int64  `json:"realizado_por_user_id,omitempty"`
}

// UpdateExamenVisualRequest representa los datos para actualizar un examen visual
type UpdateExamenVisualRequest struct {
	FechaExamen *string `json:"fecha_examen,omitempty"` // formato: YYYY-MM-DD

	// Agudeza visual sin corrección
	AvScOD *string `json:"av_sc_od,omitempty"`
	AvScOI *string `json:"av_sc_oi,omitempty"`

	// Agudeza visual con corrección
	AvCcOD *string `json:"av_cc_od,omitempty"`
	AvCcOI *string `json:"av_cc_oi,omitempty"`

	// Refracción Ojo Derecho (OD)
	ODEsfera   *float64 `json:"od_esfera,omitempty"`
	ODCilindro *float64 `json:"od_cilindro,omitempty"`
	ODEje      *int     `json:"od_eje,omitempty"`
	ODAdd      *float64 `json:"od_add,omitempty"`

	// Refracción Ojo Izquierdo (OI)
	OIEsfera   *float64 `json:"oi_esfera,omitempty"`
	OICilindro *float64 `json:"oi_cilindro,omitempty"`
	OIEje      *int     `json:"oi_eje,omitempty"`
	OIAdd      *float64 `json:"oi_add,omitempty"`

	// Tipo de lente indicado
	TipoLente     *string `json:"tipo_lente,omitempty"`
	TipoLenteOtro *string `json:"tipo_lente_otro,omitempty"`

	// Observaciones y profesional
	Observaciones      *string `json:"observaciones,omitempty"`
	RealizadoPorUserID *int64  `json:"realizado_por_user_id,omitempty"`
}

// DiferenciaRefraccion representa la diferencia en dioptrías entre dos exámenes
type DiferenciaRefraccion struct {
	PacienteID     int64     `json:"paciente_id"`
	ExamenAnterior int64     `json:"examen_anterior_id"`
	ExamenActual   int64     `json:"examen_actual_id"`
	FechaAnterior  time.Time `json:"fecha_anterior"`
	FechaActual    time.Time `json:"fecha_actual"`

	// Diferencias OD
	DiferenciaODEsfera   *float64 `json:"diferencia_od_esfera,omitempty"`
	DiferenciaODCilindro *float64 `json:"diferencia_od_cilindro,omitempty"`

	// Diferencias OI
	DiferenciaOIEsfera   *float64 `json:"diferencia_oi_esfera,omitempty"`
	DiferenciaOICilindro *float64 `json:"diferencia_oi_cilindro,omitempty"`

	// Alertas
	AlertaCambioSignificativo bool   `json:"alerta_cambio_significativo"` // > 0.25D
	Mensaje                   string `json:"mensaje,omitempty"`
}
