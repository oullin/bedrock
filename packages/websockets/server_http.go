package websockets

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
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
	conns      *ConnectionManager
	channels   *ChannelManager
	dispatcher contractsWebSockets.Dispatcher
}

// NewHTTPHandler constructs an HTTPHandler with the given managers and dispatcher.

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

// path: /apps/{id}/events  or  /apps/{id}/channels  etc.

// authenticate verifies the Pusher HMAC signature on the request.
// It returns true when authentication succeeds, and false (after writing a 401)
// when it fails.

// handleTrigger processes POST /apps/{id}/events.

// handleBatchTrigger processes POST /apps/{id}/batch_events.

// channelsResponse is the JSON shape for GET /apps/{id}/channels.
type channelsResponse struct {
	Channels map[string]map[string]any `json:"channels"`
}

func NewHTTPHandler(apps *AppManager, conns *ConnectionManager, channels *ChannelManager, dispatcher contractsWebSockets.Dispatcher) *HTTPHandler {
	return &HTTPHandler{
		apps:       apps,
		conns:      conns,
		channels:   channels,
		dispatcher: dispatcher,
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/up" {
		h.handleHealth(w)

		return
	}

	appPath, ok := appPathFromRequest(r.URL.Path)

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)

		return
	}

	parts := strings.Split(strings.TrimPrefix(appPath, "/apps/"), "/")

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

	case resource == "connections" && r.Method == http.MethodGet && len(parts) == 2:
		h.handleConnections(w, app)

	case resource == "channels" && r.Method == http.MethodGet && len(parts) == 2:
		h.handleChannels(w, r, app)

	case resource == "channels" && r.Method == http.MethodGet && len(parts) >= 4 && parts[3] == "users":
		h.handleChannelUsers(w, app, parts[2])

	case resource == "channels" && r.Method == http.MethodGet && len(parts) >= 3:
		h.handleChannel(w, r, app, parts[2])

	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func appPathFromRequest(path string) (string, bool) {
	idx := strings.Index(path, "/apps/")

	if idx < 0 {
		return "", false
	}

	return path[idx:], true
}

func (h *HTTPHandler) authenticate(w http.ResponseWriter, r *http.Request, app *App) bool {
	q := r.URL.Query()
	signature := q.Get("auth_signature")
	normalizedPath, ok := appPathFromRequest(r.URL.Path)

	if !ok {
		normalizedPath = r.URL.Path
	}

	params := make(map[string]string, len(q))

	for k, vs := range q {
		if k == "auth_signature" {
			continue
		}

		if len(vs) > 0 {
			params[k] = vs[0]
		}
	}

	if !VerifyHTTPRequest(app.Secret(), r.Method, normalizedPath, params, signature) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)

		return false
	}

	return true
}

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

		event := contractsWebSockets.Event{
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

	info := requestedInfo(r)

	if len(info) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{})

		return
	}

	writeJSON(w, http.StatusOK, h.eventResponse(app, req.Channels, info))
}

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

			event := contractsWebSockets.Event{
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

	info := requestedInfo(r)

	if len(info) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{})

		return
	}

	writeJSON(w, http.StatusOK, h.eventResponse(app, flattenBatchChannels(req.Batch), info))
}

// handleChannels lists all channels for the application.
func (h *HTTPHandler) handleChannels(w http.ResponseWriter, r *http.Request, app *App) {
	all := h.channels.All(app.ID())
	info := requestedInfo(r)

	result := make(map[string]map[string]any, len(all))

	for _, ch := range all {
		result[ch.Name()] = projectChannelListInfo(ch, info)
	}

	writeJSON(w, http.StatusOK, channelsResponse{Channels: result})
}

// handleChannel returns info for a single channel.
func (h *HTTPHandler) handleChannel(w http.ResponseWriter, r *http.Request, app *App, name string) {
	ch, ok := h.channels.Get(app.ID(), name)

	if !ok {
		http.Error(w, "channel not found", http.StatusNotFound)

		return
	}

	if len(r.URL.Query().Get("info")) == 0 {
		if len(ch.Connections()) == 0 {
			writeJSON(w, http.StatusOK, map[string]any{"occupied": false})

			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"occupied": true})

		return
	}

	writeJSON(w, http.StatusOK, projectChannelInfo(ch, requestedInfo(r)))
}

func requestedInfo(r *http.Request) map[string]struct{} {
	raw := r.URL.Query().Get("info")

	if raw == "" {
		return nil
	}

	info := make(map[string]struct{})

	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)

		if part != "" {
			info[part] = struct{}{}
		}
	}

	return info
}

