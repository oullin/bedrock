package jsonapi

import "fmt"

// ResourceIdentificationError is returned when a resource cannot determine its
// type or ID.
type ResourceIdentificationError struct {
	Field    string // "type" or "id"
	Resource any
}

func (e *ResourceIdentificationError) Error() string {
	return fmt.Sprintf("jsonapi: unable to resolve resource %s for %T", e.Field, e.Resource)
}
