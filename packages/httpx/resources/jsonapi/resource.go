package jsonapi

import (
	"encoding/json"
	"net/http"
)

// Resource is the interface that JSON:API resource types must implement.
type Resource interface {
	// ToResourceObject produces the JSON:API resource object for the given
	// request. The returned map should contain "type", "id", "attributes",
	// and optionally "relationships", "links", and "meta".
	ToResourceObject(r *http.Request) map[string]any
}

// JsonApiResource is the base JSON:API resource. It wraps a value and uses
// caller-provided functions to produce the type, ID, attributes, and
// relationships required by the JSON:API specification.
type JsonApiResource[T any] struct {
	Value             T
	TypeFunc          func(T) string
	IDFunc            func(T) string
	AttributesFunc    func(T, *http.Request) map[string]any
	RelationshipsFunc func(T, *http.Request) []Relation
	LinksFunc         func(T, *http.Request) map[string]any
	MetaFunc          func(T, *http.Request) map[string]any
}

// NewJsonApiResource creates a JSON:API resource with the required type, ID,
// and attributes functions.
func NewJsonApiResource[T any](
	value T,
	typeFn func(T) string,
	idFn func(T) string,
	attrFn func(T, *http.Request) map[string]any,
) *JsonApiResource[T] {
	return &JsonApiResource[T]{
		Value:          value,
		TypeFunc:       typeFn,
		IDFunc:         idFn,
		AttributesFunc: attrFn,
	}
}

// WithRelationships sets the relationships function.
func (r *JsonApiResource[T]) WithRelationships(fn func(T, *http.Request) []Relation) *JsonApiResource[T] {
	r.RelationshipsFunc = fn

	return r
}

// WithLinks sets the links function.
func (r *JsonApiResource[T]) WithLinks(fn func(T, *http.Request) map[string]any) *JsonApiResource[T] {
	r.LinksFunc = fn

	return r
}

// WithMeta sets the meta function.
func (r *JsonApiResource[T]) WithMeta(fn func(T, *http.Request) map[string]any) *JsonApiResource[T] {
	r.MetaFunc = fn

	return r
}

// ToResourceObject produces the JSON:API resource object. Sparse fieldsets are
// applied when the request contains "fields[type]" query parameters.
func (r *JsonApiResource[T]) ToResourceObject(req *http.Request) map[string]any {
	typeName := r.TypeFunc(r.Value)
	id := r.IDFunc(r.Value)
	attrs := r.AttributesFunc(r.Value, req)

	// Apply sparse fieldsets.
	apiReq := NewRequest(req)

	if fields := apiReq.Fields(typeName); fields != nil {
		attrs = sparseFilter(attrs, fields)
	}

	// Resolve sentinel values in attributes.
	attrs = resolveAttributes(attrs)

	obj := map[string]any{
		"type":       typeName,
		"id":         id,
		"attributes": attrs,
	}

	if r.RelationshipsFunc != nil {
		relations := r.RelationshipsFunc(r.Value, req)
		resolver := NewRelationResolver(relations...)

		if rels := resolver.Resolve(); rels != nil {
			obj["relationships"] = rels
		}
	}

	if r.LinksFunc != nil {
		if links := r.LinksFunc(r.Value, req); len(links) > 0 {
			obj["links"] = links
		}
	}

	if r.MetaFunc != nil {
		if meta := r.MetaFunc(r.Value, req); len(meta) > 0 {
			obj["meta"] = meta
		}
	}

	return obj
}

// ToJSON serialises the resource as a JSON:API document with a "data" wrapper.
func (r *JsonApiResource[T]) ToJSON(req *http.Request) ([]byte, error) {
	doc := map[string]any{
		"data": r.ToResourceObject(req),
	}

	return json.Marshal(doc)
}

// Response writes the resource as a JSON:API HTTP response with the
// application/vnd.api+json content type.
func (r *JsonApiResource[T]) Response(w http.ResponseWriter, req *http.Request, status int) error {
	b, err := r.ToJSON(req)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	_, err = w.Write(b)

	return err
}

// sparseFilter keeps only the keys present in fields.
func sparseFilter(attrs map[string]any, fields []string) map[string]any {
	allowed := make(map[string]struct{}, len(fields))

	for _, f := range fields {
		allowed[f] = struct{}{}
	}

	result := make(map[string]any, len(fields))

	for k, v := range attrs {
		if _, ok := allowed[k]; ok {
			result[k] = v
		}
	}

	return result
}

// resolveAttributes processes attribute values, removing sentinel types like
// ConditionalValue and MissingValue.
func resolveAttributes(attrs map[string]any) map[string]any {
	result := make(map[string]any, len(attrs))

	for k, v := range attrs {
		if pm, ok := v.(interface{ IsMissing() bool }); ok && pm.IsMissing() {
			continue
		}

		result[k] = v
	}

	return result
}
