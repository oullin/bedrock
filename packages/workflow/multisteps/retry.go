package multisteps

import (
	"context"
	"time"
)

// RetryPolicy describes a job's retry behaviour. Field names mirror
// packages/queue.JobOptions so the vocabulary stays consistent.
type RetryPolicy struct {
	// MaxTries is the total number of attempts (1 = no retry).
	MaxTries int
	// Backoff is the sleep between attempts. If multiple durations are
	// provided, index i is used after attempt i; the last value repeats.
	Backoff []time.Duration
	// Timeout bounds each individual attempt via context.WithTimeout.
	Timeout time.Duration
	// MaxExceptions, when > 0, caps the number of unique error types before
	// giving up — primarily useful for transient-vs-permanent classification.
	MaxExceptions int
}

func (p *RetryPolicy) backoffFor(attempt int) time.Duration {
	if p == nil || len(p.Backoff) == 0 {
		return 0
	}

	if attempt >= len(p.Backoff) {
		return p.Backoff[len(p.Backoff)-1]
	}

	return p.Backoff[attempt]
}

// runWithRetry invokes handler under the retry policy, returning the response
// and the number of attempts made.
func runWithRetry(ctx context.Context, policy *RetryPolicy, handler func(context.Context) (any, error)) (any, int, error) {
	maxTries := 1

	if policy != nil && policy.MaxTries > 0 {
		maxTries = policy.MaxTries
	}

	var (
		lastErr error
		result  any
	)

	for attempt := 0; attempt < maxTries; attempt++ {
		attemptCtx := ctx
		cancel := func() {}

		if policy != nil && policy.Timeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, policy.Timeout)
		}

		result, lastErr = handler(attemptCtx)
		cancel()

		if lastErr == nil {
			return result, attempt + 1, nil
		}

		if ctx.Err() != nil {
			return nil, attempt + 1, ctx.Err()
		}

		if attempt+1 >= maxTries {
			break
		}

		backoff := policy.backoffFor(attempt)

		if backoff > 0 {
			timer := time.NewTimer(backoff)

			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()

				return nil, attempt + 1, ctx.Err()
			}
		}
	}

	return nil, maxTries, lastErr
}
