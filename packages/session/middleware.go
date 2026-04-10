package session

import (
	"math/rand/v2"
	"net/http"
	"time"
)

// StartSessionConfig configures the StartSession middleware.
type StartSessionConfig struct {
	// CookieName is the session cookie name (default: "session").
	CookieName string
	// Lifetime is the session cookie max-age (default: 2 hours).
	Lifetime time.Duration
	// Secure sets the Secure flag on the session cookie.
	Secure bool
	// SameSite is the SameSite attribute for the session cookie.
	SameSite http.SameSite
	// GCProbability controls the chance (0–100) of running GC on a request.
	GCProbability int
	// GCMaxLifetime is the max session age in seconds for GC.
	GCMaxLifetime int
}

// sessionResponseWriter intercepts WriteHeader / Write to flush the session
// cookie before the response headers are committed.
type sessionResponseWriter struct {
	http.ResponseWriter
	store   *Store
	cfg     StartSessionConfig
	flushed bool
}

func defaultConfig() StartSessionConfig {
	return StartSessionConfig{
		CookieName:    "session",
		Lifetime:      2 * time.Hour,
		SameSite:      http.SameSiteLaxMode,
		GCProbability: 2,
		GCMaxLifetime: 7200,
	}
}

func (w *sessionResponseWriter) WriteHeader(code int) {
	w.flush()
	w.ResponseWriter.WriteHeader(code)
}

func (w *sessionResponseWriter) Write(b []byte) (int, error) {
	w.flush()

	return w.ResponseWriter.Write(b)
}

func (w *sessionResponseWriter) flush() {
	if w.flushed {
		return
	}

	w.flushed = true
	maxAge := int(w.cfg.Lifetime.Seconds())
	http.SetCookie(w.ResponseWriter, &http.Cookie{
		Name:     w.cfg.CookieName,
		Value:    w.store.GetID(),
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   w.cfg.Secure,
		SameSite: w.cfg.SameSite,
	})
}

// StartSession is HTTP middleware that manages the full session lifecycle for
// each request: read → start → handle → save → write cookie.
func StartSession(handler Handler, cfg StartSessionConfig) func(http.Handler) http.Handler {
	if cfg.CookieName == "" {
		cfg = defaultConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Resolve session ID from incoming cookie.
			var id string

			if c, err := r.Cookie(cfg.CookieName); err == nil {
				id = c.Value
			}

			var store *Store

			if id != "" {
				store = NewWithID(cfg.CookieName, handler, id)
			} else {
				store = New(cfg.CookieName, handler)
			}

			// Inform existence-aware handlers whether the session already exists.
			if ea, ok := handler.(ExistenceAware); ok {
				ea.SetExists(id != "")
			}

			if err := store.Start(r.Context()); err != nil {
				http.Error(w, "session start failed", http.StatusInternalServerError)

				return
			}

			sw := &sessionResponseWriter{ResponseWriter: w, store: store, cfg: cfg}
			next.ServeHTTP(sw, r)

			if err := store.Save(r.Context()); err != nil {
				// Best-effort; session save errors should not abort the response.
				_ = err
			}

			// If the handler never wrote, flush the cookie now.
			sw.flush()

			// Probabilistic GC.
			if cfg.GCProbability > 0 && rand.IntN(100) < cfg.GCProbability {
				_ = handler.GC(r.Context(), cfg.GCMaxLifetime)
			}
		})
	}
}
