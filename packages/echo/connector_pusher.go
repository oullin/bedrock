package echo

// PusherConnector is a stub Connector for the Pusher and Reverb transports.
// The Connect method is a no-op in this implementation; a full implementation
// would dial the Pusher/Reverb WebSocket endpoint using pusher-go or similar.
type PusherConnector struct {
	opts Options
}

// compile-time interface check.
var _ Connector = (*PusherConnector)(nil)

// NewPusherConnector creates a PusherConnector with the given options.
func NewPusherConnector(opts Options) *PusherConnector {
	return &PusherConnector{opts: opts}
}

func (c *PusherConnector) Connect() error                         { return nil }
func (c *PusherConnector) Channel(_ string) Channel               { return NewNullChannel() }
func (c *PusherConnector) PrivateChannel(_ string) PrivateChannel { return NewNullPrivateChannel() }
func (c *PusherConnector) EncryptedPrivateChannel(_ string) EncryptedPrivateChannel {
	return NewNullEncryptedPrivateChannel()
}
func (c *PusherConnector) PresenceChannel(_ string) PresenceChannel { return NewNullPresenceChannel() }
func (c *PusherConnector) Leave(_ string)                           {}
func (c *PusherConnector) LeaveChannel(_ string)                    {}
func (c *PusherConnector) LeaveAllChannels()                        {}
func (c *PusherConnector) SocketID() string                         { return "" }
func (c *PusherConnector) Disconnect()                              {}
