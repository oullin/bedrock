package scout

import (
	"context"

	cevents "github.com/bedrock/packages/contracts/events"
	contract "github.com/bedrock/packages/contracts/scout"
	scoutevents "github.com/bedrock/packages/scout/events"
)

// SearchableScope provides bulk operations for making models searchable
// or removing them from the search index. It mirrors Scout // SearchableScope which adds searchable/unsearchable macros to the
// Eloquent query builder.
type SearchableScope struct {
	engine     contract.Engine
	config     Config
	dispatcher cevents.Dispatcher
}

// NewSearchableScope creates a new SearchableScope.
func NewSearchableScope(engine contract.Engine, config Config, dispatcher ...cevents.Dispatcher) *SearchableScope {
	var d cevents.Dispatcher

	if len(dispatcher) > 0 {
		d = dispatcher[0]
	}

	return &SearchableScope{
		engine:     engine,
		config:     config,
		dispatcher: d,
	}
}

// Searchable indexes the given models in chunks.
func (s *SearchableScope) Searchable(ctx context.Context, models []contract.Searchable) error {
	if len(models) == 0 {
		return nil
	}

	chunkSize := s.config.GetChunkSearchable()

	for i := 0; i < len(models); i += chunkSize {
		end := i + chunkSize

		if end > len(models) {
			end = len(models)
		}

		chunk := models[i:end]

		// Filter to only searchable models.
		var searchable []contract.Searchable

		for _, m := range chunk {
			if m.ShouldBeSearchable() {
				searchable = append(searchable, m)
			}
		}

		if len(searchable) == 0 {
			continue
		}

		if err := s.engine.Update(ctx, searchable); err != nil {
			return err
		}

		// Dispatch ModelsImported event.
		if s.dispatcher != nil {
			_, _ = s.dispatcher.Dispatch(ctx, scoutevents.ModelsImported{
				Models: searchable,
			})
		}
	}

	return nil
}

// Unsearchable removes the given models from the search index in chunks.
func (s *SearchableScope) Unsearchable(ctx context.Context, models []contract.Searchable) error {
	if len(models) == 0 {
		return nil
	}

	chunkSize := s.config.GetChunkUnsearchable()

	for i := 0; i < len(models); i += chunkSize {
		end := i + chunkSize

		if end > len(models) {
			end = len(models)
		}

		chunk := models[i:end]

		if err := s.engine.Delete(ctx, chunk); err != nil {
			return err
		}

		// Dispatch ModelsFlushed event.
		if s.dispatcher != nil {
			_, _ = s.dispatcher.Dispatch(ctx, scoutevents.ModelsFlushed{
				Models: chunk,
			})
		}
	}

	return nil
}
