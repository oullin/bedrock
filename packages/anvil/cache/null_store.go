package cache

import (
	"context"
	"fmt"
	"time"
)

var _ Store = (*NullStore)(nil)

// NullStore is a cache store that discards all data. Every write silently
// succeeds and every read reports a miss. It is useful for disabling
// caching or as a test double.
type NullStore struct{}

// NewNullStore creates a NullStore.
func NewNullStore() *NullStore {
	return &NullStore{}
}

func (s *NullStore) Get(_ context.Context, key string) (any, error) {
	return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
}

func (s *NullStore) GetMultiple(_ context.Context, keys []string) (map[string]any, error) {
	return make(map[string]any), nil
}

func (s *NullStore) Put(_ context.Context, _ string, _ any, _ time.Duration) error {
	return nil
}

func (s *NullStore) PutMultiple(_ context.Context, _ map[string]any, _ time.Duration) error {
	return nil
}

func (s *NullStore) Add(_ context.Context, _ string, _ any, _ time.Duration) (bool, error) {
	return true, nil
}

func (s *NullStore) Forever(_ context.Context, _ string, _ any) error {
	return nil
}

func (s *NullStore) Remember(ctx context.Context, _ string, _ time.Duration, callback func() (any, error)) (any, error) {
	return callback()
}

func (s *NullStore) RememberForever(ctx context.Context, _ string, callback func() (any, error)) (any, error) {
	return callback()
}

func (s *NullStore) Touch(_ context.Context, _ string, _ time.Duration) (bool, error) {
	return false, nil
}

func (s *NullStore) Has(_ context.Context, _ string) bool {
	return false
}

func (s *NullStore) Missing(_ context.Context, _ string) bool {
	return true
}

func (s *NullStore) Increment(_ context.Context, _ string, _ int64) (int64, error) {
	return 0, nil
}

func (s *NullStore) Decrement(_ context.Context, _ string, _ int64) (int64, error) {
	return 0, nil
}

func (s *NullStore) Forget(_ context.Context, _ string) error {
	return nil
}

func (s *NullStore) Flush(_ context.Context) error {
	return nil
}
