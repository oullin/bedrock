package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/auth/events"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- SessionGuard: User resolution ---

func TestSessionGuardUserReturnsNilWhenNoUserFound(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	u, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if u != nil {
		t.Error("expected nil user when no session and no cookie")
	}
}

func TestSessionGuardUserReturnsCachedUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)

	// Calling User() multiple times should return the same cached user.
	u1, _ := guard.User(ctx)
	u2, _ := guard.User(ctx)

	if u1 != u2 {
		t.Error("User() should return cached user on subsequent calls")
	}
}

func TestSessionGuardUserIsSetToRetrievedUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	sess.Put("_auth_user", "1")
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Error("expected user from session to match provider user")
	}
}

func TestSessionGuardUserUsesRememberCookieIfItExists(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	user.SetRememberToken("recaller")
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "web_remember", Value: "1|recaller|hash"})
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Fatal("expected user from remember cookie")
	}

	if !guard.ViaRemember(context.Background()) {
		t.Error("expected ViaRemember() to be true")
	}
}

func TestSessionGuardUserDispatchesAuthenticatedEventFromSession(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	sess.Put("_auth_user", "1")
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	_, _ = guard.User(context.Background())

	dispatcher.has(t, "events.Authenticated")
}

func TestSessionGuardUserDispatchesLoginEventFromRememberCookie(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	user.SetRememberToken("tok")
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "web_remember", Value: "1|tok|hash"})
	guard.SetRequest(req)

	_, _ = guard.User(context.Background())

	dispatcher.has(t, "events.Login")
}

// --- SessionGuard: Check / Guest / ID ---

func TestSessionGuardCheckReturnsTrueWhenUserIsNotNull(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)

	if !guard.Check(ctx) {
		t.Error("Check() should return true after login")
	}

	if guard.Guest(ctx) {
		t.Error("Guest() should return false after login")
	}
}

func TestSessionGuardCheckReturnsFalseWhenUserIsNull(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	if guard.Check(ctx) {
		t.Error("Check() should return false with no user")
	}

	if !guard.Guest(ctx) {
		t.Error("Guest() should return true with no user")
	}
}

func TestSessionGuardIDReturnsIdentifier(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "42", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"42": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)

	if guard.ID(ctx) != "42" {
		t.Errorf("ID() = %v, want \"42\"", guard.ID(ctx))
	}
}

func TestSessionGuardIDReturnsNilWhenNoUser(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	if guard.ID(context.Background()) != nil {
		t.Error("ID() should return nil when no user")
	}
}

// --- SessionGuard: Attempt ---

func TestSessionGuardAttemptCallsRetrieveByCredentials(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	ok := guard.Attempt(context.Background(), map[string]string{"email": "foo@bar.com", "password": "secret"}, false)

	if ok {
		t.Error("Attempt should fail when user not found")
	}

	dispatcher.has(t, "events.Attempting")
	dispatcher.has(t, "events.Failed")
	dispatcher.hasNot(t, "events.Validated")
}

func TestSessionGuardAttemptReturnsTrue(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	ok := guard.Attempt(context.Background(), map[string]string{"email": "a@b.com", "password": "pw"}, false)

	if !ok {
		t.Error("Attempt should return true with valid credentials")
	}

	dispatcher.has(t, "events.Attempting")
	dispatcher.has(t, "events.Validated")
	dispatcher.has(t, "events.Login")
	dispatcher.hasNot(t, "events.Failed")
}

func TestSessionGuardAttemptReturnsFalseWithInvalidPassword(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	ok := guard.Attempt(context.Background(), map[string]string{"email": "a@b.com", "password": "wrong"}, false)

	if ok {
		t.Error("Attempt should return false with invalid password")
	}

	dispatcher.has(t, "events.Attempting")
	dispatcher.has(t, "events.Failed")
	dispatcher.hasNot(t, "events.Validated")
	dispatcher.hasNot(t, "events.Login")
}

