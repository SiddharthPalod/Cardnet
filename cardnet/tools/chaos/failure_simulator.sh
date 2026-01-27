#!/bin/bash

# Phase 9: Failure Simulation Scripts
# Simulates various failure scenarios for resilience testing

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../.."

echo "🔴 Phase 9: Failure Simulation"
echo "================================"

# Function to kill a container
kill_container() {
    local container=$1
    echo "🛑 Stopping container: $container"
    docker stop "$container" 2>/dev/null || true
}

# Function to start a container
start_container() {
    local container=$1
    echo "✅ Starting container: $container"
    docker start "$container" 2>/dev/null || docker-compose up -d "$container" 2>/dev/null || true
}

# Function to wait for container health
wait_for_health() {
    local container=$1
    local max_wait=${2:-30}
    echo "⏳ Waiting for $container to be healthy (max ${max_wait}s)..."
    
    for i in $(seq 1 $max_wait); do
        if docker inspect "$container" --format='{{.State.Health.Status}}' 2>/dev/null | grep -q "healthy"; then
            echo "✅ $container is healthy"
            return 0
        fi
        sleep 1
    done
    
    echo "⚠️  $container health check timeout"
    return 1
}

# Scenario 1: Simulate Issuer Timeout
simulate_issuer_timeout() {
    echo ""
    echo "📋 Scenario 1: Simulating Issuer Timeout"
    echo "----------------------------------------"
    
    kill_container "cardnet-issuer-simulator"
    echo "Issuer simulator stopped. Auth gateway should handle timeouts gracefully."
    echo "Run load test to verify: go run tools/loadtest/load_test.go -duration 10s"
    
    read -p "Press Enter to restore issuer simulator..."
    start_container "cardnet-issuer-simulator"
    wait_for_health "cardnet-issuer-simulator" 15
}

# Scenario 2: Kill Cassandra Node
simulate_cassandra_failure() {
    echo ""
    echo "📋 Scenario 2: Simulating Cassandra Failure"
    echo "------------------------------------------"
    
    kill_container "cardnet-cassandra"
    echo "Cassandra stopped. Event ledger should handle failures gracefully."
    echo "Run load test to verify: go run tools/loadtest/load_test.go -duration 10s"
    
    read -p "Press Enter to restore Cassandra..."
    start_container "cardnet-cassandra"
    wait_for_health "cardnet-cassandra" 30
}

# Scenario 3: Kill Rate Limiter
simulate_rate_limiter_failure() {
    echo ""
    echo "📋 Scenario 3: Simulating Rate Limiter Failure"
    echo "----------------------------------------------"
    
    kill_container "cardnet-rate-limiter"
    echo "Rate limiter stopped. Auth gateway should handle circuit breaker."
    echo "Run load test to verify: go run tools/loadtest/load_test.go -duration 10s"
    
    read -p "Press Enter to restore rate limiter..."
    start_container "cardnet-rate-limiter"
    sleep 5
}

# Scenario 4: Kill Risk Engine
simulate_risk_engine_failure() {
    echo ""
    echo "📋 Scenario 4: Simulating Risk Engine Failure"
    echo "--------------------------------------------"
    
    kill_container "cardnet-risk-engine"
    echo "Risk engine stopped. Auth gateway should degrade gracefully."
    echo "Run load test to verify: go run tools/loadtest/load_test.go -duration 10s"
    
    read -p "Press Enter to restore risk engine..."
    start_container "cardnet-risk-engine"
    sleep 5
}

# Scenario 5: Kill Kafka
simulate_kafka_failure() {
    echo ""
    echo "📋 Scenario 5: Simulating Kafka Failure"
    echo "--------------------------------------"
    
    kill_container "cardnet-kafka"
    echo "Kafka stopped. Event publishing should handle backpressure."
    echo "Run load test to verify: go run tools/loadtest/load_test.go -duration 10s"
    
    read -p "Press Enter to restore Kafka..."
    start_container "cardnet-kafka"
    sleep 10
}

# Scenario 6: Kill Postgres (Partial)
simulate_postgres_failure() {
    echo ""
    echo "📋 Scenario 6: Simulating Postgres Failure"
    echo "------------------------------------------"
    echo "⚠️  WARNING: This will affect all services using Postgres!"
    
    read -p "Are you sure you want to continue? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        echo "Cancelled."
        return
    fi
    
    kill_container "cardnet-postgres"
    echo "Postgres stopped. Services should handle database unavailability."
    echo "Run load test to verify: go run tools/loadtest/load_test.go -duration 10s"
    
    read -p "Press Enter to restore Postgres..."
    start_container "cardnet-postgres"
    wait_for_health "cardnet-postgres" 30
}

# Main menu
show_menu() {
    echo ""
    echo "Select failure scenario to simulate:"
    echo "1) Issuer Timeout"
    echo "2) Cassandra Failure"
    echo "3) Rate Limiter Failure"
    echo "4) Risk Engine Failure"
    echo "5) Kafka Failure"
    echo "6) Postgres Failure (WARNING: affects all services)"
    echo "7) Run All Scenarios (with delays)"
    echo "8) Restore All Services"
    echo "0) Exit"
    echo ""
}

restore_all() {
    echo ""
    echo "🔄 Restoring All Services"
    echo "------------------------"
    
    start_container "cardnet-postgres"
    start_container "cardnet-cassandra"
    start_container "cardnet-kafka"
    start_container "cardnet-rate-limiter"
    start_container "cardnet-risk-engine"
    start_container "cardnet-issuer-simulator"
    start_container "cardnet-card-network"
    start_container "cardnet-event-ledger"
    start_container "cardnet-auth-gateway"
    
    echo "✅ All services restored. Waiting for health checks..."
    sleep 10
}

run_all_scenarios() {
    echo ""
    echo "🚀 Running All Failure Scenarios"
    echo "================================="
    echo "This will run each scenario with a 15-second test window."
    
    simulate_issuer_timeout
    sleep 5
    
    simulate_rate_limiter_failure
    sleep 5
    
    simulate_risk_engine_failure
    sleep 5
    
    simulate_kafka_failure
    sleep 5
    
    simulate_cassandra_failure
    sleep 5
    
    echo ""
    echo "✅ All scenarios completed. Restoring services..."
    restore_all
}

# Main loop
while true; do
    show_menu
    read -p "Enter choice: " choice
    
    case $choice in
        1) simulate_issuer_timeout ;;
        2) simulate_cassandra_failure ;;
        3) simulate_rate_limiter_failure ;;
        4) simulate_risk_engine_failure ;;
        5) simulate_kafka_failure ;;
        6) simulate_postgres_failure ;;
        7) run_all_scenarios ;;
        8) restore_all ;;
        0) echo "Exiting..."; exit 0 ;;
        *) echo "Invalid choice. Please try again." ;;
    esac
done
