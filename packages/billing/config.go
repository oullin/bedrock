package billing

// Config holds all Billing billing configuration.
// Mirrors config/billing.php.
type Config struct {
	Path              string                    // billing portal URI prefix (default: "billing")
	Middleware        []string                  // auth middleware stack
	Prorates          bool                      // enable subscription proration
	ProrationBehavior string                    // explicit proration behavior override
	DateFormat        string                    // Go reference time layout
	AppName           string                    // application name shown in the billing portal
	DashboardURL      string                    // URL of the main dashboard
	TermsURL          string                    // URL of terms of service
	BrandLogo         string                    // path or inline SVG for brand logo
	BrandColor        string                    // CSS class or hex colour for brand
	Sandbox           bool                      // provider sandbox mode
	ClientSideToken   string                    // Paddle client-side token
	SellerID          int                       // Paddle seller ID
	RetainKey         string                    // Paddle retain key
	Billables         map[string]BillableConfig // keyed by billable type name
	FeatureFlags      map[string]FeatureConfig  // keyed by feature name
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

// DefaultConfig returns a Config with sensible defaults matching Upstream's
// billing.php defaults.

// Features provides helpers to query feature flags from the config.
type Features struct {
	cfg Config
}

func DefaultConfig() Config {
	return Config{
		Path:       "billing",
		Prorates:   true,
		DateFormat: "January 2, 2006",
		AppName:    "Upstream",
		BrandColor: "bg-gray-800",
	}
}

// NewFeatures wraps a config for feature flag queries.
func NewFeatures(cfg Config) *Features {
	return &Features{cfg: cfg}
}

// Enabled reports whether the named feature is turned on.
func (f *Features) Enabled(feature string) bool {
	fc, ok := f.cfg.FeatureFlags[feature]

	return ok && fc.Enabled
}

// Option returns the value of a feature option, or nil.
func (f *Features) Option(feature, option string) any {
	fc, ok := f.cfg.FeatureFlags[feature]

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

// SendsPaymentNotificationEmails reports whether payment notification
// emails are enabled.
func (f *Features) SendsPaymentNotificationEmails() bool {
	return f.Enabled("sends-payment-notification-emails")
}
