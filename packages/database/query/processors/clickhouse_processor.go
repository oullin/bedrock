package processors

import (
	"errors"
	"fmt"

	"github.com/bedrock/packages/database/query"
)

// ClickHouseProcessor processes query results for ClickHouse.
type ClickHouseProcessor struct{}

var _ query.Processor = (*ClickHouseProcessor)(nil)

// NewClickHouseProcessor creates a new ClickHouse processor.
func NewClickHouseProcessor() *ClickHouseProcessor { return &ClickHouseProcessor{} }

func (p *ClickHouseProcessor) ProcessSelect(_ *query.Builder, results []map[string]any) []map[string]any {
	return results
}

// ProcessInsertGetId always returns an error: ClickHouse has no auto-increment
// columns and no last-insert-id concept. Callers must assign IDs explicitly
// (e.g. snowflake, UUID, generateUUIDv4()) before inserting.
func (p *ClickHouseProcessor) ProcessInsertGetId(_ *query.Builder, _ string, _ []any, _ string) (int64, error) {
	return 0, errors.New("clickhouse: InsertGetId is not supported (no auto-increment); assign IDs explicitly")
}

func (p *ClickHouseProcessor) ProcessColumnListing(results []map[string]any) []string {
	var columns []string

	for _, row := range results {
		if col, ok := row["name"]; ok {
			columns = append(columns, fmt.Sprintf("%v", col))

			continue
		}

		if col, ok := row["column_name"]; ok {
			columns = append(columns, fmt.Sprintf("%v", col))
		}
	}

	return columns
}
