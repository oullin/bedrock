package cache_test

import (
	"testing"

	"github.com/bedrock/packages/cache"
)

func TestManagerRegisterAndStore(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()
	store := cache.NewArrayStore()
	m.Register("array", store)

	s, err := m.Store("array")

	if err != nil {
		t.Fatal(err)
	}

	if s != store {
		t.Fatal("expected same store instance")
	}
}

func TestManagerStoreNotRegistered(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()

	_, err := m.Store("missing")

	if err == nil {
		t.Fatal("expected error for unregistered store")
	}
}

func TestManagerRepository(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()
	m.Register("array", cache.NewArrayStore())

	r, err := m.Repository("array")

	if err != nil {
		t.Fatal(err)
	}

	if r == nil {
		t.Fatal("expected non-nil repository")
	}
}

func TestManagerRepositoryNotRegistered(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()

	_, err := m.Repository("missing")

	if err == nil {
		t.Fatal("expected error for unregistered store")
	}
}

func TestManagerExtend(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()

	m.Extend("custom", func(config map[string]any) (cache.Store, error) {
		return cache.NewArrayStore(), nil
	})

	// Extend registers a factory, not an instance.
	// Store should still return not registered since no store was created.
	_, err := m.Store("custom")

	if err == nil {
		t.Fatal("expected error - Extend doesn't auto-create")
	}
}

func TestManagerConcurrentAccess(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()

	done := make(chan bool, 20)

	for i := 0; i < 10; i++ {
		go func() {
			m.Register("store", cache.NewArrayStore())
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			m.Store("store") //nolint:errcheck
			done <- true
		}()
	}

	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestManagerDefaultDriver(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()
	m.Register("redis", cache.NewArrayStore())
	m.SetDefaultDriver("redis")

	if m.GetDefaultDriver() != "redis" {
		t.Fatalf("expected 'redis', got %q", m.GetDefaultDriver())
	}

	s, err := m.Driver()

	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("expected non-nil store from Driver()")
	}
}

func TestManagerPurge(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()
	m.Register("test", cache.NewArrayStore())

	_, err := m.Store("test")

	if err != nil {
		t.Fatal("expected store to exist before purge")
	}

	m.Purge("test")

	_, err = m.Store("test")

	if err == nil {
		t.Fatal("expected error after purge")
	}
}

func TestManagerForgetDriver(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()
	m.Register("test", cache.NewArrayStore())
	m.ForgetDriver("test")

	_, err := m.Store("test")

	if err == nil {
		t.Fatal("expected error after ForgetDriver")
	}
}

func TestManagerBuild(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()
	m.Extend("array", func(config map[string]any) (cache.Store, error) {
		return cache.NewArrayStore(), nil
	})

	s, err := m.Build("array", nil)

	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("expected non-nil store from Build")
	}
}

func TestManagerBuildUnregisteredDriver(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()

	_, err := m.Build("unknown", nil)

	if err == nil {
		t.Fatal("expected error for unregistered driver")
	}
}

func TestManagerMemo(t *testing.T) {
	t.Parallel()

	m := cache.NewManager()
	m.Register("test", cache.NewArrayStore())

	memo, err := m.Memo("test")

	if err != nil {
		t.Fatal(err)
	}

	if memo == nil {
		t.Fatal("expected non-nil MemoizedStore")
	}

	if memo.Inner() == nil {
		t.Fatal("expected non-nil inner store")
	}
}
