package authkit

import (
	"context"
	"time"
)

// BrowserSession represents an active user session.
type BrowserSession struct {
	ID         string
	IPAddress  string
	UserAgent  string
	LastActive time.Time
	IsCurrent  bool
}

// SessionRepository lists and purges browser sessions.
type SessionRepository interface {
	FindByUser(ctx context.Context, userID string) ([]BrowserSession, error)
	DeleteOthers(ctx context.Context, userID string, currentSessionID string) error
}
