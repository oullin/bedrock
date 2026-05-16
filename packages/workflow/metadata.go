package workflow

type MetadataStore interface {
	WorkflowMetadata(key string) (any, bool)
	PlaceMetadata(place string, key string) (any, bool)
	TransitionMetadata(transition string, key string) (any, bool)
}

type DefinitionMetadataStore struct {
	definition *Definition
}

func NewMetadataStore(definition *Definition) *DefinitionMetadataStore {
	return &DefinitionMetadataStore{definition: definition}
}

func (s *DefinitionMetadataStore) WorkflowMetadata(key string) (any, bool) {
	return s.definition.MetadataValue(key)
}

func (s *DefinitionMetadataStore) PlaceMetadata(place string, key string) (any, bool) {
	return s.definition.PlaceMetadataValue(place, key)
}

func (s *DefinitionMetadataStore) TransitionMetadata(transition string, key string) (any, bool) {
	return s.definition.TransitionMetadataValue(transition, key)
}
