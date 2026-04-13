package jsonx

// StringType represents a JSON Schema string type.
type StringType struct {
	TypeBuilder[StringType]
	minLength *int
	maxLength *int
	pattern   *string
	format    *string
}

// newStringType creates a new StringType instance.
func newStringType() *StringType {
	t := &StringType{}
	t.self = t

	return t
}

// Min sets the minimum length (inclusive).
func (t *StringType) Min(value int) *StringType {
	t.minLength = &value

	return t
}

// Max sets the maximum length (inclusive).
func (t *StringType) Max(value int) *StringType {
	t.maxLength = &value

	return t
}

// Pattern sets the regular expression pattern the value must satisfy.
func (t *StringType) Pattern(value string) *StringType {
	t.pattern = &value

	return t
}

// Format sets the format of the string.
func (t *StringType) Format(value string) *StringType {
	t.format = &value

	return t
}

// Default sets the default value for the string type.
func (t *StringType) Default(value string) *StringType {
	t.defaultVal = value

	return t
}
