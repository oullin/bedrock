package query

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	dbcontract "github.com/bedrock/packages/contracts/database"
)

// Binding category keys.
const (
	BindingSelect  = "select"
	BindingFrom    = "from"
	BindingJoin    = "join"
	BindingWhere   = "where"
	BindingGroupBy = "groupby"
	BindingHaving  = "having"
	BindingOrder   = "order"
	BindingUnion   = "union"
)

// Valid SQL operators.
var validOperators = map[string]bool{
	"=": true, "<": true, ">": true, "<=": true, ">=": true, "<>": true, "!=": true,
	"like": true, "like binary": true, "not like": true, "ilike": true,
	"&": true, "|": true, "^": true, "<<": true, ">>": true, "<=>": true,
	"rlike": true, "not rlike": true, "regexp": true, "not regexp": true,
	"~": true, "~*": true, "!~": true, "!~*": true, "similar to": true,
	"not similar to": true, "not ilike": true, "~~*": true, "!~~*": true,
}

// WhereType identifies the kind of WHERE clause.
type WhereType string

const (
	WhereBasic       WhereType = "Basic"
	WhereColumn      WhereType = "Column"
	WhereIn          WhereType = "In"
	WhereNotIn       WhereType = "NotIn"
	WhereNull        WhereType = "Null"
	WhereNotNull     WhereType = "NotNull"
	WhereBetween     WhereType = "Between"
	WhereNotBetween  WhereType = "NotBetween"
	WhereDate        WhereType = "Date"
	WhereTime        WhereType = "Time"
	WhereDay         WhereType = "Day"
	WhereMonth       WhereType = "Month"
	WhereYear        WhereType = "Year"
	WhereRaw         WhereType = "Raw"
	WhereExists      WhereType = "Exists"
	WhereNotExists   WhereType = "NotExists"
	WhereNested      WhereType = "Nested"
	WhereSub         WhereType = "Sub"
	WhereLike        WhereType = "Like"
	WhereNotLike     WhereType = "NotLike"
	WhereJsonContains    WhereType = "JsonContains"
	WhereJsonLength      WhereType = "JsonLength"
	WhereFullText        WhereType = "Fulltext"
	WhereBetweenColumns  WhereType = "BetweenColumns"
	WhereRowValues       WhereType = "RowValues"
)

// WhereClause represents a single WHERE condition.
type WhereClause struct {
	Type     WhereType
	Column   string
	Operator string
	Value    any
	Values   []any
	Boolean  string // "and" or "or"
	Not      bool
	Query    *Builder // for nested/sub/exists wheres
	SQL      string   // for raw wheres
	Columns  []string // for row values, between columns
}

// OrderClause represents an ORDER BY clause.
type OrderClause struct {
	Column    string
	Direction string
	SQL       string // for raw orders
}

// HavingClause represents a HAVING clause.
type HavingClause struct {
	Type     string
	Column   string
	Operator string
	Value    any
	Boolean  string
	SQL      string
	Not      bool
	Values   []any
}

// UnionClause represents a UNION clause.
type UnionClause struct {
	Query *Builder
	All   bool
}

// JoinType identifies a JOIN kind.
type JoinType string

const (
	JoinInner   JoinType = "inner"
	JoinLeft    JoinType = "left"
	JoinRight   JoinType = "right"
	JoinCross   JoinType = "cross"
	JoinLateral JoinType = "lateral"
)

// IndexHint represents an index hint.
type IndexHint struct {
	Type  string // "use", "force", "ignore"
	Index string
}

// ConnectionInterface is the minimal connection surface the builder needs.
type ConnectionInterface interface {
	Select(ctx context.Context, query string, bindings ...any) ([]map[string]any, error)
	Insert(ctx context.Context, query string, bindings ...any) (bool, error)
	Update(ctx context.Context, query string, bindings ...any) (int64, error)
	Delete(ctx context.Context, query string, bindings ...any) (int64, error)
	Statement(ctx context.Context, query string, bindings ...any) (bool, error)
	AffectingStatement(ctx context.Context, query string, bindings ...any) (int64, error)
	Raw(value string) dbcontract.Expression
	GetTablePrefix() string
}

// Builder provides a fluent interface for building SQL queries.
type Builder struct {
	connection ConnectionInterface
	grammar    Grammar
	processor  Processor

	// Query components.
	columns       []any    // string or Expression
	distinct      bool
	distinctColumns []string
	from          string
	fromRaw       string
	joins         []*JoinClause
	wheres        []WhereClause
	groups        []string
	havings       []HavingClause
	orders        []OrderClause
	limit_        *int
	offset_       *int
	unions        []UnionClause
	unionLimit    *int
	unionOffset   *int
	unionOrders   []OrderClause
	lock          any
	indexHint     *IndexHint
	aggregate     *AggregateClause

	// Bindings keyed by category.
	bindings map[string][]any

	// Callbacks.
	afterQueryCallbacks []func([]map[string]any)
}

