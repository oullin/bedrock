package prompts

import (
	"fmt"
	"strings"
)

func defaultTheme() *Theme {
	return &Theme{
		TextRenderer:         renderText,
		PasswordRenderer:     renderPassword,
		ConfirmRenderer:      renderConfirm,
		NumberRenderer:       renderNumber,
		PauseRenderer:        renderPause,
		SelectRenderer:       renderSelect,
		MultiSelectRenderer:  renderMultiSelect,
		TextareaRenderer:     renderTextarea,
		SuggestRenderer:      renderSuggest,
		AutocompleteRenderer: renderAutocomplete,
		SearchRenderer:       renderSearch,
		MultiSearchRenderer:  renderMultiSearch,
		NoteRenderer:         renderNote,
		TableRenderer:        renderTable,
		GridRenderer:         renderGrid,
		SpinnerRenderer:      renderSpinner,
		ProgressRenderer:     renderProgress,
		DataTableRenderer:    renderDataTable,
	}
}

// --- Symbols ---
const (
	symbolActive   = "◆"
	symbolDone     = "◇"
	symbolCancel   = "◇"
	symbolError    = "◆"
	symbolPointer  = "›"
	symbolRadio    = "●"
	symbolRadioOff = "○"
	symbolCheck    = "◼"
	symbolUncheck  = "◻"
	symbolBar      = "│"
	symbolBarEnd   = "└"
	symbolCorner   = "╭"
)

// --- Text Prompt ---

func renderText(p *TextPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.Value())))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(Strikethrough(p.Value()))))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	case StateError:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolError), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.TypedValue.AddCursor(p.Value(), 0, "")))
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolBarEnd), Yellow(p.prompt.errorMsg)))

	default: // active/initial
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		val := p.Value()

		if val == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.TypedValue.AddCursor(val, 0, "")))
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Password Prompt ---

func renderPassword(p *PasswordPrompt, state State) string {
	var b strings.Builder
	masked := strings.Repeat("•", len([]rune(p.Value())))

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(masked)))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	case StateError:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolError), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, masked+"│"))
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolBarEnd), Yellow(p.prompt.errorMsg)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))

		if masked == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, masked+"│"))
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Confirm Prompt ---

func renderConfirm(p *ConfirmPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		label := p.noLabel

		if p.confirmed {
			label = p.yesLabel
		}

		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(label)))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		yes := "  " + p.yesLabel + "  "
		no := "  " + p.noLabel + "  "

		if p.confirmed {
			yes = Inverse(" " + p.yesLabel + " ")
		} else {
			no = Inverse(" " + p.noLabel + " ")
		}

		b.WriteString(fmt.Sprintf("  %s %s / %s\n", symbolBar, yes, no))

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Number Prompt ---

func renderNumber(p *NumberPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.Value())))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(Strikethrough(p.Value()))))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	case StateError:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolError), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.TypedValue.AddCursor(p.Value(), 0, "")))
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolBarEnd), Yellow(p.prompt.errorMsg)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		val := p.Value()

		if val == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.TypedValue.AddCursor(val, 0, "")))
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Pause Prompt ---

func renderPause(p *PausePrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Dim(p.message)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.message)))
	}

	return b.String()
}

// --- Select Prompt ---

func renderSelect(p *SelectPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))

		if val := p.selectedLabel(); val != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(val)))
		}

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		start, end := p.Scrollable.Visible()

		for i := start; i < end; i++ {
			opt := p.options[i]

			if i == p.Scrollable.Highlighted() {
				b.WriteString(fmt.Sprintf("  %s %s %s\n", symbolBar, Cyan(symbolPointer), opt.Label))
			} else {
				b.WriteString(fmt.Sprintf("  %s   %s\n", symbolBar, Dim(opt.Label)))
			}
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- MultiSelect Prompt ---

func renderMultiSelect(p *MultiSelectPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		labels := p.selectedLabels()

		if len(labels) > 0 {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(strings.Join(labels, ", "))))
		}

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		start, end := p.Scrollable.Visible()

		for i := start; i < end; i++ {
			opt := p.options[i]
			selected := p.IsSelected(opt.Key)
			checkbox := symbolUncheck

			if selected {
				checkbox = Green(symbolCheck)
			}

			if i == p.Scrollable.Highlighted() {
				b.WriteString(fmt.Sprintf("  %s %s %s %s\n", symbolBar, Cyan(symbolPointer), checkbox, opt.Label))
			} else {
				b.WriteString(fmt.Sprintf("  %s   %s %s\n", symbolBar, checkbox, Dim(opt.Label)))
			}
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Textarea Prompt ---

func renderTextarea(p *TextareaPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		val := p.Value()

		if len(val) > 50 {
			val = val[:50] + "…"
		}

		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(val)))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		val := p.Value()

		if val == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else {
			lines := strings.Split(val, "\n")

			for _, line := range lines {
				b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, line))
			}
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Suggest Prompt ---

