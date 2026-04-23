package prompts

import (
	"sort"
	"strings"
)

// MultiSearchOption configures a MultiSearchPrompt.
type MultiSearchOption func(*MultiSearchPrompt)

// MultiSearchWithPlaceholder sets placeholder text.

// MultiSearchWithScroll sets visible option count.

// MultiSearchWithRequired makes the prompt required.

// MultiSearchWithValidate sets validation.

// MultiSearchWithHint sets hint text.

// MultiSearchWithTransform sets the transform function.

// MultiSearchPrompt handles dynamic searchable multi-selection.
type MultiSearchPrompt struct {
	prompt         *Prompt
	Scrollable     Scrollable
	searchTyped    TypedValue
	label          string
	placeholder    string
	hint           string
	scroll         int
	optionsFn      func(string) map[string]string
	searchValue    string
	currentMatches []OptionItem
	selected       map[string]string // key -> label
	validateMulti  func([]string) string
}

func MultiSearchWithPlaceholder(s string) MultiSearchOption {
	return func(p *MultiSearchPrompt) { p.placeholder = s }
}

func MultiSearchWithScroll(n int) MultiSearchOption {
	return func(p *MultiSearchPrompt) { p.scroll = n }
}

func MultiSearchWithRequired(v any) MultiSearchOption {
	return func(p *MultiSearchPrompt) { p.prompt.required = v }
}

func MultiSearchWithValidate(fn func([]string) string) MultiSearchOption {
	return func(p *MultiSearchPrompt) { p.validateMulti = fn }
}

func MultiSearchWithHint(s string) MultiSearchOption {
	return func(p *MultiSearchPrompt) { p.hint = s }
}

func MultiSearchWithTransform(fn TransformFunc) MultiSearchOption {
	return func(p *MultiSearchPrompt) { p.prompt.transform = fn }
}

// MultiSearch displays a dynamic searchable multi-select prompt.
func MultiSearch(label string, options func(string) map[string]string, opts ...MultiSearchOption) ([]string, error) {
	p := &MultiSearchPrompt{
		prompt:    newPrompt(),
		label:     label,
		optionsFn: options,
		selected:  make(map[string]string),
		scroll:    5,
	}

	for _, opt := range opts {
		opt(p)
	}

	// Initial search.
	p.refreshMatches()

	p.prompt.valueFn = func() string { return "" }
	p.prompt.renderer = func(state State) string {
		return getTheme().MultiSearchRenderer(p, state)
	}

	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			vals := p.SelectedValues()

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
		case key == KeyBackspace || key == KeyCtrlH:
			if len(p.searchValue) > 0 {
				runes := []rune(p.searchValue)
				p.searchValue = string(runes[:len(runes)-1])
				p.searchTyped.SetValue(p.searchValue)
				p.refreshMatches()
			}
		case key == KeyCtrlU:
			p.searchValue = ""
			p.searchTyped.SetValue("")
			p.refreshMatches()
		default:
			if len(key) > 0 && key[0] >= 32 && key[0] != 127 && !strings.HasPrefix(key, "\x1b") {
				p.searchValue += key
				p.searchTyped.SetValue(p.searchValue)
				p.refreshMatches()
			}
		}
	})

	_, err := p.prompt.run()

	if err != nil {
		return nil, err
	}

	return p.SelectedValues(), nil
}

func (p *MultiSearchPrompt) refreshMatches() {
	result := p.optionsFn(p.searchValue)
	p.currentMatches = nil

	keys := make([]string, 0, len(result))

	for key := range result {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		p.currentMatches = append(p.currentMatches, OptionItem{Key: key, Label: result[key]})
	}

	p.Scrollable.InitScrolling(len(p.currentMatches), p.scroll)
}

func (p *MultiSearchPrompt) toggleHighlighted() {
	if p.Scrollable.Highlighted() < len(p.currentMatches) {
		opt := p.currentMatches[p.Scrollable.Highlighted()]

		if _, ok := p.selected[opt.Key]; ok {
			delete(p.selected, opt.Key)
		} else {
			p.selected[opt.Key] = opt.Label
		}
	}
}

// IsSelected returns whether the given key is selected.
func (p *MultiSearchPrompt) IsSelected(key string) bool {
	_, ok := p.selected[key]

	return ok
}

// SelectedValues returns the selected keys.
func (p *MultiSearchPrompt) SelectedValues() []string {
	result := make([]string, 0, len(p.selected))

	for k := range p.selected {
		result = append(result, k)
	}

	return result
}

func (p *MultiSearchPrompt) selectedLabels() []string {
	result := make([]string, 0, len(p.selected))

	for _, v := range p.selected {
		result = append(result, v)
	}

	return result
}
