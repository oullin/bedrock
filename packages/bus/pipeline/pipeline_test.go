package pipeline_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/bus/pipeline"
)

func TestPipelineExecution(t *testing.T) {
	order := []string{}
	p := pipeline.New().Through(
		func(ctx context.Context, cmd any, next pipeline.Handler) (any, error) {
			order = append(order, "p1_before")
			result, err := next(ctx, cmd)
			order = append(order, "p1_after")

			return result, err
		},
		func(ctx context.Context, cmd any, next pipeline.Handler) (any, error) {
			order = append(order, "p2_before")
			result, err := next(ctx, cmd)
			order = append(order, "p2_after")

			return result, err
		},
	)

	result, err := p.Send(context.Background(), "command", func(_ context.Context, cmd any) (any, error) {
		order = append(order, "terminal")

		return cmd.(string) + "_done", nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if result != "command_done" {
		t.Errorf("unexpected result: %v", result)
	}

	expected := []string{"p1_before", "p2_before", "terminal", "p2_after", "p1_after"}

	for i, e := range expected {
		if i >= len(order) || order[i] != e {
			t.Errorf("order[%d]: got %q, want %q", i, func() string {
				if i < len(order) {
					return order[i]
				}

				return "<missing>"
			}(), e)
		}
	}
}

func TestEmptyPipeline(t *testing.T) {
	p := pipeline.New()
	result, err := p.Send(context.Background(), "hello", func(_ context.Context, cmd any) (any, error) {
		return cmd, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if result != "hello" {
		t.Errorf("expected 'hello', got %v", result)
	}
}
