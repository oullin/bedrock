package billing

import "github.com/bedrock/packages/config"

// Config is Billing's typed view over the shared Bedrock config repository.
type Config struct {
	repository *config.Repository
}

// BillableConfig describes a single billable type.
type BillableConfig struct {
	Model           string // identifier for the billable model (e.g. "team")
	TrialDays       int
	DefaultInterval string // "monthly" or "yearly"
}

// FeatureConfig stores a feature flag with optional key-value options.
type FeatureConfig struct {
	Enabled bool
	Options map[string]any
}

// DefaultConfigRepository returns the canonical config repository with Billing
// defaults registered under the "billing" namespace.

// DefaultConfig returns the typed Billing config backed by the shared repository.

// NewConfig wraps a shared config repository as Billing configuration.

// NewConfigFromValues returns a Billing config with the supplied values layered
// over the canonical defaults.

// Repository returns the underlying shared config repository.

// Path returns the billing portal URI prefix.

// Middleware returns the configured middleware aliases.

// Prorates reports whether subscription updates should prorate.

// ProrationBehavior returns the explicit provider proration behavior override.

// DateFormat returns the configured Go time layout.

// AppName returns the application name shown in the billing portal.

// DashboardURL returns the main application dashboard URL.

// TermsURL returns the terms of service URL.

// BrandLogo returns the configured brand logo path or inline SVG.

// BrandColor returns the billing portal brand color token.

// Sandbox reports whether the Paddle provider should use sandbox mode.

// ClientSideToken returns the Paddle client-side token.

// SellerID returns the Paddle seller ID.

// RetainKey returns the Paddle Retain key.

// Billables returns the configured billable models keyed by billable type.

// FeatureFlags returns the configured Billing feature flags.

// Features provides helpers to query feature flags from the config.
type Features struct {
	cfg *Config
}

const configPrefix = "billing."

func DefaultConfigRepository() *config.Repository {
	return config.NewWithDefaults(map[string]any{
		"billing.path":               "billing",
		"billing.middleware":         []string{},
		"billing.prorates":           true,
		"billing.proration_behavior": "",
		"billing.date_format":        "January 2, 2006",
		"billing.app_name":           "Upstream",
		"billing.dashboard_url":      "",
		"billing.terms_url":          "",
		"billing.brand_logo":         "",
		"billing.brand_color":        "bg-gray-800",
		"billing.sandbox":            false,
		"billing.client_side_token":  "",
		"billing.seller_id":          0,
		"billing.retain_key":         "",
		"billing.billables":          map[string]BillableConfig{},
		"billing.feature_flags":      map[string]FeatureConfig{},
	})
}

func DefaultConfig() *Config {
	return NewConfig(DefaultConfigRepository())
}

func NewConfig(repository *config.Repository) *Config {
	if repository == nil {
		repository = DefaultConfigRepository()
	}

	return &Config{repository: repository}
}

func NewConfigFromValues(values map[string]any) *Config {
	repository := DefaultConfigRepository()
	repository.SetMany(values)

	return NewConfig(repository)
}

func (c *Config) Repository() *config.Repository {
	if c == nil || c.repository == nil {
		return DefaultConfigRepository()
	}

	return c.repository
}

func (c *Config) string(key, fallback string) string {
	value, err := c.Repository().String(configPrefix+key, fallback)

	if err != nil {
		return fallback
	}

	return value
}

func (c *Config) integer(key string, fallback int) int {
	value, err := c.Repository().Integer(configPrefix+key, fallback)

	if err != nil {
		return fallback
	}

	return value
}

func (c *Config) boolean(key string, fallback bool) bool {
	value, err := c.Repository().Boolean(configPrefix+key, fallback)

	if err != nil {
		return fallback
	}

	return value
}

func (c *Config) Path() string { return c.string("path", "billing") }

func (c *Config) Middleware() []string {
	value := c.Repository().Get(configPrefix+"middleware", []string{})

	if typed, ok := value.([]string); ok {
		return append([]string(nil), typed...)
	}

	if values, ok := value.([]any); ok {
		out := make([]string, 0, len(values))

		for _, v := range values {
			if s, ok := v.(string); ok {
				out = append(out, s)
			}
		}

		return out
	}

	return []string{}
}

func (c *Config) Prorates() bool { return c.boolean("prorates", true) }

func (c *Config) ProrationBehavior() string { return c.string("proration_behavior", "") }

func (c *Config) DateFormat() string { return c.string("date_format", "January 2, 2006") }

