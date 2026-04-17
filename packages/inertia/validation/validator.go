package validation

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/bedrock/packages/inertia/protocol"
)

var (
	instance *validator.Validate
	once     sync.Once
)

func engine() *validator.Validate {
	once.Do(func() {
		instance = validator.New(validator.WithRequiredStructEnabled())

		instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")

			if strings.TrimSpace(name) == "" {
				name = fld.Tag.Get("form")
			}

			if strings.TrimSpace(name) == "" {
				return fld.Name
			}

			if idx := strings.IndexByte(name, ','); idx != -1 {
				name = name[:idx]
			}

			if name == "-" {
				return fld.Name
			}

			return name
		})
	})

	return instance
}

// Validate checks the given struct against its `validate` tags and returns
// Inertia-compatible validation errors keyed by json/form tag names.
// Returns nil when all rules pass.
func Validate(data any) protocol.ValidationErrors {
	err := engine().Struct(data)

	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return protocol.ValidationErrors{"_error": err.Error()}
	}

	result := make(protocol.ValidationErrors, len(validationErrors))

	for _, fe := range validationErrors {
		field := fe.Field()
		result[field] = message(fe)
	}

	return result
}

func message(fe validator.FieldError) string {
	label := humanize(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", label)
	case "email":
		return fmt.Sprintf("The %s field must be a valid email address.", label)
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters.", label, fe.Param())
	case "max":
		return fmt.Sprintf("The %s field must not exceed %s characters.", label, fe.Param())
	case "len":
		return fmt.Sprintf("The %s field must be exactly %s characters.", label, fe.Param())
	case "gte":
		return fmt.Sprintf("The %s field must be at least %s.", label, fe.Param())
	case "lte":
		return fmt.Sprintf("The %s field must be at most %s.", label, fe.Param())
	case "oneof":
		return fmt.Sprintf("The %s field must be one of: %s.", label, fe.Param())
	case "url":
		return fmt.Sprintf("The %s field must be a valid URL.", label)
	case "eqfield":
		return fmt.Sprintf("The %s field must match %s.", label, fe.Param())
	default:
		return fmt.Sprintf("The %s field is invalid.", label)
	}
}

// humanize converts a snake_case or camelCase field name to a readable label.
// e.g. "first_name" → "first name", "FirstName" → "first name"
func humanize(field string) string {
	var b strings.Builder

	for i, r := range field {
		if r == '_' {
			b.WriteByte(' ')

			continue
		}

		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(rune(field[i-1])) {
			b.WriteByte(' ')
		}

		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}
