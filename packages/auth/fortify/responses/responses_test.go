package responses_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/fortify/responses"
)

func TestDefaultRegistryAndJSONResponse(t *testing.T) {
	t.Parallel()

	registry := responses.DefaultRegistry()
	if registry.LoginResponse == nil || registry.LockoutResponse == nil {
		t.Fatal("expected default responses to be configured")
	}

	rec := httptest.NewRecorder()
	responses.JSONResponse{Status: http.StatusCreated}.ToResponse(rec, httptest.NewRequest(http.MethodGet, "/", nil), map[string]any{"status": "ok"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("unexpected body: %#v", body)
	}
}

func TestJSONAndErrorResponsesInferStatus(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/", nil)

	rec := httptest.NewRecorder()
	responses.JSONResponse{Status: http.StatusOK}.ToResponse(rec, request, auth.ErrInvalidCredentials)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected inferred status: %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	responses.ErrorResponse{Status: http.StatusBadRequest}.ToResponse(rec, request, &auth.ValidationError{Fields: map[string]string{"email": "required"}})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected validation status: %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	responses.ErrorResponse{Status: http.StatusBadRequest}.ToResponse(rec, request, &auth.ThrottleError{Scope: "login", RetryAfter: time.Minute})
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("unexpected throttle status: %d", rec.Code)
	}
}
