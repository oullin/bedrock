package drivers

import (
	"context"
	"encoding/json"
	"strings"
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

// RedisDeleter is the optional Redis capability needed by ClearQueue.
type RedisDeleter interface {
	Del(ctx context.Context, keys ...string) (int64, error)
}

// RedisScanner is the optional Redis capability needed by QueueNames.
// A driver supports queue enumeration only when its underlying client
// can iterate keys matching a pattern — go-redis exposes this via
// the SCAN command, and most cluster-aware clients fan out a SCAN
// across nodes for the caller.
type RedisScanner interface {
	ScanMatch(ctx context.Context, match string) ([]string, error)
}

// RedisListRanger is the optional capability needed by PendingJobs to
// return the raw payloads currently waiting on the queue list. Without
// it the driver can still report a size (LLen) but cannot snapshot
// payloads.
type RedisListRanger interface {
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}

// RedisSortedSetRanger is the optional capability needed by DelayedJobs.
// It returns the members of a sorted set ordered by score, which
// matches the upstream ZRANGE semantics for the delayed-job set.
type RedisSortedSetRanger interface {
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}

// RedisClusterAware is implemented by RedisClient fakes/drivers that know
// whether their underlying connection is a Redis Cluster. The Redis driver
// uses this to decide whether to wrap queue names in cluster hash tags so
// that every operation for a queue lands on the same slot.
//
// The result is cached after the first call on a given RedisDriver instance,
// mirroring the upstream `RedisQueue::isClusterConnection()` behaviour.
type RedisClusterAware interface {
	IsCluster() bool
}

// RedisDriver stores jobs in Redis lists and sorted sets.
type RedisDriver struct {
	client     RedisClient
	connection string

	// isCluster caches the cluster-connection check result. It is nil
	// until the first call to isClusterConnection, which mirrors
	// the upstream null-coalescing assignment ($this->isCluster ??= ...).
	isCluster *bool
}

// SetClusterClient forces the driver's cluster-awareness flag. Tests use
// this to assert hash-tag wrapping behaviour without needing a real Redis
// cluster connection.

// isClusterConnection returns whether the underlying Redis connection is a
// cluster. The result is cached after the first call.

// getQueue returns the plain `queues:<name>` key, unchanged regardless of
// cluster mode. Mirrors `RedisQueue::getQueue()`.

// getRedisKey returns the cluster-safe Redis key for a queue. On a cluster
// connection the queue name is wrapped in `{...}` so all keys for the same
// queue hash to the same slot. If the queue name already contains a hash
// tag (per Redis cluster semantics — `{` followed later by a `}` with at
// least one character between), the name is left unchanged.
//
// Mirrors `RedisQueue::getQueueRedisKey()` and `Connection::hasHashTag()`.

// hasHashTag reports whether key contains a valid Redis cluster hash tag
// (an opening `{` followed by a `}` with at least one character in between).
// Mirrors the upstream `@bedrock\Redis\Connections\Connection::hasHashTag()`.

// NewRedisDriver creates a RedisDriver.

// Migrate any delayed jobs that are now due.

// Already popped.

// Redis queue does not keep a reserved set by default.

type redisJob struct{ BaseJob }

func (d *RedisDriver) SetClusterClient(cluster bool) {
	v := cluster
	d.isCluster = &v
}

func (d *RedisDriver) isClusterConnection() bool {
	if d.isCluster != nil {
		return *d.isCluster
	}

	var v bool

	if c, ok := d.client.(RedisClusterAware); ok {
		v = c.IsCluster()
	}

	d.isCluster = &v

	return v
}

func (d *RedisDriver) getQueue(q string) string {
	if q == "" {
		q = "default"
	}

	return "queues:" + q
}

func (d *RedisDriver) getRedisKey(q string) string {
	if q == "" {
		q = "default"
	}

	if d.isClusterConnection() && !hasHashTag(q) {
		return "queues:{" + q + "}"
	}

	return "queues:" + q
}

func hasHashTag(key string) bool {
	open := strings.IndexByte(key, '{')

	if open < 0 {
		return false
	}

	close := strings.IndexByte(key[open+1:], '}')

	return close > 0
}

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

func (d *RedisDriver) ClearQueue(ctx context.Context, queueName string) error {
	deleter, ok := d.client.(RedisDeleter)

	if !ok {
		return nil
	}

	_, err := deleter.Del(ctx, d.queueKey(queueName), d.delayedKey(queueName), d.failedKey(queueName))

	return err
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

// QueueNames enumerates the queues currently known to the underlying
// Redis instance by scanning for keys matching `queues:*`. Names
// derived from cluster-style `queues:{name}` keys are unwrapped so the
// caller receives plain logical queue names. The driver returns
// ErrNotSupported when its client cannot SCAN — that capability is
// optional via the RedisScanner interface.
func (d *RedisDriver) QueueNames(ctx context.Context) ([]string, error) {
	scanner, ok := d.client.(RedisScanner)

	if !ok {
		return nil, queue.ErrNotSupported
	}

	keys, err := scanner.ScanMatch(ctx, "queues:*")

	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})

	var out []string

	for _, k := range keys {
		name := strings.TrimPrefix(k, "queues:")

		// Strip the per-queue ":delayed" or ":failed" suffix so we
		// only report each logical queue name once.
		if i := strings.IndexByte(name, ':'); i >= 0 {
			name = name[:i]
		}

		name = strings.TrimPrefix(strings.TrimSuffix(name, "}"), "{")

		if name == "" {
			continue
		}

		if _, already := seen[name]; already {
			continue
		}

		seen[name] = struct{}{}

		out = append(out, name)
	}

	return out, nil
}

