# Frontend Development Prompt - Iris Optics Management System

## Project Context

You are tasked with building a **modern, responsive web application** for an optics management system called **Iris**. The backend API is already fully developed and deployed. Your job is to create a complete frontend application that consumes this REST API.

---

## Business Domain

**Iris** is a multi-tenant SaaS application for optical clinics (ópticas) that allows them to:
- Manage patients with complete medical and visual history
- Schedule and manage appointments (turnos/agenda)
- Track visual exams and refraction changes over time
- Generate reports on patient activity
- Manage multiple users with role-based access control

---

## Technical Stack Requirements

### Recommended Technologies
- **Framework**: React 18+ with TypeScript
- **State Management**: React Query (TanStack Query) for server state + Zustand for client state
- **Routing**: React Router v6
- **UI Framework**: Choose one:
  - Material-UI (MUI) v5 - Recommended for quick development
  - shadcn/ui + Tailwind CSS - Recommended for modern, customizable UI
  - Ant Design - Good for admin/enterprise apps
- **Forms**: React Hook Form + Zod for validation
- **HTTP Client**: Axios or fetch with interceptors
- **Date Handling**: date-fns or Day.js
- **Charts**: Recharts or Chart.js (for visual evolution graphs)
- **Build Tool**: Vite
- **Code Quality**: ESLint + Prettier

### Nice to Have
- TypeScript strict mode
- Component documentation with Storybook
- Unit tests with Vitest + React Testing Library
- E2E tests with Playwright (optional)

---

## Backend API Specification

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Flow

All endpoints (except auth endpoints) require JWT authentication via `Authorization: Bearer <token>` header.

#### Auth Endpoints (Public)

**1. Register**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "Juan",
  "lastname": "Pérez",
  "email": "juan@example.com",
  "password": "SecurePass123",
  "role": "admin", // admin | optometrista | recepcionista
  "client_name": "Óptica Central" // Required for first user
}

Response 201:
{
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "Juan",
      "lastname": "Pérez",
      "email": "juan@example.com",
      "role": "admin",
      "client_id": 1,
      "created_at": "2025-11-01T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "client": {
      "id": 1,
      "name": "Óptica Central"
    }
  }
}
```

**2. Login**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "juan@example.com",
  "password": "SecurePass123"
}

Response 200:
{
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "Juan",
      "lastname": "Pérez",
      "email": "juan@example.com",
      "role": "admin",
      "client_id": 1
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "client": {
      "id": 1,
      "name": "Óptica Central"
    }
  }
}
```

**3. Validate Token**
```http
POST /api/v1/auth/validate
Authorization: Bearer <token>

Response 200:
{
  "message": "Token is valid",
  "data": {
    "user_id": 1,
    "client_id": 1,
    "role": "admin",
    "email": "juan@example.com"
  }
}
```

---

### User Roles & Permissions

- **admin**: Full access to all features
- **optometrista**: Can manage patients, exams, appointments, view reports
- **recepcionista**: Can manage appointments, view patient info (read-only on medical data)

---

### Module 1: Patients (Pacientes)

#### Patient Entity Structure
```typescript
interface Paciente {
  id: number;
  client_id: number;
  nombre: string;
  apellido: string;
  dni: string;
  email: string;
  telefono: string;
  fecha_nacimiento: string; // ISO date
  edad: number; // Auto-calculated
  direccion: string;
  ciudad: string;
  provincia: string;
  codigo_postal: string;
  obra_social: string;
  numero_afiliado: string;
  observaciones: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

interface AntecedentesMedicos {
  id: number;
  paciente_id: number;
  diabetes: boolean;
  hipertension: boolean;
  alergias: string;
  cirugias_previas: string;
  medicacion_actual: string;
  otras_condiciones: string;
  created_at: string;
  updated_at: string;
}

interface AntecedentesVisuales {
  id: number;
  paciente_id: number;
  usa_lentes_contacto: boolean;
  marca_lentes_contacto: string;
  complicaciones_previas: string;
  observaciones: string;
  created_at: string;
  updated_at: string;
}

interface ExamenVisual {
  id: number;
  paciente_id: number;
  fecha_examen: string; // ISO date

  // Ojo derecho (OD)
  od_esfera: number;
  od_cilindro: number;
  od_eje: number;
  od_add: number;
  od_agudeza_visual: string;

  // Ojo izquierdo (OI)
  oi_esfera: number;
  oi_cilindro: number;
  oi_eje: number;
  oi_add: number;
  oi_agudeza_visual: string;

  observaciones: string;
  created_at: string;
  updated_at: string;
}

interface ComparacionExamenes {
  cambios_significativos: boolean;
  cambios: {
    od_esfera_diff: number;
    od_cilindro_diff: number;
    oi_esfera_diff: number;
    oi_cilindro_diff: number;
  };
  alerta: string; // Mensaje si cambio > 0.25D
}
```

