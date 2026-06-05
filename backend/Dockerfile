# ========================================
# Atlas Stage
# ========================================
FROM alpine:3.20 AS atlas
RUN apk --no-cache add curl && \
    curl -sSf https://atlasgo.sh | sh && \
    apk del curl

# ========================================
# Build Stage
# ========================================
FROM --platform=$BUILDPLATFORM golang:1.25.11-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy source code (including go.mod, go.sum)
COPY . .

# Build the application binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags='-w -s -extldflags "-static"' \
        -o /app/api-server \
        ./cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags='-w -s -extldflags "-static"' \
        -o /app/aggregate-episode-stats \
        ./cmd/aggregate-episode-stats && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags='-w -s -extldflags "-static"' \
        -o /app/write-robot-status-metrics \
        ./cmd/write-robot-status-metrics

# ========================================
# Runtime Stage
# ========================================
FROM --platform=linux/amd64 alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder /app/api-server .
COPY --from=builder /app/aggregate-episode-stats .
COPY --from=builder /app/write-robot-status-metrics .

# Copy Atlas CLI and migration files
COPY --from=atlas /usr/local/bin/atlas /usr/local/bin/atlas
COPY internal/database/migrate ./migrate
COPY atlas.hcl ./atlas.hcl

# Change ownership
RUN chown -R appuser:appuser /home/appuser

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8000

ENV PORT=8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8000/health || exit 1

# Run the application
CMD ["./api-server"]
