package socialauth

// Factory resolves OAuth provider implementations by driver name.
type Factory interface {
	Driver(driver string) (Provider, error)
}
