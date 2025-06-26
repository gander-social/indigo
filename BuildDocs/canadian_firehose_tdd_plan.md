# AT Protocol Canadian Firehose TDD Implementation Plan

## Overview

This document outlines the Test Driven Development approach for implementing the Canadian Sovereign Firehose for Gander Social, ensuring AT Protocol compatibility while adding data sovereignty features.

## Phase 1: Foundation and First Test

### Step 1: Sovereignty Configuration Test

**Goal**: Verify that we can enable and configure sovereignty mode in the relay

**First Test**: Test sovereignty mode configuration

```go
// File: cmd/relay/sovereignty_test.go
package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSovereigntyConfigDefault(t *testing.T) {
	// Test that sovereignty is disabled by default
	config := DefaultServiceConfig()
	assert.False(t, config.SovereigntyEnabled, "Sovereignty should be disabled by default")
	assert.Equal(t, "", config.SovereignCountryCode, "Country code should be empty by default")
}

func TestSovereigntyConfigEnabled(t *testing.T) {
	// Test that sovereignty can be enabled with Canadian configuration
	config := &ServiceConfig{
		SovereigntyEnabled:   true,
		SovereignCountryCode: "CA",
		SovereigntyMode:      "strict",
	}
	
	assert.True(t, config.SovereigntyEnabled, "Sovereignty should be enabled")
	assert.Equal(t, "CA", config.SovereignCountryCode, "Country code should be CA")
	assert.Equal(t, "strict", config.SovereigntyMode, "Mode should be strict")
}

func TestSovereignFirehoseEndpoint(t *testing.T) {
	// Test that sovereignty mode creates dual endpoints
	config := &ServiceConfig{
		SovereigntyEnabled:   true,
		SovereignCountryCode: "CA",
	}
	
	// This test will fail initially - we need to implement the dual endpoint functionality
	expectedStandardEndpoint := "/xrpc/com.atproto.sync.subscribeRepos"
	expectedSovereignEndpoint := "/xrpc/ca.atproto.sync.subscribeRepos"
	
	// TODO: Implement endpoint detection logic
	endpoints := GetConfiguredEndpoints(config)
	
	assert.Contains(t, endpoints, expectedStandardEndpoint, "Standard endpoint should be available")
	if config.SovereigntyEnabled {
		assert.Contains(t, endpoints, expectedSovereignEndpoint, "Sovereign endpoint should be available when enabled")
	}
}
```

**Why this test first?**
- It's the simplest possible test that validates our core requirement
- It forces us to define the basic configuration structure
- It establishes the foundation for dual-mode operation
- It will fail initially, driving us to implement the minimal required functionality

### Implementation Required to Pass Test

To make this first test pass, we need to:

1. **Modify ServiceConfig struct** in `cmd/relay/service.go`:

```go
type ServiceConfig struct {
	// ... existing fields ...
	
	// Sovereignty configuration
	SovereigntyEnabled   bool
	SovereignCountryCode string
	SovereigntyMode      string // "strict", "balanced", "minimal"
	
	// ... rest of existing fields ...
}
```

2. **Update DefaultServiceConfig function**:

```go
func DefaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		ListenerBootTimeout:  5 * time.Second,
		SovereigntyEnabled:   false,  // Disabled by default
		SovereignCountryCode: "",
		SovereigntyMode:      "",
	}
}
```

3. **Create GetConfiguredEndpoints function** (minimal implementation):

```go
// File: cmd/relay/sovereignty.go
package main

func GetConfiguredEndpoints(config *ServiceConfig) []string {
	endpoints := []string{"/xrpc/com.atproto.sync.subscribeRepos"}
	
	if config.SovereigntyEnabled && config.SovereignCountryCode == "CA" {
		endpoints = append(endpoints, "/xrpc/ca.atproto.sync.subscribeRepos")
	}
	
	return endpoints
}
```

## Step 2: Dual Endpoint Handler Test

**Goal**: Verify that we can register and handle both standard and sovereign endpoints

**Test**: Verify dual endpoint registration

