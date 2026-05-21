package scout

import (
	"github.com/bedrock/packages/container"
	contract "github.com/bedrock/packages/contracts/scout"
	"github.com/bedrock/packages/scout/engines"
)

// ScoutServiceProvider registers the Scout engine manager into the container.
type ScoutServiceProvider struct {
	app    *container.Container
	config Config
}

// NewScoutServiceProvider constructs the provider.
func NewScoutServiceProvider(app *container.Container, config Config) *ScoutServiceProvider {
	return &ScoutServiceProvider{app: app, config: config}
}

// Register binds the engine manager as a singleton under "scout".
func (p *ScoutServiceProvider) Register() {
	p.app.Singleton("scout", func(_ *container.Container) (any, error) {
		m := NewEngineManager()
		m.SetDefaultDriver(p.config.Driver)

		// Register built-in engine drivers.
		m.Extend("null", func(_ map[string]any) (contract.Engine, error) {
			return engines.NewNullEngine(), nil
		})

		m.Extend("collection", func(_ map[string]any) (contract.Engine, error) {
			return engines.NewCollectionEngine(p.config.SoftDelete), nil
		})

		// Build and register the default engines that don't require external deps.
		nullEngine := engines.NewNullEngine()
		m.Register("null", nullEngine)

		collectionEngine := engines.NewCollectionEngine(p.config.SoftDelete)
		m.Register("collection", collectionEngine)

		return m, nil
	})
}

// Boot performs post-registration setup. It registers the model observer
// with the event dispatcher if available.
func (p *ScoutServiceProvider) Boot() {
	// The model observer is set up by the consumer when they have access
	// to the event dispatcher. Example:
	//
	//   manager, _ := app.Make("scout")
	//   observer := scout.NewModelObserver(manager.(*scout.EngineManager), config)
	//   dispatcher.Listen(dbevents.Saved{}, func(ctx context.Context, event any) (any, error) {
	//       observer.Saved(ctx, event.(dbevents.Saved))
	//       return nil, nil
	//   })
}

// Provides returns the abstract keys registered by this provider.
func (p *ScoutServiceProvider) Provides() []string {
	return []string{"scout"}
}
