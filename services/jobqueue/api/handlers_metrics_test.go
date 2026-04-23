package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/jobqueue"
	"github.com/bedrock/services/jobqueue/api"
)

// Port of MetricsTest::test_throughput_is_stored_per_job_class.
// Port of MetricsTest::test_average_runtime_is_stored_per_job_class_in_milliseconds.
// Port of MetricsTest::test_list_of_all_jobs_with_metric_information_is_maintained.
func TestMetricsJobsEndpointReturnsThroughputAndRuntimePerJobClass(t *testing.T) {
	t.Parallel()

	metrics := jobqueue.NewMetricsRepository()
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "SendWelcome", Queue: "default", Runtime: 200 * time.Millisecond})
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "SendWelcome", Queue: "default", Runtime: 400 * time.Millisecond})
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "RebuildIndex", Queue: "mail", Runtime: 100 * time.Millisecond})

	handler := newTestHandler(t, api.Options{Metrics: metrics})

	r := httptest.NewRequest(http.MethodGet, "/api/metrics/jobs", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(body) != 2 {
		t.Fatalf("expected 2 job classes, got %d", len(body))
	}

	summaries := map[string]map[string]any{}

	for _, entry := range body {
		summaries[entry["name"].(string)] = entry
	}

	welcome, ok := summaries["SendWelcome"]

	if !ok {
		t.Fatalf("expected SendWelcome in response, got %v", summaries)
	}

	if welcome["throughput"].(float64) != 2 {
		t.Fatalf("expected SendWelcome throughput 2, got %v", welcome["throughput"])
	}

	if welcome["averageRuntime"].(float64) != 300 {
		t.Fatalf("expected SendWelcome averageRuntime 300ms, got %v", welcome["averageRuntime"])
	}
}

// Port of MetricsTest::test_throughput_is_stored_per_queue.
// Port of MetricsTest::test_average_runtime_is_stored_per_queue_in_milliseconds.
func TestMetricsQueuesEndpointReturnsThroughputAndRuntimePerQueue(t *testing.T) {
	t.Parallel()

	metrics := jobqueue.NewMetricsRepository()
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "SendWelcome", Queue: "default", Runtime: 200 * time.Millisecond})
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "RebuildIndex", Queue: "default", Runtime: 400 * time.Millisecond})

	snapshot := metrics.SnapshotPerformance(time.Unix(1700000000, 0))

	if snapshot.JobsProcessed != 2 {
		t.Fatalf("expected snapshot to capture 2 jobs, got %d", snapshot.JobsProcessed)
	}

	handler := newTestHandler(t, api.Options{Metrics: metrics})

	r := httptest.NewRequest(http.MethodGet, "/api/metrics/queues", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if len(body) != 1 || body[0]["name"] != "default" {
		t.Fatalf("expected only default queue, got %+v", body)
	}

	if body[0]["throughput"].(float64) != 2 {
		t.Fatalf("expected default queue throughput 2, got %v", body[0]["throughput"])
	}

	if body[0]["averageRuntime"].(float64) != 300 {
		t.Fatalf("expected default queue averageRuntime 300ms, got %v", body[0]["averageRuntime"])
	}
}

// Port of MetricsTest::test_total_throughput_is_stored.
// Port of MetricsTest::test_snapshot_of_metrics_performance_can_be_stored.
func TestMetricsSnapshotEndpointRecordsSnapshot(t *testing.T) {
	t.Parallel()

	metrics := jobqueue.NewMetricsRepository()
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "SendWelcome", Queue: "default", Runtime: 200 * time.Millisecond})
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "SendWelcome", Queue: "default", Runtime: 400 * time.Millisecond})

	fixedNow := time.Unix(1700000000, 0)

	handler := newTestHandler(t, api.Options{
		Metrics: metrics,
		Now:     func() time.Time { return fixedNow },
	})

	r := httptest.NewRequest(http.MethodPost, "/api/metrics/snapshot", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var body map[string]any

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body["jobsProcessed"].(float64) != 2 {
		t.Fatalf("expected jobsProcessed 2, got %v", body["jobsProcessed"])
	}

	if len(metrics.Snapshots()) != 1 {
		t.Fatalf("expected snapshot to be retained, got %d", len(metrics.Snapshots()))
	}
}

// Port of MetricsTest::test_jobs_processed_per_minute_since_last_snapshot_is_calculable.
// Port of MetricsTest::test_only_past_24_snapshots_are_retained.
func TestMetricsSnapshotEndpointBounds24Retention(t *testing.T) {
	t.Parallel()

	metrics := jobqueue.NewMetricsRepository()

	for i := 0; i < 30; i++ {
		metrics.RecordJob(jobqueue.JobMeasurement{Job: "SendWelcome", Queue: "default", Runtime: time.Millisecond})
		metrics.SnapshotPerformance(time.Unix(1700000000+int64(i)*60, 0))
	}

	if got := len(metrics.Snapshots()); got != 24 {
		t.Fatalf("expected 24 retained snapshots, got %d", got)
	}
}

// Port of MetricsTest::test_total_throughput_is_stored (aggregate-level sanity check via Total()).
func TestMetricsTotalThroughputIsAggregatedAcrossRecordJob(t *testing.T) {
	t.Parallel()

	metrics := jobqueue.NewMetricsRepository()
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "SendWelcome", Queue: "default", Runtime: 200 * time.Millisecond})
	metrics.RecordJob(jobqueue.JobMeasurement{Job: "RebuildIndex", Queue: "mail", Runtime: 400 * time.Millisecond})

	total := metrics.Total()

	if total.Throughput != 2 {
		t.Fatalf("expected aggregate throughput 2, got %d", total.Throughput)
	}

	if total.AverageRuntime != 300*time.Millisecond {
		t.Fatalf("expected aggregate avg 300ms, got %v", total.AverageRuntime)
	}
}
