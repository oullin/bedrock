package events

import "time"

// QueryExecuted is dispatched after a SQL query has been executed.
type QueryExecuted struct {
	// SQL is the executed SQL string.
	SQL string
	// Bindings are the query parameter bindings.
	Bindings []any
	// Duration is how long the query took to execute.
	Duration time.Duration
	// ConnectionName is the name of the connection that ran the query.
	ConnectionName string
}
