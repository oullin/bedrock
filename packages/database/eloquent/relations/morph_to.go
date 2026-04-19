package relations

import (
	"context"

	"github.com/bedrock/packages/database/eloquent"
)

// MorphTo defines the inverse of a polymorphic relationship.
type MorphTo struct {
	*BaseRelation
	morphType  string
	foreignKey string
	ownerKey   string
}

// NewMorphTo creates a new MorphTo relationship.
func NewMorphTo(parent, related *eloquent.Model, morphType, foreignKey, ownerKey string) *MorphTo {
	return &MorphTo{
		BaseRelation: NewBaseRelation(nil, parent, related),
		morphType:    morphType,
		foreignKey:   foreignKey,
		ownerKey:     ownerKey,
	}
}

func (r *MorphTo) GetMorphType() string      { return r.morphType }
func (r *MorphTo) GetForeignKeyName() string { return r.foreignKey }

func (r *MorphTo) AddConstraints() {
	if r.query != nil {
		r.query.Where(r.ownerKey, r.parent.GetAttribute(r.foreignKey))
	}
}

func (r *MorphTo) AddEagerConstraints(models []*eloquent.Model) {
	keys := make([]any, 0, len(models))

	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.foreignKey))
	}

	if r.query != nil {
		r.query.WhereIn(r.ownerKey, keys)
	}
}

func (r *MorphTo) InitRelation(models []*eloquent.Model, relation string) []*eloquent.Model {
	return models
}

func (r *MorphTo) Match(models []*eloquent.Model, results []*eloquent.Model, relation string) []*eloquent.Model {
	dictionary := make(map[any]*eloquent.Model)

	for _, result := range results {
		key := result.GetAttribute(r.ownerKey)
		dictionary[key] = result
	}

	for _, model := range models {
		key := model.GetAttribute(r.foreignKey)

		if match, ok := dictionary[key]; ok {
			model.SetAttribute(relation, match)
		}
	}

	return models
}

func (r *MorphTo) GetResults() ([]*eloquent.Model, error) {
	if r.query == nil {
		return nil, nil
	}

	rows, err := r.query.Get(context.Background())

	if err != nil {
		return nil, err
	}

	var models []*eloquent.Model

	for _, row := range rows {
		m := eloquent.NewModel()
		m.SetTable(r.related.GetTable())
		m.SetRawAttributes(row, true)
		m.SetExists(true)
		models = append(models, m)
	}

	return models, nil
}

// Associate associates the morph-to with a model.
func (r *MorphTo) Associate(model *eloquent.Model) {
	r.parent.SetAttribute(r.foreignKey, model.GetAttribute(r.ownerKey))
	r.parent.SetAttribute(r.morphType, model.GetTable())
}

// Dissociate clears the polymorphic relationship.
func (r *MorphTo) Dissociate() {
	r.parent.SetAttribute(r.foreignKey, nil)
	r.parent.SetAttribute(r.morphType, nil)
}
