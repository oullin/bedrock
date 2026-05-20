package socialite

import (
	"context"
	"fmt"
)

// FakeProvider wraps a real provider and returns a preset user from User(),
// making it easy to test OAuth flows without real HTTP round-trips.
// Its fluent methods forward to the real provider for configuration
// side-effects but always return the FakeProvider so callers can chain.
//
// It mirrors upstream Socialite\Testing\FakeProvider.
type FakeProvider struct {
	driver string
	real   Provider
	user   *User
	userFn func() *User
}

func newFakeProvider(driver string, real Provider, user *User, fn func() *User) *FakeProvider {
	return &FakeProvider{driver: driver, real: real, user: user, userFn: fn}
}

// Redirect returns a deterministic fake authorization URL of the form
// "https://socialite.fake/{driver}/authorize". It mirrors FakeProvider::redirect().
func (f *FakeProvider) Redirect(_ context.Context) (string, error) {
	return fmt.Sprintf("https://socialite.fake/%s/authorize", f.driver), nil
}

// User returns the preset user without making any network calls.
func (f *FakeProvider) User(_ context.Context) (*User, error) {
	if f.userFn != nil {
		return f.userFn(), nil
	}

	return f.user, nil
}

// Stateless forwards the call to the real provider and returns the FakeProvider
// for chaining, mirroring the decorator pattern in SocialiteFakeTest.
func (f *FakeProvider) Stateless() *FakeProvider {
	if ap, ok := f.real.(interface{ Stateless() *AbstractProvider }); ok {
		ap.Stateless()
	}

	return f
}

// Scopes forwards the call to the real provider and returns the FakeProvider.
func (f *FakeProvider) Scopes(scopes []string) *FakeProvider {
	if ap, ok := f.real.(interface {
		Scopes([]string) *AbstractProvider
	}); ok {
		ap.Scopes(scopes)
	}

	return f
}

// SetScopes forwards the call to the real provider and returns the FakeProvider.
func (f *FakeProvider) SetScopes(scopes []string) *FakeProvider {
	if ap, ok := f.real.(interface {
		SetScopes([]string) *AbstractProvider
	}); ok {
		ap.SetScopes(scopes)
	}

	return f
}

// RedirectURL forwards the call to the real provider and returns the FakeProvider.
func (f *FakeProvider) RedirectURL(u string) *FakeProvider {
	if ap, ok := f.real.(interface {
		RedirectURL(string) *AbstractProvider
	}); ok {
		ap.RedirectURL(u)
	}

	return f
}

// With forwards the call to the real provider and returns the FakeProvider.
func (f *FakeProvider) With(params map[string]string) *FakeProvider {
	if ap, ok := f.real.(interface {
		With(map[string]string) *AbstractProvider
	}); ok {
		ap.With(params)
	}

	return f
}

// EnablePKCE forwards the call to the real provider and returns the FakeProvider.
func (f *FakeProvider) EnablePKCE() *FakeProvider {
	if ap, ok := f.real.(interface{ EnablePKCE() *AbstractProvider }); ok {
		ap.EnablePKCE()
	}

	return f
}
