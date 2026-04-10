package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var _ Store = (*MemcachedStore)(nil)

// MemcachedClient is the subset of Memcached operations required by MemcachedStore.
type MemcachedClient interface {
	Get(key string) ([]byte, error)
	GetMulti(keys []string) (map[string][]byte, error)
	Set(key string, value []byte, expiration int32) error
	Add(key string, value []byte, expiration int32) error
	Delete(key string) error
	FlushAll() error
	Increment(key string, delta uint64) (uint64, error)
	Decrement(key string, delta uint64) (uint64, error)
}

// MemcachedStore stores cache values in Memcached.
type MemcachedStore struct {
	client MemcachedClient
	prefix string
}

// NewMemcachedStore creates a MemcachedStore.
func NewMemcachedStore(client MemcachedClient, prefix string) *MemcachedStore {
	return &MemcachedStore{client: client, prefix: prefix}
}

func (s *MemcachedStore) GetPrefix() string { return s.prefix }

func (s *MemcachedStore) prefixed(key string) string {
	if s.prefix == "" {
		return key
	}

	return s.prefix + ":" + key
}

func (s *MemcachedStore) ttlSeconds(ttl time.Duration) int32 {
	if ttl <= 0 {
		return 0
	}

	return int32(ttl.Seconds())
}

func (s *MemcachedStore) Get(_ context.Context, key string) (any, error) {
	data, err := s.client.Get(s.prefixed(key))
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return string(data), nil
	}

	return v, nil
}

func (s *MemcachedStore) GetMany(_ context.Context, keys []string) (map[string]any, error) {
	prefixed := make([]string, len(keys))
	for i, k := range keys {
		prefixed[i] = s.prefixed(k)
	}

	data, err := s.client.GetMulti(prefixed)
	if err != nil {
		return nil, err
	}

	out := make(map[string]any, len(keys))

	for i, pk := range prefixed {
		if v, ok := data[pk]; ok {
			var decoded any
			if err := json.Unmarshal(v, &decoded); err != nil {
				out[keys[i]] = string(v)
			} else {
				out[keys[i]] = decoded
			}
		}
	}

	return out, nil
}

func (s *MemcachedStore) Put(_ context.Context, key string, value any, ttl time.Duration) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(s.prefixed(key), encoded, s.ttlSeconds(ttl))
}

func (s *MemcachedStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	for k, v := range values {
		if err := s.Put(ctx, k, v, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (s *MemcachedStore) Add(_ context.Context, key string, value any, ttl time.Duration) (bool, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	err = s.client.Add(s.prefixed(key), encoded, s.ttlSeconds(ttl))
	if err != nil {
		// ADD fails if key exists in Memcached.
		return false, nil
	}

	return true, nil
}

func (s *MemcachedStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *MemcachedStore) Increment(_ context.Context, key string, delta int64) (int64, error) {
	if delta >= 0 {
		v, err := s.client.Increment(s.prefixed(key), uint64(delta))

		return int64(v), err
	}

	v, err := s.client.Decrement(s.prefixed(key), uint64(-delta))

	return int64(v), err
}

func (s *MemcachedStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *MemcachedStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	v, err := s.Get(ctx, key)
	if err != nil {
		return false, nil
	}

	return true, s.Put(ctx, key, v, ttl)
}

func (s *MemcachedStore) Forget(_ context.Context, key string) error {
	return s.client.Delete(s.prefixed(key))
}

func (s *MemcachedStore) Flush(_ context.Context) error {
	return s.client.FlushAll()
}
