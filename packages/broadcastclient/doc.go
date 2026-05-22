// Package broadcastclient provides a Go client for the BroadcastClient JavaScript library.
// It provides real-time event broadcasting abstractions over multiple
// transport backends (Pusher, Socket.IO, Null/stub) with a uniform
// Channel and Connector interface.
//
// The main entry point is the BroadcastClient struct. Construct one via New(),
// supplying Options that specify the broadcaster and namespace.
//
// Example:
//
//	e, err := broadcastclient.New(broadcastclient.Options{Broadcaster: "null"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	ch := e.Channel("orders")
//	ch.Listen("OrderShipped", func(data any) { fmt.Println(data) })
package broadcastclient
