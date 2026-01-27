# Phase 9: Chaos Engineering & Failure Testing

This directory contains tools for testing system resilience, backpressure handling, and graceful degradation.

## Tools

### 1. Failure Simulator (`failure_simulator.sh`)

Interactive script to simulate various failure scenarios:

- **Issuer Timeout**: Stop issuer simulator to test timeout handling
- **Cassandra Failure**: Stop Cassandra to test event ledger resilience
- **Rate Limiter Failure**: Test circuit breaker behavior
- **Risk Engine Failure**: Test graceful degradation
- **Kafka Failure**: Test backpressure handling
- **Postgres Failure**: Test database unavailability handling

**Usage:**
```bash
chmod +x tools/chaos/failure_simulator.sh
./tools/chaos/failure_simulator.sh
```

### 2. Backpressure Test (`backpressure_test.go`)

Load test that sends requests at a rate higher than system capacity to verify:
- System doesn't crash under overload
- Rate limiting kicks in appropriately
- Latency increases gracefully
- Error handling is proper

**Usage:**
```bash
# High load backpressure test
go run tools/chaos/backpressure_test.go -concurrency 100 -rate 1000 -duration 60s
```

## Test Scenarios

### Scenario 1: Issuer Timeout
**Objective**: Verify auth gateway handles issuer timeouts gracefully

**Steps**:
1. Run failure simulator, select option 1
2. Run load test: `go run tools/loadtest/load_test.go -duration 10s`
3. Verify: Requests should timeout or be rejected gracefully, not crash

**Expected**: System continues operating, returns appropriate error codes

### Scenario 2: Cassandra Failure
**Objective**: Verify event ledger handles Cassandra unavailability

**Steps**:
1. Run failure simulator, select option 2
2. Run load test
3. Verify: Events are queued or handled gracefully

**Expected**: System continues, events may be buffered or logged

### Scenario 3: Rate Limiter Failure
**Objective**: Verify circuit breaker prevents cascade failures

**Steps**:
1. Run failure simulator, select option 3
2. Run load test
3. Verify: Circuit breaker opens, requests fail fast

**Expected**: Fast failure, no hanging requests

### Scenario 4: Risk Engine Failure
**Objective**: Verify graceful degradation when risk engine is unavailable

**Steps**:
1. Run failure simulator, select option 4
2. Run load test
3. Verify: System continues with degraded risk checking

**Expected**: Requests processed with default risk behavior

### Scenario 5: Kafka Failure
**Objective**: Verify backpressure when event queue is unavailable

**Steps**:
1. Run failure simulator, select option 5
2. Run load test
3. Verify: System handles queue unavailability

**Expected**: Events buffered or dropped gracefully

### Scenario 6: Backpressure Test
**Objective**: Verify system handles overload gracefully

**Steps**:
```bash
go run tools/chaos/backpressure_test.go -concurrency 100 -rate 1000 -duration 60s
```

**Expected**:
- Success rate > 95%
- Latency increases gradually
- Rate limiting activates
- No crashes or hangs

## Success Criteria

✅ **Load Test**: Average latency < 50ms at 1000 TPS
✅ **Backpressure**: System handles 10x normal load without crashing
✅ **Failure Recovery**: System recovers within 30s after service restoration
✅ **Graceful Degradation**: System continues operating with partial failures
✅ **Circuit Breakers**: Fast failure when dependencies are down

## Monitoring

While running tests, monitor:
- Service logs: `docker logs -f cardnet-auth-gateway`
- Database connections: `docker exec cardnet-postgres psql -U cardnet -c "SELECT count(*) FROM pg_stat_activity;"`
- Kafka lag: Check consumer group lag
- System resources: `docker stats`

## Cleanup

After testing, restore all services:
```bash
./tools/chaos/failure_simulator.sh
# Select option 8: Restore All Services
```

Or manually:
```bash
docker-compose restart
```
