package events

// ConnectionEstablished is dispatched when a database connection is established.
type ConnectionEstablished struct {
	ConnectionName string
	Driver         string
}
