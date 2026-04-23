// Package mock provides an in-memory Redis fake implementing the
// redis.Client interface. It is used by unit tests across this module;
// because it lives under internal/ it is not part of the public API.
package mock

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bedrock/packages/redis"
)

// Client is an in-memory implementation of redis.Client. It supports the
// commands exercised by the Upstream-parity test suite:
//   - strings: GET, SET (with PX/NX), MGET, DEL, EXISTS, INCR, INCRBY,
//     DECR, EXPIRE
//   - hashes:  HGET, HSET, HMGET, HMSET, HGETALL, HDEL, HSETNX, HINCRBY,
//     HEXISTS
//   - lists:   LPUSH, RPUSH, LPOP, RPOP, LRANGE, LREM, LLEN
//   - sets:    SADD, SREM, SMEMBERS, SISMEMBER
//   - zsets:   ZADD, ZRANGE (basic), ZRANGEBYSCORE (basic)
//   - server:  FLUSHDB, PING, KEYS
//   - scripts: EVAL (very small interpreter covering the limiter Lua
//     scripts used by this package)
//
// Unknown commands return ErrUnexpectedReply. Tests that require
// operations outside this set should use the //go:build integration suite
// against a real Redis.
type Client struct {
	mu sync.Mutex

	strings map[string]entry
	hashes  map[string]map[string]string
	lists   map[string][]string
	sets    map[string]map[string]struct{}
	zsets   map[string][]zmember
	expires map[string]time.Time

	// Pub/sub is minimal: Subscribe returns a closed channel. Tests that
	// exercise pub/sub use the integration suite.
	clock func() time.Time
}

type entry struct{ value string }

type zmember struct {
	score  float64
	member string
}

// New constructs a fresh Client.

// SetClock installs a deterministic clock. Useful for TTL tests.

// Close is a no-op for the fake.

// Do dispatches a command by name.

// Pipeline returns a best-effort pipeline: commands are executed
// immediately in FIFO order and results are collected. Tests that care
// about true pipelining semantics belong in the integration suite.

// Subscribe / PSubscribe return an immediately-closed subscription. Pub/
// sub is covered by integration tests only.

// Data exposes a shallow copy of the string store for test assertions.

// -- command dispatch --------------------------------------------------

// cursor | MATCH pat | COUNT n — return all keys in a single page

// deterministic fake sha

// expire drops any keys whose TTL has elapsed. Called under lock.

// eval is a micro-interpreter that recognises exactly the Lua shapes used
// by the limiter scripts in packages/redis/limiters/scripts.go. It does
// not attempt to be a general Lua engine.

// ConcurrencyAcquire

// ConcurrencyRelease

// DurationAcquire

// -- helpers -----------------------------------------------------------

// Very limited glob: prefix* and *suffix and exact match.

// -- pipeline ----------------------------------------------------------

type pipeline struct {
	c    *Client
	cmds []*pipeCmd
	tx   bool
	done bool
}

type pipeCmd struct {
	name   string
	args   []any
	result any
	err    error
}

// -- subscription -------------------------------------------------------

type closedSub struct{}

func New() *Client {
	return &Client{
		strings: map[string]entry{},
		hashes:  map[string]map[string]string{},
		lists:   map[string][]string{},
		sets:    map[string]map[string]struct{}{},
		zsets:   map[string][]zmember{},
		expires: map[string]time.Time{},
		clock:   time.Now,
	}
}

func (c *Client) SetClock(fn func() time.Time) { c.clock = fn }

func (c *Client) Close() error { return nil }

func (c *Client) Do(_ context.Context, args ...any) (any, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if len(args) == 0 {
		return nil, redis.ErrUnexpectedReply
	}

	name, ok := args[0].(string)

	if !ok {
		return nil, redis.ErrUnexpectedReply
	}

	c.expire()

	return c.dispatch(strings.ToUpper(name), args[1:])
}

func (c *Client) Pipeline() redis.Pipeliner   { return &pipeline{c: c} }
func (c *Client) TxPipeline() redis.Pipeliner { return &pipeline{c: c, tx: true} }

func (c *Client) Subscribe(context.Context, ...string) redis.Subscription {
	return closedSub{}
}
func (c *Client) PSubscribe(context.Context, ...string) redis.Subscription {
	return closedSub{}
}

