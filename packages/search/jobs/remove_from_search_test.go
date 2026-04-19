package jobs_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/bedrock/packages/search/jobs"
)

func TestRemoveFromSearchHandle(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	models := []contract.Searchable{
		&testModel{id: 1, table: "posts"},
		&testModel{id: 2, table: "posts"},
	}

	job := jobs.NewRemoveFromSearch(models, engine)
	err := job.Handle(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call, got %d", engine.deleteCalls)
	}

	if len(engine.lastModels) != 2 {
		t.Fatalf("expected 2 models, got %d", len(engine.lastModels))
	}
}

func TestRemoveFromSearchHandleEmpty(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	job := jobs.NewRemoveFromSearch(nil, engine)
	err := job.Handle(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.deleteCalls != 0 {
		t.Fatal("should not call delete for empty models")
	}
}

func TestRemoveFromSearchGetModels(t *testing.T) {
	t.Parallel()
	models := []contract.Searchable{
		&testModel{id: 1, table: "posts"},
	}
	job := jobs.NewRemoveFromSearch(models, &fakeEngine{})

	if len(job.GetModels()) != 1 {
		t.Fatalf("expected 1 model, got %d", len(job.GetModels()))
	}
}
