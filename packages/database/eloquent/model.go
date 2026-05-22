package eloquent

import (
	"context"
	"encoding/json"

	dbcontract "github.com/bedrock/packages/contracts/database"
	"github.com/bedrock/packages/database/query"
)

// Model is the base Eloquent model. User models embed this struct to gain
// Ref: @bedrock/code-0209
type Model struct {
	HasAttributes
	HasTimestamps
	GuardsAttributes
	HidesAttributes
	GlobalScopes
	SoftDeletes

	table              string
	primaryKey         string
	keyType            string
	incrementing       bool
	connection         string
	exists             bool
	wasRecentlyCreated bool
	perPage            int

	// Dependencies injected at construction time.
	resolver dbcontract.ConnectionResolver
}

// NewModel creates a new Model with defaults.
func NewModel() *Model {
	m := &Model{
		primaryKey:   "id",
		keyType:      "int",
		incrementing: true,
		perPage:      15,
	}
	m.InitAttributes()
	m.InitTimestamps()

	return m
}

// SetResolver sets the connection resolver (typically the DatabaseManager).
func (m *Model) SetResolver(resolver dbcontract.ConnectionResolver) {
	m.resolver = resolver
}

// GetResolver returns the connection resolver.
func (m *Model) GetResolver() dbcontract.ConnectionResolver {
	return m.resolver
}

// SetTable sets the table name.
func (m *Model) SetTable(table string) { m.table = table }

// GetTable returns the table name.
func (m *Model) GetTable() string {
	if m.table != "" {
		return m.table
	}

	return "models" // Default; in practice, derived from struct name.
}

// SetPrimaryKey sets the primary key name.
func (m *Model) SetPrimaryKey(key string) { m.primaryKey = key }

// GetKeyName returns the primary key column name.
func (m *Model) GetKeyName() string { return m.primaryKey }

// GetKey returns the value of the primary key.
func (m *Model) GetKey() any { return m.GetAttribute(m.primaryKey) }

// SetKeyType sets the primary key type.
func (m *Model) SetKeyType(keyType string) { m.keyType = keyType }

// GetKeyType returns the primary key type ("int" or "string").
func (m *Model) GetKeyType() string { return m.keyType }

// GetIncrementing returns whether the model uses auto-incrementing IDs.
func (m *Model) GetIncrementing() bool { return m.incrementing }

// SetIncrementing sets the incrementing flag.
func (m *Model) SetIncrementing(inc bool) { m.incrementing = inc }

// SetConnectionName sets the connection name for this model.
func (m *Model) SetConnectionName(name string) { m.connection = name }

// GetConnectionName returns the connection name.
func (m *Model) GetConnectionName() string { return m.connection }

// GetPerPage returns the number of results per page.
func (m *Model) GetPerPage() int { return m.perPage }

// SetPerPage sets the number of results per page.
func (m *Model) SetPerPage(n int) { m.perPage = n }

// Exists returns whether the model has been persisted.
func (m *Model) Exists() bool { return m.exists }

// SetExists sets the existence flag.
func (m *Model) SetExists(exists bool) { m.exists = exists }

// WasRecentlyCreated returns whether the model was recently created.
func (m *Model) WasRecentlyCreated() bool { return m.wasRecentlyCreated }

// GetQualifiedKeyName returns the table-qualified primary key name.
func (m *Model) GetQualifiedKeyName() string {
	return m.GetTable() + "." + m.GetKeyName()
}

// QualifyColumn returns a fully qualified column name.
func (m *Model) QualifyColumn(column string) string {
	return m.GetTable() + "." + column
}

// QualifyColumns qualifies multiple column names.
func (m *Model) QualifyColumns(columns []string) []string {
	qualified := make([]string, len(columns))

	for i, col := range columns {
		qualified[i] = m.QualifyColumn(col)
	}

	return qualified
}

// NewQueryBuilder creates a new query builder for the model's table.
func (m *Model) NewQueryBuilder() (*query.Builder, error) {
	if m.resolver == nil {
		return nil, ErrModelNotFound
	}

	conn, err := m.resolver.Connection(context.Background(), m.connection)

	if err != nil {
		return nil, err
	}

	result := conn.Table(context.Background(), m.GetTable())

	if builder, ok := result.(*query.Builder); ok {
		return builder, nil
	}

	return nil, nil
}

