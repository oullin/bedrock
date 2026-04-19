package relations

import "github.com/bedrock/packages/database/orm"

// Pivot represents a pivot table record in a many-to-many relationship.
type Pivot struct {
	*orm.Model
	pivotParent *orm.Model
	foreignKey  string
	relatedKey  string
}

// NewPivot creates a new Pivot model.

// GetForeignKey returns the foreign key column name.

// GetRelatedKey returns the related key column name.

// GetPivotParent returns the parent of this pivot.

// MorphPivot represents a pivot table record in a polymorphic many-to-many.
type MorphPivot struct {
	*Pivot
	morphType  string
	morphClass string
}

func NewPivot(parent *orm.Model, attributes map[string]any, table, foreignKey, relatedKey string) *Pivot {
	p := &Pivot{
		Model:       orm.NewModel(),
		pivotParent: parent,
		foreignKey:  foreignKey,
		relatedKey:  relatedKey,
	}
	p.SetTable(table)
	p.SetIncrementing(false)
	p.SetTimestamps(false)

	if attributes != nil {
		p.SetRawAttributes(attributes, true)
		p.SetExists(true)
	}

	return p
}

func (p *Pivot) GetForeignKey() string { return p.foreignKey }

func (p *Pivot) GetRelatedKey() string { return p.relatedKey }

func (p *Pivot) GetPivotParent() *orm.Model { return p.pivotParent }

// NewMorphPivot creates a new MorphPivot model.
func NewMorphPivot(parent *orm.Model, attributes map[string]any, table, foreignKey, relatedKey, morphType, morphClass string) *MorphPivot {
	return &MorphPivot{
		Pivot:      NewPivot(parent, attributes, table, foreignKey, relatedKey),
		morphType:  morphType,
		morphClass: morphClass,
	}
}

// GetMorphType returns the morph type column name.
func (p *MorphPivot) GetMorphType() string { return p.morphType }

// GetMorphClass returns the morph class value.
func (p *MorphPivot) GetMorphClass() string { return p.morphClass }
