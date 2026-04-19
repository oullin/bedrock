package typesense

import (
	"context"
	"fmt"
	"strings"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

// Engine is a Search search engine backed by Typesense.
// It mirrors Upstream Search's TypesenseEngine.
type Engine struct {
	client     *typesense.Client
	softDelete bool
}

// Compile-time interface check.
var _ contract.Engine = (*Engine)(nil)

// New creates a new Typesense engine.
func New(client *typesense.Client, softDelete ...bool) *Engine {
	sd := false

	if len(softDelete) > 0 {
		sd = softDelete[0]
	}

	return &Engine{client: client, softDelete: sd}
}

func (e *Engine) Update(ctx context.Context, models []contract.Searchable) error {
	if len(models) == 0 {
		return nil
	}

	collection := models[0].SearchableAs()

	for _, model := range models {
		doc := model.ToSearchableArray()
		doc["id"] = fmt.Sprintf("%v", model.GetScoutKey())

		// Add search metadata.
		for k, v := range model.GetScoutMetadata() {
			doc[k] = v
		}

		if e.softDelete && model.UsesSoftDelete() {
			if _, ok := doc["__soft_deleted"]; !ok {
				doc["__soft_deleted"] = 0
			}
		}

		_, err := e.client.Collection(collection).Documents().Upsert(ctx, doc, &api.DocumentIndexParameters{})

		if err != nil {
			return fmt.Errorf("search: typesense update failed: %w", err)
		}
	}

	return nil
}

func (e *Engine) Delete(ctx context.Context, models []contract.Searchable) error {
	if len(models) == 0 {
		return nil
	}

	collection := models[0].SearchableAs()

	for _, model := range models {
		id := fmt.Sprintf("%v", model.GetScoutKey())
		_, err := e.client.Collection(collection).Document(id).Delete(ctx)

		if err != nil {
			return fmt.Errorf("search: typesense delete failed: %w", err)
		}
	}

	return nil
}

func (e *Engine) Search(ctx context.Context, builder contract.SearchBuilder) (any, error) {
	return e.performSearch(ctx, builder, 0, 0)
}

func (e *Engine) Paginate(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	return e.performSearch(ctx, builder, perPage, page)
}

func (e *Engine) MapIds(results any) []any {
	sr, ok := results.(*api.SearchResult)

	if !ok || sr.Hits == nil {
		return []any{}
	}

	ids := make([]any, 0, len(*sr.Hits))

	for _, hit := range *sr.Hits {
		if hit.Document != nil {
			doc := *hit.Document
			ids = append(ids, doc["id"])
		}
	}

	return ids
}

func (e *Engine) Map(ctx context.Context, results any, model contract.Searchable) ([]contract.Searchable, error) {
	// Typesense returns raw hits. Model hydration is handled by the caller.
	return []contract.Searchable{}, nil
}

func (e *Engine) LazyMap(ctx context.Context, results any, model contract.Searchable) func(yield func(contract.Searchable) bool) {
	return func(yield func(contract.Searchable) bool) {}
}

func (e *Engine) GetTotalCount(results any) int64 {
	sr, ok := results.(*api.SearchResult)

	if !ok || sr.Found == nil {
		return 0
	}

	return int64(*sr.Found)
}

func (e *Engine) Flush(ctx context.Context, model contract.Searchable) error {
	collection := model.SearchableAs()

	// Delete and recreate the collection.
	_, err := e.client.Collection(collection).Delete(ctx)

	if err != nil {
		return fmt.Errorf("search: typesense flush failed: %w", err)
	}

	return nil
}

func (e *Engine) CreateIndex(ctx context.Context, name string, options map[string]any) error {
	schema := &api.CollectionSchema{
		Name: name,
	}

	// Parse fields from options if provided.
	if fields, ok := options["fields"].([]api.Field); ok {
		schema.Fields = fields
	}

	_, err := e.client.Collections().Create(ctx, schema)

	if err != nil {
		return fmt.Errorf("search: typesense create index failed: %w", err)
	}

	return nil
}

func (e *Engine) DeleteIndex(ctx context.Context, name string) error {
	_, err := e.client.Collection(name).Delete(ctx)

	if err != nil {
		return fmt.Errorf("search: typesense delete index failed: %w", err)
	}

	return nil
}

// GetClient returns the underlying Typesense client.
func (e *Engine) GetClient() *typesense.Client {
	return e.client
}

func (e *Engine) performSearch(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (*api.SearchResult, error) {
	collection := builder.GetIndex()
	query := builder.GetQuery()

	if query == "" {
		query = "*"
	}

	queryBy := e.getQueryBy(builder)

	params := &api.SearchCollectionParams{
		Q:       &query,
		QueryBy: &queryBy,
	}

	// Build filter.
	filterBy := e.buildFilterBy(builder)

	if filterBy != "" {
		params.FilterBy = &filterBy
	}

	// Build sort.
	sortBy := e.buildSortBy(builder)

	if sortBy != "" {
		params.SortBy = &sortBy
	}

	// Apply pagination.
	if perPage > 0 {
		pp := perPage
		params.PerPage = &pp
		params.Page = &page
	} else if builder.GetLimit() > 0 {
		l := builder.GetLimit()
		params.PerPage = &l
	}

	result, err := e.client.Collection(collection).Documents().Search(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("search: typesense search failed: %w", err)
	}

	return result, nil
}

func (e *Engine) getQueryBy(builder contract.SearchBuilder) string {
	opts := builder.GetOptions()

	if queryBy, ok := opts["query_by"].(string); ok {
		return queryBy
	}

	// Fall back to all searchable fields from the model.
	model := builder.GetModel()

	if model != nil {
		data := model.ToSearchableArray()
		columns := make([]string, 0, len(data))

		for key := range data {
			columns = append(columns, key)
		}

		return strings.Join(columns, ",")
	}

	return "*"
}

func (e *Engine) buildFilterBy(builder contract.SearchBuilder) string {
	var parts []string

	for key, value := range builder.GetWheres() {
		parts = append(parts, fmt.Sprintf("%s:=%v", key, value))
	}

	for key, values := range builder.GetWhereIns() {
		strVals := make([]string, len(values))

		for i, v := range values {
			strVals[i] = fmt.Sprintf("%v", v)
		}

		parts = append(parts, fmt.Sprintf("%s:[%s]", key, strings.Join(strVals, ",")))
	}

	for key, values := range builder.GetWhereNotIns() {
		strVals := make([]string, len(values))

		for i, v := range values {
			strVals[i] = fmt.Sprintf("%v", v)
		}

		parts = append(parts, fmt.Sprintf("%s:!=[%s]", key, strings.Join(strVals, ",")))
	}

	if e.softDelete {
		parts = append(parts, "__soft_deleted:=0")
	}

	return strings.Join(parts, " && ")
}

func (e *Engine) buildSortBy(builder contract.SearchBuilder) string {
	orders := builder.GetOrders()

	if len(orders) == 0 {
		return ""
	}

	parts := make([]string, len(orders))

	for i, o := range orders {
		parts[i] = o.Column + ":" + o.Direction
	}

	return strings.Join(parts, ",")
}