func (c *Client) Data() map[string]string {
	c.mu.Lock()

	defer c.mu.Unlock()

	out := make(map[string]string, len(c.strings))

	for k, v := range c.strings {
		out[k] = v.value
	}

	return out
}

func (c *Client) dispatch(cmd string, args []any) (any, error) {
	switch cmd {
	case "GET":
		k := s(args[0])
		e, ok := c.strings[k]

		if !ok {
			return nil, redis.ErrNil
		}

		return e.value, nil
	case "MGET":
		out := make([]any, len(args))

		for i, a := range args {
			if e, ok := c.strings[s(a)]; ok {
				out[i] = e.value
			} else {
				out[i] = nil
			}
		}

		return out, nil
	case "SET":
		return c.doSet(args)
	case "DEL":
		var n int64

		for _, a := range args {
			k := s(a)

			if _, ok := c.strings[k]; ok {
				delete(c.strings, k)
				delete(c.expires, k)
				n++
			}

			if _, ok := c.hashes[k]; ok {
				delete(c.hashes, k)
				delete(c.expires, k)
				n++
			}

			if _, ok := c.lists[k]; ok {
				delete(c.lists, k)
				delete(c.expires, k)
				n++
			}

			if _, ok := c.sets[k]; ok {
				delete(c.sets, k)
				delete(c.expires, k)
				n++
			}

			if _, ok := c.zsets[k]; ok {
				delete(c.zsets, k)
				delete(c.expires, k)
				n++
			}
		}

		return n, nil
	case "RENAME":
		oldKey := s(args[0])
		newKey := s(args[1])

		if v, ok := c.strings[oldKey]; ok {
			c.strings[newKey] = v
			delete(c.strings, oldKey)
		}

		if v, ok := c.hashes[oldKey]; ok {
			c.hashes[newKey] = v
			delete(c.hashes, oldKey)
		}

		if v, ok := c.lists[oldKey]; ok {
			c.lists[newKey] = v
			delete(c.lists, oldKey)
		}

		if v, ok := c.sets[oldKey]; ok {
			c.sets[newKey] = v
			delete(c.sets, oldKey)
		}

		if v, ok := c.zsets[oldKey]; ok {
			c.zsets[newKey] = v
			delete(c.zsets, oldKey)
		}

		if exp, ok := c.expires[oldKey]; ok {
			c.expires[newKey] = exp
			delete(c.expires, oldKey)
		}

		return "OK", nil
	case "EXISTS":
		var n int64

		for _, a := range args {
			if _, ok := c.strings[s(a)]; ok {
				n++
			}
		}

		return n, nil
	case "INCR":
		return c.incr(s(args[0]), 1)
	case "INCRBY":
		d, _ := strconv.ParseInt(s(args[1]), 10, 64)

		return c.incr(s(args[0]), d)
	case "DECR":
		return c.incr(s(args[0]), -1)
	case "EXPIRE":
		k := s(args[0])
		secs, _ := strconv.ParseInt(s(args[1]), 10, 64)

		if _, ok := c.strings[k]; !ok {
			return int64(0), nil
		}

		c.expires[k] = c.clock().Add(time.Duration(secs) * time.Second)

		return int64(1), nil
	case "PERSIST":
		k := s(args[0])

		if _, ok := c.expires[k]; ok {
			delete(c.expires, k)

			return int64(1), nil
		}

		return int64(0), nil
	case "HGET":
		h := c.hashes[s(args[0])]

		if h == nil {
			return nil, redis.ErrNil
		}

		v, ok := h[s(args[1])]

		if !ok {
			return nil, redis.ErrNil
		}

		return v, nil
	case "HSET":
		k := s(args[0])
		h, ok := c.hashes[k]

		if !ok {
			h = map[string]string{}
			c.hashes[k] = h
		}

		var added int64

		for i := 1; i+1 < len(args); i += 2 {
			f, v := s(args[i]), s(args[i+1])

			if _, exists := h[f]; !exists {
				added++
			}

			h[f] = v
		}

		return added, nil
	case "HMGET":
		h := c.hashes[s(args[0])]
		out := make([]any, 0, len(args)-1)

		for _, f := range args[1:] {
			if h == nil {
				out = append(out, nil)

				continue
			}

			if v, ok := h[s(f)]; ok {
				out = append(out, v)
			} else {
				out = append(out, nil)
			}
		}

		return out, nil
	case "HMSET":
		k := s(args[0])
		h, ok := c.hashes[k]

		if !ok {
			h = map[string]string{}
			c.hashes[k] = h
		}

		for i := 1; i+1 < len(args); i += 2 {
			h[s(args[i])] = s(args[i+1])
		}

		return "OK", nil
	case "HGETALL":
		h := c.hashes[s(args[0])]
		out := make([]any, 0, len(h)*2)

		for k, v := range h {
			out = append(out, k, v)
		}

		return out, nil
	case "HDEL":
		h := c.hashes[s(args[0])]

		var n int64

		for _, f := range args[1:] {
			if _, ok := h[s(f)]; ok {
				delete(h, s(f))
				n++
			}
		}

		return n, nil
	case "HSETNX":
		k := s(args[0])
		h, ok := c.hashes[k]

		if !ok {
			h = map[string]string{}
			c.hashes[k] = h
		}

		if _, exists := h[s(args[1])]; exists {
			return int64(0), nil
		}

		h[s(args[1])] = s(args[2])

		return int64(1), nil
	case "HINCRBY":
		k := s(args[0])
		h, ok := c.hashes[k]

		if !ok {
			h = map[string]string{}
			c.hashes[k] = h
		}

		cur, _ := strconv.ParseInt(h[s(args[1])], 10, 64)
		d, _ := strconv.ParseInt(s(args[2]), 10, 64)
		cur += d
		h[s(args[1])] = strconv.FormatInt(cur, 10)

		return cur, nil
	case "HEXISTS":
		h := c.hashes[s(args[0])]

		if _, ok := h[s(args[1])]; ok {
			return int64(1), nil
		}

		return int64(0), nil
	case "LPUSH":
		k := s(args[0])

		for _, v := range args[1:] {
			c.lists[k] = append([]string{s(v)}, c.lists[k]...)
		}

		return int64(len(c.lists[k])), nil
	case "RPUSH":
		k := s(args[0])

		for _, v := range args[1:] {
			c.lists[k] = append(c.lists[k], s(v))
		}

		return int64(len(c.lists[k])), nil
	case "LPOP":
		k := s(args[0])
		l := c.lists[k]

		if len(l) == 0 {
			return nil, redis.ErrNil
		}

		v := l[0]
		c.lists[k] = l[1:]

		return v, nil
	case "RPOP":
		k := s(args[0])
		l := c.lists[k]

		if len(l) == 0 {
			return nil, redis.ErrNil
		}

		v := l[len(l)-1]
		c.lists[k] = l[:len(l)-1]

		return v, nil
	case "LRANGE":
		k := s(args[0])
		start, _ := strconv.Atoi(s(args[1]))
		stop, _ := strconv.Atoi(s(args[2]))
		l := c.lists[k]

		if start < 0 {
			start = len(l) + start
		}

		if stop < 0 {
			stop = len(l) + stop
		}

		if start < 0 {
			start = 0
		}

		if stop >= len(l) {
			stop = len(l) - 1
		}

		if start > stop || len(l) == 0 {
			return []any{}, nil
		}

		out := make([]any, 0, stop-start+1)

		for i := start; i <= stop; i++ {
			out = append(out, l[i])
		}

		return out, nil
	case "LREM":
		k := s(args[0])
		count, _ := strconv.ParseInt(s(args[1]), 10, 64)
		val := s(args[2])
		l := c.lists[k]
		out := l[:0]

		var removed int64

		for _, v := range l {
			if v == val && (count == 0 || removed < count) {
				removed++

				continue
			}

			out = append(out, v)
		}

		c.lists[k] = out

		return removed, nil
	case "LLEN":
		return int64(len(c.lists[s(args[0])])), nil
	case "SADD":
		k := s(args[0])
		set, ok := c.sets[k]

		if !ok {
			set = map[string]struct{}{}
			c.sets[k] = set
		}

		var n int64

		for _, v := range args[1:] {
			key := s(v)

			if _, exists := set[key]; !exists {
				set[key] = struct{}{}
				n++
			}
		}

		return n, nil
	case "SREM":
		set := c.sets[s(args[0])]

		var n int64

		for _, v := range args[1:] {
			if _, ok := set[s(v)]; ok {
				delete(set, s(v))
				n++
			}
		}

		return n, nil
	case "SMEMBERS":
		set := c.sets[s(args[0])]
		out := make([]any, 0, len(set))

		for k := range set {
			out = append(out, k)
		}

		return out, nil
	case "SISMEMBER":
		set := c.sets[s(args[0])]

		if _, ok := set[s(args[1])]; ok {
			return int64(1), nil
		}

		return int64(0), nil
	case "SPOP":
		set := c.sets[s(args[0])]
		count := 1

		if len(args) > 1 {
			count, _ = strconv.Atoi(s(args[1]))
		}

		keys := make([]string, 0, len(set))

		for k := range set {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		if count > len(keys) {
			count = len(keys)
		}

		out := make([]any, 0, count)

		for i := 0; i < count; i++ {
			delete(set, keys[i])
			out = append(out, keys[i])
		}

		return out, nil
	case "ZADD":
		k := s(args[0])

		var added int64

		for i := 1; i+1 < len(args); i += 2 {
			score, _ := strconv.ParseFloat(s(args[i]), 64)
			m := s(args[i+1])
			found := false

			for j, zm := range c.zsets[k] {
				if zm.member == m {
					c.zsets[k][j].score = score
					found = true

					break
				}
			}

			if !found {
				c.zsets[k] = append(c.zsets[k], zmember{score: score, member: m})
				added++
			}
		}

		return added, nil
	case "ZRANGE":
		k := s(args[0])
		start, _ := strconv.Atoi(s(args[1]))
		stop, _ := strconv.Atoi(s(args[2]))
		zs := append([]zmember(nil), c.zsets[k]...)
		sortZ(zs)

		if start < 0 {
			start += len(zs)
		}

		if stop < 0 {
			stop += len(zs)
		}

		if start < 0 {
			start = 0
		}

		if stop >= len(zs) {
			stop = len(zs) - 1
		}

		if start > stop {
			return []any{}, nil
		}

		out := make([]any, 0, stop-start+1)

		for i := start; i <= stop; i++ {
			out = append(out, zs[i].member)
		}

		return out, nil
	case "ZREVRANGE":
		k := s(args[0])
		start, _ := strconv.Atoi(s(args[1]))
		stop, _ := strconv.Atoi(s(args[2]))
		zs := append([]zmember(nil), c.zsets[k]...)
		sortZ(zs)

		if start < 0 {
			start += len(zs)
		}

		if stop < 0 {
			stop += len(zs)
		}

		if start < 0 {
			start = 0
		}

		if stop >= len(zs) {
			stop = len(zs) - 1
		}

		if start > stop {
			return []any{}, nil
		}

		out := make([]any, 0, stop-start+1)

		for i := stop; i >= start; i-- {
			out = append(out, zs[i].member)
		}

		return out, nil
	case "ZRANGEBYSCORE":
		k := s(args[0])
		min, max := parseScoreBounds(args[1], args[2])
		zs := append([]zmember(nil), c.zsets[k]...)
		sortZ(zs)

		out := make([]any, 0, len(zs))

		for _, zm := range zs {
			if zm.score >= min && zm.score <= max {
				out = append(out, zm.member)
			}
		}

		return out, nil
	case "ZREVRANGEBYSCORE":
		k := s(args[0])
		max, min := parseScoreBounds(args[1], args[2])
		zs := append([]zmember(nil), c.zsets[k]...)
		sortZ(zs)

		out := make([]any, 0, len(zs))

		for i := len(zs) - 1; i >= 0; i-- {
			if zs[i].score >= min && zs[i].score <= max {
				out = append(out, zs[i].member)
			}
		}

		return out, nil
	case "ZCARD":
		return int64(len(c.zsets[s(args[0])])), nil
	case "ZINCRBY":
		k := s(args[0])
		delta, _ := strconv.ParseFloat(s(args[1]), 64)
		member := s(args[2])

		found := false

		for i, zm := range c.zsets[k] {
			if zm.member == member {
				c.zsets[k][i].score += delta
				found = true

				break
			}
		}

		if !found {
			c.zsets[k] = append(c.zsets[k], zmember{score: delta, member: member})
		}

		for _, zm := range c.zsets[k] {
			if zm.member == member {
				return zm.score, nil
			}
		}

		return delta, nil
	case "ZRANK":
		k := s(args[0])
		member := s(args[1])
		zs := append([]zmember(nil), c.zsets[k]...)
		sortZ(zs)

		for i, zm := range zs {
			if zm.member == member {
				return int64(i), nil
			}
		}

		return nil, nil
	case "ZSCORE":
		k := s(args[0])
		member := s(args[1])

		for _, zm := range c.zsets[k] {
			if zm.member == member {
				return zm.score, nil
			}
		}

		return nil, nil
	case "ZREM":
		k := s(args[0])
		zs := c.zsets[k]
		out := zs[:0]

		remove := map[string]struct{}{}

		for _, arg := range args[1:] {
			remove[s(arg)] = struct{}{}
		}

		var n int64

		for _, zm := range zs {
			if _, ok := remove[zm.member]; ok {
				n++

				continue
			}

			out = append(out, zm)
		}

		c.zsets[k] = out

		return n, nil
	case "ZREMRANGEBYSCORE":
		k := s(args[0])
		min, max := parseScoreBounds(args[1], args[2])
		zs := c.zsets[k]
		out := zs[:0]

		var n int64

		for _, zm := range zs {
			if zm.score >= min && zm.score <= max {
				n++

				continue
			}

			out = append(out, zm)
		}

		c.zsets[k] = out

		return n, nil
	case "ZREMRANGEBYRANK":
		k := s(args[0])
		start, _ := strconv.Atoi(s(args[1]))
		stop, _ := strconv.Atoi(s(args[2]))
		zs := append([]zmember(nil), c.zsets[k]...)
		sortZ(zs)

		if start < 0 {
			start += len(zs)
		}

		if stop < 0 {
			stop += len(zs)
		}

		if start < 0 {
			start = 0
		}

		if stop >= len(zs) {
			stop = len(zs) - 1
		}

		if start > stop || len(zs) == 0 {
			return int64(0), nil
		}

		remove := map[string]struct{}{}

		for i := start; i <= stop; i++ {
			remove[zs[i].member] = struct{}{}
		}

		out := c.zsets[k][:0]

		var n int64

		for _, zm := range c.zsets[k] {
			if _, ok := remove[zm.member]; ok {
				n++

				continue
			}

			out = append(out, zm)
		}

		c.zsets[k] = out

		return n, nil
	case "ZINTERSTORE":
		dest := s(args[0])
		n, _ := strconv.Atoi(s(args[1]))
		keys := make([]string, 0, n)

		for i := 0; i < n && 2+i < len(args); i++ {
			keys = append(keys, s(args[2+i]))
		}

		if len(keys) == 0 {
			c.zsets[dest] = nil

			return int64(0), nil
		}

		common := map[string]zmember{}

		for _, zm := range c.zsets[keys[0]] {
			common[zm.member] = zm
		}

		for _, key := range keys[1:] {
			next := map[string]zmember{}
			seen := map[string]struct{}{}

			for _, zm := range c.zsets[key] {
				seen[zm.member] = struct{}{}
			}

			for member, zm := range common {
				if _, ok := seen[member]; ok {
					next[member] = zm
				}
			}

			common = next
		}

		out := make([]zmember, 0, len(common))

		for _, zm := range common {
			out = append(out, zm)
		}

		sortZ(out)
		c.zsets[dest] = out

		return int64(len(out)), nil
	case "ZUNIONSTORE":
		dest := s(args[0])
		n, _ := strconv.Atoi(s(args[1]))
		keys := make([]string, 0, n)

		for i := 0; i < n && 2+i < len(args); i++ {
			keys = append(keys, s(args[2+i]))
		}

		union := map[string]zmember{}

		for _, key := range keys {
			for _, zm := range c.zsets[key] {
				union[zm.member] = zm
			}
		}

		out := make([]zmember, 0, len(union))

		for _, zm := range union {
			out = append(out, zm)
		}

		sortZ(out)
		c.zsets[dest] = out

		return int64(len(out)), nil
	case "FLUSHDB", "FLUSHALL":
		c.strings = map[string]entry{}
		c.hashes = map[string]map[string]string{}
		c.lists = map[string][]string{}
		c.sets = map[string]map[string]struct{}{}
		c.zsets = map[string][]zmember{}
		c.expires = map[string]time.Time{}

		return "OK", nil
	case "PING":
		return "PONG", nil
	case "KEYS":
		pat := s(args[0])

		var out []any

		for k := range c.strings {
			if matchPattern(pat, k) {
				out = append(out, k)
			}
		}

		return out, nil
	case "SCAN":
		return c.scanKeys(args...)
	case "HSCAN":
		return c.scanHash(args...)
	case "SSCAN":
		return c.scanSet(args...)
	case "ZSCAN":
		return c.scanZSet(args...)
	case "EVAL":
		if len(args) > 0 && strings.TrimSpace(s(args[0])) == "return 1" {
			return int64(1), nil
		}

		return c.eval(args)
	case "SCRIPT":
		sub, _ := args[0].(string)

		if strings.EqualFold(sub, "LOAD") {
			return "00", nil
		}

		if strings.EqualFold(sub, "EXISTS") {
			out := make([]any, len(args)-1)

			for i := range args[1:] {
				out[i] = int64(1)
			}

			return out, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", redis.ErrUnexpectedReply, cmd)
}

func (c *Client) doSet(args []any) (any, error) {
	if len(args) < 2 {
		return nil, redis.ErrUnexpectedReply
	}

	k := s(args[0])
	v := s(args[1])

	var px time.Duration
	nx := false

	for i := 2; i < len(args); i++ {
		switch strings.ToUpper(s(args[i])) {
		case "PX":
			i++
			ms, _ := strconv.ParseInt(s(args[i]), 10, 64)
			px = time.Duration(ms) * time.Millisecond
		case "EX":
			i++
			sec, _ := strconv.ParseInt(s(args[i]), 10, 64)
			px = time.Duration(sec) * time.Second
		case "NX":
			nx = true
		}
	}

	if nx {
		if _, exists := c.strings[k]; exists {
			return nil, nil
		}
	}

	c.strings[k] = entry{value: v}

	if px > 0 {
		c.expires[k] = c.clock().Add(px)
	} else {
		delete(c.expires, k)
	}

	return "OK", nil
}

func (c *Client) incr(k string, delta int64) (any, error) {
	cur, _ := strconv.ParseInt(c.strings[k].value, 10, 64)
	cur += delta
	c.strings[k] = entry{value: strconv.FormatInt(cur, 10)}

	return cur, nil
}

func (c *Client) expire() {
	now := c.clock()

	for k, t := range c.expires {
		if now.After(t) {
			delete(c.strings, k)
			delete(c.hashes, k)
			delete(c.lists, k)
			delete(c.sets, k)
			delete(c.zsets, k)
			delete(c.expires, k)
		}
	}
}

func (c *Client) eval(args []any) (any, error) {
	if len(args) < 2 {
		return nil, redis.ErrUnexpectedReply
	}

	script := s(args[0])
	numKeys, _ := strconv.Atoi(s(args[1]))
	keys := make([]string, 0, numKeys)

	for i := 0; i < numKeys; i++ {
		keys = append(keys, s(args[2+i]))
	}

	argv := args[2+numKeys:]

	switch {
	case strings.Contains(script, "RPUSH") && strings.Contains(script, "LLEN"):

		key := keys[0]
		max, _ := strconv.Atoi(s(argv[0]))
		ttl, _ := strconv.ParseInt(s(argv[1]), 10, 64)
		id := s(argv[2])

		if len(c.lists[key]) < max {
			c.lists[key] = append(c.lists[key], id)
			c.expires[key] = c.clock().Add(time.Duration(ttl) * time.Second)

			return id, nil
		}

		return int64(0), nil
	case strings.Contains(script, "LREM"):

		key := keys[0]
		id := s(argv[0])
		l := c.lists[key]
		out := l[:0]

		for _, v := range l {
			if v != id {
				out = append(out, v)
			}
		}

		c.lists[key] = out

		return int64(1), nil
	case strings.Contains(script, "HMSET") && strings.Contains(script, "HINCRBY"):

		key := keys[0]
		max, _ := strconv.Atoi(s(argv[0]))
		decay, _ := strconv.ParseInt(s(argv[1]), 10, 64)
		now, _ := strconv.ParseInt(s(argv[2]), 10, 64)

		h, ok := c.hashes[key]

		if !ok {
			h = map[string]string{
				"start": strconv.FormatInt(now, 10),
				"end":   strconv.FormatInt(now+decay, 10),
				"count": "1",
			}
			c.hashes[key] = h
			c.expires[key] = c.clock().Add(time.Duration(decay*2) * time.Second)

			return []any{int64(1), now + decay}, nil
		}

		endTs, _ := strconv.ParseInt(h["end"], 10, 64)

		if now >= endTs {
			h["start"] = strconv.FormatInt(now, 10)
			h["end"] = strconv.FormatInt(now+decay, 10)
			h["count"] = "1"
			c.expires[key] = c.clock().Add(time.Duration(decay*2) * time.Second)

			return []any{int64(1), now + decay}, nil
		}

		count, _ := strconv.Atoi(h["count"])

		if count < max {
			count++
			h["count"] = strconv.Itoa(count)

			return []any{int64(count), endTs}, nil
		}

		return []any{false, endTs}, nil
	case strings.Contains(script, "DEL") && strings.Contains(script, "return 1"):
		delete(c.hashes, keys[0])
		delete(c.lists, keys[0])
		delete(c.strings, keys[0])
		delete(c.expires, keys[0])

		return int64(1), nil
	}

	return nil, fmt.Errorf("%w: unsupported EVAL script in mock", redis.ErrUnexpectedReply)
}

func (c *Client) scanKeys(args ...any) (any, error) {
	pat := "*"

	for i := 1; i < len(args)-1; i++ {
		if strings.EqualFold(s(args[i]), "MATCH") {
			pat = s(args[i+1])
		}
	}

	var keys []any

	for k := range c.strings {
		if matchPattern(pat, k) {
			keys = append(keys, k)
		}
	}

	return []any{"0", keys}, nil
}

func (c *Client) scanHash(args ...any) (any, error) {
	key := s(args[0])
	h := c.hashes[key]

	var out []any

	for field, value := range h {
		out = append(out, field, value)
	}

	return []any{"0", out}, nil
}

func (c *Client) scanSet(args ...any) (any, error) {
	key := s(args[0])
	set := c.sets[key]

	var out []any

	for member := range set {
		out = append(out, member)
	}

	return []any{"0", out}, nil
}

func (c *Client) scanZSet(args ...any) (any, error) {
	key := s(args[0])
	zs := append([]zmember(nil), c.zsets[key]...)
	sortZ(zs)

	var out []any

	for _, zm := range zs {
		out = append(out, zm.member, strconv.FormatFloat(zm.score, 'f', -1, 64))
	}

	return []any{"0", out}, nil
}

func parseScoreBounds(minArg, maxArg any) (float64, float64) {
	return parseBound(minArg, true), parseBound(maxArg, false)
}

func parseBound(v any, isMin bool) float64 {
	sv := strings.TrimSpace(s(v))

	switch strings.ToLower(sv) {
	case "-inf":
		return math.Inf(-1)
	case "+inf", "inf":
		return math.Inf(1)
	}

	if strings.HasPrefix(sv, "(") {
		sv = strings.TrimPrefix(sv, "(")
	}

	f, err := strconv.ParseFloat(sv, 64)

	if err != nil {
		if isMin {
			return math.Inf(-1)
		}

		return math.Inf(1)
	}

	return f
}

func s(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "1"
		}

		return "0"
	}

	return fmt.Sprintf("%v", v)
}

func matchPattern(pattern, key string) bool {
	if pattern == "*" || pattern == "" {
		return true
	}

	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}

	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(key, strings.TrimPrefix(pattern, "*"))
	}

	return pattern == key
}

