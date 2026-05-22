package echo

import "fmt"

// Echo is the primary API for interacting with real-time broadcasting.
// It wraps a Connector and delegates channel subscription management to it.
type Echo struct {
	opts      Options
	connector Connector
}

// New creates an Echo instance and establishes the underlying connection.
// An error is returned when the broadcaster name is not recognised and no
// custom Connector was provided via Options.Connector.
func New(opts Options) (*Echo, error) {
	e := &Echo{opts: opts}

	if err := e.connect(); err != nil {
		return nil, err
	}

	return e, nil
}

// connect resolves the Connector based on opts and calls Connect().
func (e *Echo) connect() error {
	if e.opts.Connector != nil {
		e.connector = e.opts.Connector

		return e.connector.Connect()
	}

	switch e.opts.Broadcaster {
	case "websockets", "pusher":
		e.connector = NewPusherConnector(e.opts)
	case "socket.io":
		e.connector = NewSocketIOConnector(e.opts)
	case "null", "":
		e.connector = NewNullConnector()
	default:
		return fmt.Errorf("broadcaster string %s is not supported", e.opts.Broadcaster)
	}

	return e.connector.Connect()
}

// Channel returns a public channel subscription for the given name.
func (e *Echo) Channel(name string) Channel {
	return e.connector.Channel(name)
}

// PrivateChannel returns a private channel subscription for the given name.
func (e *Echo) PrivateChannel(name string) PrivateChannel {
	return e.connector.PrivateChannel(name)
}

// EncryptedPrivateChannel returns an encrypted private channel subscription.
func (e *Echo) EncryptedPrivateChannel(name string) EncryptedPrivateChannel {
	return e.connector.EncryptedPrivateChannel(name)
}

// PresenceChannel returns a presence channel subscription for the given name.
func (e *Echo) PresenceChannel(name string) PresenceChannel {
	return e.connector.PresenceChannel(name)
}

// Leave leaves a channel by name, including its private and presence variants.
func (e *Echo) Leave(name string) {
	e.connector.Leave(name)
}

// LeaveChannel leaves only the specific named channel.
func (e *Echo) LeaveChannel(name string) {
	e.connector.LeaveChannel(name)
}

// LeaveAllChannels leaves all currently subscribed channels.
func (e *Echo) LeaveAllChannels() {
	e.connector.LeaveAllChannels()
}

// SocketID returns the connection's socket ID. Returns an empty string until
// the connection is established.
func (e *Echo) SocketID() string {
	return e.connector.SocketID()
}

// Disconnect closes the transport connection.
func (e *Echo) Disconnect() {
	e.connector.Disconnect()
}

// Connector returns the underlying Connector instance.
func (e *Echo) Connector() Connector {
	return e.connector
}
