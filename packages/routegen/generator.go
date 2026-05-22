package routegen

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// generator holds the mutable state accumulated during a single Generate call.
type generator struct {
	opts    Options
	content map[string][]string            // abs path → ordered content chunks
	imports map[string]map[string][]string // abs path → importFrom → identifiers
}

// ─────────────────────────────────────────────────────────────────────────────
// Actions
// ─────────────────────────────────────────────────────────────────────────────

// generateActions builds the actions/ directory content.

// Only routes backed by a named controller.

// Group by dot-namespace (one group = one controller file).

// Write one file per controller.

// Write barrel index.ts files bottom-up.

// writeControllerFile generates the exports for a single controller file.

// Group routes by JS method name to detect multi-route-same-action cases.

// Invokable controllers are exported via `export default`, not `export const`.

// Determine whether any route in this file is invokable.

// e.g. "PostController"

// Invokable: attach sibling methods, then default-export the invokable.

// Regular controller: build an object with safe method names + original aliases.

// Add original→safe alias when the name was sanitised.

// Avoid duplicates.

// writeMethodExport generates a single route function export.
// shouldExport=false means a var declaration without "export".
// method.template.ts template.

// writeNamedMethodExport generates a route function export under its named-route name.
// Always exported (shouldExport=true).

// writeMultiRouteExport handles the case where multiple routes share the same
// controller action. It generates a keyed dictionary.

// Generate a temporary (unexported) helper for each route keyed by URI hash.

// Temporarily override the method name with the temp name.

// renderMethod builds the complete TypeScript block for a single route action.

// ── docblock ────────────────────────────────────────────────────────────

// ── main function ────────────────────────────────────────────────────────

// ── definition ──────────────────────────────────────────────────────────

// ── docblock for .url ────────────────────────────────────────────────────

// ── .url function ────────────────────────────────────────────────────────

// Single-param: allow primitive shorthand.

// Custom key resolution for single param.

// Array unpacking.

// applyUrlDefaults.

// validateParameters for optional params.

// parsedArgs object.

// Custom binding key: check if the param value is an object with the key.

// URL construction via chained .replace calls.

// No params: just strip trailing slash and append query params.

// ── per-verb shorthand methods ───────────────────────────────────────────

// ── form helpers (--with-form) ───────────────────────────────────────────

// renderFormHelpers generates the .form sub-object for a route method.

// Non-safe verb: inject _method spoofing.

// functionArgs renders the TypeScript parameter list for a generated function.

// Object form: { param: type, ... }

// Array form: [param: type, ...]

// Single-param primitive shorthand.

// ─────────────────────────────────────────────────────────────────────────────
// Barrel files (actions)
// ─────────────────────────────────────────────────────────────────────────────

// writeBarrelFiles writes index.ts barrel files for all directory levels
// within the actions/ tree.

// Collect unique namespace segments to build the directory tree.

type nsNode map[string]interface{} // either nsNode (subtree) or []*RouteInfo (leaf)

// buildNamespaceTree converts a flat dot-namespace map into a nested tree.

// writeBarrelNode recursively writes index.ts files for each directory node.

// Collect child names sorted for deterministic output.

// Import each child module.

// Recurse into sub-nodes.

// ─────────────────────────────────────────────────────────────────────────────
// Named routes
// ─────────────────────────────────────────────────────────────────────────────

// generateRoutes builds the routes/ directory content.

// Only named routes.

// Write one file per name (route functions land in the prefix/index.ts).

// Write barrel files based on the name tree.

type nameNode map[string]interface{}

func newGenerator(opts Options) *generator {
	return &generator{
		opts:    opts,
		content: make(map[string][]string),
		imports: make(map[string]map[string][]string),
	}
}

func (g *generator) generateActions(routes []*RouteInfo, base string) {

	var controlled []*RouteInfo

	for _, r := range routes {
		if r.HasController() {
			controlled = append(controlled, r)
		}
	}

	byNS := groupByDotNamespace(controlled)

	for ns, nsRoutes := range byNS {
		path := dotNamespaceToPath(base, ns) + ".ts"
		g.appendCommonImports(path, ns, nsRoutes)
		g.writeControllerFile(path, ns, nsRoutes)
	}

	g.writeBarrelFiles(base, byNS)
}

