package testing_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	packagetesting "github.com/bedrock/packages/testing"
)

type mockTB struct {
	failed   bool
	skipped  bool
	messages []string
}

func (m *mockTB) Attr(string, string)      {}
func (m *mockTB) ArtifactDir() string      { return "" }
func (m *mockTB) Cleanup(func())           {}
func (m *mockTB) Context() context.Context { return context.Background() }
func (m *mockTB) Error(args ...any) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprint(args...))
}
func (m *mockTB) Errorf(format string, args ...any) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}
func (m *mockTB) Fail()        { m.failed = true }
func (m *mockTB) FailNow()     { m.failed = true }
func (m *mockTB) Failed() bool { return m.failed }
func (m *mockTB) Fatal(args ...any) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprint(args...))
}
func (m *mockTB) Fatalf(format string, args ...any) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}
func (m *mockTB) Helper()               {}
func (m *mockTB) Log(args ...any)       {}
func (m *mockTB) Logf(string, ...any)   {}
func (m *mockTB) Name() string          { return "mockTB" }
func (m *mockTB) Setenv(string, string) {}
func (m *mockTB) Chdir(string)          {}
func (m *mockTB) Skip(args ...any)      { m.skipped = true }
func (m *mockTB) SkipNow()              { m.skipped = true }
func (m *mockTB) Skipf(string, ...any)  { m.skipped = true }
func (m *mockTB) Skipped() bool         { return m.skipped }
func (m *mockTB) TempDir() string {
	dir, _ := os.MkdirTemp("", "mocktb-*")

	return dir
}
func (m *mockTB) Output() io.Writer { return io.Discard }

type viewUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func newRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestResponseStatusAssertions(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertOk", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusOK)

		packagetesting.AssertResponse(t, rec).Ok()
	})

	t.Run("TestResponseTest::testAssertCreated", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusCreated)

		packagetesting.AssertResponse(t, rec).Created()
	})

	t.Run("TestResponseTest::testAssertAccepted", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusAccepted)

		packagetesting.AssertResponse(t, rec).Accepted()
	})

	t.Run("TestResponseTest::testAssertNoContentAsserts204StatusCodeByDefault", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusNoContent)

		packagetesting.AssertResponse(t, rec).NoContent()
	})

	t.Run("TestResponseTest::testAssertNoContentAssertsExpectedStatusCode", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusAccepted)

		packagetesting.AssertResponse(t, rec).NoContent(http.StatusAccepted)
	})

	t.Run("TestResponseTest::testAssertNoContentAssertsEmptyContent", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusNoContent)

		packagetesting.AssertResponse(t, rec).NoContent().BodyEquals("")
	})

	t.Run("TestResponseTest::testAssertNotFound", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusNotFound)

		packagetesting.AssertResponse(t, rec).NotFound()
	})

	t.Run("TestResponseTest::testAssertMovedPermanently", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/old")
		rec.WriteHeader(http.StatusMovedPermanently)

		packagetesting.AssertResponse(t, rec).MovedPermanently().Location("/old")
	})

	t.Run("TestResponseTest::testAssertFound", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/dashboard")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).Found().Location("/dashboard")
	})

	t.Run("TestResponseTest::testAssertNotModified", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusNotModified)

		packagetesting.AssertResponse(t, rec).NotModified()
	})

	t.Run("TestResponseTest::testAssertTemporaryRedirect", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/temp")
		rec.WriteHeader(http.StatusTemporaryRedirect)

		packagetesting.AssertResponse(t, rec).TemporaryRedirect().Location("/temp")
	})

	t.Run("TestResponseTest::testAssertPermanentRedirect", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/perm")
		rec.WriteHeader(http.StatusPermanentRedirect)

		packagetesting.AssertResponse(t, rec).PermanentRedirect().Location("/perm")
	})

	t.Run("TestResponseTest::testAssertBadRequest", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusBadRequest)

		packagetesting.AssertResponse(t, rec).BadRequest()
	})

	t.Run("TestResponseTest::testAssertMethodNotAllowed", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusMethodNotAllowed)

		packagetesting.AssertResponse(t, rec).MethodNotAllowed()
	})

	t.Run("TestResponseTest::testAssertNotAcceptable", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusNotAcceptable)

		packagetesting.AssertResponse(t, rec).NotAcceptable()
	})

	t.Run("TestResponseTest::testAssertForbidden", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusForbidden)

		packagetesting.AssertResponse(t, rec).Forbidden()
	})

	t.Run("TestResponseTest::testAssertUnauthorized", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusUnauthorized)

		packagetesting.AssertResponse(t, rec).Unauthorized()
	})

	t.Run("TestResponseTest::testAssertRequestTimeout", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusRequestTimeout)

		packagetesting.AssertResponse(t, rec).RequestTimeout()
	})

	t.Run("TestResponseTest::testAssertPaymentRequired", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusPaymentRequired)

		packagetesting.AssertResponse(t, rec).PaymentRequired()
	})

	t.Run("TestResponseTest::testAssertConflict", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusConflict)

		packagetesting.AssertResponse(t, rec).Conflict()
	})

	t.Run("TestResponseTest::testAssertGone", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusGone)

		packagetesting.AssertResponse(t, rec).Gone()
	})

	t.Run("TestResponseTest::testAssertTooManyRequests", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusTooManyRequests)

		packagetesting.AssertResponse(t, rec).TooManyRequests()
	})

	t.Run("TestResponseTest::testAssertFailedDependency", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusFailedDependency)

		packagetesting.AssertResponse(t, rec).FailedDependency()
	})

	t.Run("TestResponseTest::testAssertClientError", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusTooManyRequests)

		packagetesting.AssertResponse(t, rec).ClientError()
	})
}

