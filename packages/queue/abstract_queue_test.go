package queue_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// These tests are not 1:1 ports from tests/Queue — they cover Go-specific
// helpers (DisplayName, NewUUIDv4, CreatePayloadFor, ShouldDispatchAfterCommit)
// that the abstract_queue.go file introduces. The payload-hook behaviour
// they assert matches upstream Queue::createPayloadUsing's contract.

type sampleJob struct {
	_       struct{} `queue:"tries=4,timeout=45s,backoff=1s|2s,queue=mail"`
	Subject string
}

type namedJob struct{}

type afterCommitJob struct{}

type beforeCommitJob struct{}

func (namedJob) QueueDisplayName() string { return "my.custom.name" }

func (afterCommitJob) QueueAfterCommit() bool { return true }

func (beforeCommitJob) QueueBeforeCommit() bool { return true }

func TestDisplayNameFallsBackToReflectType(t *testing.T) {
	t.Parallel()

	got := queue.DisplayName(sampleJob{})

	if !strings.HasSuffix(got, ".sampleJob") {
		t.Errorf("DisplayName: got %q, want suffix .sampleJob", got)
	}
}

func TestDisplayNameRespectsNamer(t *testing.T) {
	t.Parallel()

	if got := queue.DisplayName(namedJob{}); got != "my.custom.name" {
		t.Errorf("DisplayName: got %q, want my.custom.name", got)
	}
}

func TestDisplayNamePassesStringThrough(t *testing.T) {
	t.Parallel()

	if got := queue.DisplayName("App\\Jobs\\Literal"); got != "App\\Jobs\\Literal" {
		t.Errorf("DisplayName: got %q", got)
	}
}

func TestDisplayNameNil(t *testing.T) {
	t.Parallel()

	if got := queue.DisplayName(nil); got != "" {
		t.Errorf("DisplayName: got %q, want empty", got)
	}
}

func TestNewUUIDv4Shape(t *testing.T) {
	t.Parallel()

	u := queue.NewUUIDv4()

	if len(u) != 36 {
		t.Fatalf("NewUUIDv4: length %d, want 36", len(u))
	}

	parts := strings.Split(u, "-")

	if len(parts) != 5 || len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		t.Errorf("NewUUIDv4: bad shape %q", u)
	}

	// Version nibble should be 4; variant nibble should be 8, 9, a, or b.
	if parts[2][0] != '4' {
		t.Errorf("NewUUIDv4: version nibble %c, want 4", parts[2][0])
	}

	switch parts[3][0] {
	case '8', '9', 'a', 'b':
	default:
		t.Errorf("NewUUIDv4: variant nibble %c, want 8/9/a/b", parts[3][0])
	}
}

func TestNewUUIDv4IsUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[string]struct{}, 100)

	for i := 0; i < 100; i++ {
		u := queue.NewUUIDv4()

		if _, dup := seen[u]; dup {
			t.Fatalf("NewUUIDv4: duplicate %q after %d calls", u, i)
		}

		seen[u] = struct{}{}
	}
}

func TestCreatePayloadForWritesFrameworkShape(t *testing.T) {
	// Not t.Parallel: mutates the global payload-hook list.
	queue.ClearPayloadHooks()

	defer queue.ClearPayloadHooks()

	opts, err := queue.ParseJobOptions(sampleJob{})

	if err != nil {
		t.Fatalf("ParseJobOptions: %v", err)
	}

	p, raw, err := queue.CreatePayloadFor("redis", "mail", sampleJob{Subject: "hi"}, map[string]any{"subject": "hi"}, opts)

	if err != nil {
		t.Fatalf("CreatePayloadFor: %v", err)
	}

	if p.UUID == "" {
		t.Error("UUID: empty")
	}

	if p.MaxTries != 4 {
		t.Errorf("MaxTries: got %d, want 4", p.MaxTries)
	}

	if p.Timeout != 45 {
		t.Errorf("Timeout: got %d, want 45", p.Timeout)
	}

	wantBackoff := []int{1, 2}

	if len(p.Backoff) != len(wantBackoff) || p.Backoff[0] != 1 || p.Backoff[1] != 2 {
		t.Errorf("Backoff: got %v, want %v", p.Backoff, wantBackoff)
	}

	if !strings.HasSuffix(p.DisplayName, ".sampleJob") {
		t.Errorf("DisplayName: got %q", p.DisplayName)
	}

	if p.CreatedAt == nil {
		t.Error("CreatedAt: nil")
	}

	var decoded map[string]any

	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("raw JSON unmarshal: %v", err)
	}

	if decoded["uuid"] != p.UUID {
		t.Errorf("uuid mismatch between struct and raw JSON")
	}
}

func TestCreatePayloadForRunsHooksInOrder(t *testing.T) {
	// Not t.Parallel: mutates the global payload-hook list.
	queue.ClearPayloadHooks()

	defer queue.ClearPayloadHooks()

	queue.CreatePayloadUsing(func(conn, queueName string, p *queue.Payload) {
		if p.Data == nil {
			p.Data = map[string]any{}
		}

		p.Data["tenant"] = "acme"
	})

	queue.CreatePayloadUsing(func(conn, queueName string, p *queue.Payload) {
		p.Data["hook_order"] = []string{"tenant", "audit"}
	})

	p, _, err := queue.CreatePayloadFor("sync", "default", "App\\Jobs\\X", nil, queue.JobOptions{})

	if err != nil {
		t.Fatalf("CreatePayloadFor: %v", err)
	}

	if p.Data["tenant"] != "acme" {
		t.Errorf("first hook did not run: %+v", p.Data)
	}

	order, ok := p.Data["hook_order"].([]string)

	if !ok || len(order) != 2 || order[0] != "tenant" || order[1] != "audit" {
		t.Errorf("hook ordering: got %v", p.Data["hook_order"])
	}
}

