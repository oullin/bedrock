package resources

import (
	"encoding/json"
	"net/http"
)

// Resource transforms a value into a JSON-serialisable map.
type Resource interface {
	ToMap(r *http.Request) map[string]any
}

// JsonResource is the base implementation of Resource. It wraps any value and
// provides a ToMap method for subtype customisation.
type JsonResource[T any] struct {
	Value   T
	With    map[string]any // additional top-level data
	WrapKey string         // wrapping key; empty string disables wrapping
	MapFunc func(T, *http.Request) map[string]any
}

// NewResource creates a JsonResource wrapping the given value.
func NewResource[T any](value T, mapFn func(T, *http.Request) map[string]any) *JsonResource[T] {
	return &JsonResource[T]{
		Value:   value,
		WrapKey: "data",
		MapFunc: mapFn,
	}
}

// ToMap serialises the resource value using the registered map function. It
// processes ConditionalValue, MergeValue and MissingValue sentinels.
func (r *JsonResource[T]) ToMap(req *http.Request) map[string]any {
	raw := r.MapFunc(r.Value, req)

	return resolveMap(raw)
}

// ToJSON serialises the resource to JSON bytes.
func (r *JsonResource[T]) ToJSON(req *http.Request) ([]byte, error) {
	data := r.toWrapped(req)

	return json.Marshal(data)
}

// Response writes the resource as a JSON HTTP response.
func (r *JsonResource[T]) Response(w http.ResponseWriter, req *http.Request, status int) error {
	b, err := r.ToJSON(req)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)

	return err
}

func (r *JsonResource[T]) toWrapped(req *http.Request) any {
	data := r.ToMap(req)

	if r.WrapKey != "" {
		wrapped := map[string]any{r.WrapKey: data}

		for k, v := range r.With {
			wrapped[k] = v
		}

		return wrapped
	}

	if len(r.With) > 0 {
		for k, v := range r.With {
			data[k] = v
		}
	}

	return data
}

// resolveMap processes a raw map, handling ConditionalValue, MergeValue and
// MissingValue sentinels.
func resolveMap(raw map[string]any) map[string]any {
	result := make(map[string]any, len(raw))

	for k, v := range raw {
		resolved := resolveValue(v)

		if resolved == nil {
			continue
		}

		// Handle MergeValue: merge its entries into the result.
		if mv, ok := resolved.(MergeValue); ok {
			for mk, mv := range mv.Data {
				result[mk] = mv
			}

			continue
		}

		result[k] = resolved
	}

	return result
}

// resolveValue unwraps sentinel types. Returns nil to signal omission.
func resolveValue(v any) any {
	switch val := v.(type) {
	case MissingValue:
		return nil
	case ConditionalValue:
		if val.IsMissing() {
			return nil
		}

		return resolveValue(val.Value)
	case MergeValue:
		return val
	case PotentiallyMissing:
		if val.IsMissing() {
			return nil
		}

		return v
	default:
		return v
	}
}