func TestResponseRedirectAndViewAssertions(t *testing.T) {
	t.Parallel()

	t.Run("AssertRedirectToActionTest::testAssertRedirectToActionWithoutParameters", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/dashboard")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(string, map[string]string) string { return "/dashboard" }).
			RedirectToAction("dashboard")
	})

	t.Run("AssertRedirectToActionTest::testAssertRedirectToActionWithParameters", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/users/42/edit")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(name string, params map[string]string) string {
				return "/users/" + params["id"] + "/edit"
			}).
			RedirectToAction("users.edit", map[string]string{"id": "42"})
	})

	t.Run("AssertRedirectToRouteTest::testAssertRedirectToRouteWithRouteName", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/dashboard")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(string, map[string]string) string { return "/dashboard" }).
			RedirectToRoute("dashboard")
	})

	t.Run("AssertRedirectToRouteTest::testAssertRedirectToRouteWithRouteNameAndParams", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/users/42")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(name string, params map[string]string) string {
				return "/users/" + params["id"]
			}).
			RedirectToRoute("users.show", map[string]string{"id": "42"})
	})

	t.Run("AssertRedirectToRouteTest::testAssertRedirectToRouteWithRouteNameAndParamsWhenRouteUriIsEmpty", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(string, map[string]string) string { return "" }).
			RedirectToRoute("home")
	})

	t.Run("AssertRedirectToSignedRouteTest::testAssertRedirectToSignedRouteWithoutRouteName", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/signed")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(string, map[string]string) string { return "/signed" }).
			RedirectToSignedRoute("signed")
	})

	t.Run("AssertRedirectToSignedRouteTest::testAssertRedirectToSignedRouteWithRouteName", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/signed/dashboard")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(string, map[string]string) string { return "/signed/dashboard" }).
			RedirectToSignedRoute("signed.dashboard")
	})

	t.Run("AssertRedirectToSignedRouteTest::testAssertRedirectToSignedRouteWithRouteNameAndParams", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/signed/users/42")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(name string, params map[string]string) string {
				return "/signed/users/" + params["id"]
			}).
			RedirectToSignedRoute("signed.users.show", map[string]string{"id": "42"})
	})

	t.Run("AssertRedirectToSignedRouteTest::testAssertRedirectToSignedRouteWithRouteNameToTemporarySignedRoute", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/signed/temp")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).
			WithRouteResolver(func(string, map[string]string) string { return "/signed/temp" }).
			RedirectToTemporarySignedRoute("signed.temp")
	})

	t.Run("TestResponseTest::testAssertViewIs", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("dashboard", map[string]any{}).
			ViewIs("dashboard")
	})

	t.Run("TestResponseTest::testAssertViewHas", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": map[string]any{"name": "Taylor"}}).
			ViewHas("user")
	})

	t.Run("TestResponseTest::testAssertViewHasModel", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": viewUser{Name: "Taylor", Age: 30}}).
			ViewHas("user", viewUser{Name: "Taylor", Age: 30})
	})

	t.Run("TestResponseTest::testAssertViewHasWithClosure", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"age": 30}).
			ViewHas("age", func(v any) bool {
				n, ok := v.(int)
				return ok && n >= 18
			})
	})

	t.Run("TestResponseTest::testAssertViewHasWithValue", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"name": "Taylor"}).
			ViewHas("name", "Taylor")
	})

	t.Run("TestResponseTest::testAssertViewHasNested", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": map[string]any{"name": "Taylor"}}).
			ViewHas("user.name")
	})

	t.Run("TestResponseTest::testAssertViewHasWithNestedValue", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": map[string]any{"name": "Taylor"}}).
			ViewHas("user.name", "Taylor")
	})

	t.Run("TestResponseTest::testAssertViewHasEloquentCollection", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"items": []any{map[string]any{"id": 1}, map[string]any{"id": 2}}}).
			ViewHas("items", []any{map[string]any{"id": 1}, map[string]any{"id": 2}})
	})

	t.Run("TestResponseTest::testAssertViewHasEloquentCollectionRespectsOrder", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"items": []string{"a", "b"}}).
			ViewHas("items", []string{"a", "b"})
	})

	t.Run("TestResponseTest::testAssertViewHasEloquentCollectionRespectsType", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"items": []int{1, 2}}).
			ViewHas("items", []int{1, 2})
	})

	t.Run("TestResponseTest::testAssertViewHasEloquentCollectionRespectsSize", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"items": []int{1, 2, 3}}).
			ViewHas("items", []int{1, 2, 3})
	})

	t.Run("TestResponseTest::testAssertViewHasWithArray", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"items": []any{"a", "b"}}).
			ViewHas("items", []any{"a", "b"})
	})

	t.Run("TestResponseTest::testAssertViewHasAll", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": "Taylor", "email": "taylor@example.test"}).
			ViewHasAll("user", "email")
	})

	t.Run("TestResponseTest::testAssertViewMissing", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": "Taylor"}).
			ViewMissing("email")
	})

	t.Run("TestResponseTest::testAssertViewMissingNested", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": map[string]any{"name": "Taylor"}}).
			ViewMissing("user.email")
	})

	t.Run("TestResponseTest::testViewData", func(t *testing.T) {
		rec := newRecorder()

		value, ok := packagetesting.AssertResponse(t, rec).
			WithView("profile", map[string]any{"user": map[string]any{"name": "Taylor"}}).
			ViewData("user.name")

		if !ok || value != "Taylor" {
			t.Fatalf("expected view data value Taylor, got %v (ok=%v)", value, ok)
		}
	})
}

