package bus_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/bus"
)

type batchableTestJob struct {
	bus.Batchable
	Name string
}

func TestPendingBatchFluentAPI(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"})

	result := pb.Name("my-batch").
		Add("job2", "job3").
		OnConnection("redis").
		OnQueue("emails").
		AllowFailures()

	if result == nil {
		t.Fatal("expected fluent methods to return *PendingBatch")
	}
}

func TestPendingBatchCallbackRegistration(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"})

	beforeCalled := false
	pb.Before(func(_ context.Context, _ *bus.Batch) {
		beforeCalled = true
	})
	pb.Progress(func(_ context.Context, _ *bus.Batch) {})
	pb.Then(func(_ context.Context, _ *bus.Batch) {})
	pb.Catch(func(_ context.Context, _ *bus.Batch, _ error) {})
	pb.Finally(func(_ context.Context, _ *bus.Batch) {})

	batch, err := pb.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !beforeCalled {
		t.Error("expected Before callback to be invoked during Dispatch")
	}

	if len(batch.ProgressCallbacks) != 1 {
		t.Errorf("expected 1 progress callback, got %d", len(batch.ProgressCallbacks))
	}

	if len(batch.ThenCallbacks) != 1 {
		t.Errorf("expected 1 then callback, got %d", len(batch.ThenCallbacks))
	}

	if len(batch.CatchCallbacks) != 1 {
		t.Errorf("expected 1 catch callback, got %d", len(batch.CatchCallbacks))
	}

	if len(batch.FinallyCallbacks) != 1 {
		t.Errorf("expected 1 finally callback, got %d", len(batch.FinallyCallbacks))
	}
}

func TestPendingBatchDispatchCreatesValidBatch(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1", "job2", "job3"}).Name("test-batch")

	batch, err := pb.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if batch.ID == "" {
		t.Error("expected batch ID to be non-empty")
	}

	if batch.Name != "test-batch" {
		t.Errorf("expected batch name 'test-batch', got %q", batch.Name)
	}

	if batch.TotalJobs != 3 {
		t.Errorf("expected TotalJobs 3, got %d", batch.TotalJobs)
	}

	if batch.PendingJobs != 3 {
		t.Errorf("expected PendingJobs 3, got %d", batch.PendingJobs)
	}

	if batch.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestPendingBatchDispatchCallsDispatchToQueue(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1", "job2"})

	_, err := pb.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	d.mu.Lock()
	count := len(d.dispatchedQueue)
	d.mu.Unlock()

	if count != 2 {
		t.Errorf("expected 2 DispatchToQueue calls, got %d", count)
	}
}

func TestPendingBatchDispatchStopsOnQueueError(t *testing.T) {
	d := newMockQueueingDispatcher()
	d.dispatchErrAt = 1
	d.dispatchErr = errTestFailure

	pb := bus.NewPendingBatch(d, []any{"job1", "job2", "job3"})

	_, err := pb.Dispatch(context.Background())

	if err == nil {
		t.Error("expected error from Dispatch")
	}

	d.mu.Lock()
	count := len(d.dispatchedQueue)
	d.mu.Unlock()

	if count != 2 {
		t.Errorf("expected 2 dispatched jobs (1 success + 1 failure), got %d", count)
	}
}

func TestPendingBatchAllowFailuresOption(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"}).AllowFailures()

	batch, err := pb.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	v, ok := batch.Options["allowFailures"]

	if !ok || v != true {
		t.Error("expected Options['allowFailures'] to be true")
	}
}

func TestPendingBatchDispatchStoresToRepository(t *testing.T) {
	repo := newMockBatchRepo()
	q := newMockQueue()
	d := bus.NewDispatcher(q, repo)

	pb := d.Batch([]any{"job1", "job2"}).Name("stored-batch")

	batch, err := pb.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !repo.hasCalled("Store:" + batch.ID) {
		t.Error("expected repo.Store to be called")
	}
}

func TestPendingBatchDispatchStoreErrorAborts(t *testing.T) {
	repo := newMockBatchRepo()
	repo.storeErr = errTestFailure

	q := newMockQueue()
	d := bus.NewDispatcher(q, repo)

	pb := d.Batch([]any{"job1"})
	_, err := pb.Dispatch(context.Background())

	if err == nil {
		t.Error("expected error from Store failure")
	}

	q.mu.Lock()
	count := len(q.pushes)
	q.mu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 queue pushes when Store fails, got %d", count)
	}
}

func TestPendingBatchDispatchIfTrue(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"})

	batch, err := pb.DispatchIf(context.Background(), true)

	if err != nil {
		t.Fatal(err)
	}

	if batch == nil {
		t.Error("expected non-nil batch when condition is true")
	}
}

func TestPendingBatchDispatchIfFalse(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"})

	batch, err := pb.DispatchIf(context.Background(), false)

	if err != nil {
		t.Fatal(err)
	}

	if batch != nil {
		t.Error("expected nil batch when condition is false")
	}

	d.mu.Lock()
	count := len(d.dispatchedQueue)
	d.mu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 queue dispatches when condition is false, got %d", count)
	}
}

func TestPendingBatchDispatchUnless(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"})

	batch, err := pb.DispatchUnless(context.Background(), false)

	if err != nil {
		t.Fatal(err)
	}

	if batch == nil {
		t.Error("expected non-nil batch when DispatchUnless(false)")
	}
}

func TestPendingBatchDispatchAfterResponse(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1", "job2"})

	batch, err := pb.DispatchAfterResponse(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if batch == nil {
		t.Fatal("expected non-nil batch")
	}

	d.mu.Lock()
	count := len(d.dispatchedQueue)
	d.mu.Unlock()

	if count != 2 {
		t.Errorf("expected 2 deferred dispatches, got %d", count)
	}
}

func TestPendingBatchWithOption(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"}).
		WithOption("retries", 3).
		WithOption("timeout", 60)

	batch, err := pb.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if batch.Options["retries"] != 3 {
		t.Errorf("expected retries option 3, got %v", batch.Options["retries"])
	}

	if batch.Options["timeout"] != 60 {
		t.Errorf("expected timeout option 60, got %v", batch.Options["timeout"])
	}
}

func TestPendingBatchWithOptionAndAllowFailures(t *testing.T) {
	d := newMockQueueingDispatcher()
	pb := bus.NewPendingBatch(d, []any{"job1"}).
		WithOption("custom", "value").
		AllowFailures()

	opts := pb.Options()

	if opts["custom"] != "value" {
		t.Errorf("expected custom option, got %v", opts["custom"])
	}

	if opts["allowFailures"] != true {
		t.Error("expected allowFailures to be true")
	}
}

func TestPendingBatchSetsBatchIDOnBatchableJobs(t *testing.T) {
	d := newMockQueueingDispatcher()
	job := &batchableTestJob{Name: "test"}

	pb := bus.NewPendingBatch(d, []any{job})

	batch, err := pb.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if job.BatchID != batch.ID {
		t.Errorf("expected job BatchID %q, got %q", batch.ID, job.BatchID)
	}
}

func TestPendingBatchDeletesBatchOnDispatchFailure(t *testing.T) {
	repo := newMockBatchRepo()
	q := newMockQueue()
	q.pushErr = errTestFailure

	d := bus.NewDispatcher(q, repo)

	pb := d.Batch([]any{"job1"})

	_, err := pb.Dispatch(context.Background())

	if err == nil {
		t.Error("expected error from dispatch failure")
	}

	if !repo.hasCalled("Delete:") {
		t.Error("expected repo.Delete to be called when dispatch fails after store")
	}
}
