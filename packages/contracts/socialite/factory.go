package socialite

// Factory resolves OAuth provider implementations by driver name.
// It mirrors upstream Socialite\Contracts\Factory.
type Factory interface {
	Driver(driver string) (Provider, error)
}
