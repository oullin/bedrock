package cache

import (
	"context"
	"fmt"
	"time"
)

// NullStore discards all writes and returns ErrNotFound for all reads.
// Useful for disabling caching or as a test double.
type NullStore struct{}

var _ Store = (*NullStore)(nil)
var _ Locker = (*NullStore)(nil)

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
func (s *NullStore) Lock(_, _ string, _ time.Duration) Lock { return &NoLock{} }
