package passport

import (
	"context"
	"encoding/json"
	"strings"
)

// Scope represents a named OAuth2 permission with a human-readable description.
// It maps directly to Passport Scope class.
type Scope struct {
	ID          string
	Description string
}

// ScopeRepository validates requested OAuth scopes against registered Passport
// scopes and optional client-level restrictions.
type ScopeRepository struct {
	passport *Passport
	clients  ClientStore
}

// NewScopeRepository creates a scope repository.
func NewScopeRepository(p *Passport, clients ClientStore) *ScopeRepository {
	return &ScopeRepository{passport: p, clients: clients}
}

// ToArray returns the scope as a plain map, matching the upstream toArray().
func (s Scope) ToArray() map[string]string {
	return map[string]string{
		"id":          s.ID,
		"description": s.Description,
	}
}

// MarshalJSON encodes the scope as a JSON object with id and description fields.
func (s Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToArray())
}

// FinalizeScopes filters invalid scopes and scopes disallowed for the client.
// The wildcard scope is only accepted for the personal access grant.
func (r *ScopeRepository) FinalizeScopes(ctx context.Context, requested []string, grantType, clientID string) ([]Scope, error) {
	if r == nil || r.passport == nil {
		return nil, nil
	}

	var client *Client

	if r.clients != nil && clientID != "" {
		var err error

		client, err = r.clients.Find(ctx, clientID)

		if err != nil {
			return nil, err
		}
	}

	out := make([]Scope, 0, len(requested))

	for _, id := range requested {
		scope := r.scopeForGrant(id, grantType)

		if scope == nil {
			continue
		}

		if client != nil && len(client.Scopes) > 0 && !scopeExistsIn(id, client.Scopes, r.passport.InheritedScopesEnabled()) {
			continue
		}

		out = append(out, *scope)
	}

	return out, nil
}

func (r *ScopeRepository) scopeForGrant(id, grantType string) *Scope {
	if id == "*" {
		if grantType != GrantPersonalAccess {
			return nil
		}

		return &Scope{ID: "*", Description: "All scopes"}
	}

	return r.passport.FindScope(id)
}

// resolveInheritedScopes expands a scope identifier into all ancestor scopes.
//
// For example, "admin:webhooks:read" resolves to:
//
//	["admin", "admin:webhooks", "admin:webhooks:read"]
//
// This mirrors the ResolvesInheritedScopes trait in Passport.
func resolveInheritedScopes(scope string) []string {
	parts := strings.Split(scope, ":")
	scopes := make([]string, 0, len(parts))

	for i := 1; i <= len(parts); i++ {
		scopes = append(scopes, strings.Join(parts[:i], ":"))
	}

	return scopes
}

// scopeExistsIn reports whether scope (or one of its ancestors when inherited
// scopes are enabled) is present in the haystack slice.
func scopeExistsIn(scope string, haystack []string, inherited bool) bool {
	candidates := []string{scope}

	if inherited {
		candidates = resolveInheritedScopes(scope)
	}

	for _, c := range candidates {
		for _, h := range haystack {
			if c == h {
				return true
			}
		}
	}

	return false
}
