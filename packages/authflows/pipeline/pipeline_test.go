package pipeline

import (
	"context"
	"testing"
)

func TestPipelinePassesThroughStages(t *testing.T) {
	result, err := New().
		Send(1).
		Through(
			func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				return next(passable.(int) + 1)
			},
			func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				return next(passable.(int) * 3)
			},
		).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 6 {
		t.Fatalf("expected 6, got %v", result)
	}
}

func TestPipelineStageCanShortCircuit(t *testing.T) {
	result, err := New().
		Send("hello").
		Through(
			func(_ context.Context, _ any, _ func(any) (any, error)) (any, error) {
				return "short-circuited", nil
			},
			func(_ context.Context, _ any, next func(any) (any, error)) (any, error) {
				t.Fatal("second stage should not be called")
				return next(nil)
			},
		).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "short-circuited" {
		t.Fatalf("expected short-circuited, got %v", result)
	}
}

func TestPipelineThenCallsDestination(t *testing.T) {
	called := false

	_, err := New().
		Send("value").
		Through().
		Then(context.Background(), func(passable any) (any, error) {
			called = true

			if passable != "value" {
				t.Fatalf("expected value, got %v", passable)
			}

			return passable, nil
		})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("destination was not called")
	}
}

func TestPipelineEmptyStagesPassesThrough(t *testing.T) {
	result, err := New().
		Send("unchanged").
		Through().
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "unchanged" {
		t.Fatalf("expected unchanged, got %v", result)
	}
}
