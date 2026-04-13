package queue_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestManagerRegisterAndDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q, err := m.Driver("default")

	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("expected non-nil queue")
	}

	if q.ConnectionName() != "null" {
		t.Errorf("expected 'null', got %q", q.ConnectionName())
	}
}

func TestManagerDriverCachesInstance(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	calls := 0
	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		calls++

		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q1, _ := m.Driver("default")
	q2, _ := m.Driver("default")

	if q1 != q2 {
		t.Error("expected same instance from cache")
	}

	if calls != 1 {
		t.Errorf("expected creator called once, got %d", calls)
	}
}

func TestManagerDriverInvalidDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()
	m.SetConfig("default", map[string]any{"driver": "unknown"})

	_, err := m.Driver("default")

	if err == nil {
		t.Fatal("expected error for unknown driver")
	}

	if !errors.Is(err, queue.ErrInvalidDriver) {
		t.Errorf("expected ErrInvalidDriver, got %v", err)
	}
}

func TestManagerDriverMissingConfig(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	_, err := m.Driver("unconfigured")

	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestManagerExtendIsAlias(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Extend("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q, err := m.Driver("default")

	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("expected non-nil queue from Extend")
	}
}

func TestManagerConcurrentAccess(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, _ = m.Driver("default")
		}()
	}

	wg.Wait()
}

func TestManagerConnectionIsAliasForDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q, err := m.Connection("default")

	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("expected non-nil queue from Connection")
	}
}

func TestManagerDefaultConnection(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	if m.GetDefaultConnection() != "" {
		t.Errorf("expected empty default, got %q", m.GetDefaultConnection())
	}

	m.SetDefaultConnection("redis")

	if m.GetDefaultConnection() != "redis" {
		t.Errorf("expected 'redis', got %q", m.GetDefaultConnection())
	}
}

func TestManagerPurge(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q1, _ := m.Driver("default")
	m.Purge("default")
	q2, _ := m.Driver("default")

	if q1 == q2 {
		t.Error("expected different instance after purge")
	}
}

func TestManagerForgetDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	m.ForgetDriver("null")

	_, err := m.Driver("default")

	if err == nil {
		t.Fatal("expected error after ForgetDriver")
	}
}

func TestManagerCreatorError(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("bad", func(_ map[string]any) (queue.Queue, error) {
		return nil, errors.New("create failed")
	})

	m.SetConfig("default", map[string]any{"driver": "bad"})

	_, err := m.Driver("default")

	if err == nil {
		t.Fatal("expected error from creator")
	}
}
