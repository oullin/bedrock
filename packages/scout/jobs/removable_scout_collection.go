package jobs

import contract "github.com/bedrock/packages/contracts/scout"

// RemovableScoutCollection carries enough model identity to remove records
// from a search index without depending on database rehydration.
type RemovableScoutCollection struct {
	models []contract.Searchable
}

// NewRemovableScoutCollection creates a removable collection for the models.
func NewRemovableScoutCollection(models []contract.Searchable) RemovableScoutCollection {
	return RemovableScoutCollection{models: models}
}

// Models returns the underlying searchable models.
func (c RemovableScoutCollection) Models() []contract.Searchable {
	return c.models
}

// GetQueueableIDs returns scout keys for queue payload identity.
func (c RemovableScoutCollection) GetQueueableIDs() []any {
	return c.ScoutKeys()
}

// ScoutKeys returns the document keys that should be removed from the index.
func (c RemovableScoutCollection) ScoutKeys() []any {
	keys := make([]any, 0, len(c.models))

	for _, model := range c.models {
		keys = append(keys, model.GetScoutKey())
	}

	return keys
}
