package services

import (
	"context"
	"database/sql"
	"fmt"
	"iris-api/internal/domain/entities"
	"iris-api/internal/domain/repositories"
	"strings"
	"time"
)

type PacienteService interface {
	// CRUD Pacientes
	CreatePaciente(ctx context.Context, req *entities.CreatePacienteRequest) (*entities.Paciente, error)
	GetPacienteByID(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error)
	GetPacienteComplete(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error)
	ListPacientes(ctx context.Context, clientID int64, page, pageSize int) (*entities.PacienteListResponse, error)
	SearchPacientes(ctx context.Context, clientID int64, query string, page, pageSize int) (*entities.PacienteListResponse, error)
	UpdatePaciente(ctx context.Context, id int64, clientID int64, req *entities.UpdatePacienteRequest) (*entities.Paciente, error)
	DeletePaciente(ctx context.Context, id int64, clientID int64) error

	// Antecedentes Médicos
	CreateOrUpdateAntecedentesMedicos(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesMedicosRequest) (*entities.AntecedentesMedicos, error)
	GetAntecedentesMedicos(ctx context.Context, pacienteID int64) (*entities.AntecedentesMedicos, error)

	// Antecedentes Visuales
	CreateOrUpdateAntecedentesVisuales(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesVisualesRequest) (*entities.AntecedentesVisuales, error)
	GetAntecedentesVisuales(ctx context.Context, pacienteID int64) (*entities.AntecedentesVisuales, error)

	// Exámenes Visuales
	CreateExamenVisual(ctx context.Context, req *entities.CreateExamenVisualRequest) (*entities.ExamenVisual, error)
	GetExamenVisual(ctx context.Context, id int64) (*entities.ExamenVisual, error)
	GetExamenesVisualesByPaciente(ctx context.Context, pacienteID int64) ([]entities.ExamenVisual, error)
	UpdateExamenVisual(ctx context.Context, id int64, req *entities.UpdateExamenVisualRequest) (*entities.ExamenVisual, error)
	DeleteExamenVisual(ctx context.Context, id int64) error
	CompararExamenes(ctx context.Context, examenAnteriorID, examenActualID int64) (*entities.DiferenciaRefraccion, error)
}

type pacienteService struct {
	pacienteRepo             repositories.PacienteRepository
	antecedentesMedicosRepo  repositories.AntecedentesMedicosRepository
	antecedentesVisualesRepo repositories.AntecedentesVisualesRepository
	examenVisualRepo         repositories.ExamenVisualRepository
}

func NewPacienteService(
	pacienteRepo repositories.PacienteRepository,
	antecedentesMedicosRepo repositories.AntecedentesMedicosRepository,
	antecedentesVisualesRepo repositories.AntecedentesVisualesRepository,
	examenVisualRepo repositories.ExamenVisualRepository,
) PacienteService {
	return &pacienteService{
		pacienteRepo:             pacienteRepo,
		antecedentesMedicosRepo:  antecedentesMedicosRepo,
		antecedentesVisualesRepo: antecedentesVisualesRepo,
		examenVisualRepo:         examenVisualRepo,
	}
}

// CRUD Pacientes

