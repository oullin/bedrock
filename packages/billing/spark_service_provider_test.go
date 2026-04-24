package spark_test

import (
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/billing"
)

func TestNewBillingServiceProvider_NilConfigFallsBackToDefault(t *testing.T) {
	app := container.New()

	provider := billing.NewBillingServiceProvider(app, nil)
	provider.Register()

	cfgValue, err := app.Get("billing.config")

	if err != nil {
		t.Fatalf("Get billing.config: %v", err)
	}

	cfg, ok := cfgValue.(*billing.Config)

	if !ok {
		t.Fatalf("billing.config not *Config: %T", cfgValue)
	}

	if cfg.Path() != "billing" {
		t.Fatalf("default path = %q", cfg.Path())
	}
}

func TestBillingServiceProvider_RegisterBindsManagerWithBillables(t *testing.T) {
	app := container.New()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.prorates": false,
		"billing.billables": map[string]billing.BillableConfig{
			"team": {DefaultInterval: "yearly"},
		},
	})

	provider := billing.NewBillingServiceProvider(app, cfg)
	provider.Register()

	mgrValue, err := app.Get("billing")

	if err != nil {
		t.Fatalf("Get billing: %v", err)
	}

	mgr, ok := mgrValue.(*billing.Manager)

	if !ok {
		t.Fatalf("billing not *Manager: %T", mgrValue)
	}

	if mgr.Prorates() {
		t.Fatal("manager Prorates() should match cfg.Prorates() = false")
	}

	cfgRegistered, ok := mgr.BillableConfig("team")

	if !ok {
		t.Fatal("team billable not registered")
	}

	if cfgRegistered.Model != "team" {
		t.Fatalf("billable Model = %q, want team", cfgRegistered.Model)
	}

	if cfgRegistered.DefaultInterval != "yearly" {
		t.Fatalf("billable DefaultInterval = %q", cfgRegistered.DefaultInterval)
	}
}

func TestBillingServiceProvider_Provides(t *testing.T) {
	provider := billing.NewBillingServiceProvider(container.New(), nil)

	got := provider.Provides()

	want := map[string]bool{"billing": true, "billing.config": true}

	for _, key := range got {
		if !want[key] {
			t.Fatalf("unexpected provide key: %q", key)
		}

		delete(want, key)
	}

	if len(want) != 0 {
		t.Fatalf("missing provide keys: %v", want)
	}
}
