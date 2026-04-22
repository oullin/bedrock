package bus_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

func TestBatchFinished(t *testing.T) {
	b := &bus.Batch{PendingJobs: 0}

	if !b.Finished() {
		t.Error("expected Finished() to be true when PendingJobs == 0")
	}

	b.PendingJobs = 1

	if b.Finished() {
		t.Error("expected Finished() to be false when PendingJobs > 0")
	}
}

func TestBatchCancelled(t *testing.T) {
	b := &bus.Batch{}

	if b.Cancelled() {
		t.Error("expected Cancelled() to be false when CancelledAt is nil")
	}

	now := time.Now()
	b.CancelledAt = &now

	if !b.Cancelled() {
		t.Error("expected Cancelled() to be true when CancelledAt is set")
	}
}

func TestBatchHasFailures(t *testing.T) {
	b := &bus.Batch{FailedJobs: 0}

	if b.HasFailures() {
		t.Error("expected HasFailures() to be false when FailedJobs == 0")
	}

	b.FailedJobs = 1

	if !b.HasFailures() {
		t.Error("expected HasFailures() to be true when FailedJobs > 0")
	}
}

func TestBatchCancelWithoutRepo(t *testing.T) {
	b := &bus.Batch{ID: "batch-1"}

	if err := b.Cancel(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !b.Cancelled() {
		t.Error("expected batch to be cancelled after Cancel()")
	}
}

func TestBatchCancelWithRepo(t *testing.T) {
	repo := newMockBatchRepo()
	b := bus.NewBatchWithRepo("batch-1", repo)

	if err := b.Cancel(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !repo.hasCalled("Cancel:batch-1") {
		t.Error("expected repo.Cancel to be called with batch ID")
	}
}

func TestBatchRecordSuccessfulJobWithoutRepo(t *testing.T) {
	b := &bus.Batch{PendingJobs: 3, FailedJobs: 0}
	counts, err := b.RecordSuccessfulJob(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if counts.PendingJobs != 2 {
		t.Errorf("expected PendingJobs 2, got %d", counts.PendingJobs)
	}

	if counts.FailedJobs != 0 {
		t.Errorf("expected FailedJobs 0, got %d", counts.FailedJobs)
	}
}

func TestBatchRecordSuccessfulJobNeverNegative(t *testing.T) {
	b := &bus.Batch{PendingJobs: 0}
	counts, err := b.RecordSuccessfulJob(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if counts.PendingJobs != 0 {
		t.Errorf("expected PendingJobs 0, got %d", counts.PendingJobs)
	}
}

func TestBatchRecordSuccessfulJobWithRepo(t *testing.T) {
	repo := newMockBatchRepo()
	repo.decrementResult = &bus.UpdatedBatchJobCounts{PendingJobs: 5, FailedJobs: 1}

	b := bus.NewBatchWithRepo("batch-1", repo)
	counts, err := b.RecordSuccessfulJob(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if counts.PendingJobs != 5 || counts.FailedJobs != 1 {
		t.Errorf("expected counts from repo, got pending=%d failed=%d", counts.PendingJobs, counts.FailedJobs)
	}

	if !repo.hasCalled("DecrementPendingJobs:batch-1") {
		t.Error("expected repo.DecrementPendingJobs to be called")
	}
}

func TestBatchRecordFailedJobWithoutRepo(t *testing.T) {
	b := &bus.Batch{PendingJobs: 3, FailedJobs: 0, Options: map[string]any{"allowFailures": true}}
	counts, err := b.RecordFailedJob(context.Background(), "job-42", errTestFailure)

	if err != nil {
		t.Fatal(err)
	}

	if counts.PendingJobs != 2 {
		t.Errorf("expected PendingJobs 2, got %d", counts.PendingJobs)
	}

	if counts.FailedJobs != 1 {
		t.Errorf("expected FailedJobs 1, got %d", counts.FailedJobs)
	}

	if len(b.FailedJobIDs) != 1 || b.FailedJobIDs[0] != "job-42" {
		t.Errorf("expected FailedJobIDs [job-42], got %v", b.FailedJobIDs)
	}
}

func TestBatchRecordFailedJobWithRepo(t *testing.T) {
	repo := newMockBatchRepo()
	repo.incrementFailedResult = &bus.UpdatedBatchJobCounts{PendingJobs: 2, FailedJobs: 3}

	b := bus.NewBatchWithRepo("batch-1", repo)
	b.Options["allowFailures"] = true
	counts, err := b.RecordFailedJob(context.Background(), "job-7", errTestFailure)

	if err != nil {
		t.Fatal(err)
	}

	if counts.PendingJobs != 2 || counts.FailedJobs != 3 {
		t.Errorf("expected counts from repo, got pending=%d failed=%d", counts.PendingJobs, counts.FailedJobs)
	}

	if !repo.hasCalled("IncrementFailedJobs:batch-1:job-7") {
		t.Error("expected repo.IncrementFailedJobs to be called")
	}
}

func TestBatchFreshWithRepo(t *testing.T) {
	repo := newMockBatchRepo()
	repo.batch = &bus.Batch{ID: "batch-1", Name: "refreshed", TotalJobs: 10}

	b := bus.NewBatchWithRepo("batch-1", repo)
	fresh, err := b.Fresh(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if fresh.Name != "refreshed" || fresh.TotalJobs != 10 {
		t.Errorf("expected refreshed batch, got name=%q total=%d", fresh.Name, fresh.TotalJobs)
	}
}

func TestBatchFreshWithoutRepo(t *testing.T) {
	b := &bus.Batch{ID: "batch-1", Name: "original"}
	fresh, err := b.Fresh(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if fresh != b {
		t.Error("expected Fresh to return same batch when no repo")
	}
}

func TestBatchProgress(t *testing.T) {
	b := &bus.Batch{TotalJobs: 10, PendingJobs: 3}
	p := b.Progress()

	if p != 70 {
		t.Errorf("expected progress 70, got %v", p)
	}
}

func TestBatchProgressZeroTotal(t *testing.T) {
	b := &bus.Batch{TotalJobs: 0}
	p := b.Progress()

	if p != 100 {
		t.Errorf("expected progress 100 for zero total, got %v", p)
	}
}

func TestBatchProcessedJobs(t *testing.T) {
	b := &bus.Batch{TotalJobs: 10, PendingJobs: 3}

	if b.ProcessedJobs() != 7 {
		t.Errorf("expected 7 processed jobs, got %d", b.ProcessedJobs())
	}
}

func TestBatchDeleteWithRepo(t *testing.T) {
	repo := newMockBatchRepo()
	b := bus.NewBatchWithRepo("batch-1", repo)

	if err := b.Delete(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !repo.hasCalled("Delete:batch-1") {
		t.Error("expected repo.Delete to be called")
	}
}

func TestBatchDeleteWithoutRepo(t *testing.T) {
	b := &bus.Batch{ID: "batch-1"}

	if err := b.Delete(context.Background()); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBatchAllowsFailures(t *testing.T) {
	b := &bus.Batch{Options: map[string]any{"allowFailures": true}}

	if !b.AllowsFailures() {
		t.Error("expected AllowsFailures to be true")
	}

	b2 := &bus.Batch{Options: map[string]any{}}

	if b2.AllowsFailures() {
		t.Error("expected AllowsFailures to be false")
	}
}

func TestBatchMarshalJSON(t *testing.T) {
	b := &bus.Batch{
		ID:          "batch-1",
		Name:        "test",
		TotalJobs:   10,
		PendingJobs: 3,
		FailedJobs:  1,
		Options:     map[string]any{},
		CreatedAt:   time.Now(),
	}

	data, err := b.MarshalJSON()

	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["id"] != "batch-1" {
		t.Errorf("expected id 'batch-1', got %v", m["id"])
	}

	if m["progress"] != 70.0 {
		t.Errorf("expected progress 70, got %v", m["progress"])
	}
}

func TestBatchCanceled(t *testing.T) {
	b := &bus.Batch{}

	if b.Canceled() {
		t.Error("expected Canceled() to be false when CancelledAt is nil")
	}

	now := time.Now()
	b.CancelledAt = &now

	if !b.Canceled() {
		t.Error("expected Canceled() to be true when CancelledAt is set")
	}
}

func TestBatchHasCallbackCheckers(t *testing.T) {
	b := &bus.Batch{}

	if b.HasProgressCallbacks() {
		t.Error("expected HasProgressCallbacks false with no callbacks")
	}

	if b.HasThenCallbacks() {
		t.Error("expected HasThenCallbacks false with no callbacks")
	}

	if b.HasCatchCallbacks() {
		t.Error("expected HasCatchCallbacks false with no callbacks")
	}

	if b.HasFinallyCallbacks() {
		t.Error("expected HasFinallyCallbacks false with no callbacks")
	}

	b.ProgressCallbacks = append(b.ProgressCallbacks, func(_ context.Context, _ *bus.Batch) {})
	b.ThenCallbacks = append(b.ThenCallbacks, func(_ context.Context, _ *bus.Batch) {})
	b.CatchCallbacks = append(b.CatchCallbacks, func(_ context.Context, _ *bus.Batch, _ error) {})
	b.FinallyCallbacks = append(b.FinallyCallbacks, func(_ context.Context, _ *bus.Batch) {})

	if !b.HasProgressCallbacks() {
		t.Error("expected HasProgressCallbacks true")
	}

	if !b.HasThenCallbacks() {
		t.Error("expected HasThenCallbacks true")
	}

	if !b.HasCatchCallbacks() {
		t.Error("expected HasCatchCallbacks true")
	}

	if !b.HasFinallyCallbacks() {
		t.Error("expected HasFinallyCallbacks true")
	}
}

func TestBatchInvokeProgressCallbacks(t *testing.T) {
	called := 0
	b := &bus.Batch{
		ProgressCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { called++ },
			func(_ context.Context, _ *bus.Batch) { called++ },
		},
	}

	b.InvokeProgressCallbacks(context.Background())

	if called != 2 {
		t.Errorf("expected 2 progress callbacks, got %d", called)
	}
}

func TestBatchInvokeThenCallbacks(t *testing.T) {
	called := 0
	b := &bus.Batch{
		ThenCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { called++ },
		},
	}

	b.InvokeThenCallbacks(context.Background())

	if called != 1 {
		t.Errorf("expected 1 then callback, got %d", called)
	}
}

func TestBatchInvokeCatchCallbacks(t *testing.T) {
	var receivedErr error
	b := &bus.Batch{
		CatchCallbacks: []func(context.Context, *bus.Batch, error){
			func(_ context.Context, _ *bus.Batch, err error) { receivedErr = err },
		},
	}

	b.InvokeCatchCallbacks(context.Background(), errTestFailure)

	if receivedErr != errTestFailure {
		t.Errorf("expected catch callback to receive error, got %v", receivedErr)
	}
}

func TestBatchInvokeFinallyCallbacks(t *testing.T) {
	called := 0
	b := &bus.Batch{
		FinallyCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { called++ },
		},
	}

	b.InvokeFinallyCallbacks(context.Background())

	if called != 1 {
		t.Errorf("expected 1 finally callback, got %d", called)
	}
}

func TestRecordSuccessfulJobInvokesCallbacks(t *testing.T) {
	progressCalled := false
	thenCalled := false
	finallyCalled := false

	b := &bus.Batch{
		PendingJobs: 1,
		FailedJobs:  0,
		ProgressCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { progressCalled = true },
		},
		ThenCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { thenCalled = true },
		},
		FinallyCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { finallyCalled = true },
		},
	}

	_, err := b.RecordSuccessfulJob(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !progressCalled {
		t.Error("expected progress callback to be called")
	}

	if !thenCalled {
		t.Error("expected then callback to be called when all jobs succeeded")
	}

	if !finallyCalled {
		t.Error("expected finally callback to be called when all jobs done")
	}
}

func TestRecordSuccessfulJobDoesNotInvokeThenWhenFailuresExist(t *testing.T) {
	thenCalled := false
	finallyCalled := false

	b := &bus.Batch{
		PendingJobs: 1,
		FailedJobs:  1,
		ThenCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { thenCalled = true },
		},
		FinallyCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { finallyCalled = true },
		},
	}

	_, err := b.RecordSuccessfulJob(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if thenCalled {
		t.Error("expected then callback NOT to be called when failures exist")
	}

	if !finallyCalled {
		t.Error("expected finally callback to be called when pending reaches 0")
	}
}

func TestRecordFailedJobInvokesCatchCallbacks(t *testing.T) {
	catchCalled := false
	finallyCalled := false

	b := &bus.Batch{
		PendingJobs: 1,
		FailedJobs:  0,
		CatchCallbacks: []func(context.Context, *bus.Batch, error){
			func(_ context.Context, _ *bus.Batch, _ error) { catchCalled = true },
		},
		FinallyCallbacks: []func(context.Context, *bus.Batch){
			func(_ context.Context, _ *bus.Batch) { finallyCalled = true },
		},
	}

	_, err := b.RecordFailedJob(context.Background(), "job-1", errTestFailure)

	if err != nil {
		t.Fatal(err)
	}

	if !catchCalled {
		t.Error("expected catch callback to be called")
	}

	if !finallyCalled {
		t.Error("expected finally callback when pending reaches 0")
	}
}

func TestBatchFinishedEventOnAllJobsSuccess(t *testing.T) {
	var firedEvent any

	b := &bus.Batch{
		PendingJobs: 1,
		FailedJobs:  0,
	}

	b.SetEventFunc(func(event any) { firedEvent = event })

	_, _ = b.RecordSuccessfulJob(context.Background())

	if firedEvent == nil {
		t.Fatal("expected BatchFinished event to be fired")
	}

	if _, ok := firedEvent.(bus.BatchFinished); !ok {
		t.Errorf("expected BatchFinished event, got %T", firedEvent)
	}
}

func TestBatchCanceledEvent(t *testing.T) {
	var firedEvent any

	b := &bus.Batch{ID: "batch-1"}
	b.SetEventFunc(func(event any) { firedEvent = event })

	_ = b.Cancel(context.Background())

	if firedEvent == nil {
		t.Fatal("expected BatchCanceled event to be fired")
	}

	if _, ok := firedEvent.(bus.BatchCanceled); !ok {
		t.Errorf("expected BatchCanceled event, got %T", firedEvent)
	}
}

func TestUpdatedBatchJobCountsAllJobsRanExactlyOnce(t *testing.T) {
	tests := []struct {
		name     string
		counts   bus.UpdatedBatchJobCounts
		expected bool
	}{
		{"all done", bus.UpdatedBatchJobCounts{PendingJobs: 0, FailedJobs: 0}, true},
		{"pending remain", bus.UpdatedBatchJobCounts{PendingJobs: 1, FailedJobs: 0}, false},
		{"failures exist", bus.UpdatedBatchJobCounts{PendingJobs: 0, FailedJobs: 1}, false},
		{"both nonzero", bus.UpdatedBatchJobCounts{PendingJobs: 2, FailedJobs: 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.counts.AllJobsRanExactlyOnce(); got != tt.expected {
				t.Errorf("AllJobsRanExactlyOnce() = %v, want %v", got, tt.expected)
			}
		})
	}
}