func (c *Config) AppName() string { return c.string("app_name", "Upstream") }

func (c *Config) DashboardURL() string { return c.string("dashboard_url", "") }

func (c *Config) TermsURL() string { return c.string("terms_url", "") }

func (c *Config) BrandLogo() string { return c.string("brand_logo", "") }

func (c *Config) BrandColor() string { return c.string("brand_color", "bg-gray-800") }

func (c *Config) Sandbox() bool { return c.boolean("sandbox", false) }

func (c *Config) ClientSideToken() string { return c.string("client_side_token", "") }

func (c *Config) SellerID() int { return c.integer("seller_id", 0) }

func (c *Config) RetainKey() string { return c.string("retain_key", "") }

func (c *Config) Billables() map[string]BillableConfig {
	value := c.Repository().Get(configPrefix + "billables")

	if typed, ok := value.(map[string]BillableConfig); ok {
		out := make(map[string]BillableConfig, len(typed))

		for key, billable := range typed {
			out[key] = billable
		}

		return out
	}

	if values, ok := value.(map[string]any); ok {
		return billableConfigsFromMap(values)
	}

	return map[string]BillableConfig{}
}

func (c *Config) FeatureFlags() map[string]FeatureConfig {
	value := c.Repository().Get(configPrefix + "feature_flags")

	if typed, ok := value.(map[string]FeatureConfig); ok {
		out := make(map[string]FeatureConfig, len(typed))

		for key, feature := range typed {
			out[key] = feature
		}

		return out
	}

	if values, ok := value.(map[string]any); ok {
		return featureConfigsFromMap(values)
	}

	return map[string]FeatureConfig{}
}

// NewFeatures wraps a config for feature flag queries.
func NewFeatures(cfg *Config) *Features {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &Features{cfg: cfg}
}

// Enabled reports whether the named feature is turned on.
func (f *Features) Enabled(feature string) bool {
	fc, ok := f.cfg.FeatureFlags()[feature]

	return ok && fc.Enabled
}

// Option returns the value of a feature option, or nil.
func (f *Features) Option(feature, option string) any {
	fc, ok := f.cfg.FeatureFlags()[feature]

	if !ok {
		return nil
	}

	return fc.Options[option]
}

// CollectsBillingAddress reports whether billing address collection is enabled.
func (f *Features) CollectsBillingAddress() bool {
	return f.Enabled("billing-address-collection")
}

// CollectsEuVat reports whether EU VAT collection is enabled.
func (f *Features) CollectsEuVat() bool {
	return f.Enabled("eu-vat-collection")
}

// EnforcesAcceptingTerms reports whether users must accept terms.
func (f *Features) EnforcesAcceptingTerms() bool {
	return f.Enabled("must-accept-terms")
}

// SendsInvoiceEmails reports whether invoice emails are enabled.
func (f *Features) SendsInvoiceEmails() bool {
	return f.Enabled("invoice-emails-sending")
}

// SendsPaymentNotificationEmails reports whether payment notification emails
// are enabled.
func (f *Features) SendsPaymentNotificationEmails() bool {
	return f.Enabled("sends-payment-notification-emails")
}

func billableConfigsFromMap(values map[string]any) map[string]BillableConfig {
	out := make(map[string]BillableConfig, len(values))

	for key, value := range values {
		switch typed := value.(type) {
		case BillableConfig:
			out[key] = typed
		case map[string]any:
			out[key] = BillableConfig{
				Model:           stringFromAny(typed["model"]),
				TrialDays:       intFromAny(typed["trial_days"]),
				DefaultInterval: stringFromAny(typed["default_interval"]),
			}
		}
	}

	return out
}

func featureConfigsFromMap(values map[string]any) map[string]FeatureConfig {
	out := make(map[string]FeatureConfig, len(values))

	for key, value := range values {
		switch typed := value.(type) {
		case FeatureConfig:
			out[key] = typed
		case map[string]any:
			feature := FeatureConfig{
				Enabled: boolFromAny(typed["enabled"]),
				Options: map[string]any{},
			}

			if options, ok := typed["options"].(map[string]any); ok {
				feature.Options = options
			}

			out[key] = feature
		}
	}

	return out
}

func stringFromAny(value any) string {
	if s, ok := value.(string); ok {
		return s
	}

	return ""
}

func intFromAny(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func boolFromAny(value any) bool {
	if b, ok := value.(bool); ok {
		return b
	}

	return false
}
