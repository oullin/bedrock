package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var _ Store = (*RedisStore)(nil)
var _ Locker = (*RedisStore)(nil)

// RedisClient is the subset of Redis operations required by RedisStore.
// Compatible with go-redis and similar clients.
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, keys ...string) ([]any, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
	Del(ctx context.Context, keys ...string) error
	FlushDB(ctx context.Context) error
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Eval(ctx context.Context, script string, keys []string, args ...any) (any, error)
}

// RedisStore stores cache values in Redis.
type RedisStore struct {
	client RedisClient
	prefix string
	clock  Clock
}

// NewRedisStore creates a RedisStore with the given client.
func NewRedisStore(client RedisClient, prefix string) *RedisStore {
	return &RedisStore{client: client, prefix: prefix}
}

func (s *RedisStore) GetPrefix() string { return s.prefix }

func (s *RedisStore) prefixed(key string) string {
	if s.prefix == "" {
		return key
	}

	return s.prefix + ":" + key
}

func (s *RedisStore) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.Get(ctx, s.prefixed(key))
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	var v any
	if err := json.Unmarshal([]byte(val), &v); err != nil {
		return val, nil
	}

	return v, nil
}

func (s *RedisStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	prefixed := make([]string, len(keys))
	for i, k := range keys {
		prefixed[i] = s.prefixed(k)
	}

	vals, err := s.client.MGet(ctx, prefixed...)
	if err != nil {
		return nil, err
	}

	out := make(map[string]any, len(keys))

	for i, v := range vals {
		if v == nil {
			continue
		}

		str, ok := v.(string)
		if !ok {
			out[keys[i]] = v
			continue
		}

		var decoded any
		if err := json.Unmarshal([]byte(str), &decoded); err != nil {
			out[keys[i]] = str
		} else {
			out[keys[i]] = decoded
		}
	}

	return out, nil
}

func (s *RedisStore) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, s.prefixed(key), string(encoded), ttl)
}

func (s *RedisStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	for k, v := range values {
		if err := s.Put(ctx, k, v, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (s *RedisStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return s.client.SetNX(ctx, s.prefixed(key), string(encoded), ttl)
}

func (s *RedisStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *RedisStore) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	return s.client.IncrBy(ctx, s.prefixed(key), delta)
}

func (s *RedisStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *RedisStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return s.client.Expire(ctx, s.prefixed(key), ttl)
}

func (s *RedisStore) Forget(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.prefixed(key))
}

func (s *RedisStore) Flush(ctx context.Context) error {
	return s.client.FlushDB(ctx)
}

// Lock returns a Redis-based distributed lock using a SET NX PX Lua script.
func (s *RedisStore) Lock(name, owner string, ttl time.Duration) Lock {
	return &redisLock{
		client: s.client,
		key:    s.prefixed("lock:" + name),
		owner:  owner,
		ttl:    ttl,
	}
}

// redisLock implements Lock using Redis SET NX.
type redisLock struct {
	client RedisClient
	key    string
	owner  string
	ttl    time.Duration
}

// acquireScript acquires the lock only if it's not held or has expired.
const redisAcquireScript = `
local current = redis.call("GET", KEYS[1])
if current == false or current == ARGV[1] then
  redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
  return 1
end
return 0`

// releaseScript releases the lock only if the owner matches.
const redisReleaseScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end`

func (l *redisLock) Acquire(ctx context.Context) (bool, error) {
	ms := int64(l.ttl / time.Millisecond)
	if ms <= 0 {
		ms = 0
	}

	result, err := l.client.Eval(ctx, redisAcquireScript, []string{l.key}, l.owner, ms)
	if err != nil {
		return false, err
	}

	return result == int64(1), nil
}

func (l *redisLock) Release(ctx context.Context) (bool, error) {
	result, err := l.client.Eval(ctx, redisReleaseScript, []string{l.key}, l.owner)
	if err != nil {
		return false, err
	}

	return result == int64(1), nil
}

func (l *redisLock) ForceRelease(ctx context.Context) error {
	return l.client.Del(ctx, l.key)
}

func (l *redisLock) Get(ctx context.Context, fn func() error) error {
	ok, err := l.Acquire(ctx)
	if err != nil {
		return err
	}

	if !ok {
		return ErrLockTimeout
	}

	defer l.Release(ctx) //nolint:errcheck

	return fn()
}

func (l *redisLock) Block(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		ok, err := l.Acquire(ctx)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		if time.Now().After(deadline) {
			return ErrLockTimeout
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (l *redisLock) Blocked(ctx context.Context) (bool, error) {
	val, err := l.client.Get(ctx, l.key)
	if err != nil {
		return false, nil
	}

	return val != "" && val != l.owner, nil
}
