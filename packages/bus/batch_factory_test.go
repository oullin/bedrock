package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

func TestBatchFactoryMake(t *testing.T) {
	d := newMockQueueingDispatcher()
	factory := bus.NewBatchFactory(d)
	repo := newMockBatchRepo()

	now := time.Now()
	batch := factory.Make(
		repo,
		"batch-1", "test-batch",
		10, 5, 2,
		[]string{"j1", "j2"},
		map[string]any{"allowFailures": true},
		now, nil, nil,
	)

	if batch.ID != "batch-1" {
		t.Errorf("expected ID 'batch-1', got %q", batch.ID)
	}

	if batch.Name != "test-batch" {
		t.Errorf("expected Name 'test-batch', got %q", batch.Name)
	}

	if batch.TotalJobs != 10 {
		t.Errorf("expected TotalJobs 10, got %d", batch.TotalJobs)
	}

	if batch.PendingJobs != 5 {
		t.Errorf("expected PendingJobs 5, got %d", batch.PendingJobs)
	}

	if batch.FailedJobs != 2 {
		t.Errorf("expected FailedJobs 2, got %d", batch.FailedJobs)
	}

	if len(batch.FailedJobIDs) != 2 {
		t.Errorf("expected 2 FailedJobIDs, got %d", len(batch.FailedJobIDs))
	}

	if !batch.AllowsFailures() {
		t.Error("expected AllowsFailures to be true")
	}

	if batch.CreatedAt != now {
		t.Errorf("expected CreatedAt %v, got %v", now, batch.CreatedAt)
	}
}

func TestBatchFactoryMakeSetsRepo(t *testing.T) {
	d := newMockQueueingDispatcher()
	factory := bus.NewBatchFactory(d)
	repo := newMockBatchRepo()

	batch := factory.Make(
		repo,
		"batch-1", "test",
		1, 1, 0, nil, nil,
		time.Now(), nil, nil,
	)

	// Verify the batch delegates to the repo for Cancel.
	if err := batch.Cancel(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !repo.hasCalled("Cancel:batch-1") {
		t.Error("expected batch to delegate Cancel to repo")
	}
}

func TestBatchFactoryMakeNilDefaults(t *testing.T) {
	d := newMockQueueingDispatcher()
	factory := bus.NewBatchFactory(d)
	repo := newMockBatchRepo()

	batch := factory.Make(
		repo,
		"batch-1", "test",
		0, 0, 0,
		nil, nil,
		time.Now(), nil, nil,
	)

	if batch.FailedJobIDs == nil {
		t.Error("expected FailedJobIDs to be initialized, not nil")
	}

	if batch.Options == nil {
		t.Error("expected Options to be initialized, not nil")
	}
}
