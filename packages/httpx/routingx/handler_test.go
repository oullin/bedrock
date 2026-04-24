package routingx_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/bedrock/packages/httpx/routingx"
	"github.com/bedrock/packages/routing"
)

func TestNewHandlerDispatchesRoute(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/hello", func() any {
		return "hello"
	})

	rec := perform(router, http.MethodGet, "/hello")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if rec.Body.String() != "hello" {
		t.Fatalf("body = %q, want %q", rec.Body.String(), "hello")
	}
}

func TestNewHandlerBindsRouteParametersFromHttpxRequest(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.PathValue("user")))
	})

	rec := perform(router, http.MethodGet, "/users/taylor%20otwell")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if rec.Body.String() != "taylor otwell" {
		t.Fatalf("body = %q, want decoded route parameter", rec.Body.String())
	}
}

func TestNewHandlerBindsRouteParametersFromRequestScopedRoute(t *testing.T) {
	router := routing.NewRouter(nil, nil)
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})

	var calls atomic.Int32

	router.Get("/billables/{type}/{id}", func() any {
		if calls.Add(1) == 1 {
			close(firstStarted)
			<-releaseFirst
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "%s:%s", r.PathValue("type"), r.PathValue("id"))
		})
	})

	handler := routingx.NewHandler(router)
	firstDone := make(chan *httptest.ResponseRecorder, 1)

	go func() {
		firstDone <- performWithHandler(handler, http.MethodGet, "/billables/team/1")
	}()

	<-firstStarted

	second := performWithHandler(handler, http.MethodGet, "/billables/user/2")
	close(releaseFirst)
	first := <-firstDone

	if second.Body.String() != "user:2" {
		t.Fatalf("second body = %q, want user:2", second.Body.String())
	}

	if first.Body.String() != "team:1" {
		t.Fatalf("first body = %q, want team:1", first.Body.String())
	}
}

func TestNewHandlerReturnsNotFound(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	rec := perform(router, http.MethodGet, "/missing")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestNewHandlerReturnsMethodNotAllowed(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/users", func() any {
		return "users"
	})

	rec := perform(router, http.MethodPost, "/users")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}

	if allow := rec.Header().Get("Allow"); !strings.Contains(allow, http.MethodGet) {
		t.Fatalf("Allow = %q, want GET", allow)
	}
}

func TestNewHandlerWritesRoutingHTTPResponse(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/custom", func() any {
		return &routing.HTTPResponse{
			Body:   "created",
			Status: http.StatusCreated,
			Headers: map[string][]string{
				"Content-Type": {"text/custom"},
				"X-Custom":     {"one", "two"},
			},
		}
	})

	rec := perform(router, http.MethodGet, "/custom")

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	if rec.Header().Get("Content-Type") != "text/custom" {
		t.Fatalf("Content-Type = %q, want text/custom", rec.Header().Get("Content-Type"))
	}

	if values := rec.Header().Values("X-Custom"); len(values) != 2 {
		t.Fatalf("X-Custom values = %v, want two values", values)
	}

	if rec.Body.String() != "created" {
		t.Fatalf("body = %q, want created", rec.Body.String())
	}
}

func TestNewHandlerWritesRedirectResponse(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/from", func() any {
		return &routing.RedirectResponse{
			URL:    "/to",
			Status: http.StatusMovedPermanently,
			Headers: map[string][]string{
				"X-Redirect": {"yes"},
			},
		}
	})

	rec := perform(router, http.MethodGet, "/from")

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMovedPermanently)
	}

	if rec.Header().Get("Location") != "/to" {
		t.Fatalf("Location = %q, want /to", rec.Header().Get("Location"))
	}

	if rec.Header().Get("X-Redirect") != "yes" {
		t.Fatalf("X-Redirect = %q, want yes", rec.Header().Get("X-Redirect"))
	}
}

func TestNewHandlerWritesJSONResponse(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/json", func() any {
		return map[string]any{"ok": true}
	})

	rec := perform(router, http.MethodGet, "/json")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}

	if strings.TrimSpace(rec.Body.String()) != `{"ok":true}` {
		t.Fatalf("body = %q, want JSON object", rec.Body.String())
	}
}

func TestNewHandlerDispatchesHTTPHandlerFunc(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(r.PathValue("user")))
	})

	rec := perform(router, http.MethodGet, "/users/taylor%20otwell")

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}

	if rec.Body.String() != "taylor otwell" {
		t.Fatalf("body = %q, want path value", rec.Body.String())
	}
}

func TestNewHandlerDispatchesReturnedHandlerFuncValue(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/fn/{id}", func() any {
		return func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("id=" + r.PathValue("id")))
		}
	})

	rec := perform(router, http.MethodGet, "/fn/42")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if rec.Body.String() != "id=42" {
		t.Fatalf("body = %q, want id=42", rec.Body.String())
	}
}

func TestNewHandlerDispatchesReturnedHTTPHandler(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	router.Get("/handler", func() any {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-From-Handler", "yes")
			_, _ = w.Write([]byte("handler"))
		})
	})

	rec := perform(router, http.MethodGet, "/handler")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if rec.Header().Get("X-From-Handler") != "yes" {
		t.Fatalf("missing X-From-Handler header")
	}
}

func perform(router *routing.Router, method, target string) *httptest.ResponseRecorder {
	return performWithHandler(routingx.NewHandler(router), method, target)
}

func performWithHandler(handler http.Handler, method, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, nil)
	handler.ServeHTTP(rec, req)

	return rec
}
