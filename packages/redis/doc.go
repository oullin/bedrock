// Package redis is a Go port of Illuminate/Redis package.
//
// It provides a Manager that multiplexes named Connection objects, a
// Connection type exposing the full Redis command surface (typed helpers
// plus a generic Command() escape hatch), Pipeline and Transaction
// closures, Pub/Sub via callbacks, Lua script helpers, Cluster and
// Sentinel variants, and ConcurrencyLimiter and
// DurationLimiter.
//
// The package targets 100% functional parity with upstream 13.x
// Illuminate\Redis and 1:1 test-case parity with its upstream suite.
//
// The backend is github.com/redis/go-redis/v9. Connection depends on a
// small Client interface so tests can inject in-memory fakes without a
// real Redis server. An optional integration suite (build tag
// "integration") exercises the package against a live Redis via the
// REDIS_URL environment variable.
package redis
