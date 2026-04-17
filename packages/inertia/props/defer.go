package props

import "strings"

// DeferProp is excluded from the initial response. The Inertia.js client
// loads deferred props asynchronously after the page renders.
// Props can be grouped so multiple deferred values are fetched together.

const defaultDeferGroup = "default"

type DeferProp struct {
	Value any
	Group string
	merge bool
}

// Defer creates a deferred prop. An optional group name controls which
// deferred props are loaded together; it defaults to "default".
func Defer(val any, group ...string) DeferProp {
	g := defaultDeferGroup

	if len(group) > 0 && strings.TrimSpace(group[0]) != "" {
		g = group[0]
	}

	return DeferProp{Value: val, Group: g}
}

// Merge returns a copy of the DeferProp with merge enabled, telling
// the client to merge the result into existing data.
func (p DeferProp) Merge() DeferProp {
	p.merge = true

	return p
}

// IsMerge reports whether this deferred prop should be merged.
func (p DeferProp) IsMerge() bool { return p.merge }

func (p DeferProp) Prop() any { return p.Value }
