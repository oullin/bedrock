package relations

import (
	"context"

	"github.com/bedrock/packages/database/orm"
)

// MorphMany defines a polymorphic one-to-many relationship.
type MorphMany struct {
	*BaseRelation
	morphType  string
	foreignKey string
	localKey   string
}

// NewMorphMany creates a new MorphMany relationship.
func NewMorphMany(parent, related *orm.Model, morphType, foreignKey, localKey string) *MorphMany {
	return &MorphMany{
		BaseRelation: NewBaseRelation(nil, parent, related),
		morphType:    morphType,
		foreignKey:   foreignKey,
		localKey:     localKey,
	}
}

func (r *MorphMany) AddConstraints() {
	if r.query != nil {
		r.query.Where(r.foreignKey, r.parent.GetAttribute(r.localKey))
		r.query.Where(r.morphType, r.parent.GetTable())
	}
}

func (r *MorphMany) AddEagerConstraints(models []*orm.Model) {
	keys := make([]any, 0, len(models))

	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.localKey))
	}

	if r.query != nil {
		r.query.WhereIn(r.foreignKey, keys)
		r.query.Where(r.morphType, r.parent.GetTable())
	}
}

func (r *MorphMany) InitRelation(models []*orm.Model, relation string) []*orm.Model {
	for _, model := range models {
		model.SetAttribute(relation, []*orm.Model{})
	}

	return models
}

func (r *MorphMany) Match(models []*orm.Model, results []*orm.Model, relation string) []*orm.Model {
	dictionary := make(map[any][]*orm.Model)

	for _, result := range results {
		key := result.GetAttribute(r.foreignKey)
		dictionary[key] = append(dictionary[key], result)
	}

	for _, model := range models {
		key := model.GetAttribute(r.localKey)

		if matches, ok := dictionary[key]; ok {
			model.SetAttribute(relation, matches)
		}
	}

	return models
}

func (r *MorphMany) GetResults() ([]*orm.Model, error) {
	if r.query == nil {
		return nil, nil
	}

	rows, err := r.query.Get(context.Background())

	if err != nil {
		return nil, err
	}

	var models []*orm.Model

	for _, row := range rows {
		m := orm.NewModel()
		m.SetTable(r.related.GetTable())
		m.SetRawAttributes(row, true)
		m.SetExists(true)
		models = append(models, m)
	}

	return models, nil
}
