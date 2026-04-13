package drivers_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestDatabaseDriverPush(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	_, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}

	if !strings.Contains(db.execCalls[0].Query, "INSERT INTO jobs") {
		t.Errorf("expected INSERT INTO jobs, got %q", db.execCalls[0].Query)
	}
}

func TestDatabaseDriverPushDelayed(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	_, err := drv.PushDelayed(context.Background(), "default", []byte("payload"), 10*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}

	// available_at should be in the future.
	args := db.execCalls[0].Args

	if len(args) < 3 {
		t.Fatal("expected at least 3 args")
	}

	availAt, ok := args[2].(int64)

	if !ok {
		t.Fatal("expected int64 available_at")
	}

	now := time.Now().Unix()

	if availAt <= now {
		t.Errorf("expected available_at > now, got %d <= %d", availAt, now)
	}
}

func TestDatabaseDriverPushMultiple(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	ids, err := drv.PushMultiple(context.Background(), "default", [][]byte{
		[]byte("a"), []byte("b"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 2 {
		t.Errorf("expected 2 ids, got %d", len(ids))
	}

	if len(db.execCalls) != 2 {
		t.Errorf("expected 2 exec calls, got %d", len(db.execCalls))
	}
}

func TestDatabaseDriverPopReservesJob(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(42), "test-payload", int(0))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "test-payload" {
		t.Errorf("expected 'test-payload', got %q", job.Payload())
	}

	if job.Attempts() != 1 {
		t.Errorf("expected attempts 1, got %d", job.Attempts())
	}

	// Should have called UPDATE to reserve.
	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call for reservation, got %d", len(db.execCalls))
	}

	if !strings.Contains(db.execCalls[0].Query, "UPDATE") {
		t.Errorf("expected UPDATE query, got %q", db.execCalls[0].Query)
	}
}

func TestDatabaseDriverPopEmptyReturnsErrNoJob(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addErrorRow(errors.New("no rows"))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	_, err := drv.Pop(context.Background(), "default")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestDatabaseDriverJobRelease(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(1), "payload", int(0))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	job, _ := drv.Pop(context.Background(), "default")
	db.execCalls = nil // Reset after reservation UPDATE.

	err := job.Release(5 * time.Second)

	if err != nil {
		t.Fatal(err)
	}

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}

	if !strings.Contains(db.execCalls[0].Query, "reserved_at=NULL") {
		t.Errorf("expected reserved_at=NULL in query, got %q", db.execCalls[0].Query)
	}
}

func TestDatabaseDriverJobDelete(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(1), "payload", int(0))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	job, _ := drv.Pop(context.Background(), "default")
	db.execCalls = nil

	err := job.Delete()

	if err != nil {
		t.Fatal(err)
	}

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}

	if !strings.Contains(db.execCalls[0].Query, "DELETE") {
		t.Errorf("expected DELETE query, got %q", db.execCalls[0].Query)
	}
}

func TestDatabaseDriverJobFail(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(1), "payload", int(0))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	job, _ := drv.Pop(context.Background(), "default")
	db.execCalls = nil

	err := job.Fail(errors.New("oops"))

	if err != nil {
		t.Fatal(err)
	}

	// Should INSERT into failed_jobs and DELETE from jobs.
	if len(db.execCalls) < 2 {
		t.Fatalf("expected at least 2 exec calls, got %d", len(db.execCalls))
	}

	foundInsert := false
	foundDelete := false

	for _, call := range db.execCalls {
		if strings.Contains(call.Query, "failed_jobs") {
			foundInsert = true
		}

		if strings.Contains(call.Query, "DELETE") {
			foundDelete = true
		}
	}

	if !foundInsert {
		t.Error("expected INSERT into failed_jobs")
	}

	if !foundDelete {
		t.Error("expected DELETE from jobs")
	}
}

func TestDatabaseDriverSize(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(5))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	n, err := drv.Size(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 5 {
		t.Errorf("expected 5, got %d", n)
	}
}

func TestDatabaseDriverPendingSize(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(3))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	n, err := drv.PendingSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}

func TestDatabaseDriverDelayedSize(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(2))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	n, err := drv.DelayedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestDatabaseDriverReservedSize(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	db.addRow(int64(1))
	drv := drivers.NewDatabaseDriver(db, "jobs", "database")

	n, err := drv.ReservedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
}

func TestDatabaseDriverCustomTable(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "custom_jobs", "database")

	_, _ = drv.Push(context.Background(), "default", []byte("p"))

	if !strings.Contains(db.execCalls[0].Query, "custom_jobs") {
		t.Errorf("expected custom_jobs in query, got %q", db.execCalls[0].Query)
	}
}

func TestDatabaseDriverDefaultTable(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "", "database")

	_, _ = drv.Push(context.Background(), "default", []byte("p"))

	if !strings.Contains(db.execCalls[0].Query, "INSERT INTO jobs") {
		t.Errorf("expected default 'jobs' table, got %q", db.execCalls[0].Query)
	}
}

func TestDatabaseDriverConnectionName(t *testing.T) {
	t.Parallel()

	db := newMockDBExecer()
	drv := drivers.NewDatabaseDriver(db, "jobs", "my-db")

	if drv.ConnectionName() != "my-db" {
		t.Errorf("expected 'my-db', got %q", drv.ConnectionName())
	}
}
