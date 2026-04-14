package drivers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// Partial port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest.
//
// This file ports the tests that exercise the Push / PushDelayed path
// (plus the two shared payload-failure tests that historically live in
// this PHP file but actually cover the abstract Queue::createPayload
// contract). The remaining Unit tests require driver additions that
// this step does not land yet:
//
//   - Bulk/Insert-many (testBulkBatchPushesOntoDatabase)
//   - pendingJobs / delayedJobs / reservedJobs inspection (3 tests)
//   - buildDatabaseRecord field ordering (1 test)
//   - batch-id injection for Batchable jobs (1 test, Step 12)
//   - getLockForPoppingIsCached (internal mechanism, Step 14-ish)
//
// Those tests are left in the inventory as pending work for the next
// database iteration.
//
// Adaptation rules applied (see PARITY.md §2):
//
//   - Upstream's DatabaseQueue uses a fluent query builder
//     ($db->table->insertGetId([...])). The Go driver speaks raw SQL
//     via DBExecer.Exec(query, args...). The port asserts the observable
//     side effects — one INSERT was executed against the right table
//     with args that encode the expected row — rather than matching the
//     fluent-builder call chain literally.
//   - Upstream freezes time with Carbon::setTestNow and locks the UUID
//     with Str::createUuidsUsing. Go has neither facade; the port
//     captures time.Now() once in the test, lets the driver use its own
//     clock, and asserts available_at / created_at land in a tight
//     sanity window rather than matching an exact value.

// extractInsert returns the first INSERT call recorded by the mock,
// its queue argument, its payload JSON, and the integer args in
// positional order. Fails the test immediately on mismatch.

// assertPayloadFields checks the JSON-encoded payload that landed in
// the INSERT carries the Upstream-shape uuid/displayName/job fields.

// DatabaseUpstreamTestJob is a Go analogue of the PHP `MyTestJob`
// fixture used by the Upstream data provider. The suffix of its
// reflect-derived display name (`.DatabaseUpstreamTestJob`) is what
// the port asserts on.
type DatabaseUpstreamTestJob struct{}

func extractInsert(t *testing.T, db *mockDBExecer) (queueArg, payloadJSON string, availableAt, createdAt int64) {
	t.Helper()

	if len(db.execCalls) != 1 {
		t.Fatalf("exec calls: got %d, want 1", len(db.execCalls))
	}

	call := db.execCalls[0]

	if !strings.Contains(call.Query, "INSERT INTO jobs") {
		t.Fatalf("exec query: got %q, want INSERT INTO jobs ...", call.Query)
	}

	if len(call.Args) != 4 {
		t.Fatalf("exec args: got %d, want 4 (queue, payload, available_at, created_at)", len(call.Args))
	}

	var ok bool

	if queueArg, ok = call.Args[0].(string); !ok {
		t.Fatalf("args[0] queue: got %T, want string", call.Args[0])
	}

	if payloadJSON, ok = call.Args[1].(string); !ok {
		t.Fatalf("args[1] payload: got %T, want string", call.Args[1])
	}

	if availableAt, ok = call.Args[2].(int64); !ok {
		t.Fatalf("args[2] available_at: got %T (%v), want int64", call.Args[2], call.Args[2])
	}

	if createdAt, ok = call.Args[3].(int64); !ok {
		t.Fatalf("args[3] created_at: got %T (%v), want int64", call.Args[3], call.Args[3])
	}

	return queueArg, payloadJSON, availableAt, createdAt
}

