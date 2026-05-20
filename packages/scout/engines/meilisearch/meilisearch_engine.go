package meilisearch

import (
	"context"
	"fmt"
	"strings"

	contract "github.com/bedrock/packages/contracts/scout"
	ms "github.com/meilisearch/meilisearch-go"
)

// Engine is a Scout search engine backed by Meilisearch.
// It mirrors Scout MeilisearchEngine.
type Engine struct {
	client     ms.ServiceManager
	softDelete bool
}

// Compile-time interface checks.
var _ contract.Engine = (*Engine)(nil)
var _ contract.UpdatesIndexSettings = (*Engine)(nil)

// New creates a new Meilisearch engine.
func New(client ms.ServiceManager, softDelete ...bool) *Engine {
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
	documents := make([]map[string]any, len(models))

	for i, model := range models {
		doc := model.ToSearchableArray()
		doc[model.GetScoutKeyName()] = model.GetScoutKey()

		// Add scout metadata (e.g. soft delete flag).
		for k, v := range model.GetScoutMetadata() {
			doc[k] = v
		}

		if e.softDelete && model.UsesSoftDelete() {
			if _, ok := doc["__soft_deleted"]; !ok {
				doc["__soft_deleted"] = 0
			}
		}

		documents[i] = doc
	}

	idx := e.client.Index(index)
	_, err := idx.AddDocuments(documents, models[0].GetScoutKeyName())

	if err != nil {
		return fmt.Errorf("scout: meilisearch update failed: %w", err)
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

	idx := e.client.Index(index)
	_, err := idx.DeleteDocuments(ids)

	if err != nil {
		return fmt.Errorf("scout: meilisearch delete failed: %w", err)
	}

	return nil
}

func (e *Engine) Search(ctx context.Context, builder contract.SearchBuilder) (any, error) {
	return e.performSearch(builder, nil)
}

func (e *Engine) Paginate(ctx context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	return e.performSearch(builder, &ms.SearchRequest{
		Limit:  int64(perPage),
		Offset: int64((page - 1) * perPage),
	})
}

func (e *Engine) MapIds(results any) []any {
	sr, ok := results.(*ms.SearchResponse)

	if !ok {
		return []any{}
	}

	ids := make([]any, len(sr.Hits))

	for i, hit := range sr.Hits {
		if m, ok := hit.(map[string]any); ok {
			ids[i] = m["id"]
		}
	}

	return ids
}

func (e *Engine) Map(ctx context.Context, results any, model contract.Searchable) ([]contract.Searchable, error) {
	// Meilisearch returns raw hits. Model hydration is handled by the caller.
	return []contract.Searchable{}, nil
}

func (e *Engine) LazyMap(ctx context.Context, results any, model contract.Searchable) func(yield func(contract.Searchable) bool) {
	return func(yield func(contract.Searchable) bool) {}
}

func (e *Engine) GetTotalCount(results any) int64 {
	sr, ok := results.(*ms.SearchResponse)

	if !ok {
		return 0
	}

	return sr.EstimatedTotalHits
}

func (e *Engine) Flush(ctx context.Context, model contract.Searchable) error {
	idx := e.client.Index(model.SearchableAs())
	_, err := idx.DeleteAllDocuments()

	if err != nil {
		return fmt.Errorf("scout: meilisearch flush failed: %w", err)
	}

	return nil
}

func (e *Engine) CreateIndex(ctx context.Context, name string, options map[string]any) error {
	config := &ms.IndexConfig{Uid: name}

	if pk, ok := options["primaryKey"].(string); ok {
		config.PrimaryKey = pk
	}

	_, err := e.client.CreateIndex(config)

	if err != nil {
		return fmt.Errorf("scout: meilisearch create index failed: %w", err)
	}

	return nil
}

func (e *Engine) DeleteIndex(ctx context.Context, name string) error {
	_, err := e.client.DeleteIndex(name)

	if err != nil {
		return fmt.Errorf("scout: meilisearch delete index failed: %w", err)
	}

	return nil
}

// DeleteAllIndexes removes all indexes from Meilisearch.
func (e *Engine) DeleteAllIndexes(ctx context.Context) error {
	indexes, err := e.client.ListIndexes(nil)

	if err != nil {
		return fmt.Errorf("scout: meilisearch get indexes failed: %w", err)
	}

	for _, idx := range indexes.Results {
		if _, err := e.client.DeleteIndex(idx.UID); err != nil {
			return fmt.Errorf("scout: meilisearch delete index %q failed: %w", idx.UID, err)
		}
	}

	return nil
}

// UpdateIndexSettings updates the settings for the given index.
func (e *Engine) UpdateIndexSettings(ctx context.Context, name string, settings map[string]any) error {
	idx := e.client.Index(name)
	msSettings := &ms.Settings{}

	if filterableAttrs, ok := settings["filterableAttributes"].([]string); ok {
		msSettings.FilterableAttributes = filterableAttrs
	}

	if sortableAttrs, ok := settings["sortableAttributes"].([]string); ok {
		msSettings.SortableAttributes = sortableAttrs
	}

	if searchableAttrs, ok := settings["searchableAttributes"].([]string); ok {
		msSettings.SearchableAttributes = searchableAttrs
	}

	if displayedAttrs, ok := settings["displayedAttributes"].([]string); ok {
		msSettings.DisplayedAttributes = displayedAttrs
	}

	_, err := idx.UpdateSettings(msSettings)

	if err != nil {
		return fmt.Errorf("scout: meilisearch update settings failed: %w", err)
	}

	return nil
}

// GetIndexSettings returns the settings for the given index.
func (e *Engine) GetIndexSettings(ctx context.Context, name string) (map[string]any, error) {
	idx := e.client.Index(name)
	settings, err := idx.GetSettings()

	if err != nil {
		return nil, fmt.Errorf("scout: meilisearch get settings failed: %w", err)
	}

	return map[string]any{
		"filterableAttributes": settings.FilterableAttributes,
		"sortableAttributes":   settings.SortableAttributes,
		"searchableAttributes": settings.SearchableAttributes,
		"displayedAttributes":  settings.DisplayedAttributes,
	}, nil
}

// GetClient returns the underlying Meilisearch client for advanced operations.
func (e *Engine) GetClient() ms.ServiceManager {
	return e.client
}

func (e *Engine) performSearch(builder contract.SearchBuilder, searchReq *ms.SearchRequest) (*ms.SearchResponse, error) {
	index := builder.GetIndex()
	query := builder.GetQuery()

	if searchReq == nil {
		searchReq = &ms.SearchRequest{}
	}

	// Apply limit from builder.
	if builder.GetLimit() > 0 && searchReq.Limit == 0 {
		searchReq.Limit = int64(builder.GetLimit())
	}

	// Build filters.
	filters := e.buildFilters(builder)

	if len(filters) > 0 {
		searchReq.Filter = filters
	}

	// Build sort.
	sorts := e.buildSort(builder)

	if len(sorts) > 0 {
		searchReq.Sort = sorts
	}

	idx := e.client.Index(index)
	result, err := idx.Search(query, searchReq)

	if err != nil {
		return nil, fmt.Errorf("scout: meilisearch search failed: %w", err)
	}

	return result, nil
}

func (e *Engine) buildFilters(builder contract.SearchBuilder) []string {
	var filters []string

	for key, value := range builder.GetWheres() {
		filters = append(filters, fmt.Sprintf("%s = %q", key, fmt.Sprintf("%v", value)))
	}

	for key, values := range builder.GetWhereIns() {
		parts := make([]string, len(values))

		for i, v := range values {
			parts[i] = fmt.Sprintf("%s = %q", key, fmt.Sprintf("%v", v))
		}

		filters = append(filters, "("+strings.Join(parts, " OR ")+")")
	}

	for key, values := range builder.GetWhereNotIns() {
		for _, v := range values {
			filters = append(filters, fmt.Sprintf("%s != %q", key, fmt.Sprintf("%v", v)))
		}
	}

	if e.softDelete {
		filters = append(filters, "__soft_deleted = 0")
	}

	return filters
}

func (e *Engine) buildSort(builder contract.SearchBuilder) []string {
	orders := builder.GetOrders()

	if len(orders) == 0 {
		return nil
	}

	sorts := make([]string, len(orders))

	for i, o := range orders {
		sorts[i] = o.Column + ":" + o.Direction
	}

	return sorts
}