func (g *generator) writeControllerFile(path, ns string, routes []*RouteInfo) {

	byMethod := groupByJsMethod(routes)

	for _, jsMethod := range sortedKeys(byMethod) {
		methodRoutes := byMethod[jsMethod]

		if len(methodRoutes) == 1 {

			shouldExport := !methodRoutes[0].IsInvokable
			g.writeMethodExport(path, methodRoutes[0], shouldExport)
		} else {
			g.writeMultiRouteExport(path, methodRoutes, true)
		}
	}

	var invokableRoute *RouteInfo

	var regularRoutes []*RouteInfo

	for _, r := range routes {
		if r.IsInvokable {
			invokableRoute = r
		} else {
			regularRoutes = append(regularRoutes, r)
		}
	}

	defaultExportName := lastDotSegment(ns)

	if invokableRoute != nil {

		var attachLines []string

		for _, r := range regularRoutes {
			attachLines = append(attachLines, fmt.Sprintf("%s.%s = %s",
				invokableRoute.JsMethod(), r.JsMethod(), r.JsMethod()))
		}

		if len(attachLines) > 0 {
			g.appendContent(path, strings.Join(attachLines, "\n"))
		}

		g.appendContent(path, fmt.Sprintf("\nexport default %s\n", invokableRoute.JsMethod()))
	} else {

		var parts []string
		seenMethods := map[string]bool{}

		for _, r := range routes {
			jsM := r.JsMethod()

			if !seenMethods[jsM] {
				parts = append(parts, jsM)
				seenMethods[jsM] = true
			}

			if r.OriginalJsMethod() != r.JsMethod() {
				orig := QuoteIfNeeded(r.OriginalJsMethod())
				alias := fmt.Sprintf("%s: %s", orig, r.JsMethod())

				found := false

				for _, p := range parts {
					if p == alias {
						found = true

						break
					}
				}

				if !found {
					parts = append(parts, alias)
				}
			}
		}

		g.appendContent(path, fmt.Sprintf("\nconst %s = { %s }\n\nexport default %s\n",
			defaultExportName,
			strings.Join(parts, ", "),
			defaultExportName,
		))
	}
}

func (g *generator) writeMethodExport(path string, r *RouteInfo, shouldExport bool) {
	method := r.JsMethod()
	g.appendContent(path, g.renderMethod(method, r, shouldExport))
}

func (g *generator) writeNamedMethodExport(path string, r *RouteInfo) {
	method := r.NamedMethod()
	g.appendContent(path, g.renderMethod(method, r, true))
}

func (g *generator) writeMultiRouteExport(path string, routes []*RouteInfo, shouldExport bool) {
	method := routes[0].JsMethod()

	var dictEntries []string

	for _, r := range routes {
		hash := shortHash(r.URI)
		tempMethod := method + hash

		tempRoute := *r
		tempRoute.IsInvokable = false
		g.appendContent(path, g.renderMethod(tempMethod, &tempRoute, false))

		dictEntries = append(dictEntries, fmt.Sprintf("    %s: %s", jsonString(r.URI), tempMethod))
	}

	exportKw := ""

	if shouldExport {
		exportKw = "export "
	}

	g.appendContent(path, fmt.Sprintf("\n%sconst %s = {\n%s\n}\n",
		exportKw,
		method,
		strings.Join(dictEntries, ",\n"),
	))
}

