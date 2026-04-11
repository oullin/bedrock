package cookie

import (
	"net/http"
	"sync"
)

// Jar implements QueueingFactory. It maintains a queue of cookies to be
// attached to outgoing responses and applies configurable defaults.
// Cookies are keyed by name and path, mirroring Laravel's CookieJar
// behaviour. It is safe for concurrent use.
type Jar struct {
	mu       sync.Mutex
	queued   map[string]map[string]*http.Cookie // name -> path -> cookie
	defaults Options
}

// NewJar creates a Jar with the given default options.
func NewJar(defaults Options) *Jar {
	return &Jar{
		queued:   make(map[string]map[string]*http.Cookie),
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

// Queue adds a cookie to the outgoing queue, keyed by name and path.
// Any previously queued cookie with the same name and path is replaced.
func (j *Jar) Queue(c *http.Cookie) {
	j.mu.Lock()

	defer j.mu.Unlock()

	if j.queued[c.Name] == nil {
		j.queued[c.Name] = make(map[string]*http.Cookie)
	}

	j.queued[c.Name][c.Path] = c
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

// Unqueue removes cookies from the queue. When called with only a name,
// all path entries for that name are removed. When called with a path,
// only that specific name+path entry is removed; if the name bucket
// becomes empty it is cleaned up.
func (j *Jar) Unqueue(name string, path ...string) {
	j.mu.Lock()

	defer j.mu.Unlock()

	if len(path) == 0 {
		delete(j.queued, name)

		return
	}

	bucket := j.queued[name]

	if bucket == nil {
		return
	}

	delete(bucket, path[0])

	if len(bucket) == 0 {
		delete(j.queued, name)
	}
}

// HasQueued reports whether a cookie with the given name is queued.
// When called with a path, it checks for the specific name+path entry.
func (j *Jar) HasQueued(name string, path ...string) bool {
	return j.Queued(name, path...) != nil
}

// Queued returns the queued cookie with the given name. When called
// without a path, the last entry for the name is returned (matching
// Laravel's behaviour). When called with a path, the specific
// name+path entry is returned.
func (j *Jar) Queued(name string, path ...string) *http.Cookie {
	j.mu.Lock()

	defer j.mu.Unlock()

	bucket := j.queued[name]

	if bucket == nil {
		return nil
	}

	if len(path) > 0 {
		return bucket[path[0]]
	}

	// Return the last entry (iteration order is non-deterministic in Go,
	// but consistent with "return any/last" semantics).
	var last *http.Cookie

	for _, c := range bucket {
		last = c
	}

	return last
}

// GetQueued returns all queued cookies in a flat slice.
func (j *Jar) GetQueued() []*http.Cookie {
	j.mu.Lock()

	defer j.mu.Unlock()

	var cookies []*http.Cookie

	for _, bucket := range j.queued {
		for _, c := range bucket {
			cookies = append(cookies, c)
		}
	}

	return cookies
}

// Flush clears all queued cookies.
func (j *Jar) Flush() {
	j.mu.Lock()

	defer j.mu.Unlock()

	j.queued = make(map[string]map[string]*http.Cookie)
}

// merge applies non-zero fields from opts on top of the jar's defaults.
// For boolean fields (Secure, HTTPOnly, Raw), an explicit value in opts
// takes precedence over the default. Because Go's zero-value for bool is
// false, passing Secure=false in opts cannot be distinguished from "not
// set" without a pointer or sentinel. We follow the same approach as
// Laravel: the caller-provided value wins when it differs from the
// default, and the default wins otherwise. Practically this means the
// default is used unless the caller explicitly sets the field to true.
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