func renderSuggest(p *SuggestPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.Value())))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		val := p.Value()

		if val == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.TypedValue.AddCursor(val, 0, "")))
		}

		matches := p.Matches()

		if len(matches) > 0 {
			start, end := p.Scrollable.Visible()

			for i := start; i < end; i++ {
				m := matches[i]

				if i == p.Scrollable.Highlighted() {
					b.WriteString(fmt.Sprintf("  %s %s %s\n", symbolBar, Cyan(symbolPointer), m))
				} else {
					b.WriteString(fmt.Sprintf("  %s   %s\n", symbolBar, Dim(m)))
				}
			}
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Autocomplete Prompt ---

func renderAutocomplete(p *AutocompletePrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.Value())))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		val := p.Value()
		ghost := p.GhostText()

		if val == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else if ghost != "" {
			b.WriteString(fmt.Sprintf("  %s %s%s\n", symbolBar, p.TypedValue.AddCursor(val, 0, ""), Dim(ghost)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.TypedValue.AddCursor(val, 0, "")))
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Search Prompt ---

func renderSearch(p *SearchPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))

		if lbl := p.selectedLabel(); lbl != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(lbl)))
		}

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		val := p.searchValue

		if val == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.searchTyped.AddCursor(val, 0, "")))
		}

		matches := p.currentMatches

		if len(matches) > 0 {
			start, end := p.Scrollable.Visible()

			for i := start; i < end; i++ {
				opt := matches[i]

				if i == p.Scrollable.Highlighted() {
					b.WriteString(fmt.Sprintf("  %s %s %s\n", symbolBar, Cyan(symbolPointer), opt.Label))
				} else {
					b.WriteString(fmt.Sprintf("  %s   %s\n", symbolBar, Dim(opt.Label)))
				}
			}
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- MultiSearch Prompt ---

func renderMultiSearch(p *MultiSearchPrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))
		labels := p.selectedLabels()

		if len(labels) > 0 {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(strings.Join(labels, ", "))))
		}

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))
		val := p.searchValue

		if val == "" && p.placeholder != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, Dim(p.placeholder)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, p.searchTyped.AddCursor(val, 0, "")))
		}

		matches := p.currentMatches

		if len(matches) > 0 {
			start, end := p.Scrollable.Visible()

			for i := start; i < end; i++ {
				opt := matches[i]
				selected := p.IsSelected(opt.Key)
				checkbox := symbolUncheck

				if selected {
					checkbox = Green(symbolCheck)
				}

				if i == p.Scrollable.Highlighted() {
					b.WriteString(fmt.Sprintf("  %s %s %s %s\n", symbolBar, Cyan(symbolPointer), checkbox, opt.Label))
				} else {
					b.WriteString(fmt.Sprintf("  %s   %s %s\n", symbolBar, checkbox, Dim(opt.Label)))
				}
			}
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// --- Note ---

func renderNote(message string, noteType string) string {
	var b strings.Builder

	switch noteType {
	case "error":
		b.WriteString(fmt.Sprintf("  %s %s\n", Red("ERROR"), message))
	case "warning":
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow("WARN"), message))
	case "info":
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan("INFO"), message))
	case "alert":
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow("ALERT"), message))
	case "intro":
		b.WriteString(fmt.Sprintf("\n  %s %s\n\n", Cyan(symbolCorner), Bold(message)))
	case "outro":
		b.WriteString(fmt.Sprintf("\n  %s %s\n\n", Green(symbolBarEnd), message))
	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, message))
	}

	return b.String()
}

// --- Table ---

