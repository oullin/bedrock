package events

import (
	"time"

	auth "github.com/gollin/packages/auth"
)

// Attempting is emitted before credentials are checked.
type Attempting struct {
	Guard       string
	Credentials map[string]string
}

// Login is emitted when a user logs in.
type Login struct {
	Guard    string
	User     auth.Authenticatable
	Remember bool
	At       time.Time
}

// Logout is emitted when a user logs out.
type Logout struct {
	Guard string
	User  auth.Authenticatable
	At    time.Time
}

// Failed is emitted when authentication fails.
type Failed struct {
	Guard       string
	User        auth.Authenticatable
	Credentials map[string]string
}

// Verified is emitted when a user's email is verified.
type Verified struct {
	User auth.Authenticatable
	At   time.Time
}
