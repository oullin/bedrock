// Package serialize provides PHP-compatible serialization functions for Go.
// Ported from github.com/stlong5/phpserialize v1.0.1.
package serialize

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// PHPObject represents a PHP object in Go
type PHPObject struct {
	ClassName  string
	Properties map[string]interface{}
}

type marshalConfig struct {
	phpStrict bool
	maxDepth  int
}

type unmarshalConfig struct {
	allowedClasses map[string]bool
	allowAll       bool
	maxDepth       int
}

// Option allows customization of serialize/un-serialize behavior
// This unified interface works for both Marshal and Unmarshal
type Option interface {
	applyMarshal(*marshalConfig)
	applyUnmarshal(*unmarshalConfig)
}

// maxDepthOption implements Option for depth limiting
type maxDepthOption struct {
	depth int
}

func (o maxDepthOption) applyMarshal(cfg *marshalConfig) {
	cfg.maxDepth = o.depth
}

func (o maxDepthOption) applyUnmarshal(cfg *unmarshalConfig) {
	if o.depth == 0 {
		cfg.maxDepth = 4096 // PHP default when 0 specified
	} else {
		cfg.maxDepth = o.depth
	}
}

// strictPHPOption implements Option for PHP strict mode
type strictPHPOption struct {
	strict bool
}

func (o strictPHPOption) applyMarshal(cfg *marshalConfig) {
	cfg.phpStrict = o.strict
}

func (o strictPHPOption) applyUnmarshal(*unmarshalConfig) {
	// No effect on unmarshal
}

// allowedClassesOption implements Option for class filtering
type allowedClassesOption struct {
	classes []string
}

func (o allowedClassesOption) applyMarshal(*marshalConfig) {
	// No effect on marshal
}

func (o allowedClassesOption) applyUnmarshal(cfg *unmarshalConfig) {
	if o.classes == nil {
		cfg.allowAll = false
		cfg.allowedClasses = nil
		return
	}
	cfg.allowAll = false
	allowed := make(map[string]bool)
	for _, c := range o.classes {
		allowed[c] = true
	}
	cfg.allowedClasses = allowed
}

// WithMaxDepth limits nesting depth for both Marshal and Unmarshal
// For Marshal: 0 = unlimited (default)
// For Unmarshal: 0 will use PHP default of 4096
func WithMaxDepth(depth int) Option {
	if depth < 0 {
		depth = 0
	}
	return maxDepthOption{depth: depth}
}

// WithStrictPHP ensures output matches PHP rules (default true)
func WithStrictPHP(strict bool) Option {
	return strictPHPOption{strict: strict}
}

// WithAllowedClasses restricts which PHP classes can be un-serialized
// If classes = nil, it disables object un-serialization (like PHP allowed_classes = false)
func WithAllowedClasses(classes []string) Option {
	return allowedClassesOption{classes: classes}
}

// Marshal converts a Go value to PHP serialized format
func Marshal(value interface{}, options ...Option) (string, error) {
	config := &marshalConfig{
		phpStrict: true,
		maxDepth:  0, // 0 = unlimited (PHP serialize has no max_depth)
	}
	for _, opt := range options {
		opt.applyMarshal(config)
	}

	var buf bytes.Buffer
	buf.Grow(256) // Pre-allocate reasonable size
	err := marshalValue(&buf, value, config, 0)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MarshalObject serializes a PHPObject
func MarshalObject(obj PHPObject, options ...Option) (string, error) {
	config := &marshalConfig{
		phpStrict: true,
		maxDepth:  0,
	}
	for _, opt := range options {
		opt.applyMarshal(config)
	}

	var buf bytes.Buffer
	buf.Grow(256)
	err := marshalObject(&buf, obj, config, 0)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Unmarshal converts PHP serialized data to Go values
func Unmarshal(data string, options ...Option) (interface{}, error) {
	config := &unmarshalConfig{
		allowAll: true, // PHP default = all classes allowed
		maxDepth: 4096, // PHP default max depth
	}
	for _, opt := range options {
		opt.applyUnmarshal(config)
	}

	reader := &stringReader{data: data, pos: 0}
	return unmarshalValue(reader, config, 0)
}

// stringReader helps to parse serialized data
type stringReader struct {
	data string
	pos  int
}

func (r *stringReader) read() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("unexpected end of data at position %d", r.pos)
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

func (r *stringReader) peek() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("unexpected end of data at position %d", r.pos)
	}
	return r.data[r.pos], nil
}

