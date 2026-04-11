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

// Resolve sets the response and error, unblocking any waiters.

// Wait blocks until the promise is resolved and returns the response and error.

// Then registers a callback invoked when the promise resolves. The callback
// runs in a new goroutine.

// Catch registers a callback invoked only when the promise resolves with an
// error. The callback runs in a new goroutine.

// Otherwise registers a fallback that can recover from an error. If the
// promise resolves with an error the callback is invoked and its return values
// replace the original response and error. The callback runs in a new
// goroutine; callers that need the recovered value should chain with Then or
// call Wait.

// LazyPromise defers execution of a request function until Wait is called.
type LazyPromise struct {
	once sync.Once
	fn   func() (*Response, error)
	p    *Promise
}

func NewPromise() *Promise {
	return &Promise{done: make(chan struct{})}
}

func (p *Promise) Resolve(resp *Response, err error) {
	p.once.Do(func() {
		p.response = resp
		p.err = err
		close(p.done)
	})
}

func (p *Promise) Wait() (*Response, error) {
	<-p.done

	return p.response, p.err
}

func (p *Promise) Then(fn func(*Response, error)) *Promise {
	go func() {
		<-p.done
		fn(p.response, p.err)
	}()

	return p
}

func (p *Promise) Catch(fn func(error)) *Promise {
	go func() {
		<-p.done

		if p.err != nil {
			fn(p.err)
		}
	}()

	return p
}

func (p *Promise) Otherwise(fn func(error) (*Response, error)) *Promise {
	next := NewPromise()

	go func() {
		<-p.done

		if p.err != nil {
			resp, err := fn(p.err)
			next.Resolve(resp, err)
		} else {
			next.Resolve(p.response, nil)
		}
	}()

	return next
}

// NewLazyPromise creates a promise that will not execute fn until Wait or Then
// is called.
func NewLazyPromise(fn func() (*Response, error)) *LazyPromise {
	return &LazyPromise{
		fn: fn,
		p:  NewPromise(),
	}
}

func (lp *LazyPromise) build() {
	lp.once.Do(func() {
		resp, err := lp.fn()
		lp.p.Resolve(resp, err)
	})
}

// Wait triggers the deferred execution and blocks until the result is
// available.
func (lp *LazyPromise) Wait() (*Response, error) {
	lp.build()

	return lp.p.Wait()
}

// Then triggers the deferred execution and registers a callback for when it
// completes.
func (lp *LazyPromise) Then(fn func(*Response, error)) *LazyPromise {
	lp.build()
	lp.p.Then(fn)

	return lp
}

// Catch triggers the deferred execution and registers an error-only callback.
func (lp *LazyPromise) Catch(fn func(error)) *LazyPromise {
	lp.build()
	lp.p.Catch(fn)

	return lp
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