func TestSessionGuardAttemptReturnsFalseIfUserNotFound(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	ok := guard.Attempt(context.Background(), map[string]string{"email": "unknown@b.com", "password": "pw"}, false)

	if ok {
		t.Error("Attempt should return false when user not found")
	}

	dispatcher.has(t, "events.Attempting")
	dispatcher.has(t, "events.Failed")
}

func TestSessionGuardAttemptWithRememberSetsRememberCookie(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	cookies := &stubCookieManager{}
	guard := auth.NewSessionGuard("web", provider, sess, cookies, nil)

	ok := guard.Attempt(context.Background(), map[string]string{"email": "a@b.com", "password": "pw"}, true)

	if !ok {
		t.Fatal("Attempt with remember should succeed")
	}

	if len(cookies.queued) == 0 {
		t.Error("expected remember cookie to be queued")
	}
}

func TestSessionGuardFailedEventContainsUserWhenPasswordInvalid(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	guard.Attempt(context.Background(), map[string]string{"email": "a@b.com", "password": "wrong"}, false)

	dispatcher.mu.Lock()

	defer dispatcher.mu.Unlock()

	for _, e := range dispatcher.events {
		if fe, ok := e.(events.Failed); ok {
			if fe.User == nil {
				t.Error("Failed event should contain the user when password is invalid")
			}

			if fe.Guard != "web" {
				t.Errorf("Failed.Guard = %q, want %q", fe.Guard, "web")
			}

			return
		}
	}

	t.Error("Failed event not dispatched")
}

func TestSessionGuardFailedEventHasNilUserWhenNotFound(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	guard.Attempt(context.Background(), map[string]string{"email": "missing@b.com", "password": "pw"}, false)

	dispatcher.mu.Lock()

	defer dispatcher.mu.Unlock()

	for _, e := range dispatcher.events {
		if fe, ok := e.(events.Failed); ok {
			if fe.User != nil {
				t.Error("Failed event should have nil user when user not found")
			}

			return
		}
	}

	t.Error("Failed event not dispatched")
}

// --- SessionGuard: Login ---

func TestSessionGuardLoginStoresIdentifierInSession(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "42", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"42": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	_ = guard.Login(context.Background(), user, false)

	if sess.Get("_auth_user", nil) != "42" {
		t.Errorf("session _auth_user = %v, want \"42\"", sess.Get("_auth_user", nil))
	}
}

func TestSessionGuardLoginFiresLoginEvent(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	_ = guard.Login(context.Background(), user, false)

	dispatcher.has(t, "events.Login")
}

func TestSessionGuardLoginFiresLoginEventWithRemember(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	cookies := &stubCookieManager{}
	guard := auth.NewSessionGuard("web", provider, sess, cookies, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	_ = guard.Login(context.Background(), user, true)

	dispatcher.mu.Lock()

	defer dispatcher.mu.Unlock()

	for _, e := range dispatcher.events {
		if le, ok := e.(events.Login); ok {
			if !le.Remember {
				t.Error("Login.Remember = false, want true")
			}

			return
		}
	}

	t.Error("Login event not found")
}

func TestSessionGuardLoginQueuesCookieWhenRemembering(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	cookies := &stubCookieManager{}
	guard := auth.NewSessionGuard("web", provider, sess, cookies, nil)

	_ = guard.Login(context.Background(), user, true)

	if len(cookies.queued) == 0 {
		t.Error("expected remember cookie to be queued")
	}

	if cookies.queued[0].Name != "web_remember" {
		t.Errorf("cookie name = %q, want %q", cookies.queued[0].Name, "web_remember")
	}
}

func TestSessionGuardLoginCreatesRememberTokenIfOneDoesNotExist(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	cookies := &stubCookieManager{}
	guard := auth.NewSessionGuard("web", provider, sess, cookies, nil)

	if user.GetRememberToken() != "" {
		t.Fatal("expected no remember token initially")
	}

	_ = guard.Login(context.Background(), user, true)

	if user.GetRememberToken() == "" {
		t.Error("expected remember token to be set after login with remember")
	}
}

// --- SessionGuard: LoginUsingID ---

func TestSessionGuardLoginUsingIDLogsInWithUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "10", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"10": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	got, err := guard.LoginUsingID(context.Background(), "10", false)

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Error("expected returned user to match")
	}

	if !guard.Check(context.Background()) {
		t.Error("user should be authenticated after LoginUsingID")
	}
}

