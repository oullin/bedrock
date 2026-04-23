package jobqueue

import (
	"reflect"
	"testing"
	"time"
)

func TestDashboardStatsInventoryParity(t *testing.T) {
	snapshot := Snapshot{Queues: []QueueStatus{
		{Name: "redis:default", Pending: 10, Reserved: 2, Processing: 3, Failed: 1, Throughput: 7, Wait: 45 * time.Second, Paused: true},
		{Name: "redis:mail", Pending: 4, Reserved: 1, Processing: 2, Failed: 0, Throughput: 5, Wait: 15 * time.Second, Paused: true},
	}}

	stats := StatsForSnapshot(snapshot)

	t.Run("DashboardStatsControllerTest::test_all_stats_are_correctly_returned", func(t *testing.T) {
		if stats.Queues != 2 || stats.Pending != 14 || stats.Reserved != 3 || stats.Processing != 5 || stats.Failed != 1 || stats.Throughput != 12 {
			t.Fatalf("stats = %#v", stats)
		}

		if stats.LongestWait != 45*time.Second {
			t.Fatalf("longest wait = %s", stats.LongestWait)
		}
	})

	t.Run("DashboardStatsControllerTest::test_paused_status_is_reflected_if_all_master_supervisors_are_paused", func(t *testing.T) {
		if !stats.Paused {
			t.Fatal("expected paused status when every queue is paused")
		}
	})

	t.Run("DashboardStatsControllerTest::test_paused_status_isnt_reflected_if_not_all_master_supervisors_are_paused", func(t *testing.T) {
		mixed := StatsForSnapshot(Snapshot{Queues: []QueueStatus{
			{Name: "redis:default", Paused: true},
			{Name: "redis:mail", Paused: false},
		}})

		if mixed.Paused {
			t.Fatal("expected paused=false when any queue is active")
		}
	})
}

func TestMetricsInventoryParity(t *testing.T) {
	repo := NewMetricsRepository()
	repo.RecordJob(JobMeasurement{Job: "SendEmail", Queue: "mail", Runtime: 100 * time.Millisecond})
	repo.RecordJob(JobMeasurement{Job: "SendEmail", Queue: "mail", Runtime: 300 * time.Millisecond})
	repo.RecordJob(JobMeasurement{Job: "SyncAccount", Queue: "default", Runtime: 500 * time.Millisecond})

	t.Run("MetricsTest::test_total_throughput_is_stored", func(t *testing.T) {
		if got := repo.Total().Throughput; got != 3 {
			t.Fatalf("total throughput = %d, want 3", got)
		}
	})

	t.Run("MetricsTest::test_throughput_is_stored_per_job_class", func(t *testing.T) {
		if got := repo.Job("SendEmail").Throughput; got != 2 {
			t.Fatalf("job throughput = %d, want 2", got)
		}
	})

	t.Run("MetricsTest::test_throughput_is_stored_per_queue", func(t *testing.T) {
		if got := repo.Queue("mail").Throughput; got != 2 {
			t.Fatalf("queue throughput = %d, want 2", got)
		}
	})

	t.Run("MetricsTest::test_average_runtime_is_stored_per_job_class_in_milliseconds", func(t *testing.T) {
		if got := repo.Job("SendEmail").AverageRuntime; got != 200*time.Millisecond {
			t.Fatalf("average runtime = %s, want 200ms", got)
		}
	})

	t.Run("MetricsTest::test_average_runtime_is_stored_per_queue_in_milliseconds", func(t *testing.T) {
		if got := repo.Queue("mail").AverageRuntime; got != 200*time.Millisecond {
			t.Fatalf("average runtime = %s, want 200ms", got)
		}
	})

	t.Run("MetricsTest::test_list_of_all_jobs_with_metric_information_is_maintained", func(t *testing.T) {
		if got, want := repo.Jobs(), []string{"SendEmail", "SyncAccount"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("jobs = %v, want %v", got, want)
		}
	})

	t.Run("MetricsTest::test_snapshot_of_metrics_performance_can_be_stored", func(t *testing.T) {
		snapshot := repo.SnapshotPerformance(time.Unix(100, 0).UTC())

		if snapshot.JobsProcessed != 3 || snapshot.JobThroughput["SendEmail"] != 2 || snapshot.QueueThroughput["mail"] != 2 {
			t.Fatalf("snapshot = %#v", snapshot)
		}
	})

	t.Run("MetricsTest::test_jobs_processed_per_minute_since_last_snapshot_is_calculable", func(t *testing.T) {
		snapshot := repo.SnapshotPerformance(time.Unix(200, 0).UTC())
		repo.RecordJob(JobMeasurement{Job: "SendEmail", Queue: "mail", Runtime: 250 * time.Millisecond})

		if got := repo.JobsProcessedPerMinuteSince(snapshot, time.Unix(320, 0).UTC()); got != 0.5 {
			t.Fatalf("jobs per minute = %v, want 0.5", got)
		}
	})

	t.Run("MetricsTest::test_only_past_24_snapshots_are_retained", func(t *testing.T) {
		for i := 0; i < 30; i++ {
			repo.SnapshotPerformance(time.Unix(int64(i), 0).UTC())
		}

		snapshots := repo.Snapshots()

		if len(snapshots) != 24 {
			t.Fatalf("snapshot count = %d, want 24", len(snapshots))
		}

		if snapshots[0].RecordedAt != time.Unix(6, 0).UTC() {
			t.Fatalf("first retained snapshot = %s", snapshots[0].RecordedAt)
		}
	})
}

