package client

import "sync"

// Promise represents an asynchronous HTTP response. It is resolved when the
// request completes.
type Promise struct {
	once     sync.Once
	done     chan struct{}
	response *Response
	err      error
}

// NewPromise creates a Promise that will be resolved by calling Resolve.
func NewPromise() *Promise {
	return &Promise{done: make(chan struct{})}
}

// Resolve sets the response and error, unblocking any waiters.
func (p *Promise) Resolve(resp *Response, err error) {
	p.once.Do(func() {
		p.response = resp
		p.err = err
		close(p.done)
	})
}

// Wait blocks until the promise is resolved and returns the response and error.
func (p *Promise) Wait() (*Response, error) {
	<-p.done

	return p.response, p.err
}

// Then registers a callback invoked when the promise resolves. The callback
// runs in a new goroutine.
func (p *Promise) Then(fn func(*Response, error)) *Promise {
	go func() {
		<-p.done
		fn(p.response, p.err)
	}()

	return p
}

// Async sends a request asynchronously and returns a Promise.
func (p *PendingRequest) Async(method, url string, data ...any) *Promise {
	promise := NewPromise()

	go func() {
		var d any

		if len(data) > 0 {
			d = data[0]
		}

		resp, err := p.send(method, url, d)
		promise.Resolve(resp, err)
	}()

	return promise
}
