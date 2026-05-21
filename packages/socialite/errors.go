package socialite

import "errors"

// ErrInvalidState is returned when the OAuth state parameter does not match
// the value stored in the session. It mirrors Two\InvalidStateException.
var ErrInvalidState = errors.New("socialite: invalid or missing OAuth state")

// ErrMissingVerifier is returned by an OAuth1 provider when the callback
// request does not contain an oauth_verifier parameter.
var ErrMissingVerifier = errors.New("socialite: missing OAuth verifier")

// ErrMissingTemporaryCredentials is returned by an OAuth1 provider when the
// session does not contain the temporary credentials stored during the
// redirect step. It mirrors One\MissingTemporaryCredentialsException.
var ErrMissingTemporaryCredentials = errors.New("socialite: missing temporary OAuth credentials")
