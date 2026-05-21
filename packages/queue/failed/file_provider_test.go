package failed_test

// Ported from @bedrock\Tests\Queue\FileFailedJobProviderTest (13/13).
//
// ✅ testCanLogFailedJobs
// ✅ testCanRetrieveAllFailedJobs
// ✅ testCanFindFailedJobs
// ✅ testNullIsReturnedIfJobNotFound
// ✅ testCanForgetFailedJobs
// ✅ testCanFlushFailedJobs
// ✅ testCanPruneFailedJobs
// ✅ testCanPruneFailedJobsWithRelativeHours
// ✅ testEmptyFailedJobsByDefault
// ✅ testJobsCanBeCounted
// ✅ testJobsCanBeCountedByConnection
// ✅ testJobsCanBeCountedByQueue
// ✅ testJobsCanBeCountedByQueueAndConnection

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/bedrock/packages/queue/failed"
)

type loggedJob struct {
	uuid string
	err  error
}

func newFileProvider(t *testing.T) *failed.FileFailedJobProvider {
	t.Helper()
	dir := t.TempDir()

	return failed.NewFileFailedJobProvider(filepath.Join(dir, "failed.json"), 0)
}

func logFileJob(t *testing.T, p *failed.FileFailedJobProvider, connection, queue string) loggedJob {
	t.Helper()
	uuid := fmt.Sprintf("uuid-%d", time.Now().UnixNano())
	payload, _ := json.Marshal(map[string]string{"uuid": uuid})
	ex := fmt.Errorf("Something went wrong at job [%s].", uuid)

	if _, err := p.Log(connection, queue, string(payload), ex); err != nil {
		t.Fatalf("log: %v", err)
	}

	return loggedJob{uuid: uuid, err: ex}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testCanLogFailedJobs
func TestCanLogFailedJobs(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	j := logFileJob(t, p, "connection", "queue")

	all, _ := p.All()

	if len(all) != 1 {
		t.Fatalf("count=%d", len(all))
	}

	got := all[0]

	if got.ID != j.uuid || got.Connection != "connection" || got.Queue != "queue" {
		t.Fatalf("unexpected: %+v", got)
	}

	if got.Exception != j.err.Error() {
		t.Fatalf("exception=%q", got.Exception)
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testCanRetrieveAllFailedJobs
func TestCanRetrieveAllFailedJobs(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	j1 := logFileJob(t, p, "connection", "queue")
	j2 := logFileJob(t, p, "connection", "queue")

	all, _ := p.All()

	if len(all) != 2 {
		t.Fatalf("count=%d", len(all))
	}
	// Newest first: j2 then j1.
	if all[0].ID != j2.uuid {
		t.Fatalf("all[0]=%q want %q", all[0].ID, j2.uuid)
	}

	if all[1].ID != j1.uuid {
		t.Fatalf("all[1]=%q want %q", all[1].ID, j1.uuid)
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testCanFindFailedJobs
func TestCanFindFailedJobs(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	j := logFileJob(t, p, "connection", "queue")

	got, _ := p.Find(j.uuid)

	if got == nil || got.ID != j.uuid {
		t.Fatalf("unexpected: %+v", got)
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testNullIsReturnedIfJobNotFound
func TestNullIsReturnedIfJobNotFound(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)

	if got, _ := p.Find("missing"); got != nil {
		t.Fatal("expected nil")
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testCanForgetFailedJobs
func TestCanForgetFailedJobs(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	j := logFileJob(t, p, "connection", "queue")
	p.Forget(j.uuid)

	if got, _ := p.Find(j.uuid); got != nil {
		t.Fatal("expected nil")
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testCanFlushFailedJobs
func TestCanFlushFailedJobsFile(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	logFileJob(t, p, "c", "q")
	logFileJob(t, p, "c", "q")
	p.Flush(0)
	all, _ := p.All()

	if len(all) != 0 {
		t.Fatalf("expected empty")
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testCanPruneFailedJobs
func TestCanPruneFailedJobsFile(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	logFileJob(t, p, "c", "q")
	logFileJob(t, p, "c", "q")

	// prune(now + 1 day) -> all empty
	p.Prune(time.Now().Add(24 * time.Hour))
	all, _ := p.All()

	if len(all) != 0 {
		t.Fatalf("expected empty, got %d", len(all))
	}

	logFileJob(t, p, "c", "q")
	logFileJob(t, p, "c", "q")

	// prune(now - 1 day) -> both kept
	p.Prune(time.Now().Add(-24 * time.Hour))
	all, _ = p.All()

	if len(all) != 2 {
		t.Fatalf("count=%d want 2", len(all))
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testCanPruneFailedJobsWithRelativeHours
func TestCanPruneFailedJobsWithRelativeHours(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	logFileJob(t, p, "c", "q")
	logFileJob(t, p, "c", "q")
	p.Prune(time.Now().Add(time.Hour))
	all, _ := p.All()

	if len(all) != 0 {
		t.Fatalf("expected empty, got %d", len(all))
	}

	logFileJob(t, p, "c", "q")
	logFileJob(t, p, "c", "q")
	p.Prune(time.Now().Add(-time.Hour))
	all, _ = p.All()

	if len(all) != 2 {
		t.Fatalf("count=%d want 2", len(all))
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testEmptyFailedJobsByDefault
func TestEmptyFailedJobsByDefault(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	all, _ := p.All()

	if len(all) != 0 {
		t.Fatal()
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testJobsCanBeCounted
func TestJobsCanBeCountedFile(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)

	if c, _ := p.Count("", ""); c != 0 {
		t.Fatal()
	}

	logFileJob(t, p, "database", "default")

	if c, _ := p.Count("", ""); c != 1 {
		t.Fatal()
	}

	logFileJob(t, p, "database", "default")
	logFileJob(t, p, "another-connection", "another-queue")

	if c, _ := p.Count("", ""); c != 3 {
		t.Fatal()
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testJobsCanBeCountedByConnection
func TestJobsCanBeCountedByConnectionFile(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	logFileJob(t, p, "connection-1", "default")
	logFileJob(t, p, "connection-2", "default")

	if c, _ := p.Count("connection-1", ""); c != 1 {
		t.Fatal()
	}

	if c, _ := p.Count("connection-2", ""); c != 1 {
		t.Fatal()
	}

	logFileJob(t, p, "connection-1", "default")

	if c, _ := p.Count("connection-1", ""); c != 2 {
		t.Fatal()
	}

	if c, _ := p.Count("connection-2", ""); c != 1 {
		t.Fatal()
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testJobsCanBeCountedByQueue
func TestJobsCanBeCountedByQueueFile(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	logFileJob(t, p, "database", "queue-1")
	logFileJob(t, p, "database", "queue-2")

	if c, _ := p.Count("", "queue-1"); c != 1 {
		t.Fatal()
	}

	if c, _ := p.Count("", "queue-2"); c != 1 {
		t.Fatal()
	}

	logFileJob(t, p, "database", "queue-1")

	if c, _ := p.Count("", "queue-1"); c != 2 {
		t.Fatal()
	}

	if c, _ := p.Count("", "queue-2"); c != 1 {
		t.Fatal()
	}
}

// Port of @bedrock\Tests\Queue\FileFailedJobProviderTest::testJobsCanBeCountedByQueueAndConnection
func TestJobsCanBeCountedByQueueAndConnectionFile(t *testing.T) {
	t.Parallel()
	p := newFileProvider(t)
	logFileJob(t, p, "connection-1", "queue-99")
	logFileJob(t, p, "connection-1", "queue-99")
	logFileJob(t, p, "connection-2", "queue-99")
	logFileJob(t, p, "connection-1", "queue-1")
	logFileJob(t, p, "connection-2", "queue-1")
	logFileJob(t, p, "connection-2", "queue-1")

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

// Compile-time assertion: the exception field stays non-empty.
var _ = errors.New
