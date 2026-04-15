package translation

import "github.com/bedrock/packages/container"

// TranslationServiceProvider registers the translator into the container.
// It mirrors Illuminate\Translation\TranslationServiceProvider.
type TranslationServiceProvider struct {
	app    *container.Container
	loader Loader
	locale string
}

// NewTranslationServiceProvider constructs the provider.
// loader provides translation messages; locale is the default language (e.g. "en").
func NewTranslationServiceProvider(app *container.Container, loader Loader, locale string) *TranslationServiceProvider {
	return &TranslationServiceProvider{app: app, loader: loader, locale: locale}
}

// Register binds the translator as a singleton under "translator".
func (p *TranslationServiceProvider) Register() {
	p.app.Singleton("translator", func(_ *container.Container) (any, error) {
		return NewTranslator(p.loader, p.locale), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *TranslationServiceProvider) Provides() []string {
	return []string{"translator"}
}