func TestSessionGuardLoginUsingIDFailure(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	_, err := guard.LoginUsingID(context.Background(), "11", false)

	if err == nil {
		t.Error("expected error for missing user")
	}
}

// --- SessionGuard: OnceUsingID ---

func TestSessionGuardOnceUsingIDSetsUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "10", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"10": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	got, err := guard.OnceUsingID(context.Background(), "10")

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Error("expected returned user to match")
	}

	if guard.Check(context.Background()) != true {
		t.Error("user should be authenticated for this request")
	}

	// Should NOT have stored in session.
	if sess.Get("_auth_user", nil) != nil {
		t.Error("OnceUsingID should not persist to session")
	}
}

func TestSessionGuardOnceUsingIDFailure(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	_, err := guard.OnceUsingID(context.Background(), "11")

	if err == nil {
		t.Error("expected error for missing user")
	}
}

// --- SessionGuard: Once ---

func TestSessionGuardOnceSetsUserWithoutSession(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	ok := guard.Once(context.Background(), map[string]string{"email": "a@b.com", "password": "pw"})

	if !ok {
		t.Error("Once should return true for valid credentials")
	}

	if !guard.Check(context.Background()) {
		t.Error("should be authenticated after Once")
	}

	// Should NOT persist in session.
	if sess.Get("_auth_user", nil) != nil {
		t.Error("Once should not persist to session")
	}
}

func TestSessionGuardOnceFailure(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	ok := guard.Once(context.Background(), map[string]string{"email": "a@b.com", "password": "wrong"})

	if ok {
		t.Error("Once should return false for invalid credentials")
	}
}

// --- SessionGuard: Validate ---

func TestSessionGuardValidateReturnsTrue(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	ok := guard.Validate(context.Background(), map[string]string{"email": "a@b.com", "password": "pw"})

	if !ok {
		t.Error("Validate should return true for valid credentials")
	}

	// Should NOT log in.
	if guard.Check(context.Background()) {
		t.Error("Validate should not log in the user")
	}
}

func TestSessionGuardValidateReturnsFalse(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	ok := guard.Validate(context.Background(), map[string]string{"email": "a@b.com", "password": "wrong"})

	if ok {
		t.Error("Validate should return false for invalid credentials")
	}
}

// --- SessionGuard: Logout ---

func TestSessionGuardLogoutRemovesSessionAndCookie(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	cookies := &stubCookieManager{}
	guard := auth.NewSessionGuard("web", provider, sess, cookies, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, true)

	if sess.Get("_auth_user", nil) == nil {
		t.Fatal("session should have user after login")
	}

	_ = guard.Logout(ctx)

	if sess.Get("_auth_user", nil) != nil {
		t.Error("session should be cleared after logout")
	}

	if guard.Check(ctx) {
		t.Error("should not be authenticated after logout")
	}

	// Should have a forget cookie queued.
	if len(cookies.forgotten) == 0 {
		t.Error("expected remember cookie to be forgotten")
	}
}

func TestSessionGuardLogoutFiresLogoutEvent(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)
	dispatcher.events = nil

	_ = guard.Logout(ctx)
	dispatcher.has(t, "events.Logout")
}

func TestSessionGuardLogoutDoesNotFireEventIfNoUser(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	_ = guard.Logout(context.Background())

	dispatcher.hasNot(t, "events.Logout")
}

