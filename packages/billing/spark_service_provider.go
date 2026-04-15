package billing

import "github.com/bedrock/packages/container"

// BillingServiceProvider registers the Billing billing manager into the container.
// It mirrors Billing\BillingServiceProvider.
type BillingServiceProvider struct {
	app *container.Container
	cfg Config
}

// NewBillingServiceProvider constructs the provider.
// cfg is the Billing configuration (see DefaultConfig() for sensible defaults).
func NewBillingServiceProvider(app *container.Container, cfg Config) *BillingServiceProvider {
	return &BillingServiceProvider{app: app, cfg: cfg}
}

// Register binds the Billing manager as a singleton under "billing"
// and the Config under "billing.config".
func (p *BillingServiceProvider) Register() {
	cfg := p.cfg

	p.app.Instance("billing.config", cfg)

	p.app.Singleton("billing", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetProrates(cfg.Prorates)

		for billableType, billCfg := range cfg.Billables {
			c := billCfg
			c.Model = billableType
			m.RegisterBillable(c)
		}

		return m, nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *BillingServiceProvider) Provides() []string {
	return []string{"billing", "billing.config"}
}