func sortZ(zs []zmember) {
	for i := 1; i < len(zs); i++ {
		for j := i; j > 0 && zs[j-1].score > zs[j].score; j-- {
			zs[j-1], zs[j] = zs[j], zs[j-1]
		}
	}
}

func (p *pipeline) Do(ctx context.Context, args ...any) redis.Cmder {
	cmd := &pipeCmd{args: args}

	if len(args) > 0 {
		cmd.name, _ = args[0].(string)
	}

	res, err := p.c.Do(ctx, args...)
	cmd.result = res
	cmd.err = err
	p.cmds = append(p.cmds, cmd)

	return cmd
}

func (p *pipeline) Exec(context.Context) ([]redis.Cmder, error) {
	p.done = true
	out := make([]redis.Cmder, len(p.cmds))

	var firstErr error

	for i, c := range p.cmds {
		out[i] = c

		if c.err != nil && firstErr == nil && c.err != redis.ErrNil {
			firstErr = c.err
		}
	}

	return out, firstErr
}
func (p *pipeline) Discard() { p.cmds = nil }
func (p *pipeline) Len() int { return len(p.cmds) }

func (c *pipeCmd) Name() string         { return c.name }
func (c *pipeCmd) Args() []any          { return c.args }
func (c *pipeCmd) Result() (any, error) { return c.result, c.err }
func (c *pipeCmd) Err() error           { return c.err }

func (closedSub) Channel() <-chan redis.Message {
	ch := make(chan redis.Message)
	close(ch)

	return ch
}
func (closedSub) Close() error { return nil }
