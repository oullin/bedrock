package inception

import (
	"encoding/json"
	"net/http"
)

// UpdateProfilePhotoHandler handles PUT /user/profile-photo requests.
type UpdateProfilePhotoHandler struct {
	app    *Inception
	photos UpdatesProfilePhotos
}

// NewUpdateProfilePhotoHandler creates a new update profile photo handler.
func NewUpdateProfilePhotoHandler(app *Inception, photos UpdatesProfilePhotos) *UpdateProfilePhotoHandler {
	return &UpdateProfilePhotoHandler{app: app, photos: photos}
}

// ServeHTTP handles the profile photo upload request.
func (h *UpdateProfilePhotoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.ProfilePhotos {
		http.Error(w, "profile photos are disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.app, w, r)

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

// DeleteProfilePhotoHandler handles DELETE /user/profile-photo requests.
type DeleteProfilePhotoHandler struct {
	app    *Inception
	photos DeletesProfilePhotos
}

// NewDeleteProfilePhotoHandler creates a new delete profile photo handler.
func NewDeleteProfilePhotoHandler(app *Inception, photos DeletesProfilePhotos) *DeleteProfilePhotoHandler {
	return &DeleteProfilePhotoHandler{app: app, photos: photos}
}

// ServeHTTP handles the profile photo delete request.
func (h *DeleteProfilePhotoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.ProfilePhotos {
		http.Error(w, "profile photos are disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.app, w, r)

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

// DeleteAccountHandler handles DELETE /user requests.
type DeleteAccountHandler struct {
	app     *Inception
	deletes DeletesUsers
}

// NewDeleteAccountHandler creates a new delete account handler.
func NewDeleteAccountHandler(app *Inception, deletes DeletesUsers) *DeleteAccountHandler {
	return &DeleteAccountHandler{app: app, deletes: deletes}
}

// ServeHTTP handles the delete account request.
func (h *DeleteAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.AccountDeletion {
		http.Error(w, "account deletion is disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	if err := h.deletes.Delete(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	_ = h.app.guard.Logout(r.Context(), w, r)

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventUserDeleted})
	}

	w.WriteHeader(http.StatusOK)
}
