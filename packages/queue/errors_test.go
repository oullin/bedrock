package queue_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/queue"
)

// fakeRedisJob is the Go equivalent of MyFakeRedisJob in
// Framework\Tests\Queue\QueueExceptionTest. It satisfies queue.ResolveNamer
// and returns the same display name the PHP fixture returns.
type fakeRedisJob struct{}

func (fakeRedisJob) ResolveName() string { return "App\\Jobs\\UnderlyingJob" }

// Port of Framework\Tests\Queue\QueueExceptionTest::test_it_can_create_timeout_exception_for_job
func TestItCanCreateTimeoutExceptionForJob(t *testing.T) {
	t.Parallel()

	job := fakeRedisJob{}
	e := queue.NewTimeoutExceededErrorForJob(job)

	if got, want := e.Error(), "App\\Jobs\\UnderlyingJob has timed out."; got != want {
		t.Errorf("message: got %q, want %q", got, want)
	}

	if e.Job != job {
		t.Errorf("Job: got %v, want %v", e.Job, job)
	}

	// Upstream's TimeoutExceededException extends MaxAttemptsExceededException;
	// the Go port preserves this by embedding the parent. errors.As must
	// therefore resolve to the parent type as well.
	var parent *queue.MaxAttemptsExceededError

	if !errors.As(e, &parent) {
		t.Fatal("expected TimeoutExceededError to satisfy errors.As(*MaxAttemptsExceededError)")
	}
}

// Port of Framework\Tests\Queue\QueueExceptionTest::test_it_can_create_max_attempts_exception_for_job
func TestItCanCreateMaxAttemptsExceptionForJob(t *testing.T) {
	t.Parallel()

	job := fakeRedisJob{}
	e := queue.NewMaxAttemptsExceededErrorForJob(job)

	if got, want := e.Error(), "App\\Jobs\\UnderlyingJob has been attempted too many times."; got != want {
		t.Errorf("message: got %q, want %q", got, want)
	}

	if e.Job != job {
		t.Errorf("Job: got %v, want %v", e.Job, job)
	}
}
