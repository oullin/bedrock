package failed_test

// Ported from @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest (11/11).
//
// ✅ testGettingIdsOfAllFailedJobs
// ✅ testGettingAllFailedJobs
// ✅ testFindingFailedJobsById
// ✅ testRemovingJobsById
// ✅ testRemovingAllFailedJobs
// ✅ testPruningFailedJobs
// ✅ testPruningFailedJobsWithRelativeHoursAndMinutes
// ✅ testJobsCanBeCounted
// ✅ testJobsCanBeCountedByConnection
// ✅ testJobsCanBeCountedByQueue
// ✅ testJobsCanBeCountedByQueueAndConnection

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue/failed"
)

func newUUIDProvider() *failed.DatabaseUuidFailedJobProvider {
	return failed.NewDatabaseUuidFailedJobProvider("failed_jobs")
}

func uuidPayload(uuid string) string {
	b, _ := json.Marshal(map[string]string{"uuid": uuid})

	return string(b)
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testGettingIdsOfAllFailedJobs
func TestGettingIdsOfAllFailedJobs(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()

	p.Log("connection-1", "queue-1", uuidPayload("uuid-1"), errors.New("e"))
	p.Log("connection-1", "queue-1", uuidPayload("uuid-2"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("uuid-3"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("uuid-4"), errors.New("e"))

	ids, _ := p.IDs("")
	want := []string{"uuid-1", "uuid-2", "uuid-3", "uuid-4"}
	assertStringSlice(t, ids, want)

	ids, _ = p.IDs("queue-1")
	assertStringSlice(t, ids, []string{"uuid-1", "uuid-2"})

	ids, _ = p.IDs("queue-2")
	assertStringSlice(t, ids, []string{"uuid-3", "uuid-4"})
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testGettingAllFailedJobs
func TestGettingAllFailedJobs(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()

	all, _ := p.All()

	if len(all) != 0 {
		t.Fatal("expected empty")
	}

	p.Log("connection-1", "queue-1", uuidPayload("uuid-1"), errors.New("e"))
	p.Log("connection-1", "queue-1", uuidPayload("uuid-2"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("uuid-3"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("uuid-4"), errors.New("e"))

	all, _ = p.All()

	if len(all) != 4 {
		t.Fatalf("count=%d", len(all))
	}

	want := []string{"uuid-1", "uuid-2", "uuid-3", "uuid-4"}

	for i, j := range all {
		if j.ID != want[i] {
			t.Fatalf("all[%d].ID=%q want %q", i, j.ID, want[i])
		}
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testFindingFailedJobsById
func TestFindingFailedJobsById(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	p.Log("connection-1", "queue-1", uuidPayload("uuid-1"), errors.New("e"))

	if j, _ := p.Find("uuid-2"); j != nil {
		t.Fatal("want nil")
	}

	j, _ := p.Find("uuid-1")

	if j == nil || j.ID != "uuid-1" || j.Queue != "queue-1" || j.Connection != "connection-1" {
		t.Fatalf("unexpected job: %+v", j)
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testRemovingJobsById
func TestRemovingJobsById(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	p.Log("connection-1", "queue-1", uuidPayload("uuid-1"), errors.New("e"))

	if j, _ := p.Find("uuid-1"); j == nil {
		t.Fatal("expected job")
	}

	p.Forget("uuid-1")

	if j, _ := p.Find("uuid-1"); j != nil {
		t.Fatal("want nil")
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testRemovingAllFailedJobs
func TestRemovingAllFailedJobs(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	p.Log("connection-1", "queue-1", uuidPayload("uuid-1"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("uuid-2"), errors.New("e"))

	all, _ := p.All()

	if len(all) != 2 {
		t.Fatalf("count=%d", len(all))
	}

	p.Flush(0)
	all, _ = p.All()

	if len(all) != 0 {
		t.Fatalf("expected empty, got %d", len(all))
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testPruningFailedJobs
func TestPruningFailedJobs(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	now := time.Date(2024, 4, 28, 0, 0, 0, 0, time.UTC)
	failed.SetNow(p, func() time.Time { return now })

	p.Log("connection-1", "queue-1", uuidPayload("uuid-1"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("uuid-2"), errors.New("e"))

	p.Prune(time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC))
	all, _ := p.All()

	if len(all) != 2 {
		t.Fatalf("count=%d", len(all))
	}

	p.Prune(time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC))
	all, _ = p.All()

	if len(all) != 0 {
		t.Fatalf("expected empty, got %d", len(all))
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testPruningFailedJobsWithRelativeHoursAndMinutes
func TestPruningFailedJobsWithRelativeHoursAndMinutes(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	now := time.Date(2025, 8, 24, 12, 30, 0, 0, time.UTC)
	failed.SetNow(p, func() time.Time { return now })

	p.Log("connection-1", "queue-1", uuidPayload("uuid-1"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("uuid-2"), errors.New("e"))

	p.Prune(time.Date(2025, 8, 24, 12, 30, 0, 0, time.UTC))
	all, _ := p.All()

	if len(all) != 2 {
		t.Fatalf("count=%d", len(all))
	}

	p.Prune(time.Date(2025, 8, 24, 13, 0, 0, 0, time.UTC))
	all, _ = p.All()

	if len(all) != 0 {
		t.Fatalf("expected empty, got %d", len(all))
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testJobsCanBeCounted
func TestJobsCanBeCountedUuid(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatal()
	}

	p.Log("connection-1", "queue-1", uuidPayload("a"), errors.New("e"))

	if c, _ := p.Count("", ""); c != 1 {
		t.Fatal()
	}

	p.Log("connection-1", "queue-1", uuidPayload("b"), errors.New("e"))
	p.Log("connection-2", "queue-2", uuidPayload("c"), errors.New("e"))

	if c, _ := p.Count("", ""); c != 3 {
		t.Fatal()
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testJobsCanBeCountedByConnection
func TestJobsCanBeCountedByConnectionUuid(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	p.Log("connection-1", "default", uuidPayload("a"), errors.New("e"))
	p.Log("connection-2", "default", uuidPayload("b"), errors.New("e"))

	if c, _ := p.Count("connection-1", ""); c != 1 {
		t.Fatal()
	}

	if c, _ := p.Count("connection-2", ""); c != 1 {
		t.Fatal()
	}

	p.Log("connection-1", "default", uuidPayload("c"), errors.New("e"))

	if c, _ := p.Count("connection-1", ""); c != 2 {
		t.Fatal()
	}

	if c, _ := p.Count("connection-2", ""); c != 1 {
		t.Fatal()
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testJobsCanBeCountedByQueue
func TestJobsCanBeCountedByQueueUuid(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	p.Log("database", "queue-1", uuidPayload("a"), errors.New("e"))
	p.Log("database", "queue-2", uuidPayload("b"), errors.New("e"))

	if c, _ := p.Count("", "queue-1"); c != 1 {
		t.Fatal()
	}

	if c, _ := p.Count("", "queue-2"); c != 1 {
		t.Fatal()
	}

	p.Log("database", "queue-1", uuidPayload("c"), errors.New("e"))

	if c, _ := p.Count("", "queue-1"); c != 2 {
		t.Fatal()
	}

	if c, _ := p.Count("", "queue-2"); c != 1 {
		t.Fatal()
	}
}

// Port of @bedrock\Tests\Queue\DatabaseUuidFailedJobProviderTest::testJobsCanBeCountedByQueueAndConnection
func TestJobsCanBeCountedByQueueAndConnectionUuid(t *testing.T) {
	t.Parallel()
	p := newUUIDProvider()
	p.Log("connection-1", "queue-99", uuidPayload("a"), errors.New("e"))
	p.Log("connection-1", "queue-99", uuidPayload("b"), errors.New("e"))
	p.Log("connection-2", "queue-99", uuidPayload("c"), errors.New("e"))
	p.Log("connection-1", "queue-1", uuidPayload("d"), errors.New("e"))
	p.Log("connection-2", "queue-1", uuidPayload("e"), errors.New("e"))
	p.Log("connection-2", "queue-1", uuidPayload("f"), errors.New("e"))

	if c, _ := p.Count("connection-1", "queue-99"); c != 2 {
		t.Fatal()
	}

	if c, _ := p.Count("connection-2", "queue-99"); c != 1 {
		t.Fatal()
	}

	if c, _ := p.Count("connection-1", "queue-1"); c != 1 {
		t.Fatal()
	}

	if c, _ := p.Count("connection-2", "queue-1"); c != 2 {
		t.Fatal()
	}
}

func assertStringSlice(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len=%d want %d (got %v)", len(got), len(want), got)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%q want %q", i, got[i], want[i])
		}
	}
}
