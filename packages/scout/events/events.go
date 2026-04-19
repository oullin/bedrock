package events

import contract "github.com/bedrock/packages/contracts/scout"

// ModelsImported is dispatched after models have been indexed in the
// search engine. It carries the collection of models that were imported.
type ModelsImported struct {
	Models []contract.Searchable
}

// ModelsFlushed is dispatched after all records for a model type have
// been removed from the search index.
type ModelsFlushed struct {
	Models []contract.Searchable
}

// SearchableModelUpdated is dispatched when a searchable model is
// synced to the search index after a save operation.
type SearchableModelUpdated struct {
	Model contract.Searchable
}

// SearchableModelDeleted is dispatched when a searchable model is
// removed from the search index after a delete operation.
type SearchableModelDeleted struct {
	Model contract.Searchable
}