func renderTable(headers []string, rows [][]string) string {
	if len(rows) == 0 && len(headers) == 0 {
		return ""
	}

	numCols := len(headers)

	for _, row := range rows {
		if len(row) > numCols {
			numCols = len(row)
		}
	}

	// Calculate column widths.
	widths := make([]int, numCols)

	for i, h := range headers {
		if visibleLen(h) > widths[i] {
			widths[i] = visibleLen(h)
		}
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < numCols && visibleLen(cell) > widths[i] {
				widths[i] = visibleLen(cell)
			}
		}
	}

	var b strings.Builder

	// Header separator.
	if len(headers) > 0 {
		b.WriteString("  ")

		for i, h := range headers {
			if i > 0 {
				b.WriteString("  ")
			}

			b.WriteString(Bold(padRight(h, widths[i])))
		}

		b.WriteByte('\n')

		b.WriteString("  ")

		for i, w := range widths {
			if i > 0 {
				b.WriteString("  ")
			}

			b.WriteString(strings.Repeat("─", w))
		}

		b.WriteByte('\n')
	}

	// Rows.
	for _, row := range rows {
		b.WriteString("  ")

		for i := 0; i < numCols; i++ {
			if i > 0 {
				b.WriteString("  ")
			}

			cell := ""

			if i < len(row) {
				cell = row[i]
			}

			b.WriteString(padRight(cell, widths[i]))
		}

		b.WriteByte('\n')
	}

	return b.String()
}

// --- Grid ---

func renderGrid(items []string, maxWidth int) string {
	if len(items) == 0 {
		return ""
	}

	// Find the widest item.
	maxItem := 0

	for _, item := range items {
		w := visibleLen(item)

		if w > maxItem {
			maxItem = w
		}
	}

	colWidth := maxItem + 2
	cols := maxWidth / colWidth

	if cols < 1 {
		cols = 1
	}

	var b strings.Builder

	for i, item := range items {
		if i > 0 && i%cols == 0 {
			b.WriteByte('\n')
		}

		b.WriteString(padRight(item, colWidth))
	}

	b.WriteByte('\n')

	return b.String()
}

// --- Spinner ---

func renderSpinner(message string) string {
	return fmt.Sprintf("  %s %s\n", Cyan("◒"), message)
}

// --- Progress ---

func renderProgress(label string, percentage float64, hint string) string {
	var b strings.Builder
	width := 30
	filled := int(percentage * float64(width))

	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := int(percentage * 100)

	b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(label)))
	b.WriteString(fmt.Sprintf("  %s %s %d%%\n", symbolBar, Cyan(bar), pct))

	if hint != "" {
		b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(hint)))
	}

	return b.String()
}

// --- DataTable ---

func renderDataTable(p *DataTablePrompt, state State) string {
	var b strings.Builder

	switch state {
	case StateSubmit:
		b.WriteString(fmt.Sprintf("  %s %s\n", Green(symbolDone), Bold(p.label)))

	case StateCancel:
		b.WriteString(fmt.Sprintf("  %s %s\n", Yellow(symbolCancel), Bold(p.label)))
		b.WriteString(fmt.Sprintf("  %s\n", Yellow(symbolBarEnd)))

	default:
		b.WriteString(fmt.Sprintf("  %s %s\n", Cyan(symbolActive), Bold(p.label)))

		if p.searchValue != "" || p.placeholder != "" {
			searchDisplay := p.searchValue

			if searchDisplay == "" {
				searchDisplay = Dim(p.placeholder)
			} else {
				searchDisplay = p.searchTyped.AddCursor(searchDisplay, 0, "")
			}

			b.WriteString(fmt.Sprintf("  %s %s\n", symbolBar, searchDisplay))
		}

		filtered := p.FilteredRows()

		if len(filtered) > 0 {
			start, end := p.Scrollable.Visible()

			for i := start; i < end && i < len(filtered); i++ {
				row := filtered[i]
				prefix := "  "

				if i == p.Scrollable.Highlighted() {
					prefix = Cyan(symbolPointer) + " "
				}

				b.WriteString(fmt.Sprintf("  %s %s%s\n", symbolBar, prefix, strings.Join(row, "  ")))
			}
		}

		if p.hint != "" {
			b.WriteString(fmt.Sprintf("  %s %s\n", Dim(symbolBarEnd), Dim(p.hint)))
		}
	}

	return b.String()
}

// padRight pads a string to the given visible width with spaces.
func padRight(s string, width int) string {
	vw := visibleLen(s)

	if vw >= width {
		return s
	}

	return s + strings.Repeat(" ", width-vw)
}
