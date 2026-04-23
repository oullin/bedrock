package socialite_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/bedrock/packages/socialite"
)

// SocialiteManagerTest::test_it_can_instantiate_the_github_driver
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

// SocialiteManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_from_config_array
func TestManagerAppliesConfigScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	configs := map[string]socialite.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
			Scopes:       []string{"read:user", "repo"},
		},
	}
	m := socialite.NewManager(req, session, configs)

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	got := p.(*socialite.GithubProvider).GetScopes()
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

// SocialiteManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_without_array_from_config
func TestManagerAppliesStringConfigScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	configs := map[string]socialite.ProviderConfig{
		"github": {
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://your-callback-url",
			Scopes:       "read:user repo",
		},
	}
	m := socialite.NewManager(req, session, configs)

	p, err := m.Driver("github")

	if err != nil {
		t.Fatalf("Driver('github') returned error: %v", err)
	}

	got := p.(*socialite.GithubProvider).GetScopes()
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

// SocialiteManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_from_config_array_merged_by_programmatic_scopes_using_scopes_method
func TestManagerMergesConfigScopesWithProgrammaticScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialite.NewManager(req, session, map[string]socialite.ProviderConfig{
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

	p.(*socialite.GithubProvider).Scopes([]string{"repo"})

	got := p.(*socialite.GithubProvider).GetScopes()
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

// SocialiteManagerTest::test_it_can_instantiate_the_github_driver_with_scopes_from_config_array_overwritten_by_programmatic_scopes_using_set_scopes_method
func TestManagerOverwritesConfigScopesWithProgrammaticScopes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	m := socialite.NewManager(req, session, map[string]socialite.ProviderConfig{
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

	p.(*socialite.GithubProvider).SetScopes([]string{"workflow"})

	got := p.(*socialite.GithubProvider).GetScopes()
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

// SocialiteManagerTest::test_it_throws_exception_when_client_secret_is_missing
// SocialiteManagerTest::test_it_throws_exception_when_client_id_is_missing
// SocialiteManagerTest::test_it_throws_exception_when_redirect_is_missing
// SocialiteManagerTest::test_it_throws_exception_when_configuration_is_completely_missing
func TestManagerRejectsInvalidGithubConfig(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	cases := []struct {
		name string
		cfg  socialite.ProviderConfig
		want string
	}{
		{
			name: "missing client secret",
			cfg: socialite.ProviderConfig{
				ClientID:    "id",
				RedirectURL: "http://callback",
			},
			want: "client_secret",
		},
		{
			name: "missing client id",
			cfg: socialite.ProviderConfig{
				ClientSecret: "secret",
				RedirectURL:  "http://callback",
			},
			want: "client_id",
		},
		{
			name: "missing redirect",
			cfg: socialite.ProviderConfig{
				ClientID:     "id",
				ClientSecret: "secret",
			},
			want: "redirect_url",
		},
		{
			name: "missing all",
			cfg:  socialite.ProviderConfig{},
			want: "client_id",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := socialite.NewManager(req, newTestSession(), map[string]socialite.ProviderConfig{
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
