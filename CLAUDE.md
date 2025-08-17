# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project named `iris-api` - an application for managing patient appointments and tracking for optical clinics. The project follows Go standard folder structure and uses DynamoDB for data persistence.

## Architecture

The application consists of 4 main command-line applications for patient CRUD operations:
- **create-patient**: Creates new patient records
- **get-patient**: Retrieves patient information (single patient by ID or list with pagination)
- **update-patient**: Updates existing patient records
- **delete-patient**: Removes patient records

### Project Structure
```
iris-api/
├── cmd/                    # Command entry points (main functions)
│   ├── create-patient/
│   ├── get-patient/
│   ├── update-patient/
│   └── delete-patient/
├── internal/               # Private application code
│   ├── handlers/          # Business logic handlers
│   ├── models/            # Data models
│   ├── repository/        # Data access layer with interfaces
│   └── utils/             # Utilities (DB connections, etc.)
├── pkg/                   # Public packages
│   └── response/          # HTTP response helpers
└── bin/                   # Compiled binaries
```

## Development Commands

### Building Applications
```bash
# Build all applications
make build

# Build specific application
go build -o bin/create-patient cmd/create-patient/main.go
go build -o bin/get-patient cmd/get-patient/main.go
go build -o bin/update-patient cmd/update-patient/main.go
go build -o bin/delete-patient cmd/delete-patient/main.go

# Clean build artifacts
make clean
# or
rm -rf bin/
```

### Code Quality
```bash
# Format code
make fmt
# or
go fmt ./...

# Vet code for potential issues
make vet
# or
go vet ./...

# Run tests
make test
# or
go test ./...

# Tidy modules
make tidy
# or
go mod tidy
```

### Module Management
```bash
# Add dependencies
go get <package>

# Update dependencies
go mod tidy
```

## Environment Variables

- `PATIENTS_TABLE`: DynamoDB table name for patient records
- `AWS_REGION`: AWS region for DynamoDB operations
- `IS_OFFLINE`: Set to "true" for local development with local DynamoDB



## Dependencies

Key dependencies:
- `github.com/aws/aws-lambda-go` - AWS Lambda Go runtime (for API Gateway compatibility)
- `github.com/aws/aws-sdk-go` - AWS SDK for DynamoDB operations
- `github.com/google/uuid` - UUID generation for patient IDs

## Architecture Patterns

- **Repository Pattern**: `PatientRepositoryInterface` abstracts DynamoDB operations
- **Dependency Injection**: Handlers receive repository instances for testability
- **Request/Response Models**: Separate structs for create/update requests vs. Patient model
- **Environment Configuration**: Supports both local and AWS environments