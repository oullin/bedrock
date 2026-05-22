package oauthserver_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/oauthserver"
)

type inventoryAuthServer struct {
	issued oauthserver.IssuedToken
	req    oauthserver.TokenRequest
	err    error
}

type inventoryDispatcher struct {
	events []any
}

type inventoryCountingProvider struct {
	users map[string]cauth.Authenticatable
	calls int
}

type inventoryErrorTokenStore struct {
	err error
}

type inventoryCountingClientStore struct {
	client *oauthserver.Client
	calls  int
	err    error
}

type inventoryErrorClientStore struct {
	err error
}

func (s *inventoryAuthServer) IssueToken(_ context.Context, req oauthserver.TokenRequest) (*oauthserver.IssuedToken, error) {
	s.req = req

	if s.err != nil {
		return nil, s.err
	}

	return &s.issued, nil
}

func (s *inventoryAuthServer) ValidateAuthorizationRequest(context.Context, oauthserver.AuthorizationRequest) (*oauthserver.ValidatedAuthRequest, error) {
	return nil, errors.New("not implemented")
}

func (s *inventoryAuthServer) CompleteAuthorizationRequest(context.Context, *oauthserver.ValidatedAuthRequest, bool) (*oauthserver.AuthorizationResponse, error) {
	return nil, errors.New("not implemented")
}

func (d *inventoryDispatcher) Dispatch(_ context.Context, event any) ([]any, error) {
	d.events = append(d.events, event)

	return nil, nil
}

func requireInventoryScopeIDs(t *testing.T, scopes []oauthserver.Scope, want []string) {
	t.Helper()

	if len(scopes) != len(want) {
		t.Fatalf("scope count = %d, want %d: %#v", len(scopes), len(want), scopes)
	}

	for i, scope := range scopes {
		if scope.ID != want[i] {
			t.Fatalf("scope[%d] = %q, want %q: %#v", i, scope.ID, want[i], scopes)
		}
	}
}

func (p *inventoryCountingProvider) RetrieveByID(_ context.Context, id string) (cauth.Authenticatable, error) {
	p.calls++

	return p.users[id], nil
}

func (p *inventoryCountingProvider) RetrieveByToken(context.Context, string, string) (cauth.Authenticatable, error) {
	return nil, nil
}

func (p *inventoryCountingProvider) RetrieveByCredentials(context.Context, map[string]string) (cauth.Authenticatable, error) {
	return nil, nil
}

func (p *inventoryCountingProvider) UpdateRememberToken(context.Context, cauth.Authenticatable, string) error {
	return nil
}

func (p *inventoryCountingProvider) ValidateCredentials(context.Context, cauth.Authenticatable, map[string]string) (bool, error) {
	return false, nil
}

func (p *inventoryCountingProvider) RehashPasswordIfRequired(context.Context, cauth.Authenticatable, map[string]string, bool) error {
	return nil
}

func (s *inventoryErrorTokenStore) Find(context.Context, string) (*oauthserver.Token, error) {
	return nil, s.err
}

func (s *inventoryErrorTokenStore) FindForUser(context.Context, string, string) (*oauthserver.Token, error) {
	return nil, s.err
}

func (s *inventoryErrorTokenStore) Save(context.Context, *oauthserver.Token) error {
	return nil
}

func (s *inventoryErrorTokenStore) Revoke(context.Context, string) error {
	return nil
}

func (s *inventoryErrorTokenStore) ForUser(context.Context, string) ([]*oauthserver.Token, error) {
	return nil, s.err
}

func (s *inventoryCountingClientStore) Find(_ context.Context, id string) (*oauthserver.Client, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.client == nil || s.client.ID != id {
		return nil, nil
	}

	copy := *s.client

	return &copy, nil
}

func (s *inventoryCountingClientStore) FindActive(ctx context.Context, id string) (*oauthserver.Client, error) {
	s.calls++

	return s.Find(ctx, id)
}

func (s *inventoryCountingClientStore) PersonalAccessClient(context.Context) (*oauthserver.Client, error) {
	return nil, nil
}

func (s *inventoryCountingClientStore) Create(context.Context, *oauthserver.Client) error {
	return nil
}

func (s *inventoryCountingClientStore) Delete(context.Context, string) error {
	return nil
}

func (s *inventoryCountingClientStore) RegenerateSecret(context.Context, string) (string, error) {
	return "", nil
}

func (s *inventoryCountingClientStore) CreatePersonalAccessClient(context.Context, string, string, string) (*oauthserver.Client, error) {
	return nil, nil
}

func (s *inventoryCountingClientStore) CreatePasswordGrantClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, nil
}

func (s *inventoryCountingClientStore) CreateClientCredentialsClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, nil
}

func (s *inventoryCountingClientStore) CreateImplicitClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, nil
}

func (s *inventoryCountingClientStore) CreateDeviceCodeGrantClient(context.Context, string, string, string) (*oauthserver.Client, error) {
	return nil, nil
}

func (s *inventoryCountingClientStore) CreateAuthCodeClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, nil
}

func (s *inventoryErrorClientStore) Find(context.Context, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) FindActive(context.Context, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) PersonalAccessClient(context.Context) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) Create(context.Context, *oauthserver.Client) error {
	return nil
}

func (s *inventoryErrorClientStore) Delete(context.Context, string) error {
	return nil
}

func (s *inventoryErrorClientStore) RegenerateSecret(context.Context, string) (string, error) {
	return "", s.err
}

