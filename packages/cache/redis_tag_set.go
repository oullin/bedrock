package cache

import (
	"context"
)

// RedisTagClient extends RedisClient with set operations required for
// Redis-specific tag tracking.
type RedisTagClient interface {
	RedisClient
	SAdd(ctx context.Context, key string, members ...string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SRem(ctx context.Context, key string, members ...string) error
}

// RedisTagSet extends TagSet with Redis set-based key tracking. Unlike the
// standard TagSet which orphans keys on flush, RedisTagSet tracks individual
// keys in Redis sets so they can be explicitly deleted.
type RedisTagSet struct {
	*TagSet
	client RedisTagClient
}

// NewRedisTagSet creates a RedisTagSet.
func NewRedisTagSet(store Store, names []string, client RedisTagClient) *RedisTagSet {
	return &RedisTagSet{
		TagSet: NewTagSet(store, names),
		client: client,
	}
}

// AddEntry tracks a cache key in the tag's Redis set.
func (rts *RedisTagSet) AddEntry(ctx context.Context, key string) error {
	for _, name := range rts.names {
		setKey := rts.entrySetKey(name)

		if err := rts.client.SAdd(ctx, setKey, key); err != nil {
			return err
		}
	}

	return nil
}

// FlushTaggedEntries enumerates all tracked keys and deletes them.
func (rts *RedisTagSet) FlushTaggedEntries(ctx context.Context) error {
	for _, name := range rts.names {
		setKey := rts.entrySetKey(name)

		members, err := rts.client.SMembers(ctx, setKey)

		if err != nil {
			return err
		}

		if len(members) > 0 {
			if err := rts.client.Del(ctx, members...); err != nil {
				return err
			}
		}

		if err := rts.client.Del(ctx, setKey); err != nil {
			return err
		}
	}

	return rts.Reset(ctx)
}

func (rts *RedisTagSet) entrySetKey(name string) string {
	return "tag:" + name + ":entries"
}
