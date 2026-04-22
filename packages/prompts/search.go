package prompts

import (
	"sort"
	"strings"
)

// SearchOption configures a SearchPrompt.
type SearchOption func(*SearchPrompt)

// SearchWithPlaceholder sets placeholder text.

// SearchWithScroll sets visible option count.

// SearchWithValidate sets validation.

// SearchWithHint sets hint text.

// SearchWithRequired makes the prompt required.

// SearchWithTransform sets the transform function.

// SearchPrompt handles dynamic searchable single selection.
type SearchPrompt struct {
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
}

func SearchWithPlaceholder(s string) SearchOption {
	return func(p *SearchPrompt) { p.placeholder = s }
}

func SearchWithScroll(n int) SearchOption { return func(p *SearchPrompt) { p.scroll = n } }

func SearchWithValidate(fn ValidateFunc) SearchOption {
	return func(p *SearchPrompt) { p.prompt.validate = fn }
}

func SearchWithHint(s string) SearchOption { return func(p *SearchPrompt) { p.hint = s } }

func SearchWithRequired(v any) SearchOption {
	return func(p *SearchPrompt) { p.prompt.required = v }
}

func SearchWithTransform(fn TransformFunc) SearchOption {
	return func(p *SearchPrompt) { p.prompt.transform = fn }
}

// Search displays a dynamic searchable single-select prompt.
func Search(label string, options func(string) map[string]string, opts ...SearchOption) (string, error) {
	p := &SearchPrompt{
		prompt:    newPrompt(),
		label:     label,
		optionsFn: options,
		scroll:    5,
	}

	for _, opt := range opts {
		opt(p)
	}

	// Initial search.
	p.refreshMatches()

	p.prompt.valueFn = func() string {
		if len(p.currentMatches) == 0 {
			return ""
		}

		return p.currentMatches[p.Scrollable.Highlighted()].Key
	}

	p.prompt.renderer = func(state State) string {
		return getTheme().SearchRenderer(p, state)
	}

	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			if len(p.currentMatches) > 0 {
				p.prompt.submit()
			}
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

	return p.prompt.run()
}

func (p *SearchPrompt) refreshMatches() {
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

func (p *SearchPrompt) selectedLabel() string {
	if len(p.currentMatches) == 0 {
		return ""
	}

	idx := p.Scrollable.Highlighted()

	if idx < len(p.currentMatches) {
		return p.currentMatches[idx].Label
	}

	return ""
}
