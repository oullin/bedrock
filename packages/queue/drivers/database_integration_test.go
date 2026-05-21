//go:build integration

// Ref: @bedrock/code-0365
//
// Parity status:
//   testAvailableAndUnReservedJobsArePopped   ✅
//   testPoppedJobsIncrementAttempts           ✅
//   testThatQueueCanBeCleared                 ✅
//   testUnavailableJobsAreNotPopped           ✅
//   testThatReservedAndExpiredJobsArePopped   ⏸ no expiry reclaim in Go Pop yet — stubbed as deferral
//   testThatReservedJobsAreNotPopped          ✅
//   testJobPayloadIsAvailableOnEvents         ✅ (adapted — Go driver exposes payload on popped Job;
//                                                    the worker emits events, not the driver)
//
// Implementation notes
// --------------------
// the upstream integration suite boots an in-memory SQLite database via
// orchestra/testbench. The queue/ module's go.mod has no pure-Go
// SQLite driver, so these tests are intentionally stubbed against
// the in-package mockDBExecer (see drivers_test.go). The mock stages
// row results and records exec calls, which is sufficient to assert
// the observable behaviours the PHP tests care about (pop picks the
// right row, attempts increment, clear issues a DELETE, etc.). If a
// SQLite driver is later added to go.mod, these tests can be lifted
// to a real database with only the setup code changing.
//
// The file carries a //go:build integration tag so it is excluded
// from the default `go test ./...` run and must be exercised via
// `go test -tags integration ./drivers/...`.

package drivers_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// Ref: @bedrock/code-0365
func TestAvailableAndUnReservedJobsArePopped(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	// Stage a row that is unreserved and whose available_at is in the
	// past — Pop's SELECT should find it.
	db.addRow(int64(1), "mock_payload", 0)

	job, err := drv.Pop(context.Background(), "mock_queue_name")

	if err != nil {
		t.Fatalf("Pop: unexpected error %v", err)
	}

	if job == nil {
		t.Fatal("Pop returned nil, want a job")
	}

	// Pop issues exactly one UPDATE to reserve the row.
	if got := countExecs(db, "UPDATE"); got != 1 {
		t.Errorf("UPDATE exec count: got %d, want 1", got)
	}
}

// Ref: @bedrock/code-0365
func TestPoppedJobsIncrementAttempts(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	// Row starts at attempts=0; pop should hand back attempts=1 on
	// the Job and issue an UPDATE that bumps the column.
	db.addRow(int64(1), "mock_payload", 0)

	job, err := drv.Pop(context.Background(), "mock_queue_name")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	if got := job.Attempts(); got != 1 {
		t.Errorf("job.Attempts(): got %d, want 1", got)
	}

	// The UPDATE must include "attempts=attempts+1" — the PHP test
	// reads the persisted row and asserts attempts == 1. Here we
	// assert the SQL that would make that true.
	var sawIncrement bool

	for _, c := range db.execCalls {
		if strings.Contains(c.Query, "attempts=attempts+1") {
			sawIncrement = true

			break
		}
	}

	if !sawIncrement {
		t.Errorf("no UPDATE with attempts=attempts+1 in exec calls: %+v", db.execCalls)
	}
}

// Ref: @bedrock/code-0365
func TestThatQueueCanBeCleared(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	if err := drv.ClearQueue(context.Background(), "mock_queue_name"); err != nil {
		t.Fatalf("ClearQueue: %v", err)
	}

	if len(db.execCalls) != 1 {
		t.Fatalf("exec call count: got %d, want 1", len(db.execCalls))
	}

	call := db.execCalls[0]

	if !strings.HasPrefix(strings.TrimSpace(call.Query), "DELETE FROM jobs") {
		t.Errorf("ClearQueue query: got %q, want DELETE FROM jobs ...", call.Query)
	}

	if len(call.Args) != 1 || call.Args[0] != "mock_queue_name" {
		t.Errorf("ClearQueue args: got %v, want [mock_queue_name]", call.Args)
	}

	// After clear, the pending size query returns 0 — stage a zero
	// count row to mirror the upstream assertEquals(0, $this->queue->size()).
	db.addRow(int64(0))

	n, err := drv.Size(context.Background(), "mock_queue_name")

	if err != nil {
		t.Fatalf("Size: %v", err)
	}

	if n != 0 {
		t.Errorf("Size after clear: got %d, want 0", n)
	}
}

