package stress_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock/services/brain/api/stress"
)

func TestRunHappyPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(ts.Close)

	res, err := stress.Run(context.Background(), stress.Config{
		URL: ts.URL, Concurrency: 2, Requests: 10, Timeout: time.Second,
	})

	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if res.Total != 10 || res.Succeeded != 10 || res.Failed != 0 {
		t.Errorf("counts = total=%d succeeded=%d failed=%d, want 10/10/0",
			res.Total, res.Succeeded, res.Failed)
	}

	if res.StatusCodes[200] != 10 {
		t.Errorf("statusCodes[200] = %d, want 10", res.StatusCodes[200])
	}

	if res.P50 == 0 || res.P95 == 0 || res.P99 == 0 {
		t.Errorf("percentiles must be > 0: P50=%v P95=%v P99=%v", res.P50, res.P95, res.P99)
	}

	if res.RPS <= 0 {
		t.Errorf("RPS = %v, want > 0", res.RPS)
	}
}

func TestRunCountsFailureStatusCodes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	t.Cleanup(ts.Close)

	res, err := stress.Run(context.Background(), stress.Config{
		URL: ts.URL, Concurrency: 1, Requests: 5, Timeout: time.Second,
	})

	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if res.Succeeded != 0 || res.Failed != 5 {
		t.Errorf("succeeded=%d failed=%d, want 0/5", res.Succeeded, res.Failed)
	}

	if res.StatusCodes[500] != 5 {
		t.Errorf("statusCodes[500] = %d, want 5", res.StatusCodes[500])
	}
}

func TestRunMixedStatusesSumToTotal(t *testing.T) {
	var n int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		k := atomic.AddInt32(&n, 1)

		if k%2 == 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))

	t.Cleanup(ts.Close)

	res, err := stress.Run(context.Background(), stress.Config{
		URL: ts.URL, Concurrency: 1, Requests: 6, Timeout: time.Second,
	})

	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if res.Succeeded+res.Failed != res.Total {
		t.Errorf("succeeded(%d) + failed(%d) != total(%d)",
			res.Succeeded, res.Failed, res.Total)
	}
}

func TestRunAppliesDefaults(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(ts.Close)

	res, err := stress.Run(context.Background(), stress.Config{URL: ts.URL})

	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// Defaults: Requests=1, Concurrency=1, Timeout=10s, Method=GET.
	if res.Total != 1 {
		t.Errorf("default Total = %d, want 1", res.Total)
	}

	if res.Method != "GET" {
		t.Errorf("default Method = %q, want GET", res.Method)
	}

	if res.Succeeded != 1 {
		t.Errorf("default Succeeded = %d, want 1", res.Succeeded)
	}
}

func TestRunWithCanceledContextIsSafe(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(ts.Close)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := stress.Run(ctx, stress.Config{URL: ts.URL, Requests: 5, Concurrency: 2, Timeout: time.Second})

	if err != nil {
		t.Fatalf("Run returned error on canceled context: %v", err)
	}
	// With ctx canceled, workers should exit early and skip the work.
	if res.Succeeded+res.Failed > res.Total {
		t.Errorf("succeeded+failed > total under cancel: %+v", res)
	}
}
