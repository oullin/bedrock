package limiters

// Lua scripts ported from Upstream 13.x Framework\Redis\Limiters.

// ConcurrencyAcquire is the Lua equivalent of the "increment counter then
// release slot" pattern ConcurrencyLimiter uses. The script atomically
// claims a free slot (identifier) and sets a TTL on it.
//
// KEYS[1]   - the limiter key prefix ("name:lock")
// ARGV[1]   - max locks
// ARGV[2]   - release-after seconds
// ARGV[3]   - caller identifier (unique per acquire)
//
// Returns ARGV[3] on success or 0 on failure.
const ConcurrencyAcquire = `
if redis.call('LLEN', KEYS[1]) < tonumber(ARGV[1]) then
    redis.call('RPUSH', KEYS[1], ARGV[3])
    redis.call('EXPIRE', KEYS[1], ARGV[2])
    return ARGV[3]
end
return 0
`

// ConcurrencyRelease removes a specific identifier from the slot list.
//
// KEYS[1] - limiter key
// ARGV[1] - caller identifier
const ConcurrencyRelease = `
redis.call('LREM', KEYS[1], 1, ARGV[1])
return 1
`

// DurationAcquire is the token-bucket style script used by
// DurationLimiter. It tracks a window (decaysAt) and a counter (remaining).
//
// KEYS[1] - limiter key
// ARGV[1] - max locks
// ARGV[2] - decay (seconds)
//
// Returns {remaining, decays_at}.
const DurationAcquire = `
local function reset()
    redis.call('HMSET', KEYS[1], 'start', ARGV[3], 'end', ARGV[3] + ARGV[2], 'count', 1)
    redis.call('EXPIRE', KEYS[1], ARGV[2] * 2)
    return {1, ARGV[3] + ARGV[2]}
end

if redis.call('EXISTS', KEYS[1]) == 0 then
    return reset()
end

if ARGV[3] >= tonumber(redis.call('HGET', KEYS[1], 'end')) then
    return reset()
end

if tonumber(redis.call('HGET', KEYS[1], 'count')) < tonumber(ARGV[1]) then
    return {
        redis.call('HINCRBY', KEYS[1], 'count', 1),
        tonumber(redis.call('HGET', KEYS[1], 'end')),
    }
end

return {false, tonumber(redis.call('HGET', KEYS[1], 'end'))}
`
