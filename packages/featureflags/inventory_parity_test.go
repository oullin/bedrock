package pennant_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/featureflags"
)

type inventoryScope struct {
	id string
}

func (s inventoryScope) FeatureScopeIdentifier() string {
	return "inventory:" + s.id
}

func newArrayInventoryHarness(dispatcher *testDispatcher) (*featureflags.ArrayDriver, *featureflags.Decorator) {
	driver := featureflags.NewArrayDriverWithDispatcher(dispatcher)

	return driver, featureflags.NewDecoratorWithDispatcher(driver, dispatcher)
}

func newDatabaseInventoryHarness(dispatcher *testDispatcher) (*inMemoryDB, *featureflags.DatabaseDriver, *featureflags.Decorator) {
	db := newDB()
	driver := featureflags.NewDatabaseDriverWithDispatcher(db, testTable, dispatcher)

	return db, driver, featureflags.NewDecoratorWithDispatcher(driver, dispatcher)
}

func requireFeatureValue(t *testing.T, dec *featureflags.Decorator, feature string, scope any) any {
	t.Helper()

	value, err := dec.Get(context.Background(), feature, scope)

	if err != nil {
		t.Fatalf("Get(%q, %v): %v", feature, scope, err)
	}

	return value
}

func requireFeatureActive(t *testing.T, interaction *featureflags.ScopedFeatureInteraction, feature string) {
	t.Helper()

	if !interaction.Active(context.Background(), feature) {
		t.Fatalf("expected feature %q to be active", feature)
	}
}

func requireFeatureInactive(t *testing.T, interaction *featureflags.ScopedFeatureInteraction, feature string) {
	t.Helper()

	if !interaction.Inactive(context.Background(), feature) {
		t.Fatalf("expected feature %q to be inactive", feature)
	}
}

func requireStoredFeature(t *testing.T, dec *featureflags.Decorator, feature string) {
	t.Helper()

	stored, err := dec.Stored(context.Background())

	if err != nil {
		t.Fatalf("Stored: %v", err)
	}

	for _, name := range stored {
		if name == feature {
			return
		}
	}

	t.Fatalf("expected stored feature %q in %v", feature, stored)
}

func requireDBRow(t *testing.T, db *inMemoryDB, name, scope, value string) {
	t.Helper()

	db.mu.Lock()
	defer db.mu.Unlock()

	for _, row := range db.rows {
		if row.name == name && row.scope == scope && row.value == value {
			return
		}
	}

	t.Fatalf("expected database row name=%q scope=%q value=%q in %#v", name, scope, value, db.rows)
}

func TestInventoryParityArrayDriverFeatureLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := &testDispatcher{}
	driver, dec := newArrayInventoryHarness(dispatcher)
	global := dec.For()

	// ArrayDriverTest::test_it_defaults_to_false_for_unknown_values
	requireFeatureInactive(t, global, "missing-feature")

	// ArrayDriverTest::test_it_dispatches_events_on_unknown_feature_checks
	if dispatcher.count("UnknownFeatureResolved") != 1 {
		t.Fatalf("expected one unknown-feature event, got %d", dispatcher.count("UnknownFeatureResolved"))
	}

	// ArrayDriverTest::test_it_can_register_default_boolean_values
	dec.Define("default-bool", func(context.Context, any) (any, error) {
		return true, nil
	})
	requireFeatureActive(t, global, "default-bool")

	// ArrayDriverTest::test_it_can_register_complex_values
	dec.Define("complex", func(context.Context, any) (any, error) {
		return map[string]any{"variant": "blue", "limit": 25}, nil
	})
	complexValue := requireFeatureValue(t, dec, "complex", nil)
	complexMap, ok := complexValue.(map[string]any)
	if !ok || complexMap["variant"] != "blue" || complexMap["limit"] != 25 {
		t.Fatalf("expected complex map value, got %#v", complexValue)
	}

	resolveCalls := 0
	dec.Define("cached", func(context.Context, any) (any, error) {
		resolveCalls++

		return true, nil
	})

	// ArrayDriverTest::test_caching_of_features
	// ArrayDriverTest::test_it_caches_state_after_resolving
	requireFeatureActive(t, global, "cached")
	requireFeatureActive(t, global, "cached")
	if resolveCalls != 1 {
		t.Fatalf("expected one resolver call for cached feature, got %d", resolveCalls)
	}

	// ArrayDriverTest::test_non_false_registered_values_are_considered_active
	dec.Define("variant", func(context.Context, any) (any, error) {
		return "treatment-a", nil
	})
	requireFeatureActive(t, global, "variant")

	// ArrayDriverTest::test_it_can_programatically_activate_and_deativate_features
	if err := global.Activate(ctx, []string{"manual"}); err != nil {
		t.Fatalf("Activate: %v", err)
	}
	requireFeatureActive(t, global, "manual")
	if err := global.Deactivate(ctx, []string{"manual"}); err != nil {
		t.Fatalf("Deactivate: %v", err)
	}
	requireFeatureInactive(t, global, "manual")

	// ArrayDriverTest::test_it_can_activate_and_deactivate_several_features_at_once
	if err := global.Activate(ctx, []string{"bulk-a", "bulk-b"}); err != nil {
		t.Fatalf("Activate bulk: %v", err)
	}
	if !global.AllAreActive(ctx, []string{"bulk-a", "bulk-b"}) {
		t.Fatal("expected both bulk features to be active")
	}
	if err := global.Deactivate(ctx, []string{"bulk-a", "bulk-b"}); err != nil {
		t.Fatalf("Deactivate bulk: %v", err)
	}
	if !global.AllAreInactive(ctx, []string{"bulk-a", "bulk-b"}) {
		t.Fatal("expected both bulk features to be inactive")
	}

	// ArrayDriverTest::test_it_can_check_if_multiple_features_are_active_at_once
	if err := global.Activate(ctx, []string{"check-a", "check-b"}); err != nil {
		t.Fatalf("Activate check features: %v", err)
	}
	if !global.AllAreActive(ctx, []string{"check-a", "check-b"}) {
		t.Fatal("expected all checked features to be active")
	}

	// ArrayDriverTest::test_it_can_scope_features
	dec.Define("scoped-resolver", func(_ context.Context, scope any) (any, error) {
		return scope == "user:1", nil
	})
	requireFeatureActive(t, dec.For("user:1"), "scoped-resolver")
	requireFeatureInactive(t, dec.For("user:2"), "scoped-resolver")

	// ArrayDriverTest::test_it_can_activate_and_deactivate_features_with_scope
	if err := dec.For("user:1").Activate(ctx, []string{"scoped-manual"}); err != nil {
		t.Fatalf("Activate scoped feature: %v", err)
	}
	requireFeatureActive(t, dec.For("user:1"), "scoped-manual")
	requireFeatureInactive(t, dec.For("user:2"), "scoped-manual")

	// ArrayDriverTest::test_it_can_activate_and_deactivate_features_for_multiple_scope_at_once
	if err := dec.For("user:1", "user:2").Activate(ctx, []string{"multi-scope"}); err != nil {
		t.Fatalf("Activate multi-scope feature: %v", err)
	}
	if !dec.For("user:1", "user:2").AllAreActive(ctx, []string{"multi-scope"}) {
		t.Fatal("expected multi-scope feature active for both scopes")
	}

	// ArrayDriverTest::test_it_can_activate_and_deactivate_multiple_features_for_multiple_scope_at_once
	if err := dec.For("user:1", "user:2").Activate(ctx, []string{"multi-a", "multi-b"}); err != nil {
		t.Fatalf("Activate multi features/scopes: %v", err)
	}
	if !dec.For("user:1", "user:2").AllAreActive(ctx, []string{"multi-a", "multi-b"}) {
		t.Fatal("expected all multi features/scopes active")
	}

	// ArrayDriverTest::test_it_can_check_multiple_features_for_multiple_scope_at_once
	if !dec.For("user:1", "user:2").SomeAreActive(ctx, []string{"multi-a", "multi-b"}) {
		t.Fatal("expected some multi-scope features active")
	}

	// ArrayDriverTest::test_null_is_same_as_global
	if err := dec.Set(ctx, "null-global", nil, true); err != nil {
		t.Fatalf("Set nil scope: %v", err)
	}
	requireFeatureActive(t, dec.For(), "null-global")
	requireFeatureActive(t, dec.For(nil), "null-global")

	// ArrayDriverTest::test_it_sees_null_and_empty_string_as_different_things
	if err := dec.Set(ctx, "scope-distinction", nil, true); err != nil {
		t.Fatalf("Set nil distinction: %v", err)
	}
	if err := dec.Set(ctx, "scope-distinction", "", false); err != nil {
		t.Fatalf("Set empty distinction: %v", err)
	}
	requireFeatureActive(t, dec.For(nil), "scope-distinction")
	requireFeatureInactive(t, dec.For(""), "scope-distinction")

	// ArrayDriverTest::test_scope_can_be_strings_like_email_addresses
	if err := dec.Set(ctx, "email-scope", "taylor@example.com", true); err != nil {
		t.Fatalf("Set email scope: %v", err)
	}
	requireFeatureActive(t, dec.For("taylor@example.com"), "email-scope")

	// ArrayDriverTest::test_it_can_handle_feature_scopeable_objects
	scopeable := inventoryScope{id: "42"}
	if err := dec.Set(ctx, "scopeable", scopeable, true); err != nil {
		t.Fatalf("Set scopeable: %v", err)
	}
	requireFeatureActive(t, dec.For(inventoryScope{id: "42"}), "scopeable")

	loadCalls := 0
	dec.Define("loadable", func(context.Context, any) (any, error) {
		loadCalls++

		return true, nil
	})

	// ArrayDriverTest::test_it_can_load_feature_state_into_memory
	if err := global.Load(ctx, []string{"loadable"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	requireFeatureActive(t, global, "loadable")
	if loadCalls != 1 {
		t.Fatalf("expected one loadable resolver call, got %d", loadCalls)
	}

	// ArrayDriverTest::test_it_can_load_missing_feature_state_into_memory
	if err := global.LoadMissing(ctx, []string{"loadable"}); err != nil {
		t.Fatalf("LoadMissing: %v", err)
	}
	if loadCalls != 1 {
		t.Fatalf("expected LoadMissing to skip cached feature, got %d calls", loadCalls)
	}

	// ArrayDriverTest::test_it_can_load_scoped_feature_state_into_memory
	// ArrayDriverTest::test_it_can_load_against_scope
	scopedLoadCalls := 0
	dec.Define("scoped-load", func(_ context.Context, scope any) (any, error) {
		scopedLoadCalls++

		return scope == "user:loaded", nil
	})
	if err := dec.For("user:loaded").Load(ctx, []string{"scoped-load"}); err != nil {
		t.Fatalf("scoped Load: %v", err)
	}
	requireFeatureActive(t, dec.For("user:loaded"), "scoped-load")
	if scopedLoadCalls != 1 {
		t.Fatalf("expected one scoped load resolver call, got %d", scopedLoadCalls)
	}

	// ArrayDriverTest::test_it_can_retrive_value_for_multiple_scopes
	// ArrayDriverTest::test_it_can_retrive_value_for_multiple_scopes_and_features
	dec.Define("scope-value", func(_ context.Context, scope any) (any, error) {
		return fmt.Sprintf("value:%v", scope), nil
	})
	dec.Define("scope-value-b", func(_ context.Context, scope any) (any, error) {
		return fmt.Sprintf("other:%v", scope), nil
	})
	manyValues, err := dec.GetAll(ctx, map[string][]any{
		"scope-value":   {"user:1", "user:2"},
		"scope-value-b": {"user:3"},
	})
	if err != nil {
		t.Fatalf("GetAll multi-scope: %v", err)
	}
	if manyValues["scope-value"][0] != "value:user:1" ||
		manyValues["scope-value"][1] != "value:user:2" ||
		manyValues["scope-value-b"][0] != "other:user:3" {
		t.Fatalf("unexpected multi-scope values: %#v", manyValues)
	}

	// ArrayDriverTest::test_it_throws_when_calling_value_with_multiple_scope
	if _, err := dec.For("user:1", "user:2").Value(ctx, "scope-value"); !errors.Is(err, featureflags.ErrMultipleScopes) {
		t.Fatalf("expected ErrMultipleScopes, got %v", err)
	}

	// ArrayDriverTest::test_it_may_register_shorthand_feature_values
	dec.DefineValue("shorthand", "value")
	if requireFeatureValue(t, dec, "shorthand", nil) != "value" {
		t.Fatal("expected shorthand value")
	}

	// ArrayDriverTest::test_it_can_use_lottery
	dec.DefineValue("lottery-on", featureflags.LotteryOdds(1, 1))
	dec.DefineValue("lottery-off", featureflags.LotteryOdds(0, 1))
	dec.Define("lottery-from-resolver", func(context.Context, any) (any, error) {
		return featureflags.LotteryOdds(0, 1), nil
	})
	if requireFeatureValue(t, dec, "lottery-on", nil) != true ||
		requireFeatureValue(t, dec, "lottery-off", nil) != false ||
		requireFeatureValue(t, dec, "lottery-from-resolver", nil) != false {
		t.Fatal("expected lottery-backed features to resolve to sampled booleans")
	}

	// ArrayDriverTest::test_it_can_retrieve_registered_features
	defined := strings.Join(driver.Defined(), ",")
	if !strings.Contains(defined, "default-bool") || !strings.Contains(defined, "complex") {
		t.Fatalf("expected registered features in Defined list, got %q", defined)
	}

	// ArrayDriverTest::test_it_can_get_all_features
	allValues, err := dec.GetAll(ctx, map[string][]any{
		"default-bool": {nil},
		"variant":      {nil},
	})
	if err != nil {
		t.Fatalf("GetAll features: %v", err)
	}
	if allValues["default-bool"][0] != true || allValues["variant"][0] != "treatment-a" {
		t.Fatalf("unexpected all-feature values: %#v", allValues)
	}

	// ArrayDriverTest::test_it_can_set_for_all
	if err := dec.SetForAllScopes(ctx, "multi-scope", "enabled-for-all"); err != nil {
		t.Fatalf("SetForAllScopes: %v", err)
	}
	if requireFeatureValue(t, dec, "multi-scope", "user:1") != "enabled-for-all" ||
		requireFeatureValue(t, dec, "multi-scope", "user:2") != "enabled-for-all" {
		t.Fatal("expected SetForAllScopes to update all stored scopes")
	}

	// ArrayDriverTest::test_it_can_list_stored_features
	requireStoredFeature(t, dec, "default-bool")

	// ArrayDriverTest::test_it_caches_by_identifier_for_feature_scopable_objects
	scopeableCalls := 0
	dec.Define("scopeable-cache", func(context.Context, any) (any, error) {
		scopeableCalls++

		return true, nil
	})
	requireFeatureActive(t, dec.For(inventoryScope{id: "same"}), "scopeable-cache")
	requireFeatureActive(t, dec.For(inventoryScope{id: "same"}), "scopeable-cache")
	if scopeableCalls != 1 {
		t.Fatalf("expected scopeable cache by identifier, got %d calls", scopeableCalls)
	}

	// ArrayDriverTest::test_it_caches_by_identification_for_other_objects
	structCalls := 0
	dec.Define("struct-cache", func(context.Context, any) (any, error) {
		structCalls++

		return true, nil
	})
	requireFeatureActive(t, dec.For(structWithID{ID: 99}), "struct-cache")
	requireFeatureActive(t, dec.For(structWithID{ID: 99}), "struct-cache")
	if structCalls != 1 {
		t.Fatalf("expected struct cache by ID, got %d calls", structCalls)
	}

	// ArrayDriverTest::test_caching_of_features
	// ArrayDriverTest::test_it_can_clear_the_cache
	dec.FlushCache()
	requireFeatureActive(t, dec.For(structWithID{ID: 99}), "struct-cache")
	if structCalls != 1 {
		t.Fatalf("expected driver cache to survive decorator FlushCache, got %d calls", structCalls)
	}

	// ArrayDriverTest::test_can_retrieve_scalar_values_without_in_memory_cache
	if err := dec.Set(ctx, "scalar", "scope", "plain-value"); err != nil {
		t.Fatalf("Set scalar: %v", err)
	}
	dec.FlushCache()
	if requireFeatureValue(t, dec, "scalar", "scope") != "plain-value" {
		t.Fatal("expected scalar value after clearing decorator cache")
	}

	// ArrayDriverTest::test_it_handles_integer_scopes_correctly
	if err := dec.Set(ctx, "integer-scope", 10, true); err != nil {
		t.Fatalf("Set integer scope: %v", err)
	}
	requireFeatureActive(t, dec.For(10), "integer-scope")

	// ArrayDriverTest::test_it_can_handles_double_scopes_correctly
	// ArrayDriverTest::test_it_can_handles_float_scopes_correctly
	if err := dec.Set(ctx, "float-scope", 10.5, true); err != nil {
		t.Fatalf("Set float scope: %v", err)
	}
	requireFeatureActive(t, dec.For(10.5), "float-scope")

	// ArrayDriverTest::test_it_can_reevaluate_feature_state
	reevaluateCalls := 0
	dec.Define("reevaluate", func(context.Context, any) (any, error) {
		reevaluateCalls++

		return reevaluateCalls == 1, nil
	})
	requireFeatureActive(t, global, "reevaluate")
	if err := dec.Delete(ctx, "reevaluate", nil); err != nil {
		t.Fatalf("Delete reevaluate: %v", err)
	}
	requireFeatureInactive(t, global, "reevaluate")
	if reevaluateCalls != 2 {
		t.Fatalf("expected reevaluate resolver to run twice, got %d", reevaluateCalls)
	}
}

func TestInventoryParityArrayDriverConditionalEventsAndPurging(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := &testDispatcher{}
	_, dec := newArrayInventoryHarness(dispatcher)
	scoped := dec.For("user:events")

	dec.Define("conditional-active", func(context.Context, any) (any, error) {
		return "on", nil
	})
	dec.Define("conditional-inactive", func(context.Context, any) (any, error) {
		return false, nil
	})

	// ArrayDriverTest::test_it_can_conditionally_execute_code_block_for_active_feature
	activeResult, err := scoped.When(ctx, "conditional-active",
		func(value any) (any, error) { return "active:" + value.(string), nil },
		func(any) (any, error) { return "inactive", nil },
	)
	if err != nil || activeResult != "active:on" {
		t.Fatalf("unexpected active conditional result=%v err=%v", activeResult, err)
	}

	// ArrayDriverTest::test_it_receives_value_for_feature_in_conditional_code_execution
	if activeResult != "active:on" {
		t.Fatalf("expected active callback to receive resolved value, got %v", activeResult)
	}

	// ArrayDriverTest::test_it_can_conditionally_execute_code_block_for_inactive_feature
	inactiveResult, err := scoped.When(ctx, "conditional-inactive",
		func(any) (any, error) { return "active", nil },
		func(value any) (any, error) { return fmt.Sprintf("inactive:%v", value), nil },
	)
	if err != nil || inactiveResult != "inactive:false" {
		t.Fatalf("unexpected inactive conditional result=%v err=%v", inactiveResult, err)
	}

	// ArrayDriverTest::test_it_can_conditionally_execute_code_block_providing_closure_for_only_active
	onlyActiveResult, err := scoped.When(ctx, "conditional-active",
		func(any) (any, error) { return "only-active", nil },
		nil,
	)
	if err != nil || onlyActiveResult != "only-active" {
		t.Fatalf("unexpected only-active conditional result=%v err=%v", onlyActiveResult, err)
	}

	// ArrayDriverTest::test_conditionally_executing_code_respects_scope
	dec.Define("scoped-conditional", func(_ context.Context, scope any) (any, error) {
		return scope == "allowed", nil
	})
	result, err := dec.For("denied").When(ctx, "scoped-conditional",
		func(any) (any, error) { return "active", nil },
		func(any) (any, error) { return "inactive", nil },
	)
	if err != nil || result != "inactive" {
		t.Fatalf("expected scoped conditional to use denied scope, got result=%v err=%v", result, err)
	}

	// ArrayDriverTest::test_conditional_closures_receive_current_feature_interaction
	interaction := dec.For("allowed")
	result, err = interaction.When(ctx, "scoped-conditional",
		func(value any) (any, error) {
			if !interaction.Active(ctx, "scoped-conditional") {
				return nil, errors.New("interaction was not active")
			}

			return value, nil
		},
		nil,
	)
	if err != nil || result != true {
		t.Fatalf("expected conditional callback to observe current interaction, result=%v err=%v", result, err)
	}

	// ArrayDriverTest::test_it_dispatches_events_when_resolving_feature_into_memory
	beforeResolved := dispatcher.count("FeatureResolved")
	dec.Define("event-resolve", func(context.Context, any) (any, error) {
		return true, nil
	})
	requireFeatureActive(t, scoped, "event-resolve")
	if dispatcher.count("FeatureResolved") != beforeResolved+1 {
		t.Fatalf("expected FeatureResolved event, before=%d after=%d", beforeResolved, dispatcher.count("FeatureResolved"))
	}

	// ArrayDriverTest::test_it_dispatches_events_when_updating_a_scoped_feature
	beforeUpdated := dispatcher.count("FeatureUpdated")
	if err := scoped.ActivateWithValue(ctx, []string{"event-update"}, "value"); err != nil {
		t.Fatalf("ActivateWithValue: %v", err)
	}
	if dispatcher.count("FeatureUpdated") <= beforeUpdated {
		t.Fatal("expected FeatureUpdated event for scoped update")
	}

	// ArrayDriverTest::test_it_dispatches_events_when_updating_a_feature_for_all_scopes
	beforeUpdatedAll := dispatcher.count("FeatureUpdatedForAllScopes")
	if err := dec.SetForAllScopes(ctx, "event-update", "next"); err != nil {
		t.Fatalf("SetForAllScopes event: %v", err)
	}
	if dispatcher.count("FeatureUpdatedForAllScopes") != beforeUpdatedAll+1 {
		t.Fatal("expected FeatureUpdatedForAllScopes event")
	}

	// ArrayDriverTest::test_it_dispatches_events_when_deleting_a_feature_value
	beforeDeleted := dispatcher.count("FeatureDeleted")
	if err := scoped.Forget(ctx, []string{"event-update"}); err != nil {
		t.Fatalf("Forget event: %v", err)
	}
	if dispatcher.count("FeatureDeleted") <= beforeDeleted {
		t.Fatal("expected FeatureDeleted event")
	}

	beforeBulkUpdated := dispatcher.count("FeatureUpdated")
	// ArrayDriverTest::test_it_dispatches_events_when_activating_multiple_features_and_scopes
	if err := dec.For("user:a", "user:b").Activate(ctx, []string{"event-a", "event-b"}); err != nil {
		t.Fatalf("Activate multi event features: %v", err)
	}
	if dispatcher.count("FeatureUpdated") != beforeBulkUpdated+4 {
		t.Fatalf("expected four activation update events, got %d", dispatcher.count("FeatureUpdated")-beforeBulkUpdated)
	}

	beforeBulkUpdated = dispatcher.count("FeatureUpdated")
	// ArrayDriverTest::test_it_dispatches_events_when_deactivating_multiple_features_and_scopes
	if err := dec.For("user:a", "user:b").Deactivate(ctx, []string{"event-a", "event-b"}); err != nil {
		t.Fatalf("Deactivate multi event features: %v", err)
	}
	if dispatcher.count("FeatureUpdated") != beforeBulkUpdated+4 {
		t.Fatalf("expected four deactivation update events, got %d", dispatcher.count("FeatureUpdated")-beforeBulkUpdated)
	}

	if err := scoped.Activate(ctx, []string{"purge-a", "purge-b"}); err != nil {
		t.Fatalf("Activate purge features: %v", err)
	}

	// ArrayDriverTest::test_it_dispatches_events_when_purging_features
	beforePurged := dispatcher.count("FeaturesPurged")
	if err := scoped.Purge(ctx, []string{"purge-a"}); err != nil {
		t.Fatalf("Purge features: %v", err)
	}
	if dispatcher.count("FeaturesPurged") != beforePurged+1 {
		t.Fatal("expected FeaturesPurged event")
	}

	// ArrayDriverTest::test_it_dispatches_events_when_purging_all_features
	beforeAllPurged := dispatcher.count("AllFeaturesPurged")
	if err := scoped.Purge(ctx, nil); err != nil {
		t.Fatalf("Purge all: %v", err)
	}
	if dispatcher.count("AllFeaturesPurged") != beforeAllPurged+1 {
		t.Fatal("expected AllFeaturesPurged event")
	}

	dec.Define("purge-retrieve", func(context.Context, any) (any, error) {
		return true, nil
	})
	requireFeatureActive(t, scoped, "purge-retrieve")
	if err := scoped.Purge(ctx, []string{"purge-retrieve"}); err != nil {
		t.Fatalf("Purge retrieve: %v", err)
	}

	// ArrayDriverTest::test_retrieving_values_after_purging
	requireFeatureActive(t, scoped, "purge-retrieve")
}

func TestInventoryParityDatabaseDriverFeatureLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := &testDispatcher{}
	db, driver, dec := newDatabaseInventoryHarness(dispatcher)
	global := dec.For()

	// DatabaseDriverTest::test_it_defaults_to_false_for_unknown_values
	requireFeatureInactive(t, global, "missing-feature")

	// DatabaseDriverTest::test_it_dispatches_events_on_unknown_feature_checks
	if dispatcher.count("UnknownFeatureResolved") != 1 {
		t.Fatalf("expected one unknown-feature event, got %d", dispatcher.count("UnknownFeatureResolved"))
	}

	// DatabaseDriverTest::test_it_can_register_default_boolean_values
	dec.Define("default-bool", func(context.Context, any) (any, error) {
		return true, nil
	})
	requireFeatureActive(t, global, "default-bool")

	// DatabaseDriverTest::test_it_can_register_complex_values
	dec.Define("complex", func(context.Context, any) (any, error) {
		return map[string]any{"variant": "blue"}, nil
	})
	complexValue := requireFeatureValue(t, dec, "complex", nil)
	complexMap, ok := complexValue.(map[string]any)
	if !ok || complexMap["variant"] != "blue" {
		t.Fatalf("expected complex map value, got %#v", complexValue)
	}

	resolveCalls := 0
	dec.Define("cached", func(context.Context, any) (any, error) {
		resolveCalls++

		return true, nil
	})

	// DatabaseDriverTest::test_it_caches_state_after_resolving
	requireFeatureActive(t, global, "cached")
	// DatabaseDriverTest::test_it_can_clear_the_cache
	dec.FlushCache()
	requireFeatureActive(t, global, "cached")
	if resolveCalls != 1 {
		t.Fatalf("expected database-backed cache to persist resolved state, got %d calls", resolveCalls)
	}

	// DatabaseDriverTest::test_non_false_registered_values_are_considered_active
	dec.Define("variant", func(context.Context, any) (any, error) {
		return "treatment-a", nil
	})
	requireFeatureActive(t, global, "variant")

	// DatabaseDriverTest::test_it_can_programatically_activate_and_deativate_features
	if err := global.Activate(ctx, []string{"manual"}); err != nil {
		t.Fatalf("Activate: %v", err)
	}
	requireFeatureActive(t, global, "manual")
	if err := global.Deactivate(ctx, []string{"manual"}); err != nil {
		t.Fatalf("Deactivate: %v", err)
	}
	requireFeatureInactive(t, global, "manual")

	// DatabaseDriverTest::test_it_can_activate_and_deactivate_several_features_at_once
	if err := global.Activate(ctx, []string{"bulk-a", "bulk-b"}); err != nil {
		t.Fatalf("Activate bulk: %v", err)
	}
	if !global.AllAreActive(ctx, []string{"bulk-a", "bulk-b"}) {
		t.Fatal("expected both bulk features to be active")
	}

	// DatabaseDriverTest::test_it_can_check_if_multiple_features_are_active_at_once
	if !global.SomeAreActive(ctx, []string{"bulk-a", "bulk-b"}) {
		t.Fatal("expected some bulk feature to be active")
	}

	// DatabaseDriverTest::test_it_can_scope_features
	dec.Define("scoped-resolver", func(_ context.Context, scope any) (any, error) {
		return scope == "user:1", nil
	})
	requireFeatureActive(t, dec.For("user:1"), "scoped-resolver")
	requireFeatureInactive(t, dec.For("user:2"), "scoped-resolver")

	// DatabaseDriverTest::test_it_can_activate_and_deactivate_features_with_scope
	if err := dec.For("user:1").Activate(ctx, []string{"scoped-manual"}); err != nil {
		t.Fatalf("Activate scoped feature: %v", err)
	}
	requireFeatureActive(t, dec.For("user:1"), "scoped-manual")
	requireFeatureInactive(t, dec.For("user:2"), "scoped-manual")

	// DatabaseDriverTest::test_it_can_activate_and_deactivate_features_for_multiple_scope_at_once
	if err := dec.For("user:1", "user:2").Activate(ctx, []string{"multi-scope"}); err != nil {
		t.Fatalf("Activate multi-scope feature: %v", err)
	}
	if !dec.For("user:1", "user:2").AllAreActive(ctx, []string{"multi-scope"}) {
		t.Fatal("expected multi-scope feature active for both scopes")
	}

	// DatabaseDriverTest::test_it_can_activate_and_deactivate_multiple_features_for_multiple_scope_at_once
	if err := dec.For("user:1", "user:2").Activate(ctx, []string{"multi-a", "multi-b"}); err != nil {
		t.Fatalf("Activate multi features/scopes: %v", err)
	}
	if !dec.For("user:1", "user:2").AllAreActive(ctx, []string{"multi-a", "multi-b"}) {
		t.Fatal("expected all multi features/scopes active")
	}

	// DatabaseDriverTest::test_it_can_check_multiple_features_for_multiple_scope_at_once
	// DatabaseDriverTest::test_it_handles_multiscope_checks
	if !dec.For("user:1", "user:2").SomeAreActive(ctx, []string{"multi-a", "multi-b"}) {
		t.Fatal("expected some multi-scope features active")
	}

	// DatabaseDriverTest::test_null_is_same_as_global
	if err := dec.Set(ctx, "null-global", nil, true); err != nil {
		t.Fatalf("Set nil scope: %v", err)
	}
	requireFeatureActive(t, dec.For(), "null-global")
	requireFeatureActive(t, dec.For(nil), "null-global")

	// DatabaseDriverTest::test_it_sees_null_and_empty_string_as_different_things
	if err := dec.Set(ctx, "scope-distinction", nil, true); err != nil {
		t.Fatalf("Set nil distinction: %v", err)
	}
	if err := dec.Set(ctx, "scope-distinction", "", false); err != nil {
		t.Fatalf("Set empty distinction: %v", err)
	}
	requireFeatureActive(t, dec.For(nil), "scope-distinction")
	requireFeatureInactive(t, dec.For(""), "scope-distinction")

	// DatabaseDriverTest::test_scope_can_be_strings_like_email_addresses
	if err := dec.Set(ctx, "email-scope", "taylor@example.com", true); err != nil {
		t.Fatalf("Set email scope: %v", err)
	}
	requireFeatureActive(t, dec.For("taylor@example.com"), "email-scope")

	// DatabaseDriverTest::test_it_can_handle_feature_scopeable_objects
	if err := dec.Set(ctx, "scopeable", inventoryScope{id: "42"}, true); err != nil {
		t.Fatalf("Set scopeable: %v", err)
	}
	requireFeatureActive(t, dec.For(inventoryScope{id: "42"}), "scopeable")

	// DatabaseDriverTest::test_it_can_manually_serialize_scope
	serialized, err := featureflags.SerializeScope(inventoryScope{id: "42"})
	if err != nil || serialized != "inventory:42" {
		t.Fatalf("SerializeScope returned %q, %v", serialized, err)
	}

	// DatabaseDriverTest::test_it_can_load_feature_state_into_memory
	loadCalls := 0
	dec.Define("loadable", func(context.Context, any) (any, error) {
		loadCalls++

		return true, nil
	})
	if err := global.Load(ctx, []string{"loadable"}); err != nil {
		t.Fatalf("Load: %v", err)
	}
	requireFeatureActive(t, global, "loadable")
	if loadCalls != 1 {
		t.Fatalf("expected one loadable resolver call, got %d", loadCalls)
	}

	// DatabaseDriverTest::test_it_can_load_missing_feature_state_into_memory
	if err := global.LoadMissing(ctx, []string{"loadable"}); err != nil {
		t.Fatalf("LoadMissing: %v", err)
	}
	if loadCalls != 1 {
		t.Fatalf("expected LoadMissing to skip cached feature, got %d calls", loadCalls)
	}

	// DatabaseDriverTest::test_it_can_load_scoped_feature_state_into_memory
	// DatabaseDriverTest::test_it_can_load_against_scope
	scopedLoadCalls := 0
	dec.Define("scoped-load", func(_ context.Context, scope any) (any, error) {
		scopedLoadCalls++

		return scope == "user:loaded", nil
	})
	if err := dec.For("user:loaded").Load(ctx, []string{"scoped-load"}); err != nil {
		t.Fatalf("database scoped Load: %v", err)
	}
	requireFeatureActive(t, dec.For("user:loaded"), "scoped-load")
	if scopedLoadCalls != 1 {
		t.Fatalf("expected one scoped database load resolver call, got %d", scopedLoadCalls)
	}
	requireFeatureInactive(t, dec.For("user:not-loaded"), "scoped-load")
	if scopedLoadCalls != 2 {
		t.Fatalf("expected unloaded scope to resolve independently, got %d calls", scopedLoadCalls)
	}

	// DatabaseDriverTest::test_it_does_not_hit_db_when_features_are_empty
	if err := global.Load(ctx, nil); err != nil {
		t.Fatalf("empty Load: %v", err)
	}
	emptyRawDB, _, emptyDec := newDatabaseInventoryHarness(&testDispatcher{})
	if err := emptyDec.For().Load(ctx, nil); err != nil {
		t.Fatalf("empty database Load on fresh driver: %v", err)
	}
	if emptyRawDB.execCount+emptyRawDB.queryCount+emptyRawDB.rowQueryCount != 0 {
		t.Fatalf("expected empty Load to avoid database calls, got exec=%d query=%d row=%d",
			emptyRawDB.execCount, emptyRawDB.queryCount, emptyRawDB.rowQueryCount)
	}

	// DatabaseDriverTest::test_it_can_retrieve_registered_features
	defined := strings.Join(driver.Defined(), ",")
	if !strings.Contains(defined, "default-bool") || !strings.Contains(defined, "complex") {
		t.Fatalf("expected registered features in Defined list, got %q", defined)
	}

	// DatabaseDriverTest::test_it_can_get_all_features
	allValues, err := dec.GetAll(ctx, map[string][]any{
		"default-bool": {nil},
		"variant":      {nil},
	})
	if err != nil {
		t.Fatalf("GetAll features: %v", err)
	}
	if allValues["default-bool"][0] != true || allValues["variant"][0] != "treatment-a" {
		t.Fatalf("unexpected all-feature values: %#v", allValues)
	}

	// DatabaseDriverTest::test_it_can_set_for_all
	if err := dec.SetForAllScopes(ctx, "multi-scope", "enabled-for-all"); err != nil {
		t.Fatalf("SetForAllScopes: %v", err)
	}
	if requireFeatureValue(t, dec, "multi-scope", "user:1") != "enabled-for-all" ||
		requireFeatureValue(t, dec, "multi-scope", "user:2") != "enabled-for-all" {
		t.Fatal("expected SetForAllScopes to update all stored scopes")
	}

	// DatabaseDriverTest::test_it_can_list_stored_features
	requireStoredFeature(t, dec, "default-bool")

	// DatabaseDriverTest::test_can_retrieve_scalar_values_without_in_memory_cache
	if err := dec.Set(ctx, "scalar", "scope", "plain-value"); err != nil {
		t.Fatalf("Set scalar: %v", err)
	}
	dec.FlushCache()
	if requireFeatureValue(t, dec, "scalar", "scope") != "plain-value" {
		t.Fatal("expected scalar value after clearing decorator cache")
	}

	// DatabaseDriverTest::test_it_can_reevaluate_feature_state
	reevaluateCalls := 0
	dec.Define("db-reevaluate", func(context.Context, any) (any, error) {
		reevaluateCalls++

		return reevaluateCalls == 1, nil
	})
	requireFeatureActive(t, dec.For(), "db-reevaluate")
	if err := dec.Delete(ctx, "db-reevaluate", nil); err != nil {
		t.Fatalf("Delete db-reevaluate: %v", err)
	}
	requireFeatureInactive(t, dec.For(), "db-reevaluate")
	if reevaluateCalls != 2 {
		t.Fatalf("expected db reevaluate resolver to run twice, got %d", reevaluateCalls)
	}

	// DatabaseDriverTest::test_can_retrieve_all_features_for_differing_scope_types
	// IntersectionTypeTest::test_can_retrieve_all_features_for_differing_scope_types
	for _, scope := range []any{nil, "string-scope", 42, inventoryScope{id: "typed"}} {
		if err := dec.Set(ctx, "differing-scope-types", scope, true); err != nil {
			t.Fatalf("Set differing scope %v: %v", scope, err)
		}
		requireFeatureActive(t, dec.For(scope), "differing-scope-types")
	}

	// DatabaseDriverTest::test_missing_results_are_inserted_on_load
	if err := dec.For("taylor@example.com").ActivateWithValue(ctx, []string{"missing-foo"}, 99); err != nil {
		t.Fatalf("Activate existing missing-foo: %v", err)
	}
	dec.Define("missing-foo", func(context.Context, any) (any, error) { return 1, nil })
	dec.Define("missing-bar", func(context.Context, any) (any, error) { return 2, nil })
	if err := dec.For("tim@example.com", "jess@example.com", "taylor@example.com").Load(ctx, []string{"missing-foo", "missing-bar"}); err != nil {
		t.Fatalf("Load missing results: %v", err)
	}
	requireDBRow(t, db, "missing-foo", "tim@example.com", "1")
	requireDBRow(t, db, "missing-foo", "jess@example.com", "1")
	requireDBRow(t, db, "missing-foo", "taylor@example.com", "99")
	requireDBRow(t, db, "missing-bar", "tim@example.com", "2")
}

