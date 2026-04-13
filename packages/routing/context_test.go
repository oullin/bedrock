package routing_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestContextText(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/text", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "hello")
	})

	req := httptest.NewRequest(http.MethodGet, "/text", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if rec.Body.String() != "hello" {
		t.Fatalf("expected 'hello', got %q", rec.Body.String())
	}

	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "text/plain") {
		t.Fatalf("expected text/plain, got %q", ct)
	}
}

func TestContextJSON(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/json", func(ctx *routing.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"key": "val"})
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}

	if !strings.Contains(rec.Body.String(), `"key"`) {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestContextRedirect(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/old", func(ctx *routing.Context) error {
		return ctx.Redirect("/new", http.StatusFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/new" {
		t.Fatalf("expected Location=/new, got %q", loc)
	}
}

func TestContextHTML(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/html", func(ctx *routing.Context) error {
		return ctx.HTML(http.StatusOK, "<h1>Hello</h1>")
	})

	req := httptest.NewRequest(http.MethodGet, "/html", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("expected text/html, got %q", ct)
	}

	if rec.Body.String() != "<h1>Hello</h1>" {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestContextNoContent(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Delete("/items/{id}", func(ctx *routing.Context) error {
		return ctx.NoContent()
	})

	req := httptest.NewRequest(http.MethodDelete, "/items/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestContextBlob(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/blob", func(ctx *routing.Context) error {
		return ctx.Blob(http.StatusOK, "application/octet-stream", []byte{0x01, 0x02})
	})

	req := httptest.NewRequest(http.MethodGet, "/blob", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if ct := rec.Header().Get("Content-Type"); ct != "application/octet-stream" {
		t.Fatalf("expected application/octet-stream, got %q", ct)
	}

	if rec.Body.Len() != 2 {
		t.Fatalf("expected 2 bytes, got %d", rec.Body.Len())
	}
}

func TestContextParam(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.Param("id"))
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "42" {
		t.Fatalf("expected '42', got %q", rec.Body.String())
	}
}

func TestContextQuery(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/search", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.Query("q"))
	})

	req := httptest.NewRequest(http.MethodGet, "/search?q=hello", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "hello" {
		t.Fatalf("expected 'hello', got %q", rec.Body.String())
	}
}

func TestContextQueryDefault(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/search", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.QueryDefault("q", "default"))
	})

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "default" {
		t.Fatalf("expected 'default', got %q", rec.Body.String())
	}
}

func TestContextHasQuery(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/check", func(ctx *routing.Context) error {
		if ctx.HasQuery("q") {
			return ctx.Text(http.StatusOK, "yes")
		}

		return ctx.Text(http.StatusOK, "no")
	})

	req := httptest.NewRequest(http.MethodGet, "/check?q=1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "yes" {
		t.Fatalf("expected 'yes', got %q", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/check", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "no" {
		t.Fatalf("expected 'no', got %q", rec.Body.String())
	}
}

func TestContextFormValue(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Post("/form", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.FormValue("name"))
	})

	form := url.Values{"name": {"Gustavo"}}
	req := httptest.NewRequest(http.MethodPost, "/form", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "Gustavo" {
		t.Fatalf("expected 'Gustavo', got %q", rec.Body.String())
	}
}

func TestContextMethod(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Post("/test", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.Method())
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != http.MethodPost {
		t.Fatalf("expected 'POST', got %q", rec.Body.String())
	}
}

func TestContextPath(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/some/path", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.Path())
	})

	req := httptest.NewRequest(http.MethodGet, "/some/path", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "/some/path" {
		t.Fatalf("expected '/some/path', got %q", rec.Body.String())
	}
}

func TestContextIsMethod(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Post("/test", func(ctx *routing.Context) error {
		if ctx.IsMethod(http.MethodPost) {
			return ctx.Text(http.StatusOK, "yes")
		}

		return ctx.Text(http.StatusOK, "no")
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "yes" {
		t.Fatalf("expected 'yes', got %q", rec.Body.String())
	}
}

func TestContextHeader(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/header", func(ctx *routing.Context) error {
		ctx.SetHeader("X-Custom", "value")

		return ctx.Text(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/header", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Header().Get("X-Custom") != "value" {
		t.Fatalf("expected X-Custom=value, got %q", rec.Header().Get("X-Custom"))
	}
}

func TestContextCookie(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/cookie", func(ctx *routing.Context) error {
		ctx.SetCookie(&http.Cookie{Name: "session", Value: "abc123"})

		return ctx.Text(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/cookie", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 || cookies[0].Name != "session" || cookies[0].Value != "abc123" {
		t.Fatalf("unexpected cookies: %v", cookies)
	}
}

func TestContextIP(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/ip", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.IP())
	})

	req := httptest.NewRequest(http.MethodGet, "/ip", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "192.168.1.1" {
		t.Fatalf("expected '192.168.1.1', got %q", rec.Body.String())
	}
}

func TestContextSetGetData(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/data", func(ctx *routing.Context) error {
		ctx.SetData("user", "gustavo")
		val, _ := ctx.GetData("user").(string)

		return ctx.Text(http.StatusOK, val)
	})

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "gustavo" {
		t.Fatalf("expected 'gustavo', got %q", rec.Body.String())
	}
}

func TestContextGetDataMissing(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/missing", func(ctx *routing.Context) error {
		val := ctx.GetData("nonexistent")

		if val != nil {
			t.Fatal("expected nil for missing key")
		}

		return ctx.NoContent()
	})

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
