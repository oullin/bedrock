package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/horizon"
	"github.com/bedrock/services/horizon/api"
)

func newTestHandler(t *testing.T, opts api.Options) http.Handler {
	t.Helper()

	if opts.Supervisors == nil {
		opts.Supervisors = api.NewSliceSupervisors(nil)
	}

	if opts.Batches == nil {
		opts.Batches = api.NewInMemoryBatches()
	}

	return api.NewHandler(opts)
}

// Port of DashboardStatsControllerTest::test_all_stats_are_correctly_returned.
func TestDashboardStatsReturnsAggregatedCounters(t *testing.T) {
	t.Parallel()

	repo := horizon.NewInMemoryRepository()

	if err := repo.Record(context.Background(), horizon.Snapshot{
		GeneratedAt: time.Unix(1700000000, 0),
		Queues: []horizon.QueueStatus{
			{Name: "default", Pending: 5, Processing: 2, Failed: 1, Throughput: 10, Wait: 7 * time.Second},
			{Name: "mail", Pending: 3, Throughput: 4, Wait: 2 * time.Second},
		},
	}); err != nil {
		t.Fatalf("record snapshot: %v", err)
	}

	handler := newTestHandler(t, api.Options{
		Repository: repo,
		Supervisors: api.NewSliceSupervisors([]api.MasterSupervisor{
			{Name: "horizon-1", Status: "running"},
		}),
	})

	r := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body["pendingJobs"].(float64) != 8 {
		t.Fatalf("pendingJobs = %v, want 8", body["pendingJobs"])
	}

	if body["failedJobs"].(float64) != 1 {
		t.Fatalf("failedJobs = %v, want 1", body["failedJobs"])
	}

	if body["jobsPerMinute"].(float64) != 14 {
		t.Fatalf("jobsPerMinute = %v, want 14", body["jobsPerMinute"])
	}

	if body["status"] != "running" {
		t.Fatalf("status = %v, want running", body["status"])
	}
}

// Port of DashboardStatsControllerTest::test_paused_status_is_reflected_if_all_master_supervisors_are_paused.
func TestDashboardStatsReflectsPausedStateWhenAllSupervisorsArePaused(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{
		Supervisors: api.NewSliceSupervisors([]api.MasterSupervisor{
			{Name: "horizon-1", Status: "paused"},
			{Name: "horizon-2", Status: "paused"},
		}),
	})

	r := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if body["status"] != "paused" {
		t.Fatalf("status = %v, want paused", body["status"])
	}
}

// Port of DashboardStatsControllerTest::test_paused_status_isnt_reflected_if_not_all_master_supervisors_are_paused.
func TestDashboardStatsDoesNotReflectPausedWhenSomeSupervisorsAreRunning(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{
		Supervisors: api.NewSliceSupervisors([]api.MasterSupervisor{
			{Name: "horizon-1", Status: "paused"},
			{Name: "horizon-2", Status: "running"},
		}),
	})

	r := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body map[string]any

	_ = json.NewDecoder(w.Body).Decode(&body)

	if body["status"] != "running" {
		t.Fatalf("status = %v, want running", body["status"])
	}
}

// Port of MasterSupervisorControllerTest::test_master_supervisor_listing_without_supervisors.
func TestMasterSupervisorListingWithoutSupervisorsReturnsEmpty(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{})

	r := httptest.NewRequest(http.MethodGet, "/api/master-supervisors", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Body.String() != "[]\n" {
		t.Fatalf("expected empty JSON array, got %q", w.Body.String())
	}
}

// Port of MasterSupervisorControllerTest::test_master_supervisor_listing_with_supervisors.
func TestMasterSupervisorListingWithSupervisorsReturnsThem(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{
		Supervisors: api.NewSliceSupervisors([]api.MasterSupervisor{
			{Name: "horizon-1", PID: 1234, Status: "running"},
		}),
	})

	r := httptest.NewRequest(http.MethodGet, "/api/master-supervisors", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []api.MasterSupervisor

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(body) != 1 || body[0].Name != "horizon-1" {
		t.Fatalf("got %v, want [horizon-1]", body)
	}
}

// Port of MasterSupervisorControllerTest::test_master_supervisor_with_custom_name_listing_with_supervisors.
func TestMasterSupervisorListingPreservesCustomNames(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{
		Supervisors: api.NewSliceSupervisors([]api.MasterSupervisor{
			{Name: "billing-worker", Status: "running"},
			{Name: "default-worker", Status: "running"},
		}),
	})

	r := httptest.NewRequest(http.MethodGet, "/api/master-supervisors", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var body []api.MasterSupervisor

	_ = json.NewDecoder(w.Body).Decode(&body)

	if len(body) != 2 || body[0].Name != "billing-worker" {
		t.Fatalf("got %v, want billing-worker first", body)
	}
}
