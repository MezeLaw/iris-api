# IRIS API - Postman Collection

Esta carpeta contiene la colección de Postman y los archivos de entorno para la API de IRIS.

## Estructura

```
postman/
├── iris-api.postman_collection.json  # Colección con todos los endpoints
└── environments/                      # Archivos de entorno
    ├── local.postman_environment.json
    ├── dev.postman_environment.json
    └── prod.postman_environment.json
```

## Importar en Postman

### Importar la Colección

1. Abre Postman
2. Haz clic en **Import** en la esquina superior izquierda
3. Selecciona el archivo `iris-api.postman_collection.json`
4. La colección aparecerá en tu workspace

### Importar los Entornos

1. En Postman, haz clic en el ícono de engranaje (⚙️) en la esquina superior derecha
2. Selecciona **Import**
3. Importa los archivos de entorno desde la carpeta `environments/`:
   - `local.postman_environment.json`
   - `dev.postman_environment.json`
   - `prod.postman_environment.json`

## Configuración de Entornos

### Local
- **base_url**: `http://localhost:8080`
- **auth_token**: (se establece automáticamente después de login/register)

### Dev
- **base_url**: `https://dev-api.iris.example.com`
- **auth_token**: (se establece automáticamente después de login/register)
- **Nota**: Actualiza la URL con tu dominio de desarrollo real

### Production
- **base_url**: `https://api.iris.example.com`
- **auth_token**: (se establece automáticamente después de login/register)
- **Nota**: Actualiza la URL con tu dominio de producción real

## Uso

### Seleccionar Entorno

1. En Postman, selecciona el entorno deseado desde el dropdown en la esquina superior derecha
2. Los requests utilizarán automáticamente las variables del entorno seleccionado

### Autenticación

Los endpoints protegidos requieren autenticación. Para obtener un token:

1. Ejecuta el request **Auth > Register** o **Auth > Login**
2. El token se guardará automáticamente en la variable `auth_token` del entorno
3. Los requests protegidos utilizarán automáticamente este token en el header `Authorization: Bearer {{auth_token}}`

### Endpoints Disponibles

#### Health
- `GET /health` - Health check

#### Auth
- `POST /api/v1/auth/register` - Registrar nuevo usuario (guarda token automáticamente)
- `POST /api/v1/auth/login` - Login (guarda token automáticamente)
- `GET /api/v1/auth/profile` - Obtener perfil del usuario autenticado (requiere auth)

#### Users
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users` - Listar usuarios (con paginación)
- `GET /api/v1/users/:id` - Obtener usuario por ID
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario

#### Recetas
- `POST /api/v1/recetas` - Crear receta
- `GET /api/v1/recetas` - Listar recetas (con paginación)
- `GET /api/v1/recetas/:id` - Obtener receta por ID
- `PUT /api/v1/recetas/:id` - Actualizar receta
- `DELETE /api/v1/recetas/:id` - Eliminar receta
- `GET /api/v1/pacientes/:paciente_id/recetas` - Obtener recetas de un paciente
- `GET /api/v1/pacientes/:paciente_id/recetas/historial` - Obtener historial de recetas
- `GET /api/v1/pacientes/:paciente_id/recetas/alertas` - Verificar cambios en dioptrías

#### Turnos
- `POST /api/v1/turnos` - Crear turno
- `GET /api/v1/turnos` - Listar turnos (con paginación)
- `GET /api/v1/turnos/:id` - Obtener turno por ID
- `PUT /api/v1/turnos/:id` - Actualizar turno
- `DELETE /api/v1/turnos/:id` - Eliminar turno
- `POST /api/v1/turnos/:id/cancel` - Cancelar turno
- `GET /api/v1/turnos/dia` - Obtener turnos por día
- `GET /api/v1/turnos/semana` - Obtener turnos por semana
- `GET /api/v1/turnos/alertas` - Obtener próximos turnos (alerta)
- `GET /api/v1/profesionales/:contactologo_id/turnos` - Obtener turnos de un profesional

#### Reportería
- `GET /api/v1/reportes/pacientes-activos` - Reporte de pacientes activos (requiere auth: admin/optometrista)
- `GET /api/v1/reportes/pacientes-inactivos` - Reporte de pacientes inactivos (requiere auth: admin/optometrista)

## Variables de Entorno

Cada entorno utiliza las siguientes variables:

- **base_url**: URL base de la API
- **auth_token**: Token JWT para autenticación (se establece automáticamente)

## Scripts Automáticos

La colección incluye scripts que automatizan ciertas tareas:

### Post-Request Scripts
- **Register/Login**: Guarda automáticamente el token JWT en la variable `auth_token` del entorno activo

## Notas

- Los requests de autenticación establecen automáticamente el token en el entorno
- Los requests protegidos utilizan automáticamente el token guardado
- Actualiza las URLs de los entornos Dev y Prod con tus dominios reales
- Los ejemplos de request bodies incluyen datos de muestra válidos
