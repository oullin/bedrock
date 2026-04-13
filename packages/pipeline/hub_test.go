package pipeline

import (
	"context"
	"testing"
)

func TestHubDefaultPipeline(t *testing.T) {
	hub := NewHub()
	hub.Defaults(func(p *Pipeline, object any) {
		p.Send(object).
			Through(
				Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
					return next(passable.(string) + "-default")
				}),
			)
	})

	result, err := hub.Pipe(context.Background(), "input")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "input-default" {
		t.Fatalf("expected input-default, got %v", result)
	}
}

func TestHubNamedPipeline(t *testing.T) {
	hub := NewHub()
	hub.Pipeline("custom", func(p *Pipeline, object any) {
		p.Send(object).
			Through(
				Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
					return next(passable.(string) + "-custom")
				}),
			)
	})

	result, err := hub.Pipe(context.Background(), "data", "custom")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "data-custom" {
		t.Fatalf("expected data-custom, got %v", result)
	}
}

func TestHubUnknownPipelineReturnsError(t *testing.T) {
	hub := NewHub()

	_, err := hub.Pipe(context.Background(), "data", "nonexistent")

	if err == nil {
		t.Fatal("expected error for unknown pipeline")
	}
}
