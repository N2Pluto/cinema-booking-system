# ─── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

# Install git for private modules (if needed)
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies layer separately
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build: strip debug symbols + disable CGO for static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server .

# ─── Stage 2: Production ─────────────────────────────────────────────────────
FROM alpine:3.21

# TLS certificates + timezone data
RUN apk add --no-cache ca-certificates tzdata

# Non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .

# Run as non-root
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./server"]
