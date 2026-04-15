package socialauth

// User represents an authenticated OAuth user.
// It mirrors Upstream\SocialAuth\Contracts\User.
type User interface {
	GetID() string
	GetNickname() string
	GetName() string
	GetEmail() string
	GetAvatar() string
}
