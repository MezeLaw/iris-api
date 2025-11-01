# TODO - MVP Iris API

## ✅ Módulos Completados

### 1. Sistema de Autenticación y Usuarios ✅
- [x] Registro y login con JWT
- [x] Roles (admin, optometrista, recepcionista)
- [x] Middleware de autenticación
- [x] Middleware de autorización por roles
- [x] Multi-tenant (client_id)

### 2. Módulo de Pacientes ✅ (PR #11 - Merged)
- [x] CRUD completo de pacientes
- [x] Antecedentes médicos (diabetes, hipertensión, alergias, cirugías)
- [x] Antecedentes visuales (lentes de contacto, marca, complicaciones)
- [x] Exámenes visuales/Refracciones (OD/OI, esfera, cilindro, eje, add)
- [x] Sistema de comparación de exámenes con alertas (>0.25D)
- [x] Búsqueda y paginación
- [x] Soft delete
- [x] Multi-tenant
- [x] 17 endpoints REST protegidos

### 3. Reportes Básicos ✅
- [x] Pacientes activos (con citas en últimos 60 días)
- [x] Pacientes inactivos (sin citas en 60+ días)

---

## 🚧 Módulos Pendientes del MVP

### 4. Módulo de Turnos/Agenda ⏳ (PRÓXIMO)
**Prioridad: ALTA - Core del MVP**

#### Base de Datos
- [ ] Migración para tabla `turnos`
  - Campos: id, paciente_id, profesional_user_id, tipo_servicio, fecha_hora, duracion_minutos, estado, observaciones
  - Estados: pendiente, confirmado, cancelado, completado, no_asistio
  - Relación con pacientes y users (profesionales)
  - Índices por fecha, paciente, profesional

#### Backend
- [ ] Entity: Turno
- [ ] Repository: TurnoRepository (CRUD, filtros por fecha/profesional/paciente)
- [ ] Service: TurnoService (validación de horarios, disponibilidad)
- [ ] UseCase: TurnoUseCase
- [ ] Handler: TurnoHandler
- [ ] Routes: `/api/v1/turnos/*`

#### Endpoints Necesarios
- [ ] `POST /api/v1/turnos` - Crear turno
- [ ] `GET /api/v1/turnos` - Listar turnos (con filtros: fecha_desde, fecha_hasta, profesional_id, paciente_id, estado)
- [ ] `GET /api/v1/turnos/:id` - Obtener turno por ID
- [ ] `PUT /api/v1/turnos/:id` - Actualizar turno
- [ ] `DELETE /api/v1/turnos/:id` - Cancelar turno
- [ ] `PATCH /api/v1/turnos/:id/estado` - Cambiar estado (confirmar, marcar no asistió, completar)
- [ ] `GET /api/v1/turnos/disponibilidad` - Verificar disponibilidad de horarios
- [ ] `GET /api/v1/turnos/por-dia/:fecha` - Vista por día
- [ ] `GET /api/v1/turnos/por-semana/:fecha_inicio` - Vista por semana
- [ ] `GET /api/v1/turnos/por-profesional/:user_id` - Turnos de un profesional

#### Reglas de Negocio
- [ ] Validar que no haya solapamiento de turnos para un mismo profesional
- [ ] Calcular hora_fin automáticamente (hora_inicio + duracion_minutos)
- [ ] Solo admin y optometrista pueden crear/modificar turnos
- [ ] Recepcionista puede ver turnos
- [ ] Notificaciones para turnos próximos (opcional)

---

### 5. Módulo de Ventas y Compras ⏳
**Prioridad: MEDIA**

#### Base de Datos
- [ ] Tabla `productos` (lentes, armazones, tratamientos)
  - Campos: id, client_id, nombre, tipo, marca, precio, stock, sku
- [ ] Tabla `ventas`
  - Campos: id, paciente_id, fecha_venta, total, metodo_pago, estado
