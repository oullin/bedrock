package container

import (
	"fmt"
	"sync"
	"testing"
)

// Maps to Laravel's RewindableGeneratorTest. Laravel uses RewindableGenerator
// for lazy iteration of tagged services. In Go, Tagged() returns []any (eager),
// so these tests verify correctness, repeatability, and concurrent safety.

// --- Laravel portable: testCountUsesProvidedValue ---

func TestTaggedCountMatchesRegistered(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("a", func(_ *Container) (any, error) { return 1, nil })
	c.Bind("b", func(_ *Container) (any, error) { return 2, nil })
	c.Bind("c", func(_ *Container) (any, error) { return 3, nil })
	c.Tag([]string{"a", "b", "c"}, "numbers")

	results, err := c.Tagged("numbers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("want 3 results, got %d", len(results))
	}
}

// --- Laravel portable: testCountUsesProvidedValueAsCallback ---

func TestTaggedCanBeCalledMultipleTimes(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	c.Tag([]string{"svc"}, "group")

	first, _ := c.Tagged("group")
	second, _ := c.Tagged("group")

	if len(first) != 1 || len(second) != 1 {
		t.Fatalf("want 1 result each time, got %d and %d", len(first), len(second))
	}

	if first[0] != "value" || second[0] != "value" {
		t.Fatal("both calls should return the same values")
	}
}

// --- Additional Bedrock tests ---

func TestTaggedWithSingletons(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) {
		return &struct{ Name string }{"singleton"}, nil
	})
	c.Tag([]string{"svc"}, "group")

	first, _ := c.Tagged("group")
	second, _ := c.Tagged("group")

	if first[0] != second[0] {
		t.Fatal("tagged singletons should return the same pointer")
	}
}

func TestTaggedWithScoped(t *testing.T) {
	t.Parallel()

	type svc struct{ id int }
	calls := 0

	c := New()
	c.Scoped("svc", func(_ *Container) (any, error) {
		calls++
		return &svc{id: calls}, nil
	})
	c.Tag([]string{"svc"}, "group")

	first, _ := c.Tagged("group")
	c.ForgetScopedInstances()
	second, _ := c.Tagged("group")

	if first[0] == second[0] {
		t.Fatal("tagged scoped bindings should return new instance after reset")
	}
}

func TestTaggedConcurrent(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("a", func(_ *Container) (any, error) { return "a", nil })
	c.Bind("b", func(_ *Container) (any, error) { return "b", nil })
	c.Tag([]string{"a", "b"}, "concurrent")

	const goroutines = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)

	errs := make([]error, goroutines)

	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()

			results, err := c.Tagged("concurrent")
			if err != nil {
				errs[idx] = err
				return
			}

			if len(results) != 2 {
				errs[idx] = fmt.Errorf("want 2 results, got %d", len(results))
			}
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("goroutine %d: %v", i, err)
		}
	}
}
