package engines

import (
	"context"
	"fmt"
	"strings"

	dbcontract "github.com/bedrock/packages/contracts/database"
	contract "github.com/bedrock/packages/contracts/scout"
)

// DatabaseEngine performs full-text search using the database's native
// capabilities. It leverages the query builder's WhereFullText method
// which compiles to MATCH/AGAINST (MySQL), to_tsvector (PostgreSQL),
// or LIKE (SQLite) depending on the grammar.
//
// This mirrors Laravel's Scout DatabaseEngine.
type DatabaseEngine struct {
	resolver   dbcontract.ConnectionResolver
	softDelete bool
}

// Compile-time interface checks.

// NewDatabaseEngine creates a new DatabaseEngine with the given connection resolver.

// The database engine does not maintain a separate index.
// Data is already in the database tables.

// The database engine does not maintain a separate index.
// Deleting from the database is handled by the ORM.

// Extract IDs from results to fetch full models.

// For the database engine, the rows already contain full data.
// We return them as-is via the model mapper.

// The database engine does not maintain a separate index to flush.

// Database tables serve as the index; table creation is handled by migrations.

// Database tables serve as the index; table deletion is handled by migrations.

// PaginateUsingDatabase performs paginated search using database queries directly.

// SimplePaginateUsingDatabase performs simple pagination using database queries.

// performSearch builds and executes the full-text search query.

// Build the search SQL.

// Add full-text search condition if query is not empty.

// Apply where constraints.

// Apply whereIn constraints.

// Apply whereNotIn constraints.

// Apply soft delete constraint.

// Apply ordering.

// Get count first for pagination.

// Apply limit and offset for pagination.

// Execute the search query.

// getSearchColumns returns the columns to search. It checks for fulltext
// and prefix columns from builder options, falling back to all columns
// from the model's searchable array.

// Check for explicitly configured columns.

// Fall back to all keys from the searchable array.

// buildFullTextClause builds a database-driver-specific full-text search clause.

// SQLite and others: fall back to LIKE.

// DatabaseResult holds the results of a database engine search.
type DatabaseResult struct {
	Rows       []map[string]any
	TotalCount int64
	KeyName    string
	Models     []contract.Searchable
}

var _ contract.Engine = (*DatabaseEngine)(nil)
var _ contract.PaginatesUsingDatabase = (*DatabaseEngine)(nil)

func NewDatabaseEngine(resolver dbcontract.ConnectionResolver, softDelete ...bool) *DatabaseEngine {
	sd := false

	if len(softDelete) > 0 {
		sd = softDelete[0]
	}

	return &DatabaseEngine{
		resolver:   resolver,
		softDelete: sd,
	}
}

func (e *DatabaseEngine) Update(ctx context.Context, models []contract.Searchable) error {

	return nil
}

func (e *DatabaseEngine) Delete(ctx context.Context, models []contract.Searchable) error {

	return nil
}

func (e *DatabaseEngine) Search(ctx context.Context, builder contract.SearchBuilder) (any, error) {
	return e.performSearch(ctx, builder, 0, 0)
}

func (e *DatabaseEngine) Paginate(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	return e.performSearch(ctx, builder, perPage, page)
}

func (e *DatabaseEngine) MapIds(results any) []any {
	dr, ok := results.(*DatabaseResult)

	if !ok {
		return []any{}
	}

	ids := make([]any, len(dr.Rows))
	keyName := dr.KeyName

	for i, row := range dr.Rows {
		ids[i] = row[keyName]
	}

	return ids
}

func (e *DatabaseEngine) Map(ctx context.Context, results any, model contract.Searchable) ([]contract.Searchable, error) {
	dr, ok := results.(*DatabaseResult)

	if !ok {
		return []contract.Searchable{}, nil
	}

	if len(dr.Rows) == 0 {
		return []contract.Searchable{}, nil
	}

	ids := e.MapIds(results)

	if len(ids) == 0 {
		return []contract.Searchable{}, nil
	}

	if dr.Models != nil {
		return dr.Models, nil
	}

	return []contract.Searchable{}, nil
}

func (e *DatabaseEngine) LazyMap(ctx context.Context, results any, model contract.Searchable) func(yield func(contract.Searchable) bool) {
	models, _ := e.Map(ctx, results, model)

	return func(yield func(contract.Searchable) bool) {
		for _, m := range models {
			if !yield(m) {
				return
			}
		}
	}
}

func (e *DatabaseEngine) GetTotalCount(results any) int64 {
	dr, ok := results.(*DatabaseResult)

	if !ok {
		return 0
	}

	return dr.TotalCount
}

func (e *DatabaseEngine) Flush(ctx context.Context, model contract.Searchable) error {

	return nil
}

func (e *DatabaseEngine) CreateIndex(_ context.Context, _ string, _ map[string]any) error {

	return nil
}

func (e *DatabaseEngine) DeleteIndex(_ context.Context, _ string) error {

	return nil
}

func (e *DatabaseEngine) PaginateUsingDatabase(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	return e.performSearch(ctx, builder, perPage, page)
}

func (e *DatabaseEngine) SimplePaginateUsingDatabase(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	return e.performSearch(ctx, builder, perPage+1, page)
}

