# Load Testing Tool

Load testing tool for the CardNet authorization gateway.

## Usage

```bash
# Basic load test (10 workers, 100 req/s each, 30 seconds)
go run tools/loadtest/load_test.go

# High load test (50 workers, 200 req/s each, 60 seconds)
go run tools/loadtest/load_test.go -concurrency 50 -rate 200 -duration 60s

# Test against remote server
go run tools/loadtest/load_test.go -addr localhost:50051 -concurrency 20 -duration 2m
```

## Flags

- `-addr`: Auth gateway address (default: `localhost:50051`)
- `-concurrency`: Number of concurrent workers (default: `10`)
- `-duration`: Test duration (default: `30s`)
- `-rate`: Requests per second per worker (default: `100`)
- `-merchant`: Merchant ID to use (default: `merchant_001`)

## Example Output

```
============================================================
LOAD TEST RESULTS
============================================================
Duration:           30s
Total Requests:     30000
Successful:         29850 (99.50%)
Failed:             150
Rate Limited:       50
Timeouts:           100

Throughput:         1000.00 req/s
Avg Latency:        45.23 ms
Min Latency:        12.34 ms
Max Latency:        234.56 ms
============================================================

✅ Average latency (45.23 ms) meets SLA target (50 ms)
```

## Performance Targets

- **SLA**: <50ms average latency
- **Throughput**: 1000+ TPS
- **Success Rate**: >99%
