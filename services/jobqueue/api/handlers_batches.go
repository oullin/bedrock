package api

import (
	"net/http"
	"strings"
)

// batchesIndex ports BatchesController::index with search-by-name,
// search-by-id, LIKE-wildcard escaping, and cursor-style pagination.
func batchesIndex(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("name"))
		after := strings.TrimSpace(r.URL.Query().Get("after"))
		limit := queryInt(r, "limit", 25)

		if limit < 0 {
			limit = 0
		}

		matches := opts.Batches.Search(query, after, limit)

		if matches == nil {
			matches = []Batch{}
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"batches":    matches,
			"nextCursor": nextCursor(matches, limit),
		})
	}
}

func nextCursor(batches []Batch, limit int) string {
	if limit <= 0 || len(batches) < limit {
		return ""
	}

	return batches[len(batches)-1].ID
}