// Fill fills the model with attributes, respecting mass assignment.
func (m *Model) Fill(values map[string]any) error {
	return m.GuardsAttributes.Fill(&m.HasAttributes, values)
}

// ForceFill fills the model ignoring mass assignment.
func (m *Model) ForceFill(values map[string]any) {
	m.GuardsAttributes.ForceFill(&m.HasAttributes, values)
}

// Save persists the model to the database.
func (m *Model) Save(ctx context.Context) error {
	if m.resolver == nil {
		return ErrModelNotFound
	}

	conn, err := m.resolver.Connection(ctx, m.connection)

	if err != nil {
		return err
	}

	if m.exists {
		return m.performUpdate(ctx, conn)
	}

	return m.performInsert(ctx, conn)
}

// Delete removes the model from the database.
func (m *Model) Delete(ctx context.Context) error {
	if m.resolver == nil || !m.exists {
		return nil
	}

	conn, err := m.resolver.Connection(ctx, m.connection)

	if err != nil {
		return err
	}

	keyValue := m.GetKey()

	if keyValue == nil {
		return ErrModelNotFound
	}

	_, err = conn.Delete(ctx,
		"delete from "+m.GetTable()+" where "+m.GetKeyName()+" = ?",
		keyValue,
	)

	if err != nil {
		return err
	}

	m.exists = false

	return nil
}

// Refresh reloads the model from the database.
func (m *Model) Refresh(ctx context.Context) error {
	if m.resolver == nil || !m.exists {
		return ErrModelNotFound
	}

	conn, err := m.resolver.Connection(ctx, m.connection)

	if err != nil {
		return err
	}

	row, err := conn.SelectOne(ctx,
		"select * from "+m.GetTable()+" where "+m.GetKeyName()+" = ? limit 1",
		m.GetKey(),
	)

	if err != nil {
		return err
	}

	if row == nil {
		return ErrModelNotFound
	}

	m.SetRawAttributes(row, true)

	return nil
}

// Replicate creates a copy of the model without the primary key.
func (m *Model) Replicate(except ...string) *Model {
	clone := NewModel()
	clone.table = m.table
	clone.primaryKey = m.primaryKey
	clone.keyType = m.keyType
	clone.incrementing = m.incrementing
	clone.connection = m.connection
	clone.resolver = m.resolver
	clone.fillable = m.fillable
	clone.guarded = m.guarded
	clone.casts = m.casts
	clone.HasTimestamps = m.HasTimestamps
	clone.HidesAttributes = m.HidesAttributes

	attrs := make(map[string]any, len(m.attributes))

	for k, v := range m.attributes {
		attrs[k] = v
	}

	// Remove the primary key and any excluded attributes.
	delete(attrs, m.primaryKey)

	for _, e := range except {
		delete(attrs, e)
	}

	clone.SetRawAttributes(attrs, false)

	return clone
}

// Is checks if two models represent the same database record.
func (m *Model) Is(other *Model) bool {
	if other == nil {
		return false
	}

	return m.GetTable() == other.GetTable() &&
		m.GetKeyName() == other.GetKeyName() &&
		m.GetKey() == other.GetKey() &&
		m.GetConnectionName() == other.GetConnectionName()
}

// IsNot checks if two models represent different records.
func (m *Model) IsNot(other *Model) bool {
	return !m.Is(other)
}

// ToMap returns the model as a map, applying hidden/visible rules.
func (m *Model) ToMap() map[string]any {
	attrs := m.HasAttributes.ToMap()

	return m.HidesAttributes.FilterAttributes(attrs)
}

// ToJSON serializes the model as JSON.
func (m *Model) ToJSON() ([]byte, error) {
	return json.Marshal(m.ToMap())
}

// MarshalJSON implements json.Marshaler.
func (m *Model) MarshalJSON() ([]byte, error) {
	return m.ToJSON()
}

// SaveOrFail saves the model and returns an error if it fails.
func (m *Model) SaveOrFail(ctx context.Context) error {
	return m.Save(ctx)
}