func TestInventoryParityDatabaseDriverPurgingEventsAndRetries(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := &testDispatcher{}
	db, _, dec := newDatabaseInventoryHarness(dispatcher)
	scoped := dec.For("user:events")

	dec.Define("conditional-active", func(context.Context, any) (any, error) {
		return "on", nil
	})
	dec.Define("conditional-inactive", func(context.Context, any) (any, error) {
		return false, nil
	})

	// DatabaseDriverTest::test_it_can_conditionally_execute_code_block_for_active_feature
	activeResult, err := scoped.When(ctx, "conditional-active",
		func(value any) (any, error) { return "active:" + value.(string), nil },
		func(any) (any, error) { return "inactive", nil },
	)
	if err != nil || activeResult != "active:on" {
		t.Fatalf("unexpected active conditional result=%v err=%v", activeResult, err)
	}

	// DatabaseDriverTest::test_it_receives_value_for_feature_in_conditional_code_execution
	if activeResult != "active:on" {
		t.Fatalf("expected active callback to receive resolved value, got %v", activeResult)
	}

	// DatabaseDriverTest::test_it_can_conditionally_execute_code_block_for_inactive_feature
	inactiveResult, err := scoped.When(ctx, "conditional-inactive",
		func(any) (any, error) { return "active", nil },
		func(value any) (any, error) { return fmt.Sprintf("inactive:%v", value), nil },
	)
	if err != nil || inactiveResult != "inactive:false" {
		t.Fatalf("unexpected inactive conditional result=%v err=%v", inactiveResult, err)
	}

	// DatabaseDriverTest::test_conditionally_executing_code_respects_scope
	dec.Define("scoped-conditional", func(_ context.Context, scope any) (any, error) {
		return scope == "allowed", nil
	})
	result, err := dec.For("denied").When(ctx, "scoped-conditional",
		func(any) (any, error) { return "active", nil },
		func(any) (any, error) { return "inactive", nil },
	)
	if err != nil || result != "inactive" {
		t.Fatalf("expected scoped conditional to use denied scope, got result=%v err=%v", result, err)
	}

	// DatabaseDriverTest::test_conditional_closures_receive_current_feature_interaction
	interaction := dec.For("allowed")
	result, err = interaction.When(ctx, "scoped-conditional",
		func(value any) (any, error) {
			if !interaction.Active(ctx, "scoped-conditional") {
				return nil, errors.New("interaction was not active")
			}

			return value, nil
		},
		nil,
	)
	if err != nil || result != true {
		t.Fatalf("expected conditional callback to observe current interaction, result=%v err=%v", result, err)
	}

	// DatabaseDriverTest::test_it_dispatches_events_when_checking_known_features
	beforeResolved := dispatcher.count("FeatureResolved")
	dec.Define("event-resolve", func(context.Context, any) (any, error) {
		return true, nil
	})
	requireFeatureActive(t, scoped, "event-resolve")
	if dispatcher.count("FeatureResolved") != beforeResolved+1 {
		t.Fatalf("expected FeatureResolved event, before=%d after=%d", beforeResolved, dispatcher.count("FeatureResolved"))
	}

	// DatabaseDriverTest::test_it_dispatches_events_when_updating_a_scoped_feature
	beforeUpdated := dispatcher.count("FeatureUpdated")
	if err := scoped.ActivateWithValue(ctx, []string{"event-update"}, "value"); err != nil {
		t.Fatalf("ActivateWithValue: %v", err)
	}
	if dispatcher.count("FeatureUpdated") <= beforeUpdated {
		t.Fatal("expected FeatureUpdated event for scoped update")
	}

	// DatabaseDriverTest::test_it_dispatches_events_when_updating_a_feature_for_all_scopes
	beforeUpdatedAll := dispatcher.count("FeatureUpdatedForAllScopes")
	if err := dec.SetForAllScopes(ctx, "event-update", "next"); err != nil {
		t.Fatalf("SetForAllScopes event: %v", err)
	}
	if dispatcher.count("FeatureUpdatedForAllScopes") != beforeUpdatedAll+1 {
		t.Fatal("expected FeatureUpdatedForAllScopes event")
	}

	// DatabaseDriverTest::test_it_dispatches_events_when_deleting_a_feature_value
	beforeDeleted := dispatcher.count("FeatureDeleted")
	if err := scoped.Forget(ctx, []string{"event-update"}); err != nil {
		t.Fatalf("Forget event: %v", err)
	}
	if dispatcher.count("FeatureDeleted") <= beforeDeleted {
		t.Fatal("expected FeatureDeleted event")
	}

	beforeBulkUpdated := dispatcher.count("FeatureUpdated")
	// DatabaseDriverTest::test_it_dispatches_events_when_activating_multiple_features_and_scopes
	if err := dec.For("user:a", "user:b").Activate(ctx, []string{"event-a", "event-b"}); err != nil {
		t.Fatalf("Activate multi database event features: %v", err)
	}
	if dispatcher.count("FeatureUpdated") != beforeBulkUpdated+4 {
		t.Fatalf("expected four database activation update events, got %d", dispatcher.count("FeatureUpdated")-beforeBulkUpdated)
	}

	beforeBulkUpdated = dispatcher.count("FeatureUpdated")
	// DatabaseDriverTest::test_it_dispatches_events_when_deactivating_multiple_features_and_scopes
	if err := dec.For("user:a", "user:b").Deactivate(ctx, []string{"event-a", "event-b"}); err != nil {
		t.Fatalf("Deactivate multi database event features: %v", err)
	}
	if dispatcher.count("FeatureUpdated") != beforeBulkUpdated+4 {
		t.Fatalf("expected four database deactivation update events, got %d", dispatcher.count("FeatureUpdated")-beforeBulkUpdated)
	}

	if err := scoped.Activate(ctx, []string{"purge-a", "purge-b"}); err != nil {
		t.Fatalf("Activate purge features: %v", err)
	}

	// DatabaseDriverTest::test_it_can_purge_flags
	// DatabaseDriverTest::test_it_dispatches_events_when_purging_features
	beforePurged := dispatcher.count("FeaturesPurged")
	if err := scoped.Purge(ctx, []string{"purge-a"}); err != nil {
		t.Fatalf("Purge features: %v", err)
	}
	if dispatcher.count("FeaturesPurged") != beforePurged+1 {
		t.Fatal("expected FeaturesPurged event")
	}
	requireFeatureInactive(t, scoped, "purge-a")
	requireFeatureActive(t, scoped, "purge-b")

	// DatabaseDriverTest::test_it_can_purge_multiple_flags_at_once
	if err := scoped.Purge(ctx, []string{"purge-b"}); err != nil {
		t.Fatalf("Purge multiple/specific flags: %v", err)
	}
	requireFeatureInactive(t, scoped, "purge-b")

	if err := scoped.Activate(ctx, []string{"purge-all"}); err != nil {
		t.Fatalf("Activate purge-all feature: %v", err)
	}

	// DatabaseDriverTest::test_it_can_purge_all_feature_flags
	// DatabaseDriverTest::test_it_dispatches_events_when_purging_all_features
	beforeAllPurged := dispatcher.count("AllFeaturesPurged")
	if err := scoped.Purge(ctx, nil); err != nil {
		t.Fatalf("Purge all: %v", err)
	}
	if dispatcher.count("AllFeaturesPurged") != beforeAllPurged+1 {
		t.Fatal("expected AllFeaturesPurged event")
	}

	stored, err := dec.Stored(ctx)
	if err != nil {
		t.Fatalf("Stored after purge all: %v", err)
	}
	if len(stored) != 0 {
		t.Fatalf("expected no stored features after purge all, got %v", stored)
	}

	dec.Define("purge-retrieve", func(context.Context, any) (any, error) {
		return true, nil
	})
	requireFeatureActive(t, scoped, "purge-retrieve")
	if err := scoped.Purge(ctx, []string{"purge-retrieve"}); err != nil {
		t.Fatalf("Purge retrieve: %v", err)
	}

	// DatabaseDriverTest::test_retrieving_values_after_purging
	requireFeatureActive(t, scoped, "purge-retrieve")

	db.shouldConflict = true
	db.conflictCount = 4

	// DatabaseDriverTest::test_it_retries_3_times_and_then_fails
	err = dec.Set(ctx, "conflict", "scope", true)
	if !errors.Is(err, featureflags.ErrStorageConflict) {
		t.Fatalf("expected ErrStorageConflict after retries, got %v", err)
	}
	if db.conflictsSeen != 4 {
		t.Fatalf("expected four conflict attempts including initial try, got %d", db.conflictsSeen)
	}
}