// AggregateClause holds the aggregate function and columns.
type AggregateClause struct {
	Function string
	Columns  []string
}

// NewBuilder creates a new query Builder.
func NewBuilder(connection ConnectionInterface, grammar Grammar, processor Processor) *Builder {
	return &Builder{
		connection: connection,
		grammar:    grammar,
		processor:  processor,
		bindings:   newBindings(),
	}
}

func newBindings() map[string][]any {
	return map[string][]any{
		BindingSelect:  {},
		BindingFrom:    {},
		BindingJoin:    {},
		BindingWhere:   {},
		BindingGroupBy: {},
		BindingHaving:  {},
		BindingOrder:   {},
		BindingUnion:   {},
	}
}

// NewQuery creates a fresh builder instance with the same connection/grammar.
func (b *Builder) NewQuery() *Builder {
	return NewBuilder(b.connection, b.grammar, b.processor)
}

// Clone creates a shallow copy of the builder.
func (b *Builder) Clone() *Builder {
	clone := *b
	clone.columns = append([]any(nil), b.columns...)
	clone.wheres = append([]WhereClause(nil), b.wheres...)
	clone.groups = append([]string(nil), b.groups...)
	clone.havings = append([]HavingClause(nil), b.havings...)
	clone.orders = append([]OrderClause(nil), b.orders...)
	clone.unions = append([]UnionClause(nil), b.unions...)
	clone.bindings = make(map[string][]any, len(b.bindings))
	for k, v := range b.bindings {
		clone.bindings[k] = append([]any(nil), v...)
	}
	if b.joins != nil {
		clone.joins = make([]*JoinClause, len(b.joins))
		copy(clone.joins, b.joins)
	}
	return &clone
}

// CloneWithout creates a clone that omits the given properties.
func (b *Builder) CloneWithout(properties ...string) *Builder {
	clone := b.Clone()
	for _, prop := range properties {
		switch prop {
		case "columns":
			clone.columns = nil
		case "wheres":
			clone.wheres = nil
			clone.bindings[BindingWhere] = nil
		case "orders":
			clone.orders = nil
			clone.bindings[BindingOrder] = nil
		case "limit":
			clone.limit_ = nil
		case "offset":
			clone.offset_ = nil
		case "unions":
			clone.unions = nil
			clone.bindings[BindingUnion] = nil
		case "groups":
			clone.groups = nil
			clone.bindings[BindingGroupBy] = nil
		case "havings":
			clone.havings = nil
			clone.bindings[BindingHaving] = nil
		case "joins":
			clone.joins = nil
			clone.bindings[BindingJoin] = nil
		}
	}
	return clone
}

// CloneWithoutBindings creates a clone that omits bindings for given categories.
func (b *Builder) CloneWithoutBindings(categories ...string) *Builder {
	clone := b.Clone()
	for _, cat := range categories {
		clone.bindings[cat] = nil
	}
	return clone
}

// From sets the table the query is targeting.
func (b *Builder) From(table string, as ...string) *Builder {
	if len(as) > 0 && as[0] != "" {
		b.from = table + " as " + as[0]
	} else {
		b.from = table
	}
	return b
}

// FromSub sets a subquery as the source table.
func (b *Builder) FromSub(query any, as string) *Builder {
	switch q := query.(type) {
	case *Builder:
		b.fromRaw = "(" + b.grammar.CompileSelect(q) + ") as " + b.grammar.Wrap(as)
		b.AddBinding(BindingFrom, q.GetBindings()...)
	case string:
		b.fromRaw = "(" + q + ") as " + b.grammar.Wrap(as)
	}
	return b
}

// FromRaw sets a raw FROM expression.
func (b *Builder) FromRaw(expression string, bindings ...any) *Builder {
	b.fromRaw = expression
	b.AddBinding(BindingFrom, bindings...)
	return b
}

// GetFrom returns the from table.
func (b *Builder) GetFrom() string {
	if b.fromRaw != "" {
		return b.fromRaw
	}
	return b.from
}

// Connection returns the underlying connection.
func (b *Builder) GetConnection() ConnectionInterface { return b.connection }

// Grammar returns the query grammar.
func (b *Builder) GetGrammar() Grammar { return b.grammar }

