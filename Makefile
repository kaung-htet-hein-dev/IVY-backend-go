.PHONY: dev prod build clean docker-dev docker-prod migrate seed air-dev air-prod

# Build the application
build:
	go build -o bin/app main.go

# Clean built binaries
clean:
	rm -rf bin/ tmp/

# Development server
dev:
	docker compose --env-file .env.development up -d db
	go run main.go dev

# Production server
prod:
	docker compose --env-file .env.production up -d db
	go run main.go prod

# Air development server with hot reload
air-dev:
	docker compose --env-file .env.development up -d db
	air -c .air.toml dev

# Air production server with hot reload
air-prod:
	docker compose --env-file .env.production up -d db
	GO_ENV=production air

# Run database migrations
migrate:
	go run cmd/main.go migrate

# Seed database
seed:
	go run cmd/seed/main.go

# Start development environment with Docker
docker-dev:
	docker compose --env-file .env.development up --build -d

# Start production environment with Docker
docker-prod:
	docker compose --env-file .env.production up --build -d

# Stop all containers
stop:
	docker compose down

# Show help
help:
	@echo "Usage:"
	@echo "  make dev         - Start development server with development database"
	@echo "  make prod        - Start production server with production database"
	@echo "  make air-dev     - Start development server with hot reload"
	@echo "  make air-prod    - Start production server with hot reload"
	@echo "  make build       - Build the application"
	@echo "  make clean       - Clean built binaries"
	@echo "  make migrate     - Run database migrations"
	@echo "  make seed        - Seed the database"
	@echo "  make docker-dev  - Start development environment with Docker"
	@echo "  make docker-prod - Start production environment with Docker"
	@echo "  make stop       - Stop all Docker containers"
	@echo "  make help       - Show this help message"
