// Package echo provides a Go client for the Echo JavaScript library.
// It provides real-time event broadcasting abstractions over multiple
// transport backends (Pusher, Socket.IO, Null/stub) with a uniform
// Channel and Connector interface.
//
// The main entry point is the Echo struct. Construct one via New(),
// supplying Options that specify the broadcaster and namespace.
//
// Example:
//
//	e, err := echo.New(echo.Options{Broadcaster: "null"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	ch := e.Channel("orders")
//	ch.Listen("OrderShipped", func(data any) { fmt.Println(data) })
package echo
