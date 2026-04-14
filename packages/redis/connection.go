package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Connection is the Go analogue of Illuminate\Redis\Connections\Connection.
//
// It wraps a Client and dispatches CommandExecuted events for every call
// that goes through Command. Typed helpers exist for the commands exercised
// by Laravel's test suite; for everything else, callers can invoke
// Command(ctx, "RENAME", "a", "b") directly — that pathway is Laravel's
// ->command($method, $parameters).
type Connection struct {
	name      string
	client    Client
	events    *EventDispatcher
	isCluster bool
}

// NewConnection wraps a Client as a named Connection.
func NewConnection(name string, client Client) *Connection {
	return &Connection{
		name:   name,
		client: client,
		events: NewEventDispatcher(),
	}
}

// Name returns the connection name.
func (c *Connection) Name() string { return c.name }

// SetName sets the connection name.
func (c *Connection) SetName(name string) { c.name = name }

// Client returns the underlying Client.
func (c *Connection) Client() Client { return c.client }

// IsCluster reports whether this connection is backed by a cluster client.
func (c *Connection) IsCluster() bool { return c.isCluster }

// Events returns the event dispatcher (for tests and Manager wiring).
func (c *Connection) Events() *EventDispatcher { return c.events }

// Listen registers a CommandExecuted listener and ensures events are
// enabled (parity with Laravel ->listen()).
func (c *Connection) Listen(fn func(CommandExecuted)) {
	c.events.Listen(fn)
	c.events.Enable()
}

// Close releases the underlying client.
func (c *Connection) Close() error {
	if c.client == nil {
		return ErrClosed
	}
	return c.client.Close()
}

// Command is the generic dispatcher. It executes a raw Redis command and
// fires a CommandExecuted event. All typed helpers on Connection call into
// Command so events cover every command uniformly — this matches
// Illuminate\Redis\Connections\Connection::command.
func (c *Connection) Command(ctx context.Context, name string, args ...any) (any, error) {
	start := time.Now()
	argv := make([]any, 0, len(args)+1)
	argv = append(argv, name)
	argv = append(argv, args...)

	res, err := c.client.Do(ctx, argv...)

	if c.events.Enabled() {
		c.events.Dispatch(CommandExecuted{
			Command:        strings.ToLower(name),
			Parameters:     args,
			Time:           time.Since(start),
			ConnectionName: c.name,
		})
	}
	return res, err
}

// ExecuteRaw mirrors Laravel's executeRaw($parameters): the first element
// of the slice is the command name.
func (c *Connection) ExecuteRaw(ctx context.Context, parameters []any) (any, error) {
	if len(parameters) == 0 {
		return nil, fmt.Errorf("redis: executeRaw requires at least a command name")
	}
	name, ok := parameters[0].(string)
	if !ok {
		return nil, fmt.Errorf("redis: executeRaw command must be a string")
	}
	return c.Command(ctx, name, parameters[1:]...)
}

// --------------------------------------------------------------------
// Typed helpers — strings
// --------------------------------------------------------------------

// Get returns the value at key, or ErrNil if missing.
func (c *Connection) Get(ctx context.Context, key string) (string, error) {
	v, err := c.Command(ctx, "GET", key)
	if err != nil {
		return "", err
	}
	return toString(v)
}

// MGet returns the values for the given keys. Missing keys produce nil
// entries (matching Redis's MGET semantics).
func (c *Connection) MGet(ctx context.Context, keys ...string) ([]any, error) {
	args := toAnySlice(keys)
	v, err := c.Command(ctx, "MGET", args...)
	if err != nil {
		return nil, err
	}
	return toSlice(v)
}

// Set stores value at key. expiration <= 0 means no TTL.
func (c *Connection) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := []any{key, value}
	if expiration > 0 {
		args = append(args, "PX", strconv.FormatInt(expiration.Milliseconds(), 10))
	}
	_, err := c.Command(ctx, "SET", args...)
	return err
}

