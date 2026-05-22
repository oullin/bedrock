package models

import "time"

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
