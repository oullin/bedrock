package drivers_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/bedrock/packages/queue/drivers"
)

// Continued port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest —
// covers Bulk + PendingJobs + DelayedJobs + ReservedJobs.
//
// These tests exercise the new DBExecer.Query path (multi-row result
// sets) and the new Bulk method. Both were added as part of Step 7d
// database round 2. See PARITY.md §2 for adaptation rules.
//
// Upstream's tests pipe through the fluent query builder
// ($db->table->where->whereNull->get). The Go port pre-stages rows on
// the mock via addQueryRow and asserts that:
//
//   - the driver made one Query call,
//   - the SQL mentions the expected table + WHERE clause columns,
//   - the returned slice carries the decoded InspectedJob fields.

// stagePayload returns a JSON payload string matching the shape the
// Upstream fixtures use — uuid, displayName, job, data, createdAt.
func stagePayload(t *testing.T, uuid, displayName, createdAt string) string {
	t.Helper()

	raw, err := json.Marshal(map[string]any{
		"uuid":        uuid,
		"displayName": displayName,
		"job":         "foo",
		"data":        map[string]any{},
		"createdAt":   createdAtFloat(createdAt),
	})

	if err != nil {
		t.Fatalf("stagePayload: %v", err)
	}

	return string(raw)
}

