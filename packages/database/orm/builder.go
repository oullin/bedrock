package orm

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/database/query"
	"github.com/bedrock/packages/pagination"
)

// Builder wraps a query.Builder with Orm model awareness. It provides
// model hydration, eager loading, scope application, and type-safe results.
type Builder struct {
	query      *query.Builder
	model      *Model
	eagerLoad  map[string]func(*Builder)
	scopes     map[string]Scope
	removedScopes []string
}

// NewBuilder creates a new Orm Builder.
func NewBuilder(q *query.Builder, model *Model) *Builder {
	return &Builder{
		query:     q,
		model:     model,
		eagerLoad: make(map[string]func(*Builder)),
		scopes:    make(map[string]Scope),
	}
}

// GetQuery returns the underlying query builder.
func (b *Builder) GetQuery() *query.Builder { return b.query }

// SetQuery sets the underlying query builder.
func (b *Builder) SetQuery(q *query.Builder) { b.query = q }

// GetModel returns the model.
func (b *Builder) GetModel() *Model { return b.model }

// SetModel sets the model and configures the builder.
func (b *Builder) SetModel(model *Model) *Builder {
	b.model = model
	b.query.From(model.GetTable())
	return b
}

// ---- Retrieval ----

// Find finds a model by its primary key.
func (b *Builder) Find(ctx context.Context, id any) (*Model, error) {
	return b.Where(b.model.GetQualifiedKeyName(), id).First(ctx)
}

// FindOrFail finds a model by primary key or returns an error.
func (b *Builder) FindOrFail(ctx context.Context, id any) (*Model, error) {
	result, err := b.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrModelNotFound
	}
	return result, nil
}

// FindMany finds multiple models by their primary keys.
func (b *Builder) FindMany(ctx context.Context, ids []any) ([]*Model, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return b.WhereIn(b.model.GetQualifiedKeyName(), ids).GetModels(ctx)
}

// First returns the first result of the query.
func (b *Builder) First(ctx context.Context, columns ...string) (*Model, error) {
	row, err := b.query.First(ctx, columns...)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	return b.hydrateModel(row), nil
}

// FirstOrFail returns the first result or an error.
func (b *Builder) FirstOrFail(ctx context.Context, columns ...string) (*Model, error) {
	result, err := b.First(ctx, columns...)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrModelNotFound
	}
	return result, nil
}

// Sole returns the only matching record, or errors if 0 or > 1.
func (b *Builder) Sole(ctx context.Context, columns ...string) (*Model, error) {
	results, err := b.Take(2).GetModels(ctx, columns...)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrModelNotFound
	}
	if len(results) > 1 {
		return nil, ErrMultipleRecords
	}
	return results[0], nil
}

// Get executes the query and returns all matching models.
func (b *Builder) Get(ctx context.Context, columns ...string) ([]*Model, error) {
	return b.GetModels(ctx, columns...)
}

// GetModels executes the query and hydrates models from results.
func (b *Builder) GetModels(ctx context.Context, columns ...string) ([]*Model, error) {
	rows, err := b.query.Get(ctx, columns...)
	if err != nil {
		return nil, err
	}
	models := make([]*Model, 0, len(rows))
	for _, row := range rows {
		models = append(models, b.hydrateModel(row))
	}
	return models, nil
}

// Value returns a single column value from the first row.
func (b *Builder) Value(ctx context.Context, column string) (any, error) {
	return b.query.Value(ctx, column)
}

// Pluck returns a slice of values for a single column.
func (b *Builder) Pluck(ctx context.Context, column string) ([]any, error) {
	return b.query.Pluck(ctx, column)
}

// FindOr finds by primary key or calls the fallback.
func (b *Builder) FindOr(ctx context.Context, id any, fallback func() (*Model, error)) (*Model, error) {
	result, err := b.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	return fallback()
}

// FindOrNew finds by primary key or returns a new unsaved instance.
func (b *Builder) FindOrNew(ctx context.Context, id any) (*Model, error) {
	result, err := b.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	return b.newModelInstance(), nil
}

// FirstOr returns the first result or calls the fallback.
func (b *Builder) FirstOr(ctx context.Context, fallback func() (*Model, error), columns ...string) (*Model, error) {
	result, err := b.First(ctx, columns...)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	return fallback()
}

