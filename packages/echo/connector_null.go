package echo

// NullConnector is a no-op Connector that returns Null channel instances and
// uses a fixed fake socket ID. Use it for testing without a real broadcast
// server or when the "null" broadcaster is selected.
type NullConnector struct{}

// compile-time interface check.
var _ Connector = (*NullConnector)(nil)

// NewNullConnector creates a new NullConnector.
func NewNullConnector() *NullConnector {
	return &NullConnector{}
}

func (c *NullConnector) Connect() error                         { return nil }
func (c *NullConnector) Channel(_ string) Channel               { return NewNullChannel() }
func (c *NullConnector) PrivateChannel(_ string) PrivateChannel { return NewNullPrivateChannel() }
func (c *NullConnector) EncryptedPrivateChannel(_ string) EncryptedPrivateChannel {
	return NewNullEncryptedPrivateChannel()
}
func (c *NullConnector) PresenceChannel(_ string) PresenceChannel { return NewNullPresenceChannel() }
func (c *NullConnector) Leave(_ string)                           {}
func (c *NullConnector) LeaveChannel(_ string)                    {}
func (c *NullConnector) LeaveAllChannels()                        {}
func (c *NullConnector) SocketID() string                         { return "fake-socket-id" }
func (c *NullConnector) Disconnect()                              {}
