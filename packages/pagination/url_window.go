package pagination

// UrlWindow calculates page link windows for length-aware paginators,
// determining which page numbers to display in a navigation bar.
type UrlWindow[T any] struct {
	paginator *LengthAwarePaginator[T]
}

// NewUrlWindow creates a new UrlWindow for the given paginator.
func NewUrlWindow[T any](paginator *LengthAwarePaginator[T]) *UrlWindow[T] {
	return &UrlWindow[T]{paginator: paginator}
}

// MakeUrlWindow is a convenience function that creates a UrlWindow and returns
// the computed page link window.
func MakeUrlWindow[T any](paginator *LengthAwarePaginator[T], onEachSide ...int) map[string]map[int]string {
	w := NewUrlWindow(paginator)

	return w.Get(onEachSide...)
}

// Get returns the page link window as a map with "first", "slider", and "last"
// keys. Each value is a map of page numbers to URLs, or nil if that section is
// not applicable.
func (w *UrlWindow[T]) Get(onEachSide ...int) map[string]map[int]string {
	side := w.paginator.getOnEachSide()

	if len(onEachSide) > 0 {
		side = onEachSide[0]
	}

	if w.paginator.LastPage() < (side*2)+8 {
		return w.getSmallSlider()
	}

	return w.getUrlSlider(side)
}

// HasPages reports whether the paginator has more than one page.
func (w *UrlWindow[T]) HasPages() bool {
	return w.paginator.LastPage() > 1
}

func (w *UrlWindow[T]) getSmallSlider() map[string]map[int]string {
	return map[string]map[int]string{
		"first":  w.paginator.GetUrlRange(1, w.paginator.LastPage()),
		"slider": nil,
		"last":   nil,
	}
}

func (w *UrlWindow[T]) getUrlSlider(onEachSide int) map[string]map[int]string {
	window := onEachSide + 4

	if w.paginator.CurrentPage() <= window {
		return w.getSliderTooCloseToBeginning(window, onEachSide)
	}

	if w.paginator.CurrentPage() > w.paginator.LastPage()-window {
		return w.getSliderTooCloseToEnding(window, onEachSide)
	}

	return w.getFullSlider(onEachSide)
}

func (w *UrlWindow[T]) getSliderTooCloseToBeginning(window, onEachSide int) map[string]map[int]string {
	return map[string]map[int]string{
		"first":  w.paginator.GetUrlRange(1, window+onEachSide),
		"slider": nil,
		"last":   w.paginator.GetUrlRange(w.paginator.LastPage()-1, w.paginator.LastPage()),
	}
}

func (w *UrlWindow[T]) getSliderTooCloseToEnding(window, onEachSide int) map[string]map[int]string {
	last := w.paginator.GetUrlRange(w.paginator.LastPage()-(window+onEachSide)+1, w.paginator.LastPage())

	return map[string]map[int]string{
		"first":  w.paginator.GetUrlRange(1, 2),
		"slider": nil,
		"last":   last,
	}
}

func (w *UrlWindow[T]) getFullSlider(onEachSide int) map[string]map[int]string {
	return map[string]map[int]string{
		"first":  w.paginator.GetUrlRange(1, 2),
		"slider": w.paginator.GetUrlRange(w.paginator.CurrentPage()-onEachSide, w.paginator.CurrentPage()+onEachSide),
		"last":   w.paginator.GetUrlRange(w.paginator.LastPage()-1, w.paginator.LastPage()),
	}
}
