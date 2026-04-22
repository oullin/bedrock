package jobs_test

import (
	"reflect"
	"testing"

	contract "github.com/bedrock/packages/contracts/scout"
	"github.com/bedrock/packages/scout/jobs"
)

type customScoutKeyModel struct {
	testModel
	scoutKey any
}

func (m *customScoutKeyModel) GetScoutKey() any {
	return m.scoutKey
}

func (m *customScoutKeyModel) GetScoutKeyName() string {
	return "scout_id"
}

func TestRemovableScoutCollectionGetQueueableIDs(t *testing.T) {
	t.Parallel()
	// RemovableScoutCollectionTest::test_get_queuable_ids
	models := []contract.Searchable{
		&testModel{id: 1, table: "posts"},
		&testModel{id: 2, table: "posts"},
	}

	collection := jobs.NewRemovableScoutCollection(models)

	if got := collection.GetQueueableIDs(); !reflect.DeepEqual(got, []any{1, 2}) {
		t.Fatalf("expected queueable IDs [1 2], got %v", got)
	}
}

func TestRemovableScoutCollectionGetQueueableIDsResolvesCustomScoutKeys(t *testing.T) {
	t.Parallel()
	// RemovableScoutCollectionTest::test_get_queuable_ids_resolves_custom_scout_keys
	// SearchableTest::test_overridden_remove_from_search_is_dispatched
	models := []contract.Searchable{
		&customScoutKeyModel{
			testModel: testModel{id: 1, table: "chirps"},
			scoutKey:  "chirp-uuid-1",
		},
		&customScoutKeyModel{
			testModel: testModel{id: 2, table: "chirps"},
			scoutKey:  "chirp-uuid-2",
		},
	}

	collection := jobs.NewRemovableScoutCollection(models)

	if got := collection.GetQueueableIDs(); !reflect.DeepEqual(got, []any{"chirp-uuid-1", "chirp-uuid-2"}) {
		t.Fatalf("expected custom queueable IDs, got %v", got)
	}
}

func TestRemovableScoutCollectionReturnsScoutKeys(t *testing.T) {
	t.Parallel()
	// RemovableScoutCollectionTest::test_removeable_scout_collection_returns_scout_keys
	models := []contract.Searchable{
		&testModel{id: 10, table: "posts"},
		&customScoutKeyModel{
			testModel: testModel{id: 11, table: "chirps"},
			scoutKey:  "chirp-uuid-11",
		},
	}

	collection := jobs.NewRemovableScoutCollection(models)

	if got := collection.ScoutKeys(); !reflect.DeepEqual(got, []any{10, "chirp-uuid-11"}) {
		t.Fatalf("expected scout keys [10 chirp-uuid-11], got %v", got)
	}
}