#### Patient Endpoints

**1. Create Patient**
```http
POST /api/v1/pacientes
Authorization: Bearer <token>
Content-Type: application/json

{
  "nombre": "María",
  "apellido": "González",
  "dni": "12345678",
  "email": "maria@example.com",
  "telefono": "+54 11 1234-5678",
  "fecha_nacimiento": "1990-05-15",
  "direccion": "Av. Corrientes 1234",
  "ciudad": "CABA",
  "provincia": "Buenos Aires",
  "codigo_postal": "1043",
  "obra_social": "OSDE",
  "numero_afiliado": "123456",
  "observaciones": "Paciente regular",

  // Optional nested objects
  "antecedentes_medicos": {
    "diabetes": false,
    "hipertension": true,
    "alergias": "Polen",
    "medicacion_actual": "Enalapril 10mg"
  }
}

Response 201:
{
  "message": "Patient created successfully",
  "data": { ...paciente }
}
```

**2. List Patients (with pagination)**
```http
GET /api/v1/pacientes?page=1&page_size=10
Authorization: Bearer <token>

Response 200:
{
  "message": "Patients retrieved successfully",
  "data": {
    "pacientes": [...],
    "total": 45,
    "page": 1,
    "page_size": 10
  }
}
```

**3. Search Patients**
```http
GET /api/v1/pacientes/search?q=maria
Authorization: Bearer <token>

Response 200:
{
  "message": "Patients found",
  "data": [...pacientes]
}
```

**4. Get Patient by ID**
```http
GET /api/v1/pacientes/123?complete=true
Authorization: Bearer <token>

Response 200:
{
  "message": "Patient retrieved successfully",
  "data": {
    "paciente": {...},
    "antecedentes_medicos": {...},
    "antecedentes_visuales": {...},
    "examenes_visuales": [...]
  }
}
```

**5. Update Patient**
```http
PUT /api/v1/pacientes/123
Authorization: Bearer <token>
Content-Type: application/json

{
  "telefono": "+54 11 9999-8888",
  "email": "nuevo@example.com"
}

Response 200:
{
  "message": "Patient updated successfully",
  "data": {...paciente}
}
```

**6. Delete Patient (Soft Delete)**
```http
DELETE /api/v1/pacientes/123
Authorization: Bearer <token>

Response 200:
{
  "message": "Patient deleted successfully"
}
```

**7. Medical History**
```http
PUT /api/v1/pacientes/123/antecedentes-medicos
Authorization: Bearer <token>

{
  "diabetes": true,
  "hipertension": false,
  "alergias": "Ninguna",
  "cirugias_previas": "Cataratas 2020",
  "medicacion_actual": "Metformina 500mg",
  "otras_condiciones": ""
}

GET /api/v1/pacientes/123/antecedentes-medicos
```

**8. Visual History**
```http
PUT /api/v1/pacientes/123/antecedentes-visuales
GET /api/v1/pacientes/123/antecedentes-visuales
```

**9. Visual Exams**
```http
POST /api/v1/examenes-visuales
{
  "paciente_id": 123,
  "fecha_examen": "2025-11-01",
  "od_esfera": -2.50,
  "od_cilindro": -0.75,
  "od_eje": 180,
  "od_add": 0,
  "od_agudeza_visual": "20/20",
  "oi_esfera": -2.25,
  "oi_cilindro": -0.50,
  "oi_eje": 175,
  "oi_add": 0,
  "oi_agudeza_visual": "20/20",
  "observaciones": "Examen anual de rutina"
}

GET /api/v1/pacientes/123/examenes
GET /api/v1/examenes-visuales/456
PUT /api/v1/examenes-visuales/456
DELETE /api/v1/examenes-visuales/456

GET /api/v1/examenes-visuales/comparar?anterior=455&actual=456
Response: ComparacionExamenes
```

---

### Module 2: Appointments (Turnos/Agenda)

