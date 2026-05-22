package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// recordingJob is a minimal queue.Job implementation that captures the
// argument passed to Fail so the test can inspect it. Every other
// method is a no-op — this fake only exists to exercise Fail.
type recordingJob struct {
	failErr error
}

func (j *recordingJob) UUID() string                 { return "" }
func (j *recordingJob) GetJobID() string             { return "" }
func (j *recordingJob) Payload() []byte              { return nil }
func (j *recordingJob) Fire(context.Context) error   { return nil }
func (j *recordingJob) Release(time.Duration) error  { return nil }
func (j *recordingJob) Delete() error                { return nil }
func (j *recordingJob) Fail(err error) error         { j.failErr = err; return nil }
func (j *recordingJob) MarkAsFailed(err error) error { return j.Fail(err) }
func (j *recordingJob) Attempts() int                { return 0 }
func (j *recordingJob) MaxTries() int                { return 0 }
func (j *recordingJob) MaxExceptions() int           { return 0 }
func (j *recordingJob) Timeout() time.Duration       { return 0 }
func (j *recordingJob) Backoff() []time.Duration     { return nil }
func (j *recordingJob) RetryUntil() *time.Time       { return nil }
func (j *recordingJob) IsDeleted() bool              { return false }
func (j *recordingJob) IsReleased() bool             { return false }
func (j *recordingJob) HasFailed() bool              { return j.failErr != nil }
func (j *recordingJob) GetQueue() string             { return "" }
func (j *recordingJob) GetConnectionName() string    { return "" }

// Ref: @bedrock/code-0364
func TestCreatesAnExceptionFromString(t *testing.T) {
	t.Parallel()

	rec := &recordingJob{}

	// constructed in the PHP test. Embedding is our stand-in for trait use.
	handler := struct {
		queue.InteractsWithQueue
	}{}
	handler.Job = rec

	if err := handler.Fail("Whoops!"); err != nil {
		t.Fatalf("Fail: %v", err)
	}

	if rec.failErr == nil {
		t.Fatal("expected Job.Fail to be called with a non-nil error")
	}

	if rec.failErr.Error() != "Whoops!" {
		t.Errorf("error message: got %q, want %q", rec.failErr.Error(), "Whoops!")
	}
}