func assertPayloadFields(t *testing.T, payloadJSON, wantDisplayContains string) {
	t.Helper()

	var decoded map[string]any

	if err := json.Unmarshal([]byte(payloadJSON), &decoded); err != nil {
		t.Fatalf("payload JSON: %v", err)
	}

	if uuid, _ := decoded["uuid"].(string); uuid == "" {
		t.Errorf("payload.uuid: empty")
	}

	displayName, _ := decoded["displayName"].(string)

	if !strings.Contains(displayName, wantDisplayContains) {
		t.Errorf("payload.displayName: got %q, want substring %q", displayName, wantDisplayContains)
	}

	jobName, _ := decoded["job"].(string)

	if !strings.Contains(jobName, wantDisplayContains) {
		t.Errorf("payload.job: got %q, want substring %q", jobName, wantDisplayContains)
	}
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testPushProperlyPushesJobOntoDatabase
//
// Upstream's test uses a PHP data provider to drive three variants
// (object / closure / string). The Go port runs two subtests —
// closures don't have a meaningful display name in Go, so that
// variant is dropped per the adaptation rule.
func TestPushProperlyPushesJobOntoDatabase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		job         any
		wantContain string
	}{
		{"struct", DatabaseUpstreamTestJob{}, ".DatabaseUpstreamTestJob"},
		{"string", "foo", "foo"},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := newMockDBExecer()
			drv := drivers.NewDatabaseDriver(db, "jobs", "database")

			_, payload, err := queue.CreatePayloadFor("database", "default", tc.job, map[string]any{"0": "data"}, queue.JobOptions{})

			if err != nil {
				t.Fatalf("CreatePayloadFor: %v", err)
			}

			before := time.Now().Unix()

			if _, err := drv.Push(context.Background(), "default", payload); err != nil {
				t.Fatalf("Push: %v", err)
			}

			after := time.Now().Unix()

			queueArg, payloadJSON, availableAt, createdAt := extractInsert(t, db)

			if queueArg != "default" {
				t.Errorf("queue arg: got %q, want default", queueArg)
			}

			assertPayloadFields(t, payloadJSON, tc.wantContain)

			if availableAt < before || availableAt > after {
				t.Errorf("available_at: got %d, want in [%d, %d]", availableAt, before, after)
			}

			if createdAt < before || createdAt > after {
				t.Errorf("created_at: got %d, want in [%d, %d]", createdAt, before, after)
			}
		})
	}
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testDelayedPushProperlyPushesJobOntoDatabase
func TestDelayedPushProperlyPushesJobOntoDatabase(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	_, payload, err := queue.CreatePayloadFor("database", "default", "foo", map[string]any{"0": "data"}, queue.JobOptions{})

	if err != nil {
		t.Fatalf("CreatePayloadFor: %v", err)
	}

	before := time.Now().Unix()

	if _, err := drv.PushDelayed(context.Background(), "default", payload, 10*time.Second); err != nil {
		t.Fatalf("PushDelayed: %v", err)
	}

	after := time.Now().Unix()

	queueArg, _, availableAt, createdAt := extractInsert(t, db)

	if queueArg != "default" {
		t.Errorf("queue arg: got %q, want default", queueArg)
	}

	// Upstream asserts availableAt is an int and reflects the 10-second
	// delay. The Go port checks the delta lands in a [9, 11] window to
	// cover the tiny time.Now() drift between capture points.
	if createdAt < before || createdAt > after {
		t.Errorf("created_at: got %d, want in [%d, %d]", createdAt, before, after)
	}

	if delta := availableAt - createdAt; delta < 9 || delta > 11 {
		t.Errorf("available_at - created_at: got %d, want ~10s", delta)
	}
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testFailureToCreatePayloadFromArray
//
// Upstream feeds createPayload a value containing invalid UTF-8 and
// expects an InvalidArgumentException. Go's encoding/json does not
// error on invalid UTF-8 in string values (it escapes with U+FFFD),
// so the Go port uses a value json.Marshal genuinely cannot encode
// — a channel — and asserts *InvalidPayloadError is returned.
func TestFailureToCreatePayloadFromArray(t *testing.T) {
	t.Parallel()

	_, _, err := queue.CreatePayloadFor("database", "default", "foo", map[string]any{"ch": make(chan int)}, queue.JobOptions{})

	if err == nil {
		t.Fatal("expected error from unmarshalable payload")
	}

	var invalid *queue.InvalidPayloadError

	if !errors.As(err, &invalid) {
		t.Errorf("expected *queue.InvalidPayloadError, got %T (%v)", err, err)
	}
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testFailureToCreatePayloadFromObject
//
// Same adaptation as testFailureToCreatePayloadFromArray — Go's JSON
// encoder rejects chan/func/complex values, so the port uses one of
// those to trigger the InvalidPayloadError path.
func TestFailureToCreatePayloadFromObject(t *testing.T) {
	t.Parallel()

	type invalidJob struct{ Self any }

	// A self-referential map → json.Marshal hits the recursion cap
	// and returns an UnsupportedValueError.
	cyclic := map[string]any{}
	cyclic["self"] = cyclic

	_, _, err := queue.CreatePayloadFor("database", "default", invalidJob{}, cyclic, queue.JobOptions{})

	if err == nil {
		t.Fatal("expected error from unmarshalable payload")
	}

	var invalid *queue.InvalidPayloadError

	if !errors.As(err, &invalid) {
		t.Errorf("expected *queue.InvalidPayloadError, got %T (%v)", err, err)
	}

	// Guard: fmt.Sprintf is just here so the helper import stays used
	// when the test is trimmed by a future refactor.
	_ = fmt.Sprintf("%v", invalid)
}
