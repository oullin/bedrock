package bus_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/bus"
)

func TestChainedBatchFromPendingBatch(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1", "job2"}).
		Name("test-batch").
		OnConnection("redis").
		OnQueue("high").
		AllowFailures()

	cb := bus.NewChainedBatch(pb)

	// Verify ToPendingBatch round-trips correctly.
	restored := cb.ToPendingBatch(d)
	if restored.GetName() != "test-batch" {
		t.Errorf("expected name 'test-batch', got %q", restored.GetName())
	}

	if restored.Connection() != "redis" {
		t.Errorf("expected connection 'redis', got %q", restored.Connection())
	}

	if restored.Queue() != "high" {
		t.Errorf("expected queue 'high', got %q", restored.Queue())
	}

	if !restored.AllowsFailures() {
		t.Error("expected AllowsFailures to be true")
	}
}

func TestChainedBatchToPendingBatch(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"j1"}).Name("round-trip")

	cb := bus.NewChainedBatch(pb)
	restored := cb.ToPendingBatch(d)

	if restored.GetName() != "round-trip" {
		t.Errorf("expected name 'round-trip', got %q", restored.GetName())
	}

	if len(restored.Jobs()) != 1 {
		t.Errorf("expected 1 job, got %d", len(restored.Jobs()))
	}
}

func TestChainedBatchHandle(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"j1", "j2"}).Name("handle-test")

	cb := bus.NewChainedBatch(pb)
	cb.ChainJobs = []any{"next-in-chain"}

	ctx := bus.WithDispatcher(context.Background(), d)
	_, err := cb.Handle(ctx)
	if err != nil {
		t.Fatal(err)
	}

	d.mu.Lock()
	count := len(d.dispatchedQueue)
	d.mu.Unlock()

	// Should have dispatched 2 jobs from the batch.
	if count < 2 {
		t.Errorf("expected at least 2 dispatched jobs, got %d", count)
	}
}

func TestChainedBatchHandleWithoutDispatcher(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"j1"})

	cb := bus.NewChainedBatch(pb)

	// No dispatcher in context -- should be a no-op.
	_, err := cb.Handle(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPrepareNestedBatches(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"j1"}).Name("nested")

	jobs := []any{"plain-job", pb, "another-job"}
	result := bus.PrepareNestedBatches(jobs)

	if len(result) != 3 {
		t.Fatalf("expected 3 jobs, got %d", len(result))
	}

	if result[0] != "plain-job" {
		t.Error("expected first job to pass through unchanged")
	}

	if _, ok := result[1].(*bus.ChainedBatch); !ok {
		t.Error("expected second job to be *ChainedBatch")
	}

	if result[2] != "another-job" {
		t.Error("expected third job to pass through unchanged")
	}
}

func TestPrepareNestedBatchesNoPendingBatches(t *testing.T) {
	jobs := []any{"a", "b", "c"}
	result := bus.PrepareNestedBatches(jobs)

	for i, j := range result {
		if j != jobs[i] {
			t.Errorf("expected job %d to pass through unchanged", i)
		}
	}
}