func (s *pacienteService) CreatePaciente(ctx context.Context, req *entities.CreatePacienteRequest) (*entities.Paciente, error) {
	if err := s.validateCreatePacienteRequest(req); err != nil {
		return nil, err
	}

	// Parsear fecha de nacimiento
	fechaNacimiento, err := time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		return nil, fmt.Errorf("invalid fecha_nacimiento format, expected YYYY-MM-DD: %w", err)
	}

	// Parsear fecha primera visita (o usar hoy si no se proporciona)
	var fechaPrimeraVisita time.Time
	if req.FechaPrimeraVisita != nil && *req.FechaPrimeraVisita != "" {
		fechaPrimeraVisita, err = time.Parse("2006-01-02", *req.FechaPrimeraVisita)
		if err != nil {
			return nil, fmt.Errorf("invalid fecha_primera_visita format, expected YYYY-MM-DD: %w", err)
		}
	} else {
		fechaPrimeraVisita = time.Now()
	}

	paciente := &entities.Paciente{
		ClientID:           req.ClientID,
		NombreCompleto:     strings.TrimSpace(req.NombreCompleto),
		FechaNacimiento:    fechaNacimiento,
		FechaPrimeraVisita: fechaPrimeraVisita,
	}

	// Campos opcionales
	if req.DNI != nil {
		paciente.DNI = sql.NullString{String: strings.TrimSpace(*req.DNI), Valid: true}
	}
	if req.Genero != nil {
		paciente.Genero = sql.NullString{String: *req.Genero, Valid: true}
	}
	if req.Telefono != nil {
		paciente.Telefono = sql.NullString{String: strings.TrimSpace(*req.Telefono), Valid: true}
	}
	if req.Email != nil {
		paciente.Email = sql.NullString{String: strings.ToLower(strings.TrimSpace(*req.Email)), Valid: true}
	}
	if req.Direccion != nil {
		paciente.Direccion = sql.NullString{String: strings.TrimSpace(*req.Direccion), Valid: true}
	}
	if req.Ocupacion != nil {
		paciente.Ocupacion = sql.NullString{String: strings.TrimSpace(*req.Ocupacion), Valid: true}
	}
	if req.MotivoConsulta != nil {
		paciente.MotivoConsulta = sql.NullString{String: strings.TrimSpace(*req.MotivoConsulta), Valid: true}
	}
	if req.Observaciones != nil {
		paciente.Observaciones = sql.NullString{String: strings.TrimSpace(*req.Observaciones), Valid: true}
	}

	// Crear paciente
	err = s.pacienteRepo.Create(ctx, paciente)
	if err != nil {
		return nil, fmt.Errorf("error creating paciente: %w", err)
	}

	// Crear antecedentes médicos si se proporcionan
	if req.AntecedentesMedicos != nil {
		_, err = s.CreateOrUpdateAntecedentesMedicos(ctx, paciente.ID, req.AntecedentesMedicos)
		if err != nil {
			return nil, fmt.Errorf("error creating antecedentes medicos: %w", err)
		}
	}

	// Crear antecedentes visuales si se proporcionan
	if req.AntecedentesVisuales != nil {
		_, err = s.CreateOrUpdateAntecedentesVisuales(ctx, paciente.ID, req.AntecedentesVisuales)
		if err != nil {
			return nil, fmt.Errorf("error creating antecedentes visuales: %w", err)
		}
	}

	// Cargar paciente completo
	return s.GetPacienteComplete(ctx, paciente.ID, paciente.ClientID)
}

func (s *pacienteService) GetPacienteByID(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid paciente ID")
	}

	return s.pacienteRepo.GetByID(ctx, id, clientID)
}

func (s *pacienteService) GetPacienteComplete(ctx context.Context, id int64, clientID int64) (*entities.Paciente, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid paciente ID")
	}

	return s.pacienteRepo.GetComplete(ctx, id, clientID)
}

