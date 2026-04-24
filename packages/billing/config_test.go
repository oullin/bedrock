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

func TestNewConfig_NilRepositoryFallsBackToDefaults(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfig(nil)

	if cfg.Path() != "billing" {
		t.Fatalf("Path() = %q, want billing (default)", cfg.Path())
	}
}

func TestConfigNilReceiverRepositoryIsSafe(t *testing.T) {
	t.Parallel()

	var cfg *billing.Config

	if cfg.Repository() == nil {
		t.Fatal("Repository() must not return nil even on nil receiver")
	}
}

func TestNewConfigFromValuesOverridesDefaults(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.path":      "acct",
		"billing.seller_id": 7,
	})

	if cfg.Path() != "acct" {
		t.Fatalf("Path() = %q", cfg.Path())
	}

	if cfg.SellerID() != 7 {
		t.Fatalf("SellerID() = %d", cfg.SellerID())
	}
}

func TestMiddleware_TypedSlice(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.middleware": []string{"auth", "verified"},
	})

	got := cfg.Middleware()

	if len(got) != 2 || got[0] != "auth" || got[1] != "verified" {
		t.Fatalf("Middleware() = %#v", got)
	}
}

func TestMiddleware_AnySlice(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.middleware": []any{"auth", 42, "verified"},
	})

	got := cfg.Middleware()

	if len(got) != 2 || got[0] != "auth" || got[1] != "verified" {
		t.Fatalf("Middleware() = %#v", got)
	}
}

func TestMiddleware_UnexpectedTypeReturnsEmpty(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.middleware": "not-a-slice",
	})

	if got := cfg.Middleware(); len(got) != 0 {
		t.Fatalf("Middleware() = %#v, want empty", got)
	}
}

func TestBillables_AnyMap(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.billables": map[string]any{
			"team": map[string]any{"model": "team", "trial_days": 14, "default_interval": "monthly"},
		},
	})

	got := cfg.Billables()

	if got["team"].TrialDays != 14 || got["team"].DefaultInterval != "monthly" {
		t.Fatalf("Billables() = %#v", got)
	}
}

func TestBillables_UnexpectedTypeReturnsEmpty(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.billables": []string{"team"},
	})

	if got := cfg.Billables(); len(got) != 0 {
		t.Fatalf("Billables() = %#v, want empty", got)
	}
}

func TestFeatureFlags_TypedMap(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.feature_flags": map[string]billing.FeatureConfig{
			"feature-a": {Enabled: true, Options: map[string]any{"scope": "team"}},
		},
	})

	got := cfg.FeatureFlags()

	if !got["feature-a"].Enabled {
		t.Fatalf("feature-a enabled = false")
	}

	if got["feature-a"].Options["scope"] != "team" {
		t.Fatalf("feature-a options = %#v", got["feature-a"].Options)
	}
}

func TestFeatureFlags_AnyMap(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.feature_flags": map[string]any{
			"feature-a": map[string]any{"enabled": true},
		},
	})

	got := cfg.FeatureFlags()

	if !got["feature-a"].Enabled {
		t.Fatalf("feature-a enabled = false, want true")
	}
}

func TestFeatureFlags_UnexpectedTypeReturnsEmpty(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.feature_flags": []string{"a"},
	})

	if got := cfg.FeatureFlags(); len(got) != 0 {
		t.Fatalf("FeatureFlags() = %#v, want empty", got)
	}
}

func TestFeatures_NilConfigUsesDefaults(t *testing.T) {
	t.Parallel()

	f := billing.NewFeatures(nil)

	if f.Enabled("anything") {
		t.Fatal("default features should not be enabled")
	}

	if f.Option("feature", "option") != nil {
		t.Fatal("missing feature option should return nil")
	}
}

func TestFeatures_EnabledAndOption(t *testing.T) {
	t.Parallel()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.feature_flags": map[string]billing.FeatureConfig{
			"billing-address-collection":        {Enabled: true},
			"eu-vat-collection":                 {Enabled: true},
			"must-accept-terms":                 {Enabled: true},
			"invoice-emails-sending":            {Enabled: true},
			"sends-payment-notification-emails": {Enabled: true},
			"custom-feature":                    {Enabled: true, Options: map[string]any{"opt": 42}},
		},
	})

	features := billing.NewFeatures(cfg)

	if !features.CollectsBillingAddress() {
		t.Fatal("CollectsBillingAddress = false")
	}

	if !features.CollectsEuVat() {
		t.Fatal("CollectsEuVat = false")
	}

	if !features.EnforcesAcceptingTerms() {
		t.Fatal("EnforcesAcceptingTerms = false")
	}

	if !features.SendsInvoiceEmails() {
		t.Fatal("SendsInvoiceEmails = false")
	}

	if !features.SendsPaymentNotificationEmails() {
		t.Fatal("SendsPaymentNotificationEmails = false")
	}

	if features.Option("custom-feature", "opt") != 42 {
		t.Fatalf("Option = %v", features.Option("custom-feature", "opt"))
	}

	if features.Option("unknown", "opt") != nil {
		t.Fatal("unknown feature option should be nil")
	}
}

func TestDefaultConfigExposesAllDefaults(t *testing.T) {
	t.Parallel()

	cfg := billing.DefaultConfig()

	if cfg.Path() != "billing" {
		t.Fatalf("Path = %q", cfg.Path())
	}

	if !cfg.Prorates() {
		t.Fatal("Prorates default = false")
	}

	if cfg.DateFormat() != "January 2, 2006" {
		t.Fatalf("DateFormat = %q", cfg.DateFormat())
	}

	if cfg.AppName() != "Upstream" {
		t.Fatalf("AppName = %q", cfg.AppName())
	}

	if cfg.BrandColor() != "bg-gray-800" {
		t.Fatalf("BrandColor = %q", cfg.BrandColor())
	}

	if cfg.Sandbox() {
		t.Fatal("Sandbox default = true")
	}

	if cfg.DashboardURL() != "" || cfg.TermsURL() != "" || cfg.BrandLogo() != "" || cfg.ClientSideToken() != "" || cfg.RetainKey() != "" || cfg.ProrationBehavior() != "" {
		t.Fatal("empty-string defaults diverged")
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
