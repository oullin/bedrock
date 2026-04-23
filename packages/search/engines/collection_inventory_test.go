package engines_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/bedrock/packages/search"
	"github.com/bedrock/packages/search/engines"
)

func collectionInventoryModels() []contract.Searchable {
	return []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{
			"id": 1, "title": "Alpha Go", "status": "published", "views": 5, "created_at": "2024-01-01",
		}),
		newTestModelWithData(2, "posts", map[string]any{
			"id": 2, "title": "Beta Rust", "status": "draft", "views": 10, "created_at": "2024-03-01",
		}),
		newTestModelWithData(3, "posts", map[string]any{
			"id": 3, "title": "Gamma Go", "status": "published", "views": 15, "created_at": "2024-02-01",
		}),
	}
}

// CollectionEngineTest::test_it_can_retrieve_results_with_empty_search
// CollectionEngineTest::test_it_can_retrieve_results
// CollectionEngineTest::test_it_can_retrieve_results_matching_to_custom_searchable_data
func TestCollectionEngineInventoryRetrievesResults(t *testing.T) {
	t.Parallel()

	engine := engines.NewCollectionEngine()
	model := newTestModel(0, "posts")
	result, err := engine.Search(context.Background(), makeBuilder(model, "", collectionInventoryModels()))

	if err != nil {
		t.Fatalf("unexpected empty search error: %v", err)
	}

	if total := engine.GetTotalCount(result); total != 3 {
		t.Fatalf("expected all models for empty search, got %d", total)
	}

	result, err = engine.Search(context.Background(), makeBuilder(model, "Go", collectionInventoryModels()))

	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}

	if ids := engine.MapIds(result); len(ids) != 2 || ids[0] != 1 || ids[1] != 3 {
		t.Fatalf("expected Go matches [1 3], got %v", ids)
	}
}

// CollectionEngineTest::test_it_can_paginate_results
// CollectionEngineTest::test_limit_is_applied
func TestCollectionEngineInventoryPaginationAndLimit(t *testing.T) {
	t.Parallel()

	engine := engines.NewCollectionEngine()
	model := newTestModel(0, "posts")
	result, err := engine.Paginate(context.Background(), makeBuilder(model, "", collectionInventoryModels()), 2, 2)

	if err != nil {
		t.Fatalf("unexpected paginate error: %v", err)
	}

	if total := engine.GetTotalCount(result); total != 3 {
		t.Fatalf("expected total 3 before pagination, got %d", total)
	}

	models, err := engine.Map(context.Background(), result, model)

	if err != nil {
		t.Fatalf("unexpected map error: %v", err)
	}

	if len(models) != 1 || models[0].GetScoutKey() != 3 {
		t.Fatalf("expected second page to contain id 3, got %v", engine.MapIds(result))
	}

	limited, err := engine.Search(context.Background(), makeBuilder(model, "", collectionInventoryModels()).Take(1))

	if err != nil {
		t.Fatalf("unexpected limited search error: %v", err)
	}

	if ids := engine.MapIds(limited); len(ids) != 1 || ids[0] != 1 {
		t.Fatalf("expected limit to keep first id only, got %v", ids)
	}
}

// CollectionEngineTest::test_it_can_order_results
// CollectionEngineTest::test_it_can_order_by_latest_and_oldest
// CollectionEngineTest::test_it_can_order_by_custom_model_created_at_timestamp
func TestCollectionEngineInventoryOrdering(t *testing.T) {
	t.Parallel()

	engine := engines.NewCollectionEngine()
	model := newTestModel(0, "posts")
	ordered, err := engine.Search(context.Background(), makeBuilder(model, "", collectionInventoryModels()).OrderBy("title", "desc"))

	if err != nil {
		t.Fatalf("unexpected ordered search error: %v", err)
	}

	if ids := engine.MapIds(ordered); ids[0] != 3 || ids[2] != 1 {
		t.Fatalf("expected title desc order [3 ... 1], got %v", ids)
	}

	latest, err := engine.Search(context.Background(), makeBuilder(model, "", collectionInventoryModels()).Latest())

	if err != nil {
		t.Fatalf("unexpected latest search error: %v", err)
	}

	if ids := engine.MapIds(latest); ids[0] != 2 {
		t.Fatalf("expected latest created_at id 2 first, got %v", ids)
	}

	oldest, err := engine.Search(context.Background(), makeBuilder(model, "", collectionInventoryModels()).Oldest())

	if err != nil {
		t.Fatalf("unexpected oldest search error: %v", err)
	}

	if ids := engine.MapIds(oldest); ids[0] != 1 {
		t.Fatalf("expected oldest created_at id 1 first, got %v", ids)
	}
}

// CollectionEngineTest::test_it_can_filter_with_greater_than
// CollectionEngineTest::test_it_can_filter_with_less_than
// CollectionEngineTest::test_it_can_filter_with_greater_than_or_equal
// CollectionEngineTest::test_it_can_filter_with_less_than_or_equal
// CollectionEngineTest::test_it_can_filter_with_not_equal
// CollectionEngineTest::test_it_can_filter_with_multiple_where_comparisons
func TestCollectionEngineInventoryWhereComparisons(t *testing.T) {
	t.Parallel()

	engine := engines.NewCollectionEngine()
	model := newTestModel(0, "posts")

	tests := []struct {
		name     string
		builder  *search.Builder
		expected []any
	}{
		{name: "greater than", builder: makeBuilder(model, "", collectionInventoryModels()).Where("views", ">", 10), expected: []any{3}},
		{name: "less than", builder: makeBuilder(model, "", collectionInventoryModels()).Where("views", "<", 10), expected: []any{1}},
		{name: "greater than or equal", builder: makeBuilder(model, "", collectionInventoryModels()).Where("views", ">=", 10), expected: []any{2, 3}},
		{name: "less than or equal", builder: makeBuilder(model, "", collectionInventoryModels()).Where("views", "<=", 10), expected: []any{1, 2}},
		{name: "not equal", builder: makeBuilder(model, "", collectionInventoryModels()).Where("views", "!=", 10), expected: []any{1, 3}},
		{
			name:     "multiple comparisons",
			builder:  makeBuilder(model, "", collectionInventoryModels()).Where("views", ">", 5).Where("id", "!=", 2),
			expected: []any{3},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := engine.Search(context.Background(), tt.builder)

			if err != nil {
				t.Fatalf("unexpected search error: %v", err)
			}

			ids := engine.MapIds(result)

			if len(ids) != len(tt.expected) {
				t.Fatalf("expected ids %v, got %v", tt.expected, ids)
			}

			for i := range ids {
				if ids[i] != tt.expected[i] {
					t.Fatalf("expected ids %v, got %v", tt.expected, ids)
				}
			}
		})
	}
}
