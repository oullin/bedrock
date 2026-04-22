package search

import (
	"context"

	contract "github.com/bedrock/packages/contracts/search"
)

// Builder provides a fluent API for constructing search queries.
// It mirrors Upstream's Search\Builder class.
type Builder struct {
	model       contract.Searchable
	query       string
	wheres      map[string]any
	whereIns    map[string][]any
	whereNotIns map[string][]any
	orders      []contract.Order
	limit       int
	index       string
	options     map[string]any
	callback    func(contract.Engine, string, map[string]any) any
	engine      contract.Engine

	// queryCallback is applied by the Search method on the model before
	// the builder is passed to the engine.
	queryCallback func(*Builder)
}

// NewBuilder creates a new search builder for the given model and query.

// Where adds a constraint to the search query.

// WhereIn adds a "where in" constraint to the search query.

// WhereNotIn adds a "where not in" constraint to the search query.

// OrderBy adds an order-by directive.

// Latest orders by the given column descending (defaults to "created_at").

// Oldest orders by the given column ascending (defaults to "created_at").

// Take sets the maximum number of results.

// Within sets a custom index name for this search.

// WithOptions sets engine-specific options.

// Query sets a callback for engine-specific query modification.

// WithQueryCallback sets a callback that modifies the builder before
// it is passed to the engine.

// SetEngine sets the engine instance on the builder.

// GetEngine returns the engine instance set on the builder.

// --- SearchBuilder interface implementation ---

// GetQuery returns the search query string.

// GetModel returns the model being searched.

// GetIndex returns the custom index name, or the model's default.

// GetLimit returns the result limit, or 0 for no limit.

// GetWheres returns the where constraints.

// GetWhereIns returns the whereIn constraints.

// GetWhereNotIns returns the whereNotIn constraints.

// GetOrders returns the order-by directives.

// GetOptions returns engine-specific options.

// GetCallback returns the engine-specific query callback.

// HasCallback reports whether a callback was set.

// --- Terminal methods (execute search) ---

// Get executes the search and returns the matching models.

// Raw executes the search and returns the raw engine results.

// Keys executes the search and returns only the primary keys.

// Cursor executes the search and returns a lazy iterator.

// Paginate executes the search with length-aware pagination.

// SimplePaginate executes the search with simple (next/prev) pagination.

// PaginateRaw executes the search with pagination and returns raw results.

// SimplePaginateRaw executes the search with simple pagination and raw results.

// PaginatedResult holds models and pagination metadata.
type PaginatedResult struct {
	Models      []contract.Searchable
	Total       int64
	PerPage     int
	CurrentPage int
	RawResults  any
}

// HasMorePages reports whether there are more pages of results.

// LastPage returns the last page number.

// RawPaginatedResult holds raw engine results and pagination metadata.
type RawPaginatedResult struct {
	Results     any
	Total       int64
	PerPage     int
	CurrentPage int
}

func NewBuilder(model contract.Searchable, query string, callback ...func(contract.Engine, string, map[string]any) any) *Builder {
	b := &Builder{
		model:       model,
		query:       query,
		wheres:      make(map[string]any),
		whereIns:    make(map[string][]any),
		whereNotIns: make(map[string][]any),
		options:     make(map[string]any),
	}

	if len(callback) > 0 {
		b.callback = callback[0]
	}

	return b
}

func (b *Builder) Where(key string, values ...any) *Builder {
	if len(values) == 0 {
		b.wheres[key] = nil

		return b
	}

	if len(values) == 1 {
		b.wheres[key] = values[0]

		return b
	}

	b.wheres[key] = map[string]any{
		"__operator": values[0],
		"__value":    values[1],
	}

	return b
}

func (b *Builder) WhereIn(key string, values []any) *Builder {
	b.whereIns[key] = values

	return b
}

func (b *Builder) WhereNotIn(key string, values []any) *Builder {
	b.whereNotIns[key] = values

	return b
}

func (b *Builder) OrderBy(column, direction string) *Builder {
	b.orders = append(b.orders, contract.Order{
		Column:    column,
		Direction: direction,
	})

	return b
}

func (b *Builder) Latest(column ...string) *Builder {
	col := "created_at"

	if len(column) > 0 {
		col = column[0]
	}

	return b.OrderBy(col, "desc")
}

func (b *Builder) Oldest(column ...string) *Builder {
	col := "created_at"

	if len(column) > 0 {
		col = column[0]
	}

	return b.OrderBy(col, "asc")
}

func (b *Builder) Take(limit int) *Builder {
	b.limit = limit

	return b
}

func (b *Builder) Within(index string) *Builder {
	b.index = index

	return b
}

func (b *Builder) WithOptions(options map[string]any) *Builder {
	for k, v := range options {
		b.options[k] = v
	}

	return b
}

func (b *Builder) Query(callback func(contract.Engine, string, map[string]any) any) *Builder {
	b.callback = callback

	return b
}

