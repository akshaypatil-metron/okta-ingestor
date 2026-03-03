package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	cfg := Load()

	if cfg.BatchSize != 100 {
		t.Errorf("Expected BatchSize 100, got %d", cfg.BatchSize)
	}
	if cfg.PollInterval != 30*time.Second {
		t.Errorf("Expected PollInterval 30s, got %v", cfg.PollInterval)
	}
	if cfg.OktaDomain == "" || cfg.OktaToken == "" {
		t.Errorf("Expected Okta credentials to be populated")
	}
	if cfg.StartDate == "" {
		t.Errorf("Expected StartDate to be populated")
	}
}