func (s *pacienteService) ListPacientes(ctx context.Context, clientID int64, page, pageSize int) (*entities.PacienteListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	pacientes, total, err := s.pacienteRepo.List(ctx, clientID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error listing pacientes: %w", err)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &entities.PacienteListResponse{
		Pacientes:  pacientes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *pacienteService) SearchPacientes(ctx context.Context, clientID int64, query string, page, pageSize int) (*entities.PacienteListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return s.ListPacientes(ctx, clientID, page, pageSize)
	}

	pacientes, total, err := s.pacienteRepo.Search(ctx, clientID, query, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error searching pacientes: %w", err)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &entities.PacienteListResponse{
		Pacientes:  pacientes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *pacienteService) UpdatePaciente(ctx context.Context, id int64, clientID int64, req *entities.UpdatePacienteRequest) (*entities.Paciente, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid paciente ID")
	}

	// Obtener paciente existente
	paciente, err := s.pacienteRepo.GetByID(ctx, id, clientID)
	if err != nil {
		return nil, fmt.Errorf("error getting paciente: %w", err)
	}

	// Actualizar campos proporcionados
	if req.NombreCompleto != nil {
		paciente.NombreCompleto = strings.TrimSpace(*req.NombreCompleto)
	}
	if req.DNI != nil {
		paciente.DNI = sql.NullString{String: strings.TrimSpace(*req.DNI), Valid: true}
	}
	if req.FechaNacimiento != nil {
		fechaNacimiento, err := time.Parse("2006-01-02", *req.FechaNacimiento)
		if err != nil {
			return nil, fmt.Errorf("invalid fecha_nacimiento format: %w", err)
		}
		paciente.FechaNacimiento = fechaNacimiento
	}
	if req.Genero != nil {
		paciente.Genero = sql.NullString{String: *req.Genero, Valid: true}
	}
	if req.Telefono != nil {
		paciente.Telefono = sql.NullString{String: strings.TrimSpace(*req.Telefono), Valid: true}
	}
	if req.Email != nil {
		paciente.Email = sql.NullString{String: strings.ToLower(strings.TrimSpace(*req.Email)), Valid: true}
	}
	if req.Direccion != nil {
		paciente.Direccion = sql.NullString{String: strings.TrimSpace(*req.Direccion), Valid: true}
	}
	if req.Ocupacion != nil {
		paciente.Ocupacion = sql.NullString{String: strings.TrimSpace(*req.Ocupacion), Valid: true}
	}
	if req.MotivoConsulta != nil {
		paciente.MotivoConsulta = sql.NullString{String: strings.TrimSpace(*req.MotivoConsulta), Valid: true}
	}
	if req.FechaPrimeraVisita != nil {
		fechaPrimeraVisita, err := time.Parse("2006-01-02", *req.FechaPrimeraVisita)
		if err != nil {
			return nil, fmt.Errorf("invalid fecha_primera_visita format: %w", err)
		}
		paciente.FechaPrimeraVisita = fechaPrimeraVisita
	}
	if req.Observaciones != nil {
		paciente.Observaciones = sql.NullString{String: strings.TrimSpace(*req.Observaciones), Valid: true}
	}

	err = s.pacienteRepo.Update(ctx, paciente)
	if err != nil {
		return nil, fmt.Errorf("error updating paciente: %w", err)
	}

	return paciente, nil
}

func (s *pacienteService) DeletePaciente(ctx context.Context, id int64, clientID int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid paciente ID")
	}

	return s.pacienteRepo.Delete(ctx, id, clientID)
}

// Antecedentes Médicos

func (s *pacienteService) CreateOrUpdateAntecedentesMedicos(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesMedicosRequest) (*entities.AntecedentesMedicos, error) {
	// Verificar si ya existen antecedentes
	existing, err := s.antecedentesMedicosRepo.GetByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing antecedentes medicos: %w", err)
	}

	antecedentes := &entities.AntecedentesMedicos{
		PacienteID:        pacienteID,
		TieneDiabetes:     req.TieneDiabetes,
		TieneHipertension: req.TieneHipertension,
		TieneAlergias:     req.TieneAlergias,
	}

	if req.DetalleAlergias != nil {
		antecedentes.DetalleAlergias = sql.NullString{String: strings.TrimSpace(*req.DetalleAlergias), Valid: true}
	}
	if req.OtrasEnfermedades != nil {
		antecedentes.OtrasEnfermedades = sql.NullString{String: strings.TrimSpace(*req.OtrasEnfermedades), Valid: true}
	}
	if req.MedicacionHabitual != nil {
		antecedentes.MedicacionHabitual = sql.NullString{String: strings.TrimSpace(*req.MedicacionHabitual), Valid: true}
	}
	if req.CirugiasPrevias != nil {
		antecedentes.CirugiasPrevias = sql.NullString{String: strings.TrimSpace(*req.CirugiasPrevias), Valid: true}
	}
	if req.CirugiasOculares != nil {
		antecedentes.CirugiasOculares = sql.NullString{String: strings.TrimSpace(*req.CirugiasOculares), Valid: true}
	}

	if existing != nil {
		// Actualizar
		err = s.antecedentesMedicosRepo.Update(ctx, antecedentes)
	} else {
		// Crear
		err = s.antecedentesMedicosRepo.Create(ctx, antecedentes)
	}

	if err != nil {
		return nil, err
	}

	return antecedentes, nil
}

