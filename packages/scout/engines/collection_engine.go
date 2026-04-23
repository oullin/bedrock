package engines

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	contract "github.com/bedrock/packages/contracts/scout"
)

// CollectionEngine is an in-memory search engine that filters models
// without requiring any external search service. It mirrors Laravel
// Scout's CollectionEngine.
type CollectionEngine struct {
	softDelete bool
}

type searchableCollectionTransformer interface {
	MakeSearchableUsing([]contract.Searchable) []contract.Searchable
}

// Compile-time interface check.

// NewCollectionEngine creates a new CollectionEngine.

// Collection engine operates in-memory; no external index to update.

// Collection engine operates in-memory; no external index to remove from.

// Apply pagination offset.

// No external index to flush.

// No index management needed for in-memory engine.

// No index management needed for in-memory engine.

// searchModels filters the builder's models based on the query and constraints.

// Apply text search filter.

// Apply where constraints.

// Apply whereIn constraints.

// Apply whereNotIn constraints.

// Apply ordering.

// Apply limit.

// getModels returns the base model slice from the builder. For the collection
// engine, models are expected to be passed via builder options.

// matchesQuery checks if any field in the searchable data contains the query
// string (case-insensitive), matching Laravel's collection engine behavior.

// matchesWheres checks that all where constraints are satisfied.

// matchesWhereIns checks that values are within the allowed sets.

// matchesWhereNotIns checks that values are NOT in the excluded sets.

// sortModels sorts models by the given order directives.

// CollectionResult holds the results of a collection engine search.
type CollectionResult struct {
	Models     []contract.Searchable
	TotalCount int64
}

var _ contract.Engine = (*CollectionEngine)(nil)

func NewCollectionEngine(softDelete ...bool) *CollectionEngine {
	sd := false

	if len(softDelete) > 0 {
		sd = softDelete[0]
	}

	return &CollectionEngine{softDelete: sd}
}

func (e *CollectionEngine) Update(_ context.Context, _ []contract.Searchable) error {

	return nil
}

func (e *CollectionEngine) Delete(_ context.Context, _ []contract.Searchable) error {

	return nil
}

func (e *CollectionEngine) Search(_ context.Context, builder contract.SearchBuilder) (any, error) {
	models := e.searchModels(builder)

	return &CollectionResult{Models: models, TotalCount: int64(len(models))}, nil
}

func (e *CollectionEngine) Paginate(_ context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	models := e.searchModels(builder)
	total := int64(len(models))

	offset := (page - 1) * perPage

	if offset > len(models) {
		offset = len(models)
	}

	end := offset + perPage

	if end > len(models) {
		end = len(models)
	}

	return &CollectionResult{
		Models:     models[offset:end],
		TotalCount: total,
	}, nil
}

func (e *CollectionEngine) MapIds(results any) []any {
	cr, ok := results.(*CollectionResult)

	if !ok {
		return []any{}
	}

	ids := make([]any, len(cr.Models))

	for i, m := range cr.Models {
		ids[i] = m.GetScoutKey()
	}

	return ids
}

func (e *CollectionEngine) Map(_ context.Context, results any, _ contract.Searchable) ([]contract.Searchable, error) {
	cr, ok := results.(*CollectionResult)

	if !ok {
		return []contract.Searchable{}, nil
	}

	return cr.Models, nil
}

func (e *CollectionEngine) LazyMap(_ context.Context, results any, _ contract.Searchable) func(yield func(contract.Searchable) bool) {
	cr, ok := results.(*CollectionResult)

	if !ok {
		return func(_ func(contract.Searchable) bool) {}
	}

	return func(yield func(contract.Searchable) bool) {
		for _, m := range cr.Models {
			if !yield(m) {
				return
			}
		}
	}
}

func (e *CollectionEngine) GetTotalCount(results any) int64 {
	cr, ok := results.(*CollectionResult)

	if !ok {
		return 0
	}

	return cr.TotalCount
}

func (e *CollectionEngine) Flush(_ context.Context, _ contract.Searchable) error {

	return nil
}

func (e *CollectionEngine) CreateIndex(_ context.Context, _ string, _ map[string]any) error {

	return nil
}

func (e *CollectionEngine) DeleteIndex(_ context.Context, _ string) error {

	return nil
}

