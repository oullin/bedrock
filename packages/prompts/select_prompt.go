package prompts

// SelectOption configures a SelectPrompt.
type SelectOption func(*SelectPrompt)

// SelectWithDefault sets the default selected value.

// SelectWithScroll sets the number of visible options.

// SelectWithValidate sets validation.

// SelectWithHint sets hint text.

// SelectWithRequired makes the prompt required.

// SelectWithTransform sets the transform function.

// SelectPrompt handles single-selection from a list of options.
type SelectPrompt struct {
	prompt       *Prompt
	Scrollable   Scrollable
	label        string
	options      []OptionItem
	defaultValue string
	hint         string
	scroll       int
}

func SelectWithDefault(s string) SelectOption {
	return func(p *SelectPrompt) { p.defaultValue = s }
}

func SelectWithScroll(n int) SelectOption { return func(p *SelectPrompt) { p.scroll = n } }

func SelectWithValidate(fn ValidateFunc) SelectOption {
	return func(p *SelectPrompt) { p.prompt.validate = fn }
}

func SelectWithHint(s string) SelectOption { return func(p *SelectPrompt) { p.hint = s } }

func SelectWithRequired(v any) SelectOption {
	return func(p *SelectPrompt) { p.prompt.required = v }
}

func SelectWithTransform(fn TransformFunc) SelectOption {
	return func(p *SelectPrompt) { p.prompt.transform = fn }
}

// Select displays a single-selection prompt.
// options accepts []string, map[string]string, or []OptionItem.
func Select(label string, options any, opts ...SelectOption) (string, error) {
	items, err := resolveOptions(options)

	if err != nil {
		return "", err
	}

	p := &SelectPrompt{
		prompt:  newPrompt(),
		label:   label,
		options: items,
		scroll:  5,
	}

	for _, opt := range opts {
		opt(p)
	}

	// Find default index and set on base prompt for non-interactive.
	defaultIdx := 0

	if p.defaultValue != "" {
		p.prompt.defaultValue = p.defaultValue

		for i, item := range items {
			if item.Key == p.defaultValue {
				defaultIdx = i

				break
			}
		}
	} else if len(items) > 0 {
		p.prompt.defaultValue = items[0].Key
	}

	p.Scrollable.InitScrolling(len(items), p.scroll)
	p.Scrollable.SetHighlighted(defaultIdx)

	p.prompt.valueFn = func() string {
		if len(p.options) == 0 {
			return ""
		}

		return p.options[p.Scrollable.Highlighted()].Key
	}

	p.prompt.renderer = func(state State) string {
		return getTheme().SelectRenderer(p, state)
	}

	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			p.prompt.submit()
		case IsUpKey(key):
			p.Scrollable.HighlightPrevious(1)
		case IsDownKey(key):
			p.Scrollable.HighlightNext(1)
		case OneOfKey(KeyHome, key):
			p.Scrollable.HighlightFirst()
		case OneOfKey(KeyEnd, key):
			p.Scrollable.HighlightLast()
		case key == KeyPageUp:
			p.Scrollable.HighlightPrevious(p.scroll)
		case key == KeyPageDown:
			p.Scrollable.HighlightNext(p.scroll)
		}
	})

	return p.prompt.run()
}

func (p *SelectPrompt) selectedLabel() string {
	if len(p.options) == 0 {
		return ""
	}

	return p.options[p.Scrollable.Highlighted()].Label
}
