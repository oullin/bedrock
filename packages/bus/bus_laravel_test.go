package bus_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

type laravelBatchJob struct {
	bus.Batchable
	bus.Queueable
	Name string
}

type delayedQueueableCommand struct {
	bus.Queueable
	Value string
}

func (delayedQueueableCommand) ShouldQueue() {}

// Port of Framework\Tests\Bus\BusBatchTest::test_jobs_can_be_added_to_the_batch
func TestUpstreamBusBatchAddAssignsBatchIDAndDispatches(t *testing.T) {
	d := newMockQueueingDispatcher()
	batch := &bus.Batch{
		ID:          "batch-1",
		TotalJobs:   1,
		PendingJobs: 1,
		Options:     map[string]any{},
	}
	batch.SetDispatcher(d)

	first := &laravelBatchJob{Name: "first"}
	second := &laravelBatchJob{Name: "second"}

	if err := batch.Add(context.Background(), []any{first, second}); err != nil {
		t.Fatal(err)
	}

	if batch.TotalJobs != 3 || batch.PendingJobs != 3 {
		t.Fatalf("batch counts = total:%d pending:%d, want total:3 pending:3", batch.TotalJobs, batch.PendingJobs)
	}

	if first.BatchID != "batch-1" || second.BatchID != "batch-1" {
		t.Fatalf("added jobs were not tagged with the batch id: %q %q", first.BatchID, second.BatchID)
	}

	d.mu.Lock()
	dispatched := len(d.dispatchedQueue)
	d.mu.Unlock()

	if dispatched != 2 {
		t.Fatalf("dispatched jobs = %d, want 2", dispatched)
	}
}

// Port of Framework\Tests\Bus\BusBatchTest::test_processed_jobs_can_be_calculated
func TestUpstreamBusBatchProcessedJobsAndProgress(t *testing.T) {
	batch := &bus.Batch{TotalJobs: 10, PendingJobs: 4}

	if got := batch.ProcessedJobs(); got != 6 {
		t.Fatalf("ProcessedJobs() = %d, want 6", got)
	}

	if got := batch.Progress(); got != 60 {
		t.Fatalf("Progress() = %v, want 60", got)
	}
}

// Port of Framework\Tests\Bus\BusBatchTest::test_successful_jobs_can_be_recorded
// Port of Framework\Tests\Bus\BusBatchTest::test_batch_finished_event_is_dispatched
// Port of Framework\Tests\Bus\BusBatchTest::test_batch_started_event_is_dispatched
func TestUpstreamBusBatchSuccessfulJobsLifecycle(t *testing.T) {
	var progressCount, thenCount, finallyCount int

	var events []any

	batch := &bus.Batch{
		ID:          "batch-1",
		TotalJobs:   2,
		PendingJobs: 2,
		Options:     map[string]any{},
		ProgressCallbacks: []func(context.Context, *bus.Batch){
			func(context.Context, *bus.Batch) { progressCount++ },
		},
		ThenCallbacks: []func(context.Context, *bus.Batch){
			func(context.Context, *bus.Batch) { thenCount++ },
		},
		FinallyCallbacks: []func(context.Context, *bus.Batch){
			func(context.Context, *bus.Batch) { finallyCount++ },
		},
	}
	batch.SetEventFunc(func(event any) { events = append(events, event) })

	if _, err := batch.RecordSuccessfulJob(context.Background()); err != nil {
		t.Fatal(err)
	}

	if _, err := batch.RecordSuccessfulJob(context.Background()); err != nil {
		t.Fatal(err)
	}

	if progressCount != 2 || thenCount != 1 || finallyCount != 1 {
		t.Fatalf("callback counts progress:%d then:%d finally:%d, want 2/1/1", progressCount, thenCount, finallyCount)
	}

	if len(events) != 2 {
		t.Fatalf("events = %d, want 2", len(events))
	}

	if _, ok := events[0].(bus.BatchStarted); !ok {
		t.Fatalf("first event = %T, want bus.BatchStarted", events[0])
	}

	if _, ok := events[1].(bus.BatchFinished); !ok {
		t.Fatalf("second event = %T, want bus.BatchFinished", events[1])
	}
}

