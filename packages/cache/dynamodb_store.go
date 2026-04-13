package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bedrock/packages/contracts"
)

// DynamoClient is the DynamoDB operations required by DynamoDbStore.
type DynamoClient interface {
	GetItem(ctx context.Context, table, keyAttr, key string) (map[string]any, error)
	PutItem(ctx context.Context, table string, item map[string]any) error
	DeleteItem(ctx context.Context, table, keyAttr, key string) error
	Scan(ctx context.Context, table string) ([]map[string]any, error)
}

// DynamoDbStore stores cache values in AWS DynamoDB.
// The table must have a string partition key (keyAttr) and an optional
// TTL attribute (ttlAttr) for native expiration support.
type DynamoDbStore struct {
	client  DynamoClient
	table   string
	keyAttr string // partition key attribute name (default "key")
	valAttr string // value attribute name (default "value")
	ttlAttr string // TTL attribute name (default "expiration")
	prefix  string
	clock   contracts.Clock
}

var _ Store = (*DynamoDbStore)(nil)
var _ Locker = (*DynamoDbStore)(nil)

// NewDynamoDbStore creates a DynamoDbStore with sensible defaults.
func NewDynamoDbStore(client DynamoClient, table, prefix string) *DynamoDbStore {
	return &DynamoDbStore{
		client:  client,
		table:   table,
		keyAttr: "key",
		valAttr: "value",
		ttlAttr: "expiration",
		prefix:  prefix,
	}
}

func (s *DynamoDbStore) now() time.Time {
	if s.clock != nil {
		return s.clock.Now()
	}

	return time.Now()
}

func (s *DynamoDbStore) GetPrefix() string { return s.prefix }

// SetPrefix sets the key prefix.
func (s *DynamoDbStore) SetPrefix(prefix string) { s.prefix = prefix }

// GetClient returns the underlying DynamoDB client.
func (s *DynamoDbStore) GetClient() DynamoClient { return s.client }

// Lock returns a DynamoDB-backed lock for the named resource.
func (s *DynamoDbStore) Lock(name, owner string, ttl time.Duration) Lock {
	return NewDynamoDbLock(s.client, s.table, name, owner, ttl, s.clock)
}

// RestoreLock creates a lock handle from a serialized owner without acquiring.
func (s *DynamoDbStore) RestoreLock(name, owner string) Lock {
	return s.Lock(name, owner, 0)
}

func (s *DynamoDbStore) prefixed(key string) string {
	if s.prefix == "" {
		return key
	}

	return s.prefix + ":" + key
}

func (s *DynamoDbStore) Get(ctx context.Context, key string) (any, error) {
	item, err := s.client.GetItem(ctx, s.table, s.keyAttr, s.prefixed(key))

	if err != nil || item == nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	// Check TTL.
	if exp, ok := item[s.ttlAttr]; ok {
		if expInt, err := toInt64(exp); err == nil {
			if expInt > 0 && s.now().Unix() > expInt {
				return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
			}
		}
	}

	raw, ok := item[s.valAttr]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	str, ok := raw.(string)

	if !ok {
		return raw, nil
	}

	var v any

	if err := json.Unmarshal([]byte(str), &v); err != nil {
		return str, nil
	}

	return v, nil
}

func (s *DynamoDbStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	out := make(map[string]any, len(keys))

	for _, key := range keys {
		if v, err := s.Get(ctx, key); err == nil {
			out[key] = v
		}
	}

	return out, nil
}

func (s *DynamoDbStore) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	encoded, err := json.Marshal(value)

	if err != nil {
		return err
	}

	item := map[string]any{
		s.keyAttr: s.prefixed(key),
		s.valAttr: string(encoded),
	}

	if ttl > 0 {
		item[s.ttlAttr] = s.now().Add(ttl).Unix()
	}

	return s.client.PutItem(ctx, s.table, item)
}

func (s *DynamoDbStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	for k, v := range values {
		if err := s.Put(ctx, k, v, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (s *DynamoDbStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	if _, err := s.Get(ctx, key); err == nil {
		return false, nil
	}

	return true, s.Put(ctx, key, value, ttl)
}

func (s *DynamoDbStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *DynamoDbStore) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	current, err := s.Get(ctx, key)

	var val int64

	if err == nil {
		val, err = toInt64(current)

		if err != nil {
			return 0, fmt.Errorf("%w: key %q", ErrInvalidValue, key)
		}
	}

	result := val + delta

	return result, s.Put(ctx, key, result, 0)
}

func (s *DynamoDbStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *DynamoDbStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	v, err := s.Get(ctx, key)

	if err != nil {
		return false, nil
	}

	return true, s.Put(ctx, key, v, ttl)
}

func (s *DynamoDbStore) Forget(ctx context.Context, key string) error {
	return s.client.DeleteItem(ctx, s.table, s.keyAttr, s.prefixed(key))
}

func (s *DynamoDbStore) Flush(ctx context.Context) error {
	items, err := s.client.Scan(ctx, s.table)

	if err != nil {
		return err
	}

	for _, item := range items {
		key, ok := item[s.keyAttr].(string)

		if !ok {
			continue
		}

		if err := s.client.DeleteItem(ctx, s.table, s.keyAttr, key); err != nil {
			return err
		}
	}

	return nil
}