// SetNX stores value only if key does not yet exist.
func (c *Connection) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	args := []any{key, value}
	if expiration > 0 {
		args = append(args, "PX", strconv.FormatInt(expiration.Milliseconds(), 10))
	}
	args = append(args, "NX")
	v, err := c.Command(ctx, "SET", args...)
	if err != nil {
		return false, err
	}
	// A successful SET ... NX returns "OK"; a rejected one returns nil.
	if v == nil {
		return false, nil
	}
	s, _ := toString(v)
	return strings.EqualFold(s, "OK"), nil
}

// Del removes one or more keys. Returns the count of deleted keys.
func (c *Connection) Del(ctx context.Context, keys ...string) (int64, error) {
	v, err := c.Command(ctx, "DEL", toAnySlice(keys)...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// Exists returns the count of keys that exist.
func (c *Connection) Exists(ctx context.Context, keys ...string) (int64, error) {
	v, err := c.Command(ctx, "EXISTS", toAnySlice(keys)...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// Incr increments key by 1.
func (c *Connection) Incr(ctx context.Context, key string) (int64, error) {
	v, err := c.Command(ctx, "INCR", key)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// IncrBy increments key by the given delta.
func (c *Connection) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	v, err := c.Command(ctx, "INCRBY", key, delta)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// Decr decrements key by 1.
func (c *Connection) Decr(ctx context.Context, key string) (int64, error) {
	v, err := c.Command(ctx, "DECR", key)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// Expire sets a TTL on key.
func (c *Connection) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	v, err := c.Command(ctx, "EXPIRE", key, int64(expiration.Seconds()))
	if err != nil {
		return false, err
	}
	n, _ := toInt64(v)
	return n == 1, nil
}

// --------------------------------------------------------------------
// Typed helpers — hashes
// --------------------------------------------------------------------

// HGet returns the value of a hash field.
func (c *Connection) HGet(ctx context.Context, key, field string) (string, error) {
	v, err := c.Command(ctx, "HGET", key, field)
	if err != nil {
		return "", err
	}
	return toString(v)
}

// HSet sets one or more hash fields. pairs must be an even-length slice of
// field/value pairs.
func (c *Connection) HSet(ctx context.Context, key string, pairs ...any) (int64, error) {
	if len(pairs)%2 != 0 {
		return 0, fmt.Errorf("redis: HSet requires an even number of field/value arguments")
	}
	args := append([]any{key}, pairs...)
	v, err := c.Command(ctx, "HSET", args...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// HMGet fetches multiple hash fields.
func (c *Connection) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	args := append([]any{key}, toAnySlice(fields)...)
	v, err := c.Command(ctx, "HMGET", args...)
	if err != nil {
		return nil, err
	}
	return toSlice(v)
}

// HMSet sets multiple hash fields. Deprecated in Redis 4.0 but kept for
// Laravel parity.
func (c *Connection) HMSet(ctx context.Context, key string, pairs ...any) error {
	if len(pairs)%2 != 0 {
		return fmt.Errorf("redis: HMSet requires an even number of field/value arguments")
	}
	args := append([]any{key}, pairs...)
	_, err := c.Command(ctx, "HMSET", args...)
	return err
}

// HSetNX sets a hash field only if it does not already exist.
func (c *Connection) HSetNX(ctx context.Context, key, field string, value any) (bool, error) {
	v, err := c.Command(ctx, "HSETNX", key, field, value)
	if err != nil {
		return false, err
	}
	n, _ := toInt64(v)
	return n == 1, nil
}

// HGetAll returns all fields and values.
func (c *Connection) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	v, err := c.Command(ctx, "HGETALL", key)
	if err != nil {
		return nil, err
	}
	return toStringMap(v)
}

// HDel deletes hash fields.
func (c *Connection) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	args := append([]any{key}, toAnySlice(fields)...)
	v, err := c.Command(ctx, "HDEL", args...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// --------------------------------------------------------------------
// Typed helpers — lists
// --------------------------------------------------------------------

// LPush prepends values to a list.
func (c *Connection) LPush(ctx context.Context, key string, values ...any) (int64, error) {
	v, err := c.Command(ctx, "LPUSH", append([]any{key}, values...)...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// RPush appends values to a list.
func (c *Connection) RPush(ctx context.Context, key string, values ...any) (int64, error) {
	v, err := c.Command(ctx, "RPUSH", append([]any{key}, values...)...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// LPop removes and returns the first element of a list.
func (c *Connection) LPop(ctx context.Context, key string) (string, error) {
	v, err := c.Command(ctx, "LPOP", key)
	if err != nil {
		return "", err
	}
	return toString(v)
}

// RPop removes and returns the last element of a list.
func (c *Connection) RPop(ctx context.Context, key string) (string, error) {
	v, err := c.Command(ctx, "RPOP", key)
	if err != nil {
		return "", err
	}
	return toString(v)
}

// LRem removes count occurrences of value from the list (parity with
// PhpRedisConnection::lrem).
func (c *Connection) LRem(ctx context.Context, key string, count int64, value any) (int64, error) {
	v, err := c.Command(ctx, "LREM", key, count, value)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// LRange returns a slice of list elements.
func (c *Connection) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	v, err := c.Command(ctx, "LRANGE", key, start, stop)
	if err != nil {
		return nil, err
	}
	return toStringSlice(v)
}

// BLPop is a blocking LPOP. timeout=0 blocks indefinitely.
func (c *Connection) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	args := append(toAnySlice(keys), timeout.Seconds())
	v, err := c.Command(ctx, "BLPOP", args...)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNil
	}
	return toStringSlice(v)
}

// BRPop is a blocking RPOP.
func (c *Connection) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	args := append(toAnySlice(keys), timeout.Seconds())
	v, err := c.Command(ctx, "BRPOP", args...)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNil
	}
	return toStringSlice(v)
}

// --------------------------------------------------------------------
// Typed helpers — sets
// --------------------------------------------------------------------

// SAdd adds members to a set.
func (c *Connection) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	v, err := c.Command(ctx, "SADD", append([]any{key}, members...)...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// SRem removes members from a set.
func (c *Connection) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	v, err := c.Command(ctx, "SREM", append([]any{key}, members...)...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// SPop removes and returns count random members.
func (c *Connection) SPop(ctx context.Context, key string, count int64) ([]string, error) {
	v, err := c.Command(ctx, "SPOP", key, count)
	if err != nil {
		return nil, err
	}
	return toStringSlice(v)
}

// SMembers returns all members of a set.
func (c *Connection) SMembers(ctx context.Context, key string) ([]string, error) {
	v, err := c.Command(ctx, "SMEMBERS", key)
	if err != nil {
		return nil, err
	}
	return toStringSlice(v)
}

// SIsMember reports whether value is a member of the set.
func (c *Connection) SIsMember(ctx context.Context, key string, value any) (bool, error) {
	v, err := c.Command(ctx, "SISMEMBER", key, value)
	if err != nil {
		return false, err
	}
	n, _ := toInt64(v)
	return n == 1, nil
}

// --------------------------------------------------------------------
// Typed helpers — sorted sets
// --------------------------------------------------------------------

// ZMember is a score/member pair for ZADD.
type ZMember struct {
	Score  float64
	Member any
}

// ZAdd adds one or more score/member pairs.
func (c *Connection) ZAdd(ctx context.Context, key string, members ...ZMember) (int64, error) {
	args := make([]any, 0, 1+len(members)*2)
	args = append(args, key)
	for _, m := range members {
		args = append(args, strconv.FormatFloat(m.Score, 'f', -1, 64), m.Member)
	}
	v, err := c.Command(ctx, "ZADD", args...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// ZRange returns the elements in the specified index range.
func (c *Connection) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	v, err := c.Command(ctx, "ZRANGE", key, start, stop)
	if err != nil {
		return nil, err
	}
	return toStringSlice(v)
}

// ZRangeByScore returns elements with scores between min and max (inclusive).
func (c *Connection) ZRangeByScore(ctx context.Context, key, min, max string) ([]string, error) {
	v, err := c.Command(ctx, "ZRANGEBYSCORE", key, min, max)
	if err != nil {
		return nil, err
	}
	return toStringSlice(v)
}

// ZRevRangeByScore returns elements with scores between max and min
// (inclusive), in descending order.
func (c *Connection) ZRevRangeByScore(ctx context.Context, key, max, min string) ([]string, error) {
	v, err := c.Command(ctx, "ZREVRANGEBYSCORE", key, max, min)
	if err != nil {
		return nil, err
	}
	return toStringSlice(v)
}

// ZInterStore computes the intersection of sorted sets.
func (c *Connection) ZInterStore(ctx context.Context, dest string, keys ...string) (int64, error) {
	args := make([]any, 0, 2+len(keys))
	args = append(args, dest, strconv.Itoa(len(keys)))
	for _, k := range keys {
		args = append(args, k)
	}
	v, err := c.Command(ctx, "ZINTERSTORE", args...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// ZUnionStore computes the union of sorted sets.
func (c *Connection) ZUnionStore(ctx context.Context, dest string, keys ...string) (int64, error) {
	args := make([]any, 0, 2+len(keys))
	args = append(args, dest, strconv.Itoa(len(keys)))
	for _, k := range keys {
		args = append(args, k)
	}
	v, err := c.Command(ctx, "ZUNIONSTORE", args...)
	if err != nil {
		return 0, err
	}
	return toInt64(v)
}

// --------------------------------------------------------------------
// Typed helpers — scans
// --------------------------------------------------------------------

// ScanResult is the reply shape for SCAN-family commands.
type ScanResult struct {
	Cursor uint64
	Values []string
}

// Scan performs a cursor-based iteration over the keyspace.
func (c *Connection) Scan(ctx context.Context, cursor uint64, match string, count int64) (ScanResult, error) {
	return c.doScan(ctx, "SCAN", cursor, match, count)
}

// HScan iterates hash fields.
func (c *Connection) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) (ScanResult, error) {
	return c.doScan(ctx, "HSCAN", cursor, match, count, key)
}

// SScan iterates set members.
func (c *Connection) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) (ScanResult, error) {
	return c.doScan(ctx, "SSCAN", cursor, match, count, key)
}

// ZScan iterates sorted set members.
func (c *Connection) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) (ScanResult, error) {
	return c.doScan(ctx, "ZSCAN", cursor, match, count, key)
}

func (c *Connection) doScan(ctx context.Context, cmd string, cursor uint64, match string, count int64, keyPrefix ...string) (ScanResult, error) {
	args := make([]any, 0, 6)
	for _, k := range keyPrefix {
		args = append(args, k)
	}
	args = append(args, strconv.FormatUint(cursor, 10))
	if match != "" {
		args = append(args, "MATCH", match)
	}
	if count > 0 {
		args = append(args, "COUNT", count)
	}
	v, err := c.Command(ctx, cmd, args...)
	if err != nil {
		return ScanResult{}, err
	}
	s, err := toSlice(v)
	if err != nil || len(s) != 2 {
		return ScanResult{}, ErrUnexpectedReply
	}
	curStr, _ := toString(s[0])
	nextCursor, _ := strconv.ParseUint(curStr, 10, 64)
	vals, _ := toStringSlice(s[1])
	return ScanResult{Cursor: nextCursor, Values: vals}, nil
}

// --------------------------------------------------------------------
// Typed helpers — server
// --------------------------------------------------------------------

// FlushDB flushes the current database.
func (c *Connection) FlushDB(ctx context.Context) error {
	_, err := c.Command(ctx, "FLUSHDB")
	return err
}

// Ping pings the server.
func (c *Connection) Ping(ctx context.Context) error {
	_, err := c.Command(ctx, "PING")
	return err
}