#### Turno Entity Structure
```typescript
interface Turno {
  id: number;
  client_id: number;
  paciente_id: number;
  profesional_user_id: number;
  tipo_servicio: string; // "consulta", "control", "entrega_lentes", etc.
  fecha_hora: string; // ISO datetime
  duracion_minutos: number;
  hora_fin: string; // Auto-calculated
  estado: EstadoTurno;
  observaciones: string;
  created_at: string;
  updated_at: string;
}

type EstadoTurno =
  | "pendiente"
  | "confirmado"
  | "cancelado"
  | "completado"
  | "no_asistio";

interface TurnoConDetalles extends Turno {
  paciente_nombre: string;
  paciente_apellido: string;
  paciente_email: string;
  profesional_nombre: string;
  profesional_apellido: string;
  profesional_email: string;
}

interface DisponibilidadResponse {
  disponible: boolean;
  mensaje: string;
  turnos_conflicto: Turno[];
}
```

#### Appointment Endpoints

**1. Create Appointment**
```http
POST /api/v1/turnos
Authorization: Bearer <token>

{
  "paciente_id": 123,
  "profesional_user_id": 2,
  "tipo_servicio": "consulta",
  "fecha_hora": "2025-11-05T10:00:00Z",
  "duracion_minutos": 30,
  "observaciones": "Primera consulta"
}

Response 201:
{
  "message": "Turno created successfully",
  "data": {...turno}
}
```

**2. List Appointments**
```http
GET /api/v1/turnos?page=1&page_size=20&fecha_desde=2025-11-01T00:00:00Z&fecha_hasta=2025-11-30T23:59:59Z&estado=pendiente&profesional_user_id=2
Authorization: Bearer <token>

Response 200:
{
  "message": "Turnos retrieved successfully",
  "data": {
    "turnos": [...],
    "total": 15,
    "page": 1,
    "page_size": 20
  }
}
```

**3. Get Appointment by ID**
```http
GET /api/v1/turnos/456
Authorization: Bearer <token>

Response 200:
{
  "message": "Turno retrieved successfully",
  "data": {...turnoConDetalles}
}
```

**4. Update Appointment**
```http
PUT /api/v1/turnos/456
Authorization: Bearer <token>

{
  "fecha_hora": "2025-11-05T14:00:00Z",
  "duracion_minutos": 45
}

Response 200:
{
  "message": "Turno updated successfully",
  "data": {...turno}
}
```

**5. Change Status**
```http
PATCH /api/v1/turnos/456/estado
Authorization: Bearer <token>

{
  "estado": "confirmado"
}

Response 200:
{
  "message": "Turno estado changed successfully"
}
```

**6. Delete Appointment**
```http
DELETE /api/v1/turnos/456
Authorization: Bearer <token>

Response 200:
{
  "message": "Turno deleted successfully"
}
```

**7. Check Availability**
```http
POST /api/v1/turnos/disponibilidad
Authorization: Bearer <token>

{
  "profesional_user_id": 2,
  "fecha_hora": "2025-11-05T10:00:00Z",
  "duracion_minutos": 30,
  "turno_id": null // Optional: exclude when editing
}

Response 200:
{
  "message": "Disponibilidad checked successfully",
  "data": {
    "disponible": true,
    "mensaje": "El profesional está disponible en ese horario"
  }
}
```

**8. Calendar Views**
```http
GET /api/v1/turnos/por-dia/2025-11-05?profesional_user_id=2
GET /api/v1/turnos/por-semana/2025-11-04?profesional_user_id=2
GET /api/v1/turnos/por-profesional/2?fecha_desde=2025-11-01&fecha_hasta=2025-11-30
```

---

### Module 3: Reports (Reportería)

```http
GET /api/v1/reportes/pacientes-activos
Authorization: Bearer <token>

Response 200:
{
  "message": "Report generated successfully",
  "data": [
    {
      "paciente_id": 123,
      "nombre": "María",
      "apellido": "González",
      "ultimo_turno": "2025-10-15T10:00:00Z",
      "cantidad_turnos": 3
    }
  ]
}

GET /api/v1/reportes/pacientes-inactivos
Response: Similar structure for inactive patients (no appointments in 60+ days)
```

---

### Module 4: Users (Usuarios)

```http
GET /api/v1/users?page=1&page_size=10
POST /api/v1/users
GET /api/v1/users/123
PUT /api/v1/users/123
DELETE /api/v1/users/123

User structure:
{
  "id": 1,
  "name": "Juan",
  "lastname": "Pérez",
  "email": "juan@example.com",
  "role": "admin",
  "client_id": 1,
  "created_at": "2025-11-01T10:00:00Z"
}
```

---

## UI/UX Requirements

### Overall Design Guidelines
- **Modern, clean interface** with good spacing and typography
- **Responsive design** - must work on desktop, tablet, and mobile
- **Accessibility** - WCAG 2.1 AA compliance
- **Color scheme**: Professional medical theme (blues, whites, soft grays)
- **Loading states** for all async operations
- **Error handling** with user-friendly messages
- **Success notifications** for CRUD operations

