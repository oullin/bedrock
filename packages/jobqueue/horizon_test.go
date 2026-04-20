package jobqueue

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMonitorCaptureRecordsSortedSnapshot(t *testing.T) {
	t.Parallel()

	repository := NewInMemoryRepository()
	monitor := NewMonitor(repository,
		QueueSourceFunc(func(ctx context.Context) ([]QueueStatus, error) {
			return []QueueStatus{
				{Name: "redis:mail", Pending: 3},
				{Name: "redis:default", Pending: 1, Processing: 2},
			}, nil
		}),
	)
	monitor.now = func() time.Time {
		return time.Date(2026, 4, 20, 8, 0, 0, 0, time.FixedZone("SGT", 8*60*60))
	}

	snapshot, err := monitor.Capture(context.Background())

	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}

	if snapshot.GeneratedAt.Location() != time.UTC {
		t.Fatalf("GeneratedAt should be UTC, got %s", snapshot.GeneratedAt.Location())
	}

	if got := snapshot.Queues[0].Name; got != "redis:default" {
		t.Fatalf("first queue = %q, want redis:default", got)
	}

	latest, err := repository.Latest(context.Background())

	if err != nil {
		t.Fatalf("Latest returned error: %v", err)
	}

	if latest.Queues[1].Pending != 3 {
		t.Fatalf("latest snapshot was not recorded")
	}
}

func TestInMemoryRepositoryLatestWithoutSnapshots(t *testing.T) {
	t.Parallel()

	_, err := NewInMemoryRepository().Latest(context.Background())

	if !errors.Is(err, ErrNoSnapshots) {
		t.Fatalf("Latest error = %v, want ErrNoSnapshots", err)
	}
}
