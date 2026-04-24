package spark_test

import (
	"testing"

	"github.com/bedrock/packages/config"
	"github.com/bedrock/packages/billing"
)

func TestConfigUsesSharedRepositoryDefaultsAndOverrides(t *testing.T) {
	t.Parallel()

	repository := billing.DefaultConfigRepository()
	repository.SetMany(map[string]any{
		"billing.path":          "account/billing",
		"billing.dashboard_url": "/dashboard",
		"billing.seller_id":     1234,
		"billing.prorates":      false,
		"billing.billables": map[string]billing.BillableConfig{
			"team": {DefaultInterval: "yearly"},
		},
	})

	cfg := billing.NewConfig(repository)

	if cfg.Path() != "account/billing" {
		t.Fatalf("Path() = %q, want account/billing", cfg.Path())
	}

	if cfg.DashboardURL() != "/dashboard" {
		t.Fatalf("DashboardURL() = %q, want /dashboard", cfg.DashboardURL())
	}

	if cfg.SellerID() != 1234 {
		t.Fatalf("SellerID() = %d, want 1234", cfg.SellerID())
	}

	if cfg.Prorates() {
		t.Fatal("Prorates() = true, want false")
	}

	if cfg.Billables()["team"].DefaultInterval != "yearly" {
		t.Fatalf("billable interval = %q, want yearly", cfg.Billables()["team"].DefaultInterval)
	}
}

func TestNewConfigAcceptsSharedConfigRepository(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfig(config.New(map[string]any{
		"billing": map[string]any{
			"path":        "billing",
			"brand_color": "bg-black",
		},
	}))

	if cfg.Path() != "billing" {
		t.Fatalf("Path() = %q, want billing", cfg.Path())
	}

	if cfg.BrandColor() != "bg-black" {
		t.Fatalf("BrandColor() = %q, want bg-black", cfg.BrandColor())
	}
}
