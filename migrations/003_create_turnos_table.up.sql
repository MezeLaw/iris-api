CREATE TABLE IF NOT EXISTS turnos (
    id BIGSERIAL PRIMARY KEY,
    paciente_id BIGINT NOT NULL,
    contactologo_id BIGINT NOT NULL,
    fecha_hora TIMESTAMP NOT NULL,
    duracion_minutos INTEGER NOT NULL DEFAULT 30,
    tipo_servicio VARCHAR(50) NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    motivo TEXT,
    observaciones TEXT,
    recordatorio_enviado BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_paciente FOREIGN KEY (paciente_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_contactologo FOREIGN KEY (contactologo_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_duracion CHECK (duracion_minutos >= 15 AND duracion_minutos <= 240),
    CONSTRAINT chk_tipo_servicio CHECK (tipo_servicio IN ('consulta', 'control_vision', 'adaptacion_lentes', 'seguimiento', 'emergencia')),
    CONSTRAINT chk_estado CHECK (estado IN ('pendiente', 'confirmado', 'cancelado', 'completado', 'no_asistio'))
);

-- Índices para mejorar el rendimiento de las consultas
CREATE INDEX idx_turnos_paciente_id ON turnos(paciente_id);
CREATE INDEX idx_turnos_contactologo_id ON turnos(contactologo_id);
CREATE INDEX idx_turnos_fecha_hora ON turnos(fecha_hora DESC);
CREATE INDEX idx_turnos_estado ON turnos(estado);
CREATE INDEX idx_turnos_tipo_servicio ON turnos(tipo_servicio);

-- Índice compuesto para búsquedas por profesional y fecha
CREATE INDEX idx_turnos_contactologo_fecha ON turnos(contactologo_id, fecha_hora);

-- Índice para verificación de disponibilidad
CREATE INDEX idx_turnos_disponibilidad ON turnos(contactologo_id, fecha_hora, estado) WHERE estado NOT IN ('cancelado', 'no_asistio');
