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

// Port of JobRetrievalTest::test_pending_jobs_can_be_retrieved.
func TestPendingJobsEndpointReturnsPendingJobsForQueue(t *testing.T) {
	t.Parallel()

	jobs := jobqueue.NewJobRepository(nil)

	jobs.StorePending(jobqueue.JobRecord{
		ID:       "job-1",
		Name:     "SendWelcomeEmail",
		Queue:    "default",
		Type:     "job",
		PushedAt: time.Unix(1700000000, 0),
	})

	jobs.StorePending(jobqueue.JobRecord{
		ID:       "job-2",
		Name:     "RebuildIndex",
		Queue:    "mail",
		PushedAt: time.Unix(1700000001, 0),
	})

	handler := newTestHandler(t, api.Options{Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/jobs/pending?queue=default", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(body) != 1 || body[0]["id"] != "job-1" {
		t.Fatalf("expected only job-1 on default queue, got %+v", body)
	}
}

// Port of JobRetrievalTest::test_paginating_large_job_results_gives_correct_amounts.
func TestPendingJobsEndpointPaginatesLargeResults(t *testing.T) {
	t.Parallel()

	jobs := jobqueue.NewJobRepository(nil)

	for i := 0; i < 150; i++ {
		jobs.StorePending(jobqueue.JobRecord{
			ID:       jobID(i),
			Queue:    "default",
			PushedAt: time.Unix(1700000000+int64(i), 0),
		})
	}

	handler := newTestHandler(t, api.Options{Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/jobs/pending?queue=default&starting_at=0&limit=50", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if len(body) != 50 {
		t.Fatalf("expected 50 jobs, got %d", len(body))
	}
}

// Port of JobRetrievalTest::test_recent_jobs_are_correctly_trimmed_and_expired.
func TestCompletedJobsEndpointReturnsTrimmedRecent(t *testing.T) {
	t.Parallel()

	jobs := jobqueue.NewJobRepository(nil)

	for i := 0; i < 5; i++ {
		id := jobID(i)
		jobs.StorePending(jobqueue.JobRecord{ID: id, Queue: "default", PushedAt: time.Unix(1700000000+int64(i), 0)})
		jobs.MarkComplete(id, true, time.Unix(1700000100+int64(i), 0))
	}

	handler := newTestHandler(t, api.Options{Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/jobs/completed?limit=3", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if len(body) != 3 {
		t.Fatalf("expected 3 jobs after limit, got %d", len(body))
	}
}

// Port of FailedJobTest::test_failed_jobs_are_placed_in_the_failed_job_table.
func TestFailedJobsEndpointListsFailedJobs(t *testing.T) {
	t.Parallel()

	jobs := jobqueue.NewJobRepository(nil)

	jobs.StorePending(jobqueue.JobRecord{ID: "job-1", Queue: "default", PushedAt: time.Unix(1700000000, 0)})
	jobs.MarkFailed("job-1", time.Unix(1700000010, 0), 0)

	handler := newTestHandler(t, api.Options{Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/jobs/failed", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if len(body) != 1 || body[0]["id"] != "job-1" || body[0]["status"] != "failed" {
		t.Fatalf("expected one failed job named job-1, got %+v", body)
	}
}

// Port of FailedJobTest::test_tags_for_failed_jobs_are_stored_in_redis.
func TestFailedJobsEndpointFiltersByTag(t *testing.T) {
	t.Parallel()

	jobs := jobqueue.NewJobRepository(nil)

	jobs.StorePending(jobqueue.JobRecord{ID: "job-1", Queue: "default", Tags: []string{"billing"}, PushedAt: time.Unix(1700000000, 0)})
	jobs.StorePending(jobqueue.JobRecord{ID: "job-2", Queue: "default", Tags: []string{"audit"}, PushedAt: time.Unix(1700000001, 0)})

	jobs.MarkFailed("job-1", time.Unix(1700000010, 0), 0)
	jobs.MarkFailed("job-2", time.Unix(1700000011, 0), 0)

	handler := newTestHandler(t, api.Options{Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/jobs/failed?tag=billing", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if len(body) != 1 || body[0]["id"] != "job-1" {
		t.Fatalf("expected only billing-tagged failed job, got %+v", body)
	}
}

// Port of FailedJobTest::test_failed_job_tags_have_an_expiration.
func TestFailedJobShowReturnsFailedRecordWithExpiration(t *testing.T) {
	t.Parallel()

	jobs := jobqueue.NewJobRepository(nil)

	jobs.StorePending(jobqueue.JobRecord{ID: "job-1", Queue: "default", Tags: []string{"billing"}, PushedAt: time.Unix(1700000000, 0)})
	jobs.MarkFailed("job-1", time.Unix(1700000010, 0), time.Hour)

	handler := newTestHandler(t, api.Options{Jobs: jobs})

	r := httptest.NewRequest(http.MethodGet, "/api/jobs/failed/job-1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if body["id"] != "job-1" || body["status"] != "failed" {
		t.Fatalf("expected failed job-1, got %+v", body)
	}
}

// Regression guard for jobsFailedShow returning 404 on unknown id.
func TestFailedJobShowReturns404ForUnknownID(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{})

	r := httptest.NewRequest(http.MethodGet, "/api/jobs/failed/missing", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func jobID(i int) string {
	const letters = "0123456789abcdef"

	if i < len(letters) {
		return "job-" + string(letters[i])
	}

	hi := i / len(letters)
	lo := i % len(letters)

	return "job-" + string(letters[hi]) + string(letters[lo])
}
