// Package stress runs load tests against a single target URL with a fixed
// concurrency and request budget. HTTP is delegated to
// packages/httpx/client; only the worker fan-out and the metric aggregation
// (p50/p95/p99) are owned here.
//
// The plan reserves packages/concurrency for a future swap of the worker
// pool. We keep stdlib sync here for now because concurrency.Manager.Run
// runs a fixed []Task once, while a stress test wants many short tasks
// pulled by a fixed worker count — not a clean fit. When concurrency grows
// a pool primitive, this file is the place to swap.
package stress

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bedrock/packages/httpx/client"
)

// Config is one stress-test invocation.
type Config struct {
	URL         string
	Method      string
	Concurrency int
	Requests    int
	Timeout     time.Duration
}

// Result is the report returned to the UI. Field names mirror laravel-brain
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

	factory := client.NewFactory()
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
				resp, err := factory.PendingRequest().Timeout(cfg.Timeout).Get(cfg.URL)
				dur := time.Since(t0)
				samplesMu.Lock()
				samples = append(samples, dur)
				samplesMu.Unlock()
				if err != nil {
					atomic.AddUint64(&failed, 1)
					continue
				}
				status := resp.Status()
				statusMu.Lock()
				statuses[status]++
				statusMu.Unlock()
				if status < 400 {
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
