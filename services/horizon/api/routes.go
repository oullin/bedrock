package api

import "net/http"

func registerRoutes(mux *http.ServeMux, opts Options) {
	mux.HandleFunc("GET /api/stats", dashboardStats(opts))
	mux.HandleFunc("GET /api/master-supervisors", masterSupervisors(opts))

	mux.HandleFunc("GET /api/monitoring", monitoringIndex(opts))
	mux.HandleFunc("POST /api/monitoring", monitoringStore(opts))
	mux.HandleFunc("DELETE /api/monitoring/{tag}", monitoringDelete(opts))
	mux.HandleFunc("GET /api/monitoring/{tag}", monitoringShow(opts))

	mux.HandleFunc("GET /api/batches", batchesIndex(opts))

	mux.HandleFunc("GET /api/jobs/pending", jobsPending(opts))
	mux.HandleFunc("GET /api/jobs/completed", jobsCompleted(opts))
	mux.HandleFunc("GET /api/jobs/failed", jobsFailed(opts))
	mux.HandleFunc("GET /api/jobs/silenced", jobsSilenced(opts))
	mux.HandleFunc("GET /api/jobs/failed/{id}", jobsFailedShow(opts))

	mux.HandleFunc("GET /api/metrics/jobs", metricsJobs(opts))
	mux.HandleFunc("GET /api/metrics/queues", metricsQueues(opts))
	mux.HandleFunc("POST /api/metrics/snapshot", metricsSnapshot(opts))
}
