package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type laravelConcurrencyExceptionWithoutParam struct {
	message string
}

type laravelConcurrencyAPIError struct {
	uri          string
	statusCode   int
	reason       string
	responseBody string
}

func (e laravelConcurrencyExceptionWithoutParam) Error() string {
	return e.message
}

func (e *laravelConcurrencyAPIError) Error() string {
	return fmt.Sprintf("API request to %s failed with status %d %s", e.uri, e.statusCode, e.reason)
}

// Port of Illuminate\Tests\Integration\Concurrency\ConcurrencyTest::testWorkCanBeDistributed
func TestLaravelConcurrencyWorkCanBeDistributed(t *testing.T) {
	t.Parallel()

	results, err := NewGoroutineDriver(0).Run(context.Background(), []Task{
		func() (any, error) { return 1 + 1, nil },
		func() (any, error) { return 2 + 2, nil },
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("results length = %d, want 2", len(results))
	}

	if results[0] != 2 || results[1] != 4 {
		t.Fatalf("results = %v, want [2 4]", results)
	}
}

// Port of Illuminate\Tests\Integration\Concurrency\ConcurrencyTest::testRunHandlerProcessErrorWithDefaultExceptionWithoutParam
func TestLaravelConcurrencyRunHandlerErrorWithDefaultError(t *testing.T) {
	t.Parallel()

	taskErr := errors.New("This is a different exception")

	_, err := NewGoroutineDriver(0).Run(context.Background(), []Task{
		func() (any, error) { return nil, taskErr },
	})

	if !errors.Is(err, taskErr) {
		t.Fatalf("error = %v, want %v", err, taskErr)
	}

	if err.Error() != "This is a different exception" {
		t.Fatalf("error message = %q", err.Error())
	}
}

// Port of Illuminate\Tests\Integration\Concurrency\ConcurrencyTest::testRunHandlerProcessErrorWithCustomExceptionWithoutParam
func TestLaravelConcurrencyRunHandlerErrorWithCustomError(t *testing.T) {
	t.Parallel()

	taskErr := laravelConcurrencyExceptionWithoutParam{message: "Test"}

	_, err := NewGoroutineDriver(0).Run(context.Background(), []Task{
		func() (any, error) { return nil, taskErr },
	})

	if !errors.Is(err, taskErr) {
		t.Fatalf("error = %v, want %v", err, taskErr)
	}

	if err.Error() != "Test" {
		t.Fatalf("error message = %q", err.Error())
	}
}

// Port of Illuminate\Tests\Integration\Concurrency\ConcurrencyTest::testRunHandlerProcessErrorWithCustomExceptionWithParam
func TestLaravelConcurrencyRunHandlerErrorPreservesCustomErrorParameters(t *testing.T) {
	t.Parallel()

	taskErr := &laravelConcurrencyAPIError{
		uri:          "https://api.example.com",
		statusCode:   400,
		reason:       "Bad Request",
		responseBody: "Invalid payload",
	}

	_, err := NewGoroutineDriver(0).Run(context.Background(), []Task{
		func() (any, error) { return nil, taskErr },
	})

	var apiErr *laravelConcurrencyAPIError

	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %v, want *laravelConcurrencyAPIError", err)
	}

	if apiErr.uri != taskErr.uri || apiErr.statusCode != taskErr.statusCode || apiErr.reason != taskErr.reason || apiErr.responseBody != taskErr.responseBody {
		t.Fatalf("error params = %#v, want %#v", apiErr, taskErr)
	}

	if err.Error() != "API request to https://api.example.com failed with status 400 Bad Request" {
		t.Fatalf("error message = %q", err.Error())
	}
}

// Port of Illuminate\Tests\Integration\Concurrency\ConcurrencyTest::testRunPreservesCallbackOrder
func TestLaravelConcurrencyRunPreservesCallbackOrder(t *testing.T) {
	t.Parallel()

	for name, driver := range map[string]Driver{
		"sync":      NewSyncDriver(),
		"goroutine": NewGoroutineDriver(0),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			results, err := driver.Run(context.Background(), []Task{
				func() (any, error) {
					time.Sleep(30 * time.Millisecond)

					return "first", nil
				},
				func() (any, error) {
					time.Sleep(15 * time.Millisecond)

					return "second", nil
				},
				func() (any, error) {
					time.Sleep(5 * time.Millisecond)

					return "third", nil
				},
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for i, want := range []any{"first", "second", "third"} {
				if results[i] != want {
					t.Fatalf("results[%d] = %v, want %v", i, results[i], want)
				}
			}
		})
	}
}
