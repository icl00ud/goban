.PHONY: dev dev-frontend dev-backend build build-frontend build-backend build-cli docker docker-build docker-run clean release snapshot

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

build-cli:
	CGO_ENABLED=0 go build -o bin/goban ./cmd/goban

build: build-backend build-cli

# Release (requires goreleaser installed and a v* git tag)
snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

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
	@echo "  build-cli       - Build CLI binary (cmd/goban)"
	@echo "  snapshot        - Local GoReleaser snapshot build"
	@echo "  release         - Release via GoReleaser (requires tag)"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  clean           - Remove build artifacts"
