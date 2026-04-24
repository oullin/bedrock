package spark_test

import (
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/spark"
)

func TestNewSparkServiceProvider_NilConfigFallsBackToDefault(t *testing.T) {
	app := container.New()

	provider := spark.NewSparkServiceProvider(app, nil)
	provider.Register()

	cfgValue, err := app.Get("spark.config")

	if err != nil {
		t.Fatalf("Get spark.config: %v", err)
	}

	cfg, ok := cfgValue.(*spark.Config)

	if !ok {
		t.Fatalf("spark.config not *Config: %T", cfgValue)
	}

	if cfg.Path() != "billing" {
		t.Fatalf("default path = %q", cfg.Path())
	}
}

func TestSparkServiceProvider_RegisterBindsManagerWithBillables(t *testing.T) {
	app := container.New()

	cfg := spark.NewConfigFromValues(map[string]any{
		"spark.prorates": false,
		"spark.billables": map[string]spark.BillableConfig{
			"team": {DefaultInterval: "yearly"},
		},
	})

	provider := spark.NewSparkServiceProvider(app, cfg)
	provider.Register()

	mgrValue, err := app.Get("spark")

	if err != nil {
		t.Fatalf("Get spark: %v", err)
	}

	mgr, ok := mgrValue.(*spark.Manager)

	if !ok {
		t.Fatalf("spark not *Manager: %T", mgrValue)
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

func TestSparkServiceProvider_Provides(t *testing.T) {
	provider := spark.NewSparkServiceProvider(container.New(), nil)

	got := provider.Provides()

	want := map[string]bool{"spark": true, "spark.config": true}

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