func TestResponseHeaderCookieSessionValidationAssertions(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertHeader", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("X-Custom", "value")

		packagetesting.AssertResponse(t, rec).Header("X-Custom", "value").HasHeader("X-Custom")
	})

	t.Run("TestResponseTest::testAssertHeaderMissing", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).MissingHeader("X-Custom")
	})

	t.Run("TestResponseTest::testAssertHeaderContainsSuccess", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("X-Trace", "abc-123")

		packagetesting.AssertResponse(t, rec).HeaderContains("X-Trace", "abc")
	})

	t.Run("TestResponseTest::testAssertPlainCookie", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Add("Set-Cookie", "plain=value; Path=/")

		packagetesting.AssertResponse(t, rec).PlainCookie("plain", "value")
	})

	t.Run("TestResponseTest::testAssertCookie", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Add("Set-Cookie", "session=abc; Path=/")

		packagetesting.AssertResponse(t, rec).Cookie("session", "abc")
	})

	t.Run("TestResponseTest::testAssertCookieExpired", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Add("Set-Cookie", "expired=gone; Path=/; Max-Age=-1")

		packagetesting.AssertResponse(t, rec).CookieExpired("expired")
	})

	t.Run("TestResponseTest::testAssertCookieNotExpired", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Add("Set-Cookie", "alive=1; Path=/")

		packagetesting.AssertResponse(t, rec).CookieNotExpired("alive")
	})

	t.Run("TestResponseTest::testAssertCookieMissing", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).CookieMissing("missing")
	})

	t.Run("TestResponseTest::testAssertSessionHasErrors", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"errors": map[string][]string{"email": []string{"required"}}}).
			SessionHasErrors(map[string][]string{"email": []string{"required"}})
	})

	t.Run("TestResponseTest::testAssertSessionHasNoErrors", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{}).
			SessionHasNoErrors()
	})

	t.Run("TestResponseTest::testAssertSessionHas", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"user": map[string]any{"name": "Taylor"}}).
			SessionHas("user.name", "Taylor")
	})

	t.Run("TestResponseTest::testAssertSessionDoesntHaveErrors", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{}).
			SessionDoesntHaveErrors()
	})

	t.Run("TestResponseTest::testAssertSessionHasInput", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"input": map[string]any{"name": "Taylor"}}).
			SessionHasInput("input.name", "Taylor")
	})

	t.Run("TestResponseTest::testAssertSessionMissing", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"user": map[string]any{"name": "Taylor"}}).
			SessionMissing("email")
	})

	t.Run("TestResponseTest::testAssertSessionMissingValueIsMissing", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"user": map[string]any{"name": "Taylor"}}).
			SessionMissingValue("user.name", "Alyssa")
	})

	t.Run("TestResponseTest::testAssertSessionMissingValueIsMissingClosure", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"user": map[string]any{"age": 30}}).
			SessionMissingValue("user.age", func(v any) bool {
				n, ok := v.(int)
				return ok && n < 18
			})
	})

	t.Run("TestResponseTest::testAssertSessionHasAllWithValues", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{
				"user":  map[string]any{"name": "Taylor"},
				"email": "taylor@example.test",
			}).
			SessionHasAll("user.name", "email")
	})

	t.Run("TestResponseTest::testAssertSessionHasAllWithMixedKeys", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{
				"user": map[string]any{"name": "Taylor"},
				"meta": map[string]any{"request_id": "abc123"},
			}).
			SessionHasAll("user.name", "meta.request_id")
	})

	t.Run("TestResponseTest::testAssertSessionHasAllWithClosures", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{
				"user":  map[string]any{"name": "Taylor"},
				"email": "taylor@example.test",
			}).
			SessionHas("user.name", func(v any) bool {
				return v == "Taylor"
			}).
			SessionHas("email", func(v any) bool {
				s, ok := v.(string)
				return ok && strings.Contains(s, "@")
			})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrors", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"required"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrors", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"status":"ok"}`))

		packagetesting.AssertResponse(t, rec).JSONMissingValidationErrors()
	})
}

func TestResponseStreamAndDownloadAssertions(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertStreamedAndAssertNotStreamed", func(t *testing.T) {
		streamed := newRecorder()
		streamed.Header().Set("Transfer-Encoding", "chunked")

		packagetesting.AssertResponse(t, streamed).Streamed()

		plain := newRecorder()
		packagetesting.AssertResponse(t, plain).NotStreamed()
	})

	t.Run("TestResponseTest::testAssertStreamedContent", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Transfer-Encoding", "chunked")
		rec.WriteString("alpha")

		packagetesting.AssertResponse(t, rec).StreamedContent("alpha")
	})

	t.Run("TestResponseTest::testAssertStreamedJsonContent", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Transfer-Encoding", "chunked")
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).StreamedJSONContent(map[string]any{"name": "Taylor"})
	})

	t.Run("TestResponseTest::testAssertStreamedBinaryFile", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Transfer-Encoding", "chunked")
		rec.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)
		rec.WriteString("alpha")

		packagetesting.AssertResponse(t, rec).StreamedBinaryFile("report.csv")
	})

	t.Run("TestResponseTest::testAssertStreamedJsonFile", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Transfer-Encoding", "chunked")
		rec.Header().Set("Content-Disposition", `attachment; filename="report.json"`)
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).StreamedJSONFile(map[string]any{"name": "Taylor"}, "report.json")
	})

	t.Run("TestResponseTest::testJsonAssertionsOnStreamedJsonContent", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Transfer-Encoding", "chunked")
		rec.Write([]byte(`{"name":"Taylor","roles":["admin","editor"]}`))

		packagetesting.AssertResponse(t, rec).
			StreamedJSONContent(map[string]any{"name": "Taylor", "roles": []any{"admin", "editor"}}).
			JSONPath("roles.0", "admin")
	})

	t.Run("TestResponseTest::testAssertDownloadOffered", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)

		packagetesting.AssertResponse(t, rec).Download()
	})

	t.Run("TestResponseTest::testAssertDownloadOfferedWithAFileName", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)

		packagetesting.AssertResponse(t, rec).Download("report.csv")
	})

	t.Run("TestResponseTest::testAssertDownloadOfferedWorksWithBinaryFileResponse", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Content-Disposition", `attachment; filename="binary.bin"`)

		packagetesting.AssertResponse(t, rec).Download("binary.bin")
	})

	t.Run("TestResponseTest::testAssertDownloadOfferedFailsWithInlineContentDisposition", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Content-Disposition", `inline; filename="report.csv"`)

		m := &mockTB{}

		packagetesting.AssertResponse(m, rec).Download("report.csv")

		if !m.failed {
			t.Fatal("expected download assertion to fail")
		}
	})
}

func TestResponseBodyVisibilityAssertions(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertSee", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("hello world")

		packagetesting.AssertResponse(t, rec).See("hello")
	})

	t.Run("TestResponseTest::testAssertSeeEscaped", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("1 &lt; 2")

		packagetesting.AssertResponse(t, rec).SeeEscaped("1 < 2")
	})

	t.Run("TestResponseTest::testAssertSeeHtml", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<strong>hello</strong>")

		packagetesting.AssertResponse(t, rec).SeeHtml("<strong>hello</strong>")
	})

	t.Run("TestResponseTest::testAssertSeeInOrder", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("alpha beta gamma")

		packagetesting.AssertResponse(t, rec).SeeInOrder("alpha", "beta", "gamma")
	})

	t.Run("TestResponseTest::testAssertSeeHtmlInOrder", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<p>alpha</p><p>beta</p><p>gamma</p>")

		packagetesting.AssertResponse(t, rec).SeeHtmlInOrder("<p>alpha</p>", "<p>beta</p>", "<p>gamma</p>")
	})

	t.Run("TestResponseTest::testAssertSeeText", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>Hello <strong>world</strong></div>")

		packagetesting.AssertResponse(t, rec).SeeText("Hello world")
	})

	t.Run("TestResponseTest::testAssertSeeTextWhitespace", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>Hello\n    world</div>")

		packagetesting.AssertResponse(t, rec).SeeText("Hello world")
	})

	t.Run("TestResponseTest::testAssertSeeTextEscaped", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>1 &lt; 2</div>")

		packagetesting.AssertResponse(t, rec).SeeText("1 < 2")
	})

	t.Run("TestResponseTest::testAssertSeeTextInOrder", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>alpha <strong>beta</strong> gamma</div>")

		packagetesting.AssertResponse(t, rec).SeeTextInOrder("alpha", "beta", "gamma")
	})

	t.Run("TestResponseTest::testAssertSeeTextInOrderWhitespace", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>alpha\n<strong>beta</strong>\n gamma</div>")

		packagetesting.AssertResponse(t, rec).SeeTextInOrder("alpha", "beta", "gamma")
	})

	t.Run("TestResponseTest::testAssertSeeTextInOrderEscaped", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>1 &lt; 2 <strong>and</strong> 3 &gt; 2</div>")

		packagetesting.AssertResponse(t, rec).SeeTextInOrder("1 < 2", "3 > 2")
	})

	t.Run("TestResponseTest::testAssertDontSee", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("hello world")

		packagetesting.AssertResponse(t, rec).DontSee("goodbye")
	})

	t.Run("TestResponseTest::testAssertDontSeeEscaped", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("1 &lt; 2")

		packagetesting.AssertResponse(t, rec).DontSeeEscaped("3 < 4")
	})

	t.Run("TestResponseTest::testAssertDontSeeHtml", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<strong>hello</strong>")

		packagetesting.AssertResponse(t, rec).DontSeeHtml("<em>goodbye</em>")
	})

	t.Run("TestResponseTest::testAssertDontSeeText", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>Hello <strong>world</strong></div>")

		packagetesting.AssertResponse(t, rec).DontSeeText("goodbye")
	})

	t.Run("TestResponseTest::testAssertDontSeeTextEscaped", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>1 &lt; 2</div>")

		packagetesting.AssertResponse(t, rec).DontSeeText("3 < 4")
	})

	t.Run("TestResponseTest::testAssertSeeCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("hello world")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).See("goodbye")
		if !m.failed {
			t.Fatal("expected See to fail when content is missing")
		}
	})

	t.Run("TestResponseTest::testAssertSeeEscapedCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("1 &lt; 2")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeEscaped("3 < 4")
		if !m.failed {
			t.Fatal("expected SeeEscaped to fail when escaped content is missing")
		}
	})

	t.Run("TestResponseTest::testAssertSeeHtmlCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<strong>hello</strong>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeHtml("<em>goodbye</em>")
		if !m.failed {
			t.Fatal("expected SeeHtml to fail when HTML fragment is missing")
		}
	})

	t.Run("TestResponseTest::testAssertSeeInOrderCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("alpha gamma beta")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeInOrder("alpha", "beta", "gamma")
		if !m.failed {
			t.Fatal("expected SeeInOrder to fail when order is wrong")
		}
	})

	t.Run("TestResponseTest::testAssertSeeInOrderCanFail2", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("alpha beta")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeInOrder("alpha", "gamma")
		if !m.failed {
			t.Fatal("expected SeeInOrder to fail when a fragment is missing")
		}
	})

	t.Run("TestResponseTest::testAssertSeeHtmlInOrderCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<p>alpha</p><p>gamma</p><p>beta</p>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeHtmlInOrder("<p>alpha</p>", "<p>beta</p>", "<p>gamma</p>")
		if !m.failed {
			t.Fatal("expected SeeHtmlInOrder to fail when order is wrong")
		}
	})

	t.Run("TestResponseTest::testAssertSeeHtmlInOrderCanFail2", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<p>alpha</p><p>beta</p>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeHtmlInOrder("<p>alpha</p>", "<p>gamma</p>")
		if !m.failed {
			t.Fatal("expected SeeHtmlInOrder to fail when a fragment is missing")
		}
	})

	t.Run("TestResponseTest::testAssertSeeTextCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>Hello <strong>world</strong></div>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeText("goodbye")
		if !m.failed {
			t.Fatal("expected SeeText to fail when text is missing")
		}
	})

	t.Run("TestResponseTest::testAssertSeeTextEscapedCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>1 &lt; 2</div>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeText("3 < 4")
		if !m.failed {
			t.Fatal("expected SeeText to fail when escaped text is missing")
		}
	})

	t.Run("TestResponseTest::testAssertSeeTextInOrderCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>alpha gamma beta</div>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeTextInOrder("alpha", "beta", "gamma")
		if !m.failed {
			t.Fatal("expected SeeTextInOrder to fail when order is wrong")
		}
	})

	t.Run("TestResponseTest::testAssertSeeTextInOrderCanFail2", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>alpha beta</div>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).SeeTextInOrder("alpha", "gamma")
		if !m.failed {
			t.Fatal("expected SeeTextInOrder to fail when a fragment is missing")
		}
	})

	t.Run("TestResponseTest::testAssertDontSeeCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("hello world")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).DontSee("world")
		if !m.failed {
			t.Fatal("expected DontSee to fail when the body contains the string")
		}
	})

	t.Run("TestResponseTest::testAssertDontSeeEscapedCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("1 &lt; 2")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).DontSeeEscaped("1 < 2")
		if !m.failed {
			t.Fatal("expected DontSeeEscaped to fail when the body contains the string")
		}
	})

	t.Run("TestResponseTest::testAssertDontSeeHtmlCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<strong>hello</strong>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).DontSeeHtml("<strong>hello</strong>")
		if !m.failed {
			t.Fatal("expected DontSeeHtml to fail when the HTML fragment is present")
		}
	})

	t.Run("TestResponseTest::testAssertDontSeeTextCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>Hello <strong>world</strong></div>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).DontSeeText("world")
		if !m.failed {
			t.Fatal("expected DontSeeText to fail when text is present")
		}
	})

	t.Run("TestResponseTest::testAssertDontSeeTextEscapedCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("<div>1 &lt; 2</div>")

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).DontSeeText("1 < 2")
		if !m.failed {
			t.Fatal("expected DontSeeText to fail when escaped text is present")
		}
	})
}

func TestFluentJSONAssertions(t *testing.T) {
	t.Parallel()

	t.Run("AssertTest::testAssertHas", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("name").
			Assert()
	})

	t.Run("AssertTest::testAssertHasNestedProp", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("user.name").
			Assert()
	})

	t.Run("AssertTest::testAssertHasCountItemsInProp", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Count("items", 2).
			Assert()
	})

	t.Run("AssertTest::testAssertCount", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Count("items", 2).
			Assert()
	})

	t.Run("AssertTest::testAssertBetween", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":30}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Between("age", 18, 65).
			Assert()
	})

	t.Run("AssertTest::testAssertMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("name").
			Missing("email").
			Assert()
	})

	t.Run("AssertTest::testAssertMissingAll", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("name").
			MissingAll("email", "phone").
			Assert()
	})

	t.Run("AssertTest::testAssertWhereMatchesValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Where("name", "Taylor").
			Assert()
	})

	t.Run("AssertTest::testAssertWhereNullMatchesValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":null}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereNull("name").
			Assert()
	})

	t.Run("AssertTest::testAssertWhereNotNullMatchesValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereNotNull("name").
			Assert()
	})

	t.Run("AssertTest::testAssertWhereContainsWithNestedValue", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"roles":["admin","editor"]}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereContains("roles", []any{"admin"}).
			Assert()
	})

	t.Run("AssertTest::testAssertWhereNotUsingClosure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":30}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereNot("age", func(v any) bool {
				n, ok := v.(json.Number)
				if !ok {
					return false
				}

				return n.String() == "18"
			}).
			Assert()
	})

	t.Run("AssertTest::testAssertWhereTypeString", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereType("name", "string").
			Assert()
	})

	t.Run("AssertTest::testAssertWhereTypeArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2]}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereType("items", "array").
			Assert()
	})

	t.Run("AssertTest::testAssertWhereAllMatchesValues", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","age":30}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereAll(map[string]any{"name": "Taylor", "age": 30}).
			Assert()
	})

	t.Run("AssertTest::testAssertWhereAllType", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","age":30,"active":true}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			WhereAllType(map[string]string{"name": "string", "age": "integer", "active": "boolean"}).
			Assert()
	})

	t.Run("AssertTest::testAssertHasAll", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","age":30}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			HasAll("name", "age").
			Assert()
	})

	t.Run("AssertTest::testAssertHasOnlyCounts", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","age":30}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			HasOnly("name", "age").
			Assert()
	})

	t.Run("AssertTest::testAssertCountMultipleProps", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2],"roles":["admin","editor"]}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			CountMultiple(map[string]int{"items": 2, "roles": 2}).
			Assert()
	})

	t.Run("AssertTest::testFailsWhenNotInteractingWithAllPropsInScope", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","email":"taylor@example.test"}`))

		m := &mockTB{}

		packagetesting.AssertResponse(m, rec).FluentJSON().
			Has("name").
			Assert()

		if !m.failed {
			t.Fatal("expected fluent JSON assertion to fail when top-level props remain uninspected")
		}
	})

	t.Run("AssertTest::testDisableInteractionCheckForCurrentScope", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","email":"taylor@example.test"}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Has("name").
			AllowMissingInteractions().
			Assert()
	})

	t.Run("AssertTest::testTopLevelPropInteractionDisabledByDefault", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`[{"name":"Taylor"}]`))

		packagetesting.AssertResponse(t, rec).FluentJSON().Assert()
	})
}

