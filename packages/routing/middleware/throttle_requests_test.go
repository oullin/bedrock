package middleware

import (
	"errors"
	"testing"
)

// fakeThrottleRequest is a fake ThrottleRequest used by the tests.
type fakeThrottleRequest struct {
	ip   string
	user string
	path string
}

// Hit the bucket up to the cap.

// Next call should fail.

// User A and User B share an IP but should not share a bucket.

// fakeSignatureRequest implements SignatureValidator.
type fakeSignatureRequest struct{ valid bool }

func (r fakeThrottleRequest) IP() string     { return r.ip }
func (r fakeThrottleRequest) UserID() string { return r.user }
func (r fakeThrottleRequest) Path() string   { return r.path }

func TestThrottleRequests(t *testing.T) {
	t.Run("test_under_limit_passes", func(t *testing.T) {
		mw := NewThrottleRequests(nil)
		req := fakeThrottleRequest{ip: "1.1.1.1", path: "/api/x"}
		called := false
		_, err := mw.Handle(req, func(any) any { called = true; return nil }, 3, 1, "test")

		if err != nil {
			t.Fatal(err)
		}

		if !called {
			t.Error("next not called")
		}
	})

	t.Run("test_over_limit_blocks", func(t *testing.T) {
		mw := NewThrottleRequests(nil)
		req := fakeThrottleRequest{ip: "1.1.1.1", path: "/api/x"}

		for i := 0; i < 3; i++ {
			_, _ = mw.Handle(req, func(any) any { return nil }, 3, 1, "test")
		}

		_, err := mw.Handle(req, func(any) any { return nil }, 3, 1, "test")

		var tooMany *TooManyRequestsError

		if !errors.As(err, &tooMany) {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("test_user_keyed_independent_of_ip", func(t *testing.T) {
		mw := NewThrottleRequests(nil)

		a := fakeThrottleRequest{ip: "1.1.1.1", user: "alice", path: "/x"}
		b := fakeThrottleRequest{ip: "1.1.1.1", user: "bob", path: "/x"}

		for i := 0; i < 2; i++ {
			_, _ = mw.Handle(a, func(any) any { return nil }, 2, 1, "")
		}

		_, err := mw.Handle(a, func(any) any { return nil }, 2, 1, "")

		if err == nil {
			t.Fatal("alice should be throttled")
		}

		_, err = mw.Handle(b, func(any) any { return nil }, 2, 1, "")

		if err != nil {
			t.Errorf("bob should not be throttled: %v", err)
		}
	})
}

func (r fakeSignatureRequest) HasValidSignatureWhileIgnoring(ignore []string, absolute bool) bool {
	return r.valid
}

func TestValidateSignature(t *testing.T) {
	t.Run("test_valid_signature_passes", func(t *testing.T) {
		mw := &ValidateSignature{}
		called := false
		_, err := mw.Handle(fakeSignatureRequest{valid: true}, func(any) any { called = true; return nil })

		if err != nil {
			t.Fatal(err)
		}

		if !called {
			t.Error("next not called")
		}
	})

	t.Run("test_invalid_signature_errors", func(t *testing.T) {
		mw := &ValidateSignature{}
		_, err := mw.Handle(fakeSignatureRequest{valid: false}, func(any) any { return nil })

		if err == nil {
			t.Fatal("expected error")
		}
	})
}
