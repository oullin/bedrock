package cache

import (
	"context"
	"sync"
	"time"
)

// arrayLock is the shared lock state for a named resource.
type arrayLock struct {
	mu        sync.Mutex
	store     *ArrayStore
	name      string
	owner     string
	expiresAt time.Time
}

// arrayLockHandle is a per-caller view of an arrayLock.
type arrayLockHandle struct {
	lock    *arrayLock
	owner   string
	ttl     time.Duration
	sleepMs int
}

func (l *arrayLock) isHeld() bool {
	if l.owner == "" {
		return false
	}

	if !l.expiresAt.IsZero() && l.store.now().After(l.expiresAt) {
		l.owner = ""
		l.expiresAt = time.Time{}

		return false
	}

	return true
}

func (h *arrayLockHandle) Acquire(_ context.Context) (bool, error) {
	h.lock.mu.Lock()

	defer h.lock.mu.Unlock()

	if h.lock.isHeld() {
		return false, nil
	}

	h.lock.owner = h.owner

	if h.ttl > 0 {
		h.lock.expiresAt = h.lock.store.now().Add(h.ttl)
	}

	return true, nil
}

func (h *arrayLockHandle) Release(_ context.Context) (bool, error) {
	h.lock.mu.Lock()

	defer h.lock.mu.Unlock()

	if h.lock.owner != h.owner {
		return false, nil
	}

	h.lock.owner = ""
	h.lock.expiresAt = time.Time{}

	return true, nil
}

func (h *arrayLockHandle) ForceRelease(_ context.Context) error {
	h.lock.mu.Lock()

	defer h.lock.mu.Unlock()

	h.lock.owner = ""
	h.lock.expiresAt = time.Time{}

	return nil
}

func (h *arrayLockHandle) Get(ctx context.Context, fn func() error) error {
	ok, err := h.Acquire(ctx)

	if err != nil {
		return err
	}

	if !ok {
		return ErrLockTimeout
	}

	defer h.Release(ctx) //nolint:errcheck

	return fn()
}

func (h *arrayLockHandle) Block(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		ok, err := h.Acquire(ctx)

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
		case <-time.After(time.Duration(h.sleepMs) * time.Millisecond):
		}
	}
}

func (h *arrayLockHandle) Owner() string { return h.owner }

func (h *arrayLockHandle) IsOwnedByCurrentProcess(_ context.Context) (bool, error) {
	h.lock.mu.Lock()

	defer h.lock.mu.Unlock()

	return h.lock.isHeld() && h.lock.owner == h.owner, nil
}

func (h *arrayLockHandle) IsOwnedBy(_ context.Context, owner string) (bool, error) {
	h.lock.mu.Lock()

	defer h.lock.mu.Unlock()

	return h.lock.isHeld() && h.lock.owner == owner, nil
}

func (h *arrayLockHandle) BetweenBlockedAttemptsSleepFor(ms int) Lock {
	h.sleepMs = ms

	return h
}

func (h *arrayLockHandle) Blocked(_ context.Context) (bool, error) {
	h.lock.mu.Lock()

	defer h.lock.mu.Unlock()

	held := h.lock.isHeld() && h.lock.owner != h.owner

	return held, nil
}
