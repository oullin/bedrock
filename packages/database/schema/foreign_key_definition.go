package schema

// ForeignKeyDefinition is a fluent builder for foreign key constraints.
type ForeignKeyDefinition struct {
	Columns        []string
	RefTable       string
	RefColumns     []string
	OnDeleteAction string
	OnUpdateAction string
	IndexName      string
}

// References sets the referenced columns.
func (f *ForeignKeyDefinition) References(columns ...string) *ForeignKeyDefinition {
	f.RefColumns = columns

	return f
}

// On sets the referenced table.
func (f *ForeignKeyDefinition) On(table string) *ForeignKeyDefinition {
	f.RefTable = table

	return f
}

// OnDelete sets the ON DELETE action.
func (f *ForeignKeyDefinition) OnDelete(action string) *ForeignKeyDefinition {
	f.OnDeleteAction = action

	return f
}

// OnUpdate sets the ON UPDATE action.
func (f *ForeignKeyDefinition) OnUpdate(action string) *ForeignKeyDefinition {
	f.OnUpdateAction = action

	return f
}

// CascadeOnDelete sets ON DELETE CASCADE.
func (f *ForeignKeyDefinition) CascadeOnDelete() *ForeignKeyDefinition {
	return f.OnDelete("cascade")
}

// RestrictOnDelete sets ON DELETE RESTRICT.
func (f *ForeignKeyDefinition) RestrictOnDelete() *ForeignKeyDefinition {
	return f.OnDelete("restrict")
}

// NullOnDelete sets ON DELETE SET NULL.
func (f *ForeignKeyDefinition) NullOnDelete() *ForeignKeyDefinition {
	return f.OnDelete("set null")
}

// NoActionOnDelete sets ON DELETE NO ACTION.
func (f *ForeignKeyDefinition) NoActionOnDelete() *ForeignKeyDefinition {
	return f.OnDelete("no action")
}

// CascadeOnUpdate sets ON UPDATE CASCADE.
func (f *ForeignKeyDefinition) CascadeOnUpdate() *ForeignKeyDefinition {
	return f.OnUpdate("cascade")
}

// RestrictOnUpdate sets ON UPDATE RESTRICT.
func (f *ForeignKeyDefinition) RestrictOnUpdate() *ForeignKeyDefinition {
	return f.OnUpdate("restrict")
}
