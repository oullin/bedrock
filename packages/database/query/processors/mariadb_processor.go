package processors

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/database/query"
)

// MariaDBProcessor processes query results for MariaDB.
// MariaDB 10.5+ supports RETURNING, so InsertGetId reads from the
// result set rather than calling LAST_INSERT_ID().
type MariaDBProcessor struct{}

var _ query.Processor = (*MariaDBProcessor)(nil)

// NewMariaDBProcessor creates a new MariaDB processor.
func NewMariaDBProcessor() *MariaDBProcessor { return &MariaDBProcessor{} }

func (p *MariaDBProcessor) ProcessSelect(_ *query.Builder, results []map[string]any) []map[string]any {
	return results
}

func (p *MariaDBProcessor) ProcessInsertGetId(b *query.Builder, sql string, values []any, sequence string) (int64, error) {
	// MariaDB grammar generates INSERT ... RETURNING id, so we read
	// the returned row directly — same approach as PostgreSQL.
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

func (p *MariaDBProcessor) ProcessColumnListing(results []map[string]any) []string {
	var columns []string
	for _, row := range results {
		if col, ok := row["column_name"]; ok {
			columns = append(columns, fmt.Sprintf("%v", col))
		}
	}
	return columns
}