func (g *generator) renderMethod(method string, r *RouteInfo, shouldExport bool) string {
	verbs := r.Verbs()

	if len(verbs) == 0 {
		return ""
	}

	firstVerb := verbs[0].Actual

	params := r.Params
	hasParams := len(params) > 0
	args := argsVar(method)
	opts := optionsVar(method)
	parsed := parsedArgsVar(method)

	uri := r.FullURI()
	methods := r.methodActuals()

	var b strings.Builder

	docMethod := method

	if r.IsInvokable {
		docMethod = "__invoke"
	}

	cls := r.ControllerClass()

	if cls != "" && cls != "Closure" && cls != "\\Closure" {
		cls = strings.TrimPrefix(cls, "\\")
		b.WriteString(fmt.Sprintf("/**\n * @see %s::%s\n * @route %s\n */\n", cls, docMethod, uri))
	} else {
		b.WriteString(fmt.Sprintf("/**\n * @route %s\n */\n", uri))
	}

	exportKw := ""

	if shouldExport {
		exportKw = "export "
	}

	argsCall := ""

	if hasParams {
		argsCall = args + ", "
	}

	b.WriteString(fmt.Sprintf("%sconst %s = (%s): RouteDefinition<%s> => ({\n",
		exportKw, method,
		g.functionArgs(method, params, opts),
		jsonString(firstVerb),
	))
	b.WriteString(fmt.Sprintf("    url: %s.url(%s%s),\n", method, argsCall, opts))
	b.WriteString(fmt.Sprintf("    method: %s,\n", jsonString(firstVerb)))
	b.WriteString("})\n")

	b.WriteString(fmt.Sprintf("\n%s.definition = {\n", method))
	b.WriteString(fmt.Sprintf("    methods: %s,\n", jsonStringSlice(methods)))
	b.WriteString(fmt.Sprintf("    url: %s,\n", uri))
	b.WriteString(fmt.Sprintf("} satisfies RouteDefinition<%s>\n", jsonStringSlice(methods)))

	if cls != "" && cls != "Closure" && cls != "\\Closure" {
		b.WriteString(fmt.Sprintf("\n/**\n * @see %s::%s\n * @route %s\n */\n", cls, docMethod, uri))
	} else {
		b.WriteString(fmt.Sprintf("\n/**\n * @route %s\n */\n", uri))
	}

	b.WriteString(fmt.Sprintf("%s.url = (%s) => {\n", method, g.functionArgs(method, params, opts)))

	if hasParams {

		if len(params) == 1 {
			p := params[0]
			b.WriteString(fmt.Sprintf(
				"    if (typeof %s === 'string' || typeof %s === 'number') {\n        %s = { %s: %s }\n    }\n",
				args, args, args, p.Name, args,
			))

			if p.Key != "" {
				b.WriteString(fmt.Sprintf(
					"    if (typeof %s === 'object' && !Array.isArray(%s) && %s in %s) {\n        %s = { %s: %s.%s }\n    }\n",
					args, args, jsonString(p.Key), args, args, p.Name, args, p.Key,
				))
			}
		}

		b.WriteString(fmt.Sprintf("    if (Array.isArray(%s)) {\n        %s = {\n", args, args))

		for i, p := range params {
			b.WriteString(fmt.Sprintf("            %s: %s[%d],\n", p.Name, args, i))
		}

		b.WriteString("        }\n    }\n")

		b.WriteString(fmt.Sprintf("    %s = applyUrlDefaults(%s)\n", args, args))

		if hasOptional(params) {
			optNames := optionalNames(params)
			quoted := make([]string, len(optNames))

			for i, n := range optNames {
				quoted[i] = jsonString(n)
			}

			b.WriteString(fmt.Sprintf("    validateParameters(%s, [%s])\n",
				args, strings.Join(quoted, ", ")))
		}

		allOpt := allOptional(params)
		optMark := ""

		if allOpt {
			optMark = "?"
		}

		b.WriteString(fmt.Sprintf("    const %s = {\n", parsed))

		for _, p := range params {
			if p.Key != "" {

				defaultPart := ""

				if p.Default != "" {
					defaultPart = fmt.Sprintf(" ?? %s", jsonString(p.Default))
				}

				b.WriteString(fmt.Sprintf(
					"        %s: typeof %s%s.%s === 'object'\n            ? %s%s.%s.%s\n            : %s%s.%s%s,\n",
					p.Name,
					args, optMark, p.Name,
					args, optMark, p.Name, p.Key,
					args, optMark, p.Name, defaultPart,
				))
			} else {
				defaultPart := ""

				if p.Default != "" {
					defaultPart = fmt.Sprintf(" ?? %s", jsonString(p.Default))
				}

				b.WriteString(fmt.Sprintf("        %s: %s%s.%s%s,\n", p.Name, args, optMark, p.Name, defaultPart))
			}
		}

		b.WriteString("    }\n")

		b.WriteString(fmt.Sprintf("    return %s.definition.url\n", method))

		for i, p := range params {
			if p.Optional {
				b.WriteString(fmt.Sprintf(
					"        .replace(%s, %s.%s?.toString() ?? '')\n",
					jsonString(p.Placeholder()), parsed, p.Name,
				))
			} else {
				b.WriteString(fmt.Sprintf(
					"        .replace(%s, %s.%s.toString())\n",
					jsonString(p.Placeholder()), parsed, p.Name,
				))
			}

			if i == len(params)-1 {
				b.WriteString("        .replace(/\\/+$/, '')")
			}
		}

		b.WriteString(fmt.Sprintf(" + queryParams(%s)\n", opts))
	} else {

		b.WriteString(fmt.Sprintf("    return %s.definition.url\n", method))
		b.WriteString(fmt.Sprintf("        .replace(/\\/+$/, '') + queryParams(%s)\n", opts))
	}

	b.WriteString("}\n")

	for _, v := range verbs {
		if cls != "" && cls != "Closure" && cls != "\\Closure" {
			b.WriteString(fmt.Sprintf("\n/**\n * @see %s::%s\n * @route %s\n */\n", cls, docMethod, uri))
		} else {
			b.WriteString(fmt.Sprintf("\n/**\n * @route %s\n */\n", uri))
		}

		b.WriteString(fmt.Sprintf("%s.%s = (%s): RouteDefinition<%s> => ({\n",
			method, v.Actual,
			g.functionArgs(method, params, opts),
			jsonString(v.Actual),
		))
		b.WriteString(fmt.Sprintf("    url: %s.url(%s%s),\n", method, argsCall, opts))
		b.WriteString(fmt.Sprintf("    method: %s,\n", jsonString(v.Actual)))
		b.WriteString("})\n")
	}

	if g.opts.WithForm {
		b.WriteString(g.renderFormHelpers(method, r, params, verbs))
	}

	return b.String()
}

