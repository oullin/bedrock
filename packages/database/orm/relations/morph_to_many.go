package relations

import (
	"context"

	"github.com/bedrock/packages/database/orm"
)

// MorphToMany defines a polymorphic many-to-many relationship.
type MorphToMany struct {
	*BaseRelation
	pivotTable       string
	foreignPivotKey  string
	relatedPivotKey  string
	parentKey        string
	relatedKey       string
	morphType        string
	morphClass       string
	inverse          bool
}

// NewMorphToMany creates a new MorphToMany relationship.
func NewMorphToMany(parent, related *orm.Model, pivotTable, foreignPivotKey, relatedPivotKey, parentKey, relatedKey, morphType, morphClass string, inverse bool) *MorphToMany {
	return &MorphToMany{
		BaseRelation:    NewBaseRelation(nil, parent, related),
		pivotTable:      pivotTable,
		foreignPivotKey: foreignPivotKey,
		relatedPivotKey: relatedPivotKey,
		parentKey:       parentKey,
		relatedKey:      relatedKey,
		morphType:       morphType,
		morphClass:      morphClass,
		inverse:         inverse,
	}
}

func (r *MorphToMany) GetPivotTable() string { return r.pivotTable }

func (r *MorphToMany) AddConstraints() {}

func (r *MorphToMany) AddEagerConstraints(models []*orm.Model) {
	keys := make([]any, 0, len(models))
	for _, m := range models {
		keys = append(keys, m.GetAttribute(r.parentKey))
	}
	if r.query != nil {
		r.query.WhereIn(r.pivotTable+"."+r.foreignPivotKey, keys)
		r.query.Where(r.pivotTable+"."+r.morphType, r.morphClass)
	}
}

func (r *MorphToMany) InitRelation(models []*orm.Model, relation string) []*orm.Model {
	for _, model := range models {
		model.SetAttribute(relation, []*orm.Model{})
	}
	return models
}

func (r *MorphToMany) Match(models []*orm.Model, results []*orm.Model, relation string) []*orm.Model {
	dictionary := make(map[any][]*orm.Model)
	for _, result := range results {
		key := result.GetAttribute("pivot_"+r.foreignPivotKey)
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

func (r *MorphToMany) GetResults() ([]*orm.Model, error) {
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
