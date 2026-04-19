package relations

import (
	"github.com/bedrock/packages/database/eloquent"
	"github.com/bedrock/packages/database/query"
)

// Relation is the interface all relationship types implement.
type Relation interface {
	// AddConstraints adds the base constraints to the relation query.
	AddConstraints()
	// AddEagerConstraints adds constraints for eager loading.
	AddEagerConstraints(models []*eloquent.Model)
	// InitRelation initializes the relation on a set of models.
	InitRelation(models []*eloquent.Model, relation string) []*eloquent.Model
	// Match matches eagerly loaded results to their parents.
	Match(models []*eloquent.Model, results []*eloquent.Model, relation string) []*eloquent.Model
	// GetResults returns the results of the relationship.
	GetResults() ([]*eloquent.Model, error)
	// GetQuery returns the underlying query builder.
	GetQuery() *query.Builder
}

// BaseRelation provides shared functionality for all relationship types.
type BaseRelation struct {
	query   *query.Builder
	parent  *eloquent.Model
	related *eloquent.Model
}

// NewBaseRelation creates a new base relation.
func NewBaseRelation(q *query.Builder, parent, related *eloquent.Model) *BaseRelation {
	return &BaseRelation{
		query:   q,
		parent:  parent,
		related: related,
	}
}

// GetQuery returns the relation's query builder.
func (r *BaseRelation) GetQuery() *query.Builder { return r.query }

// GetParent returns the parent model.
func (r *BaseRelation) GetParent() *eloquent.Model { return r.parent }

// GetRelated returns the related model.
func (r *BaseRelation) GetRelated() *eloquent.Model { return r.related }

// GetQualifiedParentKeyName returns the parent's qualified key name.
func (r *BaseRelation) GetQualifiedParentKeyName() string {
	return r.parent.GetQualifiedKeyName()
}
