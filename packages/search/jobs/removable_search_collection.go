package jobs

import contract "github.com/bedrock/packages/contracts/search"

// RemovableSearchCollection carries enough model identity to remove records
// from a search index without depending on database rehydration.
type RemovableSearchCollection struct {
	models []contract.Searchable
}

// NewRemovableSearchCollection creates a removable collection for the models.
func NewRemovableSearchCollection(models []contract.Searchable) RemovableSearchCollection {
	return RemovableSearchCollection{models: models}
}

// Models returns the underlying searchable models.
func (c RemovableSearchCollection) Models() []contract.Searchable {
	return c.models
}

// GetQueueableIDs returns search keys for queue payload identity.
func (c RemovableSearchCollection) GetQueueableIDs() []any {
	return c.SearchKeys()
}

// SearchKeys returns the document keys that should be removed from the index.
func (c RemovableSearchCollection) SearchKeys() []any {
	keys := make([]any, 0, len(c.models))

	for _, model := range c.models {
		keys = append(keys, model.GetSearchKey())
	}

	return keys
}
