package queue_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// The struct-tag parser replaces Laravel's PHP8 attribute system
// (#[Tries], #[Backoff], etc.). These tests are not a direct port of a
// single PHPUnit file — they validate the Go-specific options mechanism
// the rest of the queue package will rely on in later steps.

type fullyTaggedJob struct {
	_       struct{} `queue:"tries=5,max_exceptions=2,timeout=30s,delay=1s,backoff=1s|5s|10s,unique_for=5m,fail_on_timeout,queue=emails,connection=redis,retry_until=2026-05-01T12:00:00Z"`
	Payload string
}

type flagOnlyJob struct {
	_ struct{} `queue:"fail_on_timeout"`
}

type emptyJob struct {
	Payload string
}

func TestParseJobOptionsReadsEveryKey(t *testing.T) {
	t.Parallel()

	opts, err := queue.ParseJobOptions(fullyTaggedJob{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.MaxTries != 5 {
		t.Errorf("MaxTries: got %d, want 5", opts.MaxTries)
	}

	if opts.MaxExceptions != 2 {
		t.Errorf("MaxExceptions: got %d, want 2", opts.MaxExceptions)
	}

	if opts.Timeout != 30*time.Second {
		t.Errorf("Timeout: got %s, want 30s", opts.Timeout)
	}

	if opts.Delay != time.Second {
		t.Errorf("Delay: got %s, want 1s", opts.Delay)
	}

	wantBackoff := []time.Duration{time.Second, 5 * time.Second, 10 * time.Second}

	if len(opts.Backoff) != len(wantBackoff) {
		t.Fatalf("Backoff len: got %d, want %d", len(opts.Backoff), len(wantBackoff))
	}

	for i := range wantBackoff {
		if opts.Backoff[i] != wantBackoff[i] {
			t.Errorf("Backoff[%d]: got %s, want %s", i, opts.Backoff[i], wantBackoff[i])
		}
	}

	if opts.UniqueFor != 5*time.Minute {
		t.Errorf("UniqueFor: got %s, want 5m", opts.UniqueFor)
	}

	if !opts.FailOnTimeout {
		t.Error("FailOnTimeout: expected true")
	}

	if opts.Queue != "emails" {
		t.Errorf("Queue: got %q, want %q", opts.Queue, "emails")
	}

	if opts.Connection != "redis" {
		t.Errorf("Connection: got %q, want %q", opts.Connection, "redis")
	}

	wantRetry, _ := time.Parse(time.RFC3339, "2026-05-01T12:00:00Z")

	if !opts.RetryUntil.Equal(wantRetry) {
		t.Errorf("RetryUntil: got %s, want %s", opts.RetryUntil, wantRetry)
	}
}

func TestParseJobOptionsAcceptsBareBooleanFlag(t *testing.T) {
	t.Parallel()

	opts, err := queue.ParseJobOptions(flagOnlyJob{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.FailOnTimeout {
		t.Error("expected FailOnTimeout=true for bare flag")
	}
}

func TestParseJobOptionsPointerAndStruct(t *testing.T) {
	t.Parallel()

	want, err := queue.ParseJobOptions(fullyTaggedJob{})

	if err != nil {
		t.Fatalf("struct form: %v", err)
	}

	got, err := queue.ParseJobOptions(&fullyTaggedJob{})

	if err != nil {
		t.Fatalf("pointer form: %v", err)
	}

	if got.MaxTries != want.MaxTries || got.Queue != want.Queue {
		t.Errorf("pointer and struct disagreed: got %+v, want %+v", got, want)
	}
}

func TestParseJobOptionsZeroForUntaggedStruct(t *testing.T) {
	t.Parallel()

	opts, err := queue.ParseJobOptions(emptyJob{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.MaxTries != 0 || opts.MaxExceptions != 0 || opts.Timeout != 0 ||
		opts.Delay != 0 || len(opts.Backoff) != 0 || opts.UniqueFor != 0 ||
		opts.FailOnTimeout || opts.DeleteWhenMissingModels ||
		opts.Queue != "" || opts.Connection != "" || !opts.RetryUntil.IsZero() {
		t.Errorf("expected zero JobOptions, got %+v", opts)
	}
}

func TestParseJobOptionsRejectsNonStruct(t *testing.T) {
	t.Parallel()

	if _, err := queue.ParseJobOptions(42); err == nil {
		t.Fatal("expected error for non-struct input")
	}
}

func TestParseJobOptionsRejectsUnknownKey(t *testing.T) {
	t.Parallel()

	type badJob struct {
		_ struct{} `queue:"nonsense=1"`
	}

	if _, err := queue.ParseJobOptions(badJob{}); err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestParseJobOptionsRejectsMalformedDuration(t *testing.T) {
	t.Parallel()

	type badJob struct {
		_ struct{} `queue:"timeout=abc"`
	}

	if _, err := queue.ParseJobOptions(badJob{}); err == nil {
		t.Fatal("expected error for malformed duration")
	}
}
