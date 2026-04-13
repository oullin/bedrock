package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestDigestAuth(t *testing.T) {
	t.Parallel()

	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++

		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Digest ") {
			w.Header().Set("WWW-Authenticate",
				`Digest realm="test@example.com", nonce="abc123", qop="auth", opaque="xyz"`)
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		// Verify the digest response contains expected fields.
		if !strings.Contains(auth, `username="admin"`) {
			t.Fatal("expected username in digest response")
		}

		if !strings.Contains(auth, `realm="test@example.com"`) {
			t.Fatal("expected realm in digest response")
		}

		if !strings.Contains(auth, `nonce="abc123"`) {
			t.Fatal("expected nonce in digest response")
		}

		if !strings.Contains(auth, `qop=auth`) {
			t.Fatal("expected qop in digest response")
		}

		if !strings.Contains(auth, `opaque="xyz"`) {
			t.Fatal("expected opaque in digest response")
		}

		w.Write([]byte("authenticated"))
	}))

	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().
		WithDigestAuth("admin", "secret").
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "authenticated" {
		t.Fatalf("expected authenticated, got %s", resp.Body())
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts (challenge + retry), got %d", attempts)
	}
}

func TestDigestAuthNotTriggeredOnNon401(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().
		WithDigestAuth("admin", "secret").
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "ok" {
		t.Fatalf("expected ok, got %s", resp.Body())
	}
}

func TestDigestAuthWithoutQop(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Digest ") {
			w.Header().Set("WWW-Authenticate",
				`Digest realm="simple", nonce="nonce123"`)
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		if !strings.Contains(auth, `username="user"`) {
			t.Fatal("expected username")
		}

		// Without qop, nc and cnonce should not be present.
		if strings.Contains(auth, "cnonce") {
			t.Fatal("cnonce should not be present without qop")
		}

		w.Write([]byte("ok"))
	}))

	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().
		WithDigestAuth("user", "pass").
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "ok" {
		t.Fatalf("expected ok, got %s", resp.Body())
	}
}
