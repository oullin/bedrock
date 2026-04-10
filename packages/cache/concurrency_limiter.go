package cache

import (
	"context"
	"fmt"
	"time"
)

// ConcurrencyLimiter limits the number of concurrent operations by using
// cache-backed slots. Each slot is a cache key that is atomically claimed
// using Add and released using Forget.
type ConcurrencyLimiter struct {
	store        Store
	name         string
	maxSlots     int
	releaseAfter time.Duration
}

// NewConcurrencyLimiter creates a ConcurrencyLimiter.
func NewConcurrencyLimiter(store Store, name string, maxSlots int, releaseAfter time.Duration) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		store:        store,
		name:         name,
		maxSlots:     maxSlots,
		releaseAfter: releaseAfter,
	}
}

// Block acquires a slot and executes fn, waiting up to timeout if all slots
// are occupied. The slot is automatically released after fn completes.
func (cl *ConcurrencyLimiter) Block(ctx context.Context, timeout time.Duration, fn func() error) error {
	deadline := time.Now().Add(timeout)

	for {
		slot, err := cl.Acquire(ctx)

		if err != nil {
			return err
		}

		if slot != "" {
			defer cl.Release(ctx, slot) //nolint:errcheck

			return fn()
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

// Acquire tries to claim the next available slot. Returns the slot ID on
// success or an empty string if all slots are occupied.
func (cl *ConcurrencyLimiter) Acquire(ctx context.Context) (string, error) {
	for i := 0; i < cl.maxSlots; i++ {
		slot := cl.slotKey(i)

		ok, err := cl.store.Add(ctx, slot, true, cl.releaseAfter)

		if err != nil {
			return "", err
		}

		if ok {
			return slot, nil
		}
	}

	return "", nil
}

// Release frees a slot.
func (cl *ConcurrencyLimiter) Release(ctx context.Context, slot string) error {
	return cl.store.Forget(ctx, slot)
}

func (cl *ConcurrencyLimiter) slotKey(index int) string {
	return fmt.Sprintf("%s:%d", cl.name, index)
}
