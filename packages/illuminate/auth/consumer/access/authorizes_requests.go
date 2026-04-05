package access

import (
	"context"
	"strings"

	auth "github.com/gollin/packages/framework/auth"
	authaccess "github.com/gollin/packages/framework/auth/access"
)

// AuthorizesRequests adapts controller-style authorization helpers.
type AuthorizesRequests struct {
	Gate authaccess.Authorizer
}

// Authorize authorizes an ability for the current user.
func (a AuthorizesRequests) Authorize(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) error {
	return a.Gate.Authorize(ctx, user, ability, arguments...)
}

// ParseAbilityAndArguments normalizes an ability and its arguments.
func (AuthorizesRequests) ParseAbilityAndArguments(ability string, arguments ...any) (string, []any) {
	ability = strings.TrimSpace(ability)

	if ability == "" && len(arguments) > 0 {
		if guessed, ok := arguments[0].(string); ok {
			return normalizeGuessedAbilityName(guessed), arguments[1:]
		}
	}

	return normalizeGuessedAbilityName(ability), arguments
}

// ResourceAbilityMap returns the default resource ability map.
func (AuthorizesRequests) ResourceAbilityMap() map[string]string {
	return map[string]string{
		"index":   "viewAny",
		"show":    "view",
		"create":  "create",
		"store":   "create",
		"edit":    "update",
		"update":  "update",
		"destroy": "delete",
	}
}

// ResourceMethodsWithoutModels returns resource actions without models.
func (AuthorizesRequests) ResourceMethodsWithoutModels() []string {
	return []string{"index", "create", "store"}
}

func normalizeGuessedAbilityName(value string) string {
	switch strings.TrimSpace(value) {
	case "show":
		return "view"
	case "destroy":
		return "delete"
	default:
		return strings.TrimSpace(value)
	}
}
