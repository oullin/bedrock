package routing

import "testing"

// Ref: @bedrock/code-0397
// RoutingRedirectorTest::testBasicRedirectTo
// RoutingRedirectorTest::testAwayDoesntValidateTheUrl
// RoutingRedirectorTest::testSecureRedirectToHttpsUrl
// RoutingRedirectorTest::testBackRedirectToHttpReferer
// RoutingRedirectorTest::testRoute
// RoutingRedirectorTest::testIntendedRedirectToIntendedUrlInSession
// RoutingRedirectorTest::testItSetsAndGetsValidIntendedUrl

// fakeSession implements [SessionStore] for redirector tests.
type fakeSession struct {
	flashes  map[string]any
	previous string
	old      map[string]string
}

func newFakeSession() *fakeSession {
	return &fakeSession{flashes: map[string]any{}, old: map[string]string{}}
}

func (s *fakeSession) Flash(key string, value any) { s.flashes[key] = value }
func (s *fakeSession) GetOldInput(key, fallback string) string {
	if v, ok := s.old[key]; ok {
		return v
	}

	return fallback
}
func (s *fakeSession) HasOldInput(key string) bool     { _, ok := s.old[key]; return ok }
func (s *fakeSession) FlashInput(input map[string]any) {}
func (s *fakeSession) Get(key string, fallback any) any {
	if key == "_previous.url" && s.previous != "" {
		return s.previous
	}

	return fallback
}
func (s *fakeSession) Put(key string, value any) {}

func TestRedirector(t *testing.T) {
	t.Run("test_to_redirect", func(t *testing.T) {
		gen, _ := newGen(t)
		red := NewRedirector(gen)
		resp := red.To("/dashboard", 0, nil, nil)

		if resp.Status != 302 {
			t.Errorf("status = %d", resp.Status)
		}

		if resp.URL != "http://example.com/dashboard" {
			t.Errorf("url = %q", resp.URL)
		}
	})

	t.Run("test_away_passthrough", func(t *testing.T) {
		gen, _ := newGen(t)
		red := NewRedirector(gen)
		resp := red.Away("https://other.example/foo", 301, nil)

		if resp.URL != "https://other.example/foo" || resp.Status != 301 {
			t.Errorf("got %+v", resp)
		}
	})

	t.Run("test_secure_forces_https", func(t *testing.T) {
		gen, _ := newGen(t)
		red := NewRedirector(gen)
		resp := red.Secure("/foo", 0, nil)

		if resp.URL[:8] != "https://" {
			t.Errorf("url = %q", resp.URL)
		}
	})

	t.Run("test_back_uses_session_previous", func(t *testing.T) {
		gen, _ := newGen(t)
		red := NewRedirector(gen)
		s := newFakeSession()
		s.previous = "/orig"
		red.SetSession(s)
		resp := red.Back(0, nil, "")

		if resp.URL != "/orig" {
			t.Errorf("url = %q", resp.URL)
		}
	})

	t.Run("test_back_falls_back", func(t *testing.T) {
		gen, _ := newGen(t)
		red := NewRedirector(gen)
		resp := red.Back(0, nil, "/login")

		if resp.URL != "/login" {
			t.Errorf("url = %q", resp.URL)
		}
	})

	t.Run("test_route_redirect", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user}", func() {}).Name("users.show")
		red := NewRedirector(gen)
		resp, err := red.Route("users.show", map[string]any{"user": "alice"}, 302, nil)

		if err != nil {
			t.Fatal(err)
		}

		if resp.URL != "http://example.com/users/alice" {
			t.Errorf("url = %q", resp.URL)
		}
	})

	t.Run("test_with_flashes_to_session_on_send", func(t *testing.T) {
		gen, _ := newGen(t)
		red := NewRedirector(gen)
		s := newFakeSession()
		red.SetSession(s)
		resp := red.To("/x", 0, nil, nil).With("status", "saved")

		if err := resp.Send(); err != nil {
			t.Fatal(err)
		}

		if s.flashes["status"] != "saved" {
			t.Errorf("flashes = %v", s.flashes)
		}
	})

	t.Run("test_intended_uses_stored_url", func(t *testing.T) {
		gen, _ := newGen(t)
		red := NewRedirector(gen)
		red.SetIntendedUrl("/profile")
		resp := red.Intended("/", 0, nil, nil)

		if resp.URL != "http://example.com/profile" {
			t.Errorf("url = %q", resp.URL)
		}
		// Intended URL should be cleared after consumption.
		if red.GetIntendedUrl() != "" {
			t.Error("intended URL not cleared")
		}
	})
}

func TestResponseFactory(t *testing.T) {
	t.Run("test_make_returns_response", func(t *testing.T) {
		f := NewResponseFactory(nil, nil)
		r := f.Make("body", 200, nil)

		if r.Status != 200 || r.Body != "body" {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("test_no_content_default_204", func(t *testing.T) {
		f := NewResponseFactory(nil, nil)
		r := f.NoContent(0, nil)

		if r.Status != 204 {
			t.Errorf("status = %d", r.Status)
		}
	})

	t.Run("test_json_sets_content_type", func(t *testing.T) {
		f := NewResponseFactory(nil, nil)
		r := f.JSON(map[string]any{"ok": true}, 0, nil)

		if r.Headers["Content-Type"][0] != "application/json" {
			t.Errorf("headers = %v", r.Headers)
		}
	})
}
