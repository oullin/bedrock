package precognition_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/precognition"
	"github.com/bedrock/packages/routing"
)

func TestNonPrecognitivePassesThrough(t *testing.T) {
	t.Parallel()

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mw := precognition.New()
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

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)

	mw.Wrap(handler).ServeHTTP(rec, r)

	vary := rec.Header().Get("Vary")

	if vary != "Precognition" {
		t.Fatalf("expected Vary: Precognition, got %q", vary)
	}
}

func TestPrecognitiveMarksContext(t *testing.T) {
	t.Parallel()

	var isPrecognitive bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPrecognitive = precognition.IsPrecognitive(r)
		w.WriteHeader(http.StatusOK)
	})

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if !isPrecognitive {
		t.Fatal("expected request to be marked precognitive in context")
	}
}

func TestPrecognitiveAddsPrecognitionResponseHeader(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if got := rec.Header().Get("Precognition"); got != "true" {
		t.Fatalf("expected Precognition: true response header, got %q", got)
	}
}

func TestPrecognitiveAddsVaryHeader(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	values := rec.Header().Values("Vary")
	found := false

	for _, v := range values {
		if v == "Precognition" {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected Vary header to contain Precognition, got %v", values)
	}
}

// emptyBag simulates a validator with no errors.
type emptyBag struct{}

func (emptyBag) IsEmpty() bool { return true }

func TestPrecognitiveWithValidationSuccess(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hook := precognition.AfterValidationHook(r)
		hook(emptyBag{})
		w.WriteHeader(http.StatusOK)
	})

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")
	r.Header.Set("Precognition-Validate-Only", "name,email")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if got := rec.Header().Get("Precognition-Success"); got != "true" {
		t.Fatalf("expected Precognition-Success: true, got %q", got)
	}

	if got := rec.Header().Get("Precognition"); got != "true" {
		t.Fatalf("expected Precognition: true, got %q", got)
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

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")
	r.Header.Set("Precognition-Validate-Only", "name")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", rec.Code)
	}

	if got := rec.Header().Get("Precognition"); got != "true" {
		t.Fatalf("expected Precognition: true on error response, got %q", got)
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
		hook := precognition.AfterValidationHook(r)
		hook(emptyBag{})

		handlerExecuted = true
		w.WriteHeader(http.StatusCreated)
	})

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")
	r.Header.Set("Precognition-Validate-Only", "name")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if handlerExecuted {
		t.Fatal("expected handler logic after AfterValidationHook to not execute")
	}

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestRecoversPrecognitionSuccessPanic(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(precognition.SuccessResponse{})
	})

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")

	mw.Wrap(handler).ServeHTTP(rec, r)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if got := rec.Header().Get("Precognition-Success"); got != "true" {
		t.Fatalf("expected Precognition-Success: true, got %q", got)
	}
}

func TestUnexpectedPanicIsNotSwallowed(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected error")
	})

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")

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

	mw := precognition.New()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	mw.Wrap(handler).ServeHTTP(rec, r)

	values := rec.Header().Values("Vary")
	hasEncoding := false
	hasPrecognition := false

	for _, v := range values {
		if v == "Accept-Encoding" {
			hasEncoding = true
		}

		if v == "Precognition" {
			hasPrecognition = true
		}
	}

	if !hasEncoding {
		t.Fatal("expected Vary to contain Accept-Encoding")
	}

	if !hasPrecognition {
		t.Fatal("expected Vary to contain Precognition")
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

		if _, ok := cd.(*precognition.CallableDispatcher); !ok {
			t.Errorf("expected CallableDispatcher to be swapped to precognition version, got %T", cd)
		}

		ctd, _ := c.Make("routing.controller_dispatcher")

		if _, ok := ctd.(*precognition.ControllerDispatcher); !ok {
			t.Errorf("expected ControllerDispatcher to be swapped to precognition version, got %T", ctd)
		}

		w.WriteHeader(http.StatusOK)
	})

	mw := precognition.New(c)
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	r.Header.Set("Precognition", "true")

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
