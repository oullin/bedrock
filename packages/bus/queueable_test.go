package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

type chainableJob struct {
	bus.Queueable
	Name string
}

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

func TestQueueableChainConnection(t *testing.T) {
	q := &bus.Queueable{}
	q.ChainConnection = "redis"

	if q.ChainConnection != "redis" {
		t.Errorf("expected ChainConnection 'redis', got %q", q.ChainConnection)
	}
}

func TestQueueableChainQueue(t *testing.T) {
	q := &bus.Queueable{}
	q.ChainQueue = "high"

	if q.ChainQueue != "high" {
		t.Errorf("expected ChainQueue 'high', got %q", q.ChainQueue)
	}
}

func TestQueueableAfterCommit(t *testing.T) {
	q := &bus.Queueable{}

	if q.AfterCommit != nil {
		t.Error("expected AfterCommit to be nil by default")
	}

	q.SetAfterCommit()

	if q.AfterCommit == nil || !*q.AfterCommit {
		t.Error("expected AfterCommit to be true after SetAfterCommit")
	}

	q.SetBeforeCommit()

	if q.AfterCommit == nil || *q.AfterCommit {
		t.Error("expected AfterCommit to be false after SetBeforeCommit")
	}
}

func TestQueueableOnChainCatch(t *testing.T) {
	q := &bus.Queueable{}
	called := false

	q.OnChainCatch(func(_ context.Context, _ error) { called = true })

	if len(q.ChainCatchCallbacks) != 1 {
		t.Fatalf("expected 1 chain catch callback, got %d", len(q.ChainCatchCallbacks))
	}

	q.InvokeChainCatchCallbacks(context.Background(), errTestFailure)

	if !called {
		t.Error("expected chain catch callback to be called")
	}
}

func TestQueueableDispatchNextJobInChain(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	j1 := &chainableJob{Name: "first"}
	j2 := &chainableJob{Name: "second"}
	j3 := &chainableJob{Name: "third"}

	d.Map(j2, func(_ context.Context, cmd any) (any, error) {
		return cmd.(*chainableJob).Name, nil
	})

	j1.Chain(j2, j3)

	err := j1.DispatchNextJobInChain(context.Background(), d)

	if err != nil {
		t.Fatal(err)
	}

	if len(j2.ChainJobs) != 1 {
		t.Errorf("expected 1 remaining chain job on j2, got %d", len(j2.ChainJobs))
	}
}

func TestQueueableDispatchNextJobInChainEmpty(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	q := &bus.Queueable{}

	err := q.DispatchNextJobInChain(context.Background(), d)

	if err != nil {
		t.Errorf("expected no error for empty chain, got %v", err)
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

func TestBatchableWithFakeBatch(t *testing.T) {
	b := &bus.Batchable{}

	batch := b.WithFakeBatch("fake-1", "test-batch", 5)

	if batch.ID != "fake-1" {
		t.Errorf("expected ID 'fake-1', got %q", batch.ID)
	}

	if batch.Name != "test-batch" {
		t.Errorf("expected Name 'test-batch', got %q", batch.Name)
	}

	if batch.TotalJobs != 5 {
		t.Errorf("expected TotalJobs 5, got %d", batch.TotalJobs)
	}

	if b.BatchID != "fake-1" {
		t.Errorf("expected BatchID 'fake-1', got %q", b.BatchID)
	}

	if b.Batch() != batch {
		t.Error("expected Batch() to return the fake batch")
	}

	if !b.Batching() {
		t.Error("expected Batching() to be true with fake batch")
	}
}

func TestBatchableBatchFromRepo(t *testing.T) {
	repo := newMockBatchRepo()
	repo.batch = &bus.Batch{ID: "repo-1", Name: "from-repo"}

	b := &bus.Batchable{}
	b.WithBatchID("repo-1")

	batch, err := b.BatchFromRepo(context.Background(), repo)

	if err != nil {
		t.Fatal(err)
	}

	if batch.Name != "from-repo" {
		t.Errorf("expected Name 'from-repo', got %q", batch.Name)
	}

	// Should cache locally.
	if b.Batch() != batch {
		t.Error("expected Batch() to return cached repo batch")
	}
}

func TestBatchableBatchFromRepoNilWhenEmpty(t *testing.T) {
	b := &bus.Batchable{}

	batch, err := b.BatchFromRepo(context.Background(), nil)

	if err != nil {
		t.Fatal(err)
	}

	if batch != nil {
		t.Error("expected nil batch when no BatchID set")
	}
}

func TestBatchableBatchingReturnsFalseWhenCancelled(t *testing.T) {
	b := &bus.Batchable{}
	now := time.Now()

	batch := &bus.Batch{ID: "b-1", CancelledAt: &now}
	b.WithBatchID("b-1")
	b.SetBatch(batch)

	if b.Batching() {
		t.Error("expected Batching() to be false when batch is cancelled")
	}
}
