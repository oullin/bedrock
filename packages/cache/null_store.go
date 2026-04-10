package cache

import (
	"context"
	"fmt"
	"time"
)

var _ Store = (*NullStore)(nil)
var _ Locker = (*NullStore)(nil)

// NullStore discards all writes and returns ErrNotFound for all reads.
// Useful for disabling caching or as a test double.
type NullStore struct{}

// NewNullStore creates a NullStore.
func NewNullStore() *NullStore { return &NullStore{} }

func (s *NullStore) GetPrefix() string { return "" }

func (s *NullStore) Get(_ context.Context, key string) (any, error) {
	return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
}

func (s *NullStore) GetMany(_ context.Context, _ []string) (map[string]any, error) {
	return make(map[string]any), nil
}

func (s *NullStore) Put(_ context.Context, _ string, _ any, _ time.Duration) error { return nil }

func (s *NullStore) PutMany(_ context.Context, _ map[string]any, _ time.Duration) error {
	return nil
}

func (s *NullStore) Add(_ context.Context, _ string, _ any, _ time.Duration) (bool, error) {
	return true, nil
}

func (s *NullStore) Forever(_ context.Context, _ string, _ any) error { return nil }

func (s *NullStore) Increment(_ context.Context, _ string, _ int64) (int64, error) { return 0, nil }

func (s *NullStore) Decrement(_ context.Context, _ string, _ int64) (int64, error) { return 0, nil }

func (s *NullStore) Touch(_ context.Context, _ string, _ time.Duration) (bool, error) {
	return false, nil
}

func (s *NullStore) Forget(_ context.Context, _ string) error { return nil }

func (s *NullStore) Flush(_ context.Context) error { return nil }

// Lock returns a no-op lock.
func (s *NullStore) Lock(_, _ string, _ time.Duration) Lock { return &noLock{} }

// noLock is a no-op Lock that always appears to succeed.
type noLock struct{}

func (l *noLock) Acquire(_ context.Context) (bool, error)         { return true, nil }
func (l *noLock) Release(_ context.Context) (bool, error)         { return true, nil }
func (l *noLock) ForceRelease(_ context.Context) error             { return nil }
func (l *noLock) Get(ctx context.Context, fn func() error) error   { return fn() }
func (l *noLock) Block(_ context.Context, _ time.Duration) error   { return nil }
func (l *noLock) Blocked(_ context.Context) (bool, error)          { return false, nil }
