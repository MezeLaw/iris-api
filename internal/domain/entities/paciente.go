package entities

import (
	"database/sql"
	"time"
)

// Paciente representa la entidad principal de un paciente en el sistema
type Paciente struct {
	ID       int64 `json:"id" db:"id"`
	ClientID int64 `json:"client_id" db:"client_id"`

	// Información personal
	NombreCompleto  string         `json:"nombre_completo" db:"nombre_completo"`
	DNI             sql.NullString `json:"dni,omitempty" db:"dni"`
	FechaNacimiento time.Time      `json:"fecha_nacimiento" db:"fecha_nacimiento"`
	Edad            int            `json:"edad" db:"edad"`
	Genero          sql.NullString `json:"genero,omitempty" db:"genero"`

	// Contacto
	Telefono  sql.NullString `json:"telefono,omitempty" db:"telefono"`
	Email     sql.NullString `json:"email,omitempty" db:"email"`
	Direccion sql.NullString `json:"direccion,omitempty" db:"direccion"`

	// Información clínica básica
	Ocupacion      sql.NullString `json:"ocupacion,omitempty" db:"ocupacion"`
	MotivoConsulta sql.NullString `json:"motivo_consulta,omitempty" db:"motivo_consulta"`

	// Auditoría
	FechaPrimeraVisita time.Time      `json:"fecha_primera_visita" db:"fecha_primera_visita"`
	Observaciones      sql.NullString `json:"observaciones,omitempty" db:"observaciones"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt          sql.NullTime   `json:"deleted_at,omitempty" db:"deleted_at"`

	// Relaciones (cargadas opcionalmente)
	AntecedentesMedicos  *AntecedentesMedicos  `json:"antecedentes_medicos,omitempty" db:"-"`
	AntecedentesVisuales *AntecedentesVisuales `json:"antecedentes_visuales,omitempty" db:"-"`
	ExamenesVisuales     []ExamenVisual        `json:"examenes_visuales,omitempty" db:"-"`
}

// CreatePacienteRequest representa los datos necesarios para crear un paciente
type CreatePacienteRequest struct {
	ClientID int64 `json:"client_id" binding:"required"`

	// Información personal
	NombreCompleto  string  `json:"nombre_completo" binding:"required"`
	DNI             *string `json:"dni,omitempty"`
	FechaNacimiento string  `json:"fecha_nacimiento" binding:"required"` // formato: YYYY-MM-DD
	Genero          *string `json:"genero,omitempty"`

	// Contacto
	Telefono  *string `json:"telefono,omitempty"`
	Email     *string `json:"email,omitempty"`
	Direccion *string `json:"direccion,omitempty"`

	// Información clínica básica
	Ocupacion          *string `json:"ocupacion,omitempty"`
	MotivoConsulta     *string `json:"motivo_consulta,omitempty"`
	FechaPrimeraVisita *string `json:"fecha_primera_visita,omitempty"` // formato: YYYY-MM-DD
	Observaciones      *string `json:"observaciones,omitempty"`

	// Datos anidados opcionales
	AntecedentesMedicos  *CreateAntecedentesMedicosRequest  `json:"antecedentes_medicos,omitempty"`
	AntecedentesVisuales *CreateAntecedentesVisualesRequest `json:"antecedentes_visuales,omitempty"`
}

// UpdatePacienteRequest representa los datos para actualizar un paciente
type UpdatePacienteRequest struct {
	// Información personal
	NombreCompleto  *string `json:"nombre_completo,omitempty"`
	DNI             *string `json:"dni,omitempty"`
	FechaNacimiento *string `json:"fecha_nacimiento,omitempty"` // formato: YYYY-MM-DD
	Genero          *string `json:"genero,omitempty"`

	// Contacto
	Telefono  *string `json:"telefono,omitempty"`
	Email     *string `json:"email,omitempty"`
	Direccion *string `json:"direccion,omitempty"`

	// Información clínica básica
	Ocupacion          *string `json:"ocupacion,omitempty"`
	MotivoConsulta     *string `json:"motivo_consulta,omitempty"`
	FechaPrimeraVisita *string `json:"fecha_primera_visita,omitempty"` // formato: YYYY-MM-DD
	Observaciones      *string `json:"observaciones,omitempty"`
}

// PacienteListResponse representa la respuesta paginada de pacientes
type PacienteListResponse struct {
	Pacientes  []Paciente `json:"pacientes"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}
