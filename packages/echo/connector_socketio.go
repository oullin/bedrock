package echo

// SocketIOConnector is a stub Connector for the Socket.IO transport.
// The Connect method is a no-op in this implementation; a full implementation
// would dial the Socket.IO server using an appropriate Go client library.
type SocketIOConnector struct {
	opts Options
}

// compile-time interface check.
var _ Connector = (*SocketIOConnector)(nil)

// NewSocketIOConnector creates a SocketIOConnector with the given options.
func NewSocketIOConnector(opts Options) *SocketIOConnector {
	return &SocketIOConnector{opts: opts}
}

func (c *SocketIOConnector) Connect() error                         { return nil }
func (c *SocketIOConnector) Channel(_ string) Channel               { return NewNullChannel() }
func (c *SocketIOConnector) PrivateChannel(_ string) PrivateChannel { return NewNullPrivateChannel() }
func (c *SocketIOConnector) EncryptedPrivateChannel(_ string) EncryptedPrivateChannel {
	return NewNullEncryptedPrivateChannel()
}
func (c *SocketIOConnector) PresenceChannel(_ string) PresenceChannel {
	return NewNullPresenceChannel()
}
func (c *SocketIOConnector) Leave(_ string)        {}
func (c *SocketIOConnector) LeaveChannel(_ string) {}
func (c *SocketIOConnector) LeaveAllChannels()     {}
func (c *SocketIOConnector) SocketID() string      { return "" }
func (c *SocketIOConnector) Disconnect()           {}
