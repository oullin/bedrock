package cache

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/bedrock/packages/contracts"
)

// FileLock is a filesystem-based lock. It creates a lock file in the specified
// directory. The owner string is written to the file for ownership verification.
type FileLock struct {
	dir   string
	name  string
	owner string
	ttl   time.Duration
	clock contracts.Clock
}

var _ Lock = (*FileLock)(nil)

// NewFileLock creates a file-based lock.
func NewFileLock(dir, name, owner string, ttl time.Duration, clock contracts.Clock) *FileLock {
	return &FileLock{dir: dir, name: name, owner: owner, ttl: ttl, clock: clock}
}

func (l *FileLock) now() time.Time {
	if l.clock != nil {
		return l.clock.Now()
	}

	return time.Now()
}

func (l *FileLock) path() string {
	return filepath.Join(l.dir, "lock-"+l.name)
}

func (l *FileLock) Acquire(_ context.Context) (bool, error) {
	p := l.path()

	if err := os.MkdirAll(l.dir, 0o755); err != nil {
		return false, err
	}

	info, err := os.Stat(p)

	if err == nil {
		if l.ttl > 0 && l.now().Sub(info.ModTime()) > l.ttl {
			_ = os.Remove(p)
		} else {
			data, readErr := os.ReadFile(p)

			if readErr == nil && string(data) == l.owner {
				return true, nil
			}

			return false, nil
		}
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)

	if err != nil {
		if os.IsExist(err) {
			return false, nil
		}

		return false, err
	}

	defer f.Close()

	_, err = f.WriteString(l.owner)

	return err == nil, err
}

func (l *FileLock) Release(_ context.Context) (bool, error) {
	p := l.path()

	data, err := os.ReadFile(p)

	if err != nil {
		return false, nil
	}

	if string(data) != l.owner {
		return false, nil
	}

	return true, os.Remove(p)
}

func (l *FileLock) ForceRelease(_ context.Context) error {
	err := os.Remove(l.path())

	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func (l *FileLock) Get(ctx context.Context, fn func() error) error {
	ok, err := l.Acquire(ctx)

	if err != nil {
		return err
	}

	if !ok {
		return ErrLockTimeout
	}

	defer l.Release(ctx) //nolint:errcheck

	return fn()
}

func (l *FileLock) Block(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		ok, err := l.Acquire(ctx)

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
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (l *FileLock) Blocked(_ context.Context) (bool, error) {
	p := l.path()

	data, err := os.ReadFile(p)

	if err != nil {
		return false, nil
	}

	return string(data) != l.owner, nil
}
