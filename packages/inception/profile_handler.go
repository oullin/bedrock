package inception

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// UpdateProfileHandler handles PUT /user/profile-information requests.
type UpdateProfileHandler struct {
	app *Inception
}

// NewUpdateProfileHandler creates a new update profile handler.
func NewUpdateProfileHandler(app *Inception) *UpdateProfileHandler {
	return &UpdateProfileHandler{app: app}
}

// ServeHTTP handles the profile information update request.
func (h *UpdateProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.UpdateProfileInformation {
		http.Error(w, "profile updates are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.app.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "name", h.app.config.IdentifierField)

	previousEmail := ""

	if verifiable, ok := user.(cauth.MustVerifyEmail); ok {
		previousEmail = verifiable.GetEmailForVerification()
	}

	if err := h.app.updateProfile.Update(ctx, user, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.config.Features.EmailVerification && h.app.verifier != nil {
		if verifiable, ok := user.(cauth.MustVerifyEmail); ok {
			newEmail := verifiable.GetEmailForVerification()

			if previousEmail != newEmail {
				verifiable.MarkEmailAsUnverified()
				_ = h.app.verifier.SendVerificationNotification(ctx, user)
			}
		}
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(ctx, EventProfileUpdated)
	}

	h.app.responder.ProfileInformationUpdatedResponse(w, r)
}
