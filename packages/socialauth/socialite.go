package socialauth

import (
	"context"
	"fmt"
)

// SocialAuth is the package-level manager instance. It is populated by
// SetManager() — typically called from the service provider — and provides a
// convenient static API that mirrors the Upstream SocialAuth facade.

// staticFacade wraps a *Manager and exposes a static-style API.
type staticFacade struct {
	manager *Manager
}

var SocialAuth *staticFacade

// SetManager wires the package-level SocialAuth variable to a fully configured
// Manager. Call this from your bootstrap or service provider.
func SetManager(m *Manager) {
	SocialAuth = &staticFacade{manager: m}
}

// ClearResolvedInstances discards all cached and faked provider instances.
// It mirrors the same-named method on the Upstream SocialAuth facade, primarily
// used between tests to isolate state.
func ClearResolvedInstances() {
	if SocialAuth != nil {
		SocialAuth.manager.ForgetDrivers()
	}
}

// Driver resolves the named provider via the package-level manager. It panics
// if no manager has been set (call SetManager first).
func Driver(name string) Provider {
	if SocialAuth == nil {
		panic("socialauth: no manager configured — call socialauth.SetManager first")
	}

	p, err := SocialAuth.manager.Driver(name)

	if err != nil {
		panic(fmt.Sprintf("socialauth: %v", err))
	}

	return p
}

// Fake registers a preset user for the named driver on the package-level
// manager. Returns the FakeProvider so callers can chain fluent configuration.
func Fake(name string, user *User) *FakeProvider {
	if SocialAuth == nil {
		panic("socialauth: no manager configured — call socialauth.SetManager first")
	}

	return SocialAuth.manager.Fake(name, user)
}

// FakeWith registers a closure-based fake for the named driver.
func FakeWith(name string, fn func() *User) *FakeProvider {
	if SocialAuth == nil {
		panic("socialauth: no manager configured — call socialauth.SetManager first")
	}

	return SocialAuth.manager.FakeWith(name, fn)
}

// Redirect is a convenience wrapper: resolves the named driver and calls
// Redirect on it.
func Redirect(ctx context.Context, driver string) (string, error) {
	if SocialAuth == nil {
		panic("socialauth: no manager configured — call socialauth.SetManager first")
	}

	p, err := SocialAuth.manager.Driver(driver)

	if err != nil {
		return "", err
	}

	return p.Redirect(ctx)
}

// GetUser is a convenience wrapper: resolves the named driver and calls User.
func GetUser(ctx context.Context, driver string) (*User, error) {
	if SocialAuth == nil {
		panic("socialauth: no manager configured — call socialauth.SetManager first")
	}

	p, err := SocialAuth.manager.Driver(driver)

	if err != nil {
		return nil, err
	}

	return p.User(ctx)
}
