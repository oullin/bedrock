package prompts

import "strings"

// AutocompleteOption configures an AutocompletePrompt.
type AutocompleteOption func(*AutocompletePrompt)

// AutocompleteWithPlaceholder sets placeholder text.

// AutocompleteWithDefault sets the default value.

// AutocompleteWithRequired makes the prompt required.

// AutocompleteWithValidate sets validation.

// AutocompleteWithHint sets hint text.

// AutocompleteWithTransform sets the transform function.

// AutocompletePrompt handles text input with ghost-text auto-completion.
type AutocompletePrompt struct {
	prompt      *Prompt
	TypedValue  TypedValue
	label       string
	placeholder string
	hint        string
	optionsFn   func(string) []string
	staticOpts  []string
}

func AutocompleteWithPlaceholder(s string) AutocompleteOption {
	return func(p *AutocompletePrompt) { p.placeholder = s }
}

func AutocompleteWithDefault(s string) AutocompleteOption {
	return func(p *AutocompletePrompt) { p.TypedValue.SetValue(s) }
}

func AutocompleteWithRequired(v any) AutocompleteOption {
	return func(p *AutocompletePrompt) { p.prompt.required = v }
}

func AutocompleteWithValidate(fn ValidateFunc) AutocompleteOption {
	return func(p *AutocompletePrompt) { p.prompt.validate = fn }
}

func AutocompleteWithHint(s string) AutocompleteOption {
	return func(p *AutocompletePrompt) { p.hint = s }
}

func AutocompleteWithTransform(fn TransformFunc) AutocompleteOption {
	return func(p *AutocompletePrompt) { p.prompt.transform = fn }
}

// Autocomplete displays a text input with ghost-text auto-completion.
// options accepts []string or func(string) []string.
func Autocomplete(label string, options any, opts ...AutocompleteOption) (string, error) {
	p := &AutocompletePrompt{
		prompt: newPrompt(),
		label:  label,
	}

	switch v := options.(type) {
	case []string:
		p.staticOpts = v
	case func(string) []string:
		p.optionsFn = v
	}

	for _, opt := range opts {
		opt(p)
	}

	p.prompt.valueFn = p.TypedValue.Value
	p.prompt.renderer = func(state State) string {
		return getTheme().AutocompleteRenderer(p, state)
	}

	// Handle Tab to accept ghost text.
	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		if key == KeyTab || IsRightKey(key) {
			ghost := p.GhostText()

			if ghost != "" {
				p.TypedValue.SetValue(p.TypedValue.Value() + ghost)
			}
		}
	})

	p.TypedValue.TrackTypedValue(p.prompt, p.prompt.submit)

	return p.prompt.run()
}

// Matches returns the current matching options.
func (p *AutocompletePrompt) Matches() []string {
	input := p.TypedValue.Value()

	if p.optionsFn != nil {
		return p.optionsFn(input)
	}

	if input == "" {
		return p.staticOpts
	}

	lower := strings.ToLower(input)

	var matches []string

	for _, opt := range p.staticOpts {
		if strings.HasPrefix(strings.ToLower(opt), lower) {
			matches = append(matches, opt)
		}
	}

	return matches
}

// GhostText returns the completion suffix for the current best match.
func (p *AutocompletePrompt) GhostText() string {
	matches := p.Matches()

	if len(matches) == 0 {
		return ""
	}

	val := p.TypedValue.Value()
	best := matches[0]

	if len(best) > len(val) && strings.HasPrefix(strings.ToLower(best), strings.ToLower(val)) {
		return best[len(val):]
	}

	return ""
}

// Value returns the current typed value.
func (p *AutocompletePrompt) Value() string { return p.TypedValue.Value() }
