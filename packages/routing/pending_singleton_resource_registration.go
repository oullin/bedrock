package routing

// Ref: @bedrock/code-0326
// twin of [PendingResourceRegistration].
type PendingSingletonResourceRegistration struct {
	registrar  *ResourceRegistrar
	name       string
	controller string
	options    map[string]any
	registered *RouteCollection
}

// NewPendingSingletonResourceRegistration constructs the builder.
func NewPendingSingletonResourceRegistration(rr *ResourceRegistrar, name, controller string, options map[string]any) *PendingSingletonResourceRegistration {
	if options == nil {
		options = map[string]any{}
	}

	return &PendingSingletonResourceRegistration{
		registrar:  rr,
		name:       name,
		controller: controller,
		options:    options,
	}
}

// Only restricts the singleton to the named actions.
func (p *PendingSingletonResourceRegistration) Only(actions ...string) *PendingSingletonResourceRegistration {
	p.options["only"] = actions
	delete(p.options, "except")

	return p
}

// Except removes the named actions.
func (p *PendingSingletonResourceRegistration) Except(actions ...string) *PendingSingletonResourceRegistration {
	existing, _ := p.options["except"].([]string)
	p.options["except"] = appendUniqueStrings(existing, actions...)

	return p
}

// Creatable enables the create/store/destroy routes for the singleton.
func (p *PendingSingletonResourceRegistration) Creatable() *PendingSingletonResourceRegistration {
	p.options["creatable"] = true

	return p
}

// Destroyable enables the destroy route for the singleton.
func (p *PendingSingletonResourceRegistration) Destroyable() *PendingSingletonResourceRegistration {
	p.options["destroyable"] = true

	return p
}

// Register materializes the routes.
func (p *PendingSingletonResourceRegistration) Register() *RouteCollection {
	if p.registered != nil {
		return p.registered
	}

	p.registered = p.registrar.Singleton(p.name, p.controller, p.options)

	return p.registered
}

// Routes is a parity alias for [Register].
func (p *PendingSingletonResourceRegistration) Routes() *RouteCollection { return p.Register() }
