package props

// MergeProp tells the Inertia.js client to merge new data with existing
// client-side data rather than replacing it. Useful for infinite scroll
// and paginated lists.
type MergeProp struct {
	Value any
	deep  bool
}

// Merge creates a prop whose value is merged (appended) with the
// existing client-side data on navigation.
func Merge(val any) MergeProp {
	return MergeProp{Value: val}
}

// DeepMerge creates a prop whose value is deep-merged with the
// existing client-side data on navigation.
func DeepMerge(val any) MergeProp {
	return MergeProp{Value: val, deep: true}
}

// IsDeep reports whether this is a deep-merge prop.
func (p MergeProp) IsDeep() bool { return p.deep }

func (p MergeProp) Prop() any { return p.Value }
