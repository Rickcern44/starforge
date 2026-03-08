# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git if needed for dependencies
RUN apk add --no-cache git

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with a unique name to avoid conflicts with folders/existing files
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o bouncy-api ./cmd/api


# Stage 2: Runner
FROM alpine:latest

# Install curl for healthcheck
RUN apk add --no-cache curl

WORKDIR /app

# Default environment variables
ENV PORT=3000
ENV DATABASE_URL=postgres://postgres:password@db:5432/bouncy?sslmode=disable
ENV APP_ENV=development

EXPOSE ${PORT}

# Copy binary from builder
COPY --from=builder /app/bouncy-api ./bouncy-api

# Copy API documentation for Swagger/OpenAPI (if used)
COPY --from=builder /app/.docs ./.docs

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:${PORT}/api/v1/healthz || exit 1

ENTRYPOINT ["./bouncy-api"]