func TestFluentJSONScopeAssertions(t *testing.T) {
	t.Parallel()

	t.Run("AssertTest::testScope", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			Scope("user", func(scoped *packagetesting.AssertableJSON) {
				scoped.Where("name", "Taylor")
			})
	})

	t.Run("AssertTest::testFirstScope", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":[{"name":"Taylor"},{"name":"Alyssa"}]}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().
			FirstScope("users", func(scoped *packagetesting.AssertableJSON) {
				scoped.Where("name", "Taylor")
			})
	})

	t.Run("AssertTest::testEachScope", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"users":[{"name":"Taylor"},{"name":"Alyssa"}]}`))

		names := []string{"Taylor", "Alyssa"}
		i := 0

		packagetesting.AssertResponse(t, rec).FluentJSON().
			EachScope("users", func(scoped *packagetesting.AssertableJSON) {
				scoped.Where("name", names[i])
				i++
			})
	})
}

func TestResponseCoreAndJSONAssertions(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertStatus", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusTeapot)

		packagetesting.AssertResponse(t, rec).Status(http.StatusTeapot)
	})

	t.Run("TestResponseTest::testAssertRedirect", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/dashboard")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).Redirect("/dashboard").Location("/dashboard")
	})

	t.Run("TestResponseTest::testAssertRedirectBack", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/previous")
		rec.Header().Set("Referer", "/previous")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).RedirectBack().Location("/previous")
	})

	t.Run("TestResponseTest::testAssertRedirectContains", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/login?next=/dashboard")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).RedirectContains("/login")
	})

	t.Run("TestResponseTest::testAssertLocation", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("Location", "/welcome")
		rec.WriteHeader(http.StatusFound)

		packagetesting.AssertResponse(t, rec).Location("/welcome")
	})

	t.Run("TestResponseTest::testAssertBodyEquals", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("hello")

		packagetesting.AssertResponse(t, rec).BodyEquals("hello")
	})

	t.Run("TestResponseTest::testAssertBodyContains", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteString("hello world")

		packagetesting.AssertResponse(t, rec).BodyContains("world")
	})

	t.Run("TestResponseTest::testAssertJsonWithArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`["alpha","beta"]`))

		packagetesting.AssertResponse(t, rec).JSON([]string{"alpha", "beta"})
	})

	t.Run("TestResponseTest::testAssertJsonWithNull", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`null`))

		packagetesting.AssertResponse(t, rec).JSON(nil)
	})

	t.Run("TestResponseTest::testAssertJsonWithFluent", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		packagetesting.AssertResponse(t, rec).
			FluentJSON().
			Has("user.name").
			Assert()
	})

	t.Run("TestResponseTest::testAssertJsonWithFluentHasAnyPasses", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		packagetesting.AssertResponse(t, rec).
			FluentJSON().
			HasAny("missing", "user.name").
			Assert()
	})

	t.Run("TestResponseTest::testAssertJsonWithFluentHasAnyThrows", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		m := &mockTB{}

		packagetesting.AssertResponse(m, rec).
			FluentJSON().
			HasAny("missing", "still-missing").
			Assert()

		if !m.failed {
			t.Fatal("expected HasAny to fail when no JSON path exists")
		}
	})

	t.Run("TestResponseTest::testAssertJsonWithFluentFailsWhenNotInteractingWithAllProps", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","email":"taylor@example.test"}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).
			FluentJSON().
			Has("name").
			Assert()
		if !m.failed {
			t.Fatal("expected fluent JSON to fail when a top-level prop is not inspected")
		}
	})

	t.Run("TestResponseTest::testAssertJsonWithFluentSkipsInteractionWhenTopLevelKeysNonAssociative", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`[{"name":"Taylor"},{"name":"Alyssa"}]`))

		packagetesting.AssertResponse(t, rec).
			FluentJSON().
			Assert()
	})

	t.Run("TestResponseTest::testAssertSimilarJsonWithMixed", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","age":30}`))

		packagetesting.AssertResponse(t, rec).SimilarJSON(map[string]any{
			"name": "Taylor",
			"age":  30,
		})
	})

	t.Run("TestResponseTest::testAssertExactJsonWithMixedWhenDataIsExactlySame", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","age":30}`))

		packagetesting.AssertResponse(t, rec).ExactJSON(map[string]any{
			"name": "Taylor",
			"age":  30,
		})
	})

	t.Run("TestResponseTest::testAssertExactJsonWithMixedWhenDataIsSimilar", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor","age":30}`))

		packagetesting.AssertResponse(t, rec).ExactJSON(map[string]any{
			"age":  30,
			"name": "Taylor",
		})
	})

	t.Run("TestResponseTest::testAssertJsonPath", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		packagetesting.AssertResponse(t, rec).JSONPath("user.name", "Taylor")
	})

	t.Run("TestResponseTest::testAssertJsonPathWithClosure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"age":30}`))

		packagetesting.AssertResponse(t, rec).JSONPath("age", func(v any) bool {
			n, ok := v.(json.Number)
			return ok && n.String() == "30"
		})
	})

	t.Run("TestResponseTest::testAssertJsonPathCanonicalizing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"roles":["editor","admin"]}`))

		packagetesting.AssertResponse(t, rec).JSONPathCanonicalizing("roles", []any{"admin", "editor"})
	})

	t.Run("TestResponseTest::testAssertJsonFragment", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","roles":["admin"]}}`))

		packagetesting.AssertResponse(t, rec).JSONFragment(map[string]any{"name": "Taylor"})
	})

	t.Run("TestResponseTest::testAssertJsonMissingPath", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		packagetesting.AssertResponse(t, rec).JSONMissingPath("user.email")
	})

	t.Run("TestResponseTest::testAssertJsonIsArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`[]`))

		packagetesting.AssertResponse(t, rec).JSONIsArray()
	})

	t.Run("TestResponseTest::testAssertJsonIsObject", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{}`))

		packagetesting.AssertResponse(t, rec).JSONIsObject()
	})

	t.Run("TestResponseTest::testJsonHelper", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		value := packagetesting.AssertResponse(t, rec).JSONValue()
		if got := value.(map[string]any)["name"]; got != "Taylor" {
			t.Fatalf("expected decoded JSON helper to return Taylor, got %v", got)
		}
	})

	t.Run("TestResponseTest::testResponseCanBeReturnedAsCollection", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`["alpha","beta"]`))

		got := packagetesting.AssertResponse(t, rec).JSONCollection()
		if got.Count() != 2 {
			t.Fatalf("expected 2 collection items, got %d", got.Count())
		}

		if got.All()[0] != "alpha" || got.All()[1] != "beta" {
			t.Fatalf("unexpected collection contents: %v", got.All())
		}
	})

	t.Run("TestResponseTest::testItCanBeTapped", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"name":"Taylor"}`))

		var tapped bool
		packagetesting.AssertResponse(t, rec).Tap(func(a *packagetesting.Assertions) {
			tapped = true
			if a == nil {
				t.Fatal("expected assertion helper in Tap callback")
			}
		})

		if !tapped {
			t.Fatal("expected Tap callback to run")
		}
	})
}

