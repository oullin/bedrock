package props

// OnceProp is included on the first visit. Subsequent requests that list
// the prop key in X-Inertia-Except-Once-Props will exclude it, allowing
// the client to reuse the previously loaded value.
type OnceProp struct {
	Value     any
	expiresAt *int64
}

// Once creates a prop resolved only on the first visit.
func Once(val any) OnceProp {
	return OnceProp{Value: val}
}

// ExpiresAt returns a copy of the OnceProp with a Unix timestamp after
// which the client should discard its cached value.
func (p OnceProp) ExpiresAt(unix int64) OnceProp {
	p.expiresAt = &unix

	return p
}

// GetExpiresAt returns the expiration timestamp, or nil if unset.
func (p OnceProp) GetExpiresAt() *int64 { return p.expiresAt }

func (p OnceProp) Prop() any { return p.Value }
