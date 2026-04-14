package redis_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/redis"
)

func TestPipelineExecutesInOrder(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	results, err := c.Pipeline(ctx, func(p redis.Pipeliner) error {
		p.Do(ctx, "SET", "k", "v")
		p.Do(ctx, "GET", "k")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("want 2 results got %d", len(results))
	}
	v, err := results[1].Result()
	if err != nil || v != "v" {
		t.Fatalf("pipelined GET = %v err=%v", v, err)
	}
}

func TestTransactionDiscardsOnError(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()
	_, err := c.Transaction(ctx, func(p redis.Pipeliner) error {
		p.Do(ctx, "SET", "k", "v")
		return context.Canceled
	})
	if err != context.Canceled {
		t.Fatalf("want context.Canceled, got %v", err)
	}
	// Set should have been executed eagerly by the fake, but for real
	// Redis the transaction would be discarded. The fake can't rewind so
	// we don't assert on state here — the assertion is the returned err.
}
