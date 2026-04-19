package relations

import (
	"context"

	"github.com/bedrock/packages/database/eloquent"
)

// BelongsToMany defines a many-to-many relationship through a pivot table.
type BelongsToMany struct {
	*BaseRelation
	pivotTable      string
	foreignPivotKey string
	relatedPivotKey string
	parentKey       string
	relatedKey      string
	pivotColumns    []string
}

// NewBelongsToMany creates a new BelongsToMany relationship.
func NewBelongsToMany(parent, related *eloquent.Model, pivotTable, foreignPivotKey, relatedPivotKey, parentKey, relatedKey string) *BelongsToMany {
	return &BelongsToMany{
		BaseRelation:    NewBaseRelation(nil, parent, related),
		pivotTable:      pivotTable,
		foreignPivotKey: foreignPivotKey,
		relatedPivotKey: relatedPivotKey,
		parentKey:       parentKey,
		relatedKey:      relatedKey,
	}
}

// GetPivotTable returns the pivot table name.
func (r *BelongsToMany) GetPivotTable() string { return r.pivotTable }

// WithPivot specifies additional pivot columns to retrieve.
func (r *BelongsToMany) WithPivot(columns ...string) *BelongsToMany {
	r.pivotColumns = append(r.pivotColumns, columns...)

	return r
}

func (r *BelongsToMany) AddConstraints() {}

func (r *BelongsToMany) AddEagerConstraints(models []*eloquent.Model) {
	keys := make([]any, 0, len(models))

	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.parentKey))
	}

	if r.query != nil {
		r.query.WhereIn(r.pivotTable+"."+r.foreignPivotKey, keys)
	}
}

func (r *BelongsToMany) InitRelation(models []*eloquent.Model, relation string) []*eloquent.Model {
	for _, model := range models {
		model.SetAttribute(relation, []*eloquent.Model{})
	}

	return models
}

func (r *BelongsToMany) Match(models []*eloquent.Model, results []*eloquent.Model, relation string) []*eloquent.Model {
	dictionary := make(map[any][]*eloquent.Model)

	for _, result := range results {
		key := result.GetAttribute("pivot_" + r.foreignPivotKey)
		dictionary[key] = append(dictionary[key], result)
	}

	for _, model := range models {
		key := model.GetAttribute(r.parentKey)

		if matches, ok := dictionary[key]; ok {
			model.SetAttribute(relation, matches)
		}
	}

	return models
}

func (r *BelongsToMany) GetResults() ([]*eloquent.Model, error) {
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
