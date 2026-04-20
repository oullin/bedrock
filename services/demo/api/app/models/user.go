package models

import "time"

// User is the Go equivalent of laravel/laravel's default App\Models\User.
type User struct {
	ID              int64
	Name            string
	Email           string
	EmailVerifiedAt *time.Time
	Password        string
	RememberToken   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
