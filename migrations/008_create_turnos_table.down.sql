-- Drop trigger
DROP TRIGGER IF EXISTS trigger_update_turnos_updated_at ON turnos;

-- Drop function
DROP FUNCTION IF EXISTS update_turnos_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_turnos_paciente_fecha;
DROP INDEX IF EXISTS idx_turnos_profesional_fecha;
DROP INDEX IF EXISTS idx_turnos_deleted_at;
DROP INDEX IF EXISTS idx_turnos_estado;
DROP INDEX IF EXISTS idx_turnos_fecha_hora;
DROP INDEX IF EXISTS idx_turnos_profesional_user_id;
DROP INDEX IF EXISTS idx_turnos_paciente_id;
DROP INDEX IF EXISTS idx_turnos_client_id;

-- Drop table
DROP TABLE IF EXISTS turnos;