func (e *DatabaseEngine) performSearch(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (*DatabaseResult, error) {
	model := builder.GetModel()

	if model == nil {
		return nil, fmt.Errorf("scout: database engine requires a model")
	}

	conn, err := e.resolver.Connection(ctx, model.GetConnectionName())

	if err != nil {
		return nil, fmt.Errorf("scout: failed to resolve connection: %w", err)
	}

	table := model.GetTable()
	keyName := model.GetKeyName()
	query := builder.GetQuery()

	var sql strings.Builder

	var bindings []any

	sql.WriteString("select * from ")
	sql.WriteString(table)

	whereAdded := false

	if query != "" {
		columns := e.getSearchColumns(builder)

		if len(columns) > 0 {
			driver := conn.GetDriverName()
			fullTextSQL, fullTextBindings := e.buildFullTextClause(driver, columns, query)
			sql.WriteString(" where ")
			sql.WriteString(fullTextSQL)
			bindings = append(bindings, fullTextBindings...)
			whereAdded = true
		}
	}

	for key, value := range builder.GetWheres() {
		if whereAdded {
			sql.WriteString(" and ")
		} else {
			sql.WriteString(" where ")
			whereAdded = true
		}

		sql.WriteString(key)
		sql.WriteString(" = ?")
		bindings = append(bindings, value)
	}

	for key, values := range builder.GetWhereIns() {
		if len(values) == 0 {
			continue
		}

		if whereAdded {
			sql.WriteString(" and ")
		} else {
			sql.WriteString(" where ")
			whereAdded = true
		}

		sql.WriteString(key)
		sql.WriteString(" in (")
		placeholders := make([]string, len(values))

		for i, v := range values {
			placeholders[i] = "?"
			bindings = append(bindings, v)
		}

		sql.WriteString(strings.Join(placeholders, ", "))
		sql.WriteString(")")
	}

	for key, values := range builder.GetWhereNotIns() {
		if len(values) == 0 {
			continue
		}

		if whereAdded {
			sql.WriteString(" and ")
		} else {
			sql.WriteString(" where ")
			whereAdded = true
		}

		sql.WriteString(key)
		sql.WriteString(" not in (")
		placeholders := make([]string, len(values))

		for i, v := range values {
			placeholders[i] = "?"
			bindings = append(bindings, v)
		}

		sql.WriteString(strings.Join(placeholders, ", "))
		sql.WriteString(")")
	}

	if e.softDelete && model.UsesSoftDelete() {
		if whereAdded {
			sql.WriteString(" and ")
		} else {
			sql.WriteString(" where ")
			whereAdded = true
		}

		sql.WriteString("deleted_at is null")
	}

	orders := builder.GetOrders()

	if len(orders) > 0 {
		sql.WriteString(" order by ")
		orderClauses := make([]string, len(orders))

		for i, o := range orders {
			orderClauses[i] = o.Column + " " + o.Direction
		}

		sql.WriteString(strings.Join(orderClauses, ", "))
	}

	countSQL := "select count(*) as aggregate from (" + sql.String() + ") as scout_count"
	countRow, err := conn.SelectOne(ctx, countSQL, bindings...)

	if err != nil {
		return nil, fmt.Errorf("scout: count query failed: %w", err)
	}

	totalCount := int64(0)

	if countRow != nil {
		if agg, ok := countRow["aggregate"]; ok {
			switch v := agg.(type) {
			case int64:
				totalCount = v
			case float64:
				totalCount = int64(v)
			case int:
				totalCount = int64(v)
			}
		}
	}

	if perPage > 0 {
		offset := 0

		if page > 1 {
			offset = (page - 1) * perPage
		}

		sql.WriteString(fmt.Sprintf(" limit %d offset %d", perPage, offset))
	} else if builder.GetLimit() > 0 {
		sql.WriteString(fmt.Sprintf(" limit %d", builder.GetLimit()))
	}

	rows, err := conn.Select(ctx, sql.String(), bindings...)

	if err != nil {
		return nil, fmt.Errorf("scout: search query failed: %w", err)
	}

	return &DatabaseResult{
		Rows:       rows,
		TotalCount: totalCount,
		KeyName:    keyName,
	}, nil
}

func (e *DatabaseEngine) getSearchColumns(builder contract.SearchBuilder) []string {
	opts := builder.GetOptions()

	if cols, ok := opts["__fulltext_columns"].([]string); ok && len(cols) > 0 {
		return cols
	}

	model := builder.GetModel()

	if model != nil {
		data := model.ToSearchableArray()
		columns := make([]string, 0, len(data))

		for key := range data {
			columns = append(columns, key)
		}

		return columns
	}

	return nil
}

func (e *DatabaseEngine) buildFullTextClause(driver string, columns []string, query string) (string, []any) {
	colStr := strings.Join(columns, ", ")

	switch driver {
	case "mysql", "mariadb":
		return fmt.Sprintf("match (%s) against (? in natural language mode)", colStr), []any{query}
	case "pgsql":
		return fmt.Sprintf("to_tsvector(%s) @@ plainto_tsquery(?)", colStr), []any{query}
	default:

		var clauses []string

		var bindings []any
		pattern := "%" + query + "%"

		for _, col := range columns {
			clauses = append(clauses, fmt.Sprintf("%s like ?", col))
			bindings = append(bindings, pattern)
		}

		return "(" + strings.Join(clauses, " or ") + ")", bindings
	}
}
