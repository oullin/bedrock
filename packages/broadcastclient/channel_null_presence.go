package broadcastclient

// NullPresenceChannel is a no-op implementation of PresenceChannel.
// All methods produce no side effects.
type NullPresenceChannel struct{}

// compile-time interface check.
var _ PresenceChannel = (*NullPresenceChannel)(nil)

// NewNullPresenceChannel creates a new NullPresenceChannel.
func NewNullPresenceChannel() *NullPresenceChannel {
	return &NullPresenceChannel{}
}

func (c *NullPresenceChannel) Listen(_ string, _ Callback) Channel        { return c }
func (c *NullPresenceChannel) StopListening(_ string, _ Callback) Channel { return c }
func (c *NullPresenceChannel) ListenToAll(_ Callback) Channel             { return c }
func (c *NullPresenceChannel) StopListeningToAll(_ Callback) Channel      { return c }
func (c *NullPresenceChannel) Subscribed(_ Callback) Channel              { return c }
func (c *NullPresenceChannel) Error(_ Callback) Channel                   { return c }
func (c *NullPresenceChannel) On(_ string, _ Callback) Channel            { return c }
func (c *NullPresenceChannel) Leave()                                     {}
func (c *NullPresenceChannel) Here(_ func([]any)) PresenceChannel         { return c }
func (c *NullPresenceChannel) Joining(_ func(any)) PresenceChannel        { return c }
func (c *NullPresenceChannel) Leaving(_ func(any)) PresenceChannel        { return c }
func (c *NullPresenceChannel) Whisper(_ string, _ any) PresenceChannel    { return c }
