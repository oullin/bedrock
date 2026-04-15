package reverb

import (
	"encoding/json"
	"net/http"
	"strings"

	contractsReverb "github.com/bedrock/packages/contracts/reverb"
)

// HTTPHandler handles the Pusher-compatible HTTP API for triggering events.
// Mount at "/apps/".
//
// Routes:
//
//	POST /apps/{id}/events          - trigger a single event
//	POST /apps/{id}/batch_events    - trigger multiple events
//	GET  /apps/{id}/channels        - list channels
//	GET  /apps/{id}/channels/{ch}   - channel info
type HTTPHandler struct {
	apps       *AppManager
	channels   *ChannelManager
	dispatcher contractsReverb.Dispatcher
}

// NewHTTPHandler constructs an HTTPHandler with the given managers and dispatcher.
func NewHTTPHandler(apps *AppManager, channels *ChannelManager, dispatcher contractsReverb.Dispatcher) *HTTPHandler {
	return &HTTPHandler{
		apps:       apps,
		channels:   channels,
		dispatcher: dispatcher,
	}
}

// TriggerRequest is the body for POST /apps/{id}/events.
type TriggerRequest struct {
	Name     string   `json:"name"`
	Data     string   `json:"data"`
	Channels []string `json:"channels"`
	SocketID *string  `json:"socket_id,omitempty"`
}

// BatchTriggerRequest is the body for POST /apps/{id}/batch_events.
type BatchTriggerRequest struct {
	Batch []TriggerRequest `json:"batch"`
}

// ServeHTTP routes incoming requests to the appropriate handler.
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// path: /apps/{id}/events  or  /apps/{id}/channels  etc.
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/apps/"), "/")
	if len(parts) < 2 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	appID := parts[0]
	resource := parts[1]

	app, err := h.apps.FindByID(appID)
	if err != nil {
		http.Error(w, "app not found", http.StatusNotFound)
		return
	}

	if !h.authenticate(w, r, app) {
		return
	}

	ctx := r.Context()

	switch {
	case resource == "events" && r.Method == http.MethodPost:
		h.handleTrigger(w, r.WithContext(ctx), app)

	case resource == "batch_events" && r.Method == http.MethodPost:
		h.handleBatchTrigger(w, r.WithContext(ctx), app)

	case resource == "channels" && r.Method == http.MethodGet && len(parts) == 2:
		h.handleChannels(w, app)

	case resource == "channels" && r.Method == http.MethodGet && len(parts) >= 3:
		h.handleChannel(w, app, parts[2])

	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// authenticate verifies the Pusher HMAC signature on the request.
// It returns true when authentication succeeds, and false (after writing a 401)
// when it fails.
func (h *HTTPHandler) authenticate(w http.ResponseWriter, r *http.Request, app *App) bool {
	q := r.URL.Query()
	signature := q.Get("auth_signature")

	params := make(map[string]string, len(q))
	for k, vs := range q {
		if k == "auth_signature" {
			continue
		}
		if len(vs) > 0 {
			params[k] = vs[0]
		}
	}

	if !VerifyHTTPRequest(app.Secret(), r.Method, r.URL.Path, params, signature) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

// handleTrigger processes POST /apps/{id}/events.
func (h *HTTPHandler) handleTrigger(w http.ResponseWriter, r *http.Request, app *App) {
	var req TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	for _, chanName := range req.Channels {
		ch, err := h.channels.GetOrCreate(app.ID(), chanName)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		event := contractsReverb.Event{
			Event:   req.Name,
			Data:    req.Data,
			Channel: chanName,
		}

		if req.SocketID != nil {
			_ = ch.Broadcast(ctx, event, req.SocketID)
		} else {
			_ = ch.BroadcastToAll(ctx, event)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

// handleBatchTrigger processes POST /apps/{id}/batch_events.
func (h *HTTPHandler) handleBatchTrigger(w http.ResponseWriter, r *http.Request, app *App) {
	var req BatchTriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	for _, item := range req.Batch {
		for _, chanName := range item.Channels {
			ch, err := h.channels.GetOrCreate(app.ID(), chanName)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			event := contractsReverb.Event{
				Event:   item.Name,
				Data:    item.Data,
				Channel: chanName,
			}

			if item.SocketID != nil {
				_ = ch.Broadcast(ctx, event, item.SocketID)
			} else {
				_ = ch.BroadcastToAll(ctx, event)
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

// channelInfo holds the JSON shape for a single channel in list/info responses.
type channelInfo struct {
	SubscriptionCount int `json:"subscription_count"`
}

// channelsResponse is the JSON shape for GET /apps/{id}/channels.
type channelsResponse struct {
	Channels map[string]channelInfo `json:"channels"`
}

// handleChannels lists all channels for the application.
func (h *HTTPHandler) handleChannels(w http.ResponseWriter, app *App) {
	all := h.channels.All(app.ID())

	result := make(map[string]channelInfo, len(all))
	for _, ch := range all {
		result[ch.Name()] = channelInfo{
			SubscriptionCount: len(ch.Connections()),
		}
	}

	writeJSON(w, http.StatusOK, channelsResponse{Channels: result})
}

// handleChannel returns info for a single channel.
func (h *HTTPHandler) handleChannel(w http.ResponseWriter, app *App, name string) {
	ch, ok := h.channels.Get(app.ID(), name)
	if !ok {
		http.Error(w, "channel not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, channelInfo{
		SubscriptionCount: len(ch.Connections()),
	})
}

// writeJSON encodes v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// ensure HTTPHandler implements http.Handler at compile time.
var _ http.Handler = (*HTTPHandler)(nil)
