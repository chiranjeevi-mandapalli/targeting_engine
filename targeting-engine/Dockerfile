FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o targeting-engine ./cmd/server
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache curl
COPY --from=builder /app/targeting-engine .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/healthz/ready || exit 1
CMD ["./targeting-engine"]