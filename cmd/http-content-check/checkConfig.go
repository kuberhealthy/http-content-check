package main

import (
	"fmt"
	"os"
	"time"
)

// CheckConfig stores configuration for the HTTP content check.
type CheckConfig struct {
	// TargetURL is the URL to request.
	TargetURL string
	// TargetString is the string to locate in the response body.
	TargetString string
	// TimeoutDuration is the timeout for the HTTP request.
	TimeoutDuration time.Duration
}

// parseConfig loads environment variables into a CheckConfig.
func parseConfig() (*CheckConfig, error) {
	// Read required target URL.
	targetURL := os.Getenv("TARGET_URL")
	if len(targetURL) == 0 {
		return nil, fmt.Errorf("no URL provided in YAML")
	}

	// Read required target string.
	targetString := os.Getenv("TARGET_STRING")
	if len(targetString) == 0 {
		return nil, fmt.Errorf("no string provided in YAML")
	}

	// Parse timeout duration.
	timeoutEnv := os.Getenv("TIMEOUT_DURATION")
	timeoutDuration, err := time.ParseDuration(timeoutEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TIMEOUT_DURATION: %w", err)
	}

	// Assemble configuration.
	cfg := &CheckConfig{}
	cfg.TargetURL = targetURL
	cfg.TargetString = targetString
	cfg.TimeoutDuration = timeoutDuration

	return cfg, nil
}
