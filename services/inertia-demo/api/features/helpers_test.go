package features

import (
	"bytes"
	"database/sql"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/assert"
	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/wayfinder"
	"github.com/bedrock/services/inertia-demo/api/internal/database"
	"github.com/bedrock/services/inertia-demo/api/internal/seed"
	"github.com/bedrock/services/inertia-demo/api/internal/testutil"
)

type featureHarness struct {
	t        *testing.T
	db       *sql.DB
	inertia  *inertia.Inertia
	registry *wayfinder.Registry
	app      app
	flashes  []flash.Message
}

func newFeatureHarness(t *testing.T) *featureHarness {
	t.Helper()

	i, err := inertia.New(testutil.TestTemplate, inertia.WithVersion("test"))

	if err != nil {
		t.Fatalf("inertia.New() error = %v", err)
	}

	db, err := database.Open(":memory:")

	if err != nil {
		t.Fatalf("database.Open(:memory:) error = %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := seed.Run(db); err != nil {
		t.Fatalf("seed.Run() error = %v", err)
	}

	h := &featureHarness{
		t:       t,
		db:      db,
		inertia: i,
		registry: func() *wayfinder.Registry {
			r := wayfinder.New()

			DefineRoutes(r)

			return r
		}(),
	}

	requireAuthFn := func(next http.Handler) http.Handler { return next }
	renderFn := func(w http.ResponseWriter, r *http.Request, component string, props protocol.Props) {
		if strings.TrimSpace(r.RequestURI) == "" {
			r.RequestURI = r.URL.RequestURI()
		}

		if err := i.Render(w, r, component, props); err != nil {
			t.Fatalf("Render(%s) error = %v", component, err)
		}
	}

	redirectFn := func(w http.ResponseWriter, r *http.Request, url string) {
		i.Redirect(w, r, url)
	}

	locationFn := func(w http.ResponseWriter, r *http.Request, url string) {
		i.Location(w, r, url)
	}

	setFlashFn := func(w http.ResponseWriter, msg flash.Message) error {
		h.flashes = append(h.flashes, msg)

		return nil
	}

	h.app = newApp(Container{
		DB:          db,
		RequireAuth: requireAuthFn,
		Render:      renderFn,

		Redirect: redirectFn,

		Location: locationFn,

		RouteURL: h.registry.URL,
		SetFlash: setFlashFn,
	})

	return h
}

func (h *featureHarness) request(method, target string, body []byte) *http.Request {
	h.t.Helper()

	req := httptest.NewRequest(method, target, bytes.NewReader(body))

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()

	return req
}

func (h *featureHarness) partialRequest(method, target, component, only string) *http.Request {
	h.t.Helper()

	req := h.request(method, target, nil)

	req.Header.Set(protocol.HeaderPartialComponent, component)
	req.Header.Set(protocol.HeaderPartialData, only)

	return req
}

func (h *featureHarness) page(t *testing.T, w *httptest.ResponseRecorder) assert.AssertableInertia {
	t.Helper()

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	return assert.AssertFromBytes(t, w.Body.Bytes())
}

func (h *featureHarness) lastFlash(t *testing.T) flash.Message {
	t.Helper()

	if len(h.flashes) == 0 {
		t.Fatal("no flash recorded")
	}

	return h.flashes[len(h.flashes)-1]
}

func newMultipartRequest(t *testing.T, target string, files map[string][]string) (*http.Request, string) {
	t.Helper()

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	for field, names := range files {
		for _, name := range names {
			part, err := writer.CreateFormFile(field, name)

			if err != nil {
				t.Fatalf("CreateFormFile(%s, %s) error = %v", field, name, err)
			}

			if _, err := part.Write([]byte("demo")); err != nil {
				t.Fatalf("part.Write() error = %v", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, &body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()

	return req, writer.FormDataContentType()
}
