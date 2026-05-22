package socialite_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bedrock/packages/socialite"
)

// setupFakeManager creates a manager with a github and google driver for use
// in fake tests.
func setupFakeManager(t *testing.T) *socialite.Manager {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	configs := map[string]socialite.ProviderConfig{
		"github": {ClientID: "gh-id", ClientSecret: "gh-secret", RedirectURL: "http://localhost/callback"},
		"google": {ClientID: "g-id", ClientSecret: "g-secret", RedirectURL: "http://localhost/callback"},
	}
	m := socialite.NewManager(req, session, configs)
	socialite.SetManager(m)
	t.Cleanup(socialite.ClearResolvedInstances)

	return m
}

// SocialiteFakeTest::test_it_can_fake_a_driver_with_a_user
func TestFakeDriverWithUser(t *testing.T) {
	m := setupFakeManager(t)

	user := (&socialite.User{}).Map(map[string]any{
		"id":    "123",
		"name":  "Test User",
		"email": "test@example.com",
	})

	m.Fake("github", user)

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver() returned error: %v", err)
	}

	if _, ok := p.(*socialite.FakeProvider); !ok {
		t.Error("expected Driver() to return *FakeProvider after Fake()")
	}

	u, err := p.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if u.GetID() != "123" {
		t.Errorf("expected ID '123', got %q", u.GetID())
	}

	if u.GetName() != "Test User" {
		t.Errorf("expected Name 'Test User', got %q", u.GetName())
	}

	if u.GetEmail() != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got %q", u.GetEmail())
	}
}

// SocialiteFakeTest::test_it_can_fake_a_driver_with_a_closure
func TestFakeDriverWithClosure(t *testing.T) {
	m := setupFakeManager(t)

	m.FakeWith("github", func() *socialite.User {
		return (&socialite.User{}).Map(map[string]any{"id": "456", "name": "Closure User"})
	})

	p, _ := m.Driver("github")
	u, err := p.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if u.GetID() != "456" {
		t.Errorf("expected ID '456', got %q", u.GetID())
	}

	if u.GetName() != "Closure User" {
		t.Errorf("expected Name 'Closure User', got %q", u.GetName())
	}
}

// SocialiteFakeTest::test_it_can_fake_multiple_drivers
func TestFakeMultipleDrivers(t *testing.T) {
	m := setupFakeManager(t)

	m.Fake("github", (&socialite.User{}).Map(map[string]any{"id": "github-123"}))
	m.Fake("google", (&socialite.User{}).Map(map[string]any{"id": "google-456"}))

	pGH, _ := m.Driver("github")
	uGH, _ := pGH.User(context.Background())

	if uGH.GetID() != "github-123" {
		t.Errorf("expected github ID 'github-123', got %q", uGH.GetID())
	}

	pG, _ := m.Driver("google")
	uG, _ := pG.User(context.Background())

	if uG.GetID() != "google-456" {
		t.Errorf("expected google ID 'google-456', got %q", uG.GetID())
	}
}

// SocialiteFakeTest::test_it_returns_fake_redirect_response
func TestFakeReturnsRedirectURL(t *testing.T) {
	m := setupFakeManager(t)
	m.Fake("github", (&socialite.User{}).Map(map[string]any{"id": "123"}))

	p, _ := m.Driver("github")
	redirectURL, err := p.Redirect(context.Background())

	if err != nil {
		t.Fatalf("Redirect() returned error: %v", err)
	}

	if redirectURL != "https://socialite.fake/github/authorize" {
		t.Errorf("expected fake redirect URL, got %q", redirectURL)
	}
}

// SocialiteFakeTest::test_it_forwards_calls_to_the_real_provider_methods
// test_it_forwards_calls_to_the_real_provider_methods.
func TestFakeForwardsChainedCalls(t *testing.T) {
	m := setupFakeManager(t)
	fp := m.Fake("github", (&socialite.User{}).Map(map[string]any{"id": "123"}))

	// These should not panic and should return the same FakeProvider.
	fp.Stateless()
	fp.Scopes([]string{"user", "repo"})
	fp.SetScopes([]string{"user:email"})
	fp.RedirectURL("http://example.com/callback")
	fp.With(map[string]string{"custom": "param"})
	fp.EnablePKCE()

	// User() still returns the preset user.
	p, _ := m.Driver("github")
	u, err := p.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if u.GetID() != "123" {
		t.Errorf("expected ID '123' after chaining, got %q", u.GetID())
	}
}

// SocialiteFakeTest::test_it_preserves_decorator_pattern_when_chaining_methods
// test_it_preserves_decorator_pattern_when_chaining_methods.
func TestFakePreservesDecoratorPattern(t *testing.T) {
	m := setupFakeManager(t)
	fp := m.Fake("github", (&socialite.User{}).Map(map[string]any{"id": "123"}))

	// Every fluent method must return *FakeProvider (not the real provider).
	chained := fp.Stateless().
		Scopes([]string{"user", "repo"}).
		SetScopes([]string{"user:email"}).
		RedirectURL("http://example.com/callback").
		With(map[string]string{"custom": "param"}).
		EnablePKCE()

	if chained != fp {
		t.Error("chained fluent calls should return the same *FakeProvider instance")
	}

	u, err := chained.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if u.GetID() != "123" {
		t.Errorf("expected ID '123', got %q", u.GetID())
	}
}

// SocialiteFakeTest::test_it_returns_real_driver_when_not_faked
// test_it_returns_real_driver_when_not_faked.
func TestFakeReturnsRealDriverWhenNotFaked(t *testing.T) {
	m := setupFakeManager(t)

	// Fake only github.
	m.Fake("github", (&socialite.User{}).Map(map[string]any{"id": "123"}))

	pGH, _ := m.Driver("github")

	if _, ok := pGH.(*socialite.FakeProvider); !ok {
		t.Error("faked driver 'github' should return *FakeProvider")
	}

	// Google was NOT faked — should return the real provider.
	pG, _ := m.Driver("google")

	if _, ok := pG.(*socialite.FakeProvider); ok {
		t.Error("non-faked driver 'google' should not return *FakeProvider")
	}

	if _, ok := pG.(*socialite.GoogleProvider); !ok {
		t.Errorf("non-faked 'google' driver should return *GoogleProvider, got %T", pG)
	}
}