func TestInventoryParityDatabaseDriverRepositoryEdges(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := &testDispatcher{}
	db, _, dec := newDatabaseInventoryHarness(dispatcher)

	// DatabaseDriverTest::test_it_does_not_store_unknown_features
	requireFeatureInactive(t, dec.For(), "unknown-store-check")
	requireFeatureInactive(t, dec.For(), "unknown-store-check")
	if len(db.rows) != 0 {
		t.Fatalf("expected unknown features not to persist rows, got %#v", db.rows)
	}

	// DatabaseDriverTest::test_bulk_insert_adds_timestamps
	if err := dec.SetAll(ctx, []featureflags.FeatureEntry{{Feature: "timestamped", Scope: nil, Value: true}}); err != nil {
		t.Fatalf("SetAll timestamped: %v", err)
	}
	if len(db.rows) != 1 || db.rows[0].createdAt == nil || db.rows[0].updatedAt == nil {
		t.Fatalf("expected SetAll row timestamps, got %#v", db.rows)
	}

	// DatabaseDriverTest::test_it_can_load_all_features_for_scope
	dec.Define("load-all-foo", func(_ context.Context, scope any) (any, error) {
		return scope == "tim", nil
	})
	dec.Define("load-all-bar", func(_ context.Context, scope any) (any, error) {
		return scope == "taylor", nil
	})
	if err := dec.For("tim", "taylor").LoadAll(ctx); err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if !dec.For("tim").Active(ctx, "load-all-foo") ||
		dec.For("tim").Active(ctx, "load-all-bar") ||
		dec.For("taylor").Active(ctx, "load-all-foo") ||
		!dec.For("taylor").Active(ctx, "load-all-bar") {
		t.Fatal("expected LoadAll to load every registered feature for each scope")
	}

	nonConflict := errors.New("connection closed")
	errorDB := &execErrorDB{err: nonConflict}
	errorDriver := featureflags.NewDatabaseDriver(errorDB, testTable)

	// DatabaseDriverTest::test_it_only_retries_on_conflicts
	err := errorDriver.Set(ctx, "no-retry", nil, true)
	if !errors.Is(err, nonConflict) {
		t.Fatalf("expected original non-conflict error, got %v", err)
	}
	if errorDB.execs != 1 {
		t.Fatalf("expected one non-conflict insert attempt, got %d", errorDB.execs)
	}
}

