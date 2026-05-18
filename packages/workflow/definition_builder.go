package workflow

import "fmt"

// DefinitionBuilder constructs a Definition programmatically.
type DefinitionBuilder struct {
	definition *Definition
}

func NewDefinitionBuilder() *DefinitionBuilder {
	return &DefinitionBuilder{
		definition: &Definition{
			Metadata:           map[string]any{},
			PlaceMetadata:      map[string]map[string]any{},
			TransitionMetadata: map[string]map[string]any{},
		},
	}
}

func (b *DefinitionBuilder) AddPlace(place string) *DefinitionBuilder {
	if b.definition == nil {
		return b
	}

	if place == "" || b.definition.HasPlace(place) {
		return b
	}

	b.definition.Places = append(b.definition.Places, place)

	if b.definition.placesMap != nil {
		b.definition.placesMap[place] = struct{}{}
	}

	return b
}

func (b *DefinitionBuilder) SetInitialPlaces(places ...string) *DefinitionBuilder {
	if b.definition == nil {
		return b
	}

	b.definition.InitialMarking = NewMarking(places...)

	return b
}

func (b *DefinitionBuilder) AddTransition(name string, from []string, to []string) *DefinitionBuilder {
	if b.definition == nil {
		return b
	}

	b.definition.Transitions = append(b.definition.Transitions, NewTransition(name, from, to))

	return b
}

func (b *DefinitionBuilder) SetMetadata(key string, value any) *DefinitionBuilder {
	if b.definition == nil {
		return b
	}

	b.definition.Metadata[key] = value

	return b
}

func (b *DefinitionBuilder) SetPlaceMetadata(place string, key string, value any) *DefinitionBuilder {
	if b.definition == nil {
		return b
	}

	if b.definition.PlaceMetadata[place] == nil {
		b.definition.PlaceMetadata[place] = map[string]any{}
	}

	b.definition.PlaceMetadata[place][key] = value

	return b
}

func (b *DefinitionBuilder) SetTransitionMetadata(transition string, key string, value any) *DefinitionBuilder {
	if b.definition == nil {
		return b
	}

	if b.definition.TransitionMetadata[transition] == nil {
		b.definition.TransitionMetadata[transition] = map[string]any{}
	}

	b.definition.TransitionMetadata[transition][key] = value

	return b
}

func (b *DefinitionBuilder) Build() (*Definition, error) {
	if b.definition == nil {
		return nil, fmt.Errorf("definition builder is not initialized")
	}

	definition := b.definition.Clone()

	for i, transition := range definition.Transitions {
		if metadata, ok := definition.TransitionMetadata[transition.Name]; ok {
			if transition.Metadata == nil {
				transition.Metadata = map[string]any{}
			}

			for key, value := range metadata {
				transition.Metadata[key] = value
			}

			definition.Transitions[i] = transition
		}
	}

	if err := definition.Validate(); err != nil {
		return nil, err
	}

	return definition, nil
}
