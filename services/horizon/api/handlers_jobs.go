package api

import (
	"net/http"
	"strings"

	"github.com/bedrock/packages/horizon"
)

// jobsPending ports PendingJobsController::index. Laravel Horizon paginates
// pending jobs for a queue via starting_at/limit query params.
func jobsPending(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queue := strings.TrimSpace(r.URL.Query().Get("queue"))

		if queue == "" {
			queue = "default"
		}

		offset := queryInt(r, "starting_at", 0)
		limit := queryInt(r, "limit", 50)

		writeJSON(w, http.StatusOK, jobList(opts.Jobs.Pending(queue, offset, limit)))
	}
}

// jobsCompleted ports CompletedJobsController::index.
func jobsCompleted(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := queryInt(r, "limit", 50)

		writeJSON(w, http.StatusOK, jobList(opts.Jobs.Recent(limit)))
	}
}

// jobsFailed ports FailedJobsController::index. The upstream controller
// optionally filters by tag; this Go port accepts the same "tag" query value
// and matches against stored failed jobs.
func jobsFailed(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := strings.TrimSpace(r.URL.Query().Get("tag"))
		jobs := failedJobs(opts, tag)

		writeJSON(w, http.StatusOK, jobList(jobs))
	}
}

// jobsFailedShow ports FailedJobController::show by returning a single failed
// job record.
func jobsFailedShow(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		job, ok := opts.Jobs.FindFailed(id)

		if !ok {
			http.Error(w, "failed job not found", http.StatusNotFound)

			return
		}

		writeJSON(w, http.StatusOK, jobPayload(job))
	}
}

// jobsSilenced ports SilencedJobsController::index by returning the
// snapshot slice configured on the dashboard options. Laravel Horizon stores
// these in a dedicated Redis set; Bedrock delegates persistence to the
// host app and renders whatever list it provides.
func jobsSilenced(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, jobList(opts.SilencedJobs))
	}
}

func failedJobs(opts Options, tag string) []horizon.JobRecord {
	jobs := opts.Jobs.Failed()

	if tag == "" {
		return jobs
	}

	filtered := jobs[:0]

	for _, job := range jobs {
		if jobHasTag(job, tag) {
			filtered = append(filtered, job)
		}
	}

	return append([]horizon.JobRecord(nil), filtered...)
}

func jobHasTag(job horizon.JobRecord, tag string) bool {
	for _, candidate := range job.Tags {
		if candidate == tag {
			return true
		}
	}

	return false
}

func jobList(jobs []horizon.JobRecord) []map[string]any {
	out := make([]map[string]any, 0, len(jobs))

	for _, job := range jobs {
		out = append(out, jobPayload(job))
	}

	return out
}

func jobPayload(job horizon.JobRecord) map[string]any {
	return map[string]any{
		"id":          job.ID,
		"name":        job.Name,
		"queue":       job.Queue,
		"type":        job.Type,
		"tags":        append([]string(nil), job.Tags...),
		"status":      string(job.Status),
		"pushedAt":    job.PushedAt,
		"reservedAt":  job.ReservedAt,
		"completedAt": job.CompletedAt,
		"failedAt":    job.FailedAt,
	}
}