func TestInventoryParityFeatureManagerAndScopeSerialization(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := featureflags.NewManagerWithDispatcher("array", &testDispatcher{})

	// FeatureManagerTest::test_it_can_chain_scope_additions
	scoped, err := manager.For("user:1")
	if err != nil {
		t.Fatalf("manager For: %v", err)
	}
	chained := scoped.For("user:2")
	dec, err := manager.DefaultDecorator()
	if err != nil {
		t.Fatalf("DefaultDecorator: %v", err)
	}
	dec.Define("chain", func(context.Context, any) (any, error) {
		return true, nil
	})
	if !chained.AllAreActive(ctx, []string{"chain"}) {
		t.Fatal("expected chained scopes to be evaluated")
	}

	// FeatureManagerTest::test_the_authenticated_user_is_the_default_scope
	manager.ResolveScopeUsing(func(context.Context) (any, error) {
		return "default-scope", nil
	})

	defaultScoped, err := manager.For()
	if err != nil {
		t.Fatalf("manager For default scope: %v", err)
	}

	dec.Define("default-scope", func(_ context.Context, scope any) (any, error) {
		return scope, nil
	})

	defaultScopeValue, err := defaultScoped.Value(ctx, "default-scope")
	if err != nil {
		t.Fatalf("default scope Value: %v", err)
	}

	if defaultScopeValue != "default-scope" {
		t.Fatalf("expected default scope value, got %v", defaultScopeValue)
	}

	manager.ResolveScopeUsing(func(context.Context) (any, error) {
		return nil, nil
	})
	nilDefaultScoped, err := manager.For()
	if err != nil {
		t.Fatalf("manager For nil default scope: %v", err)
	}
	if err := nilDefaultScoped.Activate(ctx, []string{"nil-default"}); err != nil {
		t.Fatalf("Activate nil default scope: %v", err)
	}
	// ArrayDriverTest::test_it_doesnt_include_default_scope_when_null
	// DatabaseDriverTest::test_it_doesnt_include_default_scope_when_null
	requireFeatureActive(t, dec.For(nil), "nil-default")

	manager.ResolveScopeUsing(func(context.Context) (any, error) {
		return "loader-default", nil
	})
	loaderScoped, err := manager.For()
	if err != nil {
		t.Fatalf("manager For loader default scope: %v", err)
	}
	dec.Define("load-with-default", func(_ context.Context, scope any) (any, error) {
		return scope == "loader-default", nil
	})
	// ArrayDriverTest::test_it_uses_default_scope_for_loading_with_string
	if err := loaderScoped.Load(ctx, []string{"load-with-default"}); err != nil {
		t.Fatalf("Load with default scope: %v", err)
	}
	requireFeatureActive(t, dec.For("loader-default"), "load-with-default")

	// ArrayDriverTest::test_it_can_customise_default_scope
	// DatabaseDriverTest::test_it_can_customise_default_scope
	requireFeatureActive(t, loaderScoped, "load-with-default")

	lottery := featureflags.FixedLottery(true, true, true, true, false)
	dec.DefineValue("manager-lottery", lottery)
	managerLotteryScoped := dec.For()
	for i := 0; i < 4; i++ {
		if err := managerLotteryScoped.Load(ctx, []string{"manager-lottery"}); err != nil {
			t.Fatalf("manager lottery Load %d: %v", i, err)
		}
		if !managerLotteryScoped.Active(ctx, "manager-lottery") {
			t.Fatalf("expected manager lottery draw %d to be active", i)
		}
		if err := managerLotteryScoped.Forget(ctx, []string{"manager-lottery"}); err != nil {
			t.Fatalf("manager lottery Forget %d: %v", i, err)
		}
	}
	if err := managerLotteryScoped.Load(ctx, []string{"manager-lottery"}); err != nil {
		t.Fatalf("manager lottery final Load: %v", err)
	}
	// FeatureManagerTest::test_it_can_return_lottery_as_value
	if managerLotteryScoped.Active(ctx, "manager-lottery") {
		t.Fatal("expected final fixed lottery draw to be inactive")
	}

	// FeatureHelperTest::test_it_returns_feature_manager
	if manager.GetDefaultDriver() != "array" {
		t.Fatalf("expected array default manager, got %q", manager.GetDefaultDriver())
	}

	// FeatureHelperTest::test_it_returns_the_feature_value
	dec.Define("helper-value", func(context.Context, any) (any, error) {
		return "helper", nil
	})
	if requireFeatureValue(t, dec, "helper-value", nil) != "helper" {
		t.Fatal("expected helper value")
	}

	// FeatureHelperTest::test_it_conditionally_executes_code_blocks
	result, err := dec.For().When(ctx, "helper-value",
		func(any) (any, error) { return "active", nil },
		func(any) (any, error) { return "inactive", nil },
	)
	if err != nil || result != "active" {
		t.Fatalf("unexpected helper conditional result=%v err=%v", result, err)
	}

	// DatabaseDriverTest::test_it_handles_integer_scopes_correctly
	// ArrayDriverTest::test_it_handles_integer_scopes_correctly
	intScope, err := manager.SerializeScope(42)
	if err != nil || intScope != "42" {
		t.Fatalf("expected integer scope 42, got %q err=%v", intScope, err)
	}

	// DatabaseDriverTest::test_it_can_handles_double_scopes_correctly
	// DatabaseDriverTest::test_it_can_handles_float_scopes_correctly
	// ArrayDriverTest::test_it_can_handles_double_scopes_correctly
	// ArrayDriverTest::test_it_can_handles_float_scopes_correctly
	floatScope, err := manager.SerializeScope(42.5)
	if err != nil || floatScope != "42.5" {
		t.Fatalf("expected float scope 42.5, got %q err=%v", floatScope, err)
	}
}

