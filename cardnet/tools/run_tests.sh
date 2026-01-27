#!/bin/bash

# Quick test runner for Phase 9
# Run all tests or individual test scenarios

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "🧪 Phase 9: Load, Failure & Hardening Tests"
echo "============================================="
echo ""

show_menu() {
    echo "Select test to run:"
    echo "1) Quick Load Test (10s, 10 workers)"
    echo "2) Standard Load Test (60s, 20 workers)"
    echo "3) High Load Test (60s, 50 workers)"
    echo "4) Backpressure Test (60s, 100 workers)"
    echo "5) Failure Simulator (Interactive)"
    echo "6) Run All Tests (Sequential)"
    echo "0) Exit"
    echo ""
}

run_quick_load() {
    echo "🚀 Running Quick Load Test..."
    go run tools/loadtest/load_test.go -concurrency 10 -rate 100 -duration 10s
}

run_standard_load() {
    echo "🚀 Running Standard Load Test..."
    go run tools/loadtest/load_test.go -concurrency 20 -rate 100 -duration 60s
}

run_high_load() {
    echo "🚀 Running High Load Test..."
    go run tools/loadtest/load_test.go -concurrency 50 -rate 200 -duration 60s
}

run_backpressure() {
    echo "🚀 Running Backpressure Test..."
    go run tools/chaos/backpressure_test.go -concurrency 100 -rate 1000 -duration 60s
}

run_failure_simulator() {
    echo "🚀 Starting Failure Simulator..."
    ./tools/chaos/failure_simulator.sh
}

run_all_tests() {
    echo "🚀 Running All Tests..."
    echo ""
    
    echo "1/4: Quick Load Test"
    run_quick_load
    sleep 5
    
    echo ""
    echo "2/4: Standard Load Test"
    run_standard_load
    sleep 5
    
    echo ""
    echo "3/4: High Load Test"
    run_high_load
    sleep 5
    
    echo ""
    echo "4/4: Backpressure Test"
    run_backpressure
    
    echo ""
    echo "✅ All tests completed!"
}

# Main loop
while true; do
    show_menu
    read -p "Enter choice: " choice
    
    case $choice in
        1) run_quick_load ;;
        2) run_standard_load ;;
        3) run_high_load ;;
        4) run_backpressure ;;
        5) run_failure_simulator ;;
        6) run_all_tests ;;
        0) echo "Exiting..."; exit 0 ;;
        *) echo "Invalid choice. Please try again." ;;
    esac
    
    echo ""
    read -p "Press Enter to continue..."
done
