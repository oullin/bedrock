package props

import "github.com/bedrock/packages/inertia/protocol"

// MergeAll combines multiple prop maps in order. Later sources
// override earlier ones, giving a natural precedence chain:
// shared → context → page → validation errors.

// Proper is implemented by prop type wrappers to provide their
// underlying value during resolution.
type Proper interface {
	Prop() any
}

// TryProper extends Proper for values whose resolution may fail.
type TryProper interface {
	TryProp() (any, error)
}

func MergeAll(sources ...protocol.Props) protocol.Props {
	size := 0

	for _, s := range sources {
		size += len(s)
	}

	merged := make(protocol.Props, size)

	for _, s := range sources {
		for k, v := range s {
			merged[k] = v
		}
	}

	return merged
}

// resolve evaluates a raw prop value. It unwraps Proper/TryProper
// wrappers first, then evaluates func() any / func() (any, error)
// lazy values.
func resolve(val any) (any, error) {
	// Unwrap prop type wrappers to get the inner value.
	for {
		switch v := val.(type) {
		case TryProper:
			next, err := v.TryProp()

			if err != nil {
				return nil, err
			}

			val = next
		case Proper:
			val = v.Prop()
		default:
			goto resolved
		}
	}

	// Evaluate lazy function values.
resolved:
	switch v := val.(type) {
	case func() (any, error):
		return v()
	case func() any:
		return v(), nil
	default:
		return val, nil
	}
}
