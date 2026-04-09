package fortify

import "time"

// Config holds all Fortify configuration.
type Config struct {
	Features        Features
	Guard           string
	HomeURL         string
	LoginURL        string
	PasswordTimeout time.Duration
	LoginRateLimit  int
	LoginRateDecay  time.Duration
	IdentifierField string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Features:        DefaultFeatures(),
		Guard:           "web",
		HomeURL:         "/dashboard",
		LoginURL:        "/login",
		PasswordTimeout: 3 * time.Hour,
		LoginRateLimit:  5,
		LoginRateDecay:  time.Minute,
		IdentifierField: "email",
	}
}
