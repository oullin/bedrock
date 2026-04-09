package jetstream

// Event names dispatched by Jetstream.
const (
	EventTeamCreated       = "jetstream.team.created"
	EventTeamUpdated       = "jetstream.team.updated"
	EventTeamDeleted       = "jetstream.team.deleted"
	EventTeamMemberAdded   = "jetstream.team.member_added"
	EventTeamMemberUpdated = "jetstream.team.member_updated"
	EventTeamMemberRemoved = "jetstream.team.member_removed"
	EventTeamMemberInvited = "jetstream.team.member_invited"
	EventTeamSwitched      = "jetstream.team.switched"
	EventTokenCreated      = "jetstream.token.created"
	EventTokenUpdated      = "jetstream.token.updated"
	EventTokenDeleted      = "jetstream.token.deleted"
	EventUserDeleted       = "jetstream.user.deleted"
)