func TestResponseAdditionalAssertions(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertContent", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte("hello world"))

		packagetesting.AssertResponse(t, rec).BodyEquals("hello world")
	})

	t.Run("TestResponseTest::testAssertInternalServerError", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusInternalServerError)

		packagetesting.AssertResponse(t, rec).Status(http.StatusInternalServerError)
	})

	t.Run("TestResponseTest::testAssertServerError", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusInternalServerError)

		packagetesting.AssertResponse(t, rec).Status(http.StatusInternalServerError)
	})

	t.Run("TestResponseTest::testAssertServiceUnavailable", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusServiceUnavailable)

		packagetesting.AssertResponse(t, rec).Status(http.StatusServiceUnavailable)
	})

	t.Run("TestResponseTest::testAssertUnprocessable", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusUnprocessableEntity)

		packagetesting.AssertResponse(t, rec).Unprocessable()
	})

	t.Run("TestResponseTest::testAssertJsonCount", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[1,2,3]}`))

		packagetesting.AssertResponse(t, rec).FluentJSON().Count("items", 3)
	})

	t.Run("TestResponseTest::testAssertJsonFragments", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","email":"taylor@example.test"}}`))

		assert := packagetesting.AssertResponse(t, rec)
		assert.JSONFragment(map[string]any{"name": "Taylor"})
		assert.JSONFragment(map[string]any{"email": "taylor@example.test"})
	})

	t.Run("TestResponseTest::testAssertJsonMissing", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor"}}`))

		packagetesting.AssertResponse(t, rec).JSONMissingPath("user.email")
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsWithoutArgument", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{}`))

		packagetesting.AssertResponse(t, rec).JSONMissingValidationErrors()
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsWithoutArgumentWhenErrorsIsEmpty", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":[]}`))

		packagetesting.AssertResponse(t, rec).JSONMissingValidationErrors()
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsCustomErrorsName", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"meta":{"errors":{"email":["required"]}}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrorsAt("meta.errors", map[string][]string{
			"email": []string{"required"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsCustomNestedErrorsName", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"meta":{"validation":{"errors":{"email":["required"]}}}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrorsAt("meta.validation.errors", map[string][]string{
			"email": []string{"required"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsCustomErrorsName", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"meta":{"errors":[]}}`))

		packagetesting.AssertResponse(t, rec).JSONMissingValidationErrorsAt("meta.errors")
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsNestedCustomErrorsName1", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"meta":{"validation":{"errors":[]}}}`))

		packagetesting.AssertResponse(t, rec).JSONMissingValidationErrorsAt("meta.validation.errors")
	})

	t.Run("TestResponseTest::testAssertJsonMissingValidationErrorsNestedCustomErrorsName2", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"meta":{"validation":{"errors":{}}}}`))

		packagetesting.AssertResponse(t, rec).JSONMissingValidationErrorsAt("meta.validation.errors")
	})

	t.Run("TestResponseTest::testAssertExactJsonStructure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","roles":["admin"]}}`))

		packagetesting.AssertResponse(t, rec).ExactJSON(map[string]any{
			"user": map[string]any{
				"name":  "Taylor",
				"roles": []any{"admin"},
			},
		})
	})

	t.Run("TestResponseTest::testAssertJsonStructure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","roles":["admin"]}}`))

		packagetesting.AssertResponse(t, rec).JSONStructure(map[string]any{
			"user": map[string]any{
				"name":  "Taylor",
				"roles": []any{"admin"},
			},
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingExact", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","roles":["admin"]}}`))

		packagetesting.AssertResponse(t, rec).ExactJSON(map[string]any{
			"user": map[string]any{
				"name":  "Taylor",
				"roles": []any{"admin"},
			},
		})
	})

	t.Run("TestResponseTest::testAssertJsonMissingExactCanFail2", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"user":{"name":"Taylor","roles":["admin"]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).ExactJSON(map[string]any{
			"user": map[string]any{"name": "Taylor"},
		})
		if !m.failed {
			t.Fatal("expected ExactJSON to fail when keys are missing")
		}
	})

	t.Run("TestResponseTest::testAssertOnlyJsonValidationErrors", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"required"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsUsingAssertOnlyInvalid", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"required"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsUsingAssertInvalid", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"],"name":["too short"]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"required"},
			"name":  []string{"too short"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsFailsWhenGivenAnEmptyArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{})
		if !m.failed {
			t.Fatal("expected empty validation error expectation to fail")
		}
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorsWithArray", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required","email"]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"required", "email"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessages", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required."]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email field is required."},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorContainsMessages", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required.","The email must be valid."]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email field is required.", "The email must be valid."},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessagesCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required."]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email must be valid."},
		})
		if !m.failed {
			t.Fatal("expected validation message mismatch to fail")
		}
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessageKeyCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required."]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{
			"name": []string{"The name field is required."},
		})
		if !m.failed {
			t.Fatal("expected validation key mismatch to fail")
		}
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessagesMultipleMessages", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required.","The email must be valid."]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email field is required.", "The email must be valid."},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessagesMultipleMessagesCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required.","The email must be valid."]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email field is required."},
		})
		if !m.failed {
			t.Fatal("expected missing validation message to fail")
		}
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessagesMixed", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required."],"name":["The name must be a string."]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email field is required."},
			"name":  []string{"The name must be a string."},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessagesMixedCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required."],"name":["The name must be a string."]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email field is required."},
			"name":  []string{"The name field is required."},
		})
		if !m.failed {
			t.Fatal("expected mixed validation message mismatch to fail")
		}
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessagesMultipleErrors", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"],"name":["required"],"password":["min"]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email":    []string{"required"},
			"name":     []string{"required"},
			"password": []string{"min"},
		})
	})

	t.Run("TestResponseTest::testAssertJsonValidationErrorMessagesMultipleErrorsCanFail", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"],"name":["required"]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{
			"email":    []string{"required"},
			"name":     []string{"required"},
			"password": []string{"min"},
		})
		if !m.failed {
			t.Fatal("expected missing validation error key to fail")
		}
	})

	t.Run("TestResponseTest::testAssertSessionOnlyValidationErrorsUsingAssertOnlyInvalid", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"errors": map[string][]string{"email": []string{"required"}}}).
			SessionHasErrors(map[string][]string{"email": []string{"required"}})
	})

	t.Run("TestResponseTest::testAssertSessionValidationErrorsUsingAssertInvalid", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"errors": map[string][]string{"email": []string{"required"}, "name": []string{"required"}}}).
			SessionHasErrors(map[string][]string{"email": []string{"required"}, "name": []string{"required"}})
	})

	t.Run("TestResponseTest::testAssertSessionValidationErrorsUsingAssertValid", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{}).
			SessionHasNoErrors()
	})

	t.Run("TestResponseTest::testAssertingKeyIsInvalidErrorMessage", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["The email field is required."]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"The email field is required."},
		})
	})

	t.Run("TestResponseTest::testInvalidWithListOfErrors", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"],"password":["min"]}}`))

		packagetesting.AssertResponse(t, rec).JSONValidationErrors(map[string][]string{
			"email":    []string{"required"},
			"password": []string{"min"},
		})
	})

	t.Run("TestResponseTest::testItHandlesFalseJson", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`false`))

		if got := packagetesting.AssertResponse(t, rec).JSONValue(); got != false {
			t.Fatalf("JSONValue = %#v, want false", got)
		}
		packagetesting.AssertResponse(t, rec).JSON(false)
	})

	t.Run("TestResponseTest::testItHandlesEncodedJson", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`"{\"ok\":true}"`))

		if got := packagetesting.AssertResponse(t, rec).JSONValue(); got != `{"ok":true}` {
			t.Fatalf("JSONValue = %#v, want encoded JSON string", got)
		}
	})

	t.Run("TestResponseTest::testCanBeCreatedFromBinaryFileResponses", func(t *testing.T) {
		rec := newRecorder()
		rec.Header().Set("X-Streamed", "true")
		rec.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)
		rec.Write([]byte("id,name\n1,Taylor\n"))

		packagetesting.AssertResponse(t, rec).StreamedBinaryFile("report.csv")
	})

	t.Run("TestResponseTest::testAssertSessionCookieExpiredDoesNotTriggerOnSessionCookies", func(t *testing.T) {
		rec := newRecorder()
		http.SetCookie(rec, &http.Cookie{Name: "session", Value: "session-value"})

		packagetesting.AssertResponse(t, rec).CookieNotExpired("session")
	})

	t.Run("TestResponseTest::testAssertSessionCookieNotExpired", func(t *testing.T) {
		rec := newRecorder()
		http.SetCookie(rec, &http.Cookie{
			Name:    "session",
			Value:   "session-value",
			Expires: time.Now().Add(24 * time.Hour),
		})

		packagetesting.AssertResponse(t, rec).CookieNotExpired("session")
	})
}