// Processor returns the result processor.
func (b *Builder) GetProcessor() Processor { return b.processor }

// GetColumns returns the selected columns.
func (b *Builder) GetColumns() []any { return b.columns }

// IsDistinct returns whether the query uses DISTINCT.
func (b *Builder) IsDistinct() bool { return b.distinct }

// GetDistinctColumns returns the distinct columns (if any).
func (b *Builder) GetDistinctColumns() []string { return b.distinctColumns }

// GetJoins returns the join clauses.
func (b *Builder) GetJoins() []*JoinClause { return b.joins }

// GetWheres returns the where clauses.
func (b *Builder) GetWheres() []WhereClause { return b.wheres }

// GetGroups returns the group by columns.
func (b *Builder) GetGroups() []string { return b.groups }

// GetHavings returns the having clauses.
func (b *Builder) GetHavings() []HavingClause { return b.havings }

// GetOrders returns the order by clauses.
func (b *Builder) GetOrders() []OrderClause { return b.orders }

// GetLimit returns the limit, or -1 if unset.
func (b *Builder) GetLimit() int {
	if b.limit_ == nil {
		return -1
	}
	return *b.limit_
}

// GetOffset returns the offset, or -1 if unset.
func (b *Builder) GetOffset() int {
	if b.offset_ == nil {
		return -1
	}
	return *b.offset_
}

// GetUnions returns the union clauses.
func (b *Builder) GetUnions() []UnionClause { return b.unions }

// GetUnionLimit returns the union limit.
func (b *Builder) GetUnionLimit() *int { return b.unionLimit }

// GetUnionOffset returns the union offset.
func (b *Builder) GetUnionOffset() *int { return b.unionOffset }

// GetUnionOrders returns the union order clauses.
func (b *Builder) GetUnionOrders() []OrderClause { return b.unionOrders }

// GetLock returns the lock value.
func (b *Builder) GetLock() any { return b.lock }

// GetIndexHint returns the index hint.
func (b *Builder) GetIndexHint() *IndexHint { return b.indexHint }

// GetAggregate returns the aggregate clause.
func (b *Builder) GetAggregate() *AggregateClause { return b.aggregate }

// SetAggregate sets the aggregate function and columns.
func (b *Builder) SetAggregate(function string, columns []string) *Builder {
	b.aggregate = &AggregateClause{Function: function, Columns: columns}
	return b
}

// ---- Bindings ----

// AddBinding adds values to a binding category.
func (b *Builder) AddBinding(category string, values ...any) *Builder {
	b.bindings[category] = append(b.bindings[category], values...)
	return b
}

// SetBindings replaces bindings for a category.
func (b *Builder) SetBindings(category string, values []any) *Builder {
	b.bindings[category] = values
	return b
}

// GetBindings returns all bindings in order.
func (b *Builder) GetBindings() []any {
	order := []string{
		BindingSelect, BindingFrom, BindingJoin, BindingWhere,
		BindingGroupBy, BindingHaving, BindingOrder, BindingUnion,
	}
	var all []any
	for _, key := range order {
		all = append(all, b.bindings[key]...)
	}
	return all
}

// GetRawBindings returns the raw bindings map.
func (b *Builder) GetRawBindings() map[string][]any { return b.bindings }

// ---- Terminal methods ----

// Get executes the SELECT query and returns all rows.
func (b *Builder) Get(ctx context.Context, columns ...string) ([]map[string]any, error) {
	if len(columns) > 0 {
		original := b.columns
		b.columns = stringsToAny(columns)
		defer func() { b.columns = original }()
	}

	compiled := b.grammar.CompileSelect(b)
	results, err := b.connection.Select(ctx, compiled, b.GetBindings()...)
	if err != nil {
		return nil, err
	}

	results = b.processor.ProcessSelect(b, results)

	for _, cb := range b.afterQueryCallbacks {
		cb(results)
	}

	return results, nil
}

