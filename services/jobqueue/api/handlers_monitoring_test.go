package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/jobqueue"
	"github.com/bedrock/services/jobqueue/api"
)

// Port of MonitoringControllerTest::test_monitored_tags_and_job_counts_are_returned.
func TestMonitoringIndexReturnsMonitoredTagsWithJobCounts(t *testing.T) {
	t.Parallel()

	monitoring := jobqueue.NewMonitoringRepository()
	monitoring.Monitor([]string{"first", "second"})

	jobs := jobqueue.NewJobRepository(nil)
	jobs.StorePending(jobqueue.JobRecord{ID: "1", Queue: "default", Tags: []string{"first"}})
	jobs.StorePending(jobqueue.JobRecord{ID: "2", Queue: "default", Tags: []string{"first"}})
	jobs.StorePending(jobqueue.JobRecord{ID: "3", Queue: "default", Tags: []string{"second"}})

	handler := newTestHandler(t, api.Options{Monitoring: monitoring, Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/monitoring", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(body) != 2 {
		t.Fatalf("expected 2 monitored tags, got %d", len(body))
	}

	counts := map[string]float64{}

	for _, entry := range body {
		counts[entry["tag"].(string)] = entry["count"].(float64)
	}

	if counts["first"] != 2 || counts["second"] != 1 {
		t.Fatalf("unexpected counts: %v", counts)
	}
}

// Port of MonitoringControllerTest::test_monitored_jobs_can_be_paginated_by_tag.
func TestMonitoringShowPaginatesJobsByTag(t *testing.T) {
	t.Parallel()

	monitoring := jobqueue.NewMonitoringRepository()
	monitoring.Monitor([]string{"billing"})

	jobs := jobqueue.NewJobRepository(nil)

	for _, id := range []string{"1", "2", "3", "4", "5"} {
		jobs.StorePending(jobqueue.JobRecord{ID: id, Queue: "default", Tags: []string{"billing"}})
	}

	handler := newTestHandler(t, api.Options{Monitoring: monitoring, Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/monitoring/billing?starting_at=0&limit=2", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body struct {
		Tag  string   `json:"tag"`
		Jobs []string `json:"jobs"`
	}

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body.Tag != "billing" || len(body.Jobs) != 2 {
		t.Fatalf("got %+v, want tag=billing with 2 jobs", body)
	}
}

// Port of MonitoringControllerTest::test_can_paginate_where_jobs_dont_exist.
func TestMonitoringShowReturnsEmptyWhenNoJobsMatchTag(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{})

	r := httptest.NewRequest(http.MethodGet, "/api/monitoring/nobody", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body struct {
		Tag  string   `json:"tag"`
		Jobs []string `json:"jobs"`
	}

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body.Tag != "nobody" || len(body.Jobs) != 0 {
		t.Fatalf("got %+v, want empty jobs", body)
	}
}

// Port of MonitoringControllerTest::test_can_start_monitoring_tags.
func TestMonitoringStoreStartsMonitoringATag(t *testing.T) {
	t.Parallel()

	monitoring := jobqueue.NewMonitoringRepository()
	handler := newTestHandler(t, api.Options{Monitoring: monitoring})

	payload, _ := json.Marshal(map[string]string{"tag": "payments"})
	r := httptest.NewRequest(http.MethodPost, "/api/monitoring", bytes.NewReader(payload))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !monitoring.IsMonitoring([]string{"payments"}) {
		t.Fatalf("expected payments to be monitored")
	}
}

// Port of MonitoringControllerTest::test_can_stop_monitoring_tags.
func TestMonitoringDeleteStopsMonitoringATag(t *testing.T) {
	t.Parallel()

	monitoring := jobqueue.NewMonitoringRepository()
	monitoring.Monitor([]string{"payments"})

	handler := newTestHandler(t, api.Options{Monitoring: monitoring})

	r := httptest.NewRequest(http.MethodDelete, "/api/monitoring/payments", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	if monitoring.IsMonitoring([]string{"payments"}) {
		t.Fatalf("expected payments to be unmonitored")
	}
}