func TestFluentJSONAdditionalAssertions(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		run  func(t *testing.T)
	}

	tests := []testCase{
		{
			name: "AssertTest::testAssertMissingAllAcceptsMultipleArgumentsInsteadOfArray",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"name":"Taylor"}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().
					MissingAll("email", "roles.0").
					AllowMissingInteractions().
					Assert()
			},
		},
		{
			name: "AssertTest::testAssertWhereUsingClosure",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"age":30}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().
					Where("age", func(v any) bool {
						n, ok := v.(json.Number)
						return ok && n.String() == "30"
					}).
					Assert()
			},
		},
		{
			name: "AssertTest::testAssertHasAllAcceptsMultipleArgumentsInsteadOfArray",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"user":{"name":"Taylor","email":"taylor@example.test"}}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().HasAll("user.name", "user.email")
			},
		},
		{
			name: "AssertTest::testAssertHasWithWhereNotDoesNotFail",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"name":"Taylor"}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereNot("name", "Alyssa")
			},
		},
		{
			name: "AssertTest::testAssertHasWithWhereNotDoesNotFailClosure",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"age":30}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereNot("age", func(v any) bool {
					n, ok := v.(float64)
					return ok && n < 18
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereDoesNotMatchValue",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"name":"Taylor"}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereNot("name", "Alyssa")
			},
		},
		{
			name: "AssertTest::testAssertWhereMatchesValueUsingArrayable",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"user":{"name":"Taylor","age":30}}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().Where("user", map[string]any{
					"name": "Taylor",
					"age":  30,
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereMatchesValueUsingArrayableWhenSortedDifferently",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"user":{"name":"Taylor","age":30}}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().Where("user", map[string]any{
					"age":  30,
					"name": "Taylor",
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereClosureArrayValuesAreAutomaticallyCastedToCollections",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[{"name":"Taylor"},{"name":"Alyssa"}]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().
					Where("items", func(v any) bool {
						items, ok := v.([]any)
						return ok && len(items) == 2
					})
			},
		},
		{
			name: "AssertTest::testAssertWhereFailsWhenDoesNotMatchValueUsingArrayable",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"user":{"name":"Taylor","age":30}}`))

				m := &mockTB{}
				packagetesting.AssertResponse(m, rec).FluentJSON().
					Where("user", map[string]any{"name": "Alyssa", "age": 30})
				if !m.failed {
					t.Fatal("expected map value mismatch to fail")
				}
			},
		},
		{
			name: "AssertTest::testAssertNestedWhereMatchesValue",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"user":{"profile":{"name":"Taylor"}}}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().
					Where("user.profile.name", "Taylor").
					Assert()
			},
		},
		{
			name: "AssertTest::testAssertNestedWhereFailsWhenDoesNotMatchValue",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"user":{"profile":{"name":"Taylor"}}}`))

				m := &mockTB{}
				packagetesting.AssertResponse(m, rec).FluentJSON().
					Where("user.profile.name", "Alyssa")
				if !m.failed {
					t.Fatal("expected nested Where mismatch to fail")
				}
			},
		},
		{
			name: "AssertTest::testAssertWhereAllMatchesValues",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"name":"Taylor","age":30}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereAll(map[string]any{
					"name": "Taylor",
					"age":  30,
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereAllFailsWhenAtLeastOnePropDoesNotMatchValue",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"name":"Taylor","age":30}`))

				m := &mockTB{}
				packagetesting.AssertResponse(m, rec).FluentJSON().
					WhereAll(map[string]any{"name": "Taylor", "age": 31})
				if !m.failed {
					t.Fatal("expected WhereAll mismatch to fail")
				}
			},
		},
		{
			name: "AssertTest::testAssertWhereAllType",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"name":"Taylor","age":30,"active":true}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereAllType(map[string]string{
					"name":   "string",
					"age":    "integer",
					"active": "boolean",
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereTypeBoolean",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"active":true}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereType("active", "boolean")
			},
		},
		{
			name: "AssertTest::testAssertWhereTypeInteger",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"age":30}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereType("age", "integer")
			},
		},
		{
			name: "AssertTest::testAssertWhereTypeDouble",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"score":1.5}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereType("score", "double")
			},
		},
		{
			name: "AssertTest::testAssertWhereTypeNull",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"deleted":null}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereType("deleted", "null")
			},
		},
		{
			name: "AssertTest::testAssertWhereTypeWithUnionTypes",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"flag":true}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereType("flag", "string|boolean")
			},
		},
		{
			name: "AssertTest::testAssertWhereTypeWithPipeInUnionType",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"flag":true}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereType("flag", "boolean|string")
			},
		},
		{
			name: "AssertTest::testAssertCountMultipleProps",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[1,2],"roles":[{"name":"admin"}]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().CountMultiple(map[string]int{
					"items": 2,
					"roles": 1,
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsFailsWithEmptyValue",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":["admin"]}`))

				m := &mockTB{}
				packagetesting.AssertResponse(m, rec).FluentJSON().WhereContains("items", []any{})
				if !m.failed {
					t.Fatal("expected empty WhereContains value to fail")
				}
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsFailsWhenSatisfiesClosureButDoesNotHaveExpectedValue",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[{"name":"Taylor"}]}`))

				m := &mockTB{}
				packagetesting.AssertResponse(m, rec).FluentJSON().WhereContains("items", func(v any) bool {
					item, ok := v.(map[string]any)
					return ok && item["name"] == "Alyssa"
				})
				if !m.failed {
					t.Fatal("expected closure value mismatch to fail")
				}
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithOutOfOrderMatchingType",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":["admin","editor","viewer"]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", []any{"viewer", "admin"})
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithOutOfOrderNestedMatchingType",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[{"name":"Taylor"},{"name":"Alyssa"},{"name":"Nuno"}]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", []any{
					map[string]any{"name": "Nuno"},
					map[string]any{"name": "Taylor"},
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithClosure",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[{"name":"Taylor"},{"name":"Alyssa"}]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", func(v any) bool {
					item, ok := v.(map[string]any)
					return ok && item["name"] == "Taylor"
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithNestedClosure",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[{"profile":{"name":"Taylor"}},{"profile":{"name":"Alyssa"}}]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", func(v any) bool {
					item, ok := v.(map[string]any)
					if !ok {
						return false
					}
					profile, ok := item["profile"].(map[string]any)
					return ok && profile["name"] == "Taylor"
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithMultipleClosure",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[{"name":"Taylor","active":true},{"name":"Alyssa","active":false}]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", func(v any) bool {
					item, ok := v.(map[string]any)
					return ok && item["name"] == "Taylor" && item["active"] == true
				})
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsFailsWhenDoesNotSatisfyClosure",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[{"name":"Taylor"}]}`))

				m := &mockTB{}
				packagetesting.AssertResponse(m, rec).FluentJSON().WhereContains("items", func(v any) bool {
					item, ok := v.(map[string]any)
					return ok && item["name"] == "Alyssa"
				})
				if !m.failed {
					t.Fatal("expected WhereContains closure mismatch to fail")
				}
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithMatchingType",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[1,2,3]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", []any{1, 2})
			},
		},
		{
			name: "AssertTest::testCannotDisableInteractionCheckForDifferentScopes",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"user":{"name":"Taylor"},"meta":{"page":1}}`))

				m := &mockTB{}
				packagetesting.AssertResponse(m, rec).FluentJSON().
					Scope("user", func(scoped *packagetesting.AssertableJSON) {
						scoped.AllowMissingInteractions()
					}).
					Assert()
				if !m.failed {
					t.Fatal("expected parent scope interaction check to remain enabled")
				}
			},
		},
		{
			name: "AssertTest::testTappable",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"name":"Taylor"}`))

				var tapped bool
				packagetesting.AssertResponse(t, rec).FluentJSON().
					Tap(func(assert *packagetesting.AssertableJSON) {
						tapped = true
						assert.Where("name", "Taylor")
					}).
					Assert()
				if !tapped {
					t.Fatal("expected fluent JSON Tap callback")
				}
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithNullValue",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[null,1]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", nil)
			},
		},
		{
			name: "AssertTest::testAssertWhereContainsWithNullExpectation",
			run: func(t *testing.T) {
				rec := newRecorder()
				rec.Write([]byte(`{"items":[null,1]}`))

				packagetesting.AssertResponse(t, rec).FluentJSON().WhereContains("items", nil)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}

func TestResponseRemainingInventoryAssertions(t *testing.T) {
	t.Parallel()

	t.Run("TestResponseTest::testAssertJsonSerializedSessionHasErrors", func(t *testing.T) {
		rec := newRecorder()

		packagetesting.AssertResponse(t, rec).
			WithSession(map[string]any{"errors": `{"email":["required"]}`}).
			SessionHasErrors(map[string][]string{"email": []string{"required"}})
	})

	t.Run("TestResponseTest::testHandledExceptionIsIncludedInAssertionFailure", func(t *testing.T) {
		rec := newRecorder()
		rec.WriteHeader(http.StatusInternalServerError)

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).Status(http.StatusOK)
		if !m.failed || !strings.Contains(strings.Join(m.messages, "\n"), "500") {
			t.Fatalf("expected status failure to include handled status, got %#v", m.messages)
		}
	})

	t.Run("TestResponseTest::testValidationErrorsAreIncludedInAssertionFailure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"errors":{"email":["required"]}}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSONValidationErrors(map[string][]string{
			"email": []string{"email"},
		})
		if !m.failed || !strings.Contains(strings.Join(m.messages, "\n"), "required") {
			t.Fatalf("expected validation failure to include errors, got %#v", m.messages)
		}
	})

	t.Run("TestResponseTest::testJsonErrorsAreIncludedInAssertionFailure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`not-json`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).JSON(map[string]any{})
		if !m.failed || !strings.Contains(strings.Join(m.messages, "\n"), "invalid") {
			t.Fatalf("expected JSON failure to include decode error, got %#v", m.messages)
		}
	})

	t.Run("AssertTest::testAssertWhereContainsFailsWhenHavingExpectedValueButDoesNotSatisfyClosure", func(t *testing.T) {
		rec := newRecorder()
		rec.Write([]byte(`{"items":[{"name":"Taylor","active":false}]}`))

		m := &mockTB{}
		packagetesting.AssertResponse(m, rec).FluentJSON().WhereContains("items", func(v any) bool {
			item, ok := v.(map[string]any)
			return ok && item["name"] == "Taylor" && item["active"] == true
		})
		if !m.failed {
			t.Fatal("expected closure predicate to fail when extra expectation is not satisfied")
		}
	})
}
