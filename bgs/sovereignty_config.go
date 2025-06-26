package bgs

import (
	"errors"
	"strings"
)

// SovereignConfig holds sovereignty-specific configuration
type SovereignConfig struct {
	Enabled           bool   `json:"enabled"`
	CountryCode       string `json:"country_code"`
	DataRetentionDays int    `json:"data_retention_days"`
	PrivacyMode       string `json:"privacy_mode"`
}

// Validate checks if the sovereignty configuration is valid
func (sc *SovereignConfig) Validate() error {
	if !sc.Enabled {
		return nil // Valid when disabled
	}

	if sc.CountryCode == "" {
		return errors.New("country code is required when sovereignty is enabled")
	}

	// Validate country code format (ISO 3166-1 alpha-2)
	if len(sc.CountryCode) != 2 {
		return errors.New("country code must be 2 characters (ISO 3166-1 alpha-2)")
	}

	// Check for valid country codes (sample validation)
	validCodes := []string{"CA", "US", "GB", "FR", "DE", "AU", "NZ", "JP", "KR", "IN", "BR", "MX"}
	valid := false
	upperCode := strings.ToUpper(sc.CountryCode)
	for _, code := range validCodes {
		if upperCode == code {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("invalid country code")
	}

	return nil
}

// Config represents the main configuration with sovereignty settings
type Config struct {
	DualMode  bool            `json:"dual_mode"`
	Sovereign SovereignConfig `json:"sovereign"`
}

// IsSovereignEnabled checks if sovereignty mode is enabled
func (c *Config) IsSovereignEnabled() bool {
	return c.Sovereign.Enabled
}
