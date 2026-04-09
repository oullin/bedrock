package authkit

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- test doubles ---

type testUpdatesPhotos struct {
	called bool
	url    string
}

func (a *testUpdatesPhotos) Update(_ context.Context, _ HasTeams, _ io.Reader) (string, error) {
	a.called = true
	return a.url, nil
}

type testDeletesPhotos struct{ called bool }

func (a *testDeletesPhotos) Delete(_ context.Context, _ HasTeams) error {
	a.called = true
	return nil
}

type testDeletesUsers struct{ called bool }

func (a *testDeletesUsers) Delete(_ context.Context, _ HasTeams) error {
	a.called = true
	return nil
}

func multipartPhotoRequest(path string) *http.Request {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("photo", "avatar.jpg")
	_, _ = part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, path, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// --- update profile photo tests ---

func TestUpdateProfilePhotoSuccess(t *testing.T) {
	user := &testTeamUser{id: "1"}
	photos := &testUpdatesPhotos{url: "https://cdn.example.com/avatar.jpg"}
	js, _ := buildTestAuthKit(user, newTestTeamRepo(), newTestInvitationRepo())
	js.features.ProfilePhotos = true

	handler := NewUpdateProfilePhotoHandler(js, photos)
	w := httptest.NewRecorder()
	r := multipartPhotoRequest("/user/profile-photo")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !photos.called {
		t.Fatal("expected photos.Update to be called")
	}
}

func TestUpdateProfilePhotoDisabled(t *testing.T) {
	js, _ := buildTestAuthKit(&testTeamUser{id: "1"}, newTestTeamRepo(), newTestInvitationRepo())
	js.features.ProfilePhotos = false

	handler := NewUpdateProfilePhotoHandler(js, &testUpdatesPhotos{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/user/profile-photo", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateProfilePhotoMissingFile(t *testing.T) {
	user := &testTeamUser{id: "1"}
	js, _ := buildTestAuthKit(user, newTestTeamRepo(), newTestInvitationRepo())
	js.features.ProfilePhotos = true

	handler := NewUpdateProfilePhotoHandler(js, &testUpdatesPhotos{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/user/profile-photo", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// --- delete profile photo tests ---

func TestDeleteProfilePhotoSuccess(t *testing.T) {
	user := &testTeamUser{id: "1"}
	photos := &testDeletesPhotos{}
	js, _ := buildTestAuthKit(user, newTestTeamRepo(), newTestInvitationRepo())
	js.features.ProfilePhotos = true

	handler := NewDeleteProfilePhotoHandler(js, photos)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user/profile-photo", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !photos.called {
		t.Fatal("expected photos.Delete to be called")
	}
}

// --- delete account tests ---

func TestDeleteAccountSuccess(t *testing.T) {
	user := &testTeamUser{id: "1"}
	guard := &testGuard{user: user}
	deletes := &testDeletesUsers{}
	events := &testEvents{}

	js, _ := buildTestAuthKit(user, newTestTeamRepo(), newTestInvitationRepo())
	js.guard = guard
	js.events = events

	handler := NewDeleteAccountHandler(js, deletes)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !deletes.called {
		t.Fatal("expected DeletesUsers.Delete to be called")
	}

	if len(events.dispatched) != 1 || events.dispatched[0].Name != EventUserDeleted {
		t.Fatal("expected UserDeleted event")
	}
}

func TestDeleteAccountDisabled(t *testing.T) {
	js, _ := buildTestAuthKit(&testTeamUser{id: "1"}, newTestTeamRepo(), newTestInvitationRepo())
	js.features.AccountDeletion = false

	handler := NewDeleteAccountHandler(js, &testDeletesUsers{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
