package engines_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bedrock/packages/search"
	"github.com/bedrock/packages/search/engines"
)

// DatabaseEngineTest::test_it_can_retrieve_results_with_empty_search
// DatabaseEngineTest::test_it_does_not_add_search_where_clauses_with_empty_search
func TestDatabaseEngineInventoryEmptySearch(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver: "sqlite",
		selectRows: []map[string]any{
			{"id": int64(1), "title": "Alpha"},
			{"id": int64(2), "title": "Beta"},
		},
	}
	engine := engines.NewDatabaseEngine(&mockResolver{conn: conn})
	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	result, err := engine.Search(context.Background(), search.NewBuilder(model, ""))

	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}

	if total := engine.GetTotalCount(result); total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	if strings.Contains(conn.lastSQL, " like ?") || strings.Contains(conn.lastSQL, "match") {
		t.Fatalf("empty search should not add search clause, got SQL: %s", conn.lastSQL)
	}
}

// DatabaseEngineTest::test_it_adds_search_where_clauses_with_non_empty_search
// DatabaseEngineTest::test_it_can_retrieve_results
func TestDatabaseEngineInventorySearchClauseAndResults(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{{"id": int64(1), "title": "Alpha Go"}},
	}
	engine := engines.NewDatabaseEngine(&mockResolver{conn: conn})
	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	result, err := engine.Search(context.Background(), search.NewBuilder(model, "Go"))

	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}

	if ids := engine.MapIds(result); len(ids) != 1 || ids[0] != int64(1) {
		t.Fatalf("expected id 1, got %v", ids)
	}

	if !strings.Contains(conn.lastSQL, " like ?") {
		t.Fatalf("expected LIKE search clause, got SQL: %s", conn.lastSQL)
	}
}

// DatabaseEngineTest::test_it_can_paginate_results
// DatabaseEngineTest::test_limit_is_applied
func TestDatabaseEngineInventoryPaginationAndLimit(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{{"id": int64(1), "title": "Alpha"}},
	}
	engine := engines.NewDatabaseEngine(&mockResolver{conn: conn})
	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	_, err := engine.Paginate(context.Background(), search.NewBuilder(model, "Alpha"), 15, 3)

	if err != nil {
		t.Fatalf("unexpected paginate error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "limit 15 offset 30") {
		t.Fatalf("expected page limit/offset, got SQL: %s", conn.lastSQL)
	}

	_, err = engine.Search(context.Background(), search.NewBuilder(model, "Alpha").Take(5))

	if err != nil {
		t.Fatalf("unexpected limited search error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "limit 5") {
		t.Fatalf("expected builder limit, got SQL: %s", conn.lastSQL)
	}
}

// DatabaseEngineTest::test_it_can_order_results
// DatabaseEngineTest::test_it_uses_search_query
func TestDatabaseEngineInventoryOrderingAndSearchColumns(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{{"id": int64(1), "title": "Alpha"}},
	}
	engine := engines.NewDatabaseEngine(&mockResolver{conn: conn})
	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": "", "body": ""})
	builder := search.NewBuilder(model, "Alpha").
		OrderBy("created_at", "desc").
		WithOptions(map[string]any{"__fulltext_columns": []string{"title"}})

	if _, err := engine.Search(context.Background(), builder); err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "title like ?") || strings.Contains(conn.lastSQL, "body like ?") {
		t.Fatalf("expected custom search query columns, got SQL: %s", conn.lastSQL)
	}

	if !strings.Contains(conn.lastSQL, "order by created_at desc") {
		t.Fatalf("expected order clause, got SQL: %s", conn.lastSQL)
	}
}

// DatabaseEngineTest::test_it_can_filter_with_greater_than
// DatabaseEngineTest::test_it_can_filter_with_less_than
// DatabaseEngineTest::test_it_can_filter_with_greater_than_or_equal
// DatabaseEngineTest::test_it_can_filter_with_less_than_or_equal
// DatabaseEngineTest::test_it_can_filter_with_not_equal
// DatabaseEngineTest::test_it_can_filter_with_multiple_where_comparisons
func TestDatabaseEngineInventoryWhereComparisons(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		builder   *search.Builder
		wantSQL   []string
		wantBinds []any
	}{
		{
			name:      "greater than",
			builder:   search.NewBuilder(newTestModelWithData(0, "posts", map[string]any{"id": 0}), "").Where("views", ">", 10),
			wantSQL:   []string{"views > ?"},
			wantBinds: []any{10},
		},
		{
			name:      "less than",
			builder:   search.NewBuilder(newTestModelWithData(0, "posts", map[string]any{"id": 0}), "").Where("views", "<", 10),
			wantSQL:   []string{"views < ?"},
			wantBinds: []any{10},
		},
		{
			name:      "greater than or equal",
			builder:   search.NewBuilder(newTestModelWithData(0, "posts", map[string]any{"id": 0}), "").Where("views", ">=", 10),
			wantSQL:   []string{"views >= ?"},
			wantBinds: []any{10},
		},
		{
			name:      "less than or equal",
			builder:   search.NewBuilder(newTestModelWithData(0, "posts", map[string]any{"id": 0}), "").Where("views", "<=", 10),
			wantSQL:   []string{"views <= ?"},
			wantBinds: []any{10},
		},
		{
			name:      "not equal",
			builder:   search.NewBuilder(newTestModelWithData(0, "posts", map[string]any{"id": 0}), "").Where("views", "!=", 10),
			wantSQL:   []string{"views != ?"},
			wantBinds: []any{10},
		},
		{
			name:      "multiple comparisons",
			builder:   search.NewBuilder(newTestModelWithData(0, "posts", map[string]any{"id": 0}), "").Where("views", ">", 5).Where("id", "!=", 2),
			wantSQL:   []string{"views > ?", "id != ?"},
			wantBinds: []any{5, 2},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &mockConnection{driver: "sqlite", selectRows: []map[string]any{{"id": int64(1)}}}
			engine := engines.NewDatabaseEngine(&mockResolver{conn: conn})

			if _, err := engine.Search(context.Background(), tt.builder); err != nil {
				t.Fatalf("unexpected search error: %v", err)
			}

			for _, fragment := range tt.wantSQL {
				if !strings.Contains(conn.lastSQL, fragment) {
					t.Fatalf("expected SQL fragment %q in %s", fragment, conn.lastSQL)
				}
			}

			if len(conn.lastBinds) != len(tt.wantBinds) {
				t.Fatalf("expected bindings %v, got %v", tt.wantBinds, conn.lastBinds)
			}
		})
	}
}
