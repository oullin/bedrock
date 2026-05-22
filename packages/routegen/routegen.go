package routegen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Options controls what the generator produces.
type Options struct {
	// Path is the root output directory. The generator creates actions/,
	// routes/, and routegen/ subdirectories inside it.
	// Defaults to "resources/js" when empty.
	Path string

	// SkipActions disables generation of the actions/ directory.
	SkipActions bool

	// SkipRoutes disables generation of the routes/ directory.
	SkipRoutes bool

	// WithForm enables form-helper generation (RouteFormDefinition exports
	// and .form property on each action function).
	WithForm bool
}

// Generate writes TypeScript route helpers for all provided routes.
//
//	{Path}/actions/   — per-controller action functions
//	{Path}/routes/    — named-route helpers
//	{Path}/routegen/ — runtime TypeScript utility (index.ts)
func Generate(routes []*RouteInfo, opts Options) error {
	if opts.Path == "" {
		opts.Path = filepath.Join("resources", "js")
	}

	g := newGenerator(opts)

	if !opts.SkipActions {
		actionsBase := filepath.Join(opts.Path, "actions")

		if err := os.RemoveAll(actionsBase); err != nil {
			return fmt.Errorf("routegen: clearing actions dir: %w", err)
		}

		g.generateActions(routes, actionsBase)

		if err := g.flush(actionsBase); err != nil {
			return fmt.Errorf("routegen: writing actions: %w", err)
		}
	}

	if !opts.SkipRoutes {
		routesBase := filepath.Join(opts.Path, "routes")

		if err := os.RemoveAll(routesBase); err != nil {
			return fmt.Errorf("routegen: clearing routes dir: %w", err)
		}

		g.generateRoutes(routes, routesBase)

		if err := g.flush(routesBase); err != nil {
			return fmt.Errorf("routegen: writing routes: %w", err)
		}
	}

	// Always copy the runtime utility.
	routegenDir := filepath.Join(opts.Path, "routegen")

	if err := writeRouteGenTS(routegenDir); err != nil {
		return fmt.Errorf("routegen: writing runtime utility: %w", err)
	}

	return nil
}
