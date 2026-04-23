package api

import (
	"encoding/json"
	"net/http"
)

// monitoringIndex ports MonitoringController::index: return the monitored
// tags and the number of matching pending jobs per tag.

// monitoringShow ports MonitoringController::paginate: paginated pending jobs
// filtered by monitored tag.

type monitorRequest struct {
	Tag string `json:"tag"`
}

func monitoringIndex(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		tags := opts.Monitoring.Monitoring()
		response := make([]map[string]any, 0, len(tags))

		for _, tag := range tags {
			jobs := opts.Jobs.JobIDsForTag(tag, 0, 0)

			response = append(response, map[string]any{
				"tag":   tag,
				"count": len(jobs),
			})
		}

		writeJSON(w, http.StatusOK, response)
	}
}

func monitoringShow(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")
		offset := nonNegativeQueryInt(r, "starting_at", 0)
		limit := nonNegativeQueryInt(r, "limit", 50)

		ids := opts.Jobs.JobIDsForTag(tag, offset, limit)

		if ids == nil {
			ids = []string{}
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"tag":  tag,
			"jobs": ids,
		})
	}
}

// monitoringStore ports MonitoringController::store.
func monitoringStore(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body monitorRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Tag == "" {
			http.Error(w, "tag required", http.StatusBadRequest)

			return
		}

		opts.Monitoring.Monitor([]string{body.Tag})
		writeJSON(w, http.StatusOK, map[string]string{"tag": body.Tag})
	}
}

// monitoringDelete ports MonitoringController::destroy.
func monitoringDelete(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")

		if tag == "" {
			http.Error(w, "tag required", http.StatusBadRequest)

			return
		}

		opts.Monitoring.StopMonitoring([]string{tag})
		w.WriteHeader(http.StatusNoContent)
	}
}
