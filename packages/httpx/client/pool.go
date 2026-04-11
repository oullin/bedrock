package client

import "sync"

// PoolCallback builds and sends a request using a PendingRequest.
type PoolCallback func(p *PendingRequest) (*Response, error)

// Pool executes multiple HTTP requests concurrently.
type Pool struct {
	factory *Factory
}

// NewPool creates a Pool backed by the given factory.

// As executes named requests concurrently and returns a map of name → response.

// Concurrent executes a slice of callbacks concurrently and returns results in
// the same order.

// PoolResult holds the response and error from a pool request.
type PoolResult struct {
	Response *Response
	Err      error
}

func NewPool(factory *Factory) *Pool {
	return &Pool{factory: factory}
}

func (p *Pool) As(requests map[string]PoolCallback) map[string]*PoolResult {
	var mu sync.Mutex

	results := make(map[string]*PoolResult, len(requests))

	var wg sync.WaitGroup

	for name, cb := range requests {
		wg.Add(1)

		go func(n string, fn PoolCallback) {
			defer wg.Done()

			pending := p.factory.PendingRequest()
			resp, err := fn(pending)

			mu.Lock()
			results[n] = &PoolResult{Response: resp, Err: err}
			mu.Unlock()
		}(name, cb)
	}

	wg.Wait()

	return results
}

func (p *Pool) Concurrent(callbacks []PoolCallback) []*PoolResult {
	results := make([]*PoolResult, len(callbacks))

	var wg sync.WaitGroup

	for i, cb := range callbacks {
		wg.Add(1)

		go func(idx int, fn PoolCallback) {
			defer wg.Done()

			pending := p.factory.PendingRequest()
			resp, err := fn(pending)
			results[idx] = &PoolResult{Response: resp, Err: err}
		}(i, cb)
	}

	wg.Wait()

	return results
}

// Ok returns true when the request succeeded with a 2xx status.
func (r *PoolResult) Ok() bool {
	return r.Err == nil && r.Response != nil && r.Response.Successful()
}