// PendingJobs returns the raw payloads currently waiting on the queue
// list. The driver requires its client to satisfy the optional
// RedisListRanger interface; otherwise it returns ErrNotSupported.
func (d *RedisDriver) PendingJobs(ctx context.Context, queueName string) ([]queue.InspectedJob, error) {
	ranger, ok := d.client.(RedisListRanger)

	if !ok {
		return nil, queue.ErrNotSupported
	}

	raws, err := ranger.LRange(ctx, d.queueKey(queueName), 0, -1)

	if err != nil {
		return nil, err
	}

	return d.snapshotsFromPayloads(queueName, raws, nil), nil
}

// DelayedJobs returns the raw payloads parked in the delayed sorted
// set for queueName. The driver requires its client to satisfy the
// optional RedisSortedSetRanger interface.
func (d *RedisDriver) DelayedJobs(ctx context.Context, queueName string) ([]queue.InspectedJob, error) {
	ranger, ok := d.client.(RedisSortedSetRanger)

	if !ok {
		return nil, queue.ErrNotSupported
	}

	raws, err := ranger.ZRange(ctx, d.delayedKey(queueName), 0, -1)

	if err != nil {
		return nil, err
	}

	return d.snapshotsFromPayloads(queueName, raws, nil), nil
}

// ReservedJobs returns ErrNotSupported on the default Redis layout —
// the current driver does not maintain a reserved-set, so there are no
// in-flight jobs the snapshot could enumerate. Distinct from a clean
// empty slice so callers can decide between "no reserved jobs" and
// "this driver cannot report reserved jobs".
func (d *RedisDriver) ReservedJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return nil, queue.ErrNotSupported
}

// snapshotsFromPayloads converts raw Redis payload strings into
// InspectedJob entries, decoding the JSON envelope for displayName,
// uuid, and createdAt the same way the database driver does.
func (d *RedisDriver) snapshotsFromPayloads(queueName string, raws []string, reservedAt *time.Time) []queue.InspectedJob {
	if len(raws) == 0 {
		return nil
	}

	out := make([]queue.InspectedJob, 0, len(raws))

	for _, raw := range raws {
		job := queue.InspectedJob{
			Queue:      queueName,
			Connection: d.connection,
			Payload:    []byte(raw),
			ReservedAt: reservedAt,
		}

		var decoded map[string]any

		if err := json.Unmarshal([]byte(raw), &decoded); err == nil {
			if v, ok := decoded["displayName"].(string); ok {
				job.Name = v
			}

			if v, ok := decoded["uuid"].(string); ok {
				job.UUID = v
			}

			if v, ok := decoded["createdAt"].(float64); ok {
				job.CreatedAt = time.Unix(int64(v), 0)
			}
		}

		out = append(out, job)
	}

	return out
}

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

func (d *RedisDriver) queueKey(q string) string   { return d.getRedisKey(q) }
func (d *RedisDriver) delayedKey(q string) string { return d.getRedisKey(q) + ":delayed" }
func (d *RedisDriver) failedKey(q string) string  { return d.getRedisKey(q) + ":failed" }