// SaveQuietly saves the model without firing events.
func (m *Model) SaveQuietly(ctx context.Context) error {
	return m.Save(ctx)
}

// UpdateAttributes updates the model with the given attributes and saves.
func (m *Model) UpdateAttributes(ctx context.Context, attributes map[string]any) error {
	m.ForceFill(attributes)

	return m.Save(ctx)
}

// UpdateOrFail updates the model and returns an error if it fails.
func (m *Model) UpdateOrFail(ctx context.Context, attributes map[string]any) error {
	return m.UpdateAttributes(ctx, attributes)
}

// UpdateQuietly updates the model without firing events.
func (m *Model) UpdateQuietly(ctx context.Context, attributes map[string]any) error {
	return m.UpdateAttributes(ctx, attributes)
}

// DeleteOrFail deletes the model and returns an error if not found.
func (m *Model) DeleteOrFail(ctx context.Context) error {
	if !m.exists {
		return ErrModelNotFound
	}

	return m.Delete(ctx)
}

// DeleteQuietly deletes the model without firing events.
func (m *Model) DeleteQuietly(ctx context.Context) error {
	return m.Delete(ctx)
}

// ForceDelete hard-deletes the model (bypassing soft deletes).
func (m *Model) ForceDelete(ctx context.Context) error {
	m.SoftDeletes.forceDeleting = true

	defer func() { m.SoftDeletes.forceDeleting = false }()

	return m.Delete(ctx)
}

// Destroy deletes models by their primary keys.
func Destroy(ctx context.Context, resolver dbcontract.ConnectionResolver, table, keyName, connection string, ids ...any) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	conn, err := resolver.Connection(ctx, connection)

	if err != nil {
		return 0, err
	}

	placeholders := makePlaceholders(len(ids))
	sql := "delete from " + table + " where " + keyName + " in (" + placeholders + ")"

	return conn.Delete(ctx, sql, ids...)
}

// Touch updates the model's timestamps.
func (m *Model) Touch(ctx context.Context) error {
	if !m.UsesTimestamps() {
		return nil
	}

	m.SetAttribute(m.GetUpdatedAtColumn(), m.FreshTimestampString())

	return m.Save(ctx)
}

// Fresh returns a new model instance loaded from the database.
func (m *Model) Fresh(ctx context.Context, columns ...string) (*Model, error) {
	if !m.exists || m.resolver == nil {
		return nil, ErrModelNotFound
	}

	conn, err := m.resolver.Connection(ctx, m.connection)

	if err != nil {
		return nil, err
	}

	cols := "*"

	if len(columns) > 0 {
		cols = joinStrings(columns, ", ")
	}

	row, err := conn.SelectOne(ctx,
		"select "+cols+" from "+m.GetTable()+" where "+m.GetKeyName()+" = ? limit 1",
		m.GetKey(),
	)

	if err != nil {
		return nil, err
	}

	if row == nil {
		return nil, ErrModelNotFound
	}

	fresh := NewModel()
	fresh.table = m.table
	fresh.primaryKey = m.primaryKey
	fresh.keyType = m.keyType
	fresh.incrementing = m.incrementing
	fresh.connection = m.connection
	fresh.resolver = m.resolver
	fresh.SetRawAttributes(row, true)
	fresh.exists = true

	return fresh, nil
}

// NewInstance creates a new unsaved model instance with the same config.
func (m *Model) NewInstance(attributes ...map[string]any) *Model {
	instance := NewModel()
	instance.table = m.table
	instance.primaryKey = m.primaryKey
	instance.keyType = m.keyType
	instance.incrementing = m.incrementing
	instance.connection = m.connection
	instance.resolver = m.resolver

	if len(attributes) > 0 {
		instance.ForceFill(attributes[0])
	}

	return instance
}

// NewFromBuilder creates a model instance from raw builder results.
func (m *Model) NewFromBuilder(attributes map[string]any) *Model {
	instance := m.NewInstance()
	instance.SetRawAttributes(attributes, true)
	instance.exists = true

	return instance
}

// Push saves the model and all its loaded relationships.
func (m *Model) Push(ctx context.Context) error {
	return m.Save(ctx)
}

