package translation_test

import (
	"testing"

	"github.com/bedrock/packages/translation"
)

func TestArrayLoaderLoadReturnsMsgsForRegisteredKey(t *testing.T) {
	t.Parallel()

	l := translation.NewArrayLoader()
	l.AddMessages("en", "messages", map[string]any{"welcome": "Hello"}, nil)

	got := l.Load("en", "messages", nil)

	if got["welcome"] != "Hello" {
		t.Errorf("got %v", got)
	}
}

func TestArrayLoaderLoadReturnsEmptyMapOnMiss(t *testing.T) {
	t.Parallel()

	l := translation.NewArrayLoader()
	got := l.Load("fr", "missing", nil)

	if got == nil || len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestArrayLoaderNilNamespaceDefaultsToGlobal(t *testing.T) {
	t.Parallel()

	l := translation.NewArrayLoader()
	l.AddMessages("en", "group", map[string]any{"k": "v"}, nil)

	got := l.Load("en", "group", nil)

	if got["k"] != "v" {
		t.Errorf("got %v", got)
	}
}

func TestArrayLoaderNamespaceIsolation(t *testing.T) {
	t.Parallel()

	l := translation.NewArrayLoader()
	ns := "vendor"
	l.AddMessages("en", "group", map[string]any{"k": "namespaced"}, &ns)

	// Global namespace should be empty.
	global := l.Load("en", "group", nil)

	if len(global) != 0 {
		t.Errorf("global namespace leaked: %v", global)
	}

	// Namespaced lookup should work.
	got := l.Load("en", "group", &ns)

	if got["k"] != "namespaced" {
		t.Errorf("got %v", got)
	}
}

func TestArrayLoaderAddMessagesIsChainable(t *testing.T) {
	t.Parallel()

	l := translation.NewArrayLoader().
		AddMessages("en", "a", map[string]any{"x": "1"}, nil).
		AddMessages("en", "b", map[string]any{"y": "2"}, nil)

	if l.Load("en", "a", nil)["x"] != "1" {
		t.Error("chain call a failed")
	}

	if l.Load("en", "b", nil)["y"] != "2" {
		t.Error("chain call b failed")
	}
}

func TestArrayLoaderNamespacesReturnsRegisteredEntries(t *testing.T) {
	t.Parallel()

	l := translation.NewArrayLoader()
	l.AddNamespace("pkg", "/hint")

	ns := l.Namespaces()

	if ns["pkg"] != "/hint" {
		t.Errorf("Namespaces() = %v", ns)
	}
}
