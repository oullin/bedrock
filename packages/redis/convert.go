package redis

import (
	"fmt"
	"strconv"
)

// toString coerces a reply value to a string. nil becomes "" with ErrNil.
func toString(v any) (string, error) {
	switch x := v.(type) {
	case nil:
		return "", ErrNil
	case string:
		return x, nil
	case []byte:
		return string(x), nil
	case int64:
		return strconv.FormatInt(x, 10), nil
	case int:
		return strconv.Itoa(x), nil
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64), nil
	case bool:
		if x {
			return "1", nil
		}

		return "0", nil
	default:
		return fmt.Sprintf("%v", x), nil
	}
}

// toInt64 coerces a reply value to int64.
func toInt64(v any) (int64, error) {
	switch x := v.(type) {
	case nil:
		return 0, ErrNil
	case int64:
		return x, nil
	case int:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case uint64:
		return int64(x), nil
	case float64:
		return int64(x), nil
	case string:
		n, err := strconv.ParseInt(x, 10, 64)

		if err != nil {
			return 0, ErrUnexpectedReply
		}

		return n, nil
	case []byte:
		n, err := strconv.ParseInt(string(x), 10, 64)

		if err != nil {
			return 0, ErrUnexpectedReply
		}

		return n, nil
	case bool:
		if x {
			return 1, nil
		}

		return 0, nil
	default:
		return 0, ErrUnexpectedReply
	}
}

// toFloat64 coerces a reply value to float64.
func toFloat64(v any) (float64, error) {
	switch x := v.(type) {
	case nil:
		return 0, ErrNil
	case float64:
		return x, nil
	case int64:
		return float64(x), nil
	case int:
		return float64(x), nil
	case string:
		n, err := strconv.ParseFloat(x, 64)

		if err != nil {
			return 0, ErrUnexpectedReply
		}

		return n, nil
	case []byte:
		n, err := strconv.ParseFloat(string(x), 64)

		if err != nil {
			return 0, ErrUnexpectedReply
		}

		return n, nil
	default:
		return 0, ErrUnexpectedReply
	}
}

// toSlice coerces a reply value to []any.
func toSlice(v any) ([]any, error) {
	switch x := v.(type) {
	case nil:
		return nil, nil
	case []any:
		return x, nil
	case []string:
		out := make([]any, len(x))

		for i, s := range x {
			out[i] = s
		}

		return out, nil
	default:
		return nil, ErrUnexpectedReply
	}
}

// toStringSlice coerces a reply value to []string.
func toStringSlice(v any) ([]string, error) {
	switch x := v.(type) {
	case nil:
		return nil, nil
	case []string:
		return x, nil
	case []any:
		out := make([]string, len(x))

		for i, e := range x {
			s, err := toString(e)

			if err != nil && err != ErrNil {
				return nil, err
			}

			out[i] = s
		}

		return out, nil
	default:
		return nil, ErrUnexpectedReply
	}
}

// toStringMap coerces a flat array reply (k1, v1, k2, v2, ...) or a
// map[string]string to a map. Redis HGETALL replies use the flat form.
func toStringMap(v any) (map[string]string, error) {
	switch x := v.(type) {
	case nil:
		return nil, nil
	case map[string]string:
		return x, nil
	case map[string]any:
		out := make(map[string]string, len(x))

		for k, val := range x {
			s, _ := toString(val)
			out[k] = s
		}

		return out, nil
	case []any:
		if len(x)%2 != 0 {
			return nil, ErrUnexpectedReply
		}

		out := make(map[string]string, len(x)/2)

		for i := 0; i < len(x); i += 2 {
			k, _ := toString(x[i])
			v, _ := toString(x[i+1])
			out[k] = v
		}

		return out, nil
	}

	return nil, ErrUnexpectedReply
}

// toAnySlice converts []string to []any for variadic command builders.
func toAnySlice(ss []string) []any {
	out := make([]any, len(ss))

	for i, s := range ss {
		out[i] = s
	}

	return out
}