func (b *Builder) WithQueryCallback(fn func(*Builder)) *Builder {
	b.queryCallback = fn

	return b
}

func (b *Builder) SetEngine(engine contract.Engine) *Builder {
	b.engine = engine

	return b
}

func (b *Builder) GetEngine() contract.Engine {
	return b.engine
}

func (b *Builder) GetQuery() string { return b.query }

func (b *Builder) GetModel() contract.Searchable { return b.model }

func (b *Builder) GetIndex() string {
	if b.index != "" {
		return b.index
	}

	return b.model.SearchableAs()
}

func (b *Builder) GetLimit() int { return b.limit }

func (b *Builder) GetWheres() map[string]any { return b.wheres }

func (b *Builder) GetWhereIns() map[string][]any { return b.whereIns }

func (b *Builder) GetWhereNotIns() map[string][]any { return b.whereNotIns }

func (b *Builder) GetOrders() []contract.Order { return b.orders }

func (b *Builder) GetOptions() map[string]any { return b.options }

func (b *Builder) GetCallback() func(contract.Engine, string, map[string]any) any {
	return b.callback
}

func (b *Builder) HasCallback() bool { return b.callback != nil }

func (b *Builder) Get(ctx context.Context) ([]contract.Searchable, error) {
	b.applyQueryCallback()

	engine := b.engine

	if engine == nil {
		return nil, ErrEngineNotConfigured
	}

	results, err := engine.Search(ctx, b)

	if err != nil {
		return nil, err
	}

	return engine.Map(ctx, results, b.model)
}

func (b *Builder) Raw(ctx context.Context) (any, error) {
	b.applyQueryCallback()

	engine := b.engine

	if engine == nil {
		return nil, ErrEngineNotConfigured
	}

	return engine.Search(ctx, b)
}

func (b *Builder) Keys(ctx context.Context) ([]any, error) {
	b.applyQueryCallback()

	engine := b.engine

	if engine == nil {
		return nil, ErrEngineNotConfigured
	}

	results, err := engine.Search(ctx, b)

	if err != nil {
		return nil, err
	}

	return engine.MapIds(results), nil
}

func (b *Builder) Cursor(ctx context.Context) (func(yield func(contract.Searchable) bool), error) {
	b.applyQueryCallback()

	engine := b.engine

	if engine == nil {
		return nil, ErrEngineNotConfigured
	}

	results, err := engine.Search(ctx, b)

	if err != nil {
		return nil, err
	}

	return engine.LazyMap(ctx, results, b.model), nil
}

func (b *Builder) Paginate(ctx context.Context, perPage, page int) (*PaginatedResult, error) {
	b.applyQueryCallback()

	engine := b.engine

	if engine == nil {
		return nil, ErrEngineNotConfigured
	}

	results, err := engine.Paginate(ctx, b, perPage, page)

	if err != nil {
		return nil, err
	}

	models, err := engine.Map(ctx, results, b.model)

	if err != nil {
		return nil, err
	}

	total := engine.GetTotalCount(results)

	return &PaginatedResult{
		Models:      models,
		Total:       total,
		PerPage:     perPage,
		CurrentPage: page,
		RawResults:  results,
	}, nil
}

func (b *Builder) SimplePaginate(ctx context.Context, perPage, page int) (*PaginatedResult, error) {
	return b.Paginate(ctx, perPage, page)
}

func (b *Builder) PaginateRaw(ctx context.Context, perPage, page int) (*RawPaginatedResult, error) {
	b.applyQueryCallback()

	engine := b.engine

	if engine == nil {
		return nil, ErrEngineNotConfigured
	}

	results, err := engine.Paginate(ctx, b, perPage, page)

	if err != nil {
		return nil, err
	}

	total := engine.GetTotalCount(results)

	return &RawPaginatedResult{
		Results:     results,
		Total:       total,
		PerPage:     perPage,
		CurrentPage: page,
	}, nil
}

func (b *Builder) SimplePaginateRaw(ctx context.Context, perPage, page int) (*RawPaginatedResult, error) {
	return b.PaginateRaw(ctx, perPage, page)
}

func (p *PaginatedResult) HasMorePages() bool {
	return int64(p.CurrentPage*p.PerPage) < p.Total
}

func (p *PaginatedResult) LastPage() int {
	if p.Total <= 0 || p.PerPage <= 0 {
		return 1
	}

	last := int(p.Total) / p.PerPage

	if int(p.Total)%p.PerPage != 0 {
		last++
	}

	return last
}

// HasMorePages reports whether there are more pages of results.
func (p *RawPaginatedResult) HasMorePages() bool {
	return int64(p.CurrentPage*p.PerPage) < p.Total
}

func (b *Builder) applyQueryCallback() {
	if b.queryCallback != nil {
		b.queryCallback(b)
	}
}
