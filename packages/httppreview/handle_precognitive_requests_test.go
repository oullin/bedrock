package httppreview_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/httppreview"
	"github.com/bedrock/packages/routing"
)

// emptyBag simulates a validator with no errors.
type emptyBag struct{}

func TestNonPrecognitivePassesThrough(t *testing.T) {
	t.Parallel()

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)

	mw.Wrap(handler).ServeHTTP(rec, r)

	if !called {
		t.Fatal("expected handler to be called for non-precognitive request")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "ok" {
		t.Fatalf("expected body 'ok', got %q", rec.Body.String())
	}
}

func TestNonPrecognitiveGetsVaryHeader(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)

	mw.Wrap(handler).ServeHTTP(rec, r)

	vary := rec.Header().Get("Vary")

	if vary != "HTTPPreview" {
		t.Fatalf("expected Vary: HTTPPreview, got %q", vary)
	}
}

func TestPrecognitiveMarksContext(t *testing.T) {
	t.Parallel()

	var isPrecognitive bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPrecognitive = httppreview.IsPrecognitive(r)
		w.WriteHeader(http.StatusOK)
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if !isPrecognitive {
		t.Fatal("expected request to be marked precognitive in context")
	}
}

func TestPrecognitiveAddsHTTPPreviewResponseHeader(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if got := rec.Header().Get("HTTPPreview"); got != "true" {
		t.Fatalf("expected HTTPPreview: true response header, got %q", got)
	}
}

func TestPrecognitiveAddsVaryHeader(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	values := rec.Header().Values("Vary")
	found := false

	for _, v := range values {
		if v == "HTTPPreview" {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected Vary header to contain HTTPPreview, got %v", values)
	}
}

func (emptyBag) IsEmpty() bool { return true }

func TestPrecognitiveWithValidationSuccess(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hook := httppreview.AfterValidationHook(r)
		hook(emptyBag{})
		w.WriteHeader(http.StatusOK)
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")
	r.Header.Set("HTTPPreview-Validate-Only", "name,email")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if got := rec.Header().Get("HTTPPreview-Success"); got != "true" {
		t.Fatalf("expected HTTPPreview-Success: true, got %q", got)
	}

	if got := rec.Header().Get("HTTPPreview"); got != "true" {
		t.Fatalf("expected HTTPPreview: true, got %q", got)
	}
}

func TestPrecognitiveWithValidationFailure(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string][]string{
				"name": {"The name field is required."},
			},
		})
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")
	r.Header.Set("HTTPPreview-Validate-Only", "name")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", rec.Code)
	}

	if got := rec.Header().Get("HTTPPreview"); got != "true" {
		t.Fatalf("expected HTTPPreview: true on error response, got %q", got)
	}

	var body map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if _, ok := body["errors"]; !ok {
		t.Fatal("expected errors in response body")
	}
}

func TestPrecognitiveHandlerNotExecuted(t *testing.T) {
	t.Parallel()

	handlerExecuted := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hook := httppreview.AfterValidationHook(r)
		hook(emptyBag{})

		handlerExecuted = true
		w.WriteHeader(http.StatusCreated)
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")
	r.Header.Set("HTTPPreview-Validate-Only", "name")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if handlerExecuted {
		t.Fatal("expected handler logic after AfterValidationHook to not execute")
	}

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestRecoversHTTPPreviewSuccessPanic(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(httppreview.SuccessResponse{})
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if got := rec.Header().Get("HTTPPreview-Success"); got != "true" {
		t.Fatalf("expected HTTPPreview-Success: true, got %q", got)
	}
}

func TestUnexpectedPanicIsNotSwallowed(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected error")
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected unexpected panic to be re-raised")
		}

		if v != "unexpected error" {
			t.Fatalf("expected 'unexpected error' panic, got %v", v)
		}
	}()

	mw.Wrap(handler).ServeHTTP(rec, r)
}

func TestVaryHeaderWithExistingValues(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.WriteHeader(http.StatusOK)
	})

	mw := httppreview.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	mw.Wrap(handler).ServeHTTP(rec, r)

	values := rec.Header().Values("Vary")
	hasEncoding := false
	hasHTTPPreview := false

	for _, v := range values {
		if v == "Accept-Encoding" {
			hasEncoding = true
		}

		if v == "HTTPPreview" {
			hasHTTPPreview = true
		}
	}

	if !hasEncoding {
		t.Fatal("expected Vary to contain Accept-Encoding")
	}

	if !hasHTTPPreview {
		t.Fatal("expected Vary to contain HTTPPreview")
	}
}

func TestDispatcherSwapping(t *testing.T) {
	t.Parallel()

	c := container.New()
	origCallable := routing.NewCallableDispatcher(nil)
	origController := routing.NewControllerDispatcher(nil)
	c.Instance("routing.callable_dispatcher", origCallable)
	c.Instance("routing.controller_dispatcher", origController)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify dispatchers were swapped.
		cd, _ := c.Make("routing.callable_dispatcher")

		if _, ok := cd.(*httppreview.CallableDispatcher); !ok {
			t.Errorf("expected CallableDispatcher to be swapped to httppreview version, got %T", cd)
		}

		ctd, _ := c.Make("routing.controller_dispatcher")

		if _, ok := ctd.(*httppreview.ControllerDispatcher); !ok {
			t.Errorf("expected ControllerDispatcher to be swapped to httppreview version, got %T", ctd)
		}

		w.WriteHeader(http.StatusOK)
	})

	mw := httppreview.New(c)
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("HTTPPreview", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	// Verify dispatchers were restored after the request.
	cd, _ := c.Make("routing.callable_dispatcher")

	if cd != origCallable {
		t.Fatal("expected original callable dispatcher to be restored")
	}

	ctd, _ := c.Make("routing.controller_dispatcher")

	if ctd != origController {
		t.Fatal("expected original controller dispatcher to be restored")
	}
}
