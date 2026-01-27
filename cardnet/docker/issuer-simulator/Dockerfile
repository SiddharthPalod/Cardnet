FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o /app/bin/issuer-simulator ./cmd/issuer-simulator

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/issuer-simulator .
# Copy config files if needed
COPY --from=builder /app/internal/issuer/config /app/internal/issuer/config
CMD ["./issuer-simulator"]
