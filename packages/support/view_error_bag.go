package support

// ViewErrorBag stores named MessageBag instances for views.
type ViewErrorBag struct {
	bags map[string]*MessageBag
}

// NewViewErrorBag creates an empty view error bag.
func NewViewErrorBag() *ViewErrorBag {
	return &ViewErrorBag{bags: make(map[string]*MessageBag)}
}

// Put stores a named message bag.
func (v *ViewErrorBag) Put(name string, bag *MessageBag) *ViewErrorBag {
	if bag == nil {
		bag = NewMessageBag()
	}

	v.bags[name] = bag

	return v
}

// HasBag reports whether a named bag exists.
func (v *ViewErrorBag) HasBag(name string) bool {
	_, ok := v.bags[name]

	return ok
}

// Get returns a named bag, creating it when absent.
func (v *ViewErrorBag) Get(name string) *MessageBag {
	if bag, ok := v.bags[name]; ok {
		return bag
	}

	bag := NewMessageBag()
	v.bags[name] = bag

	return bag
}

// Bags returns a shallow copy of the named message bags.
func (v *ViewErrorBag) Bags() map[string]*MessageBag {
	bags := make(map[string]*MessageBag, len(v.bags))

	for key, bag := range v.bags {
		bags[key] = bag
	}

	return bags
}

// Any reports whether any named bag contains messages.
func (v *ViewErrorBag) Any() bool {
	for _, bag := range v.bags {
		if bag.IsNotEmpty() {
			return true
		}
	}

	return false
}

// Count returns the total number of messages across all bags.
func (v *ViewErrorBag) Count() int {
	total := 0

	for _, bag := range v.bags {
		total += bag.Count()
	}

	return total
}

// String returns the default bag string representation.
func (v *ViewErrorBag) String() string {
	return v.Get("default").String()
}
