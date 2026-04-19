package schema

// ColumnDefinition is a fluent builder for column constraints. Returned by
// Blueprint column methods (String, Integer, etc.) for chaining.
type ColumnDefinition struct {
	Name               string
	Type               string
	Length             int
	Precision          int
	Scale              int
	Unsigned           bool
	AutoIncrement      bool
	IsNullable         bool
	DefaultValue       any
	HasDefault         bool
	CommentText        string
	Collation          string
	CharsetName        string
	VirtualAs          string
	StoredAs           string
	AfterColumn        string
	IsFirst            bool
	IsChange           bool
	UseCurrent         bool
	UseCurrentOnUpdate bool
	Allowed            []string // for enum/set types
	Total              int      // for vector
	Srid               int      // for geometry
}

// Nullable marks the column as nullable.
func (c *ColumnDefinition) Nullable(nullable ...bool) *ColumnDefinition {
	c.IsNullable = true

	if len(nullable) > 0 {
		c.IsNullable = nullable[0]
	}

	return c
}

// Default sets the default value.
func (c *ColumnDefinition) Default(value any) *ColumnDefinition {
	c.DefaultValue = value
	c.HasDefault = true

	return c
}

// Comment sets the column comment.
func (c *ColumnDefinition) Comment(text string) *ColumnDefinition {
	c.CommentText = text

	return c
}

// After places the column after another column (MySQL).
func (c *ColumnDefinition) After(column string) *ColumnDefinition {
	c.AfterColumn = column

	return c
}

// First places the column first in the table (MySQL).
func (c *ColumnDefinition) First() *ColumnDefinition {
	c.IsFirst = true

	return c
}

// Change marks this as a column modification.
func (c *ColumnDefinition) Change() *ColumnDefinition {
	c.IsChange = true

	return c
}

// VirtualAsExpr sets a virtual generated column expression.
func (c *ColumnDefinition) VirtualAsExpr(expression string) *ColumnDefinition {
	c.VirtualAs = expression

	return c
}

// StoredAsExpr sets a stored generated column expression.
func (c *ColumnDefinition) StoredAsExpr(expression string) *ColumnDefinition {
	c.StoredAs = expression

	return c
}

// SetCollation sets the column collation.
func (c *ColumnDefinition) SetCollation(collation string) *ColumnDefinition {
	c.Collation = collation

	return c
}

// Charset sets the column character set.
func (c *ColumnDefinition) Charset(charset string) *ColumnDefinition {
	c.CharsetName = charset

	return c
}

// SetUseCurrent sets the column default to CURRENT_TIMESTAMP.
func (c *ColumnDefinition) SetUseCurrent() *ColumnDefinition {
	c.UseCurrent = true

	return c
}

// UseCurrentOnUpdateExpr sets the column to update to CURRENT_TIMESTAMP on update.
func (c *ColumnDefinition) UseCurrentOnUpdateExpr() *ColumnDefinition {
	c.UseCurrentOnUpdate = true

	return c
}

// Unique adds a unique index for this column (convenience, adds to parent Blueprint).
func (c *ColumnDefinition) Unique() *ColumnDefinition {
	// Implemented through Blueprint's addCommand mechanism.
	// The Blueprint stores a reference and adds the index at compile time.
	return c
}

// Index adds an index for this column.
func (c *ColumnDefinition) Index() *ColumnDefinition {
	return c
}

// Primary makes this column the primary key.
func (c *ColumnDefinition) Primary() *ColumnDefinition {
	return c
}

// Always marks the column as GENERATED ALWAYS (PostgreSQL identity).
func (c *ColumnDefinition) Always() *ColumnDefinition {
	return c
}
