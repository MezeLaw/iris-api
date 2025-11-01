-- Create turnos table for appointment management
CREATE TABLE IF NOT EXISTS turnos (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL,
    paciente_id BIGINT NOT NULL,
    profesional_user_id BIGINT NOT NULL,
    tipo_servicio VARCHAR(100) NOT NULL,
    fecha_hora TIMESTAMP NOT NULL,
    duracion_minutos INTEGER NOT NULL DEFAULT 30,
    hora_fin TIMESTAMP NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    observaciones TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_turnos_paciente
        FOREIGN KEY (paciente_id)
        REFERENCES pacientes(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_turnos_profesional
        FOREIGN KEY (profesional_user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_turnos_client
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_estado
        CHECK (estado IN ('pendiente', 'confirmado', 'cancelado', 'completado', 'no_asistio')),

    CONSTRAINT chk_duracion_positiva
        CHECK (duracion_minutos > 0),

    CONSTRAINT chk_hora_fin_mayor
        CHECK (hora_fin > fecha_hora)
);

-- Indexes for faster queries
CREATE INDEX idx_turnos_client_id ON turnos(client_id);
CREATE INDEX idx_turnos_paciente_id ON turnos(paciente_id);
CREATE INDEX idx_turnos_profesional_user_id ON turnos(profesional_user_id);
CREATE INDEX idx_turnos_fecha_hora ON turnos(fecha_hora);
CREATE INDEX idx_turnos_estado ON turnos(estado);
CREATE INDEX idx_turnos_deleted_at ON turnos(deleted_at);

-- Composite index for date range queries
CREATE INDEX idx_turnos_profesional_fecha ON turnos(profesional_user_id, fecha_hora);

-- Composite index for patient appointments
CREATE INDEX idx_turnos_paciente_fecha ON turnos(paciente_id, fecha_hora);

-- Trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_turnos_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_turnos_updated_at
    BEFORE UPDATE ON turnos
    FOR EACH ROW
    EXECUTE FUNCTION update_turnos_updated_at();

-- Comments
COMMENT ON TABLE turnos IS 'Tabla de turnos/citas para los pacientes';
COMMENT ON COLUMN turnos.tipo_servicio IS 'Tipo de servicio: consulta, control, entrega_lentes, etc.';
COMMENT ON COLUMN turnos.estado IS 'Estado del turno: pendiente, confirmado, cancelado, completado, no_asistio';
COMMENT ON COLUMN turnos.duracion_minutos IS 'Duración estimada del turno en minutos';
COMMENT ON COLUMN turnos.hora_fin IS 'Hora calculada de finalización (fecha_hora + duracion_minutos)';
