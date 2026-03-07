# =============================================================================
# ATLAS Blockchain - Multi-stage Docker Build
# =============================================================================

# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Install build dependencies (CGO for SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the blockchain binary with CGO enabled for SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o blockchain ./cmd

# =============================================================================
# Stage 2: Create minimal runtime image
# =============================================================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

# Copy the built binary from builder stage
COPY --from=builder /app/blockchain .

# Copy web frontend if exists
COPY --from=builder /app/web ./web

# Create data directory for persistence
RUN mkdir -p /data

# Environment variables
ENV API_PORT=8080
ENV DATABASE_PATH=/data/blockchain.db
ENV NODE_KEY_PATH=/data/nodekey.priv
ENV LOG_LEVEL=info

# Expose API port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/status || exit 1

# Run the blockchain node
CMD ["./blockchain", "-datadir", "/data", "-api", "8080", "-port", "8000", "-validator"]
