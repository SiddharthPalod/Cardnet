FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o /app/bin/event-ledger ./cmd/event-ledger

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/event-ledger .
CMD ["./event-ledger"]
