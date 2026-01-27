FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o /app/bin/analytics-service ./cmd/analytics-service

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/analytics-service .
COPY --from=builder /app/configs ./configs
CMD ["./analytics-service"]
