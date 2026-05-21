package queue

import "strings"

// ParseQueue splits a command-line queue argument of the form
// "connection:queue" into its two halves, falling back to
// defaultConnection and the literal "default" for missing parts.
//
// used by queue:work / queue:retry / queue:pause and friends.
//
//	ParseQueue("",                       "redis")    → ("redis",    "default")
//	ParseQueue("emails",                 "redis")    → ("redis",    "emails")
//	ParseQueue("database:notifications", "redis")    → ("database", "notifications")
//	ParseQueue("redis:foo:bar",          "sqs")      → ("redis",    "foo:bar")
func ParseQueue(raw, defaultConnection string) (connection, queue string) {
	parts := strings.SplitN(raw, ":", 2)

	if len(parts) == 2 {
		connection, queue = parts[0], parts[1]
	} else {
		queue = parts[0]
	}

	if connection == "" {
		connection = defaultConnection
	}

	if queue == "" {
		queue = "default"
	}

	return connection, queue
}
