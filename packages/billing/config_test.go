package spark_test

import (
	"testing"

	"github.com/bedrock/packages/billing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := billing.DefaultConfig()

	if cfg.Path != "billing" {
		t.Errorf("Path = %q, want %q", cfg.Path, "billing")
	}

	if cfg.DashboardURL != "/subscription" {
		t.Errorf("DashboardURL = %q, want %q", cfg.DashboardURL, "/subscription")
	}

	if !cfg.Prorates {
		t.Error("Prorates = false, want true")
	}

	if cfg.DefaultCurrency != "USD" {
		t.Errorf("DefaultCurrency = %q, want %q", cfg.DefaultCurrency, "USD")
	}

	if cfg.WebhookPath != "paddle/webhook" {
		t.Errorf("WebhookPath = %q, want %q", cfg.WebhookPath, "paddle/webhook")
	}

	if cfg.Billables == nil {
		t.Error("Billables is nil, want initialised map")
	}
}
