-- Remove from turnos
DROP INDEX IF EXISTS idx_turnos_client_id;
ALTER TABLE turnos DROP CONSTRAINT IF EXISTS fk_turnos_client;
ALTER TABLE turnos DROP COLUMN IF EXISTS client_id;

-- Remove from recetas
DROP INDEX IF EXISTS idx_recetas_client_id;
ALTER TABLE recetas DROP CONSTRAINT IF EXISTS fk_recetas_client;
ALTER TABLE recetas DROP COLUMN IF EXISTS client_id;

-- Remove from users
DROP INDEX IF EXISTS idx_users_client_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_client;
ALTER TABLE users DROP COLUMN IF EXISTS client_id;
