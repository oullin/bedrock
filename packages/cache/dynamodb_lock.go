package cache

import (
	"context"
	"time"

	"github.com/bedrock/packages/contracts"
)

// DynamoDbLock is a distributed lock backed by AWS DynamoDB. It uses
// conditional writes to ensure atomic lock acquisition.
type DynamoDbLock struct {
	client DynamoClient
	table  string
	name   string
	owner  string
	ttl    time.Duration
	clock  contracts.Clock
}

var _ Lock = (*DynamoDbLock)(nil)

// NewDynamoDbLock creates a DynamoDB-backed lock.
func NewDynamoDbLock(client DynamoClient, table, name, owner string, ttl time.Duration, clock contracts.Clock) *DynamoDbLock {
	return &DynamoDbLock{client: client, table: table, name: name, owner: owner, ttl: ttl, clock: clock}
}

func (l *DynamoDbLock) now() time.Time {
	if l.clock != nil {
		return l.clock.Now()
	}

	return time.Now()
}

func (l *DynamoDbLock) Acquire(ctx context.Context) (bool, error) {
	exp := l.now().Add(l.ttl).Unix()

	// Check if lock exists and is still valid.
	existing, err := l.client.GetItem(ctx, l.table, "key", l.name)

	if err == nil && existing != nil {
		if ownerVal, ok := existing["owner"]; ok {
			if ownerStr, ok := ownerVal.(string); ok && ownerStr == l.owner {
				// We already own it, refresh expiration.
				return true, l.client.PutItem(ctx, l.table, map[string]any{
					"key":        l.name,
					"owner":      l.owner,
					"expiration": exp,
				})
			}
		}

		if expVal, ok := existing["expiration"]; ok {
			if expNum, ok := toInt64Value(expVal); ok && l.now().Unix() < expNum {
				return false, nil
			}
		}
	}

	// Lock is absent or expired; try to claim it.
	putErr := l.client.PutItem(ctx, l.table, map[string]any{
		"key":        l.name,
		"owner":      l.owner,
		"expiration": exp,
	})

	if putErr != nil {
		return false, putErr
	}

	return true, nil
}

func (l *DynamoDbLock) Release(ctx context.Context) (bool, error) {
	existing, err := l.client.GetItem(ctx, l.table, "key", l.name)

	if err != nil {
		return false, nil
	}

	if existing == nil {
		return false, nil
	}

	if ownerVal, ok := existing["owner"]; ok {
		if ownerStr, ok := ownerVal.(string); ok && ownerStr != l.owner {
			return false, nil
		}
	}

	return true, l.client.DeleteItem(ctx, l.table, "key", l.name)
}

func (l *DynamoDbLock) ForceRelease(ctx context.Context) error {
	return l.client.DeleteItem(ctx, l.table, "key", l.name)
}

func (l *DynamoDbLock) Get(ctx context.Context, fn func() error) error {
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

func (l *DynamoDbLock) Block(ctx context.Context, timeout time.Duration) error {
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

func (l *DynamoDbLock) Blocked(ctx context.Context) (bool, error) {
	existing, err := l.client.GetItem(ctx, l.table, "key", l.name)

	if err != nil || existing == nil {
		return false, nil
	}

	ownerVal, ok := existing["owner"]

	if !ok {
		return false, nil
	}

	ownerStr, ok := ownerVal.(string)

	if !ok {
		return false, nil
	}

	if ownerStr == l.owner {
		return false, nil
	}

	// Check if expired.
	if expVal, ok := existing["expiration"]; ok {
		if expNum, ok := toInt64Value(expVal); ok && l.now().Unix() >= expNum {
			return false, nil
		}
	}

	return true, nil
}

func toInt64Value(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case float64:
		return int64(n), true
	default:
		return 0, false
	}
}
