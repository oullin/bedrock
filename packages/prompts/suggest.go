package prompts

import "strings"

// SuggestOption configures a SuggestPrompt.
type SuggestOption func(*SuggestPrompt)

// SuggestWithPlaceholder sets placeholder text.

// SuggestWithDefault sets the default value.

// SuggestWithScroll sets visible suggestion count.

// SuggestWithRequired makes the prompt required.

// SuggestWithValidate sets validation.

// SuggestWithHint sets hint text.

// SuggestWithTransform sets the transform function.

// SuggestPrompt handles text input with auto-complete suggestions.
type SuggestPrompt struct {
	prompt      *Prompt
	TypedValue  TypedValue
	Scrollable  Scrollable
	label       string
	placeholder string
	hint        string
	scroll      int
	optionsFn   func(string) []string
	staticOpts  []string
}

func SuggestWithPlaceholder(s string) SuggestOption {
	return func(p *SuggestPrompt) { p.placeholder = s }
}

func SuggestWithDefault(s string) SuggestOption {
	return func(p *SuggestPrompt) { p.TypedValue.SetValue(s) }
}

func SuggestWithScroll(n int) SuggestOption { return func(p *SuggestPrompt) { p.scroll = n } }

func SuggestWithRequired(v any) SuggestOption {
	return func(p *SuggestPrompt) { p.prompt.required = v }
}

func SuggestWithValidate(fn ValidateFunc) SuggestOption {
	return func(p *SuggestPrompt) { p.prompt.validate = fn }
}

func SuggestWithHint(s string) SuggestOption { return func(p *SuggestPrompt) { p.hint = s } }

func SuggestWithTransform(fn TransformFunc) SuggestOption {
	return func(p *SuggestPrompt) { p.prompt.transform = fn }
}

// Suggest displays a text input with auto-complete suggestions.
// options accepts []string or func(string) []string.
func Suggest(label string, options any, opts ...SuggestOption) (string, error) {
	p := &SuggestPrompt{
		prompt: newPrompt(),
		label:  label,
		scroll: 5,
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
		return getTheme().SuggestRenderer(p, state)
	}

	// Handle highlight selection.
	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		matches := p.Matches()
		p.Scrollable.UpdateTotal(len(matches))

		switch {
		case key == KeyEnter:
			// If a suggestion is highlighted, use it.
			if len(matches) > 0 && p.Scrollable.Highlighted() < len(matches) {
				// Check if the user was navigating suggestions.
				selected := matches[p.Scrollable.Highlighted()]

				if p.TypedValue.Value() != selected {
					p.TypedValue.SetValue(selected)

					return
				}
			}

			p.prompt.submit()
		case key == KeyTab && len(matches) > 0:
			p.TypedValue.SetValue(matches[p.Scrollable.Highlighted()])
		case IsUpKey(key) && len(matches) > 0:
			p.Scrollable.HighlightPrevious(1)
		case IsDownKey(key) && len(matches) > 0:
			p.Scrollable.HighlightNext(1)
		case OneOfKey(KeyHome, key) && len(matches) > 0:
			p.Scrollable.HighlightFirst()
		case OneOfKey(KeyEnd, key) && len(matches) > 0:
			p.Scrollable.HighlightLast()
		}
	})

	p.TypedValue.TrackTypedValue(p.prompt, p.prompt.submit)

	return p.prompt.run()
}

// Matches returns the current matching suggestions.
func (p *SuggestPrompt) Matches() []string {
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
		if strings.Contains(strings.ToLower(opt), lower) {
			matches = append(matches, opt)
		}
	}

	return matches
}

// Value returns the current typed value.
func (p *SuggestPrompt) Value() string { return p.TypedValue.Value() }
