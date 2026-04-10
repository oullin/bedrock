package authflows

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- test broker ---

type testBroker struct {
	sentLink    bool
	resetCalled bool
	returnError error
}

// --- test responder extensions ---

type testPasswordResponder struct {
	testResponder
	linkSentCalled bool
	resetCalled    bool
}

func (b *testBroker) SendResetLink(_ context.Context, _ map[string]string) error {
	b.sentLink = true

	return b.returnError
}

func (b *testBroker) Reset(_ context.Context, _ map[string]string, callback func(Authenticatable, string) error) error {
	if b.returnError != nil {
		return b.returnError
	}

	b.resetCalled = true

	return callback(&testUser{id: "1"}, "newpassword")
}

func (r *testPasswordResponder) PasswordResetLinkSentResponse(w http.ResponseWriter, _ *http.Request) {
	r.linkSentCalled = true
	w.WriteHeader(http.StatusOK)
}

func (r *testPasswordResponder) PasswordResetResponse(w http.ResponseWriter, _ *http.Request) {
	r.resetCalled = true
	w.WriteHeader(http.StatusOK)
}

// --- helpers ---

func forgotPasswordRequest(email string) *http.Request {
	body := strings.NewReader("email=" + email)
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.RemoteAddr = "127.0.0.1:1234"

	return req
}

func resetPasswordRequest(email string, token string, password string) *http.Request {
	body := strings.NewReader("email=" + email + "&token=" + token + "&password=" + password + "&password_confirmation=" + password)
	req := httptest.NewRequest(http.MethodPost, "/reset-password", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

// --- tests ---

func TestForgotPasswordSuccess(t *testing.T) {
	broker := &testBroker{}
	responder := &testPasswordResponder{}

	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.ResetPasswords = true
	f.broker = broker
	f.responder = responder

	handler := NewForgotPasswordHandler(f)
	w := httptest.NewRecorder()
	r := forgotPasswordRequest("user@example.com")

	handler.ServeHTTP(w, r)

	if !broker.sentLink {
		t.Fatal("expected broker.SendResetLink to be called")
	}

	if !responder.linkSentCalled {
		t.Fatal("expected link sent response")
	}
}

func TestForgotPasswordDisabled(t *testing.T) {
	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.ResetPasswords = false

	handler := NewForgotPasswordHandler(f)
	w := httptest.NewRecorder()
	r := forgotPasswordRequest("user@example.com")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestForgotPasswordBrokerFails(t *testing.T) {
	broker := &testBroker{returnError: errors.New("user not found")}
	responder := &testPasswordResponder{}

	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.ResetPasswords = true
	f.broker = broker
	f.responder = responder

	handler := NewForgotPasswordHandler(f)
	w := httptest.NewRecorder()
	r := forgotPasswordRequest("unknown@example.com")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestResetPasswordSuccess(t *testing.T) {
	broker := &testBroker{}
	resets := &stubResets{}
	events := &testEvents{}
	responder := &testPasswordResponder{}

	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, events, &responder.testResponder)
	f.config.Features.ResetPasswords = true
	f.broker = broker
	f.resetPass = resets
	f.responder = responder

	handler := NewResetPasswordHandler(f)
	w := httptest.NewRecorder()
	r := resetPasswordRequest("user@example.com", "valid-token", "newpassword")

	handler.ServeHTTP(w, r)

	if !broker.resetCalled {
		t.Fatal("expected broker.Reset to be called")
	}

	if !responder.resetCalled {
		t.Fatal("expected reset response")
	}

	if len(events.dispatched) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events.dispatched))
	}

	if events.dispatched[0].Name != EventPasswordReset {
		t.Fatalf("expected PasswordReset event, got %s", events.dispatched[0].Name)
	}
}

func TestResetPasswordDisabled(t *testing.T) {
	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.ResetPasswords = false

	handler := NewResetPasswordHandler(f)
	w := httptest.NewRecorder()
	r := resetPasswordRequest("user@example.com", "token", "pass")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestResetPasswordBrokerFails(t *testing.T) {
	broker := &testBroker{returnError: errors.New("invalid token")}
	responder := &testPasswordResponder{}

	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.ResetPasswords = true
	f.broker = broker
	f.resetPass = &stubResets{}
	f.responder = responder

	handler := NewResetPasswordHandler(f)
	w := httptest.NewRecorder()
	r := resetPasswordRequest("user@example.com", "bad-token", "pass")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}
