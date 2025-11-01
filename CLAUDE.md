# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go REST API project called "iris_api" that follows Clean Architecture principles. It uses the Gin web framework with PostgreSQL database and provides a scalable structure for multiple services.

## Project Structure

The project follows Clean Architecture with these layers:

```
iris-api/
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                 # Business logic layer
│   │   ├── entities/          # Business entities
│   │   ├── repositories/      # Repository interfaces
│   │   └── services/          # Domain services
│   ├── application/           # Application layer
│   │   └── usecases/         # Use cases (application logic)
│   ├── infrastructure/        # Infrastructure layer
│   │   ├── database/         # Database configuration
│   │   └── repositories/     # Repository implementations
│   └── presentation/          # Presentation layer
│       ├── handlers/         # HTTP handlers
│       ├── middleware/       # HTTP middlewares
│       └── routes/           # Route definitions
├── pkg/                       # Reusable packages
│   ├── config/               # Configuration management
│   ├── logger/               # Logging utilities
│   └── utils/                # General utilities
└── migrations/               # Database migrations
```

## Development Environment

- **Language**: Go 1.25.0
- **Framework**: Gin (web framework)
- **Database**: PostgreSQL with SQLx
- **Migration tool**: golang-migrate
- **IDE**: IntelliJ IDEA/GoLand

## Common Commands

### Development
```bash
# Run the application
go run cmd/api/main.go

# Build the application
go build -o bin/api cmd/api/main.go

# Run with environment file
cp .env.example .env
# Edit .env with your database credentials
go run cmd/api/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Tidy dependencies
go mod tidy
```

### Database Operations
```bash
# Create database
createdb iris_api

# Run migrations up
migrate -path migrations -database "postgres://user:password@localhost/iris_api?sslmode=disable" up

# Run migrations down
migrate -path migrations -database "postgres://user:password@localhost/iris_api?sslmode=disable" down

# Create new migration
migrate create -ext sql -dir migrations -seq migration_name
```

## Architecture Guidelines

When adding new services, follow this pattern:

1. **Domain Layer**: Create entities and repository interfaces
2. **Infrastructure Layer**: Implement repositories and external integrations
3. **Application Layer**: Define use cases that orchestrate business logic
4. **Presentation Layer**: Create handlers and routes for HTTP endpoints

### Dependencies Direction
- Domain layer should not depend on any other layer
- Application layer can depend on Domain
- Infrastructure layer can depend on Domain and Application
- Presentation layer can depend on all other layers

## Configuration

The application uses environment variables for configuration:
- Database connection settings (DB_HOST, DB_PORT, etc.)
- Server settings (SERVER_PORT, SERVER_HOST)
- Copy `.env.example` to `.env` and customize values

## Existing Services

### Authentication Service
JWT-based authentication with role-based access control:
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login (returns JWT token)
- `POST /api/v1/auth/validate` - Validate JWT token
- **Roles**: admin, optometrista, recepcionista
- **Multi-tenant**: Supports client_id isolation

### User Service
Complete CRUD operations for users with the following endpoints:
- `GET /api/v1/users` - List users with pagination
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Pacientes Service (Patient Management)
Comprehensive patient management with medical/visual history and refraction exams:

**Patients:**
- `POST /api/v1/pacientes` - Create patient (with nested antecedentes)
- `GET /api/v1/pacientes` - List with pagination
- `GET /api/v1/pacientes/search?q=query` - Search by name/DNI/email/phone
- `GET /api/v1/pacientes/:id?complete=true` - Get patient (optionally with full data)
- `PUT /api/v1/pacientes/:id` - Update patient
- `DELETE /api/v1/pacientes/:id` - Soft delete

**Medical History (Antecedentes Médicos):**
- `PUT /api/v1/pacientes/:id/antecedentes-medicos` - Create/Update medical history
- `GET /api/v1/pacientes/:id/antecedentes-medicos` - Get medical history

**Visual History (Antecedentes Visuales):**
- `PUT /api/v1/pacientes/:id/antecedentes-visuales` - Create/Update visual history
- `GET /api/v1/pacientes/:id/antecedentes-visuales` - Get visual history

**Visual Exams (Exámenes Visuales / Refracción):**
- `GET /api/v1/pacientes/:id/examenes` - List all exams for patient
- `POST /api/v1/examenes-visuales` - Create exam
- `GET /api/v1/examenes-visuales/:examenId` - Get exam
- `PUT /api/v1/examenes-visuales/:examenId` - Update exam
- `DELETE /api/v1/examenes-visuales/:examenId` - Delete exam
- `GET /api/v1/examenes-visuales/comparar?anterior=X&actual=Y` - Compare two exams (alerts for >0.25D changes)

**Features:**
- Multi-tenant (client_id isolation)
- Soft delete for patients
- Search and pagination
- Automatic age calculation
- Refraction comparison with alerts
- All routes protected with JWT authentication

### Reportería Service
Basic reporting endpoints:
- `GET /api/v1/reportes/pacientes-activos` - Patients with appointments in last 60 days
- `GET /api/v1/reportes/pacientes-inactivos` - Patients without appointments in 60+ days
- Protected routes (admin/optometrista only)

### Health Check
- `GET /health` - Server health status

## Database Schema

The project uses PostgreSQL with migration files in `/migrations/`:
- `001_create_users_table.up.sql` - Creates users table with indexes
- `002_create_clients_table.up.sql` - Creates clients table for multi-tenancy
- `003_add_role_to_users.up.sql` - Adds role column to users
- `004_create_recetas_table.up.sql` - Creates recetas (prescriptions) table
- `005_create_turnos_table.up.sql` - Creates turnos (appointments) table
- `006_add_client_id_to_existing_tables.up.sql` - Adds client_id for multi-tenancy
- `007_create_pacientes_tables.up.sql` - Creates pacientes, antecedentes_medicos, antecedentes_visuales, examenes_visuales tables

## Adding New Services

To add a new service (e.g., "products"):

1. Create `internal/domain/entities/product.go`
2. Create `internal/domain/repositories/product_repository.go` (interface)
3. Create `internal/domain/services/product_service.go`
4. Create `internal/application/usecases/product_usecase.go`
5. Create `internal/infrastructure/repositories/product_repository.go` (implementation)
6. Create `internal/presentation/handlers/product_handler.go`
7. Create `internal/presentation/routes/product_routes.go`
8. Create migration files in `/migrations/`
9. Wire dependencies in `cmd/api/main.go`

## Claude Code Configuration

- Docker exec commands are pre-approved in `.claude/settings.local.json`
- The project is on 'qa' branch with 'main' as the primary branch