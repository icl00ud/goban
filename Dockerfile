# Stage 1: Build Frontend
FROM oven/bun:1 AS frontend-builder

WORKDIR /app/web

# Copy package files
COPY web/package.json web/bun.lock* ./

# Install dependencies
RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install --frozen-lockfile

# Copy source code
COPY web/ ./

# Build for production
RUN bun run build

# Stage 2: Build Backend
FROM golang:1.23-alpine AS backend-builder

# Install build dependencies (needed for SQLite CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Copy built frontend from stage 1
COPY --from=frontend-builder /app/web/dist ./web/dist

# Build with embedded files
# CGO_ENABLED=1 needed for SQLite
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o goban ./cmd/server

# Stage 3: Production Runtime
FROM alpine:3.19 AS production

# Install ca-certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN addgroup -S goban && adduser -S goban -G goban

# Copy binary from builder
COPY --from=backend-builder /app/goban .

# Create data directory for SQLite
RUN mkdir -p /app/data && chown -R goban:goban /app

USER goban

# Environment variables
ENV PORT=8080
ENV DB_DRIVER=sqlite
ENV DATABASE_URL=/app/data/goban.db

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

CMD ["./goban"]
