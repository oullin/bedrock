package relations

import (
	"context"

	"github.com/bedrock/packages/database/eloquent"
)

// HasMany defines a one-to-many relationship.
type HasMany struct {
	*BaseRelation
	foreignKey string
	localKey   string
}

// NewHasMany creates a new HasMany relationship.
func NewHasMany(parent, related *eloquent.Model, foreignKey, localKey string) *HasMany {
	return &HasMany{
		BaseRelation: NewBaseRelation(nil, parent, related),
		foreignKey:   foreignKey,
		localKey:     localKey,
	}
}

// GetForeignKeyName returns the foreign key column name.
func (r *HasMany) GetForeignKeyName() string { return r.foreignKey }

// GetLocalKeyName returns the local key column name.
func (r *HasMany) GetLocalKeyName() string { return r.localKey }

func (r *HasMany) AddConstraints() {
	if r.query != nil {
		r.query.Where(r.foreignKey, r.parent.GetAttribute(r.localKey))
	}
}

func (r *HasMany) AddEagerConstraints(models []*eloquent.Model) {
	keys := make([]any, 0, len(models))
	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.localKey))
	}
	if r.query != nil {
		r.query.WhereIn(r.foreignKey, keys)
	}
}

func (r *HasMany) InitRelation(models []*eloquent.Model, relation string) []*eloquent.Model {
	for _, model := range models {
		model.SetAttribute(relation, []*eloquent.Model{})
	}
	return models
}

func (r *HasMany) Match(models []*eloquent.Model, results []*eloquent.Model, relation string) []*eloquent.Model {
	dictionary := make(map[any][]*eloquent.Model)
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

func (r *HasMany) GetResults() ([]*eloquent.Model, error) {
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
