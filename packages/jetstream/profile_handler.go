package jetstream

import (
	"encoding/json"
	"net/http"
)

// UpdateProfilePhotoHandler handles PUT /user/profile-photo requests.
type UpdateProfilePhotoHandler struct {
	js     *Jetstream
	photos UpdatesProfilePhotos
}

// NewUpdateProfilePhotoHandler creates a new update profile photo handler.

// ServeHTTP handles the profile photo upload request.

// DeleteProfilePhotoHandler handles DELETE /user/profile-photo requests.
type DeleteProfilePhotoHandler struct {
	js     *Jetstream
	photos DeletesProfilePhotos
}

// NewDeleteProfilePhotoHandler creates a new delete profile photo handler.

// ServeHTTP handles the profile photo delete request.

// DeleteAccountHandler handles DELETE /user requests.
type DeleteAccountHandler struct {
	js      *Jetstream
	deletes DeletesUsers
}

func NewUpdateProfilePhotoHandler(js *Jetstream, photos UpdatesProfilePhotos) *UpdateProfilePhotoHandler {
	return &UpdateProfilePhotoHandler{js: js, photos: photos}
}

func (h *UpdateProfilePhotoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.js.features.ProfilePhotos {
		http.Error(w, "profile photos are disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.js, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	file, _, err := r.FormFile("photo")

	if err != nil {
		http.Error(w, "photo is required", http.StatusUnprocessableEntity)

		return
	}

	defer file.Close()

	url, err := h.photos.Update(r.Context(), user, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func NewDeleteProfilePhotoHandler(js *Jetstream, photos DeletesProfilePhotos) *DeleteProfilePhotoHandler {
	return &DeleteProfilePhotoHandler{js: js, photos: photos}
}

func (h *DeleteProfilePhotoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.js.features.ProfilePhotos {
		http.Error(w, "profile photos are disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.js, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	if err := h.photos.Delete(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// NewDeleteAccountHandler creates a new delete account handler.
func NewDeleteAccountHandler(js *Jetstream, deletes DeletesUsers) *DeleteAccountHandler {
	return &DeleteAccountHandler{js: js, deletes: deletes}
}

// ServeHTTP handles the delete account request.
func (h *DeleteAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.js.features.AccountDeletion {
		http.Error(w, "account deletion is disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.js, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	if err := h.deletes.Delete(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	_ = h.js.guard.Logout(r.Context(), w, r)

	if h.js.events != nil {
		_ = h.js.events.Dispatch(r.Context(), Event{Name: EventUserDeleted})
	}

	w.WriteHeader(http.StatusOK)
}