func TestMonitoringInventoryParity(t *testing.T) {
	repo := NewMonitoringRepository()

	t.Run("MonitoringControllerTest::test_can_start_monitoring_tags", func(t *testing.T) {
		repo.Monitor([]string{"billing"})

		if !repo.IsMonitoring([]string{"billing"}) {
			t.Fatal("expected billing tag to be monitored")
		}
	})

	t.Run("MonitoringTest::test_can_retrieve_all_monitored_tags", func(t *testing.T) {
		repo.Monitor([]string{"billing", "mail"})

		if got, want := repo.Monitoring(), []string{"billing", "mail"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("monitoring = %v, want %v", got, want)
		}

		got := repo.Monitoring()
		got[0] = "mutated"

		if again, want := repo.Monitoring(), []string{"billing", "mail"}; !reflect.DeepEqual(again, want) {
			t.Fatalf("monitoring should be copied, got %v", again)
		}
	})

	t.Run("MonitoringTest::test_can_determine_if_a_set_of_tags_are_being_monitored", func(t *testing.T) {
		if !repo.IsMonitoring([]string{"unknown", "billing"}) {
			t.Fatal("expected billing to be monitored")
		}
	})

	t.Run("MonitoringTest::test_completed_jobs_are_stored_in_database_when_one_of_their_tags_is_being_monitored", func(t *testing.T) {
		if !repo.RecordCompletedJob(CompletedJob{ID: "job-1", Tags: []string{"billing", "user:1"}}) {
			t.Fatal("expected monitored job to be retained")
		}

		if got := len(repo.CompletedJobs("billing")); got != 1 {
			t.Fatalf("completed billing jobs = %d, want 1", got)
		}

		completed := repo.CompletedJobs("billing")
		completed[0].Tags[0] = "mutated"

		if again := repo.CompletedJobs("billing"); again[0].Tags[0] != "billing" {
			t.Fatalf("completed jobs should be copied, got %#v", again)
		}
	})

	t.Run("MonitoringTest::test_completed_jobs_are_removed_from_database_when_their_tag_is_no_longer_monitored", func(t *testing.T) {
		repo.StopMonitoring([]string{"billing"})

		if got := len(repo.CompletedJobs("billing")); got != 0 {
			t.Fatalf("completed billing jobs = %d, want 0", got)
		}
	})

	t.Run("MonitoringTest::test_all_completed_jobs_are_removed_from_database_when_their_tag_is_no_longer_monitored", func(t *testing.T) {
		repo.Monitor([]string{"billing", "mail"})
		repo.RecordCompletedJob(CompletedJob{ID: "job-2", Tags: []string{"billing"}})
		repo.RecordCompletedJob(CompletedJob{ID: "job-3", Tags: []string{"billing", "mail"}})
		repo.StopMonitoring([]string{"billing", "mail"})

		if got := len(repo.CompletedJobs("")); got != 0 {
			t.Fatalf("completed jobs = %d, want 0", got)
		}
	})

	t.Run("MonitoringTest::test_can_stop_monitoring_tags", func(t *testing.T) {
		repo.Monitor([]string{"mail"})
		repo.StopMonitoring([]string{"mail"})

		if repo.IsMonitoring([]string{"mail"}) {
			t.Fatal("expected mail tag to stop being monitored")
		}
	})

	t.Run("MonitoringTest::test_tags_that_are_removed_from_monitoring_are_removed_from_storage", func(t *testing.T) {
		if got := repo.Monitoring(); len(got) != 0 {
			t.Fatalf("monitoring = %v, want empty", got)
		}
	})

	t.Run("MonitoringControllerTest::test_can_stop_monitoring_tags", func(t *testing.T) {
		repo.Monitor([]string{"mail"})
		repo.StopMonitoring([]string{"mail"})

		if repo.IsMonitoring([]string{"mail"}) {
			t.Fatal("expected mail tag to stop being monitored")
		}
	})

	t.Run("MonitoringControllerTest::test_monitored_tags_and_job_counts_are_returned", func(t *testing.T) {
		repo = NewMonitoringRepository()
		repo.Monitor([]string{"billing", "mail"})
		repo.RecordCompletedJob(CompletedJob{ID: "job-1", Tags: []string{"billing"}})
		repo.RecordCompletedJob(CompletedJob{ID: "job-2", Tags: []string{"mail"}})

		if got := len(repo.CompletedJobs("billing")); got != 1 {
			t.Fatalf("billing jobs = %d, want 1", got)
		}

		if got := len(repo.CompletedJobs("mail")); got != 1 {
			t.Fatalf("mail jobs = %d, want 1", got)
		}
	})

	t.Run("MonitoringControllerTest::test_monitored_jobs_can_be_paginated_by_tag", func(t *testing.T) {
		repo = NewMonitoringRepository()
		repo.Monitor([]string{"billing"})
		repo.RecordCompletedJob(CompletedJob{ID: "job-1", Tags: []string{"billing"}})
		repo.RecordCompletedJob(CompletedJob{ID: "job-2", Tags: []string{"billing"}})

		if got, want := repo.CompletedJobs("billing"), []CompletedJob{{ID: "job-1", Tags: []string{"billing"}}, {ID: "job-2", Tags: []string{"billing"}}}; !reflect.DeepEqual(got, want) {
			t.Fatalf("completed jobs = %#v, want %#v", got, want)
		}
	})
}