### Required Pages/Views

#### 1. Authentication
- [ ] Login page
- [ ] Register page (for first user/clinic setup)
- [ ] Forgot password (optional for MVP)
- [ ] Protected routes with role-based access

#### 2. Dashboard (Home)
- [ ] Welcome message with user name
- [ ] Quick stats cards:
  - Total patients
  - Appointments today
  - Pending appointments
  - Active patients (last 60 days)
- [ ] Today's appointments list
- [ ] Quick actions (New patient, New appointment)

#### 3. Patients Module
- [ ] **Patient List**
  - Search bar (by name, DNI, email, phone)
  - Table with columns: Name, DNI, Email, Phone, Age, Last visit
  - Pagination
  - Actions: View, Edit, Delete
  - "New Patient" button

- [ ] **Patient Details/Profile**
  - Tabs:
    - Personal Info (editable)
    - Medical History (editable)
    - Visual History (editable)
    - Visual Exams (list + create new)
    - Appointments History

- [ ] **New/Edit Patient Form**
  - Multi-step form or single form with sections
  - Form validation
  - Optional: Create medical history in same flow

- [ ] **Visual Exam Management**
  - Create exam form with OD/OI fields
  - List all exams with dates
  - Compare two exams feature (side-by-side comparison)
  - Alert badge if changes > 0.25D
  - Chart showing visual evolution over time

#### 4. Appointments Module (Calendar/Agenda)
- [ ] **Calendar View**
  - Choose view: FullCalendar or custom implementation
  - Views: Day, Week, Month
  - Filter by professional
  - Color-coded by status:
    - Pendiente: Yellow/Orange
    - Confirmado: Blue
    - Completado: Green
    - Cancelado: Red
    - No asistió: Gray
  - Click event to view details
  - Drag-and-drop to reschedule (optional)

- [ ] **New Appointment Form**
  - Select patient (searchable dropdown)
  - Select professional (dropdown)
  - Date and time picker
  - Duration (default 30 min)
  - Service type
  - Check availability before creating
  - Observations

- [ ] **Appointment Details Modal**
  - View all details
  - Actions: Edit, Change status, Cancel/Delete
  - Show patient and professional info

- [ ] **Appointment List View** (alternative to calendar)
  - Table with filters
  - Sort by date, patient, professional, status

#### 5. Reports Module
- [ ] **Active Patients Report**
  - List of patients with appointments in last 60 days
  - Export to PDF/Excel (optional)

- [ ] **Inactive Patients Report**
  - List of patients without appointments in 60+ days
  - Action: Schedule follow-up

#### 6. Users Module (Admin only)
- [ ] User list
- [ ] Create/Edit user
- [ ] Assign roles
- [ ] Deactivate user

#### 7. Settings (Optional for MVP)
- [ ] User profile
- [ ] Change password
- [ ] Clinic info

---

## Technical Implementation Guidelines

### Authentication & Authorization
```typescript
// Example auth hook
const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(
    localStorage.getItem('token')
  );

  const login = async (email: string, password: string) => {
    const response = await api.post('/auth/login', { email, password });
    const { token, user } = response.data.data;

    localStorage.setItem('token', token);
    setToken(token);
    setUser(user);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
  };

  return { user, token, login, logout, isAuthenticated: !!token };
};
```

### API Client Setup
```typescript
// axios-client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirect to login
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

### React Query Setup
```typescript
// queries/patients.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '../api/axios-client';

export const usePacientes = (page = 1, pageSize = 10) => {
  return useQuery({
    queryKey: ['pacientes', page, pageSize],
    queryFn: async () => {
      const { data } = await apiClient.get(
        `/pacientes?page=${page}&page_size=${pageSize}`
      );
      return data.data;
    },
  });
};

export const useCreatePaciente = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (paciente: CreatePacienteRequest) => {
      const { data } = await apiClient.post('/pacientes', paciente);
      return data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pacientes'] });
    },
  });
};
```

### Form Validation Example
```typescript
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const pacienteSchema = z.object({
  nombre: z.string().min(2, 'Nombre requerido'),
  apellido: z.string().min(2, 'Apellido requerido'),
  dni: z.string().min(7, 'DNI inválido'),
  email: z.string().email('Email inválido'),
  telefono: z.string().min(10, 'Teléfono inválido'),
  fecha_nacimiento: z.string(),
});

type PacienteFormData = z.infer<typeof pacienteSchema>;

