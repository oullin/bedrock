package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestPipelineThen(t *testing.T) {
	t.Parallel()

	var order []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	pipeline := routing.NewPipeline(mw1, mw2)
	final := pipeline.Then(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	final.ServeHTTP(rec, req)

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}

	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(order), order)
	}

	for i, v := range expected {
		if order[i] != v {
			t.Fatalf("position %d: expected %q, got %q", i, v, order[i])
		}
	}
}

func TestPipelineThenFunc(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	pipeline := routing.NewPipeline(mw)
	final := pipeline.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	final.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected middleware to be called")
	}
}

func TestPipelineEmpty(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	pipeline := routing.NewPipeline()
	final := pipeline.Then(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	final.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestMiddlewarePrioritySort(t *testing.T) {
	t.Parallel()

	var order []string

	mwA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "A")
			next.ServeHTTP(w, r)
		})
	}

	mwB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "B")
			next.ServeHTTP(w, r)
		})
	}

	mwC := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "C")
			next.ServeHTTP(w, r)
		})
	}

	mp := routing.NewMiddlewarePriority()
	mp.Set(mwC, 1)
	mp.Set(mwA, 3)
	mp.Set(mwB, 2)

	sorted := mp.Sort([]routing.MiddlewareFunc{mwA, mwB, mwC})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	pipeline := routing.NewPipeline(sorted...)
	final := pipeline.Then(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	final.ServeHTTP(rec, req)

	expected := []string{"C", "B", "A", "handler"}

	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(order), order)
	}

	for i, v := range expected {
		if order[i] != v {
			t.Fatalf("position %d: expected %q, got %q", i, v, order[i])
		}
	}
}

func TestMiddlewarePrioritySortUnprioritized(t *testing.T) {
	t.Parallel()

	var order []string

	mwA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "A")
			next.ServeHTTP(w, r)
		})
	}

	mwB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "B")
			next.ServeHTTP(w, r)
		})
	}

	mwC := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "C")
			next.ServeHTTP(w, r)
		})
	}

	mp := routing.NewMiddlewarePriority()
	mp.Set(mwC, 1)

	sorted := mp.Sort([]routing.MiddlewareFunc{mwA, mwB, mwC})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	pipeline := routing.NewPipeline(sorted...)
	final := pipeline.Then(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	final.ServeHTTP(rec, req)

	if order[0] != "C" {
		t.Fatalf("expected C first (prioritized), got %q", order[0])
	}

	if order[1] != "A" || order[2] != "B" {
		t.Fatalf("expected A then B (original order), got %v", order)
	}
}

func TestMiddlewarePrioritySortEmpty(t *testing.T) {
	t.Parallel()

	mp := routing.NewMiddlewarePriority()
	sorted := mp.Sort(nil)

	if len(sorted) != 0 {
		t.Fatalf("expected empty result, got %d", len(sorted))
	}
}

func TestPipelineMiddlewareCanShortCircuit(t *testing.T) {
	t.Parallel()

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	}

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	pipeline := routing.NewPipeline(mw)
	final := pipeline.Then(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	final.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	if handlerCalled {
		t.Fatal("expected handler not to be called")
	}
}
