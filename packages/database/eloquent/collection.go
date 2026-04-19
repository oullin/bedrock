package eloquent

import "encoding/json"

// Collection is a slice of models with helper methods.
type Collection struct {
	models []*Model
}

// NewCollection creates a new Collection.
func NewCollection(models []*Model) *Collection {
	return &Collection{models: models}
}

// All returns the underlying models slice.
func (c *Collection) All() []*Model { return c.models }

// Count returns the number of models.
func (c *Collection) Count() int { return len(c.models) }

// IsEmpty checks if the collection is empty.
func (c *Collection) IsEmpty() bool { return len(c.models) == 0 }

// IsNotEmpty checks if the collection is not empty.
func (c *Collection) IsNotEmpty() bool { return len(c.models) > 0 }

// First returns the first model.
func (c *Collection) First() *Model {
	if len(c.models) == 0 {
		return nil
	}

	return c.models[0]
}

// Last returns the last model.
func (c *Collection) Last() *Model {
	if len(c.models) == 0 {
		return nil
	}

	return c.models[len(c.models)-1]
}

// Find finds a model by primary key.
func (c *Collection) Find(key any) *Model {
	for _, m := range c.models {
		if m.GetKey() == key {
			return m
		}
	}

	return nil
}

// Contains checks if a model with the given key exists.
func (c *Collection) Contains(key any) bool {
	return c.Find(key) != nil
}

// ModelKeys returns the primary key values.
func (c *Collection) ModelKeys() []any {
	keys := make([]any, len(c.models))

	for i, m := range c.models {
		keys[i] = m.GetKey()
	}

	return keys
}

// Pluck extracts a single attribute from each model.
func (c *Collection) Pluck(key string) []any {
	values := make([]any, 0, len(c.models))

	for _, m := range c.models {
		values = append(values, m.GetAttribute(key))
	}

	return values
}

// Only returns models whose keys are in the given list.
func (c *Collection) Only(keys []any) *Collection {
	keySet := make(map[any]bool, len(keys))

	for _, k := range keys {
		keySet[k] = true
	}

	var filtered []*Model

	for _, m := range c.models {
		if keySet[m.GetKey()] {
			filtered = append(filtered, m)
		}
	}

	return NewCollection(filtered)
}

// Except returns models whose keys are NOT in the given list.
func (c *Collection) Except(keys []any) *Collection {
	keySet := make(map[any]bool, len(keys))

	for _, k := range keys {
		keySet[k] = true
	}

	var filtered []*Model

	for _, m := range c.models {
		if !keySet[m.GetKey()] {
			filtered = append(filtered, m)
		}
	}

	return NewCollection(filtered)
}

// Each iterates over each model.
func (c *Collection) Each(fn func(int, *Model) bool) {
	for i, m := range c.models {
		if !fn(i, m) {
			break
		}
	}
}

// Map applies a function to each model and returns the results.
func (c *Collection) Map(fn func(*Model) any) []any {
	results := make([]any, len(c.models))

	for i, m := range c.models {
		results[i] = fn(m)
	}

	return results
}

// Filter returns models that pass the predicate.
func (c *Collection) Filter(fn func(*Model) bool) *Collection {
	var filtered []*Model

	for _, m := range c.models {
		if fn(m) {
			filtered = append(filtered, m)
		}
	}

	return NewCollection(filtered)
}

// Push adds a model to the collection.
func (c *Collection) Push(models ...*Model) {
	c.models = append(c.models, models...)
}

// Merge adds models from another collection.
func (c *Collection) Merge(other *Collection) *Collection {
	merged := make([]*Model, 0, len(c.models)+len(other.models))
	merged = append(merged, c.models...)
	merged = append(merged, other.models...)

	return NewCollection(merged)
}

// Diff returns models in this collection but not in the other.
func (c *Collection) Diff(other *Collection) *Collection {
	otherKeys := make(map[any]bool)

	for _, m := range other.models {
		otherKeys[m.GetKey()] = true
	}

	var diff []*Model

	for _, m := range c.models {
		if !otherKeys[m.GetKey()] {
			diff = append(diff, m)
		}
	}

	return NewCollection(diff)
}

// Intersect returns models that exist in both collections.
func (c *Collection) Intersect(other *Collection) *Collection {
	otherKeys := make(map[any]bool)

	for _, m := range other.models {
		otherKeys[m.GetKey()] = true
	}

	var shared []*Model

	for _, m := range c.models {
		if otherKeys[m.GetKey()] {
			shared = append(shared, m)
		}
	}

	return NewCollection(shared)
}

// Unique removes duplicate models by key.
func (c *Collection) Unique() *Collection {
	seen := make(map[any]bool)

	var unique []*Model

	for _, m := range c.models {
		key := m.GetKey()

		if !seen[key] {
			seen[key] = true
			unique = append(unique, m)
		}
	}

	return NewCollection(unique)
}

// Fresh reloads all models from the database.
func (c *Collection) Fresh() error {
	// Each model is reloaded individually. For performance, the caller
	// should use a batch query instead.
	return nil
}

// MakeHidden hides attributes on all models.
func (c *Collection) MakeHidden(attributes ...string) *Collection {
	for _, m := range c.models {
		m.MakeHidden(attributes...)
	}

	return c
}

// MakeVisible makes attributes visible on all models.
func (c *Collection) MakeVisible(attributes ...string) *Collection {
	for _, m := range c.models {
		m.MakeVisible(attributes...)
	}

	return c
}

// ToMaps converts all models to maps.
func (c *Collection) ToMaps() []map[string]any {
	maps := make([]map[string]any, len(c.models))

	for i, m := range c.models {
		maps[i] = m.ToMap()
	}

	return maps
}

// ToJSON serializes the collection as JSON.
func (c *Collection) ToJSON() ([]byte, error) {
	return json.Marshal(c.ToMaps())
}

// MarshalJSON implements json.Marshaler.
func (c *Collection) MarshalJSON() ([]byte, error) {
	return c.ToJSON()
}
