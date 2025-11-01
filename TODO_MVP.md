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

### 4. Módulo de Turnos/Agenda ✅ (PR #12 - Merged)
- [x] Migración 008: Tabla turnos con soft delete
- [x] Entity: Turno con multi-tenancy
- [x] Repository: TurnoRepository con queries optimizadas
- [x] Service: TurnoService con validación de disponibilidad
- [x] UseCase: TurnoUseCase con paginación
- [x] Handler: TurnoHandler con 10 endpoints protegidos
- [x] Routes: `/api/v1/turnos/*` con JWT authentication
- [x] CRUD completo (crear, listar, obtener, actualizar, eliminar)
- [x] Cambio de estado (PATCH /api/v1/turnos/:id/estado)
- [x] Verificación de disponibilidad (POST /api/v1/turnos/disponibilidad)
- [x] Vistas por día, semana y profesional
- [x] Validación de solapamiento de turnos
- [x] Cálculo automático de hora_fin
- [x] Soft delete
- [x] Multi-tenant

---

## 🚧 Módulos Pendientes del MVP (Para retomar luego)

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

## 📋 Orden de Implementación Sugerido (Para futuras versiones)

1. **Módulo de Ventas y Compras** ← PRÓXIMO (importante para operación completa)
2. **Reportes Avanzados** (mejora la utilidad)
3. **Comunicación WhatsApp** (nice to have)
4. **Configuración Avanzada** (opcional)

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
- **Último PR**: #12 (Módulo Turnos/Agenda)
- **Documentación**: CLAUDE.md

---

## 🎉 Estado Actual del MVP

**MVP CORE COMPLETADO** ✅

El sistema actual ya cuenta con funcionalidad completa para:
- ✅ Gestión de usuarios multi-tenant con roles
- ✅ Gestión completa de pacientes con historial médico y visual
- ✅ Sistema de exámenes visuales con comparación y alertas
- ✅ Agenda y gestión de turnos con validación de disponibilidad
- ✅ Reportes básicos de actividad

**Total de endpoints REST**: ~35+
**Autenticación**: JWT con multi-tenancy
**Base de datos**: PostgreSQL con 8 migraciones

---

**Última actualización**: 2025-11-01
**Estado**: MVP Core funcional - Listo para testing y deploy
**Próximos pasos**: Módulo de Ventas y Compras (cuando se retome el desarrollo)
