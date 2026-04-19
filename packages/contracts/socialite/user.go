package socialite

// User represents an authenticated OAuth user.
// It mirrors Laravel\Socialite\Contracts\User.
type User interface {
	GetID() string
	GetNickname() string
	GetName() string
	GetEmail() string
	GetAvatar() string
}
