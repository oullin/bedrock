package socialauth

import "errors"

// ErrInvalidState is returned when the OAuth state parameter does not match
// the value stored in the session.
var ErrInvalidState = errors.New("socialauth: invalid or missing OAuth state")

// ErrMissingVerifier is returned by an OAuth1 provider when the callback
// request does not contain an oauth_verifier parameter.
var ErrMissingVerifier = errors.New("socialauth: missing OAuth verifier")

// ErrMissingTemporaryCredentials is returned by an OAuth1 provider when the
// session does not contain the temporary credentials stored during the
// redirect step.
var ErrMissingTemporaryCredentials = errors.New("socialauth: missing temporary OAuth credentials")