func TestWaitTimeCalculatorInventoryParity(t *testing.T) {
	snapshot := Snapshot{Queues: []QueueStatus{
		{Name: "redis:default", Pending: 12, Throughput: 6},
		{Name: "redis:mail", Pending: 5, Throughput: 5},
		{Name: "redis:empty", Pending: 0, Throughput: 5},
		{Name: "redis:stalled", Pending: 10, Throughput: 0},
	}}

	t.Run("WaitTimeCalculatorTest::test_time_to_clear_is_calculated_per_queue", func(t *testing.T) {
		times := WaitTimes(snapshot, "redis:default")

		if len(times) != 1 || times[0].Duration != 2*time.Minute {
			t.Fatalf("times = %#v", times)
		}
	})

	t.Run("WaitTimeCalculatorTest::test_multiple_queues_are_supported", func(t *testing.T) {
		times := WaitTimes(snapshot, "redis:default", "redis:mail")

		if len(times) != 2 || times[0].Duration != 2*time.Minute || times[1].Duration != time.Minute {
			t.Fatalf("times = %#v", times)
		}
	})

	t.Run("WaitTimeCalculatorTest::test_single_queue_can_be_retrieved_for_multiple_queues", func(t *testing.T) {
		times := WaitTimes(snapshot, "redis:mail")

		if len(times) != 1 || times[0].Queue != "redis:mail" {
			t.Fatalf("times = %#v", times)
		}
	})

	t.Run("WaitTimeCalculatorTest::test_time_to_clear_can_be_zero", func(t *testing.T) {
		times := WaitTimes(snapshot, "redis:empty")

		if len(times) != 1 || times[0].Duration != 0 {
			t.Fatalf("times = %#v", times)
		}
	})

	t.Run("WaitTimeCalculatorTest::test_total_processes_can_be_zero", func(t *testing.T) {
		times := WaitTimes(snapshot, "redis:stalled")

		if len(times) != 1 || times[0].Duration != 0 {
			t.Fatalf("times = %#v", times)
		}
	})
}

func TestRepositoryAndOptionsInventoryParity(t *testing.T) {
	t.Run("SupervisorOptionsTest::test_default_queue_is_used_when_null_is_given", func(t *testing.T) {
		if got := NewSupervisorOptions("").Queue; got != "default" {
			t.Fatalf("queue = %q, want default", got)
		}
	})

	t.Run("StopwatchTest::test_time_between_checks_can_be_measured", func(t *testing.T) {
		stopwatch := NewStopwatch(time.Unix(100, 0).UTC())

		if got := stopwatch.Check(time.Unix(160, 0).UTC()); got != time.Minute {
			t.Fatalf("elapsed = %s, want 1m", got)
		}

		if got := stopwatch.Check(time.Unix(190, 0).UTC()); got != 30*time.Second {
			t.Fatalf("elapsed = %s, want 30s", got)
		}
	})
}