func (e *CollectionEngine) searchModels(builder contract.SearchBuilder) []contract.Searchable {
	models := e.applyMakeSearchableUsing(e.getModels(builder))
	query := strings.TrimSpace(builder.GetQuery())
	wheres := builder.GetWheres()
	whereIns := builder.GetWhereIns()
	whereNotIns := builder.GetWhereNotIns()
	orders := builder.GetOrders()
	limit := builder.GetLimit()

	var result []contract.Searchable

	for _, model := range models {
		searchableData := model.ToSearchableArray()

		if query != "" && !e.matchesQuery(searchableData, query) {
			continue
		}

		if !e.matchesWheres(searchableData, wheres) {
			continue
		}

		if !e.matchesWhereIns(searchableData, whereIns) {
			continue
		}

		if !e.matchesWhereNotIns(searchableData, whereNotIns) {
			continue
		}

		result = append(result, model)
	}

	if len(orders) > 0 {
		e.sortModels(result, orders)
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result
}

func (e *CollectionEngine) getModels(builder contract.SearchBuilder) []contract.Searchable {
	opts := builder.GetOptions()

	if models, ok := opts["__models"].([]contract.Searchable); ok {
		return models
	}

	return nil
}

func (e *CollectionEngine) applyMakeSearchableUsing(models []contract.Searchable) []contract.Searchable {
	if len(models) == 0 {
		return models
	}

	transformer, ok := models[0].(searchableCollectionTransformer)

	if !ok {
		return models
	}

	return transformer.MakeSearchableUsing(models)
}

func (e *CollectionEngine) matchesQuery(data map[string]any, query string) bool {
	lowerQuery := strings.ToLower(query)

	for _, v := range data {
		str := strings.ToLower(fmt.Sprintf("%v", v))

		if strings.Contains(str, lowerQuery) {
			return true
		}
	}

	return false
}

func (e *CollectionEngine) matchesWheres(data map[string]any, wheres map[string]any) bool {
	for key, expected := range wheres {
		actual, exists := data[key]

		if !exists {
			return false
		}

		if operator, value, ok := whereComparison(expected); ok {
			if !compareValues(actual, operator, value) {
				return false
			}

			continue
		}

		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
			return false
		}
	}

	return true
}

func (e *CollectionEngine) matchesWhereIns(data map[string]any, whereIns map[string][]any) bool {
	for key, allowed := range whereIns {
		actual, exists := data[key]

		if !exists {
			return false
		}

		found := false
		actualStr := fmt.Sprintf("%v", actual)

		for _, v := range allowed {
			if fmt.Sprintf("%v", v) == actualStr {
				found = true

				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (e *CollectionEngine) matchesWhereNotIns(data map[string]any, whereNotIns map[string][]any) bool {
	for key, excluded := range whereNotIns {
		actual, exists := data[key]

		if !exists {
			continue
		}

		actualStr := fmt.Sprintf("%v", actual)

		for _, v := range excluded {
			if fmt.Sprintf("%v", v) == actualStr {
				return false
			}
		}
	}

	return true
}

func (e *CollectionEngine) sortModels(models []contract.Searchable, orders []contract.Order) {
	sort.SliceStable(models, func(i, j int) bool {
		for _, order := range orders {
			iData := models[i].ToSearchableArray()
			jData := models[j].ToSearchableArray()
			iVal := fmt.Sprintf("%v", iData[order.Column])
			jVal := fmt.Sprintf("%v", jData[order.Column])

			if iVal == jVal {
				continue
			}

			if order.Direction == "desc" {
				return iVal > jVal
			}

			return iVal < jVal
		}

		return false
	})
}

func whereComparison(expected any) (string, any, bool) {
	comparison, ok := expected.(map[string]any)

	if !ok {
		return "", nil, false
	}

	operator, _ := comparison["__operator"].(string)
	value, exists := comparison["__value"]

	if !exists {
		return "", nil, false
	}

	switch operator {
	case "=", "!=", "<>", ">", "<", ">=", "<=":
		return operator, value, true
	default:
		return "", nil, false
	}
}

func compareValues(actual any, operator string, expected any) bool {
	if actualNumber, expectedNumber, ok := comparableNumbers(actual, expected); ok {
		switch operator {
		case "=":
			return actualNumber == expectedNumber
		case "!=", "<>":
			return actualNumber != expectedNumber
		case ">":
			return actualNumber > expectedNumber
		case "<":
			return actualNumber < expectedNumber
		case ">=":
			return actualNumber >= expectedNumber
		case "<=":
			return actualNumber <= expectedNumber
		}
	}

	actualString := fmt.Sprintf("%v", actual)
	expectedString := fmt.Sprintf("%v", expected)

	switch operator {
	case "=":
		return actualString == expectedString
	case "!=", "<>":
		return actualString != expectedString
	case ">":
		return actualString > expectedString
	case "<":
		return actualString < expectedString
	case ">=":
		return actualString >= expectedString
	case "<=":
		return actualString <= expectedString
	default:
		return false
	}
}

func comparableNumbers(actual any, expected any) (float64, float64, bool) {
	actualNumber, actualOK := toFloat(actual)
	expectedNumber, expectedOK := toFloat(expected)

	return actualNumber, expectedNumber, actualOK && expectedOK
}

func toFloat(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float32:
		return float64(typed), true
	case float64:
		return typed, true
	case string:
		parsed, err := strconv.ParseFloat(typed, 64)

		return parsed, err == nil
	default:
		return 0, false
	}
}
