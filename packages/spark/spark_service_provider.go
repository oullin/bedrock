package spark

import "github.com/bedrock/packages/container"

// SparkServiceProvider registers the Spark billing manager into the container.
type SparkServiceProvider struct {
	app *container.Container
	cfg *Config
}

// NewSparkServiceProvider constructs the provider.
// cfg is the Spark configuration (see DefaultConfig() for sensible defaults).
func NewSparkServiceProvider(app *container.Container, cfg *Config) *SparkServiceProvider {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &SparkServiceProvider{app: app, cfg: cfg}
}

// Register binds the Spark manager as a singleton under "spark"
// and the Config under "spark.config".
func (p *SparkServiceProvider) Register() {
	cfg := p.cfg

	p.app.Instance("spark.config", cfg)

	p.app.Singleton("spark", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetProrates(cfg.Prorates())

		for billableType, billCfg := range cfg.Billables() {
			c := billCfg
			c.Model = billableType
			m.RegisterBillable(c)
		}

		return m, nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *SparkServiceProvider) Provides() []string {
	return []string{"spark", "spark.config"}
}
