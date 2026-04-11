package client

import (
	"io"
	"net/http"
	"strings"
	"sync"
)

// ResponseStub defines a canned response for faking HTTP requests.
type ResponseStub struct {
	Status  int
	Headers map[string]string
	Body    string
}

// ResponseSequence provides a sequence of stubbed responses. Each call to Next
// returns the next response in the sequence. When the sequence is exhausted it
// returns a 200 empty response.
type ResponseSequence struct {
	mu       sync.Mutex
	stubs    []ResponseStub
	index    int
	fallback *ResponseStub
}

// NewResponseSequence creates a sequence from the given stubs.
func NewResponseSequence(stubs ...ResponseStub) *ResponseSequence {
	return &ResponseSequence{stubs: stubs}
}

// Push appends a stub to the sequence.
func (s *ResponseSequence) Push(stubs ...ResponseStub) *ResponseSequence {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.stubs = append(s.stubs, stubs...)

	return s
}

// WhenEmpty sets a fallback response used when the sequence is exhausted.
func (s *ResponseSequence) WhenEmpty(stub ResponseStub) *ResponseSequence {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.fallback = &stub

	return s
}

// Next returns the next stubbed response.
func (s *ResponseSequence) Next() *http.Response {
	s.mu.Lock()

	defer s.mu.Unlock()

	var stub ResponseStub

	if s.index < len(s.stubs) {
		stub = s.stubs[s.index]
		s.index++
	} else if s.fallback != nil {
		stub = *s.fallback
	} else {
		stub = ResponseStub{Status: http.StatusOK}
	}

	return stubToResponse(stub)
}

// IsEmpty returns true when all stubs have been consumed.
func (s *ResponseSequence) IsEmpty() bool {
	s.mu.Lock()

	defer s.mu.Unlock()

	return s.index >= len(s.stubs)
}

func stubToResponse(stub ResponseStub) *http.Response {
	resp := &http.Response{
		StatusCode: stub.Status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(stub.Body)),
	}

	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}

	for k, v := range stub.Headers {
		resp.Header.Set(k, v)
	}

	return resp
}
