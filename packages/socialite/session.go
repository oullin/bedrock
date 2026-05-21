package socialite

// Session is the minimal session contract required by OAuth providers to
// store and retrieve the state token and PKCE code verifier across requests.
// It is satisfied by the bedrock session store (session.Store) as well as any
// custom adapter.
type Session interface {
	// Put stores a value under the given key.
	Put(key string, value any)

	// Pull retrieves a value by key and removes it from the session in one
	// atomic step. Returns nil if the key does not exist.
	// It mirrors @bedrock\Contracts\Session\Session::pull().
	Pull(key string) any

	// Get retrieves a value by key without removing it. Returns nil if the
	// key does not exist.
	Get(key string) any
}