// FirstWhere finds the first record matching a single condition.
func (b *Builder) FirstWhere(ctx context.Context, args ...any) (*Model, error) {
	return b.Where(args...).First(ctx)
}

// SoleValue returns a single column value, ensuring exactly one row matches.
func (b *Builder) SoleValue(ctx context.Context, column string) (any, error) {
	model, err := b.Sole(ctx, column)
	if err != nil {
		return nil, err
	}
	return model.GetAttribute(column), nil
}

// ValueOrFail returns a single column value or an error if not found.
func (b *Builder) ValueOrFail(ctx context.Context, column string) (any, error) {
	val, err := b.Value(ctx, column)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, ErrModelNotFound
	}
	return val, nil
}

// CreateQuietly creates a model without firing events.
func (b *Builder) CreateQuietly(ctx context.Context, attributes map[string]any) (*Model, error) {
	return b.Create(ctx, attributes)
}

// ForceCreateQuietly creates without mass assignment protection or events.
func (b *Builder) ForceCreateQuietly(ctx context.Context, attributes map[string]any) (*Model, error) {
	return b.ForceCreate(ctx, attributes)
}

// CreateOrFirst creates a record or returns the first matching one.
func (b *Builder) CreateOrFirst(ctx context.Context, attributes map[string]any, values ...map[string]any) (*Model, error) {
	return b.FirstOrCreate(ctx, attributes, values...)
}

// IncrementOrCreate increments a column or creates the record.
func (b *Builder) IncrementOrCreate(ctx context.Context, attributes map[string]any, column string, amount ...any) (*Model, error) {
	model, err := b.FirstOrNew(ctx, attributes)
	if err != nil {
		return nil, err
	}
	if model.Exists() {
		return model, model.Increment(ctx, column, amount...)
	}
	amt := any(1)
	if len(amount) > 0 {
		amt = amount[0]
	}
	model.SetAttribute(column, amt)
	return model, model.Save(ctx)
}

// Upsert performs a bulk upsert.
func (b *Builder) Upsert(ctx context.Context, values []map[string]any, uniqueBy []string, update []string) (int64, error) {
	return b.query.Upsert(ctx, values, uniqueBy, update)
}

// Touch updates the timestamp on all matching models.
func (b *Builder) Touch(ctx context.Context, column ...string) (int64, error) {
	col := "updated_at"
	if len(column) > 0 {
		col = column[0]
	}
	return b.Update(ctx, map[string]any{
		col: b.model.FreshTimestampString(),
	})
}

// IncrementEach increments multiple columns.
func (b *Builder) IncrementEach(ctx context.Context, columns map[string]any, extra ...map[string]any) (int64, error) {
	return b.query.IncrementEach(ctx, columns, extra...)
}

// DecrementEach decrements multiple columns.
func (b *Builder) DecrementEach(ctx context.Context, columns map[string]any, extra ...map[string]any) (int64, error) {
	return b.query.DecrementEach(ctx, columns, extra...)
}

// Hydrate creates model instances from raw maps.
func (b *Builder) Hydrate(items []map[string]any) []*Model {
	models := make([]*Model, 0, len(items))
	for _, item := range items {
		models = append(models, b.hydrateModel(item))
	}
	return models
}

// FromQuery creates models from a raw SQL query.
func (b *Builder) FromQuery(ctx context.Context, sql string, bindings ...any) ([]*Model, error) {
	conn := b.query.GetConnection()
	if conn == nil {
		return nil, ErrModelNotFound
	}
	rows, err := conn.Select(ctx, sql, bindings...)
	if err != nil {
		return nil, err
	}
	return b.Hydrate(rows), nil
}

// WithCasts applies casts to the query builder for this query.
func (b *Builder) WithCasts(casts map[string]string) *Builder {
	for k, v := range casts {
		b.model.casts[k] = v
	}
	return b
}

// Cursor iterates results one at a time for memory efficiency.
func (b *Builder) Cursor(ctx context.Context, fn func(*Model) bool) error {
	return b.query.Cursor(ctx, func(row map[string]any) bool {
		model := b.hydrateModel(row)
		return fn(model)
	})
}