func projectChannelListInfo(ch contractsWebSockets.Channel, info map[string]struct{}) map[string]any {
	if len(info) == 0 {
		return map[string]any{}
	}

	attrs := make(map[string]any)

	if hasInfo(info, "subscription_count") {
		attrs["subscription_count"] = len(ch.Connections())
	}

	if hasInfo(info, "user_count") {
		if pc, ok := ch.(contractsWebSockets.PresenceChanneler); ok {
			attrs["user_count"] = pc.MemberCount()
		}
	}

	if hasInfo(info, "cache") {
		if cache, ok := ch.(contractsWebSockets.CacheableChannel); ok {
			if cached := cachedEventData(cache.LastEvent()); cached != nil {
				attrs["cache"] = cached
			}
		}
	}

	return attrs
}

func projectChannelInfo(ch contractsWebSockets.Channel, info map[string]struct{}) map[string]any {
	if len(ch.Connections()) == 0 {
		return map[string]any{"occupied": false}
	}

	attrs := map[string]any{"occupied": true}

	if hasInfo(info, "subscription_count") {
		attrs["subscription_count"] = len(ch.Connections())
	}

	if hasInfo(info, "user_count") {
		if pc, ok := ch.(contractsWebSockets.PresenceChanneler); ok {
			attrs["user_count"] = pc.MemberCount()
		}
	}

	if hasInfo(info, "cache") {
		if cache, ok := ch.(contractsWebSockets.CacheableChannel); ok {
			if cached := cachedEventData(cache.LastEvent()); cached != nil {
				attrs["cache"] = cached
			}
		}
	}

	return attrs
}

func hasInfo(info map[string]struct{}, key string) bool {
	if len(info) == 0 {
		return false
	}

	_, ok := info[key]

	return ok
}

func cachedEventData(event *contractsWebSockets.Event) any {
	if event == nil {
		return nil
	}

	var data any

	if err := json.Unmarshal([]byte(event.Data), &data); err == nil {
		return data
	}

	return event.Data
}

func parseUserID(id string) any {
	if n, err := strconv.Atoi(id); err == nil {
		return n
	}

	return id
}

func (h *HTTPHandler) handleChannelUsers(w http.ResponseWriter, app *App, name string) {
	ch, ok := h.channels.Get(app.ID(), name)

	if !ok {
		http.Error(w, "presence channel not found", http.StatusBadRequest)

		return
	}

	pc, ok := ch.(contractsWebSockets.PresenceChanneler)

	if !ok {
		http.Error(w, "presence channel not found", http.StatusBadRequest)

		return
	}

	if len(pc.Connections()) == 0 {
		http.Error(w, "presence channel not occupied", http.StatusNotFound)

		return
	}

	users := make([]map[string]any, 0, len(pc.MemberIDs()))

	for _, id := range pc.MemberIDs() {
		users = append(users, map[string]any{"id": parseUserID(id)})
	}

	writeJSON(w, http.StatusOK, map[string]any{"users": users})
}

func (h *HTTPHandler) handleConnections(w http.ResponseWriter, app *App) {
	writeJSON(w, http.StatusOK, map[string]any{"connections": h.conns.Count(app.ID())})
}

func (h *HTTPHandler) handleHealth(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, map[string]string{"health": "OK"})
}

func (h *HTTPHandler) eventResponse(app *App, channels []string, info map[string]struct{}) map[string]any {
	result := make(map[string]map[string]any)
	seen := make(map[string]struct{})

	for _, chanName := range channels {
		if _, ok := seen[chanName]; ok {
			continue
		}

		seen[chanName] = struct{}{}

		ch, ok := h.channels.Get(app.ID(), chanName)

		if !ok {
			continue
		}

		result[chanName] = projectEventInfo(chanName, ch, info)
	}

	return map[string]any{"channels": result}
}

func projectEventInfo(chanName string, ch contractsWebSockets.Channel, info map[string]struct{}) map[string]any {
	attrs := make(map[string]any)
	channelType := TypeOf(chanName)

	if hasInfo(info, "subscription_count") && !strings.HasPrefix(channelType, "presence") {
		attrs["subscription_count"] = len(ch.Connections())
	}

	if hasInfo(info, "user_count") {
		if pc, ok := ch.(contractsWebSockets.PresenceChanneler); ok {
			attrs["user_count"] = pc.MemberCount()
		}
	}

	if hasInfo(info, "cache") {
		if cache, ok := ch.(contractsWebSockets.CacheableChannel); ok {
			if cached := cachedEventData(cache.LastEvent()); cached != nil {
				attrs["cache"] = cached
			}
		}
	}

	return attrs
}

func flattenBatchChannels(batch []TriggerRequest) []string {
	channels := make([]string, 0)

	for _, item := range batch {
		channels = append(channels, item.Channels...)
	}

	return channels
}

// writeJSON encodes v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	payload, err := json.Marshal(v)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

// ensure HTTPHandler implements http.Handler at compile time.
var _ http.Handler = (*HTTPHandler)(nil)
