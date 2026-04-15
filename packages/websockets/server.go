package websockets

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
	"nhooyr.io/websocket"
)

// Server handles WebSocket connections using the Pusher protocol.
// It implements http.Handler and should be mounted at "/app/".
//
// URL pattern: GET /app/{key}
type Server struct {
	apps       *AppManager
	conns      *ConnectionManager
	channels   *ChannelManager
	dispatcher contractsWebSockets.Dispatcher
	config     Config
}

// NewServer constructs a Server with all required dependencies.
func NewServer(cfg Config, apps *AppManager, conns *ConnectionManager, channels *ChannelManager, dispatcher contractsWebSockets.Dispatcher) *Server {
	return &Server{
		apps:       apps,
		conns:      conns,
		channels:   channels,
		dispatcher: dispatcher,
		config:     cfg,
	}
}

// ServeHTTP upgrades HTTP requests to WebSocket connections.
// It expects the URL pattern /app/{key}.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/app/")

	app, err := s.apps.FindByKey(key)
	if err != nil {
		http.Error(w, "app not found", http.StatusNotFound)
		return
	}

	if !ValidateOrigin(r.Header.Get("Origin"), app.AllowedOrigins()) {
		http.Error(w, "invalid origin", http.StatusForbidden)
		return
	}

	if app.MaxConnections() > 0 && s.conns.Count(app.ID()) >= app.MaxConnections() {
		http.Error(w, "connection limit reached", http.StatusServiceUnavailable)
		return
	}

	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // origin already validated above
	})
	if err != nil {
		return
	}

	if app.MaxMessageSize() > 0 {
		ws.SetReadLimit(app.MaxMessageSize())
	}

	conn := NewConn(ws, app.ID())
	s.conns.Add(conn)
	defer s.conns.Remove(app.ID(), conn.SocketID())

	ctx := r.Context()

	if err := s.sendConnectionEstablished(ctx, conn, app); err != nil {
		return
	}

	s.readLoop(ctx, conn, app)
}

// sendConnectionEstablished sends the pusher:connection_established event to a new connection.
func (s *Server) sendConnectionEstablished(ctx context.Context, conn *Conn, app *App) error {
	data := ConnectionEstablishedData{
		SocketID:        conn.SocketID(),
		ActivityTimeout: app.ActivityTimeout(),
	}

	bytes, err := MarshalEvent("pusher:connection_established", "", data)
	if err != nil {
		return err
	}

	return conn.Send(ctx, bytes)
}

// readLoop reads messages from the connection until the connection is closed.
func (s *Server) readLoop(ctx context.Context, conn *Conn, app *App) {
	for {
		_, raw, err := conn.ws.Read(ctx)
		if err != nil {
			break
		}

		conn.Touch()
		conn.TouchMessage()

		if !Allow(conn, 100, time.Second) {
			s.disconnect(ctx, conn, CodeRateLimitExceeded, "rate limit exceeded")
			break
		}

		msg, err := Parse(raw)
		if err != nil {
			s.sendError(ctx, conn, CodeInvalidMessage, "invalid message")
			continue
		}

		if err := s.handleMessage(ctx, conn, app, msg); err != nil {
			break
		}
	}

	// Cleanup: unsubscribe from all channels.
	for _, ch := range s.channels.All(app.ID()) {
		if ch.HasConnection(conn.SocketID()) {
			ch.Unsubscribe(ctx, conn)
		}
	}

	s.channels.CleanupEmpty(app.ID())
}

// handleMessage dispatches an inbound message to the appropriate handler.
func (s *Server) handleMessage(ctx context.Context, conn *Conn, app *App, msg PusherMessage) error {
	switch msg.Event {
	case "pusher:subscribe":
		return s.handleSubscribe(ctx, conn, app, msg)
	case "pusher:unsubscribe":
		return s.handleUnsubscribe(ctx, conn, app, msg)
	case "pusher:ping":
		return s.handlePing(ctx, conn)
	case "pusher:pong":
		return s.handlePong(conn)
	default:
		if strings.HasPrefix(msg.Event, "client-") {
			return s.handleClientEvent(ctx, conn, app, msg)
		}
	}

	return nil
}

// handleSubscribe processes a pusher:subscribe message.
func (s *Server) handleSubscribe(ctx context.Context, conn *Conn, app *App, msg PusherMessage) error {
	sd, err := ParseSubscribeData(msg.Data)
	if err != nil {
		s.sendError(ctx, conn, CodeInvalidMessage, "invalid subscribe data")
		return nil
	}

	ch, err := s.channels.GetOrCreate(app.ID(), sd.Channel)
	if err != nil {
		return err
	}

	if err := ch.Subscribe(ctx, conn, sd.Auth, sd.ChannelData); err != nil {
		if err == ErrUnauthorized {
			s.sendError(ctx, conn, CodeUnauthorized, "unauthorized")
			return nil
		}

		return err
	}

	return nil
}

// handleUnsubscribe processes a pusher:unsubscribe message.
func (s *Server) handleUnsubscribe(ctx context.Context, conn *Conn, app *App, msg PusherMessage) error {
	var data struct {
		Channel string `json:"channel"`
	}

	if len(msg.Data) > 0 {
		// data field may be a JSON-encoded string (double-encoded) or direct object.
		var str string
		if err := json.Unmarshal(msg.Data, &str); err == nil {
			_ = json.Unmarshal([]byte(str), &data)
		} else {
			_ = json.Unmarshal(msg.Data, &data)
		}
	}

	if data.Channel == "" {
		return nil
	}

	ch, ok := s.channels.Get(app.ID(), data.Channel)
	if !ok {
		return nil
	}

	ch.Unsubscribe(ctx, conn)
	s.channels.CleanupEmpty(app.ID())

	return nil
}

// handlePing responds to a pusher:ping with a pusher:pong.
func (s *Server) handlePing(ctx context.Context, conn *Conn) error {
	pong := []byte(`{"event":"pusher:pong","data":"{}"}`)
	return conn.Send(ctx, pong)
}

// handlePong records that a pong was received from the client.
func (s *Server) handlePong(conn *Conn) error {
	conn.TouchPong()
	return nil
}

// handleClientEvent processes a client- prefixed event (whisper).
func (s *Server) handleClientEvent(ctx context.Context, conn *Conn, app *App, msg PusherMessage) error {
	if msg.Channel == "" {
		s.sendError(ctx, conn, CodeInvalidMessage, "client event missing channel")
		return nil
	}

	ch, ok := s.channels.Get(app.ID(), msg.Channel)
	if !ok || !ch.HasConnection(conn.SocketID()) {
		s.sendError(ctx, conn, CodeUnauthorized, "not subscribed to channel")
		return nil
	}

	switch app.ClientEventsMode() {
	case "none":
		s.sendError(ctx, conn, CodeClientEventsDisabled, "client events disabled")
		return nil

	case "members":
		pc, ok := ch.(contractsWebSockets.PresenceChanneler)
		if !ok {
			s.sendError(ctx, conn, CodeUnauthorized, "client events require presence channel")
			return nil
		}

		socketID := conn.SocketID()
		isMember := false
		for _, c := range pc.Connections() {
			if c.SocketID() == socketID {
				isMember = true
				break
			}
		}

		if !isMember {
			s.sendError(ctx, conn, CodeUnauthorized, "not a channel member")
			return nil
		}

	case "all":
		// Any subscriber may send client events.
	}

	socketID := conn.SocketID()
	event := contractsWebSockets.Event{
		Event:   msg.Event,
		Data:    string(msg.Data),
		Channel: msg.Channel,
	}

	return ch.Broadcast(ctx, event, &socketID)
}

// disconnect sends an error frame and closes the connection.
func (s *Server) disconnect(ctx context.Context, conn *Conn, code int, reason string) {
	errBytes, _ := MarshalError(code, reason)
	_ = conn.Send(ctx, errBytes)
	_ = conn.Close(ctx, code, reason)
}

// sendError sends an error event without closing the connection.
func (s *Server) sendError(ctx context.Context, conn *Conn, code int, reason string) {
	errBytes, _ := MarshalError(code, reason)
	_ = conn.Send(ctx, errBytes)
}

// ensure Server implements http.Handler at compile time.
var _ http.Handler = (*Server)(nil)
