package tools

import (
	"database/sql"
	"fmt"
	"strings"
)

// DatabaseQuery executes a read-only SQL query and returns the results.
// Mirrors Upstream\Boost\Mcp\Tools\DatabaseQuery.
// Tagged IsReadOnly.
type DatabaseQuery struct {
	// DB is the database connection to execute queries against.
	DB *sql.DB
}

func (t *DatabaseQuery) Name() string     { return "database_query" }
func (t *DatabaseQuery) IsReadOnly() bool { return true }

func (t *DatabaseQuery) Description() string {
	return "Execute a read-only SQL query and return the results as JSON. " +
		"Only SELECT statements and CTEs are permitted."
}

func (t *DatabaseQuery) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query": map[string]any{
				"type":        "string",
				"description": "The SQL SELECT query to execute.",
			},
		},
		"required": []string{"query"},
	}
}

// Handle executes the SQL query and returns rows as a list of maps.
func (t *DatabaseQuery) Handle(req McpRequest) (McpResponse, error) {
	if t.DB == nil {
		return ErrorResponse("database_query: no database connection configured"), nil
	}

	query, _ := req.Args["query"].(string)

	if query == "" {
		return ErrorResponse("database_query: query argument is required"), nil
	}

	if !isSelectQuery(query) {
		return ErrorResponse("database_query: only SELECT queries are permitted"), nil
	}

	rows, err := t.DB.Query(query) //nolint:gosec

	if err != nil {
		return ErrorResponse(fmt.Sprintf("database_query: %v", err)), nil
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		return ErrorResponse(fmt.Sprintf("database_query: columns: %v", err)), nil
	}

	var results []map[string]any

	for rows.Next() {
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(cols))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return ErrorResponse(fmt.Sprintf("database_query: scan: %v", err)), nil
		}

		row := make(map[string]any, len(cols))

		for i, col := range cols {
			row[col] = values[i]
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return ErrorResponse(fmt.Sprintf("database_query: rows: %v", err)), nil
	}

	return OkResponse(map[string]any{"rows": results, "count": len(results)}), nil
}

// isSelectQuery performs a simple check that the query is read-only.
func isSelectQuery(query string) bool {
	upper := strings.TrimSpace(strings.ToUpper(query))

	return strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "WITH") ||
		strings.HasPrefix(upper, "EXPLAIN")
}
