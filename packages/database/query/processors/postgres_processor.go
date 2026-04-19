package processors

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/database/query"
)

// PostgresProcessor processes query results for PostgreSQL.
type PostgresProcessor struct{}

var _ query.Processor = (*PostgresProcessor)(nil)

// NewPostgresProcessor creates a new PostgreSQL processor.
func NewPostgresProcessor() *PostgresProcessor { return &PostgresProcessor{} }

func (p *PostgresProcessor) ProcessSelect(_ *query.Builder, results []map[string]any) []map[string]any {
	return results
}

func (p *PostgresProcessor) ProcessInsertGetId(b *query.Builder, sql string, values []any, sequence string) (int64, error) {
	// PostgreSQL uses RETURNING clause, so the SQL already includes it.
	conn := b.GetConnection()
	rows, err := conn.Select(context.Background(), sql, values...)

	if err != nil {
		return 0, fmt.Errorf("processors: insert failed: %w", err)
	}

	if len(rows) == 0 {
		return 0, nil
	}

	return toInt64(rows[0][sequence]), nil
}

func (p *PostgresProcessor) ProcessColumnListing(results []map[string]any) []string {
	var columns []string

	for _, row := range results {
		if col, ok := row["column_name"]; ok {
			columns = append(columns, fmt.Sprintf("%v", col))
		}
	}

	return columns
}
