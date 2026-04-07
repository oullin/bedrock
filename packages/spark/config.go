package spark

// Config holds the merged configuration for the billing system. It combines
// settings from the Cashier layer (provider credentials, webhook) and the
// Spark layer (portal path, branding, billables).
type Config struct {
	// Provider credentials.
	SellerID        string
	ClientSideToken string
	APIKey          string
	RetainKey       string
	WebhookSecret   string

	// Portal settings.
	Path         string // URI prefix for the billing portal (default: "billing").
	DashboardURL string // URL to redirect to after billing actions.
	TermsURL     string // Terms of service URL shown in the portal.

	// Branding.
	BrandLogoPath string // Absolute path to SVG logo file.
	BrandColor    string // CSS color class or hex color.

	// Behaviour.
	Prorates   bool   // Whether plan changes are prorated.
	DateFormat string // Date format for the portal (Go reference time layout).
	Sandbox    bool   // Whether the provider is in sandbox/test mode.

	// Currency.
	DefaultCurrency string // ISO 4217 code (default: "USD").
	CurrencyLocale  string // Locale for formatting (default: "en").

	// Webhook.
	WebhookPath string // Route path for the provider webhook (default: "paddle/webhook").

	// Contact (for custom plan inquiries).
	ContactEmail         string
	ContactFallbackEmail string

	// Billables is the set of configured billable types.
	Billables map[string]BillableConfig
}

// BillableConfig holds configuration for a single billable type.
type BillableConfig struct {
	ModelName       string // Identifier for the billable model (e.g. "team").
	TrialDays       int    // Number of trial days granted on creation.
	DefaultInterval string // Default billing interval ("monthly" or "yearly").
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Path:            "billing",
		DashboardURL:    "/subscription",
		Prorates:        true,
		DateFormat:      "2 January 2006",
		DefaultCurrency: "USD",
		CurrencyLocale:  "en",
		WebhookPath:     "paddle/webhook",
		Billables:       make(map[string]BillableConfig),
	}
}
