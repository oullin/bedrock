package authkit

// Event names dispatched by AuthKit.
const (
	EventTeamCreated       = "authkit.team.created"
	EventTeamUpdated       = "authkit.team.updated"
	EventTeamDeleted       = "authkit.team.deleted"
	EventTeamMemberAdded   = "authkit.team.member_added"
	EventTeamMemberUpdated = "authkit.team.member_updated"
	EventTeamMemberRemoved = "authkit.team.member_removed"
	EventTeamMemberInvited = "authkit.team.member_invited"
	EventTeamSwitched      = "authkit.team.switched"
	EventTokenCreated      = "authkit.token.created"
	EventTokenUpdated      = "authkit.token.updated"
	EventTokenDeleted      = "authkit.token.deleted"
	EventUserDeleted       = "authkit.user.deleted"
)
