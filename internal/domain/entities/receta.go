package entities

import (
	"time"
)

type Receta struct {
	ID            int64     `json:"id" db:"id"`
	PacienteID    int64     `json:"paciente_id" db:"paciente_id"`
	Fecha         time.Time `json:"fecha" db:"fecha"`
	ODEsfera      float64   `json:"od_esfera" db:"od_esfera"`
	ODCilindro    float64   `json:"od_cilindro" db:"od_cilindro"`
	ODEje         int       `json:"od_eje" db:"od_eje"`
	OIEsfera      float64   `json:"oi_esfera" db:"oi_esfera"`
	OICilindro    float64   `json:"oi_cilindro" db:"oi_cilindro"`
	OIEje         int       `json:"oi_eje" db:"oi_eje"`
	TipoLente     string    `json:"tipo_lente" db:"tipo_lente"`
	Observaciones string    `json:"observaciones,omitempty" db:"observaciones"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRecetaRequest struct {
	PacienteID    int64     `json:"paciente_id" validate:"required"`
	Fecha         time.Time `json:"fecha" validate:"required"`
	ODEsfera      float64   `json:"od_esfera" validate:"required,min=-20,max=20"`
	ODCilindro    float64   `json:"od_cilindro" validate:"min=-10,max=10"`
	ODEje         int       `json:"od_eje" validate:"min=0,max=180"`
	OIEsfera      float64   `json:"oi_esfera" validate:"required,min=-20,max=20"`
	OICilindro    float64   `json:"oi_cilindro" validate:"min=-10,max=10"`
	OIEje         int       `json:"oi_eje" validate:"min=0,max=180"`
	TipoLente     string    `json:"tipo_lente" validate:"required,oneof=monofocal bifocal multifocal progresivo"`
	Observaciones string    `json:"observaciones,omitempty" validate:"max=500"`
}

type UpdateRecetaRequest struct {
	Fecha         *time.Time `json:"fecha,omitempty"`
	ODEsfera      *float64   `json:"od_esfera,omitempty" validate:"omitempty,min=-20,max=20"`
	ODCilindro    *float64   `json:"od_cilindro,omitempty" validate:"omitempty,min=-10,max=10"`
	ODEje         *int       `json:"od_eje,omitempty" validate:"omitempty,min=0,max=180"`
	OIEsfera      *float64   `json:"oi_esfera,omitempty" validate:"omitempty,min=-20,max=20"`
	OICilindro    *float64   `json:"oi_cilindro,omitempty" validate:"omitempty,min=-10,max=10"`
	OIEje         *int       `json:"oi_eje,omitempty" validate:"omitempty,min=0,max=180"`
	TipoLente     *string    `json:"tipo_lente,omitempty" validate:"omitempty,oneof=monofocal bifocal multifocal progresivo"`
	Observaciones *string    `json:"observaciones,omitempty" validate:"omitempty,max=500"`
}

type DioptriasChange struct {
	ODChange float64 `json:"od_change"`
	OIChange float64 `json:"oi_change"`
	HasAlert bool    `json:"has_alert"`
	Message  string  `json:"message"`
}
