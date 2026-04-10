package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bedrock/packages/queue"
)

// RedisClient is the minimal Redis interface required by the Redis queue driver.
type RedisClient interface {
	// LPush prepends values to a list.
	LPush(ctx context.Context, key string, values ...any) error
	// RPop removes and returns the last element of a list.
	RPop(ctx context.Context, key string) (string, error)
	// ZAdd adds a member with a score to a sorted set.
	ZAdd(ctx context.Context, key string, score float64, member string) error
	// ZRangeByScore returns members with scores in [min,max].
	ZRangeByScore(ctx context.Context, key string, min, max float64) ([]string, error)
	// ZRem removes members from a sorted set.
	ZRem(ctx context.Context, key string, members ...any) error
	// LLen returns the length of a list.
	LLen(ctx context.Context, key string) (int64, error)
	// ZCard returns the cardinality of a sorted set.
	ZCard(ctx context.Context, key string) (int64, error)
}

// RedisDriver stores jobs in Redis lists and sorted sets.
type RedisDriver struct {
	client     RedisClient
	connection string
}

// NewRedisDriver creates a RedisDriver.

// Migrate any delayed jobs that are now due.

// Already popped.

// Redis queue does not keep a reserved set by default.

type redisJob struct{ BaseJob }

func NewRedisDriver(client RedisClient, connection string) *RedisDriver {
	return &RedisDriver{client: client, connection: connection}
}

func (d *RedisDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	return "", d.client.LPush(ctx, d.queueKey(queueName), string(payload))
}

func (d *RedisDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	score := float64(time.Now().Add(delay).Unix())

	return "", d.client.ZAdd(ctx, d.delayedKey(queueName), score, string(payload))
}

func (d *RedisDriver) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
	ids := make([]string, 0, len(payloads))

	for _, p := range payloads {
		id, err := d.Push(ctx, queueName, p)

		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (d *RedisDriver) Pop(ctx context.Context, queueName string) (queue.Job, error) {

	d.migrateDue(ctx, queueName)

	raw, err := d.client.RPop(ctx, d.queueKey(queueName))

	if err != nil || raw == "" {
		return nil, queue.ErrNoJob
	}

	job := &redisJob{
		BaseJob: BaseJob{
			payload:    []byte(raw),
			queue:      queueName,
			connection: d.connection,
		},
	}
	job.deleteFunc = func() error { return nil }
	job.releaseFunc = func(delay time.Duration) error {
		if delay > 0 {
			_, err := d.PushDelayed(ctx, queueName, []byte(raw), delay)

			return err
		}

		_, err := d.Push(ctx, queueName, []byte(raw))

		return err
	}

	job.failFunc = func(err error) error {
		errMsg := ""

		if err != nil {
			errMsg = err.Error()
		}

		failed := map[string]string{"exception": errMsg, "payload": raw}
		b, _ := json.Marshal(failed)

		return d.client.LPush(ctx, d.failedKey(queueName), string(b))
	}

	return job, nil
}

func (d *RedisDriver) Size(ctx context.Context, queueName string) (int64, error) {
	return d.client.LLen(ctx, d.queueKey(queueName))
}

func (d *RedisDriver) PendingSize(ctx context.Context, queueName string) (int64, error) {
	return d.client.LLen(ctx, d.queueKey(queueName))
}

func (d *RedisDriver) DelayedSize(ctx context.Context, queueName string) (int64, error) {
	return d.client.ZCard(ctx, d.delayedKey(queueName))
}

func (d *RedisDriver) ReservedSize(_ context.Context, _ string) (int64, error) {

	return 0, nil
}

func (d *RedisDriver) ConnectionName() string { return d.connection }

func (d *RedisDriver) migrateDue(ctx context.Context, queueName string) {
	now := float64(time.Now().Unix())
	due, err := d.client.ZRangeByScore(ctx, d.delayedKey(queueName), 0, now)

	if err != nil || len(due) == 0 {
		return
	}

	members := make([]any, len(due))

	for i, m := range due {
		members[i] = m
		_ = d.client.LPush(ctx, d.queueKey(queueName), m)
	}

	_ = d.client.ZRem(ctx, d.delayedKey(queueName), members...)
}

func (d *RedisDriver) queueKey(q string) string   { return fmt.Sprintf("queues:%s", q) }
func (d *RedisDriver) delayedKey(q string) string { return fmt.Sprintf("queues:%s:delayed", q) }
func (d *RedisDriver) failedKey(q string) string  { return fmt.Sprintf("queues:%s:failed", q) }