func (g *generator) renderFormHelpers(method string, r *RouteInfo, params []Param, verbs []Verb) string {
	if len(verbs) == 0 {
		return ""
	}

	firstVerb := verbs[0]
	opts := optionsVar(method)
	args := argsVar(method)
	argsCall := ""

	if len(params) > 0 {
		argsCall = args + ", "
	}

	var b strings.Builder
	formMethod := method + "Form"

	cls := r.ControllerClass()
	docMethod := method

	if r.IsInvokable {
		docMethod = "__invoke"
	}

	uri := r.FullURI()

	docblock := func() {
		if cls != "" && cls != "Closure" && cls != "\\Closure" {
			b.WriteString(fmt.Sprintf("\n/**\n * @see %s::%s\n * @route %s\n */\n",
				strings.TrimPrefix(cls, "\\"), docMethod, uri))
		} else {
			b.WriteString(fmt.Sprintf("\n/**\n * @route %s\n */\n", uri))
		}
	}

	formArg := func(v Verb) string {
		if v.FormSafe == v.Actual {
			return opts
		}

		return fmt.Sprintf("{\n    [%s?.mergeQuery ? 'mergeQuery' : 'query']: {\n        _method: %s,\n        ...(%s?.query ?? %s?.mergeQuery ?? {}),\n    }\n}",
			opts, jsonString(strings.ToUpper(v.Actual)), opts, opts)
	}

	docblock()
	b.WriteString(fmt.Sprintf("const %s = (%s): RouteFormDefinition<%s> => ({\n",
		formMethod,
		g.functionArgs(method, params, opts),
		jsonString(firstVerb.FormSafe),
	))
	b.WriteString(fmt.Sprintf("    action: %s.url(%s%s),\n", method, argsCall, formArg(firstVerb)))
	b.WriteString(fmt.Sprintf("    method: %s,\n", jsonString(firstVerb.FormSafe)))
	b.WriteString("})\n")

	for _, v := range verbs {
		docblock()
		b.WriteString(fmt.Sprintf("%s.%s = (%s): RouteFormDefinition<%s> => ({\n",
			formMethod, v.Actual,
			g.functionArgs(method, params, opts),
			jsonString(v.FormSafe),
		))
		b.WriteString(fmt.Sprintf("    action: %s.url(%s%s),\n", method, argsCall, formArg(v)))
		b.WriteString(fmt.Sprintf("    method: %s,\n", jsonString(v.FormSafe)))
		b.WriteString("})\n")
	}

	b.WriteString(fmt.Sprintf("\n%s.form = %s\n", method, formMethod))

	return b.String()
}

