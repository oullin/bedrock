package middleware

// ThrottleRequestsWithRedis is the Redis-backed counterpart of
// Ref: @bedrock/code-0323
// In the Go port the type accepts any [RateLimiter] implementation, so the
// distinction is whether the supplied limiter wraps a Redis client. The
// service provider in M11 wires a Redis-backed limiter from
// bedrock/packages/redis when this type is requested.
type ThrottleRequestsWithRedis struct {
	*ThrottleRequests
}

// NewThrottleRequestsWithRedis constructs the Redis variant from a limiter
// that talks to Redis. Pass an in-memory limiter for tests.
func NewThrottleRequestsWithRedis(limiter RateLimiter) *ThrottleRequestsWithRedis {
	return &ThrottleRequestsWithRedis{ThrottleRequests: NewThrottleRequests(limiter)}
}
