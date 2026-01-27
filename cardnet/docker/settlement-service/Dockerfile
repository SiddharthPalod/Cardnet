FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o /app/bin/settlement-service ./cmd/settlement-service

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/settlement-service .
# Copy schema file for initialization
COPY --from=builder /app/internal/settlement/repository/postgres/schema.sql ./internal/settlement/repository/postgres/
CMD ["./settlement-service"]
