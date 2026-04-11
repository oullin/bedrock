package client_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/httpx/client"
)

func TestPromiseWait(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("async result"))
	}))

	defer server.Close()

	promise := client.NewFactory().PendingRequest().Async(http.MethodGet, server.URL)

	resp, err := promise.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "async result" {
		t.Fatalf("expected 'async result', got %s", resp.Body())
	}
}

func TestPromiseThen(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("callback"))
	}))

	defer server.Close()

	done := make(chan string, 1)

	client.NewFactory().PendingRequest().
		Async(http.MethodGet, server.URL).
		Then(func(resp *client.Response, err error) {
			done <- resp.Body()
		})

	select {
	case body := <-done:
		if body != "callback" {
			t.Fatalf("expected callback, got %s", body)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for callback")
	}
}

func TestPromiseResolve(t *testing.T) {
	t.Parallel()

	p := client.NewPromise()

	go func() {
		p.Resolve(makeResponse(200, "resolved"), nil)
	}()

	resp, err := p.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "resolved" {
		t.Fatalf("expected resolved, got %s", resp.Body())
	}
}

func TestPromiseCatchFiresOnError(t *testing.T) {
	t.Parallel()

	p := client.NewPromise()
	done := make(chan error, 1)

	p.Catch(func(err error) {
		done <- err
	})

	want := errors.New("boom")
	p.Resolve(nil, want)

	select {
	case got := <-done:
		if got != want {
			t.Fatalf("expected %v, got %v", want, got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for Catch callback")
	}
}

func TestPromiseCatchDoesNotFireOnSuccess(t *testing.T) {
	t.Parallel()

	p := client.NewPromise()
	fired := make(chan struct{}, 1)

	p.Catch(func(_ error) {
		fired <- struct{}{}
	})

	p.Resolve(makeResponse(200, "ok"), nil)

	select {
	case <-fired:
		t.Fatal("Catch should not fire on success")
	case <-time.After(100 * time.Millisecond):
		// expected
	}
}

func TestPromiseOtherwiseRecoversFromError(t *testing.T) {
	t.Parallel()

	p := client.NewPromise()

	recovered := p.Otherwise(func(_ error) (*client.Response, error) {
		return makeResponse(200, "recovered"), nil
	})

	p.Resolve(nil, errors.New("fail"))

	resp, err := recovered.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "recovered" {
		t.Fatalf("expected recovered, got %s", resp.Body())
	}
}

func TestPromiseOtherwisePassesThroughOnSuccess(t *testing.T) {
	t.Parallel()

	p := client.NewPromise()

	next := p.Otherwise(func(_ error) (*client.Response, error) {
		return makeResponse(200, "should not happen"), nil
	})

	p.Resolve(makeResponse(200, "original"), nil)

	resp, err := next.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "original" {
		t.Fatalf("expected original, got %s", resp.Body())
	}
}

func TestLazyPromiseDefersExecution(t *testing.T) {
	t.Parallel()

	executed := false

	lp := client.NewLazyPromise(func() (*client.Response, error) {
		executed = true

		return makeResponse(200, "lazy"), nil
	})

	// Should not have executed yet.
	if executed {
		t.Fatal("LazyPromise should not execute until Wait is called")
	}

	resp, err := lp.Wait()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !executed {
		t.Fatal("LazyPromise should have executed after Wait")
	}

	if resp.Body() != "lazy" {
		t.Fatalf("expected lazy, got %s", resp.Body())
	}
}

func TestLazyPromiseThen(t *testing.T) {
	t.Parallel()

	done := make(chan string, 1)

	lp := client.NewLazyPromise(func() (*client.Response, error) {
		return makeResponse(200, "deferred"), nil
	})

	lp.Then(func(resp *client.Response, err error) {
		done <- resp.Body()
	})

	select {
	case body := <-done:
		if body != "deferred" {
			t.Fatalf("expected deferred, got %s", body)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for LazyPromise Then callback")
	}
}

func TestLazyPromiseCatch(t *testing.T) {
	t.Parallel()

	done := make(chan error, 1)
	want := errors.New("lazy error")

	lp := client.NewLazyPromise(func() (*client.Response, error) {
		return nil, want
	})

	lp.Catch(func(err error) {
		done <- err
	})

	select {
	case got := <-done:
		if got != want {
			t.Fatalf("expected %v, got %v", want, got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for LazyPromise Catch callback")
	}
}
