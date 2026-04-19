package middleware

import "net/http"

// statusInterceptor wraps an http.ResponseWriter to intercept
// WriteHeader calls. It converts 302 Found to 303 See Other for
// PUT/PATCH/DELETE Inertia requests, preventing browsers from
// replaying the original method on redirect.
type statusInterceptor struct {
	http.ResponseWriter
	request *http.Request
	written bool
}

func (si *statusInterceptor) WriteHeader(code int) {
	if si.written {
		return
	}

	si.written = true

	if code == http.StatusFound {
		switch si.request.Method {
		case http.MethodPut, http.MethodPatch, http.MethodDelete:
			code = http.StatusSeeOther
		}
	}

	si.ResponseWriter.WriteHeader(code)
}

func (si *statusInterceptor) Write(b []byte) (int, error) {
	if !si.written {
		si.WriteHeader(http.StatusOK)
	}

	return si.ResponseWriter.Write(b)
}