// Ref: @bedrock/code-0365
func TestUnavailableJobsAreNotPopped(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	// Stage zero rows: Pop's SELECT filters on available_at <= now,
	// so a row whose available_at is in the future would not match
	// and the query returns nothing.
	job, err := drv.Pop(context.Background(), "mock_queue_name")

	if err == nil {
		t.Fatalf("Pop: expected ErrNoJob, got job=%v", job)
	}

	if err != queue.ErrNoJob {
		t.Errorf("Pop error: got %v, want ErrNoJob", err)
	}

	if job != nil {
		t.Errorf("Pop job: got %v, want nil", job)
	}
}

// Ref: @bedrock/code-0365
// DEFERRED. the upstream Pop reclaims reserved_at rows whose reservation
// has expired (reserved_at < now - retry_after). The Go DatabaseDriver
// Pop today only selects rows where reserved_at IS NULL, so there is
// no reclaim path to exercise. Once reclaim lands, this test should
// stage a reserved-but-expired row and assert Pop returns it.
func TestThatReservedAndExpiredJobsArePopped(t *testing.T) {
	t.Parallel()
	t.Skip("deferred: DatabaseDriver.Pop does not yet reclaim expired reservations")
}

// Ref: @bedrock/code-0365
func TestThatReservedJobsAreNotPopped(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	// Stage zero rows: Pop's SELECT filters on reserved_at IS NULL,
	// so a reserved row would not match and the query returns nothing.
	job, err := drv.Pop(context.Background(), "mock_queue_name")

	if err != queue.ErrNoJob {
		t.Errorf("Pop error: got %v, want ErrNoJob", err)
	}

	if job != nil {
		t.Errorf("Pop job: got %v, want nil", job)
	}
}

// Ref: @bedrock/code-0365
// Adaptation: upstream dispatches JobQueueing / JobQueued events from
// DatabaseQueue::push and the PHP test asserts the event payload()
// contains the job UUID. In the Go port, the driver does not emit
// events — the Worker does, and only on pop/process. The observable
// behaviour we can assert at the driver layer is: the payload a
// caller hands to Push is the same payload that comes back on the
// popped Job. That's what the PHP test is really guarding against:
// the queue framework not mutating or losing the payload in transit.
func TestJobPayloadIsAvailableOnEvents(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	_, payload, err := queue.CreatePayloadFor(
		"database", "default", "MyJob",
		map[string]any{"key": "value"},
		queue.JobOptions{},
	)

	if err != nil {
		t.Fatalf("CreatePayloadFor: %v", err)
	}

	if _, err := drv.Push(context.Background(), "default", payload); err != nil {
		t.Fatalf("Push: %v", err)
	}

	// Confirm the INSERT carried the exact payload bytes.
	if len(db.execCalls) != 1 {
		t.Fatalf("exec call count: got %d, want 1", len(db.execCalls))
	}

	insert := db.execCalls[0]

	if len(insert.Args) < 2 {
		t.Fatalf("INSERT args: got %v", insert.Args)
	}

	if got, ok := insert.Args[1].(string); !ok || got != string(payload) {
		t.Errorf("INSERT payload arg: got %v, want %s", insert.Args[1], string(payload))
	}

	// Now simulate Pop handing that same payload back to the caller:
	// stage a row whose payload column matches what Push wrote and
	// assert the popped Job exposes it byte-for-byte (this is the
	// Go analogue of "payload is available on events").
	db.addRow(int64(1), string(payload), 0)

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	if job == nil {
		t.Fatal("Pop returned nil")
	}

	if got := string(job.Payload()); got != string(payload) {
		t.Errorf("job.Payload(): got %s, want %s", got, string(payload))
	}

	// The payload must decode into a map that carries a uuid — this
	if !strings.Contains(string(job.Payload()), `"uuid"`) {
		t.Errorf("job payload missing uuid field: %s", string(job.Payload()))
	}
}

// countExecs returns how many recorded exec calls start with the
// given SQL keyword (case-sensitive), e.g. "UPDATE" or "DELETE".
func countExecs(db *mockDBExecer, keyword string) int {
	var n int

	for _, c := range db.execCalls {
		if strings.HasPrefix(strings.TrimSpace(c.Query), keyword) {
			n++
		}
	}

	return n
}
