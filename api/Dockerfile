# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o api ./cmd/api


# Stage 2: Runner
FROM alpine:latest

# Install curl for healthcheck
RUN apk add --no-cache curl

WORKDIR /app

# Default environment variables
ENV PORT=3000
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_NAME=bouncy
ENV APP_ENV=production

EXPOSE ${PORT}

# Copy binary from builder
COPY --from=builder /app/api ./api

# Copy API documentation for Swagger/OpenAPI (if used)
COPY --from=builder /app/.docs ./.docs

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:${PORT}/api/v1/healthz || exit 1

ENTRYPOINT ["./api"]
