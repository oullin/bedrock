package concurrency

import (
	"errors"
	"sync"
	"testing"
)

func TestManagerRegisterAndDriver(t *testing.T) {
	t.Parallel()

	m := NewManager()
	m.Register("goroutine", func(config map[string]any) (Driver, error) {
		return NewGoroutineDriver(0), nil
	})
	m.SetConfig("default", map[string]any{"driver": "goroutine"})

	d, err := m.Driver("default")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d == nil {
		t.Fatal("expected non-nil driver")
	}
}

func TestManagerDriverCachesInstance(t *testing.T) {
	t.Parallel()

	m := NewManager()
	calls := 0

	m.Register("goroutine", func(config map[string]any) (Driver, error) {
		calls++

		return NewGoroutineDriver(0), nil
	})
	m.SetConfig("default", map[string]any{"driver": "goroutine"})

	d1, _ := m.Driver("default")
	d2, _ := m.Driver("default")

	if d1 != d2 {
		t.Error("expected same instance from cache")
	}

	if calls != 1 {
		t.Errorf("creator called %d times, want 1", calls)
	}
}

func TestManagerDriverInvalidDriver(t *testing.T) {
	t.Parallel()

	m := NewManager()
	m.SetConfig("conn", map[string]any{"driver": "unknown"})

	_, err := m.Driver("conn")

	if !errors.Is(err, ErrInvalidDriver) {
		t.Errorf("expected ErrInvalidDriver, got %v", err)
	}
}

func TestManagerDriverMissingConfig(t *testing.T) {
	t.Parallel()

	m := NewManager()

	_, err := m.Driver("nonexistent")

	if !errors.Is(err, ErrInvalidDriver) {
		t.Errorf("expected ErrInvalidDriver, got %v", err)
	}
}

func TestManagerExtendIsAlias(t *testing.T) {
	t.Parallel()

	m := NewManager()
	m.Extend("sync", func(config map[string]any) (Driver, error) {
		return NewSyncDriver(), nil
	})
	m.SetConfig("test", map[string]any{"driver": "sync"})

	d, err := m.Driver("test")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := d.(*SyncDriver); !ok {
		t.Error("expected SyncDriver")
	}
}

func TestManagerConcurrentAccess(t *testing.T) {
	t.Parallel()

	m := NewManager()
	m.Register("goroutine", func(config map[string]any) (Driver, error) {
		return NewGoroutineDriver(0), nil
	})
	m.SetConfig("default", map[string]any{"driver": "goroutine"})

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := m.Driver("default")

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestManagerDefaultConnection(t *testing.T) {
	t.Parallel()

	m := NewManager()

	if got := m.GetDefaultConnection(); got != "" {
		t.Errorf("default = %q, want empty", got)
	}

	m.SetDefaultConnection("primary")

	if got := m.GetDefaultConnection(); got != "primary" {
		t.Errorf("default = %q, want primary", got)
	}
}

func TestManagerPurge(t *testing.T) {
	t.Parallel()

	m := NewManager()
	calls := 0

	m.Register("sync", func(config map[string]any) (Driver, error) {
		calls++

		return NewSyncDriver(), nil
	})
	m.SetConfig("conn", map[string]any{"driver": "sync"})

	_, _ = m.Driver("conn")
	m.Purge("conn")
	_, _ = m.Driver("conn")

	if calls != 2 {
		t.Errorf("creator called %d times, want 2 (after purge)", calls)
	}
}

func TestManagerForgetDriver(t *testing.T) {
	t.Parallel()

	m := NewManager()
	m.Register("sync", func(config map[string]any) (Driver, error) {
		return NewSyncDriver(), nil
	})
	m.SetConfig("conn", map[string]any{"driver": "sync"})

	m.ForgetDriver("sync")

	_, err := m.Driver("conn")

	if !errors.Is(err, ErrInvalidDriver) {
		t.Errorf("expected ErrInvalidDriver after forget, got %v", err)
	}
}

func TestManagerCreatorError(t *testing.T) {
	t.Parallel()

	creatorErr := errors.New("init failed")
	m := NewManager()

	m.Register("broken", func(config map[string]any) (Driver, error) {
		return nil, creatorErr
	})
	m.SetConfig("conn", map[string]any{"driver": "broken"})

	_, err := m.Driver("conn")

	if !errors.Is(err, creatorErr) {
		t.Errorf("expected creator error, got %v", err)
	}
}

func TestManagerConnectionIsAlias(t *testing.T) {
	t.Parallel()

	m := NewManager()
	m.Register("sync", func(config map[string]any) (Driver, error) {
		return NewSyncDriver(), nil
	})
	m.SetConfig("conn", map[string]any{"driver": "sync"})

	d1, _ := m.Driver("conn")
	d2, _ := m.Connection("conn")

	if d1 != d2 {
		t.Error("Connection and Driver should return the same instance")
	}
}
