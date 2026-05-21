package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/services/horizon/api"
)

func seededBatches() *api.InMemoryBatches {
	store := api.NewInMemoryBatches()
	base := time.Unix(1700000000, 0)

	store.Store(api.Batch{ID: "b1", Name: "Send Welcome Emails", TotalJobs: 10, CreatedAt: base})
	store.Store(api.Batch{ID: "b2", Name: "Rebuild Search Index", TotalJobs: 5, CreatedAt: base.Add(time.Minute)})
	store.Store(api.Batch{ID: "b3", Name: "Send Password Resets", TotalJobs: 3, CreatedAt: base.Add(2 * time.Minute)})
	store.Store(api.Batch{ID: "b4", Name: "100% Off Coupons", TotalJobs: 1, CreatedAt: base.Add(3 * time.Minute)})
	store.Store(api.Batch{ID: "b5", Name: "Process _private_ Data", TotalJobs: 2, CreatedAt: base.Add(4 * time.Minute)})

	return store
}

func decodeBatchResponse(t *testing.T, body []byte) (batches []api.Batch, nextCursor string) {
	t.Helper()

	var payload struct {
		Batches    []api.Batch `json:"batches"`
		NextCursor string      `json:"nextCursor"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}

	return payload.Batches, payload.NextCursor
}

func TestBatchesCanBeSearchedByName(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{Batches: seededBatches()})

	r := httptest.NewRequest(http.MethodGet, "/api/batches?name=Send", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	batches, _ := decodeBatchResponse(t, w.Body.Bytes())

	if len(batches) != 2 {
		t.Fatalf("expected 2 batches matching %q, got %d", "Send", len(batches))
	}

	for _, batch := range batches {
		if !contains(batch.Name, "Send") {
			t.Fatalf("batch %q did not match name filter", batch.Name)
		}
	}
}

func TestBatchesCanBeSearchedById(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{Batches: seededBatches()})

	r := httptest.NewRequest(http.MethodGet, "/api/batches?name=b2", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	batches, _ := decodeBatchResponse(t, w.Body.Bytes())

	if len(batches) != 1 || batches[0].ID != "b2" {
		t.Fatalf("expected batch b2, got %+v", batches)
	}
}

func TestBatchesSearchEscapesLikeWildcards(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{Batches: seededBatches()})

	r := httptest.NewRequest(http.MethodGet, "/api/batches?name=_private_", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	batches, _ := decodeBatchResponse(t, w.Body.Bytes())

	if len(batches) != 1 || batches[0].ID != "b5" {
		t.Fatalf("expected only literal _private_ match, got %+v", batches)
	}
}

func TestBatchesSearchSupportsCursorPagination(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t, api.Options{Batches: seededBatches()})

	r := httptest.NewRequest(http.MethodGet, "/api/batches?limit=2", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	batches, cursor := decodeBatchResponse(t, w.Body.Bytes())

	if len(batches) != 2 || cursor == "" {
		t.Fatalf("expected first page of 2 with cursor, got %d batches cursor=%q", len(batches), cursor)
	}

	r = httptest.NewRequest(http.MethodGet, "/api/batches?limit=2&after="+cursor, nil)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	nextPage, _ := decodeBatchResponse(t, w.Body.Bytes())

	if len(nextPage) == 0 {
		t.Fatalf("expected next page after cursor %q", cursor)
	}

	if nextPage[0].ID == batches[0].ID {
		t.Fatalf("cursor pagination returned the same first batch")
	}
}

func contains(haystack, needle string) bool {
	return len(needle) == 0 || indexFold(haystack, needle) >= 0
}

func indexFold(haystack, needle string) int {
	if len(needle) > len(haystack) {
		return -1
	}

	for i := 0; i+len(needle) <= len(haystack); i++ {
		if equalFold(haystack[i:i+len(needle)], needle) {
			return i
		}
	}

	return -1
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if lower(a[i]) != lower(b[i]) {
			return false
		}
	}

	return true
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}

	return c
}
