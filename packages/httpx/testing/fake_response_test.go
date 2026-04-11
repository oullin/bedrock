package testing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httptesting "github.com/bedrock/packages/httpx/testing"
)

func TestAssertResponseStatus(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusCreated)

	httptesting.AssertResponse(t, rec).Status(http.StatusCreated)
}

func TestAssertResponseOk(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)

	httptesting.AssertResponse(t, rec).Ok()
}

func TestAssertResponseHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.Header().Set("X-Custom", "value")
	rec.WriteHeader(http.StatusOK)

	httptesting.AssertResponse(t, rec).HasHeader("X-Custom").Header("X-Custom", "value")
}

func TestAssertResponseBodyContains(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.Write([]byte("hello world"))

	httptesting.AssertResponse(t, rec).BodyContains("hello")
}

func TestAssertResponseBodyEquals(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.Write([]byte("exact match"))

	httptesting.AssertResponse(t, rec).BodyEquals("exact match")
}

func TestAssertResponseJSONPath(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")
	rec.Write([]byte(`{"name":"Taylor","age":30}`))

	httptesting.AssertResponse(t, rec).
		HasJSON("name").
		JSONPath("name", "Taylor").
		JSONPath("age", float64(30))
}

func TestAssertResponseMissingHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)

	httptesting.AssertResponse(t, rec).MissingHeader("X-Nonexistent")
}