func (r *stringReader) readUntil(delim byte) (string, error) {
	start := r.pos
	for r.pos < len(r.data) {
		if r.data[r.pos] == delim {
			result := r.data[start:r.pos]
			r.pos++ // skip delimiter
			return result, nil
		}
		r.pos++
	}
	return "", fmt.Errorf("delimiter '%c' not found after position %d", delim, start)
}

func (r *stringReader) readBytes(n int) (string, error) {
	if r.pos+n > len(r.data) {
		return "", fmt.Errorf("not enough data at position %d: need %d bytes, have %d", r.pos, n, len(r.data)-r.pos)
	}
	result := r.data[r.pos : r.pos+n]
	r.pos += n
	return result, nil
}

// marshalValue serializes any Go value
func marshalValue(buf *bytes.Buffer, value interface{}, cfg *marshalConfig, depth int) error {
	if cfg.maxDepth > 0 && depth >= cfg.maxDepth {
		return fmt.Errorf("exceeded max depth %d", cfg.maxDepth)
	}

	if value == nil {
		buf.WriteString("N;")
		return nil
	}

	v := reflect.ValueOf(value)

	// Check for circular references in pointers
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Map || v.Kind() == reflect.Slice {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				buf.WriteString("N;")
				return nil
			}
			// Dereference pointer
			v = v.Elem()
		}
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("b:1;")
		} else {
			buf.WriteString("b:0;")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(fmt.Sprintf("i:%d;", v.Int()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := v.Uint()
		if cfg.phpStrict {
			if u > math.MaxInt64 {
				return fmt.Errorf("uint %d exceeds PHP int range", u)
			}
			buf.WriteString(fmt.Sprintf("i:%d;", int64(u)))
		} else {
			// Go-native: keep uint64 as-is in the serialized form.
			buf.WriteString(fmt.Sprintf("u:%d;", u))
		}

	case reflect.Float32, reflect.Float64:
		f := v.Float()
		// Handle special float cases like PHP does
		if math.IsNaN(f) {
			buf.WriteString("d:NAN;")
		} else if math.IsInf(f, 1) {
			buf.WriteString("d:INF;")
		} else if math.IsInf(f, -1) {
			buf.WriteString("d:-INF;")
		} else {
			buf.WriteString("d:" + strconv.FormatFloat(f, 'f', -1, 64) + ";")
		}

	case reflect.String:
		str := v.String()
		// PHP serialization uses byte length, not character count
		byteLen := len(str)
		buf.WriteString(fmt.Sprintf("s:%d:\"%s\";", byteLen, str))

	case reflect.Slice, reflect.Array:
		length := v.Len()
		buf.WriteString(fmt.Sprintf("a:%d:{", length))
		for i := 0; i < length; i++ {
			// Serialize index
			buf.WriteString(fmt.Sprintf("i:%d;", i))
			// Serialize value with incremented depth
			if err := marshalValue(buf, v.Index(i).Interface(), cfg, depth+1); err != nil {
				return err
			}
		}
		buf.WriteString("}")

	case reflect.Map:
		length := v.Len()
		buf.WriteString(fmt.Sprintf("a:%d:{", length))

		keys := v.MapKeys()

		for _, key := range keys {
			switch key.Kind() {
			case reflect.String:
				keyStr := key.String()
				buf.WriteString(fmt.Sprintf("s:%d:\"%s\";", len(keyStr), keyStr))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				buf.WriteString(fmt.Sprintf("i:%d;", key.Int()))
			default:
				return fmt.Errorf("cannot serialize map with key type %v", key.Kind())
			}

			if err := marshalValue(buf, v.MapIndex(key).Interface(), cfg, depth+1); err != nil {
				return err
			}
		}
		buf.WriteString("}")

	case reflect.Struct:
		// Check if it's a PHPObject
		if obj, ok := value.(PHPObject); ok {
			return marshalObject(buf, obj, cfg, depth)
		}
		// For other structs, convert to map
		return fmt.Errorf("cannot serialize struct type %T directly, use PHPObject or convert to map", value)

	case reflect.Ptr:
		if v.IsNil() {
			buf.WriteString("N;")
			return nil
		}
		return marshalValue(buf, v.Elem().Interface(), cfg, depth)

	default:
		return fmt.Errorf("cannot serialize type %s", v.Kind())
	}

	return nil
}

// marshalObject serializes a PHPObject
func marshalObject(buf *bytes.Buffer, obj PHPObject, cfg *marshalConfig, depth int) error {
	if cfg.maxDepth > 0 && depth >= cfg.maxDepth {
		return fmt.Errorf("exceeded max depth %d", cfg.maxDepth)
	}

	classNameLen := len(obj.ClassName)
	propCount := len(obj.Properties)

	buf.WriteString(fmt.Sprintf("O:%d:\"%s\":%d:{", classNameLen, obj.ClassName, propCount))

	for key, value := range obj.Properties {
		// Serialize property name
		buf.WriteString(fmt.Sprintf("s:%d:\"%s\";", len(key), key))
		// Serialize property value with incremented depth
		if err := marshalValue(buf, value, cfg, depth+1); err != nil {
			return err
		}
	}

	buf.WriteString("}")
	return nil
}

// unmarshalValue un-serializes a single value
func unmarshalValue(r *stringReader, cfg *unmarshalConfig, depth int) (interface{}, error) {
	if cfg.maxDepth > 0 && depth >= cfg.maxDepth {
		return nil, fmt.Errorf("exceeded max depth %d at position %d", cfg.maxDepth, r.pos)
	}

	typeChar, err := r.read()
	if err != nil {
		return nil, err
	}

	// Expect ':' after type (except for N)
	if typeChar != 'N' {
		colon, err := r.read()
		if err != nil {
			return nil, err
		}
		if colon != ':' {
			return nil, fmt.Errorf("at position %d: expected ':' after type '%c', got '%c'", r.pos-1, typeChar, colon)
		}
	}

	switch typeChar {
	case 'N': // NULL
		semicolon, err := r.read()
		if err != nil {
			return nil, err
		}
		if semicolon != ';' {
			return nil, fmt.Errorf("at position %d: expected ';' after NULL, got '%c'", r.pos-1, semicolon)
		}
		return nil, nil

	case 'b': // Boolean
		valStr, err := r.readUntil(';')
		if err != nil {
			return nil, err
		}
		return valStr == "1", nil

	case 'i': // Integer
		valStr, err := r.readUntil(';')
		if err != nil {
			return nil, err
		}
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("at position %d: invalid integer: %s", r.pos, valStr)
		}
		return val, nil

	case 'd': // Double/Float
		valStr, err := r.readUntil(';')
		if err != nil {
			return nil, err
		}
		// Handle special cases
		switch valStr {
		case "NAN":
			return math.NaN(), nil
		case "INF":
			return math.Inf(1), nil
		case "-INF":
			return math.Inf(-1), nil
		}
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("at position %d: invalid float: %s", r.pos, valStr)
		}
		return val, nil

	case 's': // String
		lenStr, err := r.readUntil(':')
		if err != nil {
			return nil, err
		}
		length, err := strconv.Atoi(lenStr)
		if err != nil {
			return nil, fmt.Errorf("at position %d: invalid string length: %s", r.pos, lenStr)
		}

		// Validate string length
		if length < 0 {
			return nil, fmt.Errorf("at position %d: negative string length: %d", r.pos, length)
		}

		// Read opening quote
		quote, err := r.read()
		if err != nil {
			return nil, err
		}
		if quote != '"' {
			return nil, fmt.Errorf("at position %d: expected '\"' before string, got '%c'", r.pos-1, quote)
		}

		// Read string bytes (not characters)
		str, err := r.readBytes(length)
		if err != nil {
			return nil, err
		}

		// Read closing quote
		quote, err = r.read()
		if err != nil {
			return nil, err
		}
		if quote != '"' {
			return nil, fmt.Errorf("at position %d: expected '\"' after string, got '%c'", r.pos-1, quote)
		}

		// Read semicolon
		semicolon, err := r.read()
		if err != nil {
			return nil, err
		}
		if semicolon != ';' {
			return nil, fmt.Errorf("at position %d: expected ';' after string, got '%c'", r.pos-1, semicolon)
		}

		return str, nil

	case 'a': // Array
		countStr, err := r.readUntil(':')
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return nil, fmt.Errorf("at position %d: invalid array count: %s", r.pos, countStr)
		}

		// Validate array size
		if count < 0 {
			return nil, fmt.Errorf("at position %d: negative array count: %d", r.pos, count)
		}

		// Read opening brace
		brace, err := r.read()
		if err != nil {
			return nil, err
		}
		if brace != '{' {
			return nil, fmt.Errorf("at position %d: expected '{' for array, got '%c'", r.pos-1, brace)
		}

		// Check if it's an indexed array (all keys are sequential integers starting from 0)
		isIndexed := true
		tempMap := make(map[string]interface{})
		indices := make([]int, 0, count)

		for i := 0; i < count; i++ {
			// Read key with incremented depth
			key, err := unmarshalValue(r, cfg, depth+1)
			if err != nil {
				return nil, err
			}

			// Read value with incremented depth
			value, err := unmarshalValue(r, cfg, depth+1)
			if err != nil {
				return nil, err
			}

			// Check if key is an integer
			switch k := key.(type) {
			case int64:
				indices = append(indices, int(k))
				tempMap[strconv.Itoa(int(k))] = value
			case string:
				isIndexed = false
				tempMap[k] = value
			default:
				keyStr := fmt.Sprintf("%v", k)
				tempMap[keyStr] = value
				isIndexed = false
			}
		}

		// Read closing brace
		brace, err = r.read()
		if err != nil {
			return nil, err
		}
		if brace != '}' {
			return nil, fmt.Errorf("at position %d: expected '}' for array, got '%c'", r.pos-1, brace)
		}

		// If it's an indexed array with sequential keys, return a slice
		if isIndexed && len(indices) > 0 {
			// Check if indices are sequential starting from 0
			sequential := true
			for i, idx := range indices {
				if idx != i {
					sequential = false
					break
				}
			}

			if sequential {
				result := make([]interface{}, len(indices))
				for i := 0; i < len(indices); i++ {
					result[i] = tempMap[strconv.Itoa(i)]
				}
				return result, nil
			}
		}

		// Otherwise, return a map
		if len(tempMap) == 0 {
			return make(map[string]interface{}), nil
		}
		return tempMap, nil

	case 'O': // Object
		classLenStr, err := r.readUntil(':')
		if err != nil {
			return nil, err
		}
		classLen, err := strconv.Atoi(classLenStr)
		if err != nil {
			return nil, fmt.Errorf("at position %d: invalid class name length: %s", r.pos, classLenStr)
		}

		if classLen < 0 {
			return nil, fmt.Errorf("at position %d: negative class name length: %d", r.pos, classLen)
		}

		// Read opening quote
		quote, err := r.read()
		if err != nil {
			return nil, err
		}
		if quote != '"' {
			return nil, fmt.Errorf("at position %d: expected '\"' before class name, got '%c'", r.pos-1, quote)
		}

		// Read class name
		className, err := r.readBytes(classLen)
		if err != nil {
			return nil, err
		}
		if !cfg.allowAll {
			if cfg.allowedClasses == nil || !cfg.allowedClasses[className] {
				return nil, fmt.Errorf("at position %d: class %q not allowed", r.pos, className)
			}
		}

		// Read closing quote
		quote, err = r.read()
		if err != nil {
			return nil, err
		}
		if quote != '"' {
			return nil, fmt.Errorf("at position %d: expected '\"' after class name, got '%c'", r.pos-1, quote)
		}

		// Read colon
		colon, err := r.read()
		if err != nil {
			return nil, err
		}
		if colon != ':' {
			return nil, fmt.Errorf("at position %d: expected ':' after class name, got '%c'", r.pos-1, colon)
		}

		// Read property count
		propCountStr, err := r.readUntil(':')
		if err != nil {
			return nil, err
		}
		propCount, err := strconv.Atoi(propCountStr)
		if err != nil {
			return nil, fmt.Errorf("at position %d: invalid property count: %s", r.pos, propCountStr)
		}

		// Validate property count
		if propCount < 0 {
			return nil, fmt.Errorf("at position %d: negative property count: %d", r.pos, propCount)
		}

		// Read opening brace
		brace, err := r.read()
		if err != nil {
			return nil, err
		}
		if brace != '{' {
			return nil, fmt.Errorf("at position %d: expected '{' for object properties, got '%c'", r.pos-1, brace)
		}

		properties := make(map[string]interface{})
		for i := 0; i < propCount; i++ {
			// Read property name with incremented depth
			propName, err := unmarshalValue(r, cfg, depth+1)
			if err != nil {
				return nil, err
			}

			// Read property value with incremented depth
			propValue, err := unmarshalValue(r, cfg, depth+1)
			if err != nil {
				return nil, err
			}

			if name, ok := propName.(string); ok {
				// Remove visibility prefixes if present (PHP private/protected properties)
				// Private: \0ClassName\0propertyName
				// Protected: \0*\0propertyName
				if strings.Contains(name, "\x00") {
					parts := strings.Split(name, "\x00")
					if len(parts) > 0 {
						name = parts[len(parts)-1]
					}
				}
				properties[name] = propValue
			} else {
				properties[fmt.Sprintf("%v", propName)] = propValue
			}
		}

		// Read closing brace
		brace, err = r.read()
		if err != nil {
			return nil, err
		}
		if brace != '}' {
			return nil, fmt.Errorf("at position %d: expected '}' for object, got '%c'", r.pos-1, brace)
		}

		return PHPObject{
			ClassName:  className,
			Properties: properties,
		}, nil

	default:
		return nil, fmt.Errorf("at position %d: unknown type '%c'", r.pos-1, typeChar)
	}
}

// Helper functions for common use cases

// IsValidMarshaled checks if a string is valid PHP serialized data
func IsValidMarshaled(data string) bool {
	_, err := Unmarshal(data)
	return err == nil
}

// MustMarshal serializes a value and panics on error
func MustMarshal(value interface{}, options ...Option) string {
	result, err := Marshal(value, options...)
	if err != nil {
		panic(err)
	}
	return result
}

// MustUnmarshal un-serializes data and panics on error
func MustUnmarshal(data string, options ...Option) interface{} {
	result, err := Unmarshal(data, options...)
	if err != nil {
		panic(err)
	}
	return result
}

// GetStringLength returns the byte length of a string as PHP would calculate it
func GetStringLength(s string) int {
	return len(s)
}

// IsUTF8 checks if a string is valid UTF-8
func IsUTF8(s string) bool {
	return utf8.ValidString(s)
}