// PushQuietly saves the model and relationships without events.
func (m *Model) PushQuietly(ctx context.Context) error {
	return m.Push(ctx)
}

// ReplicateQuietly creates a replica without events.
func (m *Model) ReplicateQuietly(except ...string) *Model {
	return m.Replicate(except...)
}

// Increment increments a column and saves.
func (m *Model) Increment(ctx context.Context, column string, amount ...any) error {
	amt := any(1)

	if len(amount) > 0 {
		amt = amount[0]
	}

	current := m.GetAttribute(column)

	var newVal any

	switch v := current.(type) {
	case int64:
		if a, ok := amt.(int); ok {
			newVal = v + int64(a)
		} else {
			newVal = v + 1
		}
	case float64:
		if a, ok := amt.(float64); ok {
			newVal = v + a
		} else {
			newVal = v + 1
		}
	default:
		newVal = amt
	}

	m.SetAttribute(column, newVal)

	return m.Save(ctx)
}

// Decrement decrements a column and saves.
func (m *Model) Decrement(ctx context.Context, column string, amount ...any) error {
	amt := any(1)

	if len(amount) > 0 {
		amt = amount[0]
	}

	current := m.GetAttribute(column)

	var newVal any

	switch v := current.(type) {
	case int64:
		if a, ok := amt.(int); ok {
			newVal = v - int64(a)
		} else {
			newVal = v - 1
		}
	case float64:
		if a, ok := amt.(float64); ok {
			newVal = v - a
		} else {
			newVal = v - 1
		}
	default:
		newVal = amt
	}

	m.SetAttribute(column, newVal)

	return m.Save(ctx)
}

// ToPrettyJSON serializes the model as indented JSON.
func (m *Model) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(m.ToMap(), "", "  ")
}

func (m *Model) performInsert(ctx context.Context, conn dbcontract.Connection) error {
	attrs := m.GetDirtyForSave()

	if m.UsesTimestamps() {
		m.HasTimestamps.UpdateTimestamps(attrs, true)

		for k, v := range attrs {
			m.SetAttribute(k, v)
		}
	}

	if len(attrs) == 0 {
		return nil
	}

	columns, values := mapToColumnsValues(attrs)
	placeholders := makePlaceholders(len(values))

	sql := "insert into " + m.GetTable() +
		" (" + joinStrings(columns, ", ") + ")" +
		" values (" + placeholders + ")"

	ok, err := conn.Insert(ctx, sql, values...)

	if err != nil {
		return err
	}

	if !ok {
		return ErrModelNotFound
	}

	m.exists = true
	m.wasRecentlyCreated = true
	m.SyncOriginal()

	return nil
}

func (m *Model) performUpdate(ctx context.Context, conn dbcontract.Connection) error {
	dirty := m.GetDirtyForSave()

	if m.UsesTimestamps() {
		m.HasTimestamps.UpdateTimestamps(dirty, false)

		for k, v := range dirty {
			m.SetAttribute(k, v)
		}
	}

	if len(dirty) == 0 {
		return nil
	}

	columns, values := mapToColumnsValues(dirty)

	var sets []string

	for _, col := range columns {
		sets = append(sets, col+" = ?")
	}

	values = append(values, m.GetKey())

	sql := "update " + m.GetTable() +
		" set " + joinStrings(sets, ", ") +
		" where " + m.GetKeyName() + " = ?"

	_, err := conn.Update(ctx, sql, values...)

	if err != nil {
		return err
	}

	m.SyncOriginal()

	return nil
}

// GetDirtyForSave returns the dirty attributes for a save operation.
func (m *Model) GetDirtyForSave() map[string]any {
	if !m.exists {
		return m.GetAttributes()
	}

	return m.GetDirty()
}

func mapToColumnsValues(m map[string]any) ([]string, []any) {
	keys := sortedMapKeys(m)
	values := make([]any, len(keys))

	for i, k := range keys {
		values[i] = m[k]
	}

	return keys, values
}

func sortedMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}

	return keys
}

func makePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	s := "?"

	for i := 1; i < n; i++ {
		s += ", ?"
	}

	return s
}

func joinStrings(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}

	result := s[0]

	for _, v := range s[1:] {
		result += sep + v
	}

	return result
}
