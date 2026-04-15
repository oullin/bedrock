package echo

// Connector is the interface every transport backend must implement.
type Connector interface {
	// Connect establishes the underlying transport connection.
	Connect() error

	// Channel returns (or creates) a public channel subscription.
	Channel(name string) Channel

	// PrivateChannel returns (or creates) a private channel subscription.
	PrivateChannel(name string) PrivateChannel

	// EncryptedPrivateChannel returns an encrypted private channel subscription.
	EncryptedPrivateChannel(name string) EncryptedPrivateChannel

	// PresenceChannel returns a presence channel subscription.
	PresenceChannel(name string) PresenceChannel

	// Leave leaves the named channel as well as its private and presence variants.
	Leave(name string)

	// LeaveChannel leaves only the specific named channel.
	LeaveChannel(name string)

	// LeaveAllChannels leaves all currently subscribed channels.
	LeaveAllChannels()

	// SocketID returns the unique socket identifier for this connection.
	// Returns an empty string if not yet connected.
	SocketID() string

	// Disconnect closes the transport connection.
	Disconnect()
}
