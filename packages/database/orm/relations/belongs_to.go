package relations

import (
	"context"

	"github.com/bedrock/packages/database/orm"
)

// BelongsTo defines the inverse of a one-to-one or one-to-many relationship.
type BelongsTo struct {
	*BaseRelation
	foreignKey string
	ownerKey   string
}

// NewBelongsTo creates a new BelongsTo relationship.
func NewBelongsTo(parent, related *orm.Model, foreignKey, ownerKey string) *BelongsTo {
	return &BelongsTo{
		BaseRelation: NewBaseRelation(nil, parent, related),
		foreignKey:   foreignKey,
		ownerKey:     ownerKey,
	}
}

// GetForeignKeyName returns the foreign key column name.
func (r *BelongsTo) GetForeignKeyName() string { return r.foreignKey }

// GetOwnerKeyName returns the owner key column name.
func (r *BelongsTo) GetOwnerKeyName() string { return r.ownerKey }

func (r *BelongsTo) AddConstraints() {
	if r.query != nil {
		r.query.Where(r.ownerKey, r.parent.GetAttribute(r.foreignKey))
	}
}

func (r *BelongsTo) AddEagerConstraints(models []*orm.Model) {
	keys := make([]any, 0, len(models))
	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.foreignKey))
	}
	if r.query != nil {
		r.query.WhereIn(r.ownerKey, keys)
	}
}

func (r *BelongsTo) InitRelation(models []*orm.Model, relation string) []*orm.Model {
	return models
}

func (r *BelongsTo) Match(models []*orm.Model, results []*orm.Model, relation string) []*orm.Model {
	dictionary := make(map[any]*orm.Model)
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

func (r *BelongsTo) GetResults() ([]*orm.Model, error) {
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
		m.SetTable(r.related.GetTable())
		m.SetRawAttributes(row, true)
		m.SetExists(true)
		models = append(models, m)
	}
	return models, nil
}

// Associate sets the foreign key on the parent model.
func (r *BelongsTo) Associate(model *orm.Model) {
	r.parent.SetAttribute(r.foreignKey, model.GetAttribute(r.ownerKey))
}

// Dissociate clears the foreign key on the parent model.
func (r *BelongsTo) Dissociate() {
	r.parent.SetAttribute(r.foreignKey, nil)
}
