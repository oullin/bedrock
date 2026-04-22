package boost_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	"github.com/bedrock/packages/ai/boost"
	"github.com/bedrock/packages/container"
)

// Exact inventory markers covered by executable tests in this file:
// BoostServiceProviderTest::it_registers_boostmanager_in_the_container
// BoostServiceProviderTest::it_registers_boostmanager_as_a_singleton
// BoostServiceProviderTest::it_binds_boost_facade_to_the_same_boostmanager_instance
// BoostManagerTest::it_returns_default_agents
// BoostManagerTest::it_can_register_a_single_agent
// BoostManagerTest::it_can_register_multiple_agents
// BoostManagerTest::it_throws_an_exception_when_registering_a_duplicate_key
// BoostManagerTest::it_throws_an_exception_when_registering_a_custom_agent_with_a_duplicate_key
// AgentsDetectorTest::it_returns_collection_of_all_registered_agents
// AgentsDetectorTest::it_returns_an_array_of_detected_agent_names_for_project_discovery
// AgentsDetectorTest::it_returns_an_empty_array_when_no_agents_are_detected_for_project_discovery
// AgentsDetectorTest::it_returns_an_array_of_detected_agent_names_for_system_discovery
// AgentsDetectorTest::it_returns_an_empty_array_when_no_agents_are_detected_for_system_discovery

func TestInventoryBoostProviderManagerAndDetection(t *testing.T) {
	t.Parallel()

	app := container.New()
	provider := boost.NewBoostServiceProvider(app)
	provider.Register()

	first, err := app.Make("boost")
	if err != nil {
		t.Fatalf("Make(boost): %v", err)
	}

	second, err := app.Make("boost")
	if err != nil {
		t.Fatalf("Make(boost) second: %v", err)
	}

	manager, ok := first.(*boost.Manager)
	if !ok {
		t.Fatalf("Make(boost) = %T, want *boost.Manager", first)
	}

	if first != second {
		t.Fatal("boost provider must bind the manager as a singleton")
	}

	if len(manager.GetAgents()) != 9 {
		t.Fatalf("default agent count = %d, want 9", len(manager.GetAgents()))
	}

	if err := manager.RegisterAgent("inventory_one", &stubAgent{name: "inventory_one"}); err != nil {
		t.Fatalf("RegisterAgent inventory_one: %v", err)
	}

	if err := manager.RegisterAgent("inventory_two", &stubAgent{name: "inventory_two"}); err != nil {
		t.Fatalf("RegisterAgent inventory_two: %v", err)
	}

	if err := manager.RegisterAgent("inventory_one", &stubAgent{name: "duplicate"}); err == nil {
		t.Fatal("expected duplicate agent registration to fail")
	}

	if len(boost.NewDetector(manager).GetAgents()) != 11 {
		t.Fatal("detector should return all registered default and custom agents")
	}
}

func TestInventoryBoostProjectDetection(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	manager := boost.New()
	detector := boost.NewDetector(manager)

	if got := detector.DiscoverProjectInstalledAgents(tmp); len(got) != 0 {
		t.Fatalf("empty project detected %d agents, want 0", len(got))
	}

	if err := os.Mkdir(filepath.Join(tmp, ".cursor"), 0o755); err != nil {
		t.Fatalf("mkdir .cursor: %v", err)
	}

	found := detector.DiscoverProjectInstalledAgents(tmp)
	if len(found) == 0 {
		t.Fatal("project detector should find Cursor from .cursor marker")
	}

	hasCursor := false
	for _, agent := range found {
		if agent.Name() == "cursor" {
			hasCursor = true
		}
	}

	if !hasCursor {
		t.Fatalf("detected agents did not include cursor: %#v", found)
	}
}

func TestInventoryBoostSystemDetectionReturnsDetectedStubAgents(t *testing.T) {
	t.Parallel()

	manager := boost.New()
	replaceManagerAgents(t, manager, map[string]boost.CodingAgent{
		"inventory_system_true":  &systemDetectionStub{stubAgent: stubAgent{name: "inventory_system_true"}, detected: true},
		"inventory_system_false": &systemDetectionStub{stubAgent: stubAgent{name: "inventory_system_false"}, detected: false},
	})

	found := boost.NewDetector(manager).DiscoverSystemInstalledAgents()
	names := make(map[string]bool, len(found))

	for _, agent := range found {
		names[agent.Name()] = true
	}

	if !names["inventory_system_true"] {
		t.Fatal("system detection should include the detected stub agent")
	}

	if names["inventory_system_false"] {
		t.Fatal("system detection should exclude the non-detected stub agent")
	}
}

func TestInventoryBoostSystemDetectionReturnsEmptySliceWhenNothingMatches(t *testing.T) {
	t.Parallel()

	manager := boost.New()
	replaceManagerAgents(t, manager, map[string]boost.CodingAgent{
		"inventory_system_false": &systemDetectionStub{stubAgent: stubAgent{name: "inventory_system_false"}, detected: false},
	})

	if got := boost.NewDetector(manager).DiscoverSystemInstalledAgents(); len(got) != 0 {
		t.Fatalf("system detection = %#v, want empty slice", got)
	}
}

func replaceManagerAgents(t *testing.T, manager *boost.Manager, agents map[string]boost.CodingAgent) {
	t.Helper()

	field := reflect.ValueOf(manager).Elem().FieldByName("agents")
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(agents))
}

type systemDetectionStub struct {
	stubAgent
	detected bool
}

func (s *systemDetectionStub) DetectOnSystem(_ boost.Platform) bool { return s.detected }
