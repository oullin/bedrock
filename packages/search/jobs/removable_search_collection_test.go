package jobs_test

import (
	"reflect"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/bedrock/packages/search/jobs"
)

type customSearchKeyModel struct {
	testModel
	searchKey any
}

func (m *customSearchKeyModel) GetSearchKey() any {
	return m.searchKey
}

func (m *customSearchKeyModel) GetSearchKeyName() string {
	return "search_id"
}

func TestRemovableSearchCollectionGetQueueableIDs(t *testing.T) {
	t.Parallel()
	// RemovableSearchCollectionTest::test_get_queuable_ids
	models := []contract.Searchable{
		&testModel{id: 1, table: "posts"},
		&testModel{id: 2, table: "posts"},
	}

	collection := jobs.NewRemovableSearchCollection(models)

	if got := collection.GetQueueableIDs(); !reflect.DeepEqual(got, []any{1, 2}) {
		t.Fatalf("expected queueable IDs [1 2], got %v", got)
	}
}

func TestRemovableSearchCollectionGetQueueableIDsResolvesCustomSearchKeys(t *testing.T) {
	t.Parallel()
	// RemovableSearchCollectionTest::test_get_queuable_ids_resolves_custom_search_keys
	// SearchableTest::test_overridden_remove_from_search_is_dispatched
	models := []contract.Searchable{
		&customSearchKeyModel{
			testModel: testModel{id: 1, table: "chirps"},
			searchKey: "chirp-uuid-1",
		},
		&customSearchKeyModel{
			testModel: testModel{id: 2, table: "chirps"},
			searchKey: "chirp-uuid-2",
		},
	}

	collection := jobs.NewRemovableSearchCollection(models)

	if got := collection.GetQueueableIDs(); !reflect.DeepEqual(got, []any{"chirp-uuid-1", "chirp-uuid-2"}) {
		t.Fatalf("expected custom queueable IDs, got %v", got)
	}
}

func TestRemovableSearchCollectionReturnsSearchKeys(t *testing.T) {
	t.Parallel()
	// RemovableSearchCollectionTest::test_removeable_search_collection_returns_search_keys
	models := []contract.Searchable{
		&testModel{id: 10, table: "posts"},
		&customSearchKeyModel{
			testModel: testModel{id: 11, table: "chirps"},
			searchKey: "chirp-uuid-11",
		},
	}

	collection := jobs.NewRemovableSearchCollection(models)

	if got := collection.SearchKeys(); !reflect.DeepEqual(got, []any{10, "chirp-uuid-11"}) {
		t.Fatalf("expected search keys [10 chirp-uuid-11], got %v", got)
	}
}
