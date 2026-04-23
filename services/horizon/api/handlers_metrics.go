package api

import (
	"net/http"
	"strings"
)

// metricsJobs ports JobMetricsController::index. The upstream controller
// returns throughput and average runtime for each job class, optionally
// filtered by name.
func metricsJobs(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := strings.TrimSpace(r.URL.Query().Get("job"))

		jobs := opts.Metrics.Jobs()
		response := make([]map[string]any, 0, len(jobs))

		for _, job := range jobs {
			if filter != "" && job != filter {
				continue
			}

			summary := opts.Metrics.Job(job)

			response = append(response, map[string]any{
				"name":           job,
				"throughput":     summary.Throughput,
				"averageRuntime": summary.AverageRuntime.Milliseconds(),
			})
		}

		writeJSON(w, http.StatusOK, response)
	}
}

// metricsQueues ports QueueMetricsController::index.
func metricsQueues(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		snapshots := opts.Metrics.Snapshots()

		if len(snapshots) == 0 {
			writeJSON(w, http.StatusOK, []map[string]any{})

			return
		}

		latest := snapshots[len(snapshots)-1]
		response := make([]map[string]any, 0, len(latest.QueueThroughput))

		for queue, throughput := range latest.QueueThroughput {
			summary := opts.Metrics.Queue(queue)

			response = append(response, map[string]any{
				"name":           queue,
				"throughput":     throughput,
				"averageRuntime": summary.AverageRuntime.Milliseconds(),
			})
		}

		writeJSON(w, http.StatusOK, response)
	}
}

// metricsSnapshot ports SnapshotCommand: capture a new metrics snapshot on
// demand. Horizon's scheduler fires this every five minutes; Bedrock exposes
// it as an HTTP endpoint the host can invoke from its own scheduler.
func metricsSnapshot(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		snapshot := opts.Metrics.SnapshotPerformance(opts.Now())

		writeJSON(w, http.StatusCreated, map[string]any{
			"recordedAt":    snapshot.RecordedAt,
			"jobsProcessed": snapshot.JobsProcessed,
		})
	}
}
