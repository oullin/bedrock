package jobs

import (
	"context"

	contract "github.com/bedrock/packages/contracts/scout"
)

// RemoveFromSearch is a queueable job that removes models from the search engine.
// It mirrors Scout Jobs\RemoveFromSearch.
type RemoveFromSearch struct {
	Models []contract.Searchable
	engine contract.Engine
}

// NewRemoveFromSearch creates a new RemoveFromSearch job.
func NewRemoveFromSearch(models []contract.Searchable, engine contract.Engine) *RemoveFromSearch {
	return &RemoveFromSearch{
		Models: models,
		engine: engine,
	}
}

// Handle executes the job, removing all models from the search engine.
func (j *RemoveFromSearch) Handle(ctx context.Context) error {
	if len(j.Models) == 0 {
		return nil
	}

	return j.engine.Delete(ctx, j.Models)
}

// GetModels returns the models to be removed.
func (j *RemoveFromSearch) GetModels() []contract.Searchable {
	return j.Models
}
