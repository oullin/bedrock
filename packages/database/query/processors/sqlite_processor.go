package processors

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/database/query"
)

// SQLiteProcessor processes query results for SQLite.
type SQLiteProcessor struct{}

var _ query.Processor = (*SQLiteProcessor)(nil)

// NewSQLiteProcessor creates a new SQLite processor.
func NewSQLiteProcessor() *SQLiteProcessor { return &SQLiteProcessor{} }

func (p *SQLiteProcessor) ProcessSelect(_ *query.Builder, results []map[string]any) []map[string]any {
	return results
}

func (p *SQLiteProcessor) ProcessInsertGetId(b *query.Builder, sql string, values []any, _ string) (int64, error) {
	conn := b.GetConnection()
	result, err := conn.Statement(context.Background(), sql, values...)

	if err != nil || !result {
		return 0, fmt.Errorf("processors: insert failed: %w", err)
	}

	// SQLite uses last_insert_rowid().
	rows, err := conn.Select(context.Background(), "SELECT last_insert_rowid() as id")

	if err != nil || len(rows) == 0 {
		return 0, err
	}

	return toInt64(rows[0]["id"]), nil
}

func (p *SQLiteProcessor) ProcessColumnListing(results []map[string]any) []string {
	var columns []string

	for _, row := range results {
		if col, ok := row["name"]; ok {
			columns = append(columns, fmt.Sprintf("%v", col))
		}
	}

	return columns
}
