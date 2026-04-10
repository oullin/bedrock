package consumer

import (
	"context"
	"net/http"

	"github.com/bedrock/packages/auth/access"
)

// AuthorizesRequests provides gate-based authorization helpers for HTTP handlers.
type AuthorizesRequests struct {
	Gate *access.Gate
}

// Authorize checks the given ability and writes a 403 if denied.
// Returns false when the request is denied (caller should return early).
func (a *AuthorizesRequests) Authorize(ctx context.Context, w http.ResponseWriter, ability string, model any) bool {
	if err := a.Gate.Authorize(ctx, ability, model); err != nil {
		if ae, ok := err.(*access.AuthorizationException); ok {
			http.Error(w, ae.Error(), ae.Response.StatusCode)
		} else {
			http.Error(w, "forbidden", http.StatusForbidden)
		}

		return false
	}

	return true
}

// Can reports whether the current user can perform the ability.
func (a *AuthorizesRequests) Can(ctx context.Context, ability string, model any) bool {
	return a.Gate.Check(ctx, ability, model)
}

// Cannot reports whether the current user cannot perform the ability.
func (a *AuthorizesRequests) Cannot(ctx context.Context, ability string, model any) bool {
	return !a.Can(ctx, ability, model)
}
