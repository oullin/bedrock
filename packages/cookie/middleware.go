package cookie

import (
	"net/http"
)

// EncryptCookies is middleware that decrypts incoming request cookies and
// encrypts outgoing response cookies using the provided Encrypter.
// Cookies whose names appear in the Except list are passed through unmodified.
type EncryptCookies struct {
	enc    Encrypter
	except map[string]bool
}

// NewEncryptCookies creates EncryptCookies middleware.

// Wrap returns an http.Handler that decrypts request cookies before calling
// next, then encrypts response cookies set by next.

// decryptRequest returns a copy of r with encrypted cookies decrypted.

// Keep original value on decryption failure.

// encryptingResponseWriter intercepts Set-Cookie headers and encrypts values.
type encryptingResponseWriter struct {
	http.ResponseWriter
	enc    Encrypter
	except map[string]bool
}

// AttachQueued is middleware that flushes a Jar's queued cookies onto the
// outgoing response.
type AttachQueued struct {
	jar *Jar
}

// NewAttachQueued creates AttachQueued middleware.

// Wrap returns an http.Handler that attaches queued cookies to the response.
// Cookies are written before the handler runs, allowing them to be overridden
// by the handler if needed.

// queueingResponseWriter intercepts the first write/header flush to inject
// queued cookies before the response headers are sent.
type queueingResponseWriter struct {
	http.ResponseWriter
	jar     *Jar
	flushed bool
}

func NewEncryptCookies(enc Encrypter, except ...string) *EncryptCookies {
	m := &EncryptCookies{enc: enc, except: make(map[string]bool, len(except))}

	for _, name := range except {
		m.except[name] = true
	}

	return m
}

func (m *EncryptCookies) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = m.decryptRequest(r)
		rw := &encryptingResponseWriter{ResponseWriter: w, enc: m.enc, except: m.except}
		next.ServeHTTP(rw, r)
	})
}

func (m *EncryptCookies) decryptRequest(r *http.Request) *http.Request {
	r2 := r.Clone(r.Context())
	r2.Header.Del("Cookie")

	for _, c := range r.Cookies() {
		if m.except[c.Name] {
			r2.AddCookie(c)

			continue
		}

		plain, err := m.enc.Decrypt(c.Value)

		if err != nil {

			r2.AddCookie(c)

			continue
		}

		nc := *c
		nc.Value = plain
		r2.AddCookie(&nc)
	}

	return r2
}

func (w *encryptingResponseWriter) WriteHeader(code int) {
	w.encryptResponseCookies()
	w.ResponseWriter.WriteHeader(code)
}

func (w *encryptingResponseWriter) Write(b []byte) (int, error) {
	w.encryptResponseCookies()

	return w.ResponseWriter.Write(b)
}

func (w *encryptingResponseWriter) encryptResponseCookies() {
	header := w.Header()
	cookies := header["Set-Cookie"]

	if len(cookies) == 0 {
		return
	}

	parsed := (&http.Response{Header: header}).Cookies()
	header.Del("Set-Cookie")

	for _, c := range parsed {
		if w.except[c.Name] {
			header.Add("Set-Cookie", c.String())

			continue
		}

		enc, err := w.enc.Encrypt(c.Value)

		if err != nil {
			header.Add("Set-Cookie", c.String())

			continue
		}

		nc := *c
		nc.Value = enc
		header.Add("Set-Cookie", nc.String())
	}
}

func NewAttachQueued(jar *Jar) *AttachQueued {
	return &AttachQueued{jar: jar}
}

func (m *AttachQueued) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &queueingResponseWriter{ResponseWriter: w, jar: m.jar}
		next.ServeHTTP(rw, r)
		rw.flush()
	})
}

func (w *queueingResponseWriter) flush() {
	if w.flushed {
		return
	}

	w.flushed = true

	for _, c := range w.jar.GetQueued() {
		http.SetCookie(w.ResponseWriter, c)
	}

	w.jar.Flush()
}

func (w *queueingResponseWriter) WriteHeader(code int) {
	w.flush()
	w.ResponseWriter.WriteHeader(code)
}

func (w *queueingResponseWriter) Write(b []byte) (int, error) {
	w.flush()

	return w.ResponseWriter.Write(b)
}
