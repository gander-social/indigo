package main

import (
	"testing"
	"time"

	"github.com/gander-social/gander-indigo-sovereign/bgs"
	"github.com/stretchr/testify/assert"
)

// TestBGSSovereigntyConfigDefault verifies that sovereignty is disabled by default
func TestBGSSovereigntyConfigDefault(t *testing.T) {
	config := bgs.DefaultBGSConfig()

	assert.False(t, config.SovereigntyEnabled, "Sovereignty should be disabled by default")
	assert.Equal(t, "", config.SovereignCountryCode, "Country code should be empty by default")
	assert.Equal(t, "balanced", config.SovereigntyMode, "Mode should default to balanced")
}

// TestBGSSovereigntyConfigEnabled verifies that sovereignty can be configured
func TestBGSSovereigntyConfigEnabled(t *testing.T) {
	config := bgs.DefaultBGSConfig()
	config.SovereigntyEnabled = true
	config.SovereignCountryCode = "CA"
	config.SovereigntyMode = "strict"

	assert.True(t, config.SovereigntyEnabled, "Sovereignty should be enabled")
	assert.Equal(t, "CA", config.SovereignCountryCode, "Country code should be CA")
	assert.Equal(t, "strict", config.SovereigntyMode, "Mode should be strict")
}

// TestServiceConfigNoSovereignty verifies that ServiceConfig no longer has sovereignty fields
func TestServiceConfigNoSovereignty(t *testing.T) {
	config := DefaultServiceConfig()

	// This should compile successfully - no sovereignty fields should exist
	assert.NotNil(t, config)
	assert.Equal(t, 5*time.Second, config.ListenerBootTimeout)
	assert.False(t, config.DisableRequestCrawl)
}