func (s *pacienteService) GetAntecedentesMedicos(ctx context.Context, pacienteID int64) (*entities.AntecedentesMedicos, error) {
	return s.antecedentesMedicosRepo.GetByPacienteID(ctx, pacienteID)
}

// Antecedentes Visuales

func (s *pacienteService) CreateOrUpdateAntecedentesVisuales(ctx context.Context, pacienteID int64, req *entities.CreateAntecedentesVisualesRequest) (*entities.AntecedentesVisuales, error) {
	// Verificar si ya existen antecedentes
	existing, err := s.antecedentesVisualesRepo.GetByPacienteID(ctx, pacienteID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing antecedentes visuales: %w", err)
	}

	antecedentes := &entities.AntecedentesVisuales{
		PacienteID:        pacienteID,
		UsaLentesContacto: req.UsaLentesContacto,
	}

	if req.TipoLenteActual != nil {
		antecedentes.TipoLenteActual = sql.NullString{String: strings.TrimSpace(*req.TipoLenteActual), Valid: true}
	}
	if req.TiempoUsoDiario != nil {
		antecedentes.TiempoUsoDiario = sql.NullString{String: strings.TrimSpace(*req.TiempoUsoDiario), Valid: true}
	}
	if req.MarcaModelo != nil {
		antecedentes.MarcaModelo = sql.NullString{String: strings.TrimSpace(*req.MarcaModelo), Valid: true}
	}
	if req.FechaUltimaAdaptacion != nil && *req.FechaUltimaAdaptacion != "" {
		fecha, err := time.Parse("2006-01-02", *req.FechaUltimaAdaptacion)
		if err != nil {
			return nil, fmt.Errorf("invalid fecha_ultima_adaptacion format: %w", err)
		}
		antecedentes.FechaUltimaAdaptacion = sql.NullTime{Time: fecha, Valid: true}
	}
	if req.MolestiaComplicaciones != nil {
		antecedentes.MolestiaComplicaciones = sql.NullString{String: strings.TrimSpace(*req.MolestiaComplicaciones), Valid: true}
	}

	if existing != nil {
		// Actualizar
		err = s.antecedentesVisualesRepo.Update(ctx, antecedentes)
	} else {
		// Crear
		err = s.antecedentesVisualesRepo.Create(ctx, antecedentes)
	}

	if err != nil {
		return nil, err
	}

	return antecedentes, nil
}

func (s *pacienteService) GetAntecedentesVisuales(ctx context.Context, pacienteID int64) (*entities.AntecedentesVisuales, error) {
	return s.antecedentesVisualesRepo.GetByPacienteID(ctx, pacienteID)
}

// Exámenes Visuales

