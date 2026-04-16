package prompts

// Scrollable provides scrolling viewport state for prompts with many options.
// Embedded by SelectPrompt, MultiSelectPrompt, SearchPrompt, etc.
type Scrollable struct {
	scroll       int
	highlighted  int
	firstVisible int
	total        int
}

// InitScrolling sets up the scrollable state.
func (s *Scrollable) InitScrolling(total, scroll int) {
	s.total = total
	s.scroll = scroll
	s.highlighted = 0
	s.firstVisible = 0
}

// Highlighted returns the currently highlighted index.
func (s *Scrollable) Highlighted() int {
	return s.highlighted
}

// SetHighlighted sets the highlighted index, clamped to valid bounds.
func (s *Scrollable) SetHighlighted(index int) {
	if index < 0 {
		index = 0
	}

	if index >= s.total {
		index = s.total - 1
	}

	s.highlighted = index
	s.scrollToHighlighted()
}

// HighlightPrevious moves highlight up by one, wrapping to the end.
func (s *Scrollable) HighlightPrevious(count int) {
	if count <= 0 {
		count = 1
	}

	s.highlighted -= count

	if s.highlighted < 0 {
		s.highlighted = s.total - 1
	}

	s.scrollToHighlighted()
}

// HighlightNext moves highlight down by one, wrapping to the start.
func (s *Scrollable) HighlightNext(count int) {
	if count <= 0 {
		count = 1
	}

	s.highlighted += count

	if s.highlighted >= s.total {
		s.highlighted = 0
	}

	s.scrollToHighlighted()
}

// HighlightFirst moves highlight to the first item.
func (s *Scrollable) HighlightFirst() {
	s.highlighted = 0
	s.scrollToHighlighted()
}

// HighlightLast moves highlight to the last item.
func (s *Scrollable) HighlightLast() {
	s.highlighted = s.total - 1
	s.scrollToHighlighted()
}

// scrollToHighlighted adjusts firstVisible to keep highlighted in view.
func (s *Scrollable) scrollToHighlighted() {
	if s.scroll <= 0 || s.total <= s.scroll {
		s.firstVisible = 0

		return
	}

	if s.highlighted < s.firstVisible {
		s.firstVisible = s.highlighted
	} else if s.highlighted >= s.firstVisible+s.scroll {
		s.firstVisible = s.highlighted - s.scroll + 1
	}
}

// Visible returns the range [start, end) of visible items.
func (s *Scrollable) Visible() (start, end int) {
	start = s.firstVisible
	end = start + s.scroll

	if end > s.total {
		end = s.total
	}

	return start, end
}

// UpdateTotal updates the total items count and re-clamps.
func (s *Scrollable) UpdateTotal(total int) {
	s.total = total

	if s.highlighted >= total {
		s.highlighted = total - 1
	}

	if s.highlighted < 0 {
		s.highlighted = 0
	}

	s.scrollToHighlighted()
}

// ReduceScrollToFitTerminal reduces the scroll if it exceeds the terminal height.
func (s *Scrollable) ReduceScrollToFitTerminal(termLines, reservedLines int) {
	maxScroll := termLines - reservedLines

	if maxScroll < 1 {
		maxScroll = 1
	}

	if s.scroll > maxScroll {
		s.scroll = maxScroll
	}
}
