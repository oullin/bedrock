package cookie

import "net/http"

// Factory creates cookies.
type Factory interface {
	Make(name, value string, opts Options) *http.Cookie
	Forever(name, value string, opts Options) *http.Cookie
	Forget(name string, opts Options) *http.Cookie
}

// QueueingFactory extends Factory with a queue for deferred attachment to
// HTTP responses.
type QueueingFactory interface {
	Factory
	Queue(c *http.Cookie)
	Unqueue(name string)
	HasQueued(name string) bool
	Queued(name string) *http.Cookie
	GetQueued() []*http.Cookie
	Flush()
}
