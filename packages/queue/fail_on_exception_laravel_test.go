package queue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock/packages/queue"
)

var errFailOnException = errors.New("fail this job")

// Port of Illuminate\Tests\Queue\FailOnExceptionMiddlewareTest::test_middleware
func TestFailOnExceptionMiddleware(t *testing.T) {
	t.Parallel()

	job := &workerLaravelJob{}
	middleware := queue.NewFailOnException(errFailOnException)

	err := middleware.Handle(context.Background(), job, func(context.Context, queue.Job) error {
		return errFailOnException
	})

	if !errors.Is(err, errFailOnException) {
		t.Fatalf("Handle error: got %v, want %v", err, errFailOnException)
	}

	if !job.HasFailed() {
		t.Fatal("job should be failed when the middleware matches")
	}
}

// Port of Illuminate\Tests\Queue\FailOnExceptionMiddlewareTest::test_can_test_against_job_properties
func TestFailOnExceptionMiddlewareCanInspectJobProperties(t *testing.T) {
	t.Parallel()

	job := &workerLaravelJob{queueName: "critical"}
	middleware := queue.NewFailOnException().When(func(err error, job queue.Job) bool {
		return errors.Is(err, errFailOnException) && job.GetQueue() == "critical"
	})

	err := middleware.Handle(context.Background(), job, func(context.Context, queue.Job) error {
		return errFailOnException
	})

	if !errors.Is(err, errFailOnException) {
		t.Fatalf("Handle error: got %v, want %v", err, errFailOnException)
	}

	if !job.HasFailed() {
		t.Fatal("job should be failed when the predicate matches job properties")
	}
}
