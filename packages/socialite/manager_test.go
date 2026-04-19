package socialite_test

import (
	"net/http"
	"testing"

	"github.com/bedrock/packages/socialite"
)

// TestManagerInstantiatesGithubDriver mirrors
// test_it_can_instantiate_the_github_driver (SocialiteManagerTest.php).
func TestManagerInstantiatesGithubDriver(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	configs := map[string]socialite.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
		},
	}
	m := socialite.NewManager(req, session, configs)

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	if _, ok := p.(*socialite.GithubProvider); !ok {
		t.Errorf("expected *GithubProvider, got %T", p)
	}
}

// TestManagerCachesResolvedDriver verifies that the same instance is returned
// on repeated Driver() calls (no double-resolution).
func TestManagerCachesResolvedDriver(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialite.NewManager(req, session, map[string]socialite.ProviderConfig{
		"github": {ClientID: "id", ClientSecret: "secret", RedirectURL: "http://cb"},
	})

	p1, _ := m.Driver("github")
	p2, _ := m.Driver("github")

	if p1 != p2 {
		t.Error("expected Driver() to return cached instance on second call")
	}
}

// TestManagerForgetDrivers verifies ForgetDrivers clears the cache.
func TestManagerForgetDrivers(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialite.NewManager(req, session, map[string]socialite.ProviderConfig{
		"github": {ClientID: "id", ClientSecret: "secret", RedirectURL: "http://cb"},
	})

	p1, _ := m.Driver("github")
	m.ForgetDrivers()
	p2, _ := m.Driver("github")

	if p1 == p2 {
		t.Error("expected fresh instance after ForgetDrivers()")
	}
}

// TestManagerUnknownDriverReturnsError verifies error on unknown driver.
func TestManagerUnknownDriverReturnsError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	m := socialite.NewManager(req, nil, nil)

	_, err := m.Driver("nonexistent")

	if err == nil {
		t.Error("expected error for unknown driver, got nil")
	}
}

// TestManagerExtend verifies custom driver registration via Extend().
func TestManagerExtend(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialite.NewManager(req, session, map[string]socialite.ProviderConfig{
		"mydriver": {ClientID: "id", ClientSecret: "secret", RedirectURL: "http://cb"},
	})

	// Register a custom driver that re-uses GithubProvider as a stand-in.
	m.Extend("mydriver", func(req *http.Request, s socialite.Session, cfg socialite.ProviderConfig) (socialite.Provider, error) {
		return socialite.NewGithubProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	})

	p, err := m.Driver("mydriver")

	if err != nil {
		t.Fatalf("extended driver returned error: %v", err)
	}

	if p == nil {
		t.Error("expected non-nil provider from extended driver")
	}
}
