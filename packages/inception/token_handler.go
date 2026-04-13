package inception

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// CreateTokenHandler handles POST /user/api-tokens requests.
type CreateTokenHandler struct {
	app    *Inception
	tokens TokenRepository
}

// NewCreateTokenHandler creates a new create token handler.

// ServeHTTP handles the create token request.

// UpdateTokenHandler handles PUT /user/api-tokens/{token} requests.
type UpdateTokenHandler struct {
	app    *Inception
	tokens TokenRepository
}

// NewUpdateTokenHandler creates a new update token handler.

// ServeHTTP handles the update token request.

// DeleteTokenHandler handles DELETE /user/api-tokens/{token} requests.
type DeleteTokenHandler struct {
	app    *Inception
	tokens TokenRepository
}

func NewCreateTokenHandler(app *Inception, tokens TokenRepository) *CreateTokenHandler {
	return &CreateTokenHandler{app: app, tokens: tokens}
}

func (h *CreateTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	input := RequestInput(r, "name", "permissions")
	name := input["name"]

	if name == "" {
		http.Error(w, "token name is required", http.StatusUnprocessableEntity)

		return
	}

	var permissions []string

	if raw := input["permissions"]; raw != "" {
		permissions = strings.Split(raw, ",")
	} else {
		permissions = []string{"*"}
	}

	plain, hash, err := GeneratePlainToken()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	token := &PersonalAccessToken{
		UserID:      user.GetAuthIdentifier(),
		Name:        name,
		TokenHash:   hash,
		Permissions: permissions,
		CreatedAt:   time.Now(),
	}

	if err := h.tokens.Create(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTokenCreated, Payload: token})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(NewTokenResult{
		PlainText: plain,
		Token:     *token,
	})
}

func NewUpdateTokenHandler(app *Inception, tokens TokenRepository) *UpdateTokenHandler {
	return &UpdateTokenHandler{app: app, tokens: tokens}
}

func (h *UpdateTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	tokenID := r.PathValue("token")

	token, err := h.tokens.FindByID(r.Context(), tokenID)

	if err != nil || token == nil {
		http.Error(w, "token not found", http.StatusNotFound)

		return
	}

	if token.UserID != user.GetAuthIdentifier() {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	input := RequestInput(r, "permissions")

	if raw := input["permissions"]; raw != "" {
		token.Permissions = strings.Split(raw, ",")
	}

	if err := h.tokens.Update(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTokenUpdated, Payload: token})
	}

	w.WriteHeader(http.StatusOK)
}

// NewDeleteTokenHandler creates a new delete token handler.
func NewDeleteTokenHandler(app *Inception, tokens TokenRepository) *DeleteTokenHandler {
	return &DeleteTokenHandler{app: app, tokens: tokens}
}

// ServeHTTP handles the delete token request.
func (h *DeleteTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	tokenID := r.PathValue("token")

	token, err := h.tokens.FindByID(r.Context(), tokenID)

	if err != nil || token == nil {
		http.Error(w, "token not found", http.StatusNotFound)

		return
	}

	if token.UserID != user.GetAuthIdentifier() {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	if err := h.tokens.Delete(r.Context(), tokenID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTokenDeleted, Payload: token})
	}

	w.WriteHeader(http.StatusOK)
}
