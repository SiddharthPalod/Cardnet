FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o /app/bin/rate-limiter ./cmd/rate-limiter

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/rate-limiter .
CMD ["./rate-limiter"]
