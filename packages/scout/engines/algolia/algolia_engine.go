package algolia

import (
	"context"
	"fmt"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	contract "github.com/bedrock/packages/contracts/scout"
)

// Engine is a Scout search engine backed by Algolia.
// It mirrors Laravel Scout's AlgoliaEngine.
type Engine struct {
	client     *search.APIClient
	softDelete bool
}

// Compile-time interface checks.
var _ contract.Engine = (*Engine)(nil)
var _ contract.UpdatesIndexSettings = (*Engine)(nil)

// New creates a new Algolia engine.
func New(client *search.APIClient, softDelete ...bool) *Engine {
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

	index := models[0].SearchableAs()
	objects := make([]map[string]any, len(models))

	for i, model := range models {
		obj := model.ToSearchableArray()
		obj["objectID"] = fmt.Sprintf("%v", model.GetScoutKey())

		// Add scout metadata.
		for k, v := range model.GetScoutMetadata() {
			obj[k] = v
		}

		if e.softDelete && model.UsesSoftDelete() {
			if _, ok := obj["__soft_deleted"]; !ok {
				obj["__soft_deleted"] = 0
			}
		}

		objects[i] = obj
	}

	_, err := e.client.SaveObjects(e.client.NewApiSaveObjectsRequest(index, objects))

	if err != nil {
		return fmt.Errorf("scout: algolia update failed: %w", err)
	}

	return nil
}

func (e *Engine) Delete(ctx context.Context, models []contract.Searchable) error {
	if len(models) == 0 {
		return nil
	}

	index := models[0].SearchableAs()
	ids := make([]string, len(models))

	for i, model := range models {
		ids[i] = fmt.Sprintf("%v", model.GetScoutKey())
	}

	_, err := e.client.DeleteObjects(e.client.NewApiDeleteObjectsRequest(index, ids))

	if err != nil {
		return fmt.Errorf("scout: algolia delete failed: %w", err)
	}

	return nil
}

func (e *Engine) Search(ctx context.Context, builder contract.SearchBuilder) (any, error) {
	return e.performSearch(builder, 0, 0)
}

func (e *Engine) Paginate(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	return e.performSearch(builder, perPage, page)
}

func (e *Engine) MapIds(results any) []any {
	sr, ok := results.(*search.SearchResponse)

	if !ok {
		return []any{}
	}

	ids := make([]any, len(sr.Hits))

	for i, hit := range sr.Hits {
		ids[i] = hit["objectID"]
	}

	return ids
}

func (e *Engine) Map(ctx context.Context, results any, model contract.Searchable) ([]contract.Searchable, error) {
	// Algolia returns raw hits. Model hydration is handled by the caller.
	return []contract.Searchable{}, nil
}

func (e *Engine) LazyMap(ctx context.Context, results any, model contract.Searchable) func(yield func(contract.Searchable) bool) {
	return func(yield func(contract.Searchable) bool) {}
}

func (e *Engine) GetTotalCount(results any) int64 {
	sr, ok := results.(*search.SearchResponse)

	if !ok {
		return 0
	}

	if sr.NbHits != nil {
		return int64(*sr.NbHits)
	}

	return 0
}

func (e *Engine) Flush(ctx context.Context, model contract.Searchable) error {
	_, err := e.client.ClearObjects(e.client.NewApiClearObjectsRequest(model.SearchableAs()))

	if err != nil {
		return fmt.Errorf("scout: algolia flush failed: %w", err)
	}

	return nil
}

func (e *Engine) CreateIndex(ctx context.Context, name string, options map[string]any) error {
	// Algolia creates indexes automatically on first use.
	// Optionally configure settings.
	if len(options) > 0 {
		return e.UpdateIndexSettings(ctx, name, options)
	}

	return nil
}

func (e *Engine) DeleteIndex(ctx context.Context, name string) error {
	_, err := e.client.DeleteIndex(e.client.NewApiDeleteIndexRequest(name))

	if err != nil {
		return fmt.Errorf("scout: algolia delete index failed: %w", err)
	}

	return nil
}

// UpdateIndexSettings updates the settings for the given Algolia index.
func (e *Engine) UpdateIndexSettings(ctx context.Context, name string, settings map[string]any) error {
	indexSettings := search.NewEmptyIndexSettings()

	if searchableAttrs, ok := settings["searchableAttributes"].([]string); ok {
		indexSettings.SearchableAttributes = searchableAttrs
	}

	if attrForFaceting, ok := settings["attributesForFaceting"].([]string); ok {
		indexSettings.AttributesForFaceting = attrForFaceting
	}

	_, err := e.client.SetSettings(e.client.NewApiSetSettingsRequest(name, indexSettings))

	if err != nil {
		return fmt.Errorf("scout: algolia update settings failed: %w", err)
	}

	return nil
}

// GetIndexSettings returns the settings for the given Algolia index.
func (e *Engine) GetIndexSettings(ctx context.Context, name string) (map[string]any, error) {
	settings, err := e.client.GetSettings(e.client.NewApiGetSettingsRequest(name))

	if err != nil {
		return nil, fmt.Errorf("scout: algolia get settings failed: %w", err)
	}

	return map[string]any{
		"searchableAttributes":  settings.SearchableAttributes,
		"attributesForFaceting": settings.AttributesForFaceting,
	}, nil
}

// GetClient returns the underlying Algolia client.
func (e *Engine) GetClient() *search.APIClient {
	return e.client
}

func (e *Engine) performSearch(builder contract.SearchBuilder, perPage, page int) (*search.SearchResponse, error) {
	index := builder.GetIndex()
	query := builder.GetQuery()

	params := &search.SearchParamsObject{
		Query: &query,
	}

	// Apply filters.
	filters := e.buildFilters(builder)

	if filters != "" {
		params.Filters = &filters
	}

	// Apply pagination.
	if perPage > 0 {
		p := int32(perPage)
		params.HitsPerPage = &p
		pg := int32(page - 1) // Algolia pages are 0-indexed
		params.Page = &pg
	} else if builder.GetLimit() > 0 {
		l := int32(builder.GetLimit())
		params.HitsPerPage = &l
	}

	result, err := e.client.SearchSingleIndex(e.client.NewApiSearchSingleIndexRequest(index).WithSearchParams(
		search.SearchParamsObjectAsSearchParams(params),
	))

	if err != nil {
		return nil, fmt.Errorf("scout: algolia search failed: %w", err)
	}

	return result, nil
}

func (e *Engine) buildFilters(builder contract.SearchBuilder) string {
	var parts []string

	for key, value := range builder.GetWheres() {
		parts = append(parts, fmt.Sprintf("%s:%v", key, value))
	}

	for key, values := range builder.GetWhereIns() {
		inParts := make([]string, len(values))

		for i, v := range values {
			inParts[i] = fmt.Sprintf("%s:%v", key, v)
		}

		parts = append(parts, "("+strings.Join(inParts, " OR ")+")")
	}

	for key, values := range builder.GetWhereNotIns() {
		for _, v := range values {
			parts = append(parts, fmt.Sprintf("NOT %s:%v", key, v))
		}
	}

	if e.softDelete {
		parts = append(parts, "__soft_deleted:0")
	}

	return strings.Join(parts, " AND ")
}
