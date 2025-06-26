package bgs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSovereignConfig_Validate(t *testing.T) {
	tests := []struct {
		name     string
		config   SovereignConfig
		valid    bool
		expected string
	}{
		{
			name: "valid canadian config",
			config: SovereignConfig{
				Enabled:           true,
				CountryCode:       "CA",
				DataRetentionDays: 30,
				PrivacyMode:       "strict",
			},
			valid: true,
		},
		{
			name: "valid us config",
			config: SovereignConfig{
				Enabled:           true,
				CountryCode:       "US",
				DataRetentionDays: 90,
				PrivacyMode:       "standard",
			},
			valid: true,
		},
		{
			name: "invalid country code",
			config: SovereignConfig{
				Enabled:     true,
				CountryCode: "XX",
			},
			valid: false,
		},
		{
			name: "empty country code",
			config: SovereignConfig{
				Enabled:     true,
				CountryCode: "",
			},
			valid: false,
		},
		{
			name: "too long country code",
			config: SovereignConfig{
				Enabled:     true,
				CountryCode: "CAN",
			},
			valid: false,
		},
		{
			name: "disabled sovereign mode",
			config: SovereignConfig{
				Enabled: false,
			},
			valid: true,
		},
		{
			name: "disabled with invalid country code",
			config: SovereignConfig{
				Enabled:     false,
				CountryCode: "XX",
			},
			valid: true, // Should be valid when disabled
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.valid {
				assert.NoError(t, err, "Config should be valid for %s", tt.name)
			} else {
				assert.Error(t, err, "Config should be invalid for %s", tt.name)
			}
		})
	}
}

func TestConfig_IsSovereignEnabled(t *testing.T) {
	config := Config{
		DualMode: true,
		Sovereign: SovereignConfig{
			Enabled:     true,
			CountryCode: "CA",
			PrivacyMode: "strict",
		},
	}

	assert.True(t, config.IsSovereignEnabled())
	assert.Equal(t, "CA", config.Sovereign.CountryCode)
}

func TestConfig_IsSovereignDisabled(t *testing.T) {
	config := Config{
		DualMode: true,
		Sovereign: SovereignConfig{
			Enabled:     false,
			CountryCode: "CA",
		},
	}

	assert.False(t, config.IsSovereignEnabled())
}

func TestSovereignConfig_ValidCountryCodes(t *testing.T) {
	validCodes := []string{"CA", "US", "GB", "FR", "DE", "AU", "NZ", "JP"}

	for _, code := range validCodes {
		t.Run("CountryCode_"+code, func(t *testing.T) {
			config := SovereignConfig{
				Enabled:     true,
				CountryCode: code,
			}

			err := config.Validate()
			assert.NoError(t, err, "Country code %s should be valid", code)
		})
	}
}

func TestSovereignConfig_CaseInsensitive(t *testing.T) {
	// Test lowercase
	config := SovereignConfig{
		Enabled:     true,
		CountryCode: "ca",
	}

	err := config.Validate()
	assert.NoError(t, err, "Lowercase country code should be valid")

	// Test mixed case
	config.CountryCode = "Ca"
	err = config.Validate()
	assert.NoError(t, err, "Mixed case country code should be valid")
}