- [ ] Tabla `items_venta`
  - Campos: id, venta_id, producto_id, cantidad, precio_unitario, subtotal
- [ ] Tabla `garantias`
  - Campos: id, venta_id, fecha_inicio, fecha_vencimiento, observaciones

#### Backend
- [ ] Entities: Producto, Venta, ItemVenta, Garantia
- [ ] Repositories para cada entidad
- [ ] Services con lógica de inventario
- [ ] UseCases
- [ ] Handlers
- [ ] Routes

#### Endpoints Necesarios
- [ ] CRUD de productos
- [ ] CRUD de ventas
- [ ] Asociar venta a paciente
- [ ] Gestión de garantías
- [ ] Alertas de renovación (12 meses)
- [ ] Reporte de productos más vendidos

---

### 6. Módulo de Comunicación (WhatsApp) ⏳
**Prioridad: BAJA - Nice to have**

#### Base de Datos
- [ ] Tabla `plantillas_mensajes`
  - Tipos: revision, seguimiento, promociones, cumpleaños
- [ ] Tabla `comunicaciones`
  - Historial de mensajes enviados

#### Backend
- [ ] Entities: PlantillaMensaje, Comunicacion
- [ ] Service para generar mensajes personalizados
- [ ] Integración con WhatsApp Web API (o manual)

#### Endpoints
- [ ] CRUD de plantillas
- [ ] Generar mensaje para paciente
- [ ] Historial de comunicaciones

---

### 7. Reportes Avanzados ⏳
**Prioridad: MEDIA**

#### Reportes Adicionales
- [ ] Evolución visual de un paciente (gráfico temporal de dioptrías)
- [ ] Pacientes sin control en 6/12 meses
- [ ] Productos más vendidos
- [ ] Tasa de retorno de pacientes
- [ ] Exportación a PDF/Excel

#### Backend
- [ ] Extender ReporteriaService con nuevos reportes
- [ ] Generación de PDFs (librería como `go-pdf`)
- [ ] Generación de Excel (librería como `excelize`)

---

### 8. Configuración Avanzada ⏳
**Prioridad: BAJA**

- [ ] Personalización por óptica (logo, nombre comercial, colores)
- [ ] Backup automático (local o cloud)
- [ ] Configuración de notificaciones
- [ ] Configuración de horarios de atención

---

## 📋 Orden de Implementación Sugerido

1. **Módulo de Turnos/Agenda** ← PRÓXIMO (crítico para MVP)
2. **Módulo de Ventas y Compras** (importante para operación completa)
3. **Reportes Avanzados** (mejora la utilidad)
4. **Comunicación WhatsApp** (nice to have)
5. **Configuración Avanzada** (opcional)

---

## 🎯 Notas Técnicas

### Testing
- Cada módulo nuevo debe incluir tests unitarios siguiendo el patrón table-driven
- Usar testify/mock para mocks
- Objetivo: >80% coverage

### Arquitectura
- Mantener Clean Architecture en todas las capas
- Seguir patrón existente: Domain → Infrastructure → Application → Presentation
- Multi-tenant en todas las tablas (client_id)
- Todas las rutas protegidas con autenticación JWT

### Base de Datos
- Usar migrations para todos los cambios
- Incluir índices para campos de búsqueda frecuente
- Soft delete donde tenga sentido
- Triggers para updated_at

### API Design
- RESTful endpoints
- Paginación en listados (page, page_size)
- Búsqueda con query param ?q=
- Filtros con query params
- Respuestas consistentes:
  ```json
  {
    "message": "Success message",
    "data": {...}
  }
  ```

---

## 🔗 Referencias

- **Repo**: https://github.com/MezeLaw/iris-api
- **Branch actual**: qa
- **Último PR**: #11 (Módulo Pacientes)
- **Documentación**: CLAUDE.md, ANALISIS_CODIGO.md

---

**Última actualización**: 2025-10-28
**Próximo paso**: Implementar Módulo de Turnos/Agenda