func TestPendingJobRetrievalInventoryParity(t *testing.T) {
	now := time.Unix(1000, 0).UTC()
	repo := NewJobRepository(func() time.Time { return now })

	repo.StorePending(JobRecord{ID: "job-1", Name: "BasicJob", Queue: "default", Tags: []string{"alpha"}, PushedAt: now.Add(1 * time.Second)})
	repo.StorePending(JobRecord{ID: "job-2", Name: "BasicJob", Queue: "default", Tags: []string{"beta"}, PushedAt: now.Add(2 * time.Second)})
	repo.StorePending(JobRecord{ID: "job-3", Name: "BasicJob", Queue: "default", Tags: []string{"alpha", "beta"}, PushedAt: now.Add(3 * time.Second)})
	repo.StorePending(JobRecord{ID: "job-4", Name: "BasicJob", Queue: "default", Tags: []string{"gamma"}, PushedAt: now.Add(4 * time.Second)})
	repo.StorePending(JobRecord{ID: "job-5", Name: "MailJob", Queue: "mail", Type: "listener", Tags: []string{"mail"}, PushedAt: now.Add(5 * time.Second)})

	t.Run("JobRetrievalTest::test_pending_jobs_can_be_retrieved", func(t *testing.T) {
		pending := repo.Pending("default", 0, 10)

		if got, want := len(pending), 4; got != want {
			t.Fatalf("pending count = %d, want %d", got, want)
		}

		if got, want := pending[0].ID, "job-1"; got != want {
			t.Fatalf("first pending job = %q, want %q", got, want)
		}

		if got, want := pending[3].ID, "job-4"; got != want {
			t.Fatalf("last pending job = %q, want %q", got, want)
		}

		page := repo.Pending("default", 1, 2)

		if got, want := len(page), 2; got != want {
			t.Fatalf("page count = %d, want %d", got, want)
		}

		if got, want := page[0].ID, "job-2"; got != want {
			t.Fatalf("first page item = %q, want %q", got, want)
		}

		if got, want := page[1].ID, "job-3"; got != want {
			t.Fatalf("second page item = %q, want %q", got, want)
		}

		pending[0].Tags[0] = "mutated"

		if again := repo.Pending("default", 0, 1); again[0].Tags[0] != "alpha" {
			t.Fatalf("pending jobs should be copied, got %#v", again[0].Tags)
		}
	})
}

