FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o /app/bin/gateway-api ./cmd/gateway-api

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/gateway-api .
CMD ["./gateway-api"]