func TestCreatePayloadForRetryUntil(t *testing.T) {
	// Not t.Parallel: mutates the global payload-hook list.
	queue.ClearPayloadHooks()

	defer queue.ClearPayloadHooks()

	deadline := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)

	p, _, err := queue.CreatePayloadFor("sync", "default", "x", nil, queue.JobOptions{RetryUntil: deadline})

	if err != nil {
		t.Fatalf("CreatePayloadFor: %v", err)
	}

	if p.RetryUntil == nil || !p.RetryUntil.Equal(deadline) {
		t.Errorf("RetryUntil: got %v, want %v", p.RetryUntil, deadline)
	}
}

func TestShouldDispatchAfterCommitPrecedence(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		job    any
		config map[string]any
		want   bool
	}{
		{
			name: "BeforeCommitMarker forces false even when config says true",
			job:  beforeCommitJob{},
			config: map[string]any{
				"after_commit": true,
			},
			want: false,
		},
		{
			name:   "AfterCommitMarker forces true",
			job:    afterCommitJob{},
			config: nil,
			want:   true,
		},
		{
			name:   "Config after_commit=true without marker",
			job:    sampleJob{},
			config: map[string]any{"after_commit": true},
			want:   true,
		},
		{
			name:   "Default is false",
			job:    sampleJob{},
			config: nil,
			want:   false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := queue.ShouldDispatchAfterCommit(tc.job, tc.config); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// Port of Illuminate\Tests\Queue\BeforeCommitContractTest::testJobWithContractDefaultsToAfterCommit
// Port of Illuminate\Tests\Queue\BeforeCommitContractTest::testJobWithContractAndAfterCommitFalseRespectsBeforeCommit
// Port of Illuminate\Tests\Queue\BeforeCommitContractTest::testJobWithContractAndExplicitAfterCommitTrueStillSchedulesAfterCommit
// Port of Illuminate\Tests\Queue\BeforeCommitContractTest::testJobWithoutContractRespectsAfterCommit
// Port of Illuminate\Tests\Queue\BeforeCommitContractTest::testJobWithoutContractRespectsBeforeCommit
func TestBeforeCommitContractDispatchPrecedence(t *testing.T) {
	t.Parallel()

	if !queue.ShouldDispatchAfterCommit(afterCommitJob{}, nil) {
		t.Fatal("after-commit marker should dispatch after commit by default")
	}

	if queue.ShouldDispatchAfterCommit(beforeCommitJob{}, map[string]any{"after_commit": true}) {
		t.Fatal("before-commit marker should override after_commit=true config")
	}

	if !queue.ShouldDispatchAfterCommit(afterCommitJob{}, map[string]any{"after_commit": false}) {
		t.Fatal("after-commit marker should override after_commit=false config")
	}

	if !queue.ShouldDispatchAfterCommit(sampleJob{}, map[string]any{"after_commit": true}) {
		t.Fatal("plain job should respect after_commit=true config")
	}

	if queue.ShouldDispatchAfterCommit(sampleJob{}, map[string]any{"after_commit": false}) {
		t.Fatal("plain job should respect after_commit=false config")
	}
}

// Port of Illuminate\Tests\Queue\QueueDelayTest::test_queue_delay
// Port of Illuminate\Tests\Queue\QueueDelayTest::test_queue_without_delay
// Port of Illuminate\Tests\Queue\QueueDelayTest::test_pending_dispatch_without_delay
func TestQueueDelayAndWithoutDelayOptions(t *testing.T) {
	t.Parallel()

	opts := queue.JobOptions{Delay: 60 * time.Second}

	if opts.Delay != 60*time.Second {
		t.Fatalf("Delay: got %s, want 60s", opts.Delay)
	}

	withoutDelay := opts.WithoutDelay()

	if withoutDelay.Delay != 0 {
		t.Errorf("WithoutDelay: got %s, want 0", withoutDelay.Delay)
	}

	if opts.Delay != 60*time.Second {
		t.Errorf("original Delay changed: got %s, want 60s", opts.Delay)
	}
}

// Port of Illuminate\Tests\Queue\QueueDatabaseQueueUnitTest::testPushIncludesBatchIdInPayloadForBatchableJob
func TestCreatePayloadForIncludesBatchID(t *testing.T) {
	t.Parallel()

	p, raw, err := queue.CreatePayloadFor("database", "default", sampleJob{}, nil, queue.JobOptions{BatchID: "batch-123"})

	if err != nil {
		t.Fatalf("CreatePayloadFor: %v", err)
	}

	if p.Data["batchId"] != "batch-123" {
		t.Fatalf("payload data: got %+v, want batchId", p.Data)
	}

	var decoded struct {
		Data map[string]any `json:"data"`
	}

	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("raw JSON: %v", err)
	}

	if decoded.Data["batchId"] != "batch-123" {
		t.Errorf("raw batchId: got %v, want batch-123", decoded.Data["batchId"])
	}
}
