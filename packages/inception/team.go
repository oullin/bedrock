package inception

import "time"

// Team represents a team that users can belong to.
type Team struct {
	ID           string
	OwnerID      string
	Name         string
	PersonalTeam bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Membership represents a user's membership in a team.
type Membership struct {
	TeamID string
	UserID string
	Role   string
}

// TeamInvitation represents a pending invitation to join a team.
type TeamInvitation struct {
	ID        string
	TeamID    string
	Email     string
	Role      string
	CreatedAt time.Time
}
