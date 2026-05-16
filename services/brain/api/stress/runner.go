// Package stress runs load tests against a single target URL with a fixed
// concurrency and request budget. Per the plan, the worker pool comes from
// packages/concurrency (Go) and the HTTP client from packages/httpx/client.
// Phase 11 implements the metric aggregation locally — the bedrock packages
// stay the only HTTP/concurrency surface this code uses.
//
// Until packages/concurrency and packages/httpx/client are wired into
// brain's go.mod (a later commit), we use the stdlib `net/http` and a
// channel-based worker pool. This file is the *single* place in services/
// brain/api/ allowed to do so — the package-level grep audit in the plan
// targets this file and cmd/brain only.
package stress

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Config is one stress-test invocation.
type Config struct {
	URL         string
	Method      string
	Concurrency int
	Requests    int
	Timeout     time.Duration
}

// Result is the report returned to the UI. Field names mirror upstream-brain
// so the existing viewer markup renders unchanged.
type Result struct {
	URL         string        `json:"url"`
	Method      string        `json:"method"`
	Total       int           `json:"total"`
	Succeeded   int           `json:"succeeded"`
	Failed      int           `json:"failed"`
	Duration    time.Duration `json:"durationNs"`
	RPS         float64       `json:"rps"`
	P50         time.Duration `json:"p50Ns"`
	P95         time.Duration `json:"p95Ns"`
	P99         time.Duration `json:"p99Ns"`
	StatusCodes map[int]int   `json:"statusCodes"`
}

// Run executes cfg and returns the aggregate metrics.
func Run(ctx context.Context, cfg Config) (Result, error) {
	if cfg.Method == "" {
		cfg.Method = "GET"
	}
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}
	if cfg.Requests < 1 {
		cfg.Requests = 1
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	client := &http.Client{Timeout: cfg.Timeout}
	jobs := make(chan struct{}, cfg.Requests)
	for i := 0; i < cfg.Requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	var (
		succeeded uint64
		failed    uint64
		samplesMu sync.Mutex
		samples   = make([]time.Duration, 0, cfg.Requests)
		statusMu  sync.Mutex
		statuses  = map[int]int{}
		wg        sync.WaitGroup
		start     = time.Now()
	)
	for w := 0; w < cfg.Concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				if ctx.Err() != nil {
					return
				}
				t0 := time.Now()
				req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.URL, nil)
				if err != nil {
					atomic.AddUint64(&failed, 1)
					continue
				}
				resp, err := client.Do(req)
				dur := time.Since(t0)
				samplesMu.Lock()
				samples = append(samples, dur)
				samplesMu.Unlock()
				if err != nil {
					atomic.AddUint64(&failed, 1)
					continue
				}
				_ = resp.Body.Close()
				statusMu.Lock()
				statuses[resp.StatusCode]++
				statusMu.Unlock()
				if resp.StatusCode < 400 {
					atomic.AddUint64(&succeeded, 1)
				} else {
					atomic.AddUint64(&failed, 1)
				}
			}
		}()
	}
	wg.Wait()
	total := time.Since(start)

	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	p := func(pct float64) time.Duration {
		if len(samples) == 0 {
			return 0
		}
		idx := int(float64(len(samples)-1) * pct)
		return samples[idx]
	}
	return Result{
		URL:         cfg.URL,
		Method:      cfg.Method,
		Total:       cfg.Requests,
		Succeeded:   int(succeeded),
		Failed:      int(failed),
		Duration:    total,
		RPS:         float64(cfg.Requests) / total.Seconds(),
		P50:         p(0.50),
		P95:         p(0.95),
		P99:         p(0.99),
		StatusCodes: statuses,
	}, nil
}
