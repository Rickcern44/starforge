# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o api ./cmd/api


# Stage 2: Runner
FROM alpine:latest

WORKDIR /app

ENV PORT=8080
ENV APP_CONFIG_PATH=/app/config/settings.dev.toml

EXPOSE ${PORT}

COPY --from=builder /app/api ./api
COPY --from=builder /app/.config/settings.dev.toml ./.config/settings.dev.toml

ENTRYPOINT ["./api"]
