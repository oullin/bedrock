package mcp

// ServerContext holds the resolved server configuration and primitive
// collections for use by method handlers. It is created fresh for each
// request from the Server's current state.
type ServerContext struct {
	SupportedVersions       []string
	ServerName              string
	ServerVersion           string
	Description             string
	Instructions            string
	MaxPaginationLength     int
	DefaultPaginationLength int

	tools     []Tool
	resources []Resource
	prompts   []Prompt
}

// ToolsList returns all registered tools.
func (c *ServerContext) ToolsList() []Tool { return c.tools }

// ResourcesList returns only static-URI resources (not templates).
func (c *ServerContext) ResourcesList() []Resource {
	out := make([]Resource, 0, len(c.resources))
	for _, r := range c.resources {
		if _, ok := r.(ResourceTemplate); !ok {
			out = append(out, r)
		}
	}
	return out
}

// ResourceTemplates returns only ResourceTemplate instances.
func (c *ServerContext) ResourceTemplates() []ResourceTemplate {
	out := make([]ResourceTemplate, 0)
	for _, r := range c.resources {
		if rt, ok := r.(ResourceTemplate); ok {
			out = append(out, rt)
		}
	}
	return out
}

// PromptsList returns all registered prompts.
func (c *ServerContext) PromptsList() []Prompt { return c.prompts }

// FindTool looks up a tool by name.
func (c *ServerContext) FindTool(name string) (Tool, bool) {
	for _, t := range c.tools {
		if t.Name() == name {
			return t, true
		}
	}
	return nil, false
}

// FindResource resolves a URI to a static resource or a template match.
// Returns the resource, the URI variable bindings (may be nil for static
// resources), and whether a match was found.
func (c *ServerContext) FindResource(uri string) (Resource, map[string]string, bool) {
	// Static resources first.
	for _, r := range c.resources {
		if _, isTemplate := r.(ResourceTemplate); !isTemplate {
			if r.URI() == uri {
				return r, nil, true
			}
		}
	}
	// Then URI templates.
	for _, r := range c.resources {
		rt, ok := r.(ResourceTemplate)
		if !ok {
			continue
		}
		tmpl, err := NewUriTemplate(rt.URITemplate())
		if err != nil {
			continue
		}
		if vars, matched := tmpl.Match(uri); matched {
			return rt, vars, true
		}
	}
	return nil, nil, false
}

// FindPrompt looks up a prompt by name.
func (c *ServerContext) FindPrompt(name string) (Prompt, bool) {
	for _, p := range c.prompts {
		if p.Name() == name {
			return p, true
		}
	}
	return nil, false
}

// PerPage clamps requested to [1, MaxPaginationLength], falling back to
// DefaultPaginationLength when requested is zero.
func (c *ServerContext) PerPage(requested int) int {
	if requested <= 0 {
		return c.DefaultPaginationLength
	}
	if requested > c.MaxPaginationLength {
		return c.MaxPaginationLength
	}
	return requested
}

// HasCompletions reports whether any registered primitive implements
// Completable, meaning the server should advertise the completions capability.
func (c *ServerContext) HasCompletions() bool {
	for _, t := range c.tools {
		if _, ok := t.(Completable); ok {
			return true
		}
	}
	for _, r := range c.resources {
		if _, ok := r.(Completable); ok {
			return true
		}
	}
	for _, p := range c.prompts {
		if _, ok := p.(Completable); ok {
			return true
		}
	}
	return false
}