func (s *pacienteService) CreateExamenVisual(ctx context.Context, req *entities.CreateExamenVisualRequest) (*entities.ExamenVisual, error) {
	if req.PacienteID <= 0 {
		return nil, fmt.Errorf("invalid paciente ID")
	}

	// Parsear fecha del examen (o usar hoy si no se proporciona)
	var fechaExamen time.Time
	var err error
	if req.FechaExamen != nil && *req.FechaExamen != "" {
		fechaExamen, err = time.Parse("2006-01-02", *req.FechaExamen)
		if err != nil {
			return nil, fmt.Errorf("invalid fecha_examen format: %w", err)
		}
	} else {
		fechaExamen = time.Now()
	}

	examen := &entities.ExamenVisual{
		PacienteID:  req.PacienteID,
		FechaExamen: fechaExamen,
	}

	// Campos opcionales
	if req.AvScOD != nil {
		examen.AvScOD = sql.NullString{String: *req.AvScOD, Valid: true}
	}
	if req.AvScOI != nil {
		examen.AvScOI = sql.NullString{String: *req.AvScOI, Valid: true}
	}
	if req.AvCcOD != nil {
		examen.AvCcOD = sql.NullString{String: *req.AvCcOD, Valid: true}
	}
	if req.AvCcOI != nil {
		examen.AvCcOI = sql.NullString{String: *req.AvCcOI, Valid: true}
	}
	if req.ODEsfera != nil {
		examen.ODEsfera = sql.NullFloat64{Float64: *req.ODEsfera, Valid: true}
	}
	if req.ODCilindro != nil {
		examen.ODCilindro = sql.NullFloat64{Float64: *req.ODCilindro, Valid: true}
	}
	if req.ODEje != nil {
		examen.ODEje = sql.NullInt64{Int64: int64(*req.ODEje), Valid: true}
	}
	if req.ODAdd != nil {
		examen.ODAdd = sql.NullFloat64{Float64: *req.ODAdd, Valid: true}
	}
	if req.OIEsfera != nil {
		examen.OIEsfera = sql.NullFloat64{Float64: *req.OIEsfera, Valid: true}
	}
	if req.OICilindro != nil {
		examen.OICilindro = sql.NullFloat64{Float64: *req.OICilindro, Valid: true}
	}
	if req.OIEje != nil {
		examen.OIEje = sql.NullInt64{Int64: int64(*req.OIEje), Valid: true}
	}
	if req.OIAdd != nil {
		examen.OIAdd = sql.NullFloat64{Float64: *req.OIAdd, Valid: true}
	}
	if req.TipoLente != nil {
		examen.TipoLente = sql.NullString{String: *req.TipoLente, Valid: true}
	}
	if req.TipoLenteOtro != nil {
		examen.TipoLenteOtro = sql.NullString{String: *req.TipoLenteOtro, Valid: true}
	}
	if req.Observaciones != nil {
		examen.Observaciones = sql.NullString{String: strings.TrimSpace(*req.Observaciones), Valid: true}
	}
	if req.RealizadoPorUserID != nil {
		examen.RealizadoPorUserID = sql.NullInt64{Int64: *req.RealizadoPorUserID, Valid: true}
	}

	err = s.examenVisualRepo.Create(ctx, examen)
	if err != nil {
		return nil, fmt.Errorf("error creating examen visual: %w", err)
	}

	return examen, nil
}

func (s *pacienteService) GetExamenVisual(ctx context.Context, id int64) (*entities.ExamenVisual, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid examen visual ID")
	}

	return s.examenVisualRepo.GetByID(ctx, id)
}

func (s *pacienteService) GetExamenesVisualesByPaciente(ctx context.Context, pacienteID int64) ([]entities.ExamenVisual, error) {
	if pacienteID <= 0 {
		return nil, fmt.Errorf("invalid paciente ID")
	}

	return s.examenVisualRepo.ListByPacienteID(ctx, pacienteID)
}

