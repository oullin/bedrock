package routing

import "strings"

// ResourceOptions configures resource route registration.
type ResourceOptions struct {
	Only              []string
	Except            []string
	Names             map[string]string
	Parameters        map[string]string
	Middleware        []MiddlewareFunc
	WithoutMiddleware []MiddlewareFunc
	Shallow           bool
}

// ResourceRegistrar registers RESTful resource routes.
type ResourceRegistrar struct {
	router *Router
}

// NewResourceRegistrar creates a resource registrar for the given router.
func NewResourceRegistrar(router *Router) *ResourceRegistrar {
	return &ResourceRegistrar{router: router}
}

// Register creates standard CRUD routes for a resource (7 actions).
func (rr *ResourceRegistrar) Register(name string, handlers ResourceHandlers, opts ...ResourceOptions) {
	opt := mergeResourceOptions(opts)
	base := "/" + strings.Trim(name, "/")
	param := rr.parameterName(name, opt)
	item := base + "/{" + param + "}"

	if rr.shouldInclude("index", opt) && handlers.Index != nil {
		route := rr.router.Get(base, handlers.Index)
		rr.applyOptions(route, "index", name, opt)
	}

	if rr.shouldInclude("create", opt) && handlers.Create != nil {
		route := rr.router.Get(base+"/create", handlers.Create)
		rr.applyOptions(route, "create", name, opt)
	}

	if rr.shouldInclude("store", opt) && handlers.Store != nil {
		route := rr.router.Post(base, handlers.Store)
		rr.applyOptions(route, "store", name, opt)
	}

	if rr.shouldInclude("show", opt) && handlers.Show != nil {
		route := rr.router.Get(item, handlers.Show)
		rr.applyOptions(route, "show", name, opt)
	}

	if rr.shouldInclude("edit", opt) && handlers.Edit != nil {
		route := rr.router.Get(item+"/edit", handlers.Edit)
		rr.applyOptions(route, "edit", name, opt)
	}

	if rr.shouldInclude("update", opt) && handlers.Update != nil {
		route := rr.router.Put(item, handlers.Update)
		rr.applyOptions(route, "update", name, opt)
	}

	if rr.shouldInclude("destroy", opt) && handlers.Destroy != nil {
		route := rr.router.Delete(item, handlers.Destroy)
		rr.applyOptions(route, "destroy", name, opt)
	}
}

// APIResource creates API resource routes (5 actions: no create/edit).
func (rr *ResourceRegistrar) APIResource(name string, handlers ResourceHandlers, opts ...ResourceOptions) {
	opt := mergeResourceOptions(opts)

	if len(opt.Except) == 0 && len(opt.Only) == 0 {
		opt.Except = []string{"create", "edit"}
	} else {
		opt.Except = append(opt.Except, "create", "edit")
	}

	rr.Register(name, handlers, opt)
}

// Singleton creates singleton resource routes (3 actions: show, edit, update).
func (rr *ResourceRegistrar) Singleton(name string, handlers SingletonHandlers, opts ...ResourceOptions) {
	opt := mergeResourceOptions(opts)
	base := "/" + strings.Trim(name, "/")

	if rr.shouldInclude("show", opt) && handlers.Show != nil {
		route := rr.router.Get(base, handlers.Show)
		rr.applyOptions(route, "show", name, opt)
	}

	if rr.shouldInclude("edit", opt) && handlers.Edit != nil {
		route := rr.router.Get(base+"/edit", handlers.Edit)
		rr.applyOptions(route, "edit", name, opt)
	}

	if rr.shouldInclude("update", opt) && handlers.Update != nil {
		route := rr.router.Put(base, handlers.Update)
		rr.applyOptions(route, "update", name, opt)
	}

	if rr.shouldInclude("destroy", opt) && handlers.Destroy != nil {
		route := rr.router.Delete(base, handlers.Destroy)
		rr.applyOptions(route, "destroy", name, opt)
	}
}

// APISingleton creates API singleton resource routes (2 actions: show, update).
func (rr *ResourceRegistrar) APISingleton(name string, handlers SingletonHandlers, opts ...ResourceOptions) {
	opt := mergeResourceOptions(opts)

	if len(opt.Except) == 0 && len(opt.Only) == 0 {
		opt.Except = []string{"create", "edit"}
	} else {
		opt.Except = append(opt.Except, "create", "edit")
	}

	rr.Singleton(name, handlers, opt)
}

func (rr *ResourceRegistrar) shouldInclude(action string, opt ResourceOptions) bool {
	if len(opt.Only) > 0 {
		return contains(opt.Only, action)
	}

	if len(opt.Except) > 0 {
		return !contains(opt.Except, action)
	}

	return true
}

func (rr *ResourceRegistrar) parameterName(resource string, opt ResourceOptions) string {
	if opt.Parameters != nil {
		if p, ok := opt.Parameters[resource]; ok {
			return p
		}
	}

	parts := strings.Split(resource, "/")

	return singularize(parts[len(parts)-1])
}

func (rr *ResourceRegistrar) applyOptions(route *Route, action, resource string, opt ResourceOptions) {
	routeName := resource + "." + action

	if opt.Names != nil {
		if custom, ok := opt.Names[action]; ok {
			routeName = custom
		}
	}

	route.Name(routeName)

	if len(opt.Middleware) > 0 {
		route.Middleware(opt.Middleware...)
	}

	if len(opt.WithoutMiddleware) > 0 {
		route.WithoutMiddleware(opt.WithoutMiddleware...)
	}
}

func mergeResourceOptions(opts []ResourceOptions) ResourceOptions {
	if len(opts) == 0 {
		return ResourceOptions{}
	}

	return opts[0]
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

// singularize performs a basic singularization by removing trailing "s".
func singularize(s string) string {
	if strings.HasSuffix(s, "ies") {
		return s[:len(s)-3] + "y"
	}

	if strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "xes") || strings.HasSuffix(s, "zes") {
		return s[:len(s)-2]
	}

	if strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss") {
		return s[:len(s)-1]
	}

	return s
}
