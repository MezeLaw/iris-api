package repositories

import (
	"context"
	"iris-api/internal/domain/entities"
)

// PacienteRepository define las operaciones de persistencia para pacientes
type PacienteRepository interface {
	// CRUD básico
	Create(ctx context.Context, paciente *entities.Paciente) error
	GetByID(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error)
	Update(ctx context.Context, paciente *entities.Paciente) error
	Delete(ctx context.Context, id int64, clientID int64) error // Soft delete

	// Listado y búsqueda
	List(ctx context.Context, clientID int64, page, pageSize int) ([]entities.Paciente, int64, error)
	Search(ctx context.Context, clientID int64, query string, page, pageSize int) ([]entities.Paciente, int64, error)
	GetByDNI(ctx context.Context, dni string, clientID int64) (*entities.Paciente, error)

	// Relaciones
	GetWithAntecedentes(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error)
	GetWithExamenes(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error)
	GetComplete(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) // Con todo

	// Estadísticas
	CountByClientID(ctx context.Context, clientID int64) (int64, error)
}

// AntecedentesMedicosRepository define las operaciones para antecedentes médicos
type AntecedentesMedicosRepository interface {
	Create(ctx context.Context, antecedentes *entities.AntecedentesMedicos) error
	GetByPacienteID(ctx context.Context, pacienteID int64) (*entities.AntecedentesMedicos, error)
	Update(ctx context.Context, antecedentes *entities.AntecedentesMedicos) error
	Delete(ctx context.Context, pacienteID int64) error
}

// AntecedentesVisualesRepository define las operaciones para antecedentes visuales
type AntecedentesVisualesRepository interface {
	Create(ctx context.Context, antecedentes *entities.AntecedentesVisuales) error
	GetByPacienteID(ctx context.Context, pacienteID int64) (*entities.AntecedentesVisuales, error)
	Update(ctx context.Context, antecedentes *entities.AntecedentesVisuales) error
	Delete(ctx context.Context, pacienteID int64) error
}

// ExamenVisualRepository define las operaciones para exámenes visuales
type ExamenVisualRepository interface {
	Create(ctx context.Context, examen *entities.ExamenVisual) error
	GetByID(ctx context.Context, id int64) (*entities.ExamenVisual, error)
	Update(ctx context.Context, examen *entities.ExamenVisual) error
	Delete(ctx context.Context, id int64) error

	// Listado por paciente
	ListByPacienteID(ctx context.Context, pacienteID int64) ([]entities.ExamenVisual, error)
	GetLatestByPacienteID(ctx context.Context, pacienteID int64) (*entities.ExamenVisual, error)
	GetLastNExamenes(ctx context.Context, pacienteID int64, n int) ([]entities.ExamenVisual, error)

	// Comparación
	GetDiferencia(ctx context.Context, examenAnteriorID, examenActualID int64) (*entities.DiferenciaRefraccion, error)
}
