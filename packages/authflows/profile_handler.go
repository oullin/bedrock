package authflows

import "net/http"

// UpdateProfileHandler handles PUT /user/profile-information requests.
type UpdateProfileHandler struct {
	authflows *AuthFlows
}

// NewUpdateProfileHandler creates a new update profile handler.
func NewUpdateProfileHandler(f *AuthFlows) *UpdateProfileHandler {
	return &UpdateProfileHandler{authflows: f}
}

// ServeHTTP handles the profile information update request.
func (h *UpdateProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authflows.config.Features.UpdateProfileInformation {
		http.Error(w, "profile updates are disabled", http.StatusNotFound)
		return
	}

	ctx := r.Context()

	user, err := h.authflows.guard.AuthenticateRequest(ctx, w, r)
	if err != nil || user == nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	input := RequestInput(r, "name", h.authflows.config.IdentifierField)

	previousEmail := ""
	if verifiable, ok := user.(MustVerifyEmail); ok {
		previousEmail = verifiable.GetEmailForVerification()
	}

	if err := h.authflows.updateProfile.Update(ctx, user, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if h.authflows.config.Features.EmailVerification && h.authflows.verifier != nil {
		if verifiable, ok := user.(MustVerifyEmail); ok {
			newEmail := verifiable.GetEmailForVerification()
			if previousEmail != newEmail {
				verifiable.MarkEmailAsUnverified()
				_ = h.authflows.verifier.SendVerificationNotification(ctx, user)
			}
		}
	}

	if h.authflows.events != nil {
		_ = h.authflows.events.Dispatch(ctx, Event{Name: "authflows.profile.updated"})
	}

	h.authflows.responder.ProfileInformationUpdatedResponse(w, r)
}
