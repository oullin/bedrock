package events

import "strings"

// matchesWildcard reports whether the given wildcard pattern matches the event
// name. Patterns use '*' as a wildcard that matches any single dot-separated
// segment. A trailing '*' matches all remaining segments.
//
// Examples:
//
//	matchesWildcard("order.*", "order.created")    => true
//	matchesWildcard("order.*", "order.shipped")    => true
//	matchesWildcard("*.created", "order.created")  => true
//	matchesWildcard("order.*", "user.created")     => false
func matchesWildcard(pattern, name string) bool {
	patternParts := strings.Split(pattern, ".")
	nameParts := strings.Split(name, ".")

	for i, pp := range patternParts {
		if pp == "*" {
			// A trailing wildcard matches all remaining segments.
			if i == len(patternParts)-1 {
				return true
			}

			// A non-trailing wildcard matches exactly one segment.
			if i >= len(nameParts) {
				return false
			}

			continue
		}

		if i >= len(nameParts) || pp != nameParts[i] {
			return false
		}
	}

	return len(patternParts) == len(nameParts)
}
