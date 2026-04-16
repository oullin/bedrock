package prompts

// MultiSelectOption configures a MultiSelectPrompt.
type MultiSelectOption func(*MultiSelectPrompt)

// MultiSelectWithDefault sets the default selected values.

// MultiSelectWithScroll sets the number of visible options.

// MultiSelectWithRequired makes the prompt required.

// MultiSelectWithValidate sets validation.

// MultiSelectWithHint sets hint text.

// MultiSelectWithTransform sets the transform function.

// MultiSelectPrompt handles multiple selection from a list of options.
type MultiSelectPrompt struct {
	prompt        *Prompt
	Scrollable    Scrollable
	label         string
	options       []OptionItem
	selected      map[string]bool
	defaultValues []string
	hint          string
	scroll        int
	validateMulti func([]string) string
}

func MultiSelectWithDefault(v []string) MultiSelectOption {
	return func(p *MultiSelectPrompt) { p.defaultValues = v }
}

func MultiSelectWithScroll(n int) MultiSelectOption {
	return func(p *MultiSelectPrompt) { p.scroll = n }
}

func MultiSelectWithRequired(v any) MultiSelectOption {
	return func(p *MultiSelectPrompt) { p.prompt.required = v }
}

func MultiSelectWithValidate(fn func([]string) string) MultiSelectOption {
	return func(p *MultiSelectPrompt) {
		p.validateMulti = fn
	}
}

func MultiSelectWithHint(s string) MultiSelectOption {
	return func(p *MultiSelectPrompt) { p.hint = s }
}

func MultiSelectWithTransform(fn TransformFunc) MultiSelectOption {
	return func(p *MultiSelectPrompt) { p.prompt.transform = fn }
}

// MultiSelect displays a multi-selection prompt.
// options accepts []string, map[string]string, or []OptionItem.
func MultiSelect(label string, options any, opts ...MultiSelectOption) ([]string, error) {
	items, err := resolveOptions(options)

	if err != nil {
		return nil, err
	}

	p := &MultiSelectPrompt{
		prompt:   newPrompt(),
		label:    label,
		options:  items,
		selected: make(map[string]bool),
		scroll:   5,
	}

	for _, opt := range opts {
		opt(p)
	}

	// Apply defaults.
	for _, d := range p.defaultValues {
		p.selected[d] = true
	}

	p.Scrollable.InitScrolling(len(items), p.scroll)

	p.prompt.valueFn = func() string { return "" } // unused, we handle values ourselves
	p.prompt.renderer = func(state State) string {
		return getTheme().MultiSelectRenderer(p, state)
	}

	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			vals := p.SelectedValues()
			// Check required.
			if p.prompt.required != nil {
				switch r := p.prompt.required.(type) {
				case bool:
					if r && len(vals) == 0 {
						p.prompt.state = StateError
						p.prompt.errorMsg = "Required."

						return
					}
				case string:
					if len(vals) == 0 {
						p.prompt.state = StateError

						if r != "" {
							p.prompt.errorMsg = r
						} else {
							p.prompt.errorMsg = "Required."
						}

						return
					}
				}
			}

			if p.validateMulti != nil {
				if msg := p.validateMulti(vals); msg != "" {
					p.prompt.state = StateError
					p.prompt.errorMsg = msg

					return
				}
			}

			p.prompt.state = StateSubmit
		case key == KeySpace:
			p.toggleHighlighted()
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
		case key == "a" || key == "A":
			p.toggleAll()
		}
	})

	_, err = p.prompt.run()

	if err != nil {
		return nil, err
	}

	return p.SelectedValues(), nil
}

func (p *MultiSelectPrompt) toggleHighlighted() {
	if p.Scrollable.Highlighted() < len(p.options) {
		key := p.options[p.Scrollable.Highlighted()].Key
		p.selected[key] = !p.selected[key]
	}
}

func (p *MultiSelectPrompt) toggleAll() {
	allSelected := true

	for _, opt := range p.options {
		if !p.selected[opt.Key] {
			allSelected = false

			break
		}
	}

	for _, opt := range p.options {
		p.selected[opt.Key] = !allSelected
	}
}

// IsSelected returns whether the given key is selected.
func (p *MultiSelectPrompt) IsSelected(key string) bool {
	return p.selected[key]
}

// SelectedValues returns the selected keys in option order.
func (p *MultiSelectPrompt) SelectedValues() []string {
	var result []string

	for _, opt := range p.options {
		if p.selected[opt.Key] {
			result = append(result, opt.Key)
		}
	}

	return result
}

func (p *MultiSelectPrompt) selectedLabels() []string {
	var result []string

	for _, opt := range p.options {
		if p.selected[opt.Key] {
			result = append(result, opt.Label)
		}
	}

	return result
}
