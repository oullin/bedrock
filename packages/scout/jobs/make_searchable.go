package jobs

import (
	"context"

	contract "github.com/bedrock/packages/contracts/scout"
)

// MakeSearchable is a queueable job that indexes models in the search engine.
type MakeSearchable struct {
	Models []contract.Searchable
	engine contract.Engine
}

// NewMakeSearchable creates a new MakeSearchable job.
func NewMakeSearchable(models []contract.Searchable, engine contract.Engine) *MakeSearchable {
	return &MakeSearchable{
		Models: models,
		engine: engine,
	}
}

// Handle executes the job, indexing all models in the search engine.
func (j *MakeSearchable) Handle(ctx context.Context) error {
	if len(j.Models) == 0 {
		return nil
	}

	// Filter to searchable models only.
	var searchable []contract.Searchable

	for _, m := range j.Models {
		if m.ShouldBeSearchable() {
			searchable = append(searchable, m)
		}
	}

	if len(searchable) == 0 {
		return nil
	}

	return j.engine.Update(ctx, searchable)
}

// GetModels returns the models to be indexed.
func (j *MakeSearchable) GetModels() []contract.Searchable {
	return j.Models
}
