package broadcastclient

// NullPrivateChannel is a no-op implementation of PrivateChannel.
// All methods produce no side effects.
type NullPrivateChannel struct{}

// compile-time interface check.
var _ PrivateChannel = (*NullPrivateChannel)(nil)

// NewNullPrivateChannel creates a new NullPrivateChannel.
func NewNullPrivateChannel() *NullPrivateChannel {
	return &NullPrivateChannel{}
}

func (c *NullPrivateChannel) Listen(_ string, _ Callback) Channel        { return c }
func (c *NullPrivateChannel) StopListening(_ string, _ Callback) Channel { return c }
func (c *NullPrivateChannel) ListenToAll(_ Callback) Channel             { return c }
func (c *NullPrivateChannel) StopListeningToAll(_ Callback) Channel      { return c }
func (c *NullPrivateChannel) Subscribed(_ Callback) Channel              { return c }
func (c *NullPrivateChannel) Error(_ Callback) Channel                   { return c }
func (c *NullPrivateChannel) On(_ string, _ Callback) Channel            { return c }
func (c *NullPrivateChannel) Leave()                                     {}
func (c *NullPrivateChannel) Whisper(_ string, _ any) PrivateChannel     { return c }
