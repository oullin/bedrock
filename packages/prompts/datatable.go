package prompts

import (
	"strings"
)

// DataTableOption configures a DataTablePrompt.
type DataTableOption func(*DataTablePrompt)

// DataTableWithScroll sets the visible row count.

// DataTableWithLabel sets the label.

// DataTableWithHint sets hint text.

// DataTableWithRequired makes the prompt required.

// DataTableWithValidate sets validation.

// DataTableWithTransform sets the transform function.

// DataTableWithFilter sets a custom filter function.

// DataTablePrompt handles interactive data table with search/filter and selection.
type DataTablePrompt struct {
	prompt      *Prompt
	Scrollable  Scrollable
	searchTyped TypedValue
	label       string
	placeholder string
	hint        string
	headers     []string
	rows        [][]string
	scroll      int
	searchValue string
	filterFn    func(query string, rows [][]string) [][]string
}

func DataTableWithScroll(n int) DataTableOption {
	return func(p *DataTablePrompt) { p.scroll = n }
}

func DataTableWithLabel(s string) DataTableOption {
	return func(p *DataTablePrompt) { p.label = s }
}

func DataTableWithHint(s string) DataTableOption {
	return func(p *DataTablePrompt) { p.hint = s }
}

func DataTableWithRequired(v any) DataTableOption {
	return func(p *DataTablePrompt) { p.prompt.required = v }
}

func DataTableWithValidate(fn ValidateFunc) DataTableOption {
	return func(p *DataTablePrompt) { p.prompt.validate = fn }
}

func DataTableWithTransform(fn TransformFunc) DataTableOption {
	return func(p *DataTablePrompt) { p.prompt.transform = fn }
}

func DataTableWithFilter(fn func(query string, rows [][]string) [][]string) DataTableOption {
	return func(p *DataTablePrompt) { p.filterFn = fn }
}

// DataTable displays an interactive data table with search and selection.
func DataTable(headers []string, rows [][]string, opts ...DataTableOption) (string, error) {
	p := &DataTablePrompt{
		prompt:      newPrompt(),
		headers:     headers,
		rows:        rows,
		scroll:      10,
		placeholder: "Search...",
	}

	for _, opt := range opts {
		opt(p)
	}

	filtered := p.FilteredRows()
	p.Scrollable.InitScrolling(len(filtered), p.scroll)

	p.prompt.valueFn = func() string {
		fr := p.FilteredRows()

		if len(fr) == 0 {
			return ""
		}

		idx := p.Scrollable.Highlighted()

		if idx < len(fr) {
			return strings.Join(fr[idx], "\t")
		}

		return ""
	}

	p.prompt.renderer = func(state State) string {
		return getTheme().DataTableRenderer(p, state)
	}

	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			if len(p.FilteredRows()) > 0 {
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
				p.refreshFiltered()
			}
		case key == KeyCtrlU:
			p.searchValue = ""
			p.searchTyped.SetValue("")
			p.refreshFiltered()
		default:
			if len(key) > 0 && key[0] >= 32 && key[0] != 127 && !strings.HasPrefix(key, "\x1b") {
				p.searchValue += key
				p.searchTyped.SetValue(p.searchValue)
				p.refreshFiltered()
			}
		}
	})

	return p.prompt.run()
}

// FilteredRows returns the rows matching the current search query.
func (p *DataTablePrompt) FilteredRows() [][]string {
	if p.searchValue == "" {
		return p.rows
	}

	if p.filterFn != nil {
		return p.filterFn(p.searchValue, p.rows)
	}

	lower := strings.ToLower(p.searchValue)

	var result [][]string

	for _, row := range p.rows {
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), lower) {
				result = append(result, row)

				break
			}
		}
	}

	return result
}

func (p *DataTablePrompt) refreshFiltered() {
	filtered := p.FilteredRows()
	p.Scrollable.UpdateTotal(len(filtered))
}