// createdAtFloat parses a Upstream-style seconds-since-epoch literal
// into a float64 (which is what json.Marshal encodes numeric values
// as — the Go decode path in fetchInspected expects float64).
func createdAtFloat(s string) float64 {
	var out float64

	for _, r := range s {
		if r < '0' || r > '9' {
			continue
		}

		out = out*10 + float64(r-'0')
	}

	return out
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testBulkBatchPushesOntoDatabase
//
// Upstream asserts that DatabaseQueue::bulk issues one $db->insert call
// with an array of records. The Go port asserts one Exec call whose
// SQL carries two VALUES tuples and 8 positional args.
func TestBulkBatchPushesOntoDatabase(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	fooPayload := []byte(`{"job":"foo"}`)
	barPayload := []byte(`{"job":"bar"}`)

	if err := drv.Bulk(context.Background(), "queue", [][]byte{fooPayload, barPayload}); err != nil {
		t.Fatalf("Bulk: %v", err)
	}

	if len(db.execCalls) != 1 {
		t.Fatalf("exec calls: got %d, want 1", len(db.execCalls))
	}

	call := db.execCalls[0]

	if !strings.Contains(call.Query, "INSERT INTO jobs") {
		t.Errorf("query: got %q, want INSERT INTO jobs ...", call.Query)
	}

	// Two VALUES tuples → the SQL must contain "),(" exactly once.
	if strings.Count(call.Query, "),(") != 1 {
		t.Errorf("query: expected two VALUES tuples, got %q", call.Query)
	}

	// 4 args per row × 2 rows = 8 positional args.
	if len(call.Args) != 8 {
		t.Fatalf("args len: got %d, want 8", len(call.Args))
	}

	// Row 1: queue, payload-foo, availableAt, createdAt.
	if q, _ := call.Args[0].(string); q != "queue" {
		t.Errorf("row1 queue: got %q, want queue", q)
	}

	if p, _ := call.Args[1].(string); p != string(fooPayload) {
		t.Errorf("row1 payload: got %q, want %q", p, string(fooPayload))
	}

	// Row 2: queue, payload-bar, availableAt, createdAt.
	if q, _ := call.Args[4].(string); q != "queue" {
		t.Errorf("row2 queue: got %q, want queue", q)
	}

	if p, _ := call.Args[5].(string); p != string(barPayload) {
		t.Errorf("row2 payload: got %q, want %q", p, string(barPayload))
	}
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testPendingJobs
func TestPendingJobs(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	payload := stagePayload(t, "test-uuid", "MyTestJob", "1000000")

	db.addQueryRow(int64(1), "default", payload, 0, (*int64)(nil))

	jobs, err := drv.PendingJobs(context.Background(), "default")

	if err != nil {
		t.Fatalf("PendingJobs: %v", err)
	}

	if len(jobs) != 1 {
		t.Fatalf("jobs: got %d, want 1", len(jobs))
	}

	j := jobs[0]

	if j.Name != "MyTestJob" {
		t.Errorf("Name: got %q, want MyTestJob", j.Name)
	}

	if j.UUID != "test-uuid" {
		t.Errorf("UUID: got %q, want test-uuid", j.UUID)
	}

	if j.Attempts != 0 {
		t.Errorf("Attempts: got %d, want 0", j.Attempts)
	}

	if j.CreatedAt.Unix() != 1000000 {
		t.Errorf("CreatedAt: got %d, want 1000000", j.CreatedAt.Unix())
	}

	// And the generated SQL must filter on queue + unreserved + available.
	if len(db.queryCalls) != 1 {
		t.Fatalf("query calls: got %d, want 1", len(db.queryCalls))
	}

	q := db.queryCalls[0].Query

	for _, expect := range []string{"queue=", "reserved_at IS NULL", "available_at<="} {
		if !strings.Contains(q, expect) {
			t.Errorf("query missing %q: %s", expect, q)
		}
	}
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testDelayedJobs
func TestDelayedJobs(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	payload := stagePayload(t, "test-uuid", "MyDelayedJob", "1000000")

	db.addQueryRow(int64(2), "default", payload, 0, (*int64)(nil))

	jobs, err := drv.DelayedJobs(context.Background(), "default")

	if err != nil {
		t.Fatalf("DelayedJobs: %v", err)
	}

	if len(jobs) != 1 {
		t.Fatalf("jobs: got %d, want 1", len(jobs))
	}

	j := jobs[0]

	if j.Name != "MyDelayedJob" {
		t.Errorf("Name: got %q, want MyDelayedJob", j.Name)
	}

	if j.UUID != "test-uuid" {
		t.Errorf("UUID: got %q, want test-uuid", j.UUID)
	}

	if j.Attempts != 0 {
		t.Errorf("Attempts: got %d, want 0", j.Attempts)
	}

	if j.CreatedAt.Unix() != 1000000 {
		t.Errorf("CreatedAt: got %d, want 1000000", j.CreatedAt.Unix())
	}

	if len(db.queryCalls) != 1 {
		t.Fatalf("query calls: got %d, want 1", len(db.queryCalls))
	}

	q := db.queryCalls[0].Query

	for _, expect := range []string{"queue=", "reserved_at IS NULL", "available_at>"} {
		if !strings.Contains(q, expect) {
			t.Errorf("query missing %q: %s", expect, q)
		}
	}
}

// Port of Framework\Tests\Queue\QueueDatabaseQueueUnitTest::testReservedJobs
func TestReservedJobs(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	payload := stagePayload(t, "test-uuid", "MyTestJob", "1000000")
	reservedAt := int64(1700000000)

	db.addQueryRow(int64(1), "default", payload, 1, &reservedAt)

	jobs, err := drv.ReservedJobs(context.Background(), "default")

	if err != nil {
		t.Fatalf("ReservedJobs: %v", err)
	}

	if len(jobs) != 1 {
		t.Fatalf("jobs: got %d, want 1", len(jobs))
	}

	j := jobs[0]

	if j.Name != "MyTestJob" {
		t.Errorf("Name: got %q, want MyTestJob", j.Name)
	}

	if j.UUID != "test-uuid" {
		t.Errorf("UUID: got %q, want test-uuid", j.UUID)
	}

	if j.Attempts != 1 {
		t.Errorf("Attempts: got %d, want 1", j.Attempts)
	}

	if j.CreatedAt.Unix() != 1000000 {
		t.Errorf("CreatedAt: got %d, want 1000000", j.CreatedAt.Unix())
	}

	if j.ReservedAt == nil || j.ReservedAt.Unix() != reservedAt {
		t.Errorf("ReservedAt: got %v, want %d", j.ReservedAt, reservedAt)
	}

	if len(db.queryCalls) != 1 {
		t.Fatalf("query calls: got %d, want 1", len(db.queryCalls))
	}

	q := db.queryCalls[0].Query

	for _, expect := range []string{"queue=", "reserved_at IS NOT NULL"} {
		if !strings.Contains(q, expect) {
			t.Errorf("query missing %q: %s", expect, q)
		}
	}
}
