package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplicationBootsAndServesWelcome(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}

	if got := strings.TrimSpace(recorder.Body.String()); got != "Bedrock" {
		t.Fatalf("expected %q, got %q", "Bedrock", got)
	}
}

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	request := httptest.NewRequest(http.MethodGet, "/up", nil)
	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}

	if got := strings.TrimSpace(recorder.Body.String()); got != "ok" {
		t.Fatalf("unexpected body: %q", got)
	}
}

func TestConsoleCommand(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	status := app.HandleCommand([]string{"inspire"}, &stdout, &stderr)

	if status != 0 {
		t.Fatalf("unexpected status: %d stderr=%q", status, stderr.String())
	}

	if strings.TrimSpace(stdout.String()) == "" {
		t.Fatal("expected console output")
	}
}

func TestAppConfigUsesEnvironmentFile(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	name, err := app.Config().String("app.name")

	if err != nil {
		t.Fatalf("Config String: %v", err)
	}

	if name != "Bedrock Demo" {
		t.Fatalf("unexpected app name: %q", name)
	}
}
