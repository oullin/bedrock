package inertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/bedrock/packages/inertia/middleware"
	"github.com/bedrock/packages/inertia/props"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/response"
)

// Inertia is the central server-side adapter for the Inertia.js
// protocol. It holds the root HTML template, asset version, shared
// props, and configuration. It is safe for concurrent use.
type Inertia struct {
	rootTemplate   *template.Template
	version        string
	containerID    string
	encryptHistory bool
	sharedProps    protocol.Props
	templateFuncs  template.FuncMap
	jsonMarshaler  protocol.JSONMarshaler
	logger         protocol.Logger
	head           protocol.Head
	headSource     headSource
	mu             sync.RWMutex
}

// New creates an Inertia instance by parsing rootTemplateHTML as a
// Go html/template.

// NewFromFile creates an Inertia instance from a template file at path.

// NewFromReader creates an Inertia instance from a template read from r.

// NewFromTemplate creates an Inertia instance from a pre-parsed template.

// Render builds an Inertia response for the given component. On XHR
// visits (X-Inertia: true) it writes JSON; on initial visits it
// renders the root HTML template with the page data embedded.

// Middleware returns the Inertia HTTP middleware configured with the
// current asset version. It handles version checking, the Vary header,
// and 302 → 303 redirect conversion.

// Redirect sends an HTTP redirect. The default status is 302 Found.

// Back redirects to the Referer URL, falling back to "/" if absent.

// Location performs an external redirect. For Inertia requests it
// responds with 409 Conflict and sets X-Inertia-Location, which tells
// the client to do a full page visit. For non-Inertia requests it
// performs a standard redirect.

// ShareProp registers a global prop included in every response.

// ShareProps registers multiple global props.

// SharedProps returns a copy of all currently registered shared props.

// Version returns the current asset version string.

// mergeProps combines shared props, context props, validation errors,
// and the props passed to Render. Later sources override earlier ones.

// Context props (set via SetProp/SetProps in middleware).

// Render-time props (the highest priority).

// Validation errors.

// Reparse template if custom template funcs were provided via options.
// This is needed because template funcs must be registered before parsing.

// StdJSONMarshaler wraps encoding/json as the default JSONMarshaler.
type StdJSONMarshaler struct{}

func New(rootTemplateHTML string, opts ...Option) (*Inertia, error) {
	i := defaults()
	tmpl, err := template.New("inertia").Funcs(i.templateFuncs).Parse(rootTemplateHTML)

	if err != nil {
		return nil, fmt.Errorf("inertia: parse template: %w", err)
	}

	i.rootTemplate = tmpl

	return i, i.apply(opts)
}

func NewFromFile(path string, opts ...Option) (*Inertia, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("inertia: read template: %w", err)
	}

	return New(string(data), opts...)
}

func NewFromReader(r io.Reader, opts ...Option) (*Inertia, error) {
	data, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("inertia: read template: %w", err)
	}

	return New(string(data), opts...)
}

func NewFromTemplate(t *template.Template, opts ...Option) (*Inertia, error) {
	i := defaults()
	i.rootTemplate = t

	return i, i.apply(opts)
}

func (i *Inertia) Render(w http.ResponseWriter, r *http.Request, component string, pageProps ...protocol.Props) error {
	if handled, err := i.HandlePrecognition(w, r); handled {
		return err
	}

	ctx := r.Context()
	merged := i.mergeProps(r, pageProps...)

	result, err := props.Resolve(r, component, merged)

	if err != nil {
		return fmt.Errorf("inertia: resolve props: %w", err)
	}

	page := &response.Page{
		Component:      component,
		Props:          result.Props,
		URL:            r.RequestURI,
		Version:        i.version,
		EncryptHistory: i.encryptHistory || encryptHistoryFromContext(r.Context()),
		ClearHistory:   clearHistoryFromContext(r.Context()),
		MergeProps:     result.MergeProps,
		DeepMergeProps: result.DeepMergeProps,
		DeferredProps:  result.DeferredProps,
		ScrollProps:    result.ScrollProps,
		OnceProps:      result.OnceProps,
	}

	if protocol.IsInertiaRequest(r) {
		return response.WriteJSON(w, page, i.jsonMarshaler)
	}

	// Build the final head by merging three layers:
	// 1. Global defaults (from WithHead / WithHeadFromFile)
	// 2. Locale head (from i18n middleware via SetLocale)
	// 3. Per-request context head (from SetHead / SetTitle / SetMeta)
	i.mu.RLock()

	finalHead := i.head

	i.mu.RUnlock()

	if locale := protocol.LocaleFromContext(ctx); locale != nil {
		finalHead = protocol.MergeHead(finalHead, locale.Head)

		if strings.TrimSpace(finalHead.Lang) == "" {
			finalHead.Lang = locale.Code
		}

		if strings.TrimSpace(finalHead.Direction) == "" {
			finalHead.Direction = locale.Direction
		}
	}

	finalHead = protocol.MergeHead(finalHead, headFromContext(ctx))

	// Auto-append CSRF token as a meta tag when present in context.
	if token := protocol.CSRFTokenFromContext(ctx); strings.TrimSpace(token) != "" {
		finalHead.Meta = append(finalHead.Meta, protocol.MetaTag{
			Name:    "csrf-token",
			Content: token,
		})
	}

	return response.WriteHTML(w, page, response.HTMLConfig{
		Template:    i.rootTemplate,
		ContainerID: i.containerID,
		Marshaler:   i.jsonMarshaler,
		ExtraData:   templateDataFromContext(ctx),
		Head:        finalHead,
	})
}

