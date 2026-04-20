package redis

import "strings"

// HasHashTag reports whether key contains a valid Redis Cluster hash tag.
//
// Redis uses the first opening brace and only treats it as a hash tag when
// a later closing brace leaves at least one byte inside the tag.
func HasHashTag(key string) bool {
	open := strings.IndexByte(key, '{')

	if open < 0 {
		return false
	}

	close := strings.IndexByte(key[open+1:], '}')

	return close > 0
}
