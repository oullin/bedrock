package inception

import (
	"context"
	"net/http"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/contracts/events"
)

// PasswordBroker sends password reset links and validates tokens.
type PasswordBroker interface {
	SendResetLink(ctx context.Context, credentials map[string]string) error
	Reset(ctx context.Context, credentials map[string]string, callback func(user cauth.Authenticatable, password string) error) error
}

// EmailVerifier handles email verification dispatch and fulfillment.
type EmailVerifier interface {
	SendVerificationNotification(ctx context.Context, user cauth.Authenticatable) error
	Verify(ctx context.Context, id string, hash string) error
}

// RateLimiter throttles attempts by key.
type RateLimiter interface {
	TooManyAttempts(key string, maxAttempts int) bool
	Hit(key string, decay time.Duration) int
	Clear(key string)
	AvailableIn(key string) time.Duration
}

// Responder customizes HTTP responses for Inception endpoints.
type Responder interface {
	LoginResponse(w http.ResponseWriter, r *http.Request)
	LogoutResponse(w http.ResponseWriter, r *http.Request)
	RegisterResponse(w http.ResponseWriter, r *http.Request)
	PasswordResetLinkSentResponse(w http.ResponseWriter, r *http.Request)
	PasswordResetResponse(w http.ResponseWriter, r *http.Request)
	PasswordUpdateResponse(w http.ResponseWriter, r *http.Request)
	PasswordConfirmResponse(w http.ResponseWriter, r *http.Request)
	ProfileInformationUpdatedResponse(w http.ResponseWriter, r *http.Request)
	EmailVerificationSentResponse(w http.ResponseWriter, r *http.Request)
	TwoFactorChallengeResponse(w http.ResponseWriter, r *http.Request)
	TwoFactorEnabledResponse(w http.ResponseWriter, r *http.Request)
	TwoFactorDisabledResponse(w http.ResponseWriter, r *http.Request)
}

// HasTeams is implemented by the user model to expose team membership.
type HasTeams interface {
	cauth.Authenticatable
	Teams() []Team
	CurrentTeam() *Team
	SetCurrentTeam(team *Team)
	OwnsTeam(team *Team) bool
	BelongsToTeam(team *Team) bool
	TeamRole(team *Team) string
	HasTeamPermission(team *Team, permission string) bool
}

// TeamRepository persists teams and memberships.
type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	FindByID(ctx context.Context, id string) (*Team, error)
	Update(ctx context.Context, team *Team) error
	Delete(ctx context.Context, id string) error
	FindByOwner(ctx context.Context, ownerID string) ([]Team, error)
	Members(ctx context.Context, teamID string) ([]Membership, error)
	AddMember(ctx context.Context, membership *Membership) error
	UpdateMemberRole(ctx context.Context, teamID string, userID string, role string) error
	RemoveMember(ctx context.Context, teamID string, userID string) error
}

// InvitationRepository persists team invitations.
type InvitationRepository interface {
	Create(ctx context.Context, invitation *TeamInvitation) error
	FindByID(ctx context.Context, id string) (*TeamInvitation, error)
	FindByTeam(ctx context.Context, teamID string) ([]TeamInvitation, error)
	FindByEmail(ctx context.Context, teamID string, email string) (*TeamInvitation, error)
	Delete(ctx context.Context, id string) error
}

// EventDispatcher dispatches domain events.
type EventDispatcher = events.Dispatcher

// Event represents a domain event fired by Inception.
type Event struct {
	Name    string
	Payload any
}
