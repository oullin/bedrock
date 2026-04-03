package kv

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Query encodes the map as a URL query string.
func Query(items map[string]any) string {
	params := url.Values{}

	for k, v := range items {
		params.Set(k, fmt.Sprint(v))
	}

	return params.Encode()
}

// ToCssClasses returns a space-separated string of CSS class names
// whose corresponding boolean values are true. Classes are sorted alphabetically.
func ToCssClasses(classes map[string]bool) string {
	result := make([]string, 0)

	for class, condition := range classes {
		if condition {
			result = append(result, class)
		}
	}

	sort.Strings(result)

	return strings.Join(result, " ")
}

// ToCssStyles returns a space-separated string of CSS style declarations
// whose corresponding boolean values are true. A trailing semicolon is appended
// to each style if not already present. Styles are sorted alphabetically.
func ToCssStyles(styles map[string]bool) string {
	result := make([]string, 0)

	for style, condition := range styles {
		if condition {
			if !strings.HasSuffix(style, ";") {
				style += ";"
			}

			result = append(result, style)
		}
	}

	sort.Strings(result)

	return strings.Join(result, " ")
}
