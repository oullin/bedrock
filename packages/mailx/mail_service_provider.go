package mailx

import (
	"github.com/bedrock/packages/config"
	"github.com/bedrock/packages/container"
)

// MailFromConfig is the default From address applied to outgoing mail when
// a Mailable does not specify one.
type MailFromConfig struct {
	Address string
	Name    string
}

// MailProviderConfig is the typed configuration consumed by MailServiceProvider.
// It exposes only the keys that the MailManager reads from the underlying
// config.Repository so callers don't need to know about config.Repository.
type MailProviderConfig struct {
	// Default is the default mailer name (e.g. "smtp", "log", "array").
	Default string

	// From is the default sender address/name.
	From MailFromConfig

	// Mailers maps mailer name → transport options (e.g. "transport", "host").
	Mailers map[string]map[string]any
}

// MailServiceProvider registers the mail manager into the container.
// Ref: @bedrock/code-0226
type MailServiceProvider struct {
	app  *container.Container
	cfg  MailProviderConfig
	opts []ManagerOption
}

// NewMailServiceProvider constructs the provider from a typed config.
// opts are forwarded to NewManager (e.g. WithManagerDispatcher, WithLogger).
func NewMailServiceProvider(app *container.Container, cfg MailProviderConfig, opts ...ManagerOption) *MailServiceProvider {
	return &MailServiceProvider{app: app, cfg: cfg, opts: opts}
}

// Register binds the mail manager as a singleton under "mailer".
func (p *MailServiceProvider) Register() {
	p.app.Singleton("mailer", func(_ *container.Container) (any, error) {
		repo := p.cfg.toRepository()

		return NewManager(repo, p.opts...), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *MailServiceProvider) Provides() []string {
	return []string{"mailer"}
}

func (c MailProviderConfig) toRepository() *config.Repository {
	def := c.Default

	if def == "" {
		def = "smtp"
	}

	mailers := c.Mailers

	if mailers == nil {
		mailers = map[string]map[string]any{}
	}

	mail := map[string]any{
		"default": def,
		"from": map[string]any{
			"address": c.From.Address,
			"name":    c.From.Name,
		},
		"mailers": mailers,
	}

	return config.New(map[string]any{"mail": mail})
}
