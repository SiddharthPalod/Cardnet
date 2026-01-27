FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o /app/bin/card-network ./cmd/card-network

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/card-network .
# Copy issuer config file (required for BIN routing)
RUN mkdir -p ./internal/issuer/config
COPY --from=builder /app/internal/issuer/config/issuers.yaml ./internal/issuer/config/issuers.yaml
CMD ["./card-network"]
