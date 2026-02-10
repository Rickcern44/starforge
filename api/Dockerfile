# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies - using go mod download for better caching
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for static linking, making the binary portable
# -o /app/api specifies the output path and name of the binary
# ./cmd/api specifies the main package
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/api ./cmd/api

# Stage 2: Runner
FROM alpine:latest

# Set environment variables
ENV PORT=8080

# Expose the port the application listens on
EXPOSE ${PORT}

# Set the working directory
WORKDIR /appc

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/api /app/api

# Run the application
ENTRYPOINT ["/app/api"]