func (s *pacienteService) UpdateExamenVisual(ctx context.Context, id int64, req *entities.UpdateExamenVisualRequest) (*entities.ExamenVisual, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid examen visual ID")
	}

	// Obtener examen existente
	examen, err := s.examenVisualRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting examen visual: %w", err)
	}

	// Actualizar campos proporcionados
	if req.FechaExamen != nil {
		fechaExamen, err := time.Parse("2006-01-02", *req.FechaExamen)
		if err != nil {
			return nil, fmt.Errorf("invalid fecha_examen format: %w", err)
		}
		examen.FechaExamen = fechaExamen
	}
	if req.AvScOD != nil {
		examen.AvScOD = sql.NullString{String: *req.AvScOD, Valid: true}
	}
	if req.AvScOI != nil {
		examen.AvScOI = sql.NullString{String: *req.AvScOI, Valid: true}
	}
	if req.AvCcOD != nil {
		examen.AvCcOD = sql.NullString{String: *req.AvCcOD, Valid: true}
	}
	if req.AvCcOI != nil {
		examen.AvCcOI = sql.NullString{String: *req.AvCcOI, Valid: true}
	}
	if req.ODEsfera != nil {
		examen.ODEsfera = sql.NullFloat64{Float64: *req.ODEsfera, Valid: true}
	}
	if req.ODCilindro != nil {
		examen.ODCilindro = sql.NullFloat64{Float64: *req.ODCilindro, Valid: true}
	}
	if req.ODEje != nil {
		examen.ODEje = sql.NullInt64{Int64: int64(*req.ODEje), Valid: true}
	}
	if req.ODAdd != nil {
		examen.ODAdd = sql.NullFloat64{Float64: *req.ODAdd, Valid: true}
	}
	if req.OIEsfera != nil {
		examen.OIEsfera = sql.NullFloat64{Float64: *req.OIEsfera, Valid: true}
	}
	if req.OICilindro != nil {
		examen.OICilindro = sql.NullFloat64{Float64: *req.OICilindro, Valid: true}
	}
	if req.OIEje != nil {
		examen.OIEje = sql.NullInt64{Int64: int64(*req.OIEje), Valid: true}
	}
	if req.OIAdd != nil {
		examen.OIAdd = sql.NullFloat64{Float64: *req.OIAdd, Valid: true}
	}
	if req.TipoLente != nil {
		examen.TipoLente = sql.NullString{String: *req.TipoLente, Valid: true}
	}
	if req.TipoLenteOtro != nil {
		examen.TipoLenteOtro = sql.NullString{String: *req.TipoLenteOtro, Valid: true}
	}
	if req.Observaciones != nil {
		examen.Observaciones = sql.NullString{String: strings.TrimSpace(*req.Observaciones), Valid: true}
	}
	if req.RealizadoPorUserID != nil {
		examen.RealizadoPorUserID = sql.NullInt64{Int64: *req.RealizadoPorUserID, Valid: true}
	}

	err = s.examenVisualRepo.Update(ctx, examen)
	if err != nil {
		return nil, fmt.Errorf("error updating examen visual: %w", err)
	}

	return examen, nil
}

func (s *pacienteService) DeleteExamenVisual(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid examen visual ID")
	}

	return s.examenVisualRepo.Delete(ctx, id)
}

func (s *pacienteService) CompararExamenes(ctx context.Context, examenAnteriorID, examenActualID int64) (*entities.DiferenciaRefraccion, error) {
	if examenAnteriorID <= 0 || examenActualID <= 0 {
		return nil, fmt.Errorf("invalid examen IDs")
	}

	return s.examenVisualRepo.GetDiferencia(ctx, examenAnteriorID, examenActualID)
}

// Validaciones

func (s *pacienteService) validateCreatePacienteRequest(req *entities.CreatePacienteRequest) error {
	if req.ClientID <= 0 {
		return fmt.Errorf("client_id is required")
	}
	if strings.TrimSpace(req.NombreCompleto) == "" {
		return fmt.Errorf("nombre_completo is required")
	}
	if len(req.NombreCompleto) < 2 {
		return fmt.Errorf("nombre_completo must be at least 2 characters long")
	}
	if req.FechaNacimiento == "" {
		return fmt.Errorf("fecha_nacimiento is required")
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email != "" && !strings.Contains(email, "@") {
			return fmt.Errorf("invalid email format")
		}
	}
	return nil
}
