.PHONY: dev dev-frontend dev-backend build build-frontend build-backend docker docker-build docker-run clean

# Development
dev-frontend:
	cd web && bun install && bun run dev

dev-backend:
	go run ./cmd/server

dev: dev-backend

# Build
build-frontend:
	cd web && bun install && bun run build

build-backend: build-frontend
	go build -o bin/goban ./cmd/server

build: build-backend

# Docker
docker-build:
	docker build -t goban:latest .

docker-run:
	docker run -p 8080:8080 -v goban-data:/app/data -e JWT_SECRET=your-secret-key goban:latest

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

# Clean
clean:
	rm -rf bin/
	rm -rf web/dist/
	rm -rf web/node_modules/
	rm -f goban.db

# Help
help:
	@echo "Available targets:"
	@echo "  dev-frontend    - Run frontend dev server"
	@echo "  dev-backend     - Run backend dev server"
	@echo "  dev             - Run backend dev server (alias)"
	@echo "  build-frontend  - Build frontend for production"
	@echo "  build-backend   - Build backend with embedded frontend"
	@echo "  build           - Build everything"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  clean           - Remove build artifacts"
