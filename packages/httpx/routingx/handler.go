// Package routingx adapts the Bedrock routing package to net/http using httpx
// request and response primitives.
package routingx

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/bedrock/packages/httpx"
	"github.com/bedrock/packages/routing"
)

// NewHandler returns an http.Handler that dispatches requests through router.
func NewHandler(router *routing.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := httpx.NewRequest(r)
		req.SetRouteResolver(router)

		dispatch, err := router.Dispatch(req)

		if err != nil {
			writeError(w, err)

			return
		}

		applyRouteParameters(r, dispatch.Route)

		if err := writeResult(w, r, dispatch.Value); err != nil {
			writeError(w, err)
		}
	})
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, routing.ErrRouteNotFound) {
		http.NotFound(w, nil)

		return
	}

	var methodNotAllowed *routing.MethodNotAllowedError

	if errors.As(err, &methodNotAllowed) {
		w.Header().Set("Allow", strings.Join(methodNotAllowed.Allowed, ", "))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	var responseError *httpx.HttpResponseError

	if errors.As(err, &responseError) {
		applyHeaderMap(w.Header(), responseError.Headers)
		http.Error(w, responseError.Message, responseError.StatusCode)

		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func writeResult(w http.ResponseWriter, r *http.Request, result any) error {
	switch value := result.(type) {
	case nil:
		return httpx.NewResponse(w).NoContent()
	case http.Handler:
		value.ServeHTTP(w, r)

		return nil
	case func(http.ResponseWriter, *http.Request):
		value(w, r)

		return nil
	case *routing.HTTPResponse:
		return writeRoutingResponse(w, value)
	case routing.HTTPResponse:
		return writeRoutingResponse(w, &value)
	case *routing.RedirectResponse:
		return writeRoutingRedirect(w, r, value)
	case routing.RedirectResponse:
		return writeRoutingRedirect(w, r, &value)
	case *httpx.Response:
		return nil
	case string:
		return httpx.NewResponse(w).SendString(value)
	case []byte:
		return httpx.NewResponse(w).Send(value)
	default:
		if shouldWriteJSON(value) {
			return httpx.NewResponse(w).JSON(value)
		}

		return httpx.NewResponse(w).SendString(fmt.Sprint(value))
	}
}

func applyRouteParameters(r *http.Request, route *routing.Route) {
	if route == nil {
		return
	}

	for name, value := range route.ParametersWithoutNulls() {
		r.SetPathValue(name, value)
	}
}

func writeRoutingResponse(w http.ResponseWriter, response *routing.HTTPResponse) error {
	status := normalizeStatus(response.Status, http.StatusOK)
	resp := httpx.NewResponse(w).Status(status)
	applyHeaderValues(resp, response.Headers)

	switch body := response.Body.(type) {
	case nil:
		if response.Status == 0 {
			return resp.NoContent()
		}

		return resp.NoContent(status)
	case string:
		return resp.SendString(body)
	case []byte:
		return resp.Send(body)
	default:
		if shouldWriteJSON(body) {
			return resp.JSON(body)
		}

		return resp.SendString(fmt.Sprint(body))
	}
}

func writeRoutingRedirect(w http.ResponseWriter, r *http.Request, response *routing.RedirectResponse) error {
	if err := response.Send(); err != nil {
		return err
	}

	status := normalizeStatus(response.Status, http.StatusFound)

	applyHeaderMap(w.Header(), response.Headers)

	return httpx.NewRedirectResponse(w, r, response.URL, status).Send()
}

func normalizeStatus(status, fallback int) int {
	if status == 0 {
		return fallback
	}

	return status
}

func applyHeaderValues(response *httpx.Response, headers map[string][]string) {
	for key, values := range headers {
		if len(values) == 0 {
			continue
		}

		response.Header(key, values[0])

		for _, value := range values[1:] {
			response.Writer().Header().Add(key, value)
		}
	}
}

func applyHeaderMap(header http.Header, headers map[string][]string) {
	for key, values := range headers {
		for _, value := range values {
			header.Add(key, value)
		}
	}
}

func shouldWriteJSON(value any) bool {
	if value == nil {
		return false
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
		return true
	default:
		return false
	}
}
