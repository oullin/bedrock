package queue

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// The `queue` struct tag encodes job dispatch options in a single string.
//
//	type SendEmail struct {
//	    _       struct{} `queue:"tries=3,backoff=1s|5s|10s,timeout=60s,unique_for=5m,fail_on_timeout"`
//	    To      string
//	    Subject string
//	}
//
// This is the Go analogue of the upstream PHP8 attributes (#[Tries], #[Backoff],
// #[Timeout], #[Queue], #[Connection], #[Delay], #[UniqueFor], #[FailOnTimeout],
// #[MaxExceptions], #[DeleteWhenMissingModels]).
//
// Supported keys:
//
//	tries=N                    int, maps to JobOptions.MaxTries
//	max_tries=N                alias of tries
//	max_exceptions=N           int, maps to JobOptions.MaxExceptions
//	timeout=DURATION           time.Duration, e.g. "60s", "2m"
//	delay=DURATION             time.Duration
//	backoff=DUR|DUR|...        pipe-separated time.Duration list
//	retry_until=RFC3339        absolute time, e.g. "2026-05-01T12:00:00Z"
//	unique_for=DURATION        time.Duration
//	fail_on_timeout            bool flag (presence = true)
//	delete_when_missing_models bool flag (presence = true)
//	queue=NAME                 string
//	connection=NAME            string
const jobOptionsTag = "queue"

// ParseJobOptions reads the `queue` struct tag from each field of v and
// returns the decoded JobOptions. v may be a struct or a pointer to a struct.
// If multiple fields carry a `queue` tag the values are merged in declaration
// order (later fields override earlier ones for scalar keys). If no tags are
// present the zero-valued JobOptions is returned.
//
// Callers typically wire this into a HandlerRegistry so each handler parses
// its options exactly once.
func ParseJobOptions(v any) (JobOptions, error) {
	rv := reflect.ValueOf(v)

	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return JobOptions{}, nil
		}

		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return JobOptions{}, fmt.Errorf("queue: ParseJobOptions requires a struct, got %s", rv.Kind())
	}

	var opts JobOptions

	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		tag, ok := rt.Field(i).Tag.Lookup(jobOptionsTag)

		if !ok || tag == "" {
			continue
		}

		if err := applyJobOptionsTag(&opts, tag); err != nil {
			return JobOptions{}, fmt.Errorf("queue: field %s: %w", rt.Field(i).Name, err)
		}
	}

	return opts, nil
}

// applyJobOptionsTag merges one comma-separated tag value into opts.
func applyJobOptionsTag(opts *JobOptions, tag string) error {
	for _, part := range strings.Split(tag, ",") {
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		key, value, hasValue := strings.Cut(part, "=")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if err := applyJobOptionsKey(opts, key, value, hasValue); err != nil {
			return err
		}
	}

	return nil
}

func applyJobOptionsKey(opts *JobOptions, key, value string, hasValue bool) error {
	switch key {
	case "tries", "max_tries":
		return parseIntInto(&opts.MaxTries, key, value, hasValue)
	case "max_exceptions":
		return parseIntInto(&opts.MaxExceptions, key, value, hasValue)
	case "timeout":
		return parseDurationInto(&opts.Timeout, key, value, hasValue)
	case "delay":
		return parseDurationInto(&opts.Delay, key, value, hasValue)
	case "backoff":
		if !hasValue {
			return fmt.Errorf("key %q requires a value", key)
		}

		list, err := parseDurationList(value)

		if err != nil {
			return fmt.Errorf("key %q: %w", key, err)
		}

		opts.Backoff = list

		return nil
	case "retry_until":
		if !hasValue {
			return fmt.Errorf("key %q requires a value", key)
		}

		t, err := time.Parse(time.RFC3339, value)

		if err != nil {
			return fmt.Errorf("key %q: %w", key, err)
		}

		opts.RetryUntil = t

		return nil
	case "unique_for":
		return parseDurationInto(&opts.UniqueFor, key, value, hasValue)
	case "fail_on_timeout":
		return parseBoolFlagInto(&opts.FailOnTimeout, key, value, hasValue)
	case "delete_when_missing_models":
		return parseBoolFlagInto(&opts.DeleteWhenMissingModels, key, value, hasValue)
	case "queue":
		if !hasValue {
			return fmt.Errorf("key %q requires a value", key)
		}

		opts.Queue = value

		return nil
	case "connection":
		if !hasValue {
			return fmt.Errorf("key %q requires a value", key)
		}

		opts.Connection = value

		return nil
	default:
		return fmt.Errorf("unknown key %q", key)
	}
}

func parseIntInto(dst *int, key, value string, hasValue bool) error {
	if !hasValue {
		return fmt.Errorf("key %q requires a value", key)
	}

	n, err := strconv.Atoi(value)

	if err != nil {
		return fmt.Errorf("key %q: %w", key, err)
	}

	*dst = n

	return nil
}

func parseDurationInto(dst *time.Duration, key, value string, hasValue bool) error {
	if !hasValue {
		return fmt.Errorf("key %q requires a value", key)
	}

	d, err := time.ParseDuration(value)

	if err != nil {
		return fmt.Errorf("key %q: %w", key, err)
	}

	*dst = d

	return nil
}

func parseDurationList(value string) ([]time.Duration, error) {
	parts := strings.Split(value, "|")
	out := make([]time.Duration, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if p == "" {
			continue
		}

		d, err := time.ParseDuration(p)

		if err != nil {
			return nil, err
		}

		out = append(out, d)
	}

	return out, nil
}

// parseBoolFlagInto treats a bare key ("fail_on_timeout") as true and a
// key=value ("fail_on_timeout=false") as the parsed boolean.
func parseBoolFlagInto(dst *bool, key, value string, hasValue bool) error {
	if !hasValue {
		*dst = true

		return nil
	}

	b, err := strconv.ParseBool(value)

	if err != nil {
		return fmt.Errorf("key %q: %w", key, err)
	}

	*dst = b

	return nil
}
