package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

func TestQueueableOnConnection(t *testing.T) {
	q := &bus.Queueable{}
	q.OnConnection("redis")

	if q.Connection != "redis" {
		t.Errorf("expected Connection 'redis', got %q", q.Connection)
	}
}

func TestQueueableOnQueue(t *testing.T) {
	q := &bus.Queueable{}
	q.OnQueue("emails")

	if q.Queue != "emails" {
		t.Errorf("expected Queue 'emails', got %q", q.Queue)
	}
}

func TestQueueableWithDelay(t *testing.T) {
	q := &bus.Queueable{}
	q.WithDelay(5 * time.Second)

	if q.Delay != 5*time.Second {
		t.Errorf("expected Delay 5s, got %v", q.Delay)
	}
}

func TestQueueableChain(t *testing.T) {
	q := &bus.Queueable{}
	q.Chain("job1", "job2")

	if len(q.ChainJobs) != 2 {
		t.Errorf("expected 2 chain jobs, got %d", len(q.ChainJobs))
	}
}

func TestQueueableAppendToChain(t *testing.T) {
	q := &bus.Queueable{}
	q.Chain("job1")
	q.AppendToChain("job2", "job3")

	if len(q.ChainJobs) != 3 {
		t.Errorf("expected 3 chain jobs, got %d", len(q.ChainJobs))
	}
}

func TestQueueableFluentChaining(t *testing.T) {
	q := &bus.Queueable{}
	result := q.OnConnection("redis").OnQueue("emails").WithDelay(10 * time.Second)

	if result.Connection != "redis" || result.Queue != "emails" || result.Delay != 10*time.Second {
		t.Error("fluent chaining did not set all fields correctly")
	}
}

func TestBatchableBatching(t *testing.T) {
	b := &bus.Batchable{}
	if b.Batching() {
		t.Error("expected Batching() to be false when BatchID is empty")
	}

	b.WithBatchID("batch-123")
	if !b.Batching() {
		t.Error("expected Batching() to be true when BatchID is set")
	}
}

func TestBatchableWithBatchID(t *testing.T) {
	b := &bus.Batchable{}
	b.WithBatchID("abc-def")

	if b.BatchID != "abc-def" {
		t.Errorf("expected BatchID 'abc-def', got %q", b.BatchID)
	}
}

func TestQueueableGetQueue(t *testing.T) {
	q := &bus.Queueable{}
	q.OnQueue("emails")

	if q.GetQueue() != "emails" {
		t.Errorf("expected GetQueue() 'emails', got %q", q.GetQueue())
	}
}

func TestQueueableGetConnection(t *testing.T) {
	q := &bus.Queueable{}
	q.OnConnection("redis")

	if q.GetConnection() != "redis" {
		t.Errorf("expected GetConnection() 'redis', got %q", q.GetConnection())
	}
}

func TestQueueableWithoutDelay(t *testing.T) {
	q := &bus.Queueable{}
	q.WithDelay(5 * time.Second)
	q.WithoutDelay()

	if q.Delay != 0 {
		t.Errorf("expected Delay 0 after WithoutDelay, got %v", q.Delay)
	}
}

func TestQueueablePrependToChain(t *testing.T) {
	q := &bus.Queueable{}
	q.Chain("B", "C")
	q.PrependToChain("A")

	if len(q.ChainJobs) != 3 {
		t.Fatalf("expected 3 chain jobs, got %d", len(q.ChainJobs))
	}

	if q.ChainJobs[0] != "A" {
		t.Errorf("expected first job 'A', got %v", q.ChainJobs[0])
	}
}

func TestQueueableThrough(t *testing.T) {
	q := &bus.Queueable{}
	pipe := bus.Pipe(func(_ context.Context, cmd any, next bus.Handler) (any, error) {
		return next(context.Background(), cmd)
	})
	q.Through(pipe)

	if len(q.Middleware) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(q.Middleware))
	}
}

type chainableJob struct {
	bus.Queueable
	Name string
}

func TestQueueableAllOnConnection(t *testing.T) {
	j1 := &chainableJob{Name: "j1"}
	j2 := &chainableJob{Name: "j2"}

	q := &bus.Queueable{}
	q.Chain(j1, j2)
	q.AllOnConnection("sqs")

	if q.Connection != "sqs" {
		t.Errorf("expected connection 'sqs', got %q", q.Connection)
	}
}

func TestQueueableAllOnQueue(t *testing.T) {
	j1 := &chainableJob{Name: "j1"}
	j2 := &chainableJob{Name: "j2"}

	q := &bus.Queueable{}
	q.Chain(j1, j2)
	q.AllOnQueue("high")

	if q.Queue != "high" {
		t.Errorf("expected queue 'high', got %q", q.Queue)
	}
}

func TestBatchableReturnsBatchInstance(t *testing.T) {
	b := &bus.Batchable{}
	if b.Batch() != nil {
		t.Error("expected nil Batch before SetBatch")
	}

	batch := &bus.Batch{ID: "b-1"}
	b.SetBatch(batch)

	if b.Batch() != batch {
		t.Error("expected SetBatch/Batch round-trip to work")
	}
}