// HandlePrecognition writes the validation-only response for a precognitive
// request. It returns handled=true when the request was precognitive.
func (i *Inertia) HandlePrecognition(w http.ResponseWriter, r *http.Request) (handled bool, err error) {
	if !protocol.IsPrecognition(r.Context()) {
		return false, nil
	}

	return true, i.writePrecognitionResponse(w, r)
}

func (i *Inertia) writePrecognitionResponse(w http.ResponseWriter, r *http.Request) error {
	errors := validationErrorsFromContext(r.Context())

	w.Header().Set(protocol.HeaderPrecognition, "true")

	// Filter errors to only requested fields when Validate-Only is present.
	if only := protocol.ValidateOnly(r); len(only) > 0 && len(errors) > 0 {
		filtered := make(protocol.ValidationErrors, len(only))

		for _, field := range only {
			if v, ok := errors[field]; ok {
				filtered[field] = v
			}
		}

		errors = filtered
	}

	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusUnprocessableEntity)

		data, err := i.jsonMarshaler.Marshal(map[string]any{"errors": errors})

		if err != nil {
			return fmt.Errorf("inertia: marshal precognition errors: %w", err)
		}

		_, err = w.Write(data)

		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (i *Inertia) Middleware(next http.Handler) http.Handler {
	i.mu.RLock()

	version := i.version

	i.mu.RUnlock()

	return middleware.New(middleware.Config{
		Version: version,
	})(next)
}

func (i *Inertia) Redirect(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	code := http.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	http.Redirect(w, r, url, code)
}

func (i *Inertia) Back(w http.ResponseWriter, r *http.Request, status ...int) {
	url := r.Header.Get("Referer")

	if strings.TrimSpace(url) == "" {
		url = "/"
	}

	i.Redirect(w, r, url, status...)
}

func (i *Inertia) Location(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	if protocol.IsInertiaRequest(r) {
		w.Header().Set(protocol.HeaderLocation, url)

		w.WriteHeader(http.StatusConflict)

		return
	}

	i.Redirect(w, r, url, status...)
}

func (i *Inertia) ShareProp(key string, val any) {
	i.mu.Lock()

	i.sharedProps[key] = val

	i.mu.Unlock()
}

func (i *Inertia) ShareProps(p protocol.Props) {
	i.mu.Lock()

	for k, v := range p {
		i.sharedProps[k] = v
	}

	i.mu.Unlock()
}

func (i *Inertia) SharedProps() protocol.Props {
	i.mu.RLock()

	defer i.mu.RUnlock()

	out := make(protocol.Props, len(i.sharedProps))

	for k, v := range i.sharedProps {
		out[k] = v
	}

	return out
}

func (i *Inertia) Version() string {
	i.mu.RLock()

	defer i.mu.RUnlock()

	return i.version
}

func (i *Inertia) mergeProps(r *http.Request, pageProps ...protocol.Props) protocol.Props {
	i.mu.RLock()

	shared := make(protocol.Props, len(i.sharedProps))

	for k, v := range i.sharedProps {
		shared[k] = v
	}

	i.mu.RUnlock()

	ctx := r.Context()

	sources := make([]protocol.Props, 0, 2+len(pageProps)+1)
	sources = append(sources, shared, propsFromContext(ctx))
	sources = append(sources, pageProps...)

	if errors := validationErrorsFromContext(ctx); len(errors) > 0 {
		sources = append(sources, protocol.Props{"errors": errors})
	}

	return props.MergeAll(sources...)
}

func defaults() *Inertia {
	return &Inertia{
		containerID:   "app",
		sharedProps:   make(protocol.Props),
		jsonMarshaler: &StdJSONMarshaler{},
	}
}

func (i *Inertia) apply(opts []Option) error {
	for _, opt := range opts {
		if err := opt(i); err != nil {
			return err
		}
	}

	if i.templateFuncs != nil && i.rootTemplate != nil {
		i.rootTemplate = i.rootTemplate.Funcs(i.templateFuncs)
	}

	return nil
}

func (m *StdJSONMarshaler) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (m *StdJSONMarshaler) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
