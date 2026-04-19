package props

import (
	"net/http"
	"strings"

	"github.com/bedrock/packages/inertia/protocol"
)

// Result holds the output of prop resolution: the evaluated props,
// the list of merge-prop keys, and the deferred-prop grouping.
type Result struct {
	Props          map[string]any
	MergeProps     []string
	DeepMergeProps []string
	DeferredProps  map[string][]string
	ScrollProps    map[string]protocol.Scroll
	OnceProps      map[string]protocol.Once
}

// Resolve filters and evaluates the merged props map according to the
// Inertia.js partial-reload protocol headers found on r.

// filterContext holds the request-scoped filter state parsed from
// Inertia protocol headers. It is built once per Resolve call and
// passed to shouldInclude so the method signature stays small.
type filterContext struct {
	isPartial     bool
	onlySet       map[string]struct{}
	exceptSet     map[string]struct{}
	exceptOnceSet map[string]struct{}
	mergeIntent   string
}

// propTraits holds metadata discovered by walking a nested prop
// wrapper chain. Each flag indicates whether the corresponding
// wrapper type was found; the associated fields carry its metadata.
type propTraits struct {
	hasAlways   bool
	hasDefer    bool
	hasOptional bool
	hasOnce     bool
	hasScroll   bool
	hasMerge    bool

	deferGroup string
	deferMerge bool

	scrollMeta  protocol.Scroll
	scrollMerge bool

	mergeDeep bool

	onceExpiresAt *int64
}

func Resolve(r *http.Request, component string, merged protocol.Props) (*Result, error) {
	res := &Result{
		Props:         make(map[string]any, len(merged)),
		DeferredProps: make(map[string][]string),
		ScrollProps:   make(map[string]protocol.Scroll),
		OnceProps:     make(map[string]protocol.Once),
	}

	partialComponent := r.Header.Get(protocol.HeaderPartialComponent)

	fc := filterContext{
		isPartial:     protocol.IsInertiaRequest(r) && partialComponent == component,
		onlySet:       toSet(splitHeader(r.Header.Get(protocol.HeaderPartialData))),
		exceptSet:     toSet(splitHeader(r.Header.Get(protocol.HeaderPartialExcept))),
		exceptOnceSet: toSet(splitHeader(r.Header.Get(protocol.HeaderExceptOnceProps))),
		mergeIntent:   r.Header.Get(protocol.HeaderInfiniteScroll),
	}

	// Stage 1: filter — decide which props to include and collect metadata.
	included := make(map[string]any, len(merged))

	for key, val := range merged {
		ok, err := res.shouldInclude(key, val, fc)

		if err != nil {
			return nil, err
		}

		if ok {
			included[key] = val
		}
	}

	// Stage 2: evaluate — resolve the included props into final values.
	if err := res.evaluate(included); err != nil {
		return nil, err
	}

	return res, nil
}

// evaluate resolves every included prop into its final value and
// stores it in res.Props.
func (res *Result) evaluate(included map[string]any) error {
	for key, val := range included {
		resolved, err := resolve(val)

		if err != nil {
			return err
		}

		res.Props[key] = resolved
	}

	return nil
}

// walkPropChain iterates through the Proper wrapper chain collecting
// trait flags and metadata from every layer. When the same wrapper
// type appears more than once the outermost instance wins.
func walkPropChain(val any) propTraits {
	var t propTraits

	cur := val

	for {
		switch v := cur.(type) {
		case AlwaysProp:
			t.hasAlways = true
			cur = v.Value

		case DeferProp:
			if !t.hasDefer {
				t.hasDefer = true
				t.deferGroup = v.Group
				t.deferMerge = v.IsMerge()
			}

			cur = v.Value

		case OptionalProp:
			t.hasOptional = true
			cur = v.Value

		case OnceProp:
			if !t.hasOnce {
				t.hasOnce = true
				t.onceExpiresAt = v.GetExpiresAt()
			}

			cur = v.Value

		case ScrollProp:
			if !t.hasScroll {
				t.hasScroll = true
				t.scrollMeta = protocol.Scroll{
					PageName:     v.PageName,
					PreviousPage: v.PreviousPage,
					NextPage:     v.NextPage,
					CurrentPage:  v.CurrentPage,
					Reset:        v.IsReset(),
				}
				t.scrollMerge = v.IsMerge()
			}

			cur = v.Value

		case MergeProp:
			if !t.hasMerge {
				t.hasMerge = true
				t.mergeDeep = v.IsDeep()
			}

			cur = v.Value

		default:
			return t
		}
	}
}

// shouldInclude determines whether a prop should be included in the
// response and records metadata (merge keys, deferred groups) as a
// side effect. It returns true when the prop should be resolved and
// added to the output.
func (res *Result) shouldInclude(key string, val any, fc filterContext) (bool, error) {
	traits := walkPropChain(val)

	// AlwaysProp bypasses all filters.
	if traits.hasAlways {
		res.recordMetadata(key, traits, fc.mergeIntent)

		return true, nil
	}

	// Partial reload filtering.
	if fc.isPartial {
		// exceptSet takes precedence.
		if len(fc.exceptSet) > 0 {
			if _, excluded := fc.exceptSet[key]; excluded {
				return false, nil
			}
		}

		if len(fc.onlySet) > 0 {
			if _, included := fc.onlySet[key]; !included {
				return false, nil
			}
		}
	}

	// DeferProp: excluded on initial, included on partial reload.
	if traits.hasDefer && !fc.isPartial {
		res.DeferredProps[traits.deferGroup] = append(res.DeferredProps[traits.deferGroup], key)

		if traits.deferMerge {
			res.MergeProps = append(res.MergeProps, key)
		}

		return false, nil
	}

	// OptionalProp: excluded on initial, only on explicit request.
	if traits.hasOptional {
		if !fc.isPartial {
			return false, nil
		}

		if len(fc.onlySet) > 0 {
			if _, requested := fc.onlySet[key]; !requested {
				return false, nil
			}
		} else {
			return false, nil
		}
	}

	// OnceProp: record metadata, skip if in except-once set.
	if traits.hasOnce {
		res.OnceProps[key] = protocol.Once{Prop: key, ExpiresAt: traits.onceExpiresAt}

		if _, skip := fc.exceptOnceSet[key]; skip {
			return false, nil
		}
	}

	// Record all remaining metadata for included props.
	res.recordMetadata(key, traits, fc.mergeIntent)

	return true, nil
}

// recordMetadata records scroll, merge, and deferred-merge metadata
// from the collected traits for an included prop.
func (res *Result) recordMetadata(key string, traits propTraits, mergeIntent string) {
	if traits.hasScroll {
		res.ScrollProps[key] = traits.scrollMeta

		if traits.scrollMerge || mergeIntent == "append" {
			res.MergeProps = append(res.MergeProps, key)
		}
	}

	if traits.hasMerge {
		if traits.mergeDeep {
			res.DeepMergeProps = append(res.DeepMergeProps, key)
		} else {
			res.MergeProps = append(res.MergeProps, key)
		}

		return
	}

	// DeferProp merge on partial reload (prop is included).
	if traits.hasDefer && traits.deferMerge {
		res.MergeProps = append(res.MergeProps, key)
	}
}

// splitHeader splits a comma-separated header value into trimmed,
// non-empty tokens.
func splitHeader(val string) []string {
	if strings.TrimSpace(val) == "" {
		return nil
	}

	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}

	return out
}

func toSet(items []string) map[string]struct{} {
	if len(items) == 0 {
		return nil
	}

	s := make(map[string]struct{}, len(items))

	for _, item := range items {
		s[item] = struct{}{}
	}

	return s
}
