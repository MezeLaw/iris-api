-- Drop triggers
DROP TRIGGER IF EXISTS update_examenes_visuales_updated_at ON examenes_visuales;
DROP TRIGGER IF EXISTS update_antecedentes_visuales_updated_at ON antecedentes_visuales;
DROP TRIGGER IF EXISTS update_antecedentes_medicos_updated_at ON antecedentes_medicos;
DROP TRIGGER IF EXISTS update_pacientes_updated_at ON pacientes;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_examenes_visuales_fecha;
DROP INDEX IF EXISTS idx_examenes_visuales_paciente;
DROP INDEX IF EXISTS idx_antecedentes_visuales_paciente;
DROP INDEX IF EXISTS idx_antecedentes_medicos_paciente;
DROP INDEX IF EXISTS idx_pacientes_deleted_at;
DROP INDEX IF EXISTS idx_pacientes_nombre;
DROP INDEX IF EXISTS idx_pacientes_dni;
DROP INDEX IF EXISTS idx_pacientes_client_id;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS examenes_visuales;
DROP TABLE IF EXISTS antecedentes_visuales;
DROP TABLE IF EXISTS antecedentes_medicos;
DROP TABLE IF EXISTS pacientes;
