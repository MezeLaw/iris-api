-- Add client_id to users table (pacientes)
ALTER TABLE users ADD COLUMN client_id BIGINT;
ALTER TABLE users ADD CONSTRAINT fk_users_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE;
CREATE INDEX idx_users_client_id ON users(client_id);

-- Add client_id to recetas table
ALTER TABLE recetas ADD COLUMN client_id BIGINT;
ALTER TABLE recetas ADD CONSTRAINT fk_recetas_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE;
CREATE INDEX idx_recetas_client_id ON recetas(client_id);

-- Add client_id to turnos table
ALTER TABLE turnos ADD COLUMN client_id BIGINT;
ALTER TABLE turnos ADD CONSTRAINT fk_turnos_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE;
CREATE INDEX idx_turnos_client_id ON turnos(client_id);
