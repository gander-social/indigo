#!/bin/bash

# Test Sovereignty Implementation
echo "🇨🇦 Testing Gander Social Sovereignty Implementation"
echo "====================================================="

# Set up Go environment
export PATH=/usr/local/go/bin:$PATH
export GO111MODULE=on

# Function to run test and capture result
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo "🧪 Testing: $test_name"
    
    if eval "$test_command" >/dev/null 2>&1; then
        echo "✅ PASSED: $test_name"
        return 0
    else
        echo "❌ FAILED: $test_name"
        echo "Error output:"
        eval "$test_command" 2>&1 | head -n 10
        echo ""
        return 1
    fi
}

# Test 1: BGS Package Compilation
run_test "BGS Package Compilation" "go build ./bgs/"

# Test 2: Geographic Filter Tests
run_test "Geographic Filter Tests" "go test ./bgs/ -run TestGeographicFilter -v"

# Test 3: Sovereignty Config Tests
run_test "Sovereignty Config Tests" "go test ./bgs/ -run TestSovereignConfig -v"

# Test 4: Relay Package Compilation
run_test "Relay Package Compilation" "go build ./cmd/relay/"

# Test 5: Relay Sovereignty Tests
run_test "Relay Sovereignty Tests" "go test ./cmd/relay/ -run TestBGSSovereignty -v"

# Test 6: All BGS Tests
run_test "All BGS Tests" "go test ./bgs/ -v"

echo ""
echo "🇨🇦 Test Results Summary"
echo "========================="
echo "If all tests pass, the sovereignty implementation is ready for integration!"
