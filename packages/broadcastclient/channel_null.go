package broadcastclient

// NullChannel is a no-op implementation of Channel. All methods are safe to
// call and produce no side effects. Use it as a stub or in tests that do not
// require real broadcasting.
type NullChannel struct{}

// compile-time interface check.
var _ Channel = (*NullChannel)(nil)

// NewNullChannel creates a new NullChannel.
func NewNullChannel() *NullChannel {
	return &NullChannel{}
}

func (c *NullChannel) Listen(_ string, _ Callback) Channel        { return c }
func (c *NullChannel) StopListening(_ string, _ Callback) Channel { return c }
func (c *NullChannel) ListenToAll(_ Callback) Channel             { return c }
func (c *NullChannel) StopListeningToAll(_ Callback) Channel      { return c }
func (c *NullChannel) Subscribed(_ Callback) Channel              { return c }
func (c *NullChannel) Error(_ Callback) Channel                   { return c }
func (c *NullChannel) On(_ string, _ Callback) Channel            { return c }
func (c *NullChannel) Leave()                                     {}
