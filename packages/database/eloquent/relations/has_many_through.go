package relations

import (
	"context"

	"github.com/bedrock/packages/database/eloquent"
)

// HasManyThrough defines a one-to-many relationship through an intermediate table.
type HasManyThrough struct {
	*BaseRelation
	throughParent  *eloquent.Model
	farParent      *eloquent.Model
	firstKey       string
	secondKey      string
	localKey       string
	secondLocalKey string
}

// NewHasManyThrough creates a new HasManyThrough relationship.
func NewHasManyThrough(parent, throughParent, farParent *eloquent.Model, firstKey, secondKey, localKey, secondLocalKey string) *HasManyThrough {
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

func (r *HasManyThrough) AddEagerConstraints(models []*eloquent.Model) {
	keys := make([]any, 0, len(models))

	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.localKey))
	}

	if r.query != nil {
		r.query.WhereIn(r.throughParent.GetTable()+"."+r.firstKey, keys)
	}
}

func (r *HasManyThrough) InitRelation(models []*eloquent.Model, relation string) []*eloquent.Model {
	for _, model := range models {
		model.SetAttribute(relation, []*eloquent.Model{})
	}

	return models
}

func (r *HasManyThrough) Match(models []*eloquent.Model, results []*eloquent.Model, relation string) []*eloquent.Model {
	dictionary := make(map[any][]*eloquent.Model)

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

func (r *HasManyThrough) GetResults() ([]*eloquent.Model, error) {
	if r.query == nil {
		return nil, nil
	}

	rows, err := r.query.Get(context.Background())

	if err != nil {
		return nil, err
	}

	models := make([]*eloquent.Model, 0, len(rows))

	for _, row := range rows {
		m := eloquent.NewModel()
		m.SetTable(r.farParent.GetTable())
		m.SetRawAttributes(row, true)
		m.SetExists(true)
		models = append(models, m)
	}

	return models, nil
}
