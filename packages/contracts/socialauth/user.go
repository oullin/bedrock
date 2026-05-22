package socialauth

// User represents an authenticated OAuth user.
type User interface {
	GetID() string
	GetNickname() string
	GetName() string
	GetEmail() string
	GetAvatar() string
}
