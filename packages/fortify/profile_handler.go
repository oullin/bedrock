package fortify

import (
	"net/http"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// UpdateProfileHandler handles PUT /user/profile-information requests.
type UpdateProfileHandler struct {
	fortify *Fortify
}

// NewUpdateProfileHandler creates a new update profile handler.
func NewUpdateProfileHandler(f *Fortify) *UpdateProfileHandler {
	return &UpdateProfileHandler{fortify: f}
}

// ServeHTTP handles the profile information update request.
func (h *UpdateProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.fortify.config.Features.UpdateProfileInformation {
		http.Error(w, "profile updates are disabled", http.StatusNotFound)

		return
	}

	ctx := r.Context()

	user, err := h.fortify.guard.AuthenticateRequest(ctx, w, r)

	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)

		return
	}

	input := RequestInput(r, "name", h.fortify.config.IdentifierField)

	previousEmail := ""

	if verifiable, ok := user.(cauth.MustVerifyEmail); ok {
		previousEmail = verifiable.GetEmailForVerification()
	}

	if err := h.fortify.updateProfile.Update(ctx, user, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.fortify.config.Features.EmailVerification && h.fortify.verifier != nil {
		if verifiable, ok := user.(cauth.MustVerifyEmail); ok {
			newEmail := verifiable.GetEmailForVerification()

			if previousEmail != newEmail {
				verifiable.MarkEmailAsUnverified()
				_ = h.fortify.verifier.SendVerificationNotification(ctx, user)
			}
		}
	}

	if h.fortify.events != nil {
		_ = h.fortify.events.Dispatch(ctx, "fortify.profile.updated")
	}

	h.fortify.responder.ProfileInformationUpdatedResponse(w, r)
}
