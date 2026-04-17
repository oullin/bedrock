package props

// ScrollProp wraps a value with infinite-scroll metadata expected by
// the Inertia v3 client.
type ScrollProp struct {
	Value        any
	PageName     string
	PreviousPage any
	NextPage     any
	CurrentPage  any
	reset        bool
	merge        bool
}

// Scroll creates a prop with scroll metadata for infinite scroll pages.
func Scroll(val any, pageName string, currentPage, previousPage, nextPage any) ScrollProp {
	return ScrollProp{
		Value:        val,
		PageName:     pageName,
		PreviousPage: previousPage,
		NextPage:     nextPage,
		CurrentPage:  currentPage,
	}
}

// Reset returns a copy of the prop marked as a reset.
func (p ScrollProp) Reset() ScrollProp {
	p.reset = true

	return p
}

// Merge returns a copy of the prop marked to merge on partial reloads.
func (p ScrollProp) Merge() ScrollProp {
	p.merge = true

	return p
}

func (p ScrollProp) IsReset() bool { return p.reset }

func (p ScrollProp) IsMerge() bool { return p.merge }

func (p ScrollProp) Prop() any { return p.Value }
