package cookie

import (
	"net/http"
	"sync"
)

// Jar implements QueueingFactory. It maintains a queue of cookies to be
// attached to outgoing responses and applies configurable defaults.
// It is safe for concurrent use.
type Jar struct {
	mu       sync.Mutex
	queued   map[string]*http.Cookie
	defaults Options
}

// NewJar creates a Jar with the given default options.
func NewJar(defaults Options) *Jar {
	return &Jar{
		queued:   make(map[string]*http.Cookie),
		defaults: defaults,
	}
}

// SetDefaults replaces the default cookie options.
func (j *Jar) SetDefaults(opts Options) {
	j.mu.Lock()

	defer j.mu.Unlock()

	j.defaults = opts
}

// Defaults returns the current default options.
func (j *Jar) Defaults() Options {
	j.mu.Lock()

	defer j.mu.Unlock()

	return j.defaults
}

// Make creates an *http.Cookie using the given options merged on top of
// the jar's defaults. Fields in opts with zero values inherit from defaults.
func (j *Jar) Make(name, value string, opts Options) *http.Cookie {
	return Make(name, value, j.merge(opts))
}

// Forever creates a 400-day cookie.
func (j *Jar) Forever(name, value string, opts Options) *http.Cookie {
	return Forever(name, value, j.merge(opts))
}

// Forget creates an expired cookie to delete the named cookie from the client.
func (j *Jar) Forget(name string, opts Options) *http.Cookie {
	return Forget(name, j.merge(opts))
}

// Queue adds a cookie to the outgoing queue. Any previously queued cookie
// with the same name is replaced.
func (j *Jar) Queue(c *http.Cookie) {
	j.mu.Lock()

	defer j.mu.Unlock()

	j.queued[c.Name] = c
}

// QueueMake creates a cookie from name, value, and opts, then queues it.
func (j *Jar) QueueMake(name, value string, opts Options) {
	j.Queue(j.Make(name, value, opts))
}

// QueueForever creates a 400-day cookie and queues it.
func (j *Jar) QueueForever(name, value string, opts Options) {
	j.Queue(j.Forever(name, value, opts))
}

// Expire queues a deletion cookie for the named cookie.
func (j *Jar) Expire(name string, opts Options) {
	j.Queue(j.Forget(name, opts))
}

// Unqueue removes the cookie with the given name from the queue.
func (j *Jar) Unqueue(name string) {
	j.mu.Lock()

	defer j.mu.Unlock()

	delete(j.queued, name)
}

// HasQueued reports whether a cookie with the given name is queued.
func (j *Jar) HasQueued(name string) bool {
	j.mu.Lock()

	defer j.mu.Unlock()

	_, ok := j.queued[name]

	return ok
}

// Queued returns the queued cookie with the given name, or nil.
func (j *Jar) Queued(name string) *http.Cookie {
	j.mu.Lock()

	defer j.mu.Unlock()

	return j.queued[name]
}

// GetQueued returns all queued cookies in an unspecified order.
func (j *Jar) GetQueued() []*http.Cookie {
	j.mu.Lock()

	defer j.mu.Unlock()

	cookies := make([]*http.Cookie, 0, len(j.queued))

	for _, c := range j.queued {
		cookies = append(cookies, c)
	}

	return cookies
}

// Flush clears all queued cookies.
func (j *Jar) Flush() {
	j.mu.Lock()

	defer j.mu.Unlock()

	j.queued = make(map[string]*http.Cookie)
}

// merge applies non-zero fields from opts on top of the jar's defaults.
func (j *Jar) merge(opts Options) Options {
	d := j.defaults

	if opts.Path != "" {
		d.Path = opts.Path
	}

	if opts.Domain != "" {
		d.Domain = opts.Domain
	}

	if opts.MaxAge != 0 {
		d.MaxAge = opts.MaxAge
	}

	if opts.SameSite != 0 {
		d.SameSite = opts.SameSite
	}

	d.Secure = d.Secure || opts.Secure
	d.HTTPOnly = d.HTTPOnly || opts.HTTPOnly
	d.Raw = d.Raw || opts.Raw

	return d
}
