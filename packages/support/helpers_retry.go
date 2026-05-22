package support

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Retry attempts to execute the given callback the specified number of times.
// If all attempts fail, the last error is returned.
// The sleep argument controls delay between retries and can be:
//   - int: milliseconds to sleep between attempts
//   - []int: backoff schedule in milliseconds (cycles if exhausted)
//   - func(int) time.Duration: callback receiving attempt number, returns sleep duration
func Retry(times int, fn func(attempt int) error, sleep ...any) error {
	var lastErr error

	for attempt := 1; attempt <= times; attempt++ {
		lastErr = fn(attempt)

		if lastErr == nil {
			return nil
		}

		// Don't sleep after the last attempt
		if attempt == times {
			break
		}

		if len(sleep) > 0 {
			d := resolveSleepDuration(sleep[0], attempt)

			if d > 0 {
				time.Sleep(d)
			}
		}
	}

	return lastErr
}

// RetryWhen retries the callback only when the given condition returns true.
func RetryWhen(times int, fn func(attempt int) error, when func(error) bool, sleep ...any) error {
	var lastErr error

	for attempt := 1; attempt <= times; attempt++ {
		lastErr = fn(attempt)

		if lastErr == nil {
			return nil
		}

		if !when(lastErr) {
			return lastErr
		}

		if attempt == times {
			break
		}

		if len(sleep) > 0 {
			d := resolveSleepDuration(sleep[0], attempt)

			if d > 0 {
				time.Sleep(d)
			}
		}
	}

	return lastErr
}

// Throw returns an error if condition is true.
//
// This helper supports common upstream throw_if/throw_unless variants by accepting
// an error, string message, or builder function.
func Throw(condition bool, throwable any, args ...any) error {
	if !condition {
		return nil
	}

	return resolveThrowable(throwable, args...)
}

// ThrowUnless returns an error when the condition is false.
//
// It accepts the same throwable formats as Throw().
func ThrowUnless(condition bool, throwable any, args ...any) error {
	return Throw(!condition, throwable, args...)
}

func resolveSleepDuration(sleep any, attempt int) time.Duration {
	switch v := sleep.(type) {
	case int:
		return time.Duration(v) * time.Millisecond
	case []int:
		if len(v) == 0 {
			return 0
		}

		idx := attempt - 1

		if idx >= len(v) {
			idx = len(v) - 1
		}

		return time.Duration(v[idx]) * time.Millisecond
	case func(int) time.Duration:
		return v(attempt)
	case time.Duration:
		return v
	}

	return 0
}

// ThrowIf throws the given exception if the condition is true.
func ThrowIf(condition bool, err error) error {
	return Throw(condition, err)
}

// PanicIf panics with the given error if the condition is true.
func PanicIf(condition bool, err error) {
	if condition {
		panic(err)
	}
}

// PanicUnless panics with the given error unless the condition is true.
func PanicUnless(condition bool, err error) {
	PanicIf(!condition, err)
}

func resolveThrowable(throwable any, args ...any) error {
	if throwable == nil {
		return errors.New("RuntimeException")
	}

	switch value := throwable.(type) {
	case error:
		return value
	case func(string) error:
		var arg string

		if len(args) > 0 {
			arg = fmt.Sprint(args[0])
		}

		return value(arg)
	case string:
		if len(args) == 0 {
			return errors.New(value)
		}

		return errors.New(fmt.Sprintf(value, args...))
	case func() error:
		return value()
	case func(string) string:
		var arg string

		if len(args) > 0 {
			arg = fmt.Sprint(args[0])
		}

		return errors.New(value(arg))
	case func(...any) error:
		return value(args...)
	case func(any) error:
		arg := firstArg(args)

		return value(arg)
	case func() string:
		return errors.New(value())
	case func(...any) string:
		return errors.New(value(args...))
	case func(any) string:
		arg := firstArg(args)

		return errors.New(value(arg))
	default:
		if len(args) > 0 {
			reflected, ok := callSingleArgThrowable(reflect.ValueOf(throwable), args[0])

			if ok {
				return reflected
			}

			return fmt.Errorf("%v", throwable)
		}

		if err, ok := throwable.(fmt.Stringer); ok {
			return errors.New(err.String())
		}

		return fmt.Errorf("%v", throwable)
	}
}

func callSingleArgThrowable(fn reflect.Value, arg any) (error, bool) {
	if !fn.IsValid() || fn.Kind() != reflect.Func {
		return nil, false
	}

	fnType := fn.Type()

	if fnType.NumIn() != 1 || fnType.NumOut() != 1 {
		return nil, false
	}

	inType := fnType.In(0)

	var callArg reflect.Value

	if arg == nil {
		callArg = reflect.Zero(inType)
	} else {
		argValue := reflect.ValueOf(arg)

		if !argValue.IsValid() {
			callArg = reflect.Zero(inType)
		} else if argValue.Type().AssignableTo(inType) {
			callArg = argValue
		} else if argValue.Type().ConvertibleTo(inType) {
			callArg = argValue.Convert(inType)
		} else if inType.Kind() == reflect.Interface && argValue.Type().Implements(inType) {
			callArg = argValue
		} else if argValue.CanConvert(inType) {
			callArg = argValue.Convert(inType)
		} else {
			return nil, false
		}
	}

	results := fn.Call([]reflect.Value{callArg})
	result := results[0].Interface()

	err, ok := result.(error)

	if ok {
		return err, true
	}

	if str, ok := result.(string); ok {
		return errors.New(str), true
	}

	return nil, false
}

func firstArg(args []any) any {
	if len(args) == 0 {
		return nil
	}

	return args[0]
}

// Rescue executes the given callback and recovers from any panic,
// returning the error. Useful for wrapping code that may panic.
func Rescue[T any](fn func() T, def ...T) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			case string:
				err = errors.New(e)
			default:
				err = errors.New("unknown panic")
			}

			if len(def) > 0 {
				result = def[0]
			}
		}
	}()

	return fn(), nil
}
