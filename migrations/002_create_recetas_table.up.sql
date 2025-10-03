CREATE TABLE IF NOT EXISTS recetas (
    id BIGSERIAL PRIMARY KEY,
    paciente_id BIGINT NOT NULL,
    fecha TIMESTAMP NOT NULL,
    od_esfera DECIMAL(4,2) NOT NULL,
    od_cilindro DECIMAL(4,2) DEFAULT 0,
    od_eje INTEGER DEFAULT 0,
    oi_esfera DECIMAL(4,2) NOT NULL,
    oi_cilindro DECIMAL(4,2) DEFAULT 0,
    oi_eje INTEGER DEFAULT 0,
    tipo_lente VARCHAR(50) NOT NULL,
    observaciones TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_paciente FOREIGN KEY (paciente_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_recetas_paciente_id ON recetas(paciente_id);
CREATE INDEX idx_recetas_fecha ON recetas(fecha DESC);