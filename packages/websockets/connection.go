package websockets

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
	"nhooyr.io/websocket"
)

// Conn wraps an nhooyr.io/websocket.Conn as a WebSockets connection.
// Socket IDs are generated in the format "<rand.Int31>.<rand.Int31>".
// All methods are safe for concurrent use.
type Conn struct {
	mu             sync.Mutex
	id             string
	appID          string
	ws             *websocket.Conn
	lastSeenAt     time.Time
	lastMsgAt      time.Time
	pongedAt       time.Time
	msgCount       int64
	msgWindowStart time.Time
}

var _ contractsWebSockets.Connection = (*Conn)(nil)

// NewConn creates a new Conn wrapping the given websocket connection for the
// specified application. A unique socket ID is generated automatically.
func NewConn(ws *websocket.Conn, appID string) *Conn {
	now := time.Now()

	return &Conn{
		id:             fmt.Sprintf("%d.%d", rand.Int31(), rand.Int31()),
		appID:          appID,
		ws:             ws,
		lastSeenAt:     now,
		lastMsgAt:      now,
		msgWindowStart: now,
	}
}

// SocketID returns the unique identifier for this connection.
func (c *Conn) SocketID() string {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.id
}

// AppID returns the application ID this connection belongs to.
func (c *Conn) AppID() string {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.appID
}

// Send writes a raw text message to the client.
func (c *Conn) Send(ctx context.Context, msg []byte) error {
	return c.ws.Write(ctx, websocket.MessageText, msg)
}

// Close terminates the WebSocket connection with the given Pusher error code
// and human-readable reason. The code is mapped directly to a WebSocket status
// code via websocket.StatusCode.
func (c *Conn) Close(ctx context.Context, code int, reason string) error {
	return c.ws.Close(websocket.StatusCode(code), reason)
}

// LastSeenAt returns the time the connection last sent any message.
func (c *Conn) LastSeenAt() time.Time {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.lastSeenAt
}

// Touch updates the LastSeenAt timestamp to now.
func (c *Conn) Touch() {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.lastSeenAt = time.Now()
}

// TouchMessage updates the last-message-at timestamp to now.
func (c *Conn) TouchMessage() {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.lastMsgAt = time.Now()
}

// TouchPong updates the last-pong-at timestamp to now.
func (c *Conn) TouchPong() {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.pongedAt = time.Now()
}

// IncrMessageCount increments the per-window message counter and returns the
// new count.
func (c *Conn) IncrMessageCount() int64 {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.msgCount++

	return c.msgCount
}

// ResetMessageWindow resets the sliding-window counter and sets the window
// start time to t.
func (c *Conn) ResetMessageWindow(t time.Time) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.msgCount = 0
	c.msgWindowStart = t
}

// MessageWindowStart returns the start time of the current message rate window.
func (c *Conn) MessageWindowStart() time.Time {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.msgWindowStart
}