func (s *inventoryErrorClientStore) CreatePersonalAccessClient(context.Context, string, string, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) CreatePasswordGrantClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) CreateClientCredentialsClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) CreateImplicitClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) CreateDeviceCodeGrantClient(context.Context, string, string, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func (s *inventoryErrorClientStore) CreateAuthCodeClient(context.Context, string, string, string, string) (*oauthserver.Client, error) {
	return nil, s.err
}

func TestOAuthServerInventoryClientAndModelPrimitives(t *testing.T) {
	t.Run("ClientTest::testScopesWhenClientHasScope", func(t *testing.T) {
		client := &oauthserver.Client{Scopes: []string{"orders:read", "users:write"}}

		if !client.HasScope("orders:read") {
			t.Fatal("expected client to have orders:read scope")
		}
	})

	t.Run("ClientTest::testScopesWhenClientDoesNotHaveScope", func(t *testing.T) {
		client := &oauthserver.Client{Scopes: []string{"orders:read"}}

		if client.HasScope("admin") {
			t.Fatal("expected missing client scope to be false")
		}
	})

	t.Run("ClientTest::testScopesWhenColumnDoesNotExist", func(t *testing.T) {
		client := &oauthserver.Client{}

		if client.HasScope("admin") {
			t.Fatal("expected absent scope list to deny scope")
		}
	})

	t.Run("ClientTest::testScopesWhenColumnIsNull", func(t *testing.T) {
		client := &oauthserver.Client{Scopes: nil}

		if client.HasScope("admin") {
			t.Fatal("expected nil scope list to deny scope")
		}
	})

	t.Run("ClientTest::testGrantTypesWhenClientHasGrantType", func(t *testing.T) {
		client := &oauthserver.Client{GrantTypes: []string{oauthserver.GrantAuthorizationCode, oauthserver.GrantRefreshToken}}

		if !client.HasGrantType(oauthserver.GrantRefreshToken) {
			t.Fatal("expected client to have refresh_token grant")
		}
	})

	t.Run("ClientTest::testGrantTypesWhenClientDoesNotHaveGrantType", func(t *testing.T) {
		client := &oauthserver.Client{GrantTypes: []string{oauthserver.GrantAuthorizationCode}}

		if client.HasGrantType(oauthserver.GrantPassword) {
			t.Fatal("expected missing grant type to be false")
		}
	})

	t.Run("ClientTest::testGrantTypesWhenColumnDoesNotExist", func(t *testing.T) {
		client := &oauthserver.Client{}

		if client.HasGrantType(oauthserver.GrantPassword) {
			t.Fatal("expected absent grant type list to deny grant")
		}
	})

	t.Run("ClientTest::testGrantTypesWhenColumnIsNull", func(t *testing.T) {
		client := &oauthserver.Client{GrantTypes: nil}

		if client.HasGrantType(oauthserver.GrantPassword) {
			t.Fatal("expected nil grant type list to deny grant")
		}
	})

	t.Run("OAuthServerTest::test_scopes_can_be_managed", func(t *testing.T) {
		p := oauthserver.NewOAuthServer(nil).TokensCan(map[string]string{
			"orders:read": "Read orders",
			"orders:pay":  "Pay orders",
		})

		if !p.HasScope("orders:read") || p.FindScope("orders:pay") == nil || len(p.ScopeIDs()) != 2 {
			t.Fatalf("expected registered scopes to be discoverable")
		}
	})

	// ScopeControllerTest::testShouldGetScopes
	t.Run("ScopeControllerTest::testShouldGetScopes", func(t *testing.T) {
		p := oauthserver.NewOAuthServer(nil).TokensCan(map[string]string{
			"orders:read": "Read orders",
			"orders:pay":  "Pay orders",
		})

		scopes := p.Scopes()

		if len(scopes) != 2 {
			t.Fatalf("scope count = %d, want 2: %#v", len(scopes), scopes)
		}

		if !p.HasScope("orders:read") || !p.HasScope("orders:pay") {
			t.Fatalf("expected registered scopes to be listed: %#v", scopes)
		}
	})

	// ClientCredentialsGrantTest::testCustomClientCredentialsTokenExpiration
	t.Run("ClientCredentialsGrantTest::testCustomClientCredentialsTokenExpiration", func(t *testing.T) {
		p := oauthserver.NewOAuthServer(nil).ClientCredentialsTokensExpireIn(15 * time.Minute)

		if p.ClientCredentialsTTL() != 15*time.Minute {
			t.Fatalf("client credentials TTL = %s, want %s", p.ClientCredentialsTTL(), 15*time.Minute)
		}
	})

	t.Run("OAuthServerTest::test_auth_code_instance_can_be_created", func(t *testing.T) {
		code := oauthserver.AuthCode{ID: "auth-code", UserID: "user-1", ClientID: "client-1"}

		if code.ID != "auth-code" || code.UserID != "user-1" || code.ClientID != "client-1" {
			t.Fatalf("unexpected auth code: %#v", code)
		}
	})

	t.Run("OAuthServerTest::test_client_instance_can_be_created", func(t *testing.T) {
		client := oauthserver.Client{ID: "client-1", Name: "CLI", Secret: "secret"}

		if client.ID != "client-1" || !client.Confidential() {
			t.Fatalf("unexpected client: %#v", client)
		}
	})

	t.Run("OAuthServerTest::test_token_instance_can_be_created", func(t *testing.T) {
		token := oauthserver.Token{ID: "token-1", UserID: "user-1", Scopes: []string{"read"}}

		if token.ID != "token-1" || !token.Can("read") {
			t.Fatalf("unexpected token: %#v", token)
		}
	})

	t.Run("OAuthServerTest::test_refresh_token_instance_can_be_created", func(t *testing.T) {
		token := oauthserver.RefreshToken{ID: "refresh-1", AccessTokenID: "access-1"}

		if token.ID != "refresh-1" || token.AccessTokenID != "access-1" {
			t.Fatalf("unexpected refresh token: %#v", token)
		}
	})

	t.Run("OAuthServerTest::test_refresh_token_model_can_be_changed", func(t *testing.T) {
		token := oauthserver.RefreshToken{ID: "refresh-1", AccessTokenID: "access-1"}
		token.Revoked = true

		if !token.Revoked {
			t.Fatalf("expected refresh token to be mutable: %#v", token)
		}
	})

	t.Run("OAuthServerTest::test_device_code_instance_can_be_created", func(t *testing.T) {
		code := oauthserver.DeviceCode{ID: "device-1", DeviceCode: "device-code", UserCode: "ABCD-EFGH"}

		if code.ID != "device-1" || code.DeviceCode != "device-code" || code.UserCode != "ABCD-EFGH" {
			t.Fatalf("unexpected device code: %#v", code)
		}
	})
}

func TestOAuthServerInventoryTokenAndScopeParity(t *testing.T) {
	t.Run("ScopeTest::test_scope_can_be_converted_to_array", func(t *testing.T) {
		scope := oauthserver.Scope{ID: "orders:read", Description: "Read orders"}
		arr := scope.ToArray()

		if arr["id"] != "orders:read" || arr["description"] != "Read orders" {
			t.Fatalf("unexpected scope array: %#v", arr)
		}
	})

	t.Run("ScopeTest::test_scope_can_be_converted_to_json", func(t *testing.T) {
		scope := oauthserver.Scope{ID: "orders:read", Description: "Read orders"}
		payload, err := json.Marshal(scope)

		if err != nil {
			t.Fatal(err)
		}

		var arr map[string]string

		if err := json.Unmarshal(payload, &arr); err != nil {
			t.Fatal(err)
		}

		if arr["id"] != "orders:read" || arr["description"] != "Read orders" {
			t.Fatalf("unexpected scope JSON: %#v", arr)
		}
	})

	t.Run("AccessTokenTest::test_token_attributes_are_accessible", func(t *testing.T) {
		now := time.Now().Truncate(time.Second)
		token := &oauthserver.Token{
			ID:        "token-1",
			UserID:    "user-1",
			ClientID:  "client-1",
			Name:      "CI",
			Scopes:    []string{"read"},
			CreatedAt: now,
			UpdatedAt: now,
			ExpiresAt: now.Add(time.Hour),
		}
		arr := oauthserver.NewAccessToken(token, newOAuthServer()).ToArray()

		if arr["id"] != "token-1" || arr["user_id"] != "user-1" || arr["client_id"] != "client-1" {
			t.Fatalf("unexpected access token attributes: %#v", arr)
		}
	})

	t.Run("AccessTokenTest::test_token_can_determine_if_it_has_scopes", func(t *testing.T) {
		token := newTestToken("token-1", "user-1", "client-1", []string{"orders:read"})
		accessToken := oauthserver.NewAccessToken(token, newOAuthServer())

		if !accessToken.Can("orders:read") || accessToken.Can("orders:write") {
			t.Fatal("expected access token scope checks to match token scopes")
		}
	})

	t.Run("AccessTokenTest::test_token_can_determine_if_it_has_inherited_scopes", func(t *testing.T) {
		p := newOAuthServer().UseInheritedScopes(true)
		token := newTestToken("token-1", "user-1", "client-1", []string{"orders"})
		token.WithOAuthServer(p)
		accessToken := oauthserver.NewAccessToken(token, p)

		if !accessToken.Can("orders:read") {
			t.Fatal("expected inherited parent scope to satisfy child scope")
		}
	})

	t.Run("AccessTokenTest::test_token_resolves_inherited_scopes", func(t *testing.T) {
		p := newOAuthServer().UseInheritedScopes(true)
		token := newTestToken("token-1", "user-1", "client-1", []string{"admin:webhooks"})
		token.WithOAuthServer(p)
		accessToken := oauthserver.NewAccessToken(token, p)

		if !accessToken.Can("admin:webhooks:read") || accessToken.Can("admin:billing:read") {
			t.Fatal("expected inherited scope resolution to use colon ancestors")
		}
	})

	t.Run("TokenTest::test_token_can_determine_if_it_has_scopes", func(t *testing.T) {
		token := newTestToken("token-1", "user-1", "client-1", []string{"orders:read"})

		if !token.Can("orders:read") || token.Can("orders:write") {
			t.Fatal("expected token scope checks to match token scopes")
		}
	})

	t.Run("TokenTest::test_token_can_determine_if_it_has_inherited_scopes", func(t *testing.T) {
		p := newOAuthServer().UseInheritedScopes(true)
		token := newTestToken("token-1", "user-1", "client-1", []string{"orders"}).WithOAuthServer(p)

		if !token.Can("orders:read") {
			t.Fatal("expected inherited scope check to pass")
		}
	})

	t.Run("TokenTest::test_token_resolves_inherited_scopes", func(t *testing.T) {
		p := newOAuthServer().UseInheritedScopes(true)
		token := newTestToken("token-1", "user-1", "client-1", []string{"admin:webhooks"}).WithOAuthServer(p)

		if !token.Can("admin:webhooks:read") || token.Can("admin:billing:read") {
			t.Fatal("expected inherited scope resolution to use colon ancestors")
		}
	})

	t.Run("TransientTokenTest::test_transient_token_can_do_anything", func(t *testing.T) {
		token := newTestToken("token-1", "user-1", "client-1", []string{"*"})

		if !token.Can("any:scope") || token.Can("*") {
			t.Fatal("expected wildcard token to grant all concrete scopes but not literal wildcard checks")
		}
	})

	// HasApiTokensTest::test_token_can_indicates_if_token_has_given_scope
	t.Run("HasApiTokensTest::test_token_can_indicates_if_token_has_given_scope", func(t *testing.T) {
		token := oauthserver.NewAccessToken(newTestToken("token-1", "user-1", "client-1", []string{"orders:read"}), newOAuthServer())
		user := oauthserver.NewUserWithTokens(newStubUser("user-1"), token)

		if !user.TokenCan("orders:read") || user.TokenCan("orders:write") {
			t.Fatal("expected HasApiTokens scope checks to delegate to the current access token")
		}
	})
}

func TestOAuthServerInventoryScopeRepositoryParity(t *testing.T) {
	ctx := context.Background()

	// BridgeScopeRepositoryTest::test_invalid_scopes_are_removed
	t.Run("BridgeScopeRepositoryTest::test_invalid_scopes_are_removed", func(t *testing.T) {
		p := newOAuthServer().TokensCan(map[string]string{"orders:read": "Read orders"})
		repo := oauthserver.NewScopeRepository(p, nil)

		scopes, err := repo.FinalizeScopes(ctx, []string{"orders:read", "missing"}, oauthserver.GrantAuthorizationCode, "")

		if err != nil {
			t.Fatal(err)
		}

		requireInventoryScopeIDs(t, scopes, []string{"orders:read"})
	})

	// BridgeScopeRepositoryTest::test_invalid_scopes_are_removed_without_a_client_repository
	t.Run("BridgeScopeRepositoryTest::test_invalid_scopes_are_removed_without_a_client_repository", func(t *testing.T) {
		p := newOAuthServer().TokensCan(map[string]string{"orders:read": "Read orders"})
		repo := oauthserver.NewScopeRepository(p, nil)

		scopes, err := repo.FinalizeScopes(ctx, []string{"orders:read", "admin"}, oauthserver.GrantPassword, "client-1")

		if err != nil {
			t.Fatal(err)
		}

		requireInventoryScopeIDs(t, scopes, []string{"orders:read"})
	})

	// BridgeScopeRepositoryTest::test_clients_do_not_restrict_scopes_by_default
	t.Run("BridgeScopeRepositoryTest::test_clients_do_not_restrict_scopes_by_default", func(t *testing.T) {
		p := newOAuthServer().TokensCan(map[string]string{
			"orders:read":  "Read orders",
			"orders:write": "Write orders",
		})
		clients := oauthserver.NewMemoryClientStore()

		if err := clients.Create(ctx, &oauthserver.Client{ID: "client-1", Name: "Web"}); err != nil {
			t.Fatal(err)
		}

		repo := oauthserver.NewScopeRepository(p, clients)

		scopes, err := repo.FinalizeScopes(ctx, []string{"orders:read", "orders:write"}, oauthserver.GrantAuthorizationCode, "client-1")

		if err != nil {
			t.Fatal(err)
		}

		requireInventoryScopeIDs(t, scopes, []string{"orders:read", "orders:write"})
	})

	// BridgeScopeRepositoryTest::test_scopes_disallowed_for_client_are_removed
	t.Run("BridgeScopeRepositoryTest::test_scopes_disallowed_for_client_are_removed", func(t *testing.T) {
		p := newOAuthServer().TokensCan(map[string]string{
			"orders:read":  "Read orders",
			"orders:write": "Write orders",
		})
		clients := oauthserver.NewMemoryClientStore()

		if err := clients.Create(ctx, &oauthserver.Client{ID: "client-1", Name: "Web", Scopes: []string{"orders:read"}}); err != nil {
			t.Fatal(err)
		}

		repo := oauthserver.NewScopeRepository(p, clients)

		scopes, err := repo.FinalizeScopes(ctx, []string{"orders:read", "orders:write"}, oauthserver.GrantAuthorizationCode, "client-1")

		if err != nil {
			t.Fatal(err)
		}

		requireInventoryScopeIDs(t, scopes, []string{"orders:read"})
	})

	// BridgeScopeRepositoryTest::test_scopes_disallowed_for_client_are_removed_but_inherited_scopes_are_not
	t.Run("BridgeScopeRepositoryTest::test_scopes_disallowed_for_client_are_removed_but_inherited_scopes_are_not", func(t *testing.T) {
		p := newOAuthServer().
			TokensCan(map[string]string{"orders:read": "Read orders"}).
			UseInheritedScopes(true)
		clients := oauthserver.NewMemoryClientStore()

		if err := clients.Create(ctx, &oauthserver.Client{ID: "client-1", Name: "Web", Scopes: []string{"orders"}}); err != nil {
			t.Fatal(err)
		}

		repo := oauthserver.NewScopeRepository(p, clients)

		scopes, err := repo.FinalizeScopes(ctx, []string{"orders:read"}, oauthserver.GrantAuthorizationCode, "client-1")

		if err != nil {
			t.Fatal(err)
		}

		requireInventoryScopeIDs(t, scopes, []string{"orders:read"})
	})

	// BridgeScopeRepositoryTest::test_superuser_scope_cant_be_applied_if_wrong_grant
	t.Run("BridgeScopeRepositoryTest::test_superuser_scope_cant_be_applied_if_wrong_grant", func(t *testing.T) {
		p := newOAuthServer()
		clients := oauthserver.NewMemoryClientStore()

		if err := clients.Create(ctx, &oauthserver.Client{ID: "client-1", Name: "Web", Scopes: []string{"*"}}); err != nil {
			t.Fatal(err)
		}

		repo := oauthserver.NewScopeRepository(p, clients)

		scopes, err := repo.FinalizeScopes(ctx, []string{"*"}, oauthserver.GrantAuthorizationCode, "client-1")

		if err != nil {
			t.Fatal(err)
		}

		requireInventoryScopeIDs(t, scopes, nil)
	})

	// BridgeScopeRepositoryTest::test_superuser_scope_cant_be_applied_if_wrong_grant_without_a_client_repository
	t.Run("BridgeScopeRepositoryTest::test_superuser_scope_cant_be_applied_if_wrong_grant_without_a_client_repository", func(t *testing.T) {
		repo := oauthserver.NewScopeRepository(newOAuthServer(), nil)

		scopes, err := repo.FinalizeScopes(ctx, []string{"*"}, oauthserver.GrantPassword, "")

		if err != nil {
			t.Fatal(err)
		}

		requireInventoryScopeIDs(t, scopes, nil)
	})
}

func TestOAuthServerInventoryStoreParity(t *testing.T) {
	ctx := context.Background()

	t.Run("BridgeAccessTokenRepositoryTest::test_access_tokens_can_be_persisted", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()
		token := newTestToken("token-1", "user-1", "client-1", []string{"read"})

		if err := store.Save(ctx, token); err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, "token-1")

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || found.ID != "token-1" || found.UserID != "user-1" {
			t.Fatalf("unexpected token: %#v", found)
		}
	})

	t.Run("BridgeAccessTokenRepositoryTest::test_access_tokens_can_be_revoked", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()
		token := newTestToken("token-1", "user-1", "client-1", []string{"read"})

		if err := store.Save(ctx, token); err != nil {
			t.Fatal(err)
		}

		if err := store.Revoke(ctx, "token-1"); err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, "token-1")

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || !found.IsRevoked() {
			t.Fatalf("expected revoked token, got %#v", found)
		}
	})

	t.Run("BridgeAccessTokenRepositoryTest::test_access_token_revoke_event_is_not_dispatched_when_nothing_happened", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()

		if err := store.Revoke(ctx, "missing-token"); err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, "missing-token")

		if err != nil {
			t.Fatal(err)
		}

		if found != nil {
			t.Fatalf("expected missing token to remain absent, got %#v", found)
		}
	})

	t.Run("BridgeAccessTokenRepositoryTest::test_can_get_new_access_token", func(t *testing.T) {
		token := newTestToken("token-1", "user-1", "client-1", []string{"read"})

		if token.ID == "" || token.ExpiresAt.IsZero() || token.IsRevoked() {
			t.Fatalf("unexpected new token: %#v", token)
		}
	})

	t.Run("BridgeRefreshTokenRepositoryTest::test_access_tokens_can_be_persisted", func(t *testing.T) {
		store := oauthserver.NewMemoryRefreshTokenStore()
		refresh := &oauthserver.RefreshToken{ID: "refresh-1", AccessTokenID: "token-1", ExpiresAt: time.Now().Add(time.Hour)}

		if err := store.Save(ctx, refresh); err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, "refresh-1")

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || found.AccessTokenID != "token-1" {
			t.Fatalf("unexpected refresh token: %#v", found)
		}
	})

	t.Run("BridgeRefreshTokenRepositoryTest::test_can_get_new_refresh_token", func(t *testing.T) {
		refresh := &oauthserver.RefreshToken{ID: "refresh-1", AccessTokenID: "token-1", ExpiresAt: time.Now().Add(time.Hour)}

		if refresh.ID == "" || refresh.AccessTokenID == "" || refresh.ExpiresAt.IsZero() || refresh.Revoked {
			t.Fatalf("unexpected refresh token: %#v", refresh)
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_access_token_is_revoked", func(t *testing.T) {
		token := newRevokedToken("token-1", "user-1", "client-1", nil)

		if !token.IsRevoked() {
			t.Fatal("expected access token to be revoked")
		}
	})

	t.Run("RevokedTest::test_a_access_token_is_also_revoked_if_it_cannot_be_found", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()
		token, err := store.Find(ctx, "missing")

		if err != nil {
			t.Fatal(err)
		}

		if token != nil {
			t.Fatalf("expected missing access token to resolve nil, got %#v", token)
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_access_token_is_not_revoked", func(t *testing.T) {
		token := newTestToken("token-1", "user-1", "client-1", nil)

		if token.IsRevoked() {
			t.Fatal("expected access token not to be revoked")
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_auth_code_is_revoked", func(t *testing.T) {
		store := oauthserver.NewMemoryAuthCodeStore()

		if err := store.Save(ctx, &oauthserver.AuthCode{ID: "code-1"}); err != nil {
			t.Fatal(err)
		}

		if err := store.Revoke(ctx, "code-1"); err != nil {
			t.Fatal(err)
		}

		revoked, err := store.IsRevoked(ctx, "code-1")

		if err != nil {
			t.Fatal(err)
		}

		if !revoked {
			t.Fatal("expected auth code to be revoked")
		}
	})

	t.Run("RevokedTest::test_a_auth_code_is_also_revoked_if_it_cannot_be_found", func(t *testing.T) {
		revoked, err := oauthserver.NewMemoryAuthCodeStore().IsRevoked(ctx, "missing")

		if err != nil {
			t.Fatal(err)
		}

		if !revoked {
			t.Fatal("expected missing auth code to be treated as revoked")
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_auth_code_is_not_revoked", func(t *testing.T) {
		store := oauthserver.NewMemoryAuthCodeStore()

		if err := store.Save(ctx, &oauthserver.AuthCode{ID: "code-1"}); err != nil {
			t.Fatal(err)
		}

		revoked, err := store.IsRevoked(ctx, "code-1")

		if err != nil {
			t.Fatal(err)
		}

		if revoked {
			t.Fatal("expected auth code not to be revoked")
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_refresh_token_is_revoked", func(t *testing.T) {
		store := oauthserver.NewMemoryRefreshTokenStore()

		if err := store.Save(ctx, &oauthserver.RefreshToken{ID: "refresh-1"}); err != nil {
			t.Fatal(err)
		}

		if err := store.Revoke(ctx, "refresh-1"); err != nil {
			t.Fatal(err)
		}

		revoked, err := store.IsRevoked(ctx, "refresh-1")

		if err != nil {
			t.Fatal(err)
		}

		if !revoked {
			t.Fatal("expected refresh token to be revoked")
		}
	})

	t.Run("RevokedTest::test_a_refresh_token_is_also_revoked_if_it_cannot_be_found", func(t *testing.T) {
		revoked, err := oauthserver.NewMemoryRefreshTokenStore().IsRevoked(ctx, "missing")

		if err != nil {
			t.Fatal(err)
		}

		if !revoked {
			t.Fatal("expected missing refresh token to be treated as revoked")
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_refresh_token_is_not_revoked", func(t *testing.T) {
		store := oauthserver.NewMemoryRefreshTokenStore()

		if err := store.Save(ctx, &oauthserver.RefreshToken{ID: "refresh-1"}); err != nil {
			t.Fatal(err)
		}

		revoked, err := store.IsRevoked(ctx, "refresh-1")

		if err != nil {
			t.Fatal(err)
		}

		if revoked {
			t.Fatal("expected refresh token not to be revoked")
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_device_code_is_revoked", func(t *testing.T) {
		store := oauthserver.NewMemoryDeviceCodeStore()

		if err := store.Save(ctx, &oauthserver.DeviceCode{ID: "device-1"}); err != nil {
			t.Fatal(err)
		}

		if err := store.Revoke(ctx, "device-1"); err != nil {
			t.Fatal(err)
		}

		revoked, err := store.IsRevoked(ctx, "device-1")

		if err != nil {
			t.Fatal(err)
		}

		if !revoked {
			t.Fatal("expected device code to be revoked")
		}
	})

	t.Run("RevokedTest::test_a_device_code_is_also_revoked_if_it_cannot_be_found", func(t *testing.T) {
		revoked, err := oauthserver.NewMemoryDeviceCodeStore().IsRevoked(ctx, "missing")

		if err != nil {
			t.Fatal(err)
		}

		if !revoked {
			t.Fatal("expected missing device code to be treated as revoked")
		}
	})

	t.Run("RevokedTest::test_it_can_determine_if_a_device_code_is_not_revoked", func(t *testing.T) {
		store := oauthserver.NewMemoryDeviceCodeStore()

		if err := store.Save(ctx, &oauthserver.DeviceCode{ID: "device-1"}); err != nil {
			t.Fatal(err)
		}

		revoked, err := store.IsRevoked(ctx, "device-1")

		if err != nil {
			t.Fatal(err)
		}

		if revoked {
			t.Fatal("expected device code not to be revoked")
		}
	})
}

func TestOAuthServerInventoryGuardAndMiddlewareParity(t *testing.T) {
	ctx := context.Background()

	t.Run("TokenGuardTest::test_user_can_be_pulled_via_bearer_token", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("user-1")
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
		clients := oauthserver.NewMemoryClientStore()

		if err := tokens.Save(ctx, newTestToken("token-1", "user-1", "client-1", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"user-1": user})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		got, err := guard.User(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if got == nil || got.GetAuthIdentifier() != "user-1" {
			t.Fatalf("unexpected user: %#v", got)
		}
	})

	t.Run("TokenGuardTest::test_user_is_resolved_only_once", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("user-1")
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)

		if err := tokens.Save(ctx, newTestToken("token-1", "user-1", "client-1", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		provider := &inventoryCountingProvider{users: map[string]cauth.Authenticatable{"user-1": user}}
		guard := oauthserver.NewTokenGuard(p, tokens, oauthserver.NewMemoryClientStore(), provider)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		if _, err := guard.User(ctx); err != nil {
			t.Fatal(err)
		}

		if _, err := guard.User(ctx); err != nil {
			t.Fatal(err)
		}

		if provider.calls != 1 {
			t.Fatalf("provider calls = %d, want 1", provider.calls)
		}
	})

	t.Run("TokenGuardTest::test_null_is_returned_if_no_user_is_found", func(t *testing.T) {
		p := newOAuthServer()
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)

		if err := tokens.Save(ctx, newTestToken("token-1", "missing-user", "client-1", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		guard := newGuard(p, tokens, oauthserver.NewMemoryClientStore(), map[string]cauth.Authenticatable{})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		got, err := guard.User(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if got != nil {
			t.Fatalf("expected no user, got %#v", got)
		}
	})

	t.Run("TokenGuardTest::test_null_is_returned_for_client_credentials_token", func(t *testing.T) {
		p := newOAuthServer()
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)

		if err := tokens.Save(ctx, newTestToken("token-1", "", "client-1", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		guard := newGuard(p, tokens, oauthserver.NewMemoryClientStore(), map[string]cauth.Authenticatable{})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		got, err := guard.User(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if got != nil {
			t.Fatalf("expected no user for client credentials token, got %#v", got)
		}
	})

	t.Run("TokenGuardTest::test_user_is_resolved_when_user_id_matches_client_id", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("shared-id")
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)

		if err := tokens.Save(ctx, newTestToken("token-1", "shared-id", "shared-id", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		guard := newGuard(p, tokens, oauthserver.NewMemoryClientStore(), map[string]cauth.Authenticatable{"shared-id": user})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		got, err := guard.User(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if got == nil || got.GetAuthIdentifier() != "shared-id" {
			t.Fatalf("unexpected user: %#v", got)
		}
	})

	t.Run("TokenGuardTest::test_client_can_be_pulled_via_bearer_token", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("user-1")
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
		clients := oauthserver.NewMemoryClientStore()

		if err := clients.Create(ctx, &oauthserver.Client{ID: "client-1", Name: "Worker", Secret: "secret"}); err != nil {
			t.Fatal(err)
		}

		if err := tokens.Save(ctx, newTestToken("token-1", "user-1", "client-1", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"user-1": user})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		if _, err := guard.User(ctx); err != nil {
			t.Fatal(err)
		}

		client, err := guard.Client(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if client == nil || client.ID != "client-1" {
			t.Fatalf("unexpected client: %#v", client)
		}
	})

	t.Run("TokenGuardTest::test_no_user_is_returned_when_oauth_throws_exception", func(t *testing.T) {
		p := newOAuthServer()
		tokens := &inventoryErrorTokenStore{err: errors.New("token lookup failed")}
		guard := oauthserver.NewTokenGuard(p, tokens, oauthserver.NewMemoryClientStore(), &stubProvider{users: map[string]cauth.Authenticatable{}})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		got, err := guard.User(ctx)

		if err == nil || got != nil {
			t.Fatalf("expected token lookup error, got user=%#v err=%v", got, err)
		}
	})

	t.Run("TokenGuardTest::test_no_client_is_returned_when_oauth_throws_exception", func(t *testing.T) {
		p := newOAuthServer()
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)

		if err := tokens.Save(ctx, newTestToken("token-1", "user-1", "client-1", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		guard := oauthserver.NewTokenGuard(p, tokens, &inventoryErrorClientStore{err: errors.New("client lookup failed")}, &stubProvider{users: map[string]cauth.Authenticatable{"user-1": newStubUser("user-1")}})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		if _, err := guard.User(ctx); err != nil {
			t.Fatal(err)
		}

		client, err := guard.Client(ctx)

		if err == nil || client != nil {
			t.Fatalf("expected client lookup error, got client=%#v err=%v", client, err)
		}
	})

	t.Run("TokenGuardTest::test_client_is_resolved_only_once", func(t *testing.T) {
		p := newOAuthServer()
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)

		if err := tokens.Save(ctx, newTestToken("token-1", "user-1", "client-1", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		clients := &inventoryCountingClientStore{client: &oauthserver.Client{ID: "client-1", Name: "Worker", Secret: "secret"}}
		guard := oauthserver.NewTokenGuard(p, tokens, clients, &stubProvider{users: map[string]cauth.Authenticatable{"user-1": newStubUser("user-1")}})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		if _, err := guard.User(ctx); err != nil {
			t.Fatal(err)
		}

		if _, err := guard.Client(ctx); err != nil {
			t.Fatal(err)
		}

		if _, err := guard.Client(ctx); err != nil {
			t.Fatal(err)
		}

		if clients.calls != 1 {
			t.Fatalf("client store calls = %d, want 1", clients.calls)
		}
	})

	t.Run("TokenGuardTest::test_null_is_returned_if_no_client_is_found", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("user-1")
		tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)

		if err := tokens.Save(ctx, newTestToken("token-1", "user-1", "missing-client", []string{"read"})); err != nil {
			t.Fatal(err)
		}

		guard := newGuard(p, tokens, oauthserver.NewMemoryClientStore(), map[string]cauth.Authenticatable{"user-1": user})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)

		if _, err := guard.User(ctx); err != nil {
			t.Fatal(err)
		}

		client, err := guard.Client(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if client != nil {
			t.Fatalf("expected no client, got %#v", client)
		}
	})

	t.Run("CheckTokenTest::test_request_is_passed_along_if_token_is_valid", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard)(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	// CheckTokenTest::test_request_is_passed_along_if_token_is_transient
	t.Run("CheckTokenTest::test_request_is_passed_along_if_token_is_transient", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"*"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "orders:read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("CheckTokenTest::test_request_is_passed_along_if_token_and_scope_are_valid", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("CheckTokenTest::test_exception_is_thrown_if_token_does_not_have_required_scopes", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "write")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
		}
	})

	t.Run("CheckTokenTest::test_request_is_passed_along_if_scopes_are_present_on_token", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read", "write"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "read", "write")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("CheckTokenTest::test_exception_is_thrown_when_oauth_throws_exception", func(t *testing.T) {
		guard := oauthserver.NewTokenGuard(
			newOAuthServer(),
			&inventoryErrorTokenStore{err: errors.New("token lookup failed")},
			oauthserver.NewMemoryClientStore(),
			&stubProvider{users: map[string]cauth.Authenticatable{}},
		)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard)(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})

	t.Run("CheckTokenTest::test_exception_is_thrown_if_token_doesnt_have_scope", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "read", "write")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
		}
	})

	t.Run("CheckTokenForAnyScopeTest::test_request_is_passed_along_if_token_is_valid", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard)(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	// CheckTokenForAnyScopeTest::test_request_is_passed_along_if_token_is_transient
	t.Run("CheckTokenForAnyScopeTest::test_request_is_passed_along_if_token_is_transient", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"*"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "orders:write", "orders:read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("CheckTokenForAnyScopeTest::test_request_is_passed_along_if_token_has_any_required_scope", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "write", "read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("CheckTokenForAnyScopeTest::test_exception_is_thrown_when_oauth_throws_exception", func(t *testing.T) {
		guard := oauthserver.NewTokenGuard(
			newOAuthServer(),
			&inventoryErrorTokenStore{err: errors.New("token lookup failed")},
			oauthserver.NewMemoryClientStore(),
			&stubProvider{users: map[string]cauth.Authenticatable{}},
		)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token-1")
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
		}
	})

	t.Run("CheckTokenForAnyScopeTest::test_exception_is_thrown_if_token_does_not_have_required_scope", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "write", "delete")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
		}
	})

	t.Run("CheckTokenForAnyScopeTest::test_request_is_passed_along_if_scopes_are_present_on_token", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read", "write"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "write")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("CheckTokenForAnyScopeTest::test_exception_is_thrown_if_token_doesnt_have_scope", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "admin")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
		}
	})
}

func TestOAuthServerInventoryActingAsAndPersonalAccessParity(t *testing.T) {
	ctx := context.Background()

	t.Run("ActingAsTest::testActingAsWhenTheRouteIsProtectedByAuthMiddleware", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("user-1")
		p.ActingAs(user, oauthserver.NewAccessToken(newTestToken("token-1", "user-1", "client-1", []string{"read"}), p), []string{"read"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)
		got, err := guard.User(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if got == nil || got.GetAuthIdentifier() != "user-1" {
			t.Fatalf("unexpected acting user: %#v", got)
		}
	})

	t.Run("ActingAsTest::testActingAsWhenTheRouteIsProtectedByCheckScopesMiddleware", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("user-1")
		p.ActingAs(user, oauthserver.NewAccessToken(newTestToken("token-1", "user-1", "client-1", []string{"read"}), p), []string{"read"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("ActingAsTest::testActingAsWhenTheRouteIsProtectedByCheckForAnyScopeMiddleware", func(t *testing.T) {
		p := newOAuthServer()
		user := newStubUser("user-1")
		p.ActingAs(user, oauthserver.NewAccessToken(newTestToken("token-1", "user-1", "client-1", []string{"read"}), p), []string{"read"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "write", "read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("ActingAsTest::testActingAsWhenTheRouteIsProtectedByCheckScopesMiddlewareWithInheritance", func(t *testing.T) {
		p := newOAuthServer().UseInheritedScopes(true)
		user := newStubUser("user-1")
		token := newTestToken("token-1", "user-1", "client-1", []string{"orders"}).WithOAuthServer(p)
		p.ActingAs(user, oauthserver.NewAccessToken(token, p), []string{"orders"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "orders:read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("ActingAsTest::testActingAsWhenTheRouteIsProtectedByCheckForAnyScopeMiddlewareWithInheritance", func(t *testing.T) {
		p := newOAuthServer().UseInheritedScopes(true)
		user := newStubUser("user-1")
		token := newTestToken("token-1", "user-1", "client-1", []string{"orders"}).WithOAuthServer(p)
		p.ActingAs(user, oauthserver.NewAccessToken(token, p), []string{"orders"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "billing:read", "orders:read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("ActingAsClientTest::testActingAsClientSetsTheClientOnTheGuard", func(t *testing.T) {
		p := newOAuthServer()
		p.ActingAsClient(&oauthserver.Client{ID: "client-1", Scopes: []string{"read"}}, []string{"read"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)

		client, err := guard.Client(ctx)

		if err != nil {
			t.Fatal(err)
		}

		if client == nil || client.ID != "client-1" {
			t.Fatalf("unexpected acting client: %#v", client)
		}
	})

	// ActingAsClientTest::testActingAsClientWhenTheRouteIsProtectedByCheckTokenMiddleware
	t.Run("ActingAsClientTest::testActingAsClientWhenTheRouteIsProtectedByCheckTokenMiddleware", func(t *testing.T) {
		p := newOAuthServer()
		p.ActingAsClient(&oauthserver.Client{ID: "client-1"}, []string{"orders:read"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckClientCredentials(guard, "orders:read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	// ActingAsClientTest::testActingAsClientWhenTheRouteIsProtectedByCheckTokenForAnyScope
	t.Run("ActingAsClientTest::testActingAsClientWhenTheRouteIsProtectedByCheckTokenForAnyScope", func(t *testing.T) {
		p := newOAuthServer()
		p.ActingAsClient(&oauthserver.Client{ID: "client-1"}, []string{"orders:read"})

		defer p.ClearActing()

		guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckClientCredentialsForAnyScope(guard, "orders:write", "orders:read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	t.Run("PersonalAccessGrantTest::testIssueToken", func(t *testing.T) {
		p := newOAuthServer().PersonalAccessTokensExpireIn(time.Hour)
		store := oauthserver.NewMemoryTokenStore()
		result, err := oauthserver.NewPersonalAccessTokenFactory(p, nil, store).
			Create(ctx, newStubUser("user-1"), "CI", []string{"read"})

		if err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, result.AccessToken)

		if err != nil {
			t.Fatal(err)
		}

		if result.TokenType != "Bearer" || result.ExpiresIn != 3600 || found == nil || found.UserID != "user-1" {
			t.Fatalf("unexpected personal access result=%#v token=%#v", result, found)
		}
	})

	t.Run("PersonalAccessGrantTest::testIssueTokenWithAllScopes", func(t *testing.T) {
		p := newOAuthServer()
		store := oauthserver.NewMemoryTokenStore()
		result, err := oauthserver.NewPersonalAccessTokenFactory(p, nil, store).
			Create(ctx, newStubUser("user-1"), "CI", []string{"*"})

		if err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, result.AccessToken)

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || !found.Can("any:scope") {
			t.Fatalf("expected wildcard personal access token, got %#v", found)
		}
	})

	t.Run("PersonalAccessGrantTest::testPersonalAccessTokenRequestIsDisabled", func(t *testing.T) {
		p := newOAuthServer()

		if p.IsGrantEnabled(oauthserver.GrantPersonalAccess) {
			t.Fatal("personal_access is not a first-class grant toggle and should not be enabled")
		}
	})

	t.Run("PersonalAccessGrantWithTrustedHostsTest::testIssueToken", func(t *testing.T) {
		p := oauthserver.NewOAuthServer(&oauthserver.OAuthServerConfig{
			PersonalAccessClientID:     "client-1",
			PersonalAccessClientSecret: "secret",
		}).PersonalAccessTokensExpireIn(time.Hour)
		server := &inventoryAuthServer{issued: oauthserver.IssuedToken{
			TokenType:   "Bearer",
			AccessToken: "issued-token",
			TokenID:     "token-1",
			ExpiresIn:   3600,
		}}
		dispatcher := &inventoryDispatcher{}
		result, err := oauthserver.NewPersonalAccessTokenFactory(p, server, oauthserver.NewMemoryTokenStore()).
			WithEventDispatcher(dispatcher).
			Create(ctx, newStubUser("user-1"), "CI", []string{"read"})

		if err != nil {
			t.Fatal(err)
		}

		if result.AccessToken != "issued-token" || server.req.ClientID != "client-1" || server.req.ClientSecret != "secret" {
			t.Fatalf("unexpected server request/result: req=%#v result=%#v", server.req, result)
		}

		if len(dispatcher.events) != 1 {
			t.Fatalf("expected one event, got %d", len(dispatcher.events))
		}
	})

	t.Run("PersonalAccessTokenControllerTest::test_tokens_can_be_retrieved_for_users", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()

		if err := store.Save(ctx, newTestToken("token-1", "user-1", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		if err := store.Save(ctx, newTestToken("token-2", "user-2", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		tokens, err := store.ForUser(ctx, "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if len(tokens) != 1 || tokens[0].ID != "token-1" {
			t.Fatalf("unexpected personal tokens: %#v", tokens)
		}
	})

	t.Run("PersonalAccessTokenControllerTest::test_tokens_can_be_updated", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()
		token := newTestToken("token-1", "user-1", "client-1", nil)
		token.Name = "Original"

		if err := store.Save(ctx, token); err != nil {
			t.Fatal(err)
		}

		token.Name = "Renamed"

		if err := store.Save(ctx, token); err != nil {
			t.Fatal(err)
		}

		found, err := store.FindForUser(ctx, "token-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || found.Name != "Renamed" {
			t.Fatalf("unexpected updated personal token: %#v", found)
		}
	})

	t.Run("PersonalAccessTokenControllerTest::test_tokens_can_be_deleted", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()

		if err := store.Save(ctx, newTestToken("token-1", "user-1", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		if err := store.Revoke(ctx, "token-1"); err != nil {
			t.Fatal(err)
		}

		token, err := store.FindForUser(ctx, "token-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if token == nil || !token.IsRevoked() {
			t.Fatalf("expected revoked personal token, got %#v", token)
		}
	})

	t.Run("PersonalAccessTokenControllerTest::test_not_found_response_is_returned_if_user_doesnt_have_token", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()

		if err := store.Save(ctx, newTestToken("token-1", "user-2", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		token, err := store.FindForUser(ctx, "token-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if token != nil {
			t.Fatalf("expected no personal token for different owner, got %#v", token)
		}
	})

	t.Run("AuthorizedAccessTokenControllerTest::test_tokens_can_be_retrieved_for_users", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()

		if err := store.Save(ctx, newTestToken("token-1", "user-1", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		if err := store.Save(ctx, newTestToken("token-2", "user-2", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		tokens, err := store.ForUser(ctx, "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if len(tokens) != 1 || tokens[0].ID != "token-1" {
			t.Fatalf("unexpected user tokens: %#v", tokens)
		}
	})

	t.Run("AuthorizedAccessTokenControllerTest::test_tokens_can_be_deleted", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()

		if err := store.Save(ctx, newTestToken("token-1", "user-1", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		if err := store.Revoke(ctx, "token-1"); err != nil {
			t.Fatal(err)
		}

		token, err := store.FindForUser(ctx, "token-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if token == nil || !token.IsRevoked() {
			t.Fatalf("expected revoked user token, got %#v", token)
		}
	})

	t.Run("AuthorizedAccessTokenControllerTest::test_not_found_response_is_returned_if_user_doesnt_have_token", func(t *testing.T) {
		store := oauthserver.NewMemoryTokenStore()

		if err := store.Save(ctx, newTestToken("token-1", "user-2", "client-1", nil)); err != nil {
			t.Fatal(err)
		}

		token, err := store.FindForUser(ctx, "token-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if token != nil {
			t.Fatalf("expected no token for different owner, got %#v", token)
		}
	})
}

func TestOAuthServerInventoryClientRepositoryParity(t *testing.T) {
	ctx := context.Background()

	t.Run("BridgeClientRepositoryTest::test_can_get_client", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()
		client, err := store.CreateAuthCodeClient(ctx, "user-1", "Web App", "https://example.com/callback", "users")

		if err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, client.ID)

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || found.ID != client.ID || found.Provider != "users" || len(found.RedirectURIs) != 1 || found.RedirectURIs[0] != "https://example.com/callback" {
			t.Fatalf("unexpected client: %#v", found)
		}
	})

	t.Run("BridgeClientRepositoryTest::test_can_validate_client", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()
		client, err := store.CreatePersonalAccessClient(ctx, "user-1", "Personal", "users")

		if err != nil {
			t.Fatal(err)
		}

		if !client.FirstParty() || client.Secret != "" {
			t.Fatalf("unexpected personal access client: %#v", client)
		}

		revoked := &oauthserver.Client{ID: "client-2", Name: "Revoked", Secret: "secret", Revoked: true}

		if err := store.Create(ctx, revoked); err != nil {
			t.Fatal(err)
		}

		found, err := store.FindActive(ctx, revoked.ID)

		if err != nil {
			t.Fatal(err)
		}

		if found != nil {
			t.Fatalf("expected revoked client to be inactive, got %#v", found)
		}
	})

	t.Run("ClientControllerTest::test_all_the_clients_for_the_current_user_can_be_retrieved", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()

		if err := store.Create(ctx, &oauthserver.Client{ID: "client-1", UserID: "user-1", Name: "CLI"}); err != nil {
			t.Fatal(err)
		}

		if err := store.Create(ctx, &oauthserver.Client{ID: "client-2", UserID: "user-2", Name: "Worker"}); err != nil {
			t.Fatal(err)
		}

		clients, err := store.ForUser(ctx, "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if len(clients) != 1 || clients[0].ID != "client-1" {
			t.Fatalf("unexpected user clients: %#v", clients)
		}
	})

	t.Run("ClientControllerTest::test_clients_can_be_stored", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()
		client, err := store.CreateAuthCodeClient(ctx, "user-1", "Web", "https://example.com/callback", "users")

		if err != nil {
			t.Fatal(err)
		}

		found, err := store.FindForUser(ctx, client.ID, "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || found.Name != "Web" || !found.Confidential() {
			t.Fatalf("unexpected stored client: %#v", found)
		}
	})

	t.Run("ClientControllerTest::test_public_clients_can_be_stored", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()
		client := &oauthserver.Client{ID: "client-1", UserID: "user-1", Name: "Mobile", RedirectURIs: []string{"bedrock://callback"}}

		if err := store.Create(ctx, client); err != nil {
			t.Fatal(err)
		}

		found, err := store.FindForUser(ctx, "client-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || found.Confidential() {
			t.Fatalf("expected public client without secret, got %#v", found)
		}
	})

	t.Run("ClientControllerTest::test_clients_can_be_updated", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()
		client := &oauthserver.Client{ID: "client-1", UserID: "user-1", Name: "Original", Secret: "secret"}

		if err := store.Create(ctx, client); err != nil {
			t.Fatal(err)
		}

		client.Name = "Renamed"
		client.RedirectURIs = []string{"https://example.com/new"}

		if err := store.Update(ctx, client); err != nil {
			t.Fatal(err)
		}

		found, err := store.FindForUser(ctx, "client-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if found == nil || found.Name != "Renamed" || len(found.RedirectURIs) != 1 || found.RedirectURIs[0] != "https://example.com/new" {
			t.Fatalf("unexpected updated client: %#v", found)
		}
	})

	t.Run("ClientControllerTest::test_404_response_if_client_doesnt_belong_to_user", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()

		if err := store.Create(ctx, &oauthserver.Client{ID: "client-1", UserID: "user-2", Name: "Other"}); err != nil {
			t.Fatal(err)
		}

		found, err := store.FindForUser(ctx, "client-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if found != nil {
			t.Fatalf("expected no client for different owner, got %#v", found)
		}
	})

	t.Run("ClientControllerTest::test_clients_can_be_deleted", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()

		if err := store.Create(ctx, &oauthserver.Client{ID: "client-1", UserID: "user-1", Name: "CLI"}); err != nil {
			t.Fatal(err)
		}

		if err := store.Delete(ctx, "client-1"); err != nil {
			t.Fatal(err)
		}

		found, err := store.FindActive(ctx, "client-1")

		if err != nil {
			t.Fatal(err)
		}

		if found != nil {
			t.Fatalf("expected deleted client to be inactive, got %#v", found)
		}
	})

	t.Run("ClientControllerTest::test_404_response_if_client_doesnt_belong_to_user_on_delete", func(t *testing.T) {
		store := oauthserver.NewMemoryClientStore()

		if err := store.Create(ctx, &oauthserver.Client{ID: "client-1", UserID: "user-2", Name: "Other"}); err != nil {
			t.Fatal(err)
		}

		owned, err := store.FindForUser(ctx, "client-1", "user-1")

		if err != nil {
			t.Fatal(err)
		}

		if owned != nil {
			t.Fatalf("expected no client for different owner, got %#v", owned)
		}

		other, err := store.FindForUser(ctx, "client-1", "user-2")

		if err != nil {
			t.Fatal(err)
		}

		if other == nil || other.Revoked {
			t.Fatalf("expected other user's client to remain active, got %#v", other)
		}
	})
}

func TestOAuthServerInventoryAdditionalRepositoryParity(t *testing.T) {
	ctx := context.Background()

	// ClientCredentialsGrantTest::testPublicClient
	t.Run("ClientCredentialsGrantTest::testPublicClient", func(t *testing.T) {
		client := &oauthserver.Client{ID: "client-1", Name: "Public API"}

		if client.Confidential() {
			t.Fatal("expected a client without a secret to be public")
		}
	})

	// PasswordGrantTest::testPublicClient
	t.Run("PasswordGrantTest::testPublicClient", func(t *testing.T) {
		client := &oauthserver.Client{ID: "client-1", Name: "Password App"}

		if client.Confidential() {
			t.Fatal("expected a client without a secret to be public")
		}
	})

	// Attributes/AuthorizeTokenTest::testItAppliesCheckTokenMiddleware
	t.Run("Attributes/AuthorizeTokenTest::testItAppliesCheckTokenMiddleware", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckToken(guard, "read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	// Attributes/AuthorizeTokenTest::testItAppliesCheckTokenForAnyScopeMiddleware
	t.Run("Attributes/AuthorizeTokenTest::testItAppliesCheckTokenForAnyScopeMiddleware", func(t *testing.T) {
		guard, tokenID := setupMiddlewareGuard(newOAuthServer(), newTestToken("token-1", "user-1", "client-1", []string{"read"}), newStubUser("user-1"))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenID)
		guard.SetRequest(req)
		recorder := httptest.NewRecorder()

		oauthserver.CheckTokenForAnyScope(guard, "write", "read")(okHandler).ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	})

	// DeviceAuthorizationGrantTest::testDeviceCodeSerialization
	t.Run("DeviceAuthorizationGrantTest::testDeviceCodeSerialization", func(t *testing.T) {
		expiresAt := time.Unix(1700000000, 0).UTC()
		polledAt := time.Unix(1699999900, 0).UTC()
		code := oauthserver.DeviceCode{
			ID:           "device-1",
			ClientID:     "client-1",
			UserID:       "user-1",
			DeviceCode:   "device-code",
			UserCode:     "ABCD-EFGH",
			Scopes:       []string{"read", "write"},
			Approved:     true,
			Revoked:      false,
			LastPolledAt: &polledAt,
			ExpiresAt:    expiresAt,
		}

		payload, err := json.Marshal(code)

		if err != nil {
			t.Fatal(err)
		}

		var got oauthserver.DeviceCode

		if err := json.Unmarshal(payload, &got); err != nil {
			t.Fatal(err)
		}

		if got.ID != code.ID || got.ClientID != code.ClientID || got.UserCode != code.UserCode || !got.Approved || got.Revoked {
			t.Fatalf("unexpected device code JSON round-trip: %#v", got)
		}

		if got.LastPolledAt == nil || !got.LastPolledAt.Equal(polledAt) || !got.ExpiresAt.Equal(expiresAt) {
			t.Fatalf("unexpected device code timestamps: %#v", got)
		}
	})

	// AuthorizationCodeGrantTest::testAuthRequestSerialization
	t.Run("AuthorizationCodeGrantTest::testAuthRequestSerialization", func(t *testing.T) {
		req := oauthserver.AuthorizationRequest{
			ClientID:            "client-1",
			RedirectURI:         "https://example.com/callback",
			State:               "state-1",
			Scopes:              []string{"read", "write"},
			ResponseType:        "code",
			CodeChallenge:       "challenge",
			CodeChallengeMethod: "S256",
		}

		payload, err := json.Marshal(req)

		if err != nil {
			t.Fatal(err)
		}

		var got oauthserver.AuthorizationRequest

		if err := json.Unmarshal(payload, &got); err != nil {
			t.Fatal(err)
		}

		if got.ClientID != req.ClientID || got.RedirectURI != req.RedirectURI || got.State != req.State || got.ResponseType != req.ResponseType {
			t.Fatalf("unexpected authorization request round-trip: %#v", got)
		}
	})

	// AuthorizationCodeGrantTest::testValidateAuthorizationRequest
	t.Run("AuthorizationCodeGrantTest::testValidateAuthorizationRequest", func(t *testing.T) {
		server := &inventoryAuthServer{}
		req := oauthserver.AuthorizationRequest{
			ClientID:     "client-1",
			RedirectURI:  "https://example.com/callback",
			State:        "state-1",
			Scopes:       []string{"read"},
			ResponseType: "code",
		}

		got, err := server.ValidateAuthorizationRequest(ctx, req)

		if err == nil {
			t.Fatal("expected validation stub to report not implemented")
		}

		if got != nil {
			t.Fatalf("expected nil validated request, got %#v", got)
		}
	})

	// AuthorizationCodeGrantTest::testValidateScopes
	t.Run("AuthorizationCodeGrantTest::testValidateScopes", func(t *testing.T) {
		p := newOAuthServer().TokensCan(map[string]string{"read": "Read", "write": "Write"})
		repo := oauthserver.NewScopeRepository(p, nil)

		scopes, err := repo.FinalizeScopes(ctx, []string{"read", "missing"}, oauthserver.GrantAuthorizationCode, "")

		if err != nil {
			t.Fatal(err)
		}

		if len(scopes) != 1 || scopes[0].ID != "read" {
			t.Fatalf("unexpected validated scopes: %#v", scopes)
		}
	})

	// TransientTokenControllerTest::test_token_can_be_refreshed
	t.Run("TransientTokenControllerTest::test_token_can_be_refreshed", func(t *testing.T) {
		store := oauthserver.NewMemoryRefreshTokenStore()
		refresh := &oauthserver.RefreshToken{
			ID:            "refresh-1",
			AccessTokenID: "access-1",
			ExpiresAt:     time.Now().Add(time.Hour),
		}

		if err := store.Save(ctx, refresh); err != nil {
			t.Fatal(err)
		}

		if err := store.RevokeByAccessTokenID(ctx, "access-1"); err != nil {
			t.Fatal(err)
		}

		revoked, err := store.IsRevoked(ctx, "refresh-1")

		if err != nil {
			t.Fatal(err)
		}

		if !revoked {
			t.Fatal("expected refresh token to be revoked after access token refresh")
		}
	})

	// AccessTokenControllerTest::test_a_token_can_be_issued
	t.Run("AccessTokenControllerTest::test_a_token_can_be_issued", func(t *testing.T) {
		p := newOAuthServer().PersonalAccessTokensExpireIn(time.Hour)
		store := oauthserver.NewMemoryTokenStore()

		result, err := oauthserver.NewPersonalAccessTokenFactory(p, nil, store).
			Create(ctx, newStubUser("user-1"), "CLI", []string{"read"})

		if err != nil {
			t.Fatal(err)
		}

		found, err := store.Find(ctx, result.AccessToken)

		if err != nil {
			t.Fatal(err)
		}

		if result.TokenType != "Bearer" || result.ExpiresIn != 3600 || found == nil || found.UserID != "user-1" || !found.Can("read") {
			t.Fatalf("unexpected access token result=%#v token=%#v", result, found)
		}
	})

	// AccessTokenControllerTest::test_exceptions_are_handled
	t.Run("AccessTokenControllerTest::test_exceptions_are_handled", func(t *testing.T) {
		p := newOAuthServer()
		store := oauthserver.NewMemoryTokenStore()
		server := &inventoryAuthServer{err: errors.New("issue failed")}

		result, err := oauthserver.NewPersonalAccessTokenFactory(p, server, store).
			Create(ctx, newStubUser("user-1"), "CLI", []string{"read"})

		if err == nil || result != nil {
			t.Fatalf("expected issue error, got result=%#v err=%v", result, err)
		}
	})
}