```go
func TestDualEndpointRegistration(t *testing.T) {
	// Setup test relay with sovereignty enabled
	config := &ServiceConfig{
		SovereigntyEnabled:   true,
		SovereignCountryCode: "CA",
		SovereigntyMode:      "strict",
	}
	
	// Create mock relay and service
	relay := createTestRelay(t)
	service, err := NewService(relay, config)
	require.NoError(t, err)
	
	// Create test echo instance
	e := echo.New()
	
	// Register handlers (this will need to be implemented)
	service.RegisterFirehoseHandlers(e)
	
	// Check that both endpoints are registered
	routes := getRegisteredRoutes(e)
	
	assert.Contains(t, routes, "GET /xrpc/com.atproto.sync.subscribeRepos")
	assert.Contains(t, routes, "GET /xrpc/ca.atproto.sync.subscribeRepos")
}

func TestSovereignEndpointAuthentication(t *testing.T) {
	// Test that sovereign endpoint requires authentication
	config := &ServiceConfig{
		SovereigntyEnabled:   true,
		SovereignCountryCode: "CA",
	}
	
	service := createTestService(t, config)
	
	// Test unauthenticated request to sovereign endpoint
	req := createTestFirehoseRequest("/xrpc/ca.atproto.sync.subscribeRepos", nil)
	resp := service.handleSovereignFirehose(req)
	
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	
	// Test authenticated request
	authReq := createTestFirehoseRequest("/xrpc/ca.atproto.sync.subscribeRepos", &AuthCredentials{
		Token: "valid-canadian-token",
	})
	authResp := service.handleSovereignFirehose(authReq)
	
	assert.Equal(t, http.StatusSwitchingProtocols, authResp.StatusCode) // WebSocket upgrade
}
```

## Step 3: Geographic Filter Test

**Goal**: Test the core filtering logic for Canadian content

**Test**: Basic geographic filtering

```go
func TestCanadianContentFilter(t *testing.T) {
	filter := NewGeographicFilter("CA")
	
	testCases := []struct {
		name     string
		event    *StreamEvent
		expected bool
	}{
		{
			name: "Canadian DID passes filter",
			event: &StreamEvent{
				Repo: "did:plc:canadian123",
				Kind: "commit",
			},
			expected: true,
		},
		{
			name: "US DID blocked by filter",
			event: &StreamEvent{
				Repo: "did:plc:american456",
				Kind: "commit",
			},
			expected: false,
		},
		{
			name: "Unknown DID defaults to blocked",
			event: &StreamEvent{
				Repo: "did:plc:unknown789",
				Kind: "commit",
			},
			expected: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.ShouldIncludeInSovereignFeed(tc.event)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGeographicFilterWithDatabase(t *testing.T) {
	// Test database-backed filtering
	db := setupTestDatabase(t)
	defer cleanupTestDatabase(db)
	
	// Insert test data
	insertCanadianDID(db, "did:plc:canadian123")
	insertNonCanadianDID(db, "did:plc:american456")
	
	filter := NewDatabaseGeographicFilter(db, "CA")
	
	// Test known Canadian DID
	canadianEvent := &StreamEvent{Repo: "did:plc:canadian123"}
	assert.True(t, filter.ShouldIncludeInSovereignFeed(canadianEvent))
	
	// Test known non-Canadian DID
	americanEvent := &StreamEvent{Repo: "did:plc:american456"}
	assert.False(t, filter.ShouldIncludeInSovereignFeed(americanEvent))
}
```

## Implementation Strategy

### Test Execution Order

1. **Run the first test** - it will fail
2. **Implement minimal code** to make it pass
3. **Run test again** - it should now pass
4. **Move to next test** - repeat cycle

### Key Principles

- **Red-Green-Refactor**: Write failing test → Make it pass → Clean up code
- **Minimal Implementation**: Only write enough code to make the current test pass
- **Incremental Complexity**: Each test adds one small piece of functionality
- **Continuous Integration**: All existing tests must continue to pass

### Development Environment Setup

Before starting, ensure:

```bash
# Clone the repository
git clone https://github.com/gander-social/indigo-sovereign
cd indigo-sovereign

# Create feature branch
git checkout -b feature/canadian-sovereign-firehose

# Run existing tests to ensure baseline
go test ./cmd/relay/...

# Create test file
touch cmd/relay/sovereignty_test.go
```

### Command to Run Tests

```bash
# Run sovereignty tests specifically
go test ./cmd/relay/ -v -run TestSovereignty

# Run all relay tests
go test ./cmd/relay/... -v

# Run with coverage
go test ./cmd/relay/ -v -cover -coverprofile=coverage.out
```

## Next Steps After First Test Passes

1. **Step 4**: Database schema tests for geographic data
2. **Step 5**: Member registry integration tests
3. **Step 6**: Filtering performance tests
4. **Step 7**: End-to-end firehose integration tests
5. **Step 8**: Authentication and authorization tests
6. **Step 9**: Privacy compliance tests
7. **Step 10**: Production deployment tests

Each step will be implemented using the same TDD approach, ensuring that:
- The AT Protocol core functionality continues to work
- New sovereign features integrate seamlessly
- All security and privacy requirements are met
- Performance requirements are satisfied

## Success Criteria for First Implementation

✅ **Test passes**: Basic sovereignty configuration works  
✅ **No regressions**: All existing relay tests continue to pass  
✅ **Clean code**: Implementation follows Go best practices  
✅ **Documentation**: Code is well-commented and clear  

Ready to begin implementation? Let's start with the first test!
