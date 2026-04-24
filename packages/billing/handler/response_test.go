package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJsonResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	jsonResponse(rec, http.StatusCreated, map[string]any{"foo": "bar"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("code = %d, want 201", rec.Code)
	}

	var body map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body["foo"] != "bar" {
		t.Fatalf("body.foo = %v", body["foo"])
	}
}

func TestNoContent(t *testing.T) {
	rec := httptest.NewRecorder()

	noContent(rec)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("code = %d, want 204", rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Fatalf("body not empty: %q", rec.Body.String())
	}
}

func TestErrorResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	errorResponse(rec, http.StatusBadRequest, "bad input")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "bad input") {
		t.Fatalf("body missing message: %q", rec.Body.String())
	}
}

func TestValidationErrors(t *testing.T) {
	rec := httptest.NewRecorder()

	validationErrors(rec, map[string][]string{"name": {"required"}})

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("code = %d, want 422", rec.Code)
	}

	var body map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	errs, ok := body["errors"].(map[string]any)

	if !ok {
		t.Fatalf("errors = %#v", body["errors"])
	}

	list, ok := errs["name"].([]any)

	if !ok || len(list) != 1 || list[0] != "required" {
		t.Fatalf("name errors = %#v", errs["name"])
	}
}
