package routing

// PendingResourceRegistration is a fluent builder returned by
// [Router.Resource]. The PHP equivalent applies its options on destruct;
// the Go form applies them when [PendingResourceRegistration.Register] is
// called explicitly, or implicitly when the user reads the resulting
// collection via [PendingResourceRegistration.Routes].
//
// Ref: @bedrock/code-0325
type PendingResourceRegistration struct {
	registrar  *ResourceRegistrar
	name       string
	controller string
	options    map[string]any
	registered *RouteCollection
}

// NewPendingResourceRegistration constructs the builder.
func NewPendingResourceRegistration(rr *ResourceRegistrar, name, controller string, options map[string]any) *PendingResourceRegistration {
	if options == nil {
		options = map[string]any{}
	}

	return &PendingResourceRegistration{
		registrar:  rr,
		name:       name,
		controller: controller,
		options:    options,
	}
}

// Only restricts the resource to the named actions.
func (p *PendingResourceRegistration) Only(actions ...string) *PendingResourceRegistration {
	p.options["only"] = actions
	delete(p.options, "except")

	return p
}

// Except removes the named actions from the resource.
func (p *PendingResourceRegistration) Except(actions ...string) *PendingResourceRegistration {
	existing, _ := p.options["except"].([]string)
	p.options["except"] = appendUniqueStrings(existing, actions...)

	return p
}

// Names overrides the route names for individual actions.
func (p *PendingResourceRegistration) Names(names map[string]string) *PendingResourceRegistration {
	p.options["names"] = names

	return p
}

// Parameters overrides the URL placeholder names.
func (p *PendingResourceRegistration) Parameters(parameters map[string]string) *PendingResourceRegistration {
	p.options["parameters"] = parameters

	for k, v := range parameters {
		p.registrar.parameters[k] = v
	}

	return p
}

// Middleware applies middleware to every action in the resource.
func (p *PendingResourceRegistration) Middleware(middleware ...any) *PendingResourceRegistration {
	p.options["middleware"] = middleware

	return p
}

// As prepends a name prefix to every action route name.
func (p *PendingResourceRegistration) As(prefix string) *PendingResourceRegistration {
	p.options["as"] = prefix

	return p
}

// Shallow toggles shallow-nested resource routing for nested resources.
func (p *PendingResourceRegistration) Shallow(shallow ...bool) *PendingResourceRegistration {
	v := true

	if len(shallow) > 0 {
		v = shallow[0]
	}

	p.options["shallow"] = v

	return p
}

// Register applies the configured options and registers the routes.
func (p *PendingResourceRegistration) Register() *RouteCollection {
	if p.registered != nil {
		return p.registered
	}

	p.registered = p.registrar.Register(p.name, p.controller, p.options)

	return p.registered
}

// Routes is a parity alias for [Register].
func (p *PendingResourceRegistration) Routes() *RouteCollection { return p.Register() }