// First returns the first row of the query.
func (b *Builder) First(ctx context.Context, columns ...string) (map[string]any, error) {
	original := b.limit_
	one := 1
	b.limit_ = &one
	defer func() { b.limit_ = original }()

	results, err := b.Get(ctx, columns...)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

// Value returns a single column value from the first row.
func (b *Builder) Value(ctx context.Context, column string) (any, error) {
	row, err := b.First(ctx, column)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	return row[column], nil
}

// Pluck returns a slice of values for a single column.
func (b *Builder) Pluck(ctx context.Context, column string) ([]any, error) {
	rows, err := b.Get(ctx, column)
	if err != nil {
		return nil, err
	}
	var values []any
	for _, row := range rows {
		values = append(values, row[column])
	}
	return values, nil
}

// PluckMap returns a map of key-value pairs from two columns.
func (b *Builder) PluckMap(ctx context.Context, column, key string) (map[any]any, error) {
	rows, err := b.Get(ctx, column, key)
	if err != nil {
		return nil, err
	}
	m := make(map[any]any, len(rows))
	for _, row := range rows {
		m[row[key]] = row[column]
	}
	return m, nil
}

// Chunk processes results in chunks of the given size.
func (b *Builder) Chunk(ctx context.Context, count int, fn func([]map[string]any, int) bool) error {
	page := 1
	for {
		results, err := b.ForPage(page, count).Get(ctx)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			break
		}
		if !fn(results, page) {
			break
		}
		if len(results) < count {
			break
		}
		page++
	}
	return nil
}

// ChunkByID processes results in chunks using an ID column for pagination.
func (b *Builder) ChunkByID(ctx context.Context, count int, column string, fn func([]map[string]any) bool) error {
	var lastID any

	for {
		clone := b.Clone()
		if lastID != nil {
			clone.Where(column, ">", lastID)
		}
		results, err := clone.OrderBy(column).Limit(count).Get(ctx)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			break
		}
		if !fn(results) {
			break
		}
		lastID = results[len(results)-1][column]
		if len(results) < count {
			break
		}
	}
	return nil
}

// ToSQL returns the compiled SQL and bindings.
func (b *Builder) ToSQL() (string, []any) {
	return b.grammar.CompileSelect(b), b.GetBindings()
}

// Dump returns the SQL and bindings as a formatted string.
func (b *Builder) Dump() string {
	sql, bindings := b.ToSQL()
	return fmt.Sprintf("SQL: %s\nBindings: %v", sql, bindings)
}

// AfterQuery registers a callback to run after query execution.
func (b *Builder) AfterQuery(fn func([]map[string]any)) *Builder {
	b.afterQueryCallbacks = append(b.afterQueryCallbacks, fn)
	return b
}

// ---- Utility ----

// UseIndex adds a USE INDEX hint.
func (b *Builder) UseIndex(index string) *Builder {
	b.indexHint = &IndexHint{Type: "use", Index: index}
	return b
}

// ForceIndex adds a FORCE INDEX hint.
func (b *Builder) ForceIndex(index string) *Builder {
	b.indexHint = &IndexHint{Type: "force", Index: index}
	return b
}

// IgnoreIndex adds an IGNORE INDEX hint.
func (b *Builder) IgnoreIndex(index string) *Builder {
	b.indexHint = &IndexHint{Type: "ignore", Index: index}
	return b
}

// LockForUpdate adds a "for update" lock to the query.
func (b *Builder) LockForUpdate() *Builder {
	b.lock = "for update"
	return b
}

// SharedLock adds a shared lock to the query.
func (b *Builder) SharedLock() *Builder {
	b.lock = "lock in share mode"
	return b
}

// Lock sets a custom lock expression.
func (b *Builder) Lock(value any) *Builder {
	b.lock = value
	return b
}

// When applies a callback conditionally.
func (b *Builder) When(condition bool, callback func(*Builder) *Builder, otherwise ...func(*Builder) *Builder) *Builder {
	if condition {
		return callback(b)
	}
	if len(otherwise) > 0 {
		return otherwise[0](b)
	}
	return b
}

// Unless applies a callback when the condition is false.
func (b *Builder) Unless(condition bool, callback func(*Builder) *Builder, otherwise ...func(*Builder) *Builder) *Builder {
	return b.When(!condition, callback, otherwise...)
}

// Tap passes the builder to a callback for inspection without modifying it.
func (b *Builder) Tap(callback func(*Builder)) *Builder {
	callback(b)
	return b
}

// invalidOperator reports whether an operator is invalid.
func invalidOperator(operator string) bool {
	return !validOperators[strings.ToLower(operator)]
}

// stringsToAny converts a string slice to an any slice.
func stringsToAny(s []string) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

// Reorder clears all order by clauses.
func (b *Builder) Reorder(columns ...string) *Builder {
	b.orders = nil
	b.bindings[BindingOrder] = nil
	if len(columns) > 0 {
		for _, col := range columns {
			b.OrderBy(col)
		}
	}
	return b
}

// getDB returns the underlying *sql.DB for raw operations.
func (b *Builder) getDB() *sql.DB {
	type dbGetter interface{ DB() *sql.DB }
	if dg, ok := b.connection.(dbGetter); ok {
		return dg.DB()
	}
	return nil
}
