package jsonapi

// Relation represents a JSON:API relationship linkage object.
type Relation struct {
	Name   string
	Type   string
	IDs    []string
	IsMany bool
}

// HasOne creates a to-one relationship linkage.

// HasMany creates a to-many relationship linkage.

// RelationResolver maps a set of Relation values into JSON:API relationship
// objects suitable for the "relationships" key of a resource object.
type RelationResolver struct {
	relations []Relation
}

func HasOne(name, resourceType, id string) Relation {
	return Relation{
		Name:   name,
		Type:   resourceType,
		IDs:    []string{id},
		IsMany: false,
	}
}

func HasMany(name, resourceType string, ids []string) Relation {
	return Relation{
		Name:   name,
		Type:   resourceType,
		IDs:    ids,
		IsMany: true,
	}
}

// NewRelationResolver creates a RelationResolver from the given relations.
func NewRelationResolver(relations ...Relation) *RelationResolver {
	return &RelationResolver{relations: relations}
}

// Resolve produces the "relationships" map for a JSON:API resource object. Each
// entry is keyed by the relation name and contains a "data" key with either a
// single linkage object or an array of linkage objects.
func (rr *RelationResolver) Resolve() map[string]any {
	if len(rr.relations) == 0 {
		return nil
	}

	result := make(map[string]any, len(rr.relations))

	for _, rel := range rr.relations {
		if rel.IsMany {
			data := make([]map[string]any, len(rel.IDs))

			for i, id := range rel.IDs {
				data[i] = map[string]any{"type": rel.Type, "id": id}
			}

			result[rel.Name] = map[string]any{"data": data}
		} else {
			var data any

			if len(rel.IDs) > 0 {
				data = map[string]any{"type": rel.Type, "id": rel.IDs[0]}
			}

			result[rel.Name] = map[string]any{"data": data}
		}
	}

	return result
}
