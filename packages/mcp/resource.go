package mcp

import "context"

// Resource is an MCP resource accessible at a static URI.
type Resource interface {
	Name() string
	Description() string
	URI() string
	MIMEType() string
	Read(ctx context.Context, req *Request) (*Response, error)
}

// ResourceTemplate is a Resource whose URI contains {varname} placeholders.
// The template string is used for matching and extracting URI variables.
type ResourceTemplate interface {
	Resource
	// URITemplate returns the URI template string, e.g. "file://users/{id}".
	URITemplate() string
}

// funcResource implements Resource from plain functions.
type funcResource struct {
	name        string
	description string
	uri         string
	mimeType    string
	reader      func(ctx context.Context, req *Request) (*Response, error)
}

func (r *funcResource) Name() string        { return r.name }
func (r *funcResource) Description() string { return r.description }
func (r *funcResource) URI() string         { return r.uri }
func (r *funcResource) MIMEType() string    { return r.mimeType }
func (r *funcResource) Read(ctx context.Context, req *Request) (*Response, error) {
	return r.reader(ctx, req)
}

// NewResource creates a Resource from plain functions.
func NewResource(
	name, description, uri, mimeType string,
	reader func(ctx context.Context, req *Request) (*Response, error),
) Resource {
	return &funcResource{
		name:        name,
		description: description,
		uri:         uri,
		mimeType:    mimeType,
		reader:      reader,
	}
}

// funcResourceTemplate implements ResourceTemplate from plain functions.
type funcResourceTemplate struct {
	funcResource
	uriTemplate string
}

func (r *funcResourceTemplate) URITemplate() string { return r.uriTemplate }

// NewResourceTemplate creates a ResourceTemplate. uriTemplate is the pattern
// string (e.g. "file://users/{id}").
func NewResourceTemplate(
	name, description, uriTemplate, mimeType string,
	reader func(ctx context.Context, req *Request) (*Response, error),
) ResourceTemplate {
	return &funcResourceTemplate{
		funcResource: funcResource{
			name:        name,
			description: description,
			uri:         uriTemplate,
			mimeType:    mimeType,
			reader:      reader,
		},
		uriTemplate: uriTemplate,
	}
}

// resourceToMap serialises a Resource for inclusion in a resources/list response.
func resourceToMap(r Resource) map[string]any {
	return map[string]any{
		"name":        r.Name(),
		"description": r.Description(),
		"uri":         r.URI(),
		"mimeType":    r.MIMEType(),
	}
}

// resourceTemplateToMap serialises a ResourceTemplate for resources/templates/list.
func resourceTemplateToMap(r ResourceTemplate) map[string]any {
	return map[string]any{
		"name":        r.Name(),
		"description": r.Description(),
		"uriTemplate": r.URITemplate(),
		"mimeType":    r.MIMEType(),
	}
}
