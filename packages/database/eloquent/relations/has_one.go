package relations

import (
	"context"

	"github.com/bedrock/packages/database/eloquent"
)

// HasOne defines a one-to-one relationship.
type HasOne struct {
	*BaseRelation
	foreignKey string
	localKey   string
}

// NewHasOne creates a new HasOne relationship.
func NewHasOne(parent, related *eloquent.Model, foreignKey, localKey string) *HasOne {
	return &HasOne{
		BaseRelation: NewBaseRelation(nil, parent, related),
		foreignKey:   foreignKey,
		localKey:     localKey,
	}
}

// GetForeignKeyName returns the foreign key column name.
func (r *HasOne) GetForeignKeyName() string { return r.foreignKey }

// GetLocalKeyName returns the local key column name.
func (r *HasOne) GetLocalKeyName() string { return r.localKey }

func (r *HasOne) AddConstraints() {
	if r.query != nil {
		r.query.Where(r.foreignKey, r.parent.GetAttribute(r.localKey))
	}
}

func (r *HasOne) AddEagerConstraints(models []*eloquent.Model) {
	keys := make([]any, 0, len(models))

	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.localKey))
	}

	if r.query != nil {
		r.query.WhereIn(r.foreignKey, keys)
	}
}

func (r *HasOne) InitRelation(models []*eloquent.Model, relation string) []*eloquent.Model {
	return models
}

func (r *HasOne) Match(models []*eloquent.Model, results []*eloquent.Model, relation string) []*eloquent.Model {
	dictionary := make(map[any]*eloquent.Model)

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

func (r *HasOne) GetResults() ([]*eloquent.Model, error) {
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
		m.SetTable(r.related.GetTable())
		m.SetRawAttributes(row, true)
		m.SetExists(true)
		models = append(models, m)
	}

	return models, nil
}
