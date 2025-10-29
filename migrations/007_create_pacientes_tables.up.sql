-- Tabla principal de pacientes
CREATE TABLE IF NOT EXISTS pacientes (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL,

    -- Información personal
    nombre_completo VARCHAR(255) NOT NULL,
    dni VARCHAR(50),
    fecha_nacimiento DATE NOT NULL,
    edad INT GENERATED ALWAYS AS (EXTRACT(YEAR FROM AGE(fecha_nacimiento))) STORED,
    genero VARCHAR(20) CHECK (genero IN ('Masculino', 'Femenino', 'Otro', 'Prefiero no decir')),

    -- Contacto
    telefono VARCHAR(50),
    email VARCHAR(255),
    direccion TEXT,

    -- Información clínica básica
    ocupacion VARCHAR(255),
    motivo_consulta TEXT,

    -- Auditoría
    fecha_primera_visita DATE DEFAULT CURRENT_DATE,
    observaciones TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Foreign keys
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

-- Tabla de antecedentes médicos (1:1 con paciente)
CREATE TABLE IF NOT EXISTS antecedentes_medicos (
    id BIGSERIAL PRIMARY KEY,
    paciente_id BIGINT NOT NULL UNIQUE,

    -- Enfermedades sistémicas (array de strings o JSON)
    tiene_diabetes BOOLEAN DEFAULT FALSE,
    tiene_hipertension BOOLEAN DEFAULT FALSE,
    tiene_alergias BOOLEAN DEFAULT FALSE,
    detalle_alergias TEXT,
    otras_enfermedades TEXT,

    -- Medicación
    medicacion_habitual TEXT,

    -- Cirugías
    cirugias_previas TEXT,
    cirugias_oculares TEXT,

    -- Auditoría
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    FOREIGN KEY (paciente_id) REFERENCES pacientes(id) ON DELETE CASCADE
);

-- Tabla de antecedentes visuales (1:1 con paciente)
CREATE TABLE IF NOT EXISTS antecedentes_visuales (
    id BIGSERIAL PRIMARY KEY,
    paciente_id BIGINT NOT NULL UNIQUE,

    -- Uso de lentes de contacto
    usa_lentes_contacto BOOLEAN DEFAULT FALSE,
    tipo_lente_actual VARCHAR(100),
    tiempo_uso_diario VARCHAR(50),
    marca_modelo VARCHAR(255),
    fecha_ultima_adaptacion DATE,
    molestias_complicaciones TEXT,

    -- Auditoría
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    FOREIGN KEY (paciente_id) REFERENCES pacientes(id) ON DELETE CASCADE
);

-- Tabla de exámenes visuales / refracciones (1:N con paciente - historial)
CREATE TABLE IF NOT EXISTS examenes_visuales (
    id BIGSERIAL PRIMARY KEY,
    paciente_id BIGINT NOT NULL,

    -- Fecha del examen
    fecha_examen DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Agudeza visual sin corrección
    av_sc_od VARCHAR(20), -- Ojo derecho
    av_sc_oi VARCHAR(20), -- Ojo izquierdo

    -- Agudeza visual con corrección
    av_cc_od VARCHAR(20),
    av_cc_oi VARCHAR(20),

    -- Refracción Ojo Derecho (OD)
    od_esfera DECIMAL(5,2),
    od_cilindro DECIMAL(5,2),
    od_eje INT,
    od_add DECIMAL(5,2),

    -- Refracción Ojo Izquierdo (OI)
    oi_esfera DECIMAL(5,2),
    oi_cilindro DECIMAL(5,2),
    oi_eje INT,
    oi_add DECIMAL(5,2),

    -- Tipo de lente indicado
    tipo_lente VARCHAR(50) CHECK (tipo_lente IN ('Monofocal', 'Bifocal', 'Progresivo', 'Lente de contacto', 'Ninguno', 'Otro')),
    tipo_lente_otro VARCHAR(255),

    -- Observaciones del examen
    observaciones TEXT,

    -- Profesional que realizó el examen
    realizado_por_user_id BIGINT,

    -- Auditoría
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    FOREIGN KEY (paciente_id) REFERENCES pacientes(id) ON DELETE CASCADE,
    FOREIGN KEY (realizado_por_user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Índices para mejorar performance
CREATE INDEX idx_pacientes_client_id ON pacientes(client_id);
CREATE INDEX idx_pacientes_dni ON pacientes(dni);
CREATE INDEX idx_pacientes_nombre ON pacientes(nombre_completo);
CREATE INDEX idx_pacientes_deleted_at ON pacientes(deleted_at);

CREATE INDEX idx_antecedentes_medicos_paciente ON antecedentes_medicos(paciente_id);
CREATE INDEX idx_antecedentes_visuales_paciente ON antecedentes_visuales(paciente_id);

CREATE INDEX idx_examenes_visuales_paciente ON examenes_visuales(paciente_id);
CREATE INDEX idx_examenes_visuales_fecha ON examenes_visuales(fecha_examen DESC);

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_pacientes_updated_at BEFORE UPDATE ON pacientes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_antecedentes_medicos_updated_at BEFORE UPDATE ON antecedentes_medicos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_antecedentes_visuales_updated_at BEFORE UPDATE ON antecedentes_visuales
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_examenes_visuales_updated_at BEFORE UPDATE ON examenes_visuales
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
