package access_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bedrock/packages/auth/access"
)

// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testAllowMethod
func TestLaravelAuthAccessResponseAllowMethod(t *testing.T) {
	resp := access.Allow("allowed")

	if !resp.Allowed {
		t.Fatal("Allow should create an allowed response")
	}

	if resp.Message != "allowed" {
		t.Errorf("Message = %q, want %q", resp.Message, "allowed")
	}
}

// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testDenyMethod
// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testDenyMethodWithNoMessageReturnsNull
func TestLaravelAuthAccessResponseDenyMethod(t *testing.T) {
	resp := access.Deny("", 0)

	if resp.Allowed {
		t.Fatal("Deny should create a denied response")
	}

	if resp.Message != "" {
		t.Errorf("Message = %q, want empty", resp.Message)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
}

// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testItSetsEmptyStatusOnExceptionWhenAuthorizing
// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testItSetsStatusOnExceptionWhenAuthorizing
// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testAuthorizeMethodThrowsAuthorizationExceptionWhenResponseDenied
// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testAuthorizeMethodThrowsAuthorizationExceptionWithDefaultMessage
// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testThrowIfNeededDoesntThrowAuthorizationExceptionWhenResponseAllowed
func TestLaravelAuthAccessResponseAuthorize(t *testing.T) {
	if err := access.Allow("ok").Authorize(); err != nil {
		t.Fatalf("allowed response should not error: %v", err)
	}

	err := access.Deny("no", http.StatusTeapot).Authorize()

	if err == nil {
		t.Fatal("denied response should authorize with an error")
	}

	var authErr *access.AuthorizationException

	if !errors.As(err, &authErr) {
		t.Fatalf("error = %T, want AuthorizationException", err)
	}

	if authErr.Response.StatusCode != http.StatusTeapot {
		t.Errorf("StatusCode = %d, want %d", authErr.Response.StatusCode, http.StatusTeapot)
	}
}

// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testCastingToStringReturnsMessage
func TestLaravelAuthAccessResponseStringReturnsMessage(t *testing.T) {
	resp := access.Deny("not allowed", 0)

	if resp.String() != "not allowed" {
		t.Errorf("String() = %q, want %q", resp.String(), "not allowed")
	}
}

// Port of Illuminate\Tests\Auth\AuthAccessResponseTest::testResponseToArrayMethod
func TestLaravelAuthAccessResponseToMap(t *testing.T) {
	resp := access.DenyWithStatus(http.StatusNotFound, "missing")
	got := resp.ToMap()

	if got["allowed"] != false || got["message"] != "missing" || got["status"] != http.StatusNotFound {
		t.Errorf("ToMap() = %#v", got)
	}
}