func (g *generator) functionArgs(method string, params []Param, opts string) string {
	if len(params) == 0 {
		return opts + "?: RouteQueryOptions"
	}

	args := argsVar(method)
	allOpt := allOptional(params)
	optMark := ""

	if allOpt {
		optMark = "?"
	}

	var typeLines []string

	var objFields []string

	for _, p := range params {
		fieldOpt := ""

		if p.Optional {
			fieldOpt = "?"
		}

		line := fmt.Sprintf("    %s%s: %s", p.Name, fieldOpt, p.TSTypes())

		if p.Key != "" {
			line += fmt.Sprintf(" | { %s: %s }", p.Key, p.TSTypes())
		}

		objFields = append(objFields, line)
	}

	typeLines = append(typeLines, "{\n"+strings.Join(objFields, ",\n")+"\n}")

	var arrFields []string

	for _, p := range params {
		arrFields = append(arrFields, fmt.Sprintf("%s: %s", p.SafeName(), p.TSTypes()))

		if p.Key != "" {
			arrFields[len(arrFields)-1] += fmt.Sprintf(" | { %s: %s }", p.Key, p.TSTypes())
		}
	}

	typeLines = append(typeLines, "[\n"+strings.Join(arrFields, ", ")+"\n]")

	if len(params) == 1 {
		p := params[0]
		prim := p.TSTypes()
		typeLines = append(typeLines, prim)

		if p.Key != "" {
			typeLines = append(typeLines, fmt.Sprintf("{ %s: %s }", p.Key, p.TSTypes()))
		}
	}

	unionType := strings.Join(typeLines, "\n| ")

	return fmt.Sprintf("%s%s: %s,\n%s?: RouteQueryOptions", args, optMark, unionType, opts)
}

func (g *generator) writeBarrelFiles(base string, byNS map[string][]*RouteInfo) {

	tree := buildNamespaceTree(byNS)
	g.writeBarrelNode(base, "", tree)
}

func buildNamespaceTree(byNS map[string][]*RouteInfo) nsNode {
	root := nsNode{}

	for ns := range byNS {
		parts := strings.Split(ns, ".")
		node := root

		for _, p := range parts {
			if _, ok := node[p]; !ok {
				node[p] = nsNode{}
			}

			if sub, ok := node[p].(nsNode); ok {
				node = sub
			}
		}
	}

	return root
}

func (g *generator) writeBarrelNode(base, parent string, node nsNode) {
	if len(node) == 0 {
		return
	}

	var dir string

	if parent == "" {
		dir = base
	} else {
		dir = filepath.Join(base, filepath.FromSlash(strings.ReplaceAll(parent, ".", "/")))
	}

	indexPath := filepath.Join(dir, "index.ts")

	children := sortedKeys2(node)

	var importLines []string

	var objectEntries []string

	for _, child := range children {
		safeName := SafeMethod(child, "Method")
		importLines = append(importLines, fmt.Sprintf("import %s from './%s'", safeName, child))
		normalised := toNormalisedKey(child)

		if normalised != safeName {
			objectEntries = append(objectEntries, fmt.Sprintf("    %s: %s", normalised, safeName))
		} else {
			objectEntries = append(objectEntries, fmt.Sprintf("    %s", safeName))
		}
	}

	varName := SafeMethod(lastPathSegment(parent, base), "Method")

	if varName == "" {
		varName = "actions"
	}

	var buf strings.Builder

	buf.WriteString(strings.Join(importLines, "\n"))

	if len(importLines) > 0 {
		buf.WriteString("\n\n")
	}

	buf.WriteString(fmt.Sprintf("const %s = {\n%s,\n}\n\nexport default %s\n",
		varName,
		strings.Join(objectEntries, ",\n"),
		varName,
	))

	g.prependContent(indexPath, buf.String())

	for _, child := range children {
		if sub, ok := node[child].(nsNode); ok {
			childParent := child

			if parent != "" {
				childParent = parent + "." + child
			}

			g.writeBarrelNode(base, childParent, sub)
		}
	}
}

