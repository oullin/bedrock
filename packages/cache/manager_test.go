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
