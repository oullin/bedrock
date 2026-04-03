package authflows

import "testing"

func TestFeaturesAndRoutePath(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Features: []string{FeatureTwoFactorAuthentication, FeatureRegistration},
		Options: map[string]map[string]any{
			FeatureTwoFactorAuthentication: {
				"confirm": true,
				"window":  "2",
			},
		},
		Paths: map[string]string{
			"login": "/sign-in",
		},
	}

	features := NewFeatures(cfg)
	if !features.Enabled(FeatureRegistration) {
		t.Fatal("expected registration feature to be enabled")
	}
	if !features.OptionEnabled(FeatureTwoFactorAuthentication, "confirm") {
		t.Fatal("expected confirm option to be enabled")
	}
	if got := features.OptionInt(FeatureTwoFactorAuthentication, "window", 1); got != 2 {
		t.Fatalf("unexpected window: %d", got)
	}

	paths := NewRoutePath(cfg.Paths)
	if got := paths.For("login", "/login"); got != "/sign-in" {
		t.Fatalf("unexpected custom path: %s", got)
	}
	if got := paths.For("logout", "/logout"); got != "/logout" {
		t.Fatalf("unexpected default path: %s", got)
	}
}