// Chunk processes models in chunks.
func (b *Builder) Chunk(ctx context.Context, count int, fn func([]*Model, int) bool) error {
	return b.query.Chunk(ctx, count, func(rows []map[string]any, page int) bool {
		models := make([]*Model, 0, len(rows))
		for _, row := range rows {
			models = append(models, b.hydrateModel(row))
		}
		return fn(models, page)
	})
}

// ChunkByID processes models in ID-based chunks.
func (b *Builder) ChunkByID(ctx context.Context, count int, column string, fn func([]*Model) bool) error {
	return b.query.ChunkByID(ctx, count, column, func(rows []map[string]any) bool {
		models := make([]*Model, 0, len(rows))
		for _, row := range rows {
			models = append(models, b.hydrateModel(row))
		}
		return fn(models)
	})
}

// ---- Creation ----

// Create creates a new model and persists it.
func (b *Builder) Create(ctx context.Context, attributes map[string]any) (*Model, error) {
	model := b.newModelInstance()
	if err := model.Fill(attributes); err != nil {
		return nil, err
	}
	if err := model.Save(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

// ForceCreate creates a model without mass assignment protection.
func (b *Builder) ForceCreate(ctx context.Context, attributes map[string]any) (*Model, error) {
	model := b.newModelInstance()
	model.ForceFill(attributes)
	if err := model.Save(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

// FirstOrNew finds the first record matching attributes, or creates a new instance.
func (b *Builder) FirstOrNew(ctx context.Context, attributes map[string]any, values ...map[string]any) (*Model, error) {
	result, err := b.Where(attributes).First(ctx)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	model := b.newModelInstance()
	model.ForceFill(attributes)
	if len(values) > 0 {
		model.ForceFill(values[0])
	}
	return model, nil
}

// FirstOrCreate finds the first record matching attributes, or creates it.
func (b *Builder) FirstOrCreate(ctx context.Context, attributes map[string]any, values ...map[string]any) (*Model, error) {
	result, err := b.Where(attributes).First(ctx)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	merged := make(map[string]any, len(attributes))
	for k, v := range attributes {
		merged[k] = v
	}
	if len(values) > 0 {
		for k, v := range values[0] {
			merged[k] = v
		}
	}

	return b.Create(ctx, merged)
}

// UpdateOrCreate finds the first record matching attributes or creates it, then updates.
func (b *Builder) UpdateOrCreate(ctx context.Context, attributes map[string]any, values map[string]any) (*Model, error) {
	result, err := b.FirstOrNew(ctx, attributes)
	if err != nil {
		return nil, err
	}
	result.ForceFill(values)
	if err := result.Save(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

// ---- Modification ----

// Update updates all rows matching the query.
func (b *Builder) Update(ctx context.Context, values map[string]any) (int64, error) {
	return b.query.Update(ctx, values)
}

// Delete deletes all rows matching the query.
func (b *Builder) Delete(ctx context.Context) (int64, error) {
	return b.query.Delete(ctx)
}

// ForceDelete force-deletes (hard delete) all matching rows.
func (b *Builder) ForceDelete(ctx context.Context) (int64, error) {
	return b.query.Delete(ctx)
}

// Increment increments a column.
func (b *Builder) Increment(ctx context.Context, column string, amount ...any) (int64, error) {
	return b.query.Increment(ctx, column, amount...)
}

// Decrement decrements a column.
func (b *Builder) Decrement(ctx context.Context, column string, amount ...any) (int64, error) {
	return b.query.Decrement(ctx, column, amount...)
}

// ---- Aggregates ----

// Count returns the count of matching rows.
func (b *Builder) Count(ctx context.Context, columns ...string) (int64, error) {
	return b.query.Count(ctx, columns...)
}

// Max returns the max value of a column.
func (b *Builder) Max(ctx context.Context, column string) (any, error) {
	return b.query.Max(ctx, column)
}

// Min returns the min value of a column.
func (b *Builder) Min(ctx context.Context, column string) (any, error) {
	return b.query.Min(ctx, column)
}

// Sum returns the sum of a column.
func (b *Builder) Sum(ctx context.Context, column string) (float64, error) {
	return b.query.Sum(ctx, column)
}

// Avg returns the average of a column.
func (b *Builder) Avg(ctx context.Context, column string) (float64, error) {
	return b.query.Avg(ctx, column)
}

// Exists checks if any rows match.
func (b *Builder) Exists(ctx context.Context) (bool, error) {
	return b.query.Exists(ctx)
}

// DoesntExist checks if no rows match.
func (b *Builder) DoesntExist(ctx context.Context) (bool, error) {
	return b.query.DoesntExist(ctx)
}

// ---- Pagination ----

// Paginate paginates the query results.
func (b *Builder) Paginate(ctx context.Context, perPage, page int, columns ...string) (*pagination.LengthAwarePaginator[map[string]any], error) {
	return b.query.Paginate(ctx, perPage, page, columns...)
}

// SimplePaginate paginates without total count.
func (b *Builder) SimplePaginate(ctx context.Context, perPage, page int, columns ...string) (*pagination.Paginator[map[string]any], error) {
	return b.query.SimplePaginate(ctx, perPage, page, columns...)
}

// ---- Query Builder Proxies ----

// Select sets the columns.
func (b *Builder) Select(columns ...any) *Builder {
	b.query.Select(columns...)
	return b
}

// Where adds a where clause.
func (b *Builder) Where(args ...any) *Builder {
	b.query.Where(args...)
	return b
}

// OrWhere adds an or where clause.
func (b *Builder) OrWhere(args ...any) *Builder {
	b.query.OrWhere(args...)
	return b
}

// WhereIn adds a where in clause.
func (b *Builder) WhereIn(column string, values []any) *Builder {
	b.query.WhereIn(column, values)
	return b
}

// WhereNotIn adds a where not in clause.
func (b *Builder) WhereNotIn(column string, values []any) *Builder {
	b.query.WhereNotIn(column, values)
	return b
}

// WhereNull adds a where null clause.
func (b *Builder) WhereNull(columns ...string) *Builder {
	b.query.WhereNull(columns...)
	return b
}

// WhereNotNull adds a where not null clause.
func (b *Builder) WhereNotNull(columns ...string) *Builder {
	b.query.WhereNotNull(columns...)
	return b
}

// WhereBetween adds a where between clause.
func (b *Builder) WhereBetween(column string, values [2]any) *Builder {
	b.query.WhereBetween(column, values)
	return b
}

// WhereDate adds a where date clause.
func (b *Builder) WhereDate(column, operator, value string) *Builder {
	b.query.WhereDate(column, operator, value)
	return b
}

// WhereRaw adds a raw where clause.
func (b *Builder) WhereRaw(sql string, bindings ...any) *Builder {
	b.query.WhereRaw(sql, bindings...)
	return b
}

// OrderBy adds an order by clause.
func (b *Builder) OrderBy(column string, direction ...string) *Builder {
	b.query.OrderBy(column, direction...)
	return b
}

// OrderByDesc adds a descending order by clause.
func (b *Builder) OrderByDesc(column string) *Builder {
	b.query.OrderByDesc(column)
	return b
}

// Latest orders by created_at descending.
func (b *Builder) Latest(column ...string) *Builder {
	b.query.Latest(column...)
	return b
}

// Oldest orders by created_at ascending.
func (b *Builder) Oldest(column ...string) *Builder {
	b.query.Oldest(column...)
	return b
}

// Limit sets the limit.
func (b *Builder) Limit(n int) *Builder {
	b.query.Limit(n)
	return b
}

// Take is an alias for Limit.
func (b *Builder) Take(n int) *Builder {
	b.query.Take(n)
	return b
}

// Offset sets the offset.
func (b *Builder) Offset(n int) *Builder {
	b.query.Offset(n)
	return b
}

// Skip is an alias for Offset.
func (b *Builder) Skip(n int) *Builder {
	b.query.Skip(n)
	return b
}

// GroupBy adds a group by clause.
func (b *Builder) GroupBy(groups ...string) *Builder {
	b.query.GroupBy(groups...)
	return b
}

// Having adds a having clause.
func (b *Builder) Having(column string, args ...any) *Builder {
	b.query.Having(column, args...)
	return b
}

// Join adds a join clause.
func (b *Builder) Join(table string, args ...any) *Builder {
	b.query.Join(table, args...)
	return b
}

// LeftJoin adds a left join.
func (b *Builder) LeftJoin(table string, args ...any) *Builder {
	b.query.LeftJoin(table, args...)
	return b
}

// Distinct marks the query as distinct.
func (b *Builder) Distinct(columns ...string) *Builder {
	b.query.Distinct(columns...)
	return b
}

// ForPage sets limit/offset for pagination.
func (b *Builder) ForPage(page, perPage int) *Builder {
	b.query.ForPage(page, perPage)
	return b
}

// ---- Eager Loading ----

// With registers relationships for eager loading.
func (b *Builder) With(relations ...string) *Builder {
	for _, rel := range relations {
		b.eagerLoad[rel] = nil
	}
	return b
}

// Without removes relationships from eager loading.
func (b *Builder) Without(relations ...string) *Builder {
	for _, rel := range relations {
		delete(b.eagerLoad, rel)
	}
	return b
}

// WithOnly replaces all eager loads with just these.
func (b *Builder) WithOnly(relations ...string) *Builder {
	b.eagerLoad = make(map[string]func(*Builder))
	return b.With(relations...)
}

// GetEagerLoads returns the registered eager loads.
func (b *Builder) GetEagerLoads() map[string]func(*Builder) {
	return b.eagerLoad
}

// ---- Scopes ----

// WithGlobalScope registers a named global scope.
func (b *Builder) WithGlobalScope(name string, scope Scope) *Builder {
	b.scopes[name] = scope
	return b
}

// WithoutGlobalScope removes a named global scope.
func (b *Builder) WithoutGlobalScope(names ...string) *Builder {
	for _, name := range names {
		delete(b.scopes, name)
		b.removedScopes = append(b.removedScopes, name)
	}
	return b
}

// WithoutGlobalScopes removes all global scopes.
func (b *Builder) WithoutGlobalScopes() *Builder {
	b.scopes = make(map[string]Scope)
	return b
}

// ApplyScopes applies all registered scopes to the query.
func (b *Builder) ApplyScopes() *Builder {
	for _, scope := range b.scopes {
		scope.Apply(b.query)
	}
	return b
}

// Scopes applies local scope functions.
func (b *Builder) Scopes(scopes ...func(*Builder) *Builder) *Builder {
	for _, scope := range scopes {
		scope(b)
	}
	return b
}

// ---- Conditional ----

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

// Unless applies a callback when condition is false.
func (b *Builder) Unless(condition bool, callback func(*Builder) *Builder, otherwise ...func(*Builder) *Builder) *Builder {
	return b.When(!condition, callback, otherwise...)
}

// ---- SQL Output ----

// ToSQL returns the compiled SQL and bindings.
func (b *Builder) ToSQL() (string, []any) {
	return b.query.ToSQL()
}

// Dump returns the SQL and bindings as a string.
func (b *Builder) Dump() string {
	return b.query.Dump()
}

// ---- Hydration ----

func (b *Builder) hydrateModel(row map[string]any) *Model {
	model := b.newModelInstance()
	model.SetRawAttributes(row, true)
	model.SetExists(true)
	return model
}

func (b *Builder) newModelInstance() *Model {
	model := NewModel()
	model.SetTable(b.model.GetTable())
	model.SetPrimaryKey(b.model.GetKeyName())
	model.SetKeyType(b.model.GetKeyType())
	model.SetIncrementing(b.model.GetIncrementing())
	model.SetConnectionName(b.model.GetConnectionName())
	model.SetResolver(b.model.GetResolver())
	model.SetFillable(b.model.GetFillable())
	model.SetGuarded(b.model.GetGuarded())
	model.SetCasts(b.model.GetCasts())
	model.SetHidden(b.model.GetHidden())
	model.SetVisible(b.model.GetVisible())
	if b.model.UsesTimestamps() {
		model.InitTimestamps()
	} else {
		model.SetTimestamps(false)
	}
	return model
}

// String implements fmt.Stringer.
func (b *Builder) String() string {
	sql, bindings := b.ToSQL()
	return fmt.Sprintf("SQL: %s | Bindings: %v", sql, bindings)
}
