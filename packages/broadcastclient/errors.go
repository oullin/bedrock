package broadcastclient

import "errors"

var (
	// ErrUnsupportedBroadcaster is returned by New when the broadcaster
	// string does not match any built-in or custom connector.
	ErrUnsupportedBroadcaster = errors.New("broadcastclient: unsupported broadcaster")
)
