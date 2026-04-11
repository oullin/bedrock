package httpx

// Old retrieves a previously flashed input value from the session. Returns the
// optional fallback (or nil) when no session is attached or the key is absent.
func (r *Request) Old(key string, fallback ...any) any {
	if r.session == nil {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return nil
	}

	var fb any

	if len(fallback) > 0 {
		fb = fallback[0]
	}

	return r.session.GetOldInput(key, fb)
}

// HasOld returns true when the session has old input for the given key.
func (r *Request) HasOld(key string) bool {
	if r.session == nil {
		return false
	}

	return r.session.HasOldInput(key)
}

// Flash flashes all current input to the session.
func (r *Request) Flash() {
	if r.session == nil {
		return
	}

	r.session.FlashInput(r.All())
}

// FlashOnly flashes only the specified input keys to the session.
func (r *Request) FlashOnly(keys ...string) {
	if r.session == nil {
		return
	}

	r.session.FlashInput(r.Only(keys...))
}

// FlashExcept flashes all input except the specified keys to the session.
func (r *Request) FlashExcept(keys ...string) {
	if r.session == nil {
		return
	}

	r.session.FlashInput(r.Except(keys...))
}

// Flush removes all old input from the session.
func (r *Request) Flush() {
	if r.session == nil {
		return
	}

	r.session.Remove("_old_input")
}
