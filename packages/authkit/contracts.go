package authkit

import (
	"context"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/contracts/events"
)

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

// Event represents a domain event fired by AuthKit.
type Event struct {
	Name    string
	Payload any
}