const PacienteForm = () => {
  const { register, handleSubmit, formState: { errors } } = useForm<PacienteFormData>({
    resolver: zodResolver(pacienteSchema),
  });

  const onSubmit = (data: PacienteFormData) => {
    // Handle form submission
  };

  return <form onSubmit={handleSubmit(onSubmit)}>...</form>;
};
```

---

## Key Features to Implement

### Must Have (MVP)
- ✅ Authentication (login/register)
- ✅ Patient CRUD
- ✅ Medical & Visual history management
- ✅ Visual exam creation and comparison
- ✅ Appointment calendar/agenda
- ✅ Appointment status management
- ✅ Basic reports
- ✅ Responsive design
- ✅ Role-based access control

### Nice to Have (Post-MVP)
- ⭐ Drag-and-drop appointment rescheduling
- ⭐ Real-time notifications
- ⭐ Export to PDF/Excel
- ⭐ Dark mode
- ⭐ Multi-language support (ES/EN)
- ⭐ Appointment reminders via email/SMS
- ⭐ Visual evolution charts (line chart showing diopter changes over time)

---

## Multi-Tenancy Considerations

- All API calls automatically filtered by `client_id` (handled by backend)
- No need to pass `client_id` in requests (extracted from JWT token)
- Each clinic (client) sees only their own data

---

## Error Handling

### Common Error Responses
```json
// 400 Bad Request
{
  "error": "Invalid request body",
  "details": "Validation error: email is required"
}

// 401 Unauthorized
{
  "error": "Unauthorized",
  "details": "Invalid or expired token"
}

// 404 Not Found
{
  "error": "Resource not found",
  "details": "Patient with ID 999 not found"
}

// 500 Internal Server Error
{
  "error": "Internal server error",
  "details": "An unexpected error occurred"
}
```

### Frontend Error Handling Strategy
- Show toast notifications for errors
- Display inline validation errors on forms
- Provide retry mechanism for failed requests
- Graceful degradation for network issues

---

## Development Workflow

1. **Setup**
   - Initialize React + TypeScript project with Vite
   - Install dependencies
   - Configure ESLint + Prettier
   - Setup folder structure

2. **Project Structure**
   ```
   src/
   ├── api/
   │   ├── axios-client.ts
   │   └── endpoints/
   │       ├── auth.ts
   │       ├── patients.ts
   │       ├── appointments.ts
   │       └── reports.ts
   ├── components/
   │   ├── common/
   │   ├── patients/
   │   ├── appointments/
   │   └── layout/
   ├── pages/
   │   ├── Login.tsx
   │   ├── Dashboard.tsx
   │   ├── Patients/
   │   ├── Appointments/
   │   └── Reports/
   ├── hooks/
   │   ├── useAuth.ts
   │   └── queries/
   ├── types/
   │   └── api.ts
   ├── utils/
   ├── contexts/
   └── App.tsx
   ```

3. **Development Order**
   - Start with authentication
   - Implement patient module
   - Add appointments calendar
   - Add reports
   - Polish UI/UX

4. **Testing Strategy**
   - Unit tests for utilities and hooks
   - Integration tests for API calls
   - E2E tests for critical user flows

---

## Deliverables

1. Fully functional React application
2. TypeScript types for all API entities
3. Responsive design for desktop and mobile
4. README with setup instructions
5. Environment variables configuration (.env.example)
6. Build-ready for production deployment

---

## Questions to Clarify with Backend Team (Already Answered)

✅ What is the base URL? → `http://localhost:8080/api/v1`
✅ Are all endpoints paginated? → Yes, with `page` and `page_size`
✅ What date format is used? → ISO 8601 (e.g., "2025-11-01T10:00:00Z")
✅ Is multi-tenancy automatic? → Yes, via JWT token's `client_id`
✅ What are valid roles? → `admin`, `optometrista`, `recepcionista`

---

## Success Criteria

The frontend is considered complete when:
- ✅ All MVP features are implemented and working
- ✅ UI is responsive and accessible
- ✅ Authentication works correctly with role-based access
- ✅ All CRUD operations work as expected
- ✅ Error handling is user-friendly
- ✅ Code is well-organized and maintainable
- ✅ Application builds without errors
- ✅ Basic tests are passing

---

## Additional Resources

- **Backend Repository**: https://github.com/MezeLaw/iris-api
- **API Base URL**: http://localhost:8080
- **Health Check**: GET http://localhost:8080/health

---

## Notes

- The backend uses JWT tokens with expiration, so implement token refresh logic if needed
- All dates from backend are in ISO 8601 format
- The backend uses soft deletes, so deleted records have `deleted_at` field
- Pagination starts at page 1 (not 0)
- Search is case-insensitive and searches across multiple fields

---

Good luck building an amazing frontend for Iris! 🚀👓
