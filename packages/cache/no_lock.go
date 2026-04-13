package cache

import (
	"context"
	"time"
)

// NoLock is a null-object Lock that always succeeds. Useful for stores that
// do not support locking or for testing.
type NoLock struct{}

var _ Lock = (*NoLock)(nil)

func (l *NoLock) Acquire(_ context.Context) (bool, error)                 { return true, nil }
func (l *NoLock) Release(_ context.Context) (bool, error)                 { return true, nil }
func (l *NoLock) ForceRelease(_ context.Context) error                    { return nil }
func (l *NoLock) Get(ctx context.Context, fn func() error) error          { return fn() }
func (l *NoLock) Block(_ context.Context, _ time.Duration) error          { return nil }
func (l *NoLock) Blocked(_ context.Context) (bool, error)                 { return false, nil }
func (l *NoLock) Owner() string                                           { return "" }
func (l *NoLock) IsOwnedByCurrentProcess(_ context.Context) (bool, error) { return true, nil }
func (l *NoLock) IsOwnedBy(_ context.Context, _ string) (bool, error)     { return true, nil }
func (l *NoLock) BetweenBlockedAttemptsSleepFor(_ int) Lock               { return l }
