package passport

import (
	"encoding/json"
	"strings"
)

// Scope represents a named OAuth2 permission with a human-readable description.
// It maps directly to Laravel Passport's Scope class.
type Scope struct {
	ID          string
	Description string
}

// ToArray returns the scope as a plain map, matching Laravel's toArray().
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

// resolveInheritedScopes expands a scope identifier into all ancestor scopes.
//
// For example, "admin:webhooks:read" resolves to:
//
//	["admin", "admin:webhooks", "admin:webhooks:read"]
//
// This mirrors the ResolvesInheritedScopes trait in Laravel Passport.
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
