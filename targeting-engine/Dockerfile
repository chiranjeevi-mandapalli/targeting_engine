# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o targeting-engine ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Install curl for health checks
RUN apk add --no-cache curl

# Copy binaries and support files
COPY --from=builder /app/targeting-engine .
COPY --from=builder /app/migrations ./migrations

# Expose ports
EXPOSE 8080

# Health check configuration
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/healthz/ready || exit 1

# Application startup
CMD ["./targeting-engine"]