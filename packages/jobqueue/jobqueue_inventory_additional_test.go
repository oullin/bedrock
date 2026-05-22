package jobqueue

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type failingRepository struct {
	err error
}

func (r failingRepository) Record(context.Context, Snapshot) error {
	return r.err
}

func (r failingRepository) Latest(context.Context) (Snapshot, error) {
	return Snapshot{}, r.err
}

func TestMonitorInventoryParityAdditional(t *testing.T) {
	t.Run("monitor can capture without persisting when no repository is configured", func(t *testing.T) {
		monitor := NewMonitor(nil,
			QueueSourceFunc(func(ctx context.Context) ([]QueueStatus, error) {
				return []QueueStatus{
					{Name: "redis:mail", Pending: 3},
					{Name: "redis:default", Pending: 1},
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

		if got, want := snapshot.Queues[0].Name, "redis:default"; got != want {
			t.Fatalf("first queue = %q, want %q", got, want)
		}
	})

	t.Run("monitor capture returns source errors", func(t *testing.T) {
		want := errors.New("source unavailable")
		monitor := NewMonitor(nil, QueueSourceFunc(func(context.Context) ([]QueueStatus, error) {
			return nil, want
		}))

		_, err := monitor.Capture(context.Background())

		if !errors.Is(err, want) {
			t.Fatalf("Capture error = %v, want %v", err, want)
		}
	})

	t.Run("monitor capture returns repository errors", func(t *testing.T) {
		want := errors.New("record failed")
		monitor := NewMonitor(failingRepository{err: want}, QueueSourceFunc(func(context.Context) ([]QueueStatus, error) {
			return []QueueStatus{{Name: "redis:default"}}, nil
		}))

		_, err := monitor.Capture(context.Background())

		if !errors.Is(err, want) {
			t.Fatalf("Capture error = %v, want %v", err, want)
		}
	})

	t.Run("record and latest respect context cancellation", func(t *testing.T) {
		repository := NewInMemoryRepository()

		canceledRecord, cancelRecord := context.WithCancel(context.Background())
		cancelRecord()

		if err := repository.Record(canceledRecord, Snapshot{}); !errors.Is(err, context.Canceled) {
			t.Fatalf("Record error = %v, want context.Canceled", err)
		}

		canceledLatest, cancelLatest := context.WithCancel(context.Background())
		cancelLatest()

		if _, err := repository.Latest(canceledLatest); !errors.Is(err, context.Canceled) {
			t.Fatalf("Latest error = %v, want context.Canceled", err)
		}
	})
}

func TestMetricsInventoryParityAdditional(t *testing.T) {
	t.Run("empty metric summaries are zero", func(t *testing.T) {
		repository := NewMetricsRepository()

		if got := repository.Total(); got != (MetricSummary{}) {
			t.Fatalf("total metrics = %#v, want zero value", got)
		}

		if got := repository.Job("missing"); got != (MetricSummary{}) {
			t.Fatalf("job metrics = %#v, want zero value", got)
		}

		if got := repository.Queue("missing"); got != (MetricSummary{}) {
			t.Fatalf("queue metrics = %#v, want zero value", got)
		}
	})

	t.Run("jobs processed per minute since last snapshot is zero without new jobs", func(t *testing.T) {
		repository := NewMetricsRepository()
		repository.RecordJob(JobMeasurement{Job: "SendEmail", Queue: "mail", Runtime: 100 * time.Millisecond})

		snapshot := repository.SnapshotPerformance(time.Unix(100, 0).UTC())

		if got := repository.JobsProcessedPerMinuteSince(snapshot, time.Unix(100, 0).UTC()); got != 0 {
			t.Fatalf("jobs per minute = %v, want 0", got)
		}

		if got := repository.JobsProcessedPerMinuteSince(snapshot, time.Unix(160, 0).UTC()); got != 0 {
			t.Fatalf("jobs per minute = %v, want 0", got)
		}
	})
}

func TestMonitoringInventoryParityAdditional(t *testing.T) {
	t.Run("completed jobs are only stored when a monitored tag is present", func(t *testing.T) {
		repository := NewMonitoringRepository()
		repository.Monitor([]string{"billing"})

		if repository.RecordCompletedJob(CompletedJob{ID: "job-1", Tags: []string{"mail"}}) {
			t.Fatal("expected unmonitored job not to be retained")
		}

		if !repository.RecordCompletedJob(CompletedJob{ID: "job-2", Tags: []string{"billing", "mail"}}) {
			t.Fatal("expected monitored job to be retained")
		}

		if !repository.RecordCompletedJob(CompletedJob{ID: "job-3", Tags: []string{"billing"}}) {
			t.Fatal("expected second monitored job to be retained")
		}

		if got := repository.CompletedJobs("mail"); len(got) != 1 || got[0].ID != "job-2" {
			t.Fatalf("completed mail jobs = %#v", got)
		}

		if got := repository.CompletedJobs(""); len(got) != 2 || got[0].ID != "job-2" || got[1].ID != "job-3" {
			t.Fatalf("completed jobs = %#v", got)
		}

		if got := repository.CompletedJobs("invoice"); len(got) != 0 {
			t.Fatalf("completed invoice jobs = %#v, want empty", got)
		}
	})

	t.Run("MonitoringControllerTest::test_monitored_jobs_can_be_paginated_by_tag", func(t *testing.T) {
		repository := NewMonitoringRepository()
		repository.Monitor([]string{"billing"})
		repository.RecordCompletedJob(CompletedJob{ID: "job-1", Tags: []string{"billing"}})
		repository.RecordCompletedJob(CompletedJob{ID: "job-2", Tags: []string{"billing", "mail"}})

		if got, want := repository.CompletedJobsPage("billing", 1, 1), []CompletedJob{{ID: "job-2", Tags: []string{"billing", "mail"}}}; !reflect.DeepEqual(got, want) {
			t.Fatalf("paged completed jobs = %#v, want %#v", got, want)
		}
	})

	t.Run("MonitoringControllerTest::test_can_paginate_where_jobs_dont_exist", func(t *testing.T) {
		repository := NewMonitoringRepository()
		repository.Monitor([]string{"billing"})
		repository.RecordCompletedJob(CompletedJob{ID: "job-1", Tags: []string{"billing"}})

		if got := repository.CompletedJobsPage("billing", 1, 1); got != nil {
			t.Fatalf("paged completed jobs = %#v, want nil", got)
		}
	})
}

func TestWaitTimeCalculatorInventoryParityAdditional(t *testing.T) {
	snapshot := Snapshot{Queues: []QueueStatus{
		{Name: "redis:default", Pending: 3, Throughput: 1, Wait: 5 * time.Minute},
		{Name: "redis:mail", Pending: 5, Throughput: 5},
	}}

	t.Run("time to clear uses the recorded wait when present", func(t *testing.T) {
		times := WaitTimes(snapshot, "redis:default")

		if len(times) != 1 || times[0].Duration != 5*time.Minute {
			t.Fatalf("times = %#v", times)
		}
	})

	t.Run("time to clear filters out unknown queues", func(t *testing.T) {
		times := WaitTimes(snapshot, "redis:missing")

		if len(times) != 0 {
			t.Fatalf("times = %#v, want empty", times)
		}
	})
}

func TestJobRepositoryInventoryParityAdditional(t *testing.T) {
	now := time.Unix(1_000, 0).UTC()

	t.Run("job repository uses the current time when no clock is provided", func(t *testing.T) {
		fallback := NewJobRepository(nil)

		if fallback == nil || fallback.now == nil {
			t.Fatal("expected fallback repository clock to be initialized")
		}
	})

	t.Run("pending jobs default to the pending status and push time", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "job-default", Queue: "default"})

		pending := repository.Pending("default", 0, 10)

		if len(pending) != 1 {
			t.Fatalf("pending jobs = %#v", pending)
		}

		if pending[0].Status != JobPending {
			t.Fatalf("status = %q, want %q", pending[0].Status, JobPending)
		}

		if pending[0].PushedAt.IsZero() {
			t.Fatal("expected pushed_at to be set")
		}
	})

	t.Run("pending jobs are sorted by id when they share the same pushed at", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "job-b", Queue: "default", PushedAt: now.Add(time.Second)})
		repository.StorePending(JobRecord{ID: "job-a", Queue: "default", PushedAt: now.Add(time.Second)})

		pending := repository.Pending("default", 0, 10)

		if got, want := pending[0].ID, "job-a"; got != want {
			t.Fatalf("first pending job = %q, want %q", got, want)
		}
	})

	t.Run("recent jobs can be limited to the latest records", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "job-recent-1", Queue: "default", PushedAt: now.Add(2 * time.Second)})
		repository.StorePending(JobRecord{ID: "job-recent-2", Queue: "default", PushedAt: now.Add(3 * time.Second)})
		repository.MarkComplete("job-recent-1", true, now.Add(4*time.Second))
		repository.MarkComplete("job-recent-2", true, now.Add(5*time.Second))

		recent := repository.Recent(1)

		if len(recent) != 1 || recent[0].ID != "job-recent-2" {
			t.Fatalf("recent jobs = %#v", recent)
		}
	})

	t.Run("recent jobs can be trimmed with a negative limit", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "job-recent-1", Queue: "default", PushedAt: now.Add(2 * time.Second)})
		repository.MarkComplete("job-recent-1", true, now.Add(4*time.Second))

		repository.TrimRecent(-1)

		if got := repository.Recent(10); len(got) != 0 {
			t.Fatalf("recent jobs = %#v, want empty", got)
		}
	})

	t.Run("mark reserved returns false when the job does not exist", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })

		if repository.MarkReserved("missing", now) {
			t.Fatal("expected missing job not to be reserved")
		}
	})

	t.Run("mark complete returns false when the job does not exist", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })

		if repository.MarkComplete("missing", true, now) {
			t.Fatal("expected missing job not to be completed")
		}
	})

	t.Run("mark failed returns false when the job does not exist", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })

		if repository.MarkFailed("missing", now, time.Hour) {
			t.Fatal("expected missing job not to be failed")
		}
	})

	t.Run("release returns false when the job does not exist", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })

		if repository.Release("missing", now.Add(time.Minute)) {
			t.Fatal("expected missing job not to be released")
		}
	})

	t.Run("QueueProcessingTest::test_pending_delayed_jobs_are_stored_in_pending_job_database", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "job-delayed", Queue: "default", PushedAt: now.Add(6 * time.Second)})

		if !repository.Release("job-delayed", now.Add(10*time.Minute)) {
			t.Fatal("expected delayed job to be released")
		}

		if got := repository.Pending("default", 0, 10); len(got) != 0 {
			t.Fatalf("pending jobs = %#v, want delayed job to be hidden", got)
		}

		if got := repository.MigrateReleased(now.Add(11 * time.Minute)); got != 1 {
			t.Fatalf("migrated jobs = %d, want 1", got)
		}

		if got := repository.Pending("default", 0, 10); len(got) != 1 || got[0].ID != "job-delayed" {
			t.Fatalf("pending jobs after migrate = %#v", got)
		}
	})

	t.Run("pagination of job ids returns nil when the offset is out of range", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "tagged-job", Queue: "default", Tags: []string{"billing"}, PushedAt: now.Add(6 * time.Second)})

		if got := repository.JobIDsForTag("billing", 1, 10); got != nil {
			t.Fatalf("tagged jobs = %#v, want nil", got)
		}
	})

	t.Run("pagination of job ids supports non zero offsets", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "tagged-job-1", Queue: "default", Tags: []string{"billing"}, PushedAt: now.Add(6 * time.Second)})
		repository.StorePending(JobRecord{ID: "tagged-job-2", Queue: "default", Tags: []string{"billing"}, PushedAt: now.Add(7 * time.Second)})

		if got, want := repository.JobIDsForTag("billing", 1, 1), []string{"tagged-job-2"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("tagged jobs = %#v, want %#v", got, want)
		}
	})

	t.Run("pagination of job ids can be limited", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "tagged-job-1", Queue: "default", Tags: []string{"billing"}, PushedAt: now.Add(6 * time.Second)})
		repository.StorePending(JobRecord{ID: "tagged-job-2", Queue: "default", Tags: []string{"billing"}, PushedAt: now.Add(7 * time.Second)})

		if got, want := repository.JobIDsForTag("billing", 0, 1), []string{"tagged-job-1"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("tagged jobs = %#v, want %#v", got, want)
		}
	})

	t.Run("purging a queue that has no recent jobs leaves the repository unchanged", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "job-recent", Queue: "default", PushedAt: now.Add(2 * time.Second)})
		repository.MarkComplete("job-recent", true, now.Add(3*time.Second))

		if got := repository.PurgeQueue("missing"); got != 0 {
			t.Fatalf("purged count = %d, want 0", got)
		}

		if got := repository.Recent(10); len(got) != 1 {
			t.Fatalf("recent jobs = %#v, want kept", got)
		}
	})

	t.Run("FailedJobTest::test_temporary_failed_job_should_be_deleted_when_the_main_job_is_deleted", func(t *testing.T) {
		repository := NewJobRepository(func() time.Time { return now })
		repository.StorePending(JobRecord{ID: "job-failed", Queue: "default", Tags: []string{"failed-tag"}, PushedAt: now.Add(7 * time.Second)})

		if !repository.MarkFailed("job-failed", now.Add(8*time.Second), time.Hour) {
			t.Fatal("expected failed job to be stored")
		}

		if got, ok := repository.FindFailed("job-failed"); !ok || !reflect.DeepEqual(got.Tags, []string{"failed-tag"}) {
			t.Fatalf("failed job = %#v, ok=%v", got, ok)
		}

		if !repository.DeleteFailed("job-failed") {
			t.Fatal("expected failed job to be deleted")
		}
	})
}

func TestRedisJobRepositoryInventoryParityAdditional(t *testing.T) {
	t.Run("RedisJobRepositoryTest::test_it_saves_microseconds_as_a_float_and_disregards_the_locale", func(t *testing.T) {
		at := time.Unix(1, 234_567_890).UTC()

		if got, want := microsecondsFloat(at), 1_234_567.89; got != want {
			t.Fatalf("microseconds float = %v, want %v", got, want)
		}

		if got, want := strconv.FormatFloat(microsecondsFloat(at), 'f', 2, 64), "1234567.89"; got != want {
			t.Fatalf("formatted microseconds = %q, want %q", got, want)
		}
	})
}