func (g *generator) generateRoutes(routes []*RouteInfo, base string) {

	byName := make(map[string]*RouteInfo)

	for _, r := range routes {
		name := cleanRouteName(r.Name)

		if name == "" {
			continue
		}

		byName[name] = r
	}

	for name, r := range byName {
		parts := routeNameToFileParts(name)
		path := filepath.Join(append([]string{base}, parts...)...) + ".ts"
		g.appendCommonImports(path, name, []*RouteInfo{r})
		g.writeNamedMethodExport(path, r)
	}

	nameTree := buildNameTree(byName)
	g.writeNameBarrelNode(base, "", nameTree, byName)
}

func buildNameTree(byName map[string]*RouteInfo) nameNode {
	root := nameNode{}

	for name := range byName {
		parts := strings.Split(name, ".")
		node := root

		for _, p := range parts[:len(parts)-1] { // all but the leaf
			if _, ok := node[p]; !ok {
				node[p] = nameNode{}
			}

			if sub, ok := node[p].(nameNode); ok {
				node = sub
			}
		}
	}

	return root
}

// writeNameBarrelNode writes barrel index.ts files for routes/ directory.
func (g *generator) writeNameBarrelNode(base, prefix string, node nameNode, byName map[string]*RouteInfo) {
	children := sortedKeys3(node)

	if len(children) == 0 {
		return
	}

	var dir string

	if prefix == "" {
		dir = base
	} else {
		dir = filepath.Join(base, filepath.FromSlash(strings.ReplaceAll(prefix, ".", "/")))
	}

	indexPath := filepath.Join(dir, "index.ts")

	var importLines []string

	var objectEntries []string

	for _, child := range children {
		safeName := SafeMethod(child, "Method")
		importLines = append(importLines, fmt.Sprintf("import %s from './%s'", safeName, child))
		normalised := toNormalisedKey(child)

		if normalised != safeName {
			objectEntries = append(objectEntries, fmt.Sprintf("    %s: %s", normalised, safeName))
		} else {
			objectEntries = append(objectEntries, fmt.Sprintf("    %s", safeName))
		}
	}

	varName := SafeMethod(lastPathSegment(prefix, base), "Method")

	if varName == "" {
		varName = "routes"
	}

	var buf strings.Builder

	buf.WriteString(strings.Join(importLines, "\n"))

	if len(importLines) > 0 {
		buf.WriteString("\n\n")
	}

	buf.WriteString(fmt.Sprintf("const %s = {\n%s,\n}\n\nexport default %s\n",
		varName,
		strings.Join(objectEntries, ",\n"),
		varName,
	))

	g.prependContent(indexPath, buf.String())

	for _, child := range children {
		if sub, ok := node[child].(nameNode); ok {
			childPrefix := child

			if prefix != "" {
				childPrefix = prefix + "." + child
			}

			g.writeNameBarrelNode(base, childPrefix, sub, byName)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Import management
// ─────────────────────────────────────────────────────────────────────────────

// appendCommonImports adds the routegen runtime imports to a file's import map.
func (g *generator) appendCommonImports(path, namespace string, routes []*RouteInfo) {
	imports := []string{"queryParams", "type RouteQueryOptions", "type RouteDefinition"}

	if g.opts.WithForm {
		imports = append(imports, "type RouteFormDefinition")
	}

	// applyUrlDefaults needed when any route has params.
	for _, r := range routes {
		if len(r.Params) > 0 {
			imports = append(imports, "applyUrlDefaults")

			break
		}
	}

	// validateParameters needed when any route has optional params.
outer:
	for _, r := range routes {
		for _, p := range r.Params {
			if p.Optional {
				imports = append(imports, "validateParameters")

				break outer
			}
		}
	}

	// Compute the relative import path depth.
	dotCount := strings.Count(namespace, ".")
	ups := strings.Repeat("/..", dotCount+1)
	importFrom := "." + ups + "/routegen"

	if g.imports[path] == nil {
		g.imports[path] = make(map[string][]string)
	}

	g.imports[path][importFrom] = append(g.imports[path][importFrom], imports...)
}

// ─────────────────────────────────────────────────────────────────────────────
// Content management
// ─────────────────────────────────────────────────────────────────────────────

func (g *generator) appendContent(path, content string) {
	g.content[path] = append(g.content[path], content)
}

func (g *generator) prependContent(path, content string) {
	g.content[path] = append([]string{content}, g.content[path]...)
}

// flush writes all accumulated content to disk, prepending imports.
func (g *generator) flush(base string) error {
	for path, chunks := range g.content {
		var b strings.Builder

		// Prepend imports.
		if imp, ok := g.imports[path]; ok {
			// De-duplicate and sort import identifiers.
			for from, ids := range imp {
				unique := deduplicateStrings(ids)

				sort.Strings(unique)

				b.WriteString(fmt.Sprintf("import { %s } from '%s'\n", strings.Join(unique, ", "), from))
			}

			b.WriteByte('\n')
		}

		b.WriteString(strings.Join(chunks, "\n"))

		cleaned := CleanUp(b.String())

		if err := writeFile(path, cleaned); err != nil {
			return err
		}
	}

	// Reset state for next flush call.
	g.content = make(map[string][]string)
	g.imports = make(map[string]map[string][]string)

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// groupByDotNamespace groups routes by their dot-separated controller namespace.
func groupByDotNamespace(routes []*RouteInfo) map[string][]*RouteInfo {
	m := make(map[string][]*RouteInfo)

	for _, r := range routes {
		ns := r.DotNamespace()
		m[ns] = append(m[ns], r)
	}

	return m
}

// groupByJsMethod groups routes by their sanitised JS method name.
func groupByJsMethod(routes []*RouteInfo) map[string][]*RouteInfo {
	m := make(map[string][]*RouteInfo)

	for _, r := range routes {
		js := r.JsMethod()
		m[js] = append(m[js], r)
	}

	return m
}

// dotNamespaceToPath converts a dot-separated namespace to an absolute file path.
// "App.Http.Controllers.PostController" → "{base}/App/Http/Controllers/PostController"
func dotNamespaceToPath(base, ns string) string {
	parts := strings.Split(ns, ".")

	return filepath.Join(append([]string{base}, parts...)...)
}

// lastDotSegment returns the portion after the last dot.
// "App.Http.Controllers.PostController" → "PostController"
func lastDotSegment(s string) string {
	if idx := strings.LastIndex(s, "."); idx >= 0 {
		return s[idx+1:]
	}

	return s
}

// lastPathSegment returns the last component of a dot-path, or the base
// directory name when parent is empty.
func lastPathSegment(parent, base string) string {
	if parent == "" {
		return filepath.Base(base)
	}

	return lastDotSegment(parent)
}

// toNormalisedKey converts a child key to a normalised object property name.
// Kebab-case names are converted to camelCase; others are unchanged.
func toNormalisedKey(child string) string {
	if strings.Contains(child, "-") {
		return toCamel(child)
	}

	return child
}

// cleanRouteName strips invalid route name forms (generated::, trailing dot).
func cleanRouteName(name string) string {
	if name == "" {
		return ""
	}

	if strings.HasPrefix(name, "generated::") {
		return ""
	}

	if strings.HasSuffix(name, ".") {
		return ""
	}

	if strings.Contains(name, "::") {
		name = "namespaced." + strings.ReplaceAll(name, "::", ".")
	}

	return name
}

// shortHash returns an 8-character hash of s for use as a unique suffix.
func shortHash(s string) string {
	h := md5.Sum([]byte(s))

	return fmt.Sprintf("%x", h[:4])
}

// deduplicateStrings returns a new slice with duplicate strings removed,
// preserving the first occurrence order.
func deduplicateStrings(ss []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(ss))

	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}

	return out
}

func sortedKeys(m map[string][]*RouteInfo) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func sortedKeys2(m nsNode) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func sortedKeys3(m nameNode) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
