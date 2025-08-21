# Iris API

API para la gestión de pacientes y seguimiento de citas en clínicas ópticas.

## Estructura del Proyecto

```
iris-api/
├── cmd/                    # Aplicaciones de línea de comandos
│   ├── create-patient/     # Crear pacientes
│   ├── get-patient/        # Obtener pacientes
│   ├── update-patient/     # Actualizar pacientes
│   └── delete-patient/     # Eliminar pacientes
├── internal/               # Código privado de la aplicación
│   ├── handlers/          # Lógica de negocio
│   ├── models/            # Modelos de datos
│   ├── repository/        # Capa de acceso a datos
│   └── utils/             # Utilidades (conexiones DB, etc.)
├── pkg/                   # Paquetes públicos
│   └── response/          # Helpers para respuestas HTTP
└── localstack-init/       # Scripts de inicialización LocalStack
```

## Desarrollo Local con DynamoDB

Este proyecto usa LocalStack para emular DynamoDB localmente durante el desarrollo.

### Prerequisitos
- Docker y Docker Compose
- Go 1.19+

### Configuración Inicial

1. **Iniciar LocalStack:**
   ```bash
   docker-compose up -d
   ```

2. **Verificar que DynamoDB está funcionando:**
   ```bash
   docker-compose ps
   # Debe mostrar localstack como "healthy"
   ```

3. **Ver logs (opcional):**
   ```bash
   docker-compose logs localstack
   ```

### Comandos Útiles DynamoDB Local

#### Gestión de Tablas
```bash
# Listar tablas
docker exec localstack awslocal dynamodb list-tables --region us-east-1

# Describir tabla Patients
docker exec localstack awslocal dynamodb describe-table --table-name Patients --region us-east-1

# Ver todos los elementos de la tabla
docker exec localstack awslocal dynamodb scan --table-name Patients --region us-east-1
```

#### Operaciones CRUD desde CLI

**PowerShell (Windows):**
```powershell
# Crear un paciente de prueba
docker exec localstack awslocal dynamodb put-item --table-name Patients --item "{\"PatientID\":{\"S\":\"P001\"},\"Name\":{\"S\":\"Juan Perez\"},\"Email\":{\"S\":\"juan@example.com\"}}" --region us-east-1

# Obtener un paciente específico
docker exec localstack awslocal dynamodb get-item --table-name Patients --key "{\"PatientID\":{\"S\":\"P001\"}}" --region us-east-1

# Eliminar un paciente
docker exec localstack awslocal dynamodb delete-item --table-name Patients --key "{\"PatientID\":{\"S\":\"P001\"}}" --region us-east-1
```

**Bash/Terminal (Linux/Mac):**
```bash
# Crear un paciente de prueba
docker exec localstack awslocal dynamodb put-item \
  --table-name Patients \
  --item '{"PatientID":{"S":"P001"},"Name":{"S":"Juan Perez"},"Email":{"S":"juan@example.com"}}' \
  --region us-east-1

# Obtener un paciente específico
docker exec localstack awslocal dynamodb get-item \
  --table-name Patients \
  --key '{"PatientID":{"S":"P001"}}' \
  --region us-east-1

# Eliminar un paciente
docker exec localstack awslocal dynamodb delete-item \
  --table-name Patients \
  --key '{"PatientID":{"S":"P001"}}' \
  --region us-east-1
```

### Conexión desde DataGrip/IDE

Para conectarte a DynamoDB local desde DataGrip:

1. **Crear nueva conexión:**
   - Data Source: `Amazon DynamoDB`
   - Authentication: `No Auth`
   - Region: `us-east-1`
   - Endpoint URL: `http://localhost:4566`

2. **Test Connection** debería ser exitoso

### Variables de Entorno

Para desarrollo local, configura estas variables:

```bash
# Archivo .env.localstack
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
AWS_DEFAULT_REGION=us-east-1
AWS_ENDPOINT_URL=http://localhost:4566
DYNAMODB_ENDPOINT=http://localhost:4566
PATIENTS_TABLE=Patients
IS_OFFLINE=true
```

### Comandos de Desarrollo

#### Build y Ejecutar
```bash
# Construir todas las aplicaciones
make build

# Construir aplicación específica
go build -o bin/create-patient cmd/create-patient/main.go

# Ejecutar aplicación
./bin/create-patient
```

#### Calidad de Código
```bash
# Formatear código
make fmt

# Verificar código
make vet

# Ejecutar tests
make test

# Ordenar dependencias
make tidy
```

### Gestión de LocalStack

#### Comandos Básicos
```bash
# Iniciar LocalStack
docker-compose up -d

# Parar LocalStack
docker-compose down

# Reiniciar (borra todos los datos)
docker-compose down && docker-compose up -d

# Ver estado de servicios
curl -s http://localhost:4566/_localstack/health
```

#### Reset de Datos
```bash
# Eliminar todos los datos y reiniciar
docker-compose down -v
docker-compose up -d
```

### Arquitectura de la Aplicación

- **Repository Pattern**: `PatientRepositoryInterface` abstrae las operaciones de DynamoDB
- **Dependency Injection**: Los handlers reciben instancias del repository para facilitar testing
- **Request/Response Models**: Estructuras separadas para requests de creación/actualización vs. modelo Patient
- **Configuración por Entorno**: Soporta tanto desarrollo local como entornos AWS

### Dependencias Principales

- `github.com/aws/aws-lambda-go` - Runtime de AWS Lambda para Go
- `github.com/aws/aws-sdk-go` - SDK de AWS para operaciones DynamoDB
- `github.com/google/uuid` - Generación de UUIDs para IDs de pacientes

## Producción

Para despliegue en AWS, las aplicaciones se configuran automáticamente para usar DynamoDB real estableciendo `IS_OFFLINE=false` y eliminando `AWS_ENDPOINT_URL`.