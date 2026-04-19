package socialauth

// Factory resolves OAuth provider implementations by driver name.
// It mirrors Upstream\SocialAuth\Contracts\Factory.
type Factory interface {
	Driver(driver string) (Provider, error)
}
