package relations

import (
	"github.com/bedrock/packages/database/orm"
	"github.com/bedrock/packages/database/query"
)

// Relation is the interface all relationship types implement.
type Relation interface {
	// AddConstraints adds the base constraints to the relation query.
	AddConstraints()
	// AddEagerConstraints adds constraints for eager loading.
	AddEagerConstraints(models []*orm.Model)
	// InitRelation initializes the relation on a set of models.
	InitRelation(models []*orm.Model, relation string) []*orm.Model
	// Match matches eagerly loaded results to their parents.
	Match(models []*orm.Model, results []*orm.Model, relation string) []*orm.Model
	// GetResults returns the results of the relationship.
	GetResults() ([]*orm.Model, error)
	// GetQuery returns the underlying query builder.
	GetQuery() *query.Builder
}

// BaseRelation provides shared functionality for all relationship types.
type BaseRelation struct {
	query   *query.Builder
	parent  *orm.Model
	related *orm.Model
}

// NewBaseRelation creates a new base relation.
func NewBaseRelation(q *query.Builder, parent, related *orm.Model) *BaseRelation {
	return &BaseRelation{
		query:   q,
		parent:  parent,
		related: related,
	}
}

// GetQuery returns the relation's query builder.
func (r *BaseRelation) GetQuery() *query.Builder { return r.query }

// GetParent returns the parent model.
func (r *BaseRelation) GetParent() *orm.Model { return r.parent }

// GetRelated returns the related model.
func (r *BaseRelation) GetRelated() *orm.Model { return r.related }

// GetQualifiedParentKeyName returns the parent's qualified key name.
func (r *BaseRelation) GetQualifiedParentKeyName() string {
	return r.parent.GetQualifiedKeyName()
}