// Port of Framework\Tests\Bus\BusBatchTest::test_batch_started_event_is_dispatched_when_first_job_fails
// Port of Framework\Tests\Bus\BusBatchTest::test_failed_jobs_can_be_recorded_while_not_allowing_failures
// Port of Framework\Tests\Bus\BusBatchTest::test_failed_jobs_can_be_recorded_while_allowing_failures
// Port of Framework\Tests\Bus\BusBatchTest::test_failure_callbacks_execute_correctly
func TestUpstreamBusBatchFailedJobsLifecycle(t *testing.T) {
	t.Run("non tolerant batch cancels and finishes on first failure", func(t *testing.T) {
		var catchCount, progressCount, finallyCount int

		var receivedErr error

		var events []any

		batch := &bus.Batch{
			ID:          "batch-1",
			TotalJobs:   2,
			PendingJobs: 2,
			Options:     map[string]any{},
			ProgressCallbacks: []func(context.Context, *bus.Batch){
				func(context.Context, *bus.Batch) { progressCount++ },
			},
			CatchCallbacks: []func(context.Context, *bus.Batch, error){
				func(_ context.Context, _ *bus.Batch, err error) {
					catchCount++
					receivedErr = err
				},
			},
			FinallyCallbacks: []func(context.Context, *bus.Batch){
				func(context.Context, *bus.Batch) { finallyCount++ },
			},
		}
		batch.SetEventFunc(func(event any) { events = append(events, event) })

		err := errors.New("boom")
		counts, recordErr := batch.RecordFailedJob(context.Background(), "job-1", err)

		if recordErr != nil {
			t.Fatal(recordErr)
		}

		if counts.PendingJobs != 0 || counts.FailedJobs != 1 {
			t.Fatalf("counts = pending:%d failed:%d, want pending:0 failed:1", counts.PendingJobs, counts.FailedJobs)
		}

		if !batch.Finished() || !batch.Cancelled() {
			t.Fatalf("batch finished/cancelled = %v/%v, want true/true", batch.Finished(), batch.Cancelled())
		}

		if catchCount != 1 || !errors.Is(receivedErr, err) || progressCount != 0 || finallyCount != 1 {
			t.Fatalf("callback counts catch:%d progress:%d finally:%d err:%v", catchCount, progressCount, finallyCount, receivedErr)
		}

		if len(events) != 2 {
			t.Fatalf("events = %d, want 2", len(events))
		}

		if _, ok := events[0].(bus.BatchStarted); !ok {
			t.Fatalf("first event = %T, want bus.BatchStarted", events[0])
		}

		if _, ok := events[1].(bus.BatchCanceled); !ok {
			t.Fatalf("second event = %T, want bus.BatchCanceled", events[1])
		}
	})

	t.Run("tolerant batch records failures without cancellation", func(t *testing.T) {
		var catchCount, progressCount, finallyCount int

		batch := &bus.Batch{
			ID:          "batch-2",
			TotalJobs:   2,
			PendingJobs: 2,
			Options:     map[string]any{"allowFailures": true},
			ProgressCallbacks: []func(context.Context, *bus.Batch){
				func(context.Context, *bus.Batch) { progressCount++ },
			},
			CatchCallbacks: []func(context.Context, *bus.Batch, error){
				func(context.Context, *bus.Batch, error) { catchCount++ },
			},
			FinallyCallbacks: []func(context.Context, *bus.Batch){
				func(context.Context, *bus.Batch) { finallyCount++ },
			},
		}

		if _, err := batch.RecordFailedJob(context.Background(), "job-1", errors.New("boom")); err != nil {
			t.Fatal(err)
		}

		if batch.Finished() || batch.Cancelled() {
			t.Fatalf("batch finished/cancelled = %v/%v, want false/false", batch.Finished(), batch.Cancelled())
		}

		if catchCount != 1 || progressCount != 1 || finallyCount != 0 {
			t.Fatalf("callback counts catch:%d progress:%d finally:%d, want 1/1/0", catchCount, progressCount, finallyCount)
		}
	})
}

