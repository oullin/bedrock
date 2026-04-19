package scout_test

import (
	"context"
	"testing"

	cevents "github.com/bedrock/packages/contracts/events"
	contract "github.com/bedrock/packages/contracts/scout"
	"github.com/bedrock/packages/scout"
)

// fakeDispatcher records dispatched events.
type fakeDispatcher struct {
	dispatched []any
}

func (d *fakeDispatcher) Listen(_ any, _ ...cevents.Listener) {}
func (d *fakeDispatcher) HasListeners(_ any) bool             { return false }
func (d *fakeDispatcher) HasWildcardListeners(_ any) bool     { return false }
func (d *fakeDispatcher) Subscribe(_ cevents.Subscriber)      {}
func (d *fakeDispatcher) Until(_ context.Context, _ any) (any, error) {
	return nil, nil
}
func (d *fakeDispatcher) Dispatch(_ context.Context, event any) ([]any, error) {
	d.dispatched = append(d.dispatched, event)

	return nil, nil
}
func (d *fakeDispatcher) Push(_ context.Context, _ any)           {}
func (d *fakeDispatcher) Flush(_ context.Context, _ string) error { return nil }
func (d *fakeDispatcher) Forget(_ any)                            {}
func (d *fakeDispatcher) ForgetPushed()                           {}
func (d *fakeDispatcher) GetListeners(_ any) []cevents.Listener   { return nil }

func TestSearchableScopeSearchable(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	config := scout.Config{Chunk: scout.ChunkConfig{Searchable: 5}}
	dispatcher := &fakeDispatcher{}
	scope := scout.NewSearchableScope(engine, config, dispatcher)

	models := make([]contract.Searchable, 12)

	for i := range models {
		models[i] = newTestModel(i+1, "posts")
	}

	err := scope.Searchable(context.Background(), models)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 12 models / chunk 5 = 3 chunks
	if engine.updateCalls != 3 {
		t.Fatalf("expected 3 update calls, got %d", engine.updateCalls)
	}

	// 3 ModelsImported events dispatched.
	if len(dispatcher.dispatched) != 3 {
		t.Fatalf("expected 3 dispatched events, got %d", len(dispatcher.dispatched))
	}
}

func TestSearchableScopeSearchableEmpty(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	config := scout.DefaultConfig()
	scope := scout.NewSearchableScope(engine, config)

	err := scope.Searchable(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.updateCalls != 0 {
		t.Fatal("should not call update for empty models")
	}
}

func TestSearchableScopeUnsearchable(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	config := scout.Config{Chunk: scout.ChunkConfig{Unsearchable: 3}}
	dispatcher := &fakeDispatcher{}
	scope := scout.NewSearchableScope(engine, config, dispatcher)

	models := make([]contract.Searchable, 7)

	for i := range models {
		models[i] = newTestModel(i+1, "posts")
	}

	err := scope.Unsearchable(context.Background(), models)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 7 models / chunk 3 = 3 chunks
	if engine.deleteCalls != 3 {
		t.Fatalf("expected 3 delete calls, got %d", engine.deleteCalls)
	}

	// 3 ModelsFlushed events dispatched.
	if len(dispatcher.dispatched) != 3 {
		t.Fatalf("expected 3 dispatched events, got %d", len(dispatcher.dispatched))
	}
}

func TestSearchableScopeUnsearchableEmpty(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	config := scout.DefaultConfig()
	scope := scout.NewSearchableScope(engine, config)

	err := scope.Unsearchable(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.deleteCalls != 0 {
		t.Fatal("should not call delete for empty models")
	}
}

func TestSearchableScopeWithoutDispatcher(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	config := scout.Config{Chunk: scout.ChunkConfig{Searchable: 100}}
	scope := scout.NewSearchableScope(engine, config) // no dispatcher

	models := []contract.Searchable{newTestModel(1, "posts")}
	err := scope.Searchable(context.Background(), models)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call, got %d", engine.updateCalls)
	}
}
