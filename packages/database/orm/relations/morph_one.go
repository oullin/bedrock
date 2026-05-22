package relations

import (
	"context"

	"github.com/bedrock/packages/database/orm"
)

// MorphOne defines a polymorphic one-to-one relationship.
type MorphOne struct {
	*BaseRelation
	morphType  string
	foreignKey string
	localKey   string
}

// NewMorphOne creates a new MorphOne relationship.
func NewMorphOne(parent, related *orm.Model, morphType, foreignKey, localKey string) *MorphOne {
	return &MorphOne{
		BaseRelation: NewBaseRelation(nil, parent, related),
		morphType:    morphType,
		foreignKey:   foreignKey,
		localKey:     localKey,
	}
}

func (r *MorphOne) GetMorphType() string      { return r.morphType }
func (r *MorphOne) GetForeignKeyName() string { return r.foreignKey }

func (r *MorphOne) AddConstraints() {
	if r.query != nil {
		r.query.Where(r.foreignKey, r.parent.GetAttribute(r.localKey))
		r.query.Where(r.morphType, r.parent.GetTable())
	}
}

func (r *MorphOne) AddEagerConstraints(models []*orm.Model) {
	keys := make([]any, 0, len(models))

	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.localKey))
	}

	if r.query != nil {
		r.query.WhereIn(r.foreignKey, keys)
		r.query.Where(r.morphType, r.parent.GetTable())
	}
}

func (r *MorphOne) InitRelation(models []*orm.Model, relation string) []*orm.Model {
	return models
}

func (r *MorphOne) Match(models []*orm.Model, results []*orm.Model, relation string) []*orm.Model {
	dictionary := make(map[any]*orm.Model)

	for _, result := range results {
		key := result.GetAttribute(r.foreignKey)
		dictionary[key] = result
	}

	for _, model := range models {
		key := model.GetAttribute(r.localKey)

		if match, ok := dictionary[key]; ok {
			model.SetAttribute(relation, match)
		}
	}

	return models
}

func (r *MorphOne) GetResults() ([]*orm.Model, error) {
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
