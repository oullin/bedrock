package processors

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/database/query"
)

// MySQLProcessor processes query results for MySQL.
type MySQLProcessor struct{}

var _ query.Processor = (*MySQLProcessor)(nil)

// NewMySQLProcessor creates a new MySQL processor.
func NewMySQLProcessor() *MySQLProcessor { return &MySQLProcessor{} }

func (p *MySQLProcessor) ProcessSelect(_ *query.Builder, results []map[string]any) []map[string]any {
	return results
}

func (p *MySQLProcessor) ProcessInsertGetId(b *query.Builder, sql string, values []any, _ string) (int64, error) {
	conn := b.GetConnection()
	result, err := conn.Statement(context.Background(), sql, values...)
	if err != nil || !result {
		return 0, fmt.Errorf("processors: insert failed: %w", err)
	}

	type lastInsertIDer interface {
		DB() interface {
			QueryRowContext(ctx context.Context, query string, args ...any) interface{ Scan(dest ...any) error }
		}
	}

	// MySQL uses LAST_INSERT_ID() — retrieve via the connection.
	row, err := conn.Select(context.Background(), "SELECT LAST_INSERT_ID() as id")
	if err != nil || len(row) == 0 {
		return 0, err
	}
	return toInt64(row[0]["id"]), nil
}

func (p *MySQLProcessor) ProcessColumnListing(results []map[string]any) []string {
	var columns []string
	for _, row := range results {
		if col, ok := row["column_name"]; ok {
			columns = append(columns, fmt.Sprintf("%v", col))
		}
	}
	return columns
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case string:
		var i int64
		fmt.Sscanf(n, "%d", &i)
		return i
	default:
		return 0
	}
}
