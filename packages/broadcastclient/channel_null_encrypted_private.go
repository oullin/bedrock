package broadcastclient

// NullEncryptedPrivateChannel is a no-op implementation of EncryptedPrivateChannel.
// All methods produce no side effects.
type NullEncryptedPrivateChannel struct{}

// compile-time interface check.
var _ EncryptedPrivateChannel = (*NullEncryptedPrivateChannel)(nil)

// NewNullEncryptedPrivateChannel creates a new NullEncryptedPrivateChannel.
func NewNullEncryptedPrivateChannel() *NullEncryptedPrivateChannel {
	return &NullEncryptedPrivateChannel{}
}

func (c *NullEncryptedPrivateChannel) Listen(_ string, _ Callback) Channel        { return c }
func (c *NullEncryptedPrivateChannel) StopListening(_ string, _ Callback) Channel { return c }
func (c *NullEncryptedPrivateChannel) ListenToAll(_ Callback) Channel             { return c }
func (c *NullEncryptedPrivateChannel) StopListeningToAll(_ Callback) Channel      { return c }
func (c *NullEncryptedPrivateChannel) Subscribed(_ Callback) Channel              { return c }
func (c *NullEncryptedPrivateChannel) Error(_ Callback) Channel                   { return c }
func (c *NullEncryptedPrivateChannel) On(_ string, _ Callback) Channel            { return c }
func (c *NullEncryptedPrivateChannel) Leave()                                     {}
func (c *NullEncryptedPrivateChannel) Whisper(_ string, _ any) PrivateChannel     { return c }
