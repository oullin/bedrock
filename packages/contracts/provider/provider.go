package provider

// ServiceProvider is implemented by packages that register their services
// with the application container. It mirrors @bedrock\Support\ServiceProvider.
type ServiceProvider interface {
	Register()
}

// Bootable is implemented by providers that require post-registration setup,
// such as subscribing event listeners or configuring guards. Boot is called
// after all providers have been registered.
type Bootable interface {
	Boot()
}

// Provides is an optional hint from a provider declaring which container
// abstract keys it binds. Consumers can use this for introspection.
type Provides interface {
	Provides() []string
}

// Deferred is a marker interface. When a provider implements both Deferred
// and Provides, the Application defers calling Register() until one of the
// abstract keys returned by Provides() is actually resolved. This avoids
// paying the cost of constructing services that may never be used.
//
// Deferred providers MUST also implement Provides — otherwise the
// Application has no way to know which keys to watch.
type Deferred interface {
	Deferred() bool
}

// DependsOn is an optional hint declaring abstract keys that must be
// registered before this provider's Register() runs. The Application uses
// this to topologically sort providers at registration time.
//
// Cycles are reported as an error/panic. Missing dependencies are
// permitted: if you depend on "foo" but no provider binds "foo", you will
// be registered after every provider that does declare DependsOn.
type DependsOn interface {
	DependsOn() []string
}
