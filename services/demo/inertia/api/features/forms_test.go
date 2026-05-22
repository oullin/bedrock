package features

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestFormHandlers_RenderAndMethodGuards(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	tests := []struct {
		name      string
		component string
		target    string
		handler   func(http.ResponseWriter, *http.Request)
	}{
		{name: "use form", component: "Features/Forms/UseForm", target: "/features/forms/use-form", handler: h.app.useFormHandler},
		{name: "form component", component: "Features/Forms/FormComponent", target: "/features/forms/form-component", handler: h.app.formComponentHandler},
		{name: "file uploads", component: "Features/Forms/FileUploads", target: "/features/forms/file-uploads", handler: h.app.fileUploadsHandler},
		{name: "validation", component: "Features/Forms/Validation", target: "/features/forms/validation", handler: h.app.validationHandler},
		{name: "httppreview", component: "Features/Forms/HTTPPreview", target: "/features/forms/httppreview", handler: h.app.precognitionHandler},
		{name: "dotted keys", component: "Features/Forms/DottedKeys", target: "/features/forms/dotted-keys", handler: h.app.dottedKeysHandler},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := h.request(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			page := h.page(t, w)

			page.AssertComponent(t, tt.component)

			req = h.request(http.MethodPut, tt.target, nil)
			w = httptest.NewRecorder()

			tt.handler(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
			}
		})
	}

	for _, tt := range []struct {
		name      string
		component string
		target    string
		handler   func(http.ResponseWriter, *http.Request)
	}{
		{name: "use form context", component: "Features/Forms/UseFormContext", target: "/features/forms/use-form-context", handler: h.app.formContextHandler},
		{name: "routegen", component: "Features/Forms/RouteGen", target: "/features/forms/routegen", handler: h.app.routegenHandler},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := h.request(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			page := h.page(t, w)

			page.AssertComponent(t, tt.component)
		})
	}
}