func TestSessionGuardLogoutRefreshesRememberToken(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	user.SetRememberToken("old-token")
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)
	_ = guard.Logout(ctx)

	if user.GetRememberToken() == "old-token" {
		t.Error("remember token should be refreshed on logout")
	}
}

// --- SessionGuard: LogoutOtherDevices ---

func TestSessionGuardLogoutOtherDevicesDispatchesEvent(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)
	dispatcher.events = nil

	err := guard.LogoutOtherDevices(ctx)

	if err != nil {
		t.Fatal(err)
	}

	dispatcher.has(t, "events.OtherDeviceLogout")
}

// --- SessionGuard: ViaRemember ---

func TestSessionGuardViaRememberReturnsFalseByDefault(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)

	if guard.ViaRemember(ctx) {
		t.Error("ViaRemember should be false when not authenticated via cookie")
	}
}

// --- SessionGuard: SetUser / HasUser / ForgetUser ---

func TestSessionGuardSetUserFiresAuthenticatedEvent(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	guard.SetUser(context.Background(), user)

	dispatcher.has(t, "events.Authenticated")
}

func TestSessionGuardHasUserReturnsTrueWhenUserIsSet(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	if guard.HasUser() {
		t.Error("HasUser should be false initially")
	}

	guard.SetUser(context.Background(), user)

	if !guard.HasUser() {
		t.Error("HasUser should be true after SetUser")
	}
}

func TestSessionGuardForgetUserClearsUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	guard.SetUser(context.Background(), user)
	guard.ForgetUser()

	if guard.HasUser() {
		t.Error("HasUser should be false after ForgetUser")
	}
}

// --- SessionGuard: LogoutCurrentDevice ---

func TestSessionGuardLogoutCurrentDeviceFiresEvent(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	cookies := &stubCookieManager{}
	guard := auth.NewSessionGuard("web", provider, sess, cookies, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)
	dispatcher.events = nil

	err := guard.LogoutCurrentDevice(ctx)

	if err != nil {
		t.Fatal(err)
	}

	dispatcher.has(t, "events.CurrentDeviceLogout")
	dispatcher.hasNot(t, "events.Logout")
}

func TestSessionGuardLogoutCurrentDeviceRemovesSessionAndCookie(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	cookies := &stubCookieManager{}
	guard := auth.NewSessionGuard("web", provider, sess, cookies, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, true)
	_ = guard.LogoutCurrentDevice(ctx)

	if sess.Get("_auth_user", nil) != nil {
		t.Error("session should be cleared")
	}

	if guard.HasUser() {
		t.Error("user should be cleared")
	}

	if len(cookies.forgotten) == 0 {
		t.Error("remember cookie should be forgotten")
	}
}

func TestSessionGuardLogoutCurrentDeviceDoesNotRefreshRememberToken(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	user.SetRememberToken("original-token")
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	_ = guard.Login(ctx, user, false)

	// Reset token back after login may have changed it.
	user.SetRememberToken("original-token")

	_ = guard.LogoutCurrentDevice(ctx)

	if user.GetRememberToken() != "original-token" {
		t.Error("LogoutCurrentDevice should NOT refresh remember token (unlike Logout)")
	}
}

// --- SessionGuard: Login fires both Login and Authenticated ---

func TestSessionGuardLoginFiresBothLoginAndAuthenticatedEvents(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	dispatcher := &recordingDispatcher{}
	guard.SetEventDispatcher(dispatcher)

	_ = guard.Login(context.Background(), user, false)

	dispatcher.has(t, "events.Login")
	dispatcher.has(t, "events.Authenticated")
}

// --- SessionGuard: No dispatcher (backward compat) ---

func TestSessionGuardWorksWithoutDispatcher(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "a@b.com", "password": "pw"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()

	// All operations should work without panic.
	ok := guard.Attempt(ctx, map[string]string{"email": "a@b.com", "password": "pw"}, false)

	if !ok {
		t.Fatal("Attempt should succeed without dispatcher")
	}

	_ = guard.Logout(ctx)
}
