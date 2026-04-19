package relations

import (
	"context"

	"github.com/bedrock/packages/database/eloquent"
)

// HasOneThrough defines a one-to-one relationship through an intermediate table.
type HasOneThrough struct {
	*BaseRelation
	throughParent  *eloquent.Model
	farParent      *eloquent.Model
	firstKey       string
	secondKey      string
	localKey       string
	secondLocalKey string
}

// NewHasOneThrough creates a new HasOneThrough relationship.
func NewHasOneThrough(parent, throughParent, farParent *eloquent.Model, firstKey, secondKey, localKey, secondLocalKey string) *HasOneThrough {
	return &HasOneThrough{
		BaseRelation:   NewBaseRelation(nil, parent, farParent),
		throughParent:  throughParent,
		farParent:      farParent,
		firstKey:       firstKey,
		secondKey:      secondKey,
		localKey:       localKey,
		secondLocalKey: secondLocalKey,
	}
}

func (r *HasOneThrough) GetFirstKeyName() string  { return r.firstKey }
func (r *HasOneThrough) GetSecondKeyName() string { return r.secondKey }

func (r *HasOneThrough) AddConstraints() {}

func (r *HasOneThrough) AddEagerConstraints(models []*eloquent.Model) {
	keys := make([]any, 0, len(models))

	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.localKey))
	}

	if r.query != nil {
		r.query.WhereIn(r.throughParent.GetTable()+"."+r.firstKey, keys)
	}
}

func (r *HasOneThrough) InitRelation(models []*eloquent.Model, relation string) []*eloquent.Model {
	return models
}

func (r *HasOneThrough) Match(models []*eloquent.Model, results []*eloquent.Model, relation string) []*eloquent.Model {
	dictionary := make(map[any]*eloquent.Model)

	for _, result := range results {
		key := result.GetAttribute("laravel_through_key")
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

func (r *HasOneThrough) GetResults() ([]*eloquent.Model, error) {
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
