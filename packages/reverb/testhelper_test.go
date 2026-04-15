package reverb_test

import (
	"context"
	"sync"
	"time"

	contractsReverb "github.com/bedrock/packages/contracts/reverb"
	. "github.com/bedrock/packages/reverb"
)

// fakeConn is a test double for contractsReverb.Connection.
type fakeConn struct {
	mu       sync.Mutex
	id       string
	appID    string
	sent     [][]byte
	closed   bool
	lastSeen time.Time
}

func newFakeConn(id, appID string) *fakeConn {
	return &fakeConn{id: id, appID: appID, lastSeen: time.Now()}
}

var _ contractsReverb.Connection = (*fakeConn)(nil)

func (f *fakeConn) SocketID() string { return f.id }
func (f *fakeConn) AppID() string    { return f.appID }
func (f *fakeConn) Send(_ context.Context, msg []byte) error {
	f.mu.Lock()
	f.sent = append(f.sent, msg)
	f.mu.Unlock()
	return nil
}
func (f *fakeConn) Close(_ context.Context, _ int, _ string) error {
	f.mu.Lock()
	f.closed = true
	f.mu.Unlock()
	return nil
}
func (f *fakeConn) LastSeenAt() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.lastSeen
}
func (f *fakeConn) Touch() {
	f.mu.Lock()
	f.lastSeen = time.Now()
	f.mu.Unlock()
}
func (f *fakeConn) TouchMessage() {}
func (f *fakeConn) TouchPong()    {}

func (f *fakeConn) SentMessages() [][]byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([][]byte(nil), f.sent...)
}

func (f *fakeConn) WasClosed() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.closed
}

// newTestApp creates an App for tests.
func newTestApp() *App {
	return NewApp(AppConfig{
		ID:             "app-1",
		Key:            "key-1",
		Secret:         "secret-1",
		AllowedOrigins: []string{"*"},
		ClientEvents:   ClientEventsConfig{Mode: "all"},
	})
}
