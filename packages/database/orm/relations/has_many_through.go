package relations

import (
	"context"

	"github.com/bedrock/packages/database/orm"
)

// HasManyThrough defines a one-to-many relationship through an intermediate table.
type HasManyThrough struct {
	*BaseRelation
	throughParent *orm.Model
	farParent     *orm.Model
	firstKey      string
	secondKey     string
	localKey      string
	secondLocalKey string
}

// NewHasManyThrough creates a new HasManyThrough relationship.
func NewHasManyThrough(parent, throughParent, farParent *orm.Model, firstKey, secondKey, localKey, secondLocalKey string) *HasManyThrough {
	return &HasManyThrough{
		BaseRelation:   NewBaseRelation(nil, parent, farParent),
		throughParent:  throughParent,
		farParent:      farParent,
		firstKey:       firstKey,
		secondKey:      secondKey,
		localKey:       localKey,
		secondLocalKey: secondLocalKey,
	}
}

func (r *HasManyThrough) AddConstraints() {}

func (r *HasManyThrough) AddEagerConstraints(models []*orm.Model) {
	keys := make([]any, 0, len(models))
	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.localKey))
	}
	if r.query != nil {
		r.query.WhereIn(r.throughParent.GetTable()+"."+r.firstKey, keys)
	}
}

func (r *HasManyThrough) InitRelation(models []*orm.Model, relation string) []*orm.Model {
	for _, model := range models {
		model.SetAttribute(relation, []*orm.Model{})
	}
	return models
}

func (r *HasManyThrough) Match(models []*orm.Model, results []*orm.Model, relation string) []*orm.Model {
	dictionary := make(map[any][]*orm.Model)
	for _, result := range results {
		key := result.GetAttribute("laravel_through_key")
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

func (r *HasManyThrough) GetResults() ([]*orm.Model, error) {
	if r.query == nil {
		return nil, nil
	}
	rows, err := r.query.Get(context.Background())
	if err != nil {
		return nil, err
	}
	models := make([]*orm.Model, 0, len(rows))
	for _, row := range rows {
		m := orm.NewModel()
		m.SetTable(r.farParent.GetTable())
		m.SetRawAttributes(row, true)
		m.SetExists(true)
		models = append(models, m)
	}
	return models, nil
}