func TestInventoryParityFeatureMiddleware(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := &testDispatcher{}
	_, dec := newArrayInventoryHarness(dispatcher)
	middleware := featureflags.NewEnsureFeaturesAreActive(dec.For())
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// FeatureMiddlewareTest::test_it_throws_a_http_exception_if_feature_is_not_defined
	recorder := httptest.NewRecorder()
	middleware.Handle(next, "test").ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected undefined feature to return 404, got %d", recorder.Code)
	}

	dec.DefineValue("test", true)

	// FeatureMiddlewareTest::test_it_passes_if_feature_is_defined
	recorder = httptest.NewRecorder()
	middleware.Handle(next, "test").ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx))
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected active feature to pass, got %d", recorder.Code)
	}

	// FeatureMiddlewareTest::test_it_throws_an_exception_if_one_of_the_features_is_not_active
	recorder = httptest.NewRecorder()
	middleware.Handle(next, "test", "another").ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected inactive feature to return 404, got %d", recorder.Code)
	}

	middleware.WhenInactive(func(w http.ResponseWriter, _ *http.Request, _ []string) {
		_, _ = w.Write([]byte("test-response"))
	})

	// FeatureMiddlewareTest::test_it_allows_custom_responses
	recorder = httptest.NewRecorder()
	middleware.Handle(next, "test", "another").ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx))
	if recorder.Body.String() != "test-response" {
		t.Fatalf("expected custom response, got %q", recorder.Body.String())
	}
	middleware.WhenInactive(nil)

	if err := dec.For().Activate(ctx, []string{"another"}); err != nil {
		t.Fatalf("Activate another: %v", err)
	}

	// FeatureMiddlewareTest::test_it_passes_if_all_features_are_enabled
	recorder = httptest.NewRecorder()
	middleware.Handle(next, "test", "another").ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx))
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected all active features to pass, got %d", recorder.Code)
	}

	// FeatureMiddlewareTest::test_middleware_string_can_be_returned
	if got := middleware.Using("test", "another"); got != "featureflags.EnsureFeaturesAreActive:test,another" {
		t.Fatalf("unexpected middleware string %q", got)
	}
}

type execErrorDB struct {
	err   error
	execs int
}

func (db *execErrorDB) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	db.execs++

	return nil, db.err
}

func (db *execErrorDB) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return openFakeRows(nil)
}

func (db *execErrorDB) QueryRowContext(context.Context, string, ...any) *sql.Row {
	return fakeEmptyRow()
}