// Port of Framework\Tests\Bus\BusBatchTest::test_batch_can_be_cancelled
// Port of Framework\Tests\Bus\BusBatchTest::test_batch_cancelled_event_is_dispatched
// Port of Framework\Tests\Bus\BusBatchTest::test_batch_can_be_deleted
// Port of Framework\Tests\Bus\BusBatchTest::test_batch_state_can_be_inspected
func TestUpstreamBusBatchStateAndDeletion(t *testing.T) {
	repo := newMockBatchRepo()
	batch := bus.NewBatchWithRepo("batch-1", repo)

	var event any

	batch.SetEventFunc(func(e any) { event = e })

	if err := batch.Cancel(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !repo.hasCalled("Cancel:batch-1") {
		t.Fatal("expected repository cancel")
	}

	if _, ok := event.(bus.BatchCanceled); !ok {
		t.Fatalf("event = %T, want bus.BatchCanceled", event)
	}

	if err := batch.Delete(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !repo.hasCalled("Delete:batch-1") {
		t.Fatal("expected repository delete")
	}

	local := &bus.Batch{PendingJobs: 1, Options: map[string]any{}}

	if local.Finished() || local.HasFailures() || local.Cancelled() || local.AllowsFailures() {
		t.Fatal("new batch should start unfinished, not failed, not cancelled, and not failure tolerant")
	}

	local.Options["allowFailures"] = true
	local.FailedJobs = 1
	now := time.Now()
	local.CancelledAt = &now

	if !local.HasFailures() || !local.Cancelled() || !local.Canceled() || !local.AllowsFailures() {
		t.Fatal("batch state accessors did not reflect updated state")
	}
}

// Port of Framework\Tests\Bus\BusBatchTest::test_chain_can_be_added_to_batch
// Port of Framework\Tests\Bus\BusBatchTest::test_chained_jobs_in_batch_preserve_their_queue_when_batch_has_no_queue
func TestUpstreamBusBatchAddChainJobs(t *testing.T) {
	d := newMockQueueingDispatcher()
	batch := &bus.Batch{ID: "batch-1", Options: map[string]any{}}
	batch.SetDispatcher(d)

	first := &laravelBatchJob{Name: "first"}
	second := &laravelBatchJob{Name: "second"}
	first.Chain(second)
	first.AllOnQueue("custom-queue")

	if err := batch.Add(context.Background(), []any{first, second}); err != nil {
		t.Fatal(err)
	}

	if first.BatchID != "batch-1" || second.BatchID != "batch-1" {
		t.Fatalf("chained jobs were not tagged with batch id: %q %q", first.BatchID, second.BatchID)
	}

	if first.Queue != "custom-queue" || second.Queue != "custom-queue" {
		t.Fatalf("chain queues = %q/%q, want custom-queue/custom-queue", first.Queue, second.Queue)
	}
}

// Port of Framework\Tests\Bus\BusBatchableTest::test_batch_may_be_retrieved
// Port of Framework\Tests\Bus\BusBatchableTest::test_with_fake_batch_sets_and_returns_fake
// Port of Framework\Tests\Bus\BusBatchableTest::test_batching_reflects_cancelled_state
func TestUpstreamBusBatchableHelpers(t *testing.T) {
	repo := newMockBatchRepo()
	repo.batch = &bus.Batch{ID: "batch-1", Name: "from-repo"}

	job := &laravelBatchJob{}
	job.WithBatchID("batch-1")

	batch, err := job.BatchFromRepo(context.Background(), repo)

	if err != nil {
		t.Fatal(err)
	}

	if batch.Name != "from-repo" || job.Batch() != batch {
		t.Fatalf("BatchFromRepo returned %#v and cached %#v", batch, job.Batch())
	}

	fake := job.WithFakeBatch("fake-id", "fake-name", 3)

	if fake.ID != "fake-id" || fake.Name != "fake-name" || fake.TotalJobs != 3 || job.Batch() != fake {
		t.Fatalf("fake batch = %#v", fake)
	}

	if !job.Batching() {
		t.Fatal("expected active fake batch to report batching")
	}

	if err := fake.Cancel(context.Background()); err != nil {
		t.Fatal(err)
	}

	if job.Batching() {
		t.Fatal("expected cancelled fake batch to stop batching")
	}
}

// Port of Framework\Tests\Bus\BusDispatcherTest::testCommandsThatShouldQueueIsQueued
// Port of Framework\Tests\Bus\BusDispatcherTest::testCommandsThatShouldQueueIsQueuedUsingCustomHandler
// Port of Framework\Tests\Bus\BusDispatcherTest::testCommandsThatShouldQueueIsQueuedUsingCustomQueueAndDelay
// Port of Framework\Tests\Bus\BusDispatcherTest::testCommandsAreDispatchedWithQueueRoute
// Port of Framework\Tests\Bus\BusDispatcherTest::testDispatchNowShouldNeverQueue
// Port of Framework\Tests\Bus\BusDispatcherTest::testDispatcherCanDispatchStandAloneHandler
// Port of Framework\Tests\Bus\BusDispatcherTest::testOnConnectionOnJobWhenDispatching
func TestUpstreamBusDispatcherQueueRoutingAndSync(t *testing.T) {
	q := newMockQueue()
	d := bus.NewDispatcher(q, nil)

	queued := &delayedQueueableCommand{Value: "queued"}
	queued.OnQueue("emails").WithDelay(10 * time.Second)

	if _, err := d.Dispatch(context.Background(), queued); err != nil {
		t.Fatal(err)
	}

	q.mu.Lock()
	pushes := append([]mockPush(nil), q.pushes...)
	delayed := append([]time.Duration(nil), q.delays...)
	q.mu.Unlock()

	if len(pushes) != 1 || pushes[0].Queue != "emails" || delayed[0] != 10*time.Second {
		t.Fatalf("queue pushes = %#v delays = %#v, want one delayed email push", pushes, delayed)
	}

	d.Map(testCommand{}, func(_ context.Context, cmd any) (any, error) {
		return cmd.(testCommand).Value, nil
	})

	result, err := d.DispatchNow(context.Background(), testCommand{Value: "sync"})

	if err != nil {
		t.Fatal(err)
	}

	if result != "sync" {
		t.Fatalf("DispatchNow result = %v, want sync", result)
	}

	q.mu.Lock()
	count := len(q.pushes)
	q.mu.Unlock()

	if count != 1 {
		t.Fatalf("DispatchNow pushed to queue; pushes = %d, want 1", count)
	}
}

// Port of Framework\Tests\Bus\BusPendingBatchTest::test_pending_batch_may_be_configured_and_dispatched
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_batch_is_deleted_from_storage_if_exception_thrown_during_batching
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_batch_is_dispatched_when_dispatchif_is_true
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_batch_is_not_dispatched_when_dispatchif_is_false
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_batch_is_dispatched_when_dispatchunless_is_false
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_batch_is_not_dispatched_when_dispatchunless_is_true
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_batch_before_event_is_called
func TestUpstreamBusPendingBatchDispatchControls(t *testing.T) {
	d := newMockQueueingDispatcher()

	var beforeCalled bool

	pending := bus.NewPendingBatch(d, []any{&laravelBatchJob{Name: "one"}}).
		Before(func(context.Context, *bus.Batch) { beforeCalled = true }).
		Progress(func(context.Context, *bus.Batch) {}).
		Then(func(context.Context, *bus.Batch) {}).
		Catch(func(context.Context, *bus.Batch, error) {}).
		AllowFailures().
		OnConnection("redis").
		OnQueue("high").
		WithOption("extra-option", 123)

	batch, err := pending.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !beforeCalled || batch.Options["extra-option"] != 123 || !batch.AllowsFailures() {
		t.Fatalf("pending batch configuration was not preserved: before=%v options=%#v", beforeCalled, batch.Options)
	}

	skipped, err := pending.DispatchIf(context.Background(), false)

	if err != nil {
		t.Fatal(err)
	}

	if skipped != nil {
		t.Fatal("DispatchIf(false) should skip")
	}

	dispatched, err := pending.DispatchUnless(context.Background(), false)

	if err != nil {
		t.Fatal(err)
	}

	if dispatched == nil {
		t.Fatal("DispatchUnless(false) should dispatch")
	}
}

// Port of Framework\Tests\Bus\BusPendingBatchTest::test_allow_failures_with_boolean_true_enables_failure_tolerance
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_allow_failures_with_boolean_false_disables_failure_tolerance
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_allow_failures_with_single_closure_registers_callback
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_allow_failures_with_array_of_callables_registers_multiple_callbacks
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_allow_failures_with_empty_array_enables_tolerance_without_callbacks
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_allow_failures_is_chainable
// Port of Framework\Tests\Bus\BusPendingBatchTest::test_failure_callbacks_accessor_returns_registered_callbacks
func TestUpstreamBusPendingBatchAllowFailuresTypedAPI(t *testing.T) {
	d := newMockQueueingDispatcher()
	pending := bus.NewPendingBatch(d, []any{&laravelBatchJob{Name: "one"}})

	if pending.AllowFailures() != pending {
		t.Fatal("AllowFailures should be chainable")
	}

	if !pending.AllowsFailures() {
		t.Fatal("AllowFailures should enable tolerance")
	}

	pending.DisallowFailures()

	if pending.AllowsFailures() {
		t.Fatal("DisallowFailures should disable tolerance")
	}

	callbackOne := func(context.Context, *bus.Batch, error) {}
	callbackTwo := func(context.Context, *bus.Batch, error) {}

	if pending.OnFailure(callbackOne, callbackTwo) != pending {
		t.Fatal("OnFailure should be chainable")
	}

	if !pending.AllowsFailures() {
		t.Fatal("OnFailure should enable tolerance")
	}

	if got := len(pending.FailureCallbacks()); got != 2 {
		t.Fatalf("FailureCallbacks length = %d, want 2", got)
	}

	batch, err := pending.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got := len(batch.CatchCallbacks); got != 2 {
		t.Fatalf("batch catch callbacks = %d, want 2", got)
	}
}

// Port of Framework\Tests\Bus\BusPendingDispatchTest::testOnConnection
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testOnQueue
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testAllOnConnection
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testAllOnQueue
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testDelay
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testWithoutDelay
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testAfterCommit
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testBeforeCommit
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testChain
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testAfterResponse
// Port of Framework\Tests\Bus\BusPendingDispatchTest::testGetJob
// Port of Framework\Tests\Bus\QueueableTest::testOnConnection
// Port of Framework\Tests\Bus\QueueableTest::testAllOnConnection
// Port of Framework\Tests\Bus\QueueableTest::testOnQueue
// Port of Framework\Tests\Bus\QueueableTest::testAllOnQueue
func TestUpstreamBusQueueableAndPendingDispatchEquivalentFluentAPI(t *testing.T) {
	first := &laravelBatchJob{Name: "first"}
	second := &laravelBatchJob{Name: "second"}

	first.Chain(second).
		OnConnection("redis").
		OnQueue("emails").
		WithDelay(60 * time.Second).
		SetAfterCommit()

	if first.Connection != "redis" || first.Queue != "emails" || first.GetDelay() != 60*time.Second {
		t.Fatalf("queueable fields = connection:%q queue:%q delay:%s", first.Connection, first.Queue, first.GetDelay())
	}

	if first.AfterCommit == nil || !*first.AfterCommit {
		t.Fatal("SetAfterCommit did not set true")
	}

	first.WithoutDelay().SetBeforeCommit().AllOnConnection("sqs").AllOnQueue("critical")

	if first.GetDelay() != 0 || first.AfterCommit == nil || *first.AfterCommit {
		t.Fatalf("delay/commit fields = %s/%v", first.GetDelay(), first.AfterCommit)
	}

	if first.Connection != "sqs" || first.ChainConnection != "sqs" || second.Connection != "sqs" {
		t.Fatalf("chain connection propagation failed: first=%q chain=%q second=%q", first.Connection, first.ChainConnection, second.Connection)
	}

	if first.Queue != "critical" || first.ChainQueue != "critical" || second.Queue != "critical" {
		t.Fatalf("chain queue propagation failed: first=%q chain=%q second=%q", first.Queue, first.ChainQueue, second.Queue)
	}

	d := bus.NewDispatcher(nil, nil)
	d.Map(first, func(context.Context, any) (any, error) { return "ok", nil })
	d.WithDispatchingAfterResponses()

	if _, err := d.Dispatch(context.Background(), first); err != nil {
		t.Fatal(err)
	}

	if err := d.FlushDeferred(context.Background()); err != nil {
		t.Fatal(err)
	}
}