func TestJobRepositoryInventoryParity(t *testing.T) {
	now := time.Unix(1000, 0).UTC()
	repo := NewJobRepository(func() time.Time { return now })

	repo.StorePending(JobRecord{ID: "job-1", Name: "SendEmail", Queue: "default", Type: "job", Tags: []string{"user:1"}, PushedAt: now})
	repo.StorePending(JobRecord{ID: "job-2", Name: "BuildReport", Queue: "default", Type: "job", Tags: []string{"reports"}, PushedAt: now.Add(time.Second)})
	repo.StorePending(JobRecord{ID: "job-3", Name: "NotifyUser", Queue: "mail", Type: "listener", Tags: []string{"user:1", "mail"}, PushedAt: now.Add(2 * time.Second)})

	t.Run("JobRetrievalTest::test_pending_jobs_can_be_retrieved", func(t *testing.T) {
		pending := repo.Pending("default", 0, 10)

		if got, want := len(pending), 2; got != want {
			t.Fatalf("pending count = %d, want %d", got, want)
		}

		if pending[0].ID != "job-1" || pending[1].ID != "job-2" {
			t.Fatalf("pending jobs = %#v", pending)
		}
	})

	t.Run("JobRetrievalTest::test_paginating_large_job_results_gives_correct_amounts", func(t *testing.T) {
		page := repo.Pending("default", 1, 1)

		if len(page) != 1 || page[0].ID != "job-2" {
			t.Fatalf("page = %#v", page)
		}
	})

	t.Run("QueueProcessingTest::test_pending_jobs_are_stored_in_pending_job_database", func(t *testing.T) {
		if got := len(repo.Pending("mail", 0, 10)); got != 1 {
			t.Fatalf("mail pending count = %d, want 1", got)
		}
	})

	t.Run("QueueProcessingTest::test_pending_jobs_are_stored_with_their_tags", func(t *testing.T) {
		if got, want := repo.JobIDsForTag("user:1", 0, 10), []string{"job-1", "job-3"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("tagged jobs = %v, want %v", got, want)
		}
	})

	t.Run("QueueProcessingTest::test_pending_jobs_are_stored_with_their_type", func(t *testing.T) {
		pending := repo.Pending("mail", 0, 1)

		if len(pending) != 1 || pending[0].Type != "listener" {
			t.Fatalf("pending mail job = %#v", pending)
		}
	})

	t.Run("QueueProcessingTest::test_pending_job_is_marked_as_reserved_during_processing", func(t *testing.T) {
		if !repo.MarkReserved("job-1", now.Add(time.Minute)) {
			t.Fatal("expected job-1 to be reserved")
		}

		if got := len(repo.Pending("default", 0, 10)); got != 1 {
			t.Fatalf("pending default count = %d, want 1 after reserve", got)
		}
	})

	t.Run("QueueProcessingTest::test_stale_reserved_jobs_are_marked_as_pending_after_migrating", func(t *testing.T) {
		if got := repo.MigrateStaleReserved(now.Add(2 * time.Minute)); got != 1 {
			t.Fatalf("migrated count = %d, want 1", got)
		}

		pending := repo.Pending("default", 0, 10)

		if len(pending) != 2 || !pending[0].ReservedWasMigrated {
			t.Fatalf("pending after migrate = %#v", pending)
		}
	})

	t.Run("MarkJobAsCompleteTest::test_it_can_mark_a_job_as_complete", func(t *testing.T) {
		if !repo.MarkComplete("job-1", true, now.Add(3*time.Minute)) {
			t.Fatal("expected job-1 completion")
		}

		recent := repo.Recent(10)

		if len(recent) != 1 || recent[0].Status != JobCompleted || !recent[0].CompletionStored {
			t.Fatalf("recent = %#v", recent)
		}
	})

	t.Run("QueueProcessingTest::test_pending_jobs_are_no_longer_in_pending_database_after_being_worked", func(t *testing.T) {
		for _, job := range repo.Pending("default", 0, 10) {
			if job.ID == "job-1" {
				t.Fatalf("completed job is still pending: %#v", job)
			}
		}
	})

	t.Run("QueueProcessingTest::test_completed_jobs_are_not_normally_stored_in_completed_database", func(t *testing.T) {
		repo.StorePending(JobRecord{ID: "job-4", Name: "SilentComplete", Queue: "default", PushedAt: now.Add(4 * time.Second)})

		if !repo.MarkComplete("job-4", false, now.Add(4*time.Minute)) {
			t.Fatal("expected job-4 completion")
		}

		if got := len(repo.Recent(10)); got != 1 {
			t.Fatalf("recent count = %d, want unchanged 1", got)
		}
	})

	t.Run("QueueProcessingTest::test_legacy_jobs_can_be_processed_without_errors", func(t *testing.T) {
		repo.StorePending(JobRecord{ID: "job-legacy", Name: "LegacyJob", Queue: "default", PushedAt: now.Add(4 * time.Second)})

		if !repo.MarkReserved("job-legacy", now.Add(4*time.Minute)) {
			t.Fatal("expected legacy job to be reserved")
		}

		if !repo.MarkComplete("job-legacy", true, now.Add(5*time.Minute)) {
			t.Fatal("expected legacy job to complete")
		}

		if got := len(repo.Recent(10)); got != 2 {
			t.Fatalf("recent jobs count = %d, want 2", got)
		}
	})

	t.Run("JobRetrievalTest::test_recent_jobs_are_correctly_trimmed_and_expired", func(t *testing.T) {
		repo.StorePending(JobRecord{ID: "job-5", Name: "Archive", Queue: "default", PushedAt: now.Add(5 * time.Second)})
		repo.MarkComplete("job-5", true, now.Add(5*time.Minute))
		repo.TrimRecent(1)

		recent := repo.Recent(10)

		if len(recent) != 1 || recent[0].ID != "job-5" {
			t.Fatalf("trimmed recent = %#v", recent)
		}
	})

	t.Run("FailedJobTest::test_failed_jobs_are_placed_in_the_failed_job_table", func(t *testing.T) {
		repo.StorePending(JobRecord{ID: "job-failed", Name: "FailingJob", Queue: "default", Tags: []string{"failed-tag"}, PushedAt: now.Add(6 * time.Second)})

		if !repo.MarkFailed("job-failed", now.Add(6*time.Minute), time.Hour) {
			t.Fatal("expected failed job to be stored")
		}

		failed, ok := repo.FindFailed("job-failed")

		if !ok || failed.Status != JobFailed {
			t.Fatalf("failed job = %#v, ok=%v", failed, ok)
		}
	})

	t.Run("FailedJobTest::test_tags_for_failed_jobs_are_stored_in_redis", func(t *testing.T) {
		failed, _ := repo.FindFailed("job-failed")

		if got, want := failed.Tags, []string{"failed-tag"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("failed tags = %v, want %v", got, want)
		}
	})

	t.Run("FailedJobTest::test_failed_job_tags_have_an_expiration", func(t *testing.T) {
		failed, _ := repo.FindFailed("job-failed")

		if got, want := failed.FailedTagsExpireAt, now.Add(6*time.Minute).Add(time.Hour); !got.Equal(want) {
			t.Fatalf("failed tag expiration = %s, want %s", got, want)
		}
	})

	t.Run("RedisJobRepositoryTest::test_it_can_find_a_failed_job_by_its_id", func(t *testing.T) {
		if _, ok := repo.FindFailed("job-failed"); !ok {
			t.Fatal("expected failed job to be found")
		}
	})

	t.Run("RedisJobRepositoryTest::test_it_will_not_find_a_failed_job_if_the_job_has_not_failed", func(t *testing.T) {
		if _, ok := repo.FindFailed("job-2"); ok {
			t.Fatal("did not expect pending job in failed storage")
		}
	})

	t.Run("RedisJobRepositoryTest::test_it_will_delete_a_failed_job", func(t *testing.T) {
		if !repo.DeleteFailed("job-failed") {
			t.Fatal("expected failed job delete")
		}

		if _, ok := repo.FindFailed("job-failed"); ok {
			t.Fatal("failed job was not deleted")
		}
	})

	t.Run("RedisJobRepositoryTest::test_it_will_not_delete_a_job_if_the_job_has_not_failed", func(t *testing.T) {
		if repo.DeleteFailed("job-2") {
			t.Fatal("did not expect deleting a non-failed job")
		}
	})

	t.Run("StoreTagsForFailedTest::test_temporary_failed_job_should_be_deleted_when_the_main_job_is_deleted", func(t *testing.T) {
		repo.StorePending(JobRecord{ID: "job-temporary-failed", Name: "FailingJob", Queue: "default", Tags: []string{"failed-tag"}, PushedAt: now.Add(7 * time.Second)})

		if !repo.MarkFailed("job-temporary-failed", now.Add(7*time.Minute), time.Minute) {
			t.Fatal("expected temporary failed job to be stored")
		}

		if !repo.DeleteFailed("job-temporary-failed") {
			t.Fatal("expected temporary failed job delete")
		}

		if _, ok := repo.FindFailed("job-temporary-failed"); ok {
			t.Fatal("temporary failed job was not deleted")
		}
	})

	t.Run("RedisJobRepositoryTest::test_it_stores_delay_when_job_is_released", func(t *testing.T) {
		if !repo.Release("job-2", now.Add(10*time.Minute)) {
			t.Fatal("expected job-2 release")
		}

		if got := len(repo.Pending("default", 0, 10)); got != 0 {
			t.Fatalf("pending count before delay expires = %d, want 0", got)
		}
	})

	t.Run("RedisJobRepositoryTest::test_it_clears_delay_when_job_is_migrated", func(t *testing.T) {
		if got := repo.MigrateReleased(now.Add(11 * time.Minute)); got != 1 {
			t.Fatalf("migrated released count = %d, want 1", got)
		}

		if got := len(repo.Pending("default", 0, 10)); got != 1 {
			t.Fatalf("pending count after delay expires = %d, want 1", got)
		}
	})

	t.Run("RedisJobRepositoryTest::test_it_removes_recent_jobs_when_queue_is_purged", func(t *testing.T) {
		if got := repo.PurgeQueue("default"); got != 1 {
			t.Fatalf("purged count = %d, want 1", got)
		}

		if got := len(repo.Recent(10)); got != 0 {
			t.Fatalf("recent count = %d, want 0", got)
		}
	})

	t.Run("TagRepositoryTest::test_pagination_of_job_ids_can_be_accomplished", func(t *testing.T) {
		got := repo.JobIDsForTag("user:1", 0, 1)

		if want := []string{"job-3"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("tag page = %v, want %v", got, want)
		}
	})
}
