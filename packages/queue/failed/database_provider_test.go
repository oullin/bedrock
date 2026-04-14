package failed_test

// Ported from Framework\Tests\Queue\DatabaseFailedJobProviderTest (12/12).
//
// ✅ testCanGetAllFailedJobIds
// ✅ testCanGetAllFailedJobs
// ✅ testCanRetrieveFailedJobsById
// ✅ testCanRemoveFailedJobsById
// ✅ testCanPruneFailedJobs
// ✅ testCanPruneFailedJobsWithRelativeHoursAndMinutes
// ✅ testCanFlushFailedJobs
// ✅ testCanProperlyLogFailedJob
// ✅ testJobsCanBeCounted
// ✅ testJobsCanBeCountedByConnection
// ✅ testJobsCanBeCountedByQueue
// ✅ testJobsCanBeCountedByQueueAndConnection

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bedrock/packages/queue/failed"
)

func newIntProvider() *failed.DatabaseFailedJobProvider {
	return failed.NewDatabaseFailedJobProvider("failed_jobs")
}

func createIntFailedJob(t *testing.T, p *failed.DatabaseFailedJobProvider, failedAt time.Time) {
	t.Helper()
	uuid := fmt.Sprintf("uuid-%d", time.Now().UnixNano())
	payload, _ := json.Marshal(map[string]string{"uuid": uuid})
	failed.SetNow(p, func() time.Time { return failedAt })

	if _, err := p.Log("database", "default", string(payload), errors.New("Whoops!")); err != nil {
		t.Fatalf("log: %v", err)
	}

	failed.SetNow(p, time.Now)
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanGetAllFailedJobIds
func TestCanGetAllFailedJobIds(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	ids, _ := p.IDs("")

	if len(ids) != 0 {
		t.Fatalf("expected empty, got %v", ids)
	}

	for i := 0; i < 4; i++ {
		createIntFailedJob(t, p, time.Now())
	}

	ids, _ = p.IDs("")

	if len(ids) != 4 {
		t.Fatalf("expected 4 ids, got %d", len(ids))
	}

	want := []string{"4", "3", "2", "1"}

	for i, id := range ids {
		if id != want[i] {
			t.Fatalf("ids[%d]=%s want %s", i, id, want[i])
		}
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanGetAllFailedJobs
func TestCanGetAllFailedJobs(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	all, _ := p.All()

	if len(all) != 0 {
		t.Fatalf("expected empty, got %d", len(all))
	}

	for i := 0; i < 4; i++ {
		createIntFailedJob(t, p, time.Now())
	}

	all, _ = p.All()

	if len(all) != 4 {
		t.Fatalf("expected 4, got %d", len(all))
	}

	if all[1].ID != "3" {
		t.Fatalf("all[1].ID=%q want 3", all[1].ID)
	}

	if all[1].Queue != "default" {
		t.Fatalf("all[1].Queue=%q want default", all[1].Queue)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanRetrieveFailedJobsById
func TestCanRetrieveFailedJobsById(t *testing.T) {
	t.Parallel()
	p := newIntProvider()
	createIntFailedJob(t, p, time.Now())
	createIntFailedJob(t, p, time.Now())

	if j, _ := p.Find("1"); j == nil {
		t.Fatal("find 1 nil")
	}

	if j, _ := p.Find("2"); j == nil {
		t.Fatal("find 2 nil")
	}

	if j, _ := p.Find("3"); j != nil {
		t.Fatal("find 3 not nil")
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanRemoveFailedJobsById
func TestCanRemoveFailedJobsById(t *testing.T) {
	t.Parallel()
	p := newIntProvider()
	createIntFailedJob(t, p, time.Now())

	ok, _ := p.Forget("2")

	if ok {
		t.Fatal("forget(2) should be false")
	}

	if c, _ := p.Count("", ""); c != 1 {
		t.Fatalf("count=%d want 1", c)
	}

	ok, _ = p.Forget("1")

	if !ok {
		t.Fatal("forget(1) should be true")
	}

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatalf("count=%d want 0", c)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanPruneFailedJobs
func TestCanPruneFailedJobs(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	createIntFailedJob(t, p, time.Date(2024, 4, 24, 0, 0, 0, 0, time.UTC))
	createIntFailedJob(t, p, time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC))

	if _, err := p.Prune(time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}

	if c, _ := p.Count("", ""); c != 2 {
		t.Fatalf("count=%d want 2", c)
	}

	if _, err := p.Prune(time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}

	if c, _ := p.Count("", ""); c != 1 {
		t.Fatalf("count=%d want 1", c)
	}

	if _, err := p.Prune(time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatalf("count=%d want 0", c)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanPruneFailedJobsWithRelativeHoursAndMinutes
func TestCanPruneFailedJobsWithRelativeHoursAndMinutes(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	createIntFailedJob(t, p, time.Date(2025, 8, 24, 11, 45, 0, 0, time.UTC))
	createIntFailedJob(t, p, time.Date(2025, 8, 24, 13, 0, 0, 0, time.UTC))

	// prune at failed_at of first row: row 1 has failed_at == cutoff, not strictly <, so kept.
	p.Prune(time.Date(2025, 8, 24, 11, 45, 0, 0, time.UTC))

	if c, _ := p.Count("", ""); c != 2 {
		t.Fatalf("count=%d want 2", c)
	}

	p.Prune(time.Date(2025, 8, 24, 14, 0, 0, 0, time.UTC))

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatalf("count=%d want 0", c)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanFlushFailedJobs
func TestCanFlushFailedJobs(t *testing.T) {
	t.Parallel()
	p := newIntProvider()
	now := time.Now()
	failed.SetNow(p, func() time.Time { return now })

	createIntFailedJob(t, p, now.AddDate(0, 0, -10))
	failed.SetNow(p, func() time.Time { return now })
	p.Flush(0)

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatalf("count=%d want 0", c)
	}

	createIntFailedJob(t, p, now.AddDate(0, 0, -10))
	failed.SetNow(p, func() time.Time { return now })
	p.Flush(15 * 24)

	if c, _ := p.Count("", ""); c != 1 {
		t.Fatalf("count=%d want 1", c)
	}

	createIntFailedJob(t, p, now.AddDate(0, 0, -10))
	failed.SetNow(p, func() time.Time { return now })
	p.Flush(10 * 24)

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatalf("count=%d want 0", c)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testCanProperlyLogFailedJob
func TestCanProperlyLogFailedJob(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	uuid := "abc-123"
	payload, _ := json.Marshal(map[string]string{"uuid": uuid})
	exception := errors.New("something went wrong")

	if _, err := p.Log("database", "default", string(payload), exception); err != nil {
		t.Fatalf("log: %v", err)
	}

	if c, _ := p.Count("", ""); c != 1 {
		t.Fatalf("count=%d want 1", c)
	}

	all, _ := p.All()

	if all[0].Exception != exception.Error() {
		t.Fatalf("exception=%q want %q", all[0].Exception, exception.Error())
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testJobsCanBeCounted
func TestJobsCanBeCounted(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatalf("count=%d want 0", c)
	}

	createIntFailedJob(t, p, time.Now())

	if c, _ := p.Count("", ""); c != 1 {
		t.Fatalf("count=%d want 1", c)
	}

	createIntFailedJob(t, p, time.Now())
	// another-connection row
	payload, _ := json.Marshal(map[string]string{"uuid": "z"})
	p.Log("another-connection", "another-queue", string(payload), errors.New("e"))

	if c, _ := p.Count("", ""); c != 3 {
		t.Fatalf("count=%d want 3", c)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testJobsCanBeCountedByConnection
func TestJobsCanBeCountedByConnection(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	payload, _ := json.Marshal(map[string]string{"uuid": "a"})
	p.Log("connection-1", "default", string(payload), errors.New("e"))
	p.Log("connection-2", "default", string(payload), errors.New("e"))

	if c, _ := p.Count("connection-1", ""); c != 1 {
		t.Fatalf("c1=%d want 1", c)
	}

	if c, _ := p.Count("connection-2", ""); c != 1 {
		t.Fatalf("c2=%d want 1", c)
	}

	p.Log("connection-1", "default", string(payload), errors.New("e"))

	if c, _ := p.Count("connection-1", ""); c != 2 {
		t.Fatalf("c1=%d want 2", c)
	}

	if c, _ := p.Count("connection-2", ""); c != 1 {
		t.Fatalf("c2=%d want 1", c)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testJobsCanBeCountedByQueue
func TestJobsCanBeCountedByQueue(t *testing.T) {
	t.Parallel()
	p := newIntProvider()

	payload, _ := json.Marshal(map[string]string{"uuid": "a"})
	p.Log("database", "queue-1", string(payload), errors.New("e"))
	p.Log("database", "queue-2", string(payload), errors.New("e"))

	if c, _ := p.Count("", "queue-1"); c != 1 {
		t.Fatalf("q1=%d want 1", c)
	}

	if c, _ := p.Count("", "queue-2"); c != 1 {
		t.Fatalf("q2=%d want 1", c)
	}

	p.Log("database", "queue-1", string(payload), errors.New("e"))

	if c, _ := p.Count("", "queue-1"); c != 2 {
		t.Fatalf("q1=%d want 2", c)
	}

	if c, _ := p.Count("", "queue-2"); c != 1 {
		t.Fatalf("q2=%d want 1", c)
	}
}

// Port of Framework\Tests\Queue\DatabaseFailedJobProviderTest::testJobsCanBeCountedByQueueAndConnection
func TestJobsCanBeCountedByQueueAndConnection(t *testing.T) {
	t.Parallel()
	p := newIntProvider()
	payload, _ := json.Marshal(map[string]string{"uuid": "a"})

	p.Log("connection-1", "queue-99", string(payload), errors.New("e"))
	p.Log("connection-1", "queue-99", string(payload), errors.New("e"))
	p.Log("connection-2", "queue-99", string(payload), errors.New("e"))
	p.Log("connection-1", "queue-1", string(payload), errors.New("e"))
	p.Log("connection-2", "queue-1", string(payload), errors.New("e"))
	p.Log("connection-2", "queue-1", string(payload), errors.New("e"))

	if c, _ := p.Count("connection-1", "queue-99"); c != 2 {
		t.Fatalf("got %d", c)
	}

	if c, _ := p.Count("connection-2", "queue-99"); c != 1 {
		t.Fatalf("got %d", c)
	}

	if c, _ := p.Count("connection-1", "queue-1"); c != 1 {
		t.Fatalf("got %d", c)
	}

	if c, _ := p.Count("connection-2", "queue-1"); c != 2 {
		t.Fatalf("got %d", c)
	}
}
