package socialite_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/bedrock/packages/socialauth"
)

// SocialAuthManagerTest::test_it_can_instantiate_the_github_driver
// TestManagerInstantiatesGithubDriver mirrors
// test_it_can_instantiate_the_github_driver (SocialAuthManagerTest.php).
func TestManagerInstantiatesGithubDriver(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	configs := map[string]socialauth.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
		},
	}
	m := socialauth.NewManager(req, session, configs)

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	if _, ok := p.(*socialauth.GithubProvider); !ok {
		t.Errorf("expected *GithubProvider, got %T", p)
	}
}

// SocialAuthManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_from_config_array
func TestManagerAppliesConfigScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	configs := map[string]socialauth.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
			Scopes:       []string{"read:user", "repo"},
		},
	}
	m := socialauth.NewManager(req, session, configs)

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	got := p.(*socialauth.GithubProvider).GetScopes()
	want := []string{"read:user", "repo"}

	if len(got) != len(want) {
		t.Fatalf("unexpected scopes length: got %v want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected scope at %d: got %q want %q", i, got[i], want[i])
		}
	}
}

// SocialAuthManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_without_array_from_config
func TestManagerAppliesStringConfigScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	configs := map[string]socialauth.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
			Scopes:       "read:user repo",
		},
	}
	m := socialauth.NewManager(req, session, configs)

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	got := p.(*socialauth.GithubProvider).GetScopes()
	want := []string{"read:user", "repo"}

	if len(got) != len(want) {
		t.Fatalf("unexpected scopes length: got %v want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected scope at %d: got %q want %q", i, got[i], want[i])
		}
	}
}

// SocialAuthManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_from_config_array_merged_by_programmatic_scopes_using_scopes_method
func TestManagerMergesConfigScopesWithProgrammaticScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialauth.NewManager(req, session, map[string]socialauth.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
			Scopes:       []string{"read:user"},
		},
	})

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	p.(*socialauth.GithubProvider).Scopes([]string{"repo"})

	got := p.(*socialauth.GithubProvider).GetScopes()
	want := []string{"read:user", "repo"}

	if len(got) != len(want) {
		t.Fatalf("unexpected scopes length: got %v want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected scope at %d: got %q want %q", i, got[i], want[i])
		}
	}
}

// SocialAuthManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_from_config_array_overwritten_by_programmatic_scopes_using_set_scopes_method
func TestManagerOverwritesConfigScopesWithProgrammaticScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialauth.NewManager(req, session, map[string]socialauth.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
			Scopes:       []string{"read:user", "repo"},
		},
	})

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	p.(*socialauth.GithubProvider).SetScopes([]string{"workflow"})

	got := p.(*socialauth.GithubProvider).GetScopes()
	want := []string{"workflow"}

	if len(got) != len(want) {
		t.Fatalf("unexpected scopes length: got %v want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected scope at %d: got %q want %q", i, got[i], want[i])
		}
	}
}

// SocialAuthManagerTest::test_it_throws_exception_when_client_secret_is_missing
// SocialAuthManagerTest::test_it_throws_exception_when_client_id_is_missing
// SocialAuthManagerTest::test_it_throws_exception_when_redirect_is_missing
// SocialAuthManagerTest::test_it_throws_exception_when_configuration_is_completely_missing
func TestManagerRejectsInvalidGithubConfig(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	cases := []struct {
		name string
		cfg  socialauth.ProviderConfig
		want string
	}{
		{
			name: "missing client secret",
			cfg: socialauth.ProviderConfig{
				ClientID:    "id",
				RedirectURL: "http://callback",
			},
			want: "client_secret",
		},
		{
			name: "missing client id",
			cfg: socialauth.ProviderConfig{
				ClientSecret: "secret",
				RedirectURL:  "http://callback",
			},
			want: "client_id",
		},
		{
			name: "missing redirect",
			cfg: socialauth.ProviderConfig{
				ClientID:     "id",
				ClientSecret: "secret",
			},
			want: "redirect_url",
		},
		{
			name: "missing all",
			cfg:  socialauth.ProviderConfig{},
			want: "client_id",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := socialauth.NewManager(req, newTestSession(), map[string]socialauth.ProviderConfig{
				"github": tc.cfg,
			})

			_, err := m.Driver("github")

			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error to mention %q, got %v", tc.want, err)
			}
		})
	}
}

// TestManagerCachesResolvedDriver verifies that the same instance is returned
// on repeated Driver() calls (no double-resolution).
func TestManagerCachesResolvedDriver(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialauth.NewManager(req, session, map[string]socialauth.ProviderConfig{
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
	m := socialauth.NewManager(req, session, map[string]socialauth.ProviderConfig{
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
	m := socialauth.NewManager(req, nil, nil)

	_, err := m.Driver("nonexistent")

	if err == nil {
		t.Error("expected error for unknown driver, got nil")
	}
}

// TestManagerExtend verifies custom driver registration via Extend().
func TestManagerExtend(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialauth.NewManager(req, session, map[string]socialauth.ProviderConfig{
		"mydriver": {ClientID: "id", ClientSecret: "secret", RedirectURL: "http://cb"},
	})

	// Register a custom driver that re-uses GithubProvider as a stand-in.
	m.Extend("mydriver", func(req *http.Request, s socialauth.Session, cfg socialauth.ProviderConfig) (socialauth.Provider, error) {
		return socialauth.NewGithubProvider(req, s, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL), nil
	})

	p, err := m.Driver("mydriver")

	if err != nil {
		t.Fatalf("extended driver returned error: %v", err)
	}

	if p == nil {
		t.Error("expected non-nil provider from extended driver")
	}
}
