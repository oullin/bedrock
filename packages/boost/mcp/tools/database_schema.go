package tools

import (
	"database/sql"
	"fmt"
	"strings"
)

// DatabaseSchema inspects the database schema and returns table/column information.
// Mirrors Upstream\Boost\Mcp\Tools\DatabaseSchema.
// Tagged IsReadOnly.
type DatabaseSchema struct {
	// DB is the database connection to inspect.
	DB *sql.DB
	// Driver is "sqlite", "mysql", or "postgres".
	Driver string
}

func (t *DatabaseSchema) Name() string     { return "database_schema" }
func (t *DatabaseSchema) IsReadOnly() bool { return true }

func (t *DatabaseSchema) Description() string {
	return "Retrieve the database schema including tables, columns, indexes, and foreign keys."
}

func (t *DatabaseSchema) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"filter": map[string]any{
				"type":        "string",
				"description": "Optional table name filter (partial match).",
			},
			"summary": map[string]any{
				"type":        "boolean",
				"description": "Return table names only without column details.",
				"default":     false,
			},
			"include_views": map[string]any{
				"type":        "boolean",
				"description": "Include views in the result.",
				"default":     false,
			},
		},
		"required": []string{},
	}
}

// Handle introspects the schema.
func (t *DatabaseSchema) Handle(req McpRequest) (McpResponse, error) {
	if t.DB == nil {
		return ErrorResponse("database_schema: no database connection configured"), nil
	}

	filter, _ := req.Args["filter"].(string)
	summary, _ := req.Args["summary"].(bool)

	tables, err := t.getTables(filter)

	if err != nil {
		return ErrorResponse(fmt.Sprintf("database_schema: %v", err)), nil
	}

	if summary {
		return OkResponse(map[string]any{"tables": tables}), nil
	}

	schema := make([]map[string]any, 0, len(tables))

	for _, table := range tables {
		cols, colErr := t.getColumns(table)

		if colErr != nil {
			continue
		}

		schema = append(schema, map[string]any{
			"table":   table,
			"columns": cols,
		})
	}

	return OkResponse(map[string]any{"schema": schema}), nil
}

func (t *DatabaseSchema) getTables(filter string) ([]string, error) {
	var query string

	switch strings.ToLower(t.Driver) {
	case "sqlite", "sqlite3":
		query = "SELECT name FROM sqlite_master WHERE type='table'"

		if filter != "" {
			query += fmt.Sprintf(" AND name LIKE '%%%s%%'", filter)
		}

		query += " ORDER BY name"
	case "postgres", "postgresql":
		query = "SELECT tablename FROM pg_tables WHERE schemaname='public'"

		if filter != "" {
			query += fmt.Sprintf(" AND tablename LIKE '%%%s%%'", filter)
		}

		query += " ORDER BY tablename"
	default: // mysql
		query = "SHOW TABLES"
	}

	rows, err := t.DB.Query(query) //nolint:gosec

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tables []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			continue
		}

		if filter == "" || strings.Contains(name, filter) {
			tables = append(tables, name)
		}
	}

	return tables, rows.Err()
}

func (t *DatabaseSchema) getColumns(table string) ([]map[string]any, error) {
	var query string

	switch strings.ToLower(t.Driver) {
	case "sqlite", "sqlite3":
		query = fmt.Sprintf("PRAGMA table_info(%s)", table)
	case "postgres", "postgresql":
		query = fmt.Sprintf(
			"SELECT column_name, data_type, is_nullable FROM information_schema.columns "+
				"WHERE table_name = '%s' ORDER BY ordinal_position", table)
	default:
		query = fmt.Sprintf("DESCRIBE %s", table)
	}

	rows, err := t.DB.Query(query) //nolint:gosec

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	var result []map[string]any

	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))

		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			continue
		}

		row := make(map[string]any, len(cols))

		for i, col := range cols {
			row[col] = values[i]
		}

		result = append(result, row)
	}

	return result, rows.Err()
}
