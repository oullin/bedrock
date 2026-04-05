package configuration

import (
	"fmt"
	nethttp "net/http"
)

// Exceptions configures error reporting and rendering.
type Exceptions struct {
	reporter func(error)
	renderer func(nethttp.ResponseWriter, *nethttp.Request, error)
}

// NewExceptions creates a default exceptions configuration.
func NewExceptions() *Exceptions {
	return &Exceptions{
		reporter: func(error) {},
		renderer: func(writer nethttp.ResponseWriter, _ *nethttp.Request, err error) {
			nethttp.Error(writer, err.Error(), nethttp.StatusInternalServerError)
		},
	}
}

// ReportUsing overrides the default error reporter.
func (e *Exceptions) ReportUsing(reporter func(error)) {
	if reporter != nil {
		e.reporter = reporter
	}
}

// RenderUsing overrides the default error renderer.
func (e *Exceptions) RenderUsing(renderer func(nethttp.ResponseWriter, *nethttp.Request, error)) {
	if renderer != nil {
		e.renderer = renderer
	}
}

// Wrap applies panic recovery to an HTTP handler.
func (e *Exceptions) Wrap(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(writer nethttp.ResponseWriter, request *nethttp.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err := fmt.Errorf("foundation: panic recovered: %v", recovered)
				e.reporter(err)
				e.renderer(writer, request, err)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}
