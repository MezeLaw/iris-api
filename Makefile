.PHONY: docker-up docker-down migrate-up migrate-down setup

# Levantar docker compose
docker-up:
	docker-compose up -d
	@echo "Esperando a que PostgreSQL esté listo..."
	@timeout /t 5 /nobreak > nul
	@echo "PostgreSQL listo!"

# Bajar docker compose
docker-down:
	docker-compose down

# Ejecutar migraciones
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/iris?sslmode=disable" up

# Revertir migraciones
migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/iris?sslmode=disable" down

# Setup completo: levantar docker + migraciones
setup: docker-up migrate-up
	@echo "Proyecto listo!"

# Ejecutar la aplicación
run:
	go run cmd/api/main.go

# Setup y run
start: setup run