func TestUseFormHandler_PostBranches(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req := h.request(http.MethodPost, "/features/forms/use-form", []byte(url.Values{
		"name":  {""},
		"email": {""},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.app.useFormHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/Forms/UseForm")

	errors, ok := page.Props["errors"].(map[string]any)

	if !ok || errors["name"] != "Name is required." || errors["email"] != "Email is required." {
		t.Fatalf("errors = %#v", page.Props["errors"])
	}

	req = h.request(http.MethodPost, "/features/forms/use-form", []byte(url.Values{
		"name":  {"Ada"},
		"email": {"ada@example.com"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	h.app.useFormHandler(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/features/forms/use-form" {
		t.Fatalf("status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	if flash := h.lastFlash(t); flash.Title != "Form submitted" {
		t.Fatalf("flash = %#v, want Form submitted", flash)
	}
}

func TestFormComponentAndValidationHandlers_PostBranches(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req := h.request(http.MethodPost, "/features/forms/form-component", []byte(url.Values{
		"name": {""},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.app.formComponentHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/Forms/FormComponent")

	req = h.request(http.MethodPost, "/features/forms/validation", []byte(url.Values{
		"name":  {""},
		"email": {""},
		"age":   {"17"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	h.app.validationHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Features/Forms/Validation")

	errors, ok := page.Props["errors"].(map[string]any)

	if !ok || len(errors) != 3 {
		t.Fatalf("errors = %#v, want three validation errors", page.Props["errors"])
	}

	req = h.request(http.MethodPost, "/features/forms/validation", []byte(url.Values{
		"name":  {"Valid Name"},
		"email": {"valid@example.com"},
		"age":   {"25"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	h.app.validationHandler(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/features/forms/validation" {
		t.Fatalf("status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}
}

func TestFileUploadsHandler_PostBranches(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req, _ := newMultipartRequest(t, "/features/forms/file-uploads", map[string][]string{
		"files": {"first.txt", "second.txt"},
		"photo": {"avatar.png"},
	})

	w := httptest.NewRecorder()

	h.app.fileUploadsHandler(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}

	if flash := h.lastFlash(t); flash.Message != "3 file(s) received." {
		t.Fatalf("flash = %#v, want file count 3", flash)
	}

	req = h.request(http.MethodPost, "/features/forms/file-uploads", nil)
	w = httptest.NewRecorder()

	h.app.fileUploadsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestValidationSecondaryHandler_PostBranches(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req := h.request(http.MethodGet, "/features/forms/validation/secondary", nil)
	w := httptest.NewRecorder()

	h.app.validationSecondaryHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = h.request(http.MethodPost, "/features/forms/validation/secondary", []byte(url.Values{
		"feedback": {""},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	h.app.validationSecondaryHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/Forms/Validation")

	req = h.request(http.MethodPost, "/features/forms/validation/secondary", []byte(url.Values{
		"feedback": {"Looks good."},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	h.app.validationSecondaryHandler(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}
}

func TestHTTPPreviewAndOptimisticHandlers(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req := h.request(http.MethodPost, "/features/forms/httppreview", []byte(url.Values{
		"username":              {"ab"},
		"email":                 {"invalid"},
		"password":              {"short"},
		"password_confirmation": {"mismatch"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.app.precognitionHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/Forms/HTTPPreview")

	req = h.request(http.MethodPost, "/features/forms/httppreview", []byte(url.Values{
		"username":              {"valid-user"},
		"email":                 {"valid@example.com"},
		"password":              {"password123"},
		"password_confirmation": {"password123"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	h.app.precognitionHandler(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/features/forms/httppreview" {
		t.Fatalf("status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = h.request(http.MethodGet, "/features/forms/optimistic-updates", nil)
	w = httptest.NewRecorder()

	h.app.optimisticUpdatesHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Features/Forms/OptimisticUpdates")

	page.AssertHasProp(t, "contacts")

	req = h.request(http.MethodPost, "/features/forms/optimistic-toggle/1", nil)

	req.SetPathValue("id", "1")

	w = httptest.NewRecorder()

	h.app.optimisticToggleHandler(w, req)

	if flash := h.lastFlash(t); flash.Kind != "success" {
		t.Fatalf("flash = %#v, want success", flash)
	}

	req = h.request(http.MethodPost, "/features/forms/optimistic-toggle/999999", nil)

	req.SetPathValue("id", "999999")

	w = httptest.NewRecorder()

	h.app.optimisticToggleHandler(w, req)

	if flash := h.lastFlash(t); flash.Kind != "error" {
		t.Fatalf("flash = %#v, want error", flash)
	}

	req = h.request(http.MethodPost, "/features/forms/optimistic-toggle/bad", nil)

	req.SetPathValue("id", "bad")

	w = httptest.NewRecorder()

	h.app.optimisticToggleHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestDottedKeysHandler_PostBranch(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req := h.request(http.MethodPost, "/features/forms/dotted-keys", []byte(url.Values{
		"user.name": {"Ada"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.app.dottedKeysHandler(w, req)

	if w.Code != http.StatusFound || w.Header().Get("Location") != "/features/forms/dotted-keys" {
		t.Fatalf("status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}
}

func TestFormHandlers_BadRequestAndDatabaseErrors(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	for _, tt := range []struct {
		name    string
		target  string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{name: "use form", target: "/features/forms/use-form", handler: h.app.useFormHandler},
		{name: "form component", target: "/features/forms/form-component", handler: h.app.formComponentHandler},
		{name: "validation", target: "/features/forms/validation", handler: h.app.validationHandler},
		{name: "validation secondary", target: "/features/forms/validation/secondary", handler: h.app.validationSecondaryHandler},
		{name: "httppreview", target: "/features/forms/httppreview", handler: h.app.precognitionHandler},
		{name: "dotted keys", target: "/features/forms/dotted-keys", handler: h.app.dottedKeysHandler},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.target, strings.NewReader("bad"))

			req.Header.Set("Content-Type", "multipart/form-data; boundary=bad")
			req.Header.Set(protocol.HeaderInertia, "true")

			req.RequestURI = req.URL.RequestURI()
			w := httptest.NewRecorder()

			tt.handler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}

	h.db.Close()

	req := h.request(http.MethodGet, "/features/forms/optimistic-updates", nil)
	w := httptest.NewRecorder()

	h.app.optimisticUpdatesHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("optimisticUpdatesHandler status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
