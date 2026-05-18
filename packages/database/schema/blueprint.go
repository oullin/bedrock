package schema

// BlueprintCommand represents a DDL command (create index, add foreign key, etc.)
type BlueprintCommand struct {
	Name      string
	Index     string
	Columns   []string
	Algorithm string
	On        string // for rename
}

// Blueprint is a fluent table definition builder. It collects columns,
// indexes, and commands and is compiled by a SchemaGrammar into DDL.
type Blueprint struct {
	Table       string
	Columns     []*ColumnDefinition
	Commands    []BlueprintCommand
	ForeignKeys []*ForeignKeyDefinition
	Temporary   bool
	Charset     string
	Collation   string
	Engine      string
	// OrderBy is the table-level sort key. Required by ClickHouse MergeTree
	// engines; ignored by row-store grammars.
	OrderBy []string
	// PartitionBy is the table-level partition key expression. Optional for
	// ClickHouse; ignored by other grammars.
	PartitionBy string
}

// NewBlueprint creates a new Blueprint for the given table.
func NewBlueprint(table string) *Blueprint {
	return &Blueprint{Table: table}
}

// ---- Auto-Increment / Primary Key ----

// ID adds a big auto-incrementing ID column (alias for BigIncrements("id")).
func (bp *Blueprint) ID(column ...string) *ColumnDefinition {
	col := "id"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.BigIncrements(col)
}

// Increments adds an auto-incrementing unsigned integer column.
func (bp *Blueprint) Increments(column string) *ColumnDefinition {
	return bp.addColumn(column, "integer", 0, 0, 0, true, true)
}

// TinyIncrements adds a tiny auto-incrementing column.
func (bp *Blueprint) TinyIncrements(column string) *ColumnDefinition {
	return bp.addColumn(column, "tinyInteger", 0, 0, 0, true, true)
}

// SmallIncrements adds a small auto-incrementing column.
func (bp *Blueprint) SmallIncrements(column string) *ColumnDefinition {
	return bp.addColumn(column, "smallInteger", 0, 0, 0, true, true)
}

// MediumIncrements adds a medium auto-incrementing column.
func (bp *Blueprint) MediumIncrements(column string) *ColumnDefinition {
	return bp.addColumn(column, "mediumInteger", 0, 0, 0, true, true)
}

// BigIncrements adds a big auto-incrementing column.
func (bp *Blueprint) BigIncrements(column string) *ColumnDefinition {
	return bp.addColumn(column, "bigInteger", 0, 0, 0, true, true)
}

// IntegerIncrements adds an integer auto-incrementing column.
func (bp *Blueprint) IntegerIncrements(column string) *ColumnDefinition {
	return bp.Increments(column)
}

// ---- Integer Types ----

// TinyInteger adds a tiny integer column.
func (bp *Blueprint) TinyInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "tinyInteger", 0, 0, 0, false, false)
}

// SmallInteger adds a small integer column.
func (bp *Blueprint) SmallInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "smallInteger", 0, 0, 0, false, false)
}

// MediumInteger adds a medium integer column.
func (bp *Blueprint) MediumInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "mediumInteger", 0, 0, 0, false, false)
}

// Integer adds an integer column.
func (bp *Blueprint) Integer(column string) *ColumnDefinition {
	return bp.addColumn(column, "integer", 0, 0, 0, false, false)
}

// BigInteger adds a big integer column.
func (bp *Blueprint) BigInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "bigInteger", 0, 0, 0, false, false)
}

// UnsignedTinyInteger adds an unsigned tiny integer column.
func (bp *Blueprint) UnsignedTinyInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "tinyInteger", 0, 0, 0, true, false)
}

// UnsignedSmallInteger adds an unsigned small integer column.
func (bp *Blueprint) UnsignedSmallInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "smallInteger", 0, 0, 0, true, false)
}

// UnsignedMediumInteger adds an unsigned medium integer column.
func (bp *Blueprint) UnsignedMediumInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "mediumInteger", 0, 0, 0, true, false)
}

// UnsignedInteger adds an unsigned integer column.
func (bp *Blueprint) UnsignedInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "integer", 0, 0, 0, true, false)
}

// UnsignedBigInteger adds an unsigned big integer column.
func (bp *Blueprint) UnsignedBigInteger(column string) *ColumnDefinition {
	return bp.addColumn(column, "bigInteger", 0, 0, 0, true, false)
}

// ---- String Types ----

// Char adds a char column.
func (bp *Blueprint) Char(column string, length ...int) *ColumnDefinition {
	l := 255

	if len(length) > 0 {
		l = length[0]
	}

	return bp.addColumn(column, "char", l, 0, 0, false, false)
}

// String adds a varchar column.
func (bp *Blueprint) String(column string, length ...int) *ColumnDefinition {
	l := 255

	if len(length) > 0 {
		l = length[0]
	}

	return bp.addColumn(column, "string", l, 0, 0, false, false)
}

// TinyText adds a tiny text column.
func (bp *Blueprint) TinyText(column string) *ColumnDefinition {
	return bp.addColumn(column, "tinyText", 0, 0, 0, false, false)
}

// Text adds a text column.
func (bp *Blueprint) Text(column string) *ColumnDefinition {
	return bp.addColumn(column, "text", 0, 0, 0, false, false)
}

// MediumText adds a medium text column.
func (bp *Blueprint) MediumText(column string) *ColumnDefinition {
	return bp.addColumn(column, "mediumText", 0, 0, 0, false, false)
}

// LongText adds a long text column.
func (bp *Blueprint) LongText(column string) *ColumnDefinition {
	return bp.addColumn(column, "longText", 0, 0, 0, false, false)
}

// ---- Numeric Types ----

// Float adds a float column.
func (bp *Blueprint) Float(column string, precision ...int) *ColumnDefinition {
	p := 53

	if len(precision) > 0 {
		p = precision[0]
	}

	return bp.addColumn(column, "float", 0, p, 0, false, false)
}

// Double adds a double column.
func (bp *Blueprint) Double(column string) *ColumnDefinition {
	return bp.addColumn(column, "double", 0, 0, 0, false, false)
}

// Decimal adds a decimal column.
func (bp *Blueprint) Decimal(column string, args ...int) *ColumnDefinition {
	precision := 8
	scale := 2

	if len(args) > 0 {
		precision = args[0]
	}

	if len(args) > 1 {
		scale = args[1]
	}

	return bp.addColumn(column, "decimal", 0, precision, scale, false, false)
}

// ---- Boolean ----

// Boolean adds a boolean column.
func (bp *Blueprint) Boolean(column string) *ColumnDefinition {
	return bp.addColumn(column, "boolean", 0, 0, 0, false, false)
}

// ---- Date/Time Types ----

// Date adds a date column.
func (bp *Blueprint) Date(column string) *ColumnDefinition {
	return bp.addColumn(column, "date", 0, 0, 0, false, false)
}

// DateTime adds a datetime column.
func (bp *Blueprint) DateTime(column string, precision ...int) *ColumnDefinition {
	p := 0

	if len(precision) > 0 {
		p = precision[0]
	}

	return bp.addColumn(column, "dateTime", 0, p, 0, false, false)
}

// DateTimeTz adds a datetime with timezone column.
func (bp *Blueprint) DateTimeTz(column string, precision ...int) *ColumnDefinition {
	p := 0

	if len(precision) > 0 {
		p = precision[0]
	}

	return bp.addColumn(column, "dateTimeTz", 0, p, 0, false, false)
}

// Time adds a time column.
func (bp *Blueprint) Time(column string, precision ...int) *ColumnDefinition {
	p := 0

	if len(precision) > 0 {
		p = precision[0]
	}

	return bp.addColumn(column, "time", 0, p, 0, false, false)
}

// TimeTz adds a time with timezone column.
func (bp *Blueprint) TimeTz(column string, precision ...int) *ColumnDefinition {
	p := 0

	if len(precision) > 0 {
		p = precision[0]
	}

	return bp.addColumn(column, "timeTz", 0, p, 0, false, false)
}

// Timestamp adds a timestamp column.
func (bp *Blueprint) Timestamp(column string, precision ...int) *ColumnDefinition {
	p := 0

	if len(precision) > 0 {
		p = precision[0]
	}

	return bp.addColumn(column, "timestamp", 0, p, 0, false, false)
}

// TimestampTz adds a timestamp with timezone column.
func (bp *Blueprint) TimestampTz(column string, precision ...int) *ColumnDefinition {
	p := 0

	if len(precision) > 0 {
		p = precision[0]
	}

	return bp.addColumn(column, "timestampTz", 0, p, 0, false, false)
}

// Timestamps adds created_at and updated_at nullable timestamp columns.
func (bp *Blueprint) Timestamps(precision ...int) {
	bp.Timestamp("created_at", precision...).Nullable()
	bp.Timestamp("updated_at", precision...).Nullable()
}

// TimestampsTz adds created_at and updated_at nullable timestampTz columns.
func (bp *Blueprint) TimestampsTz(precision ...int) {
	bp.TimestampTz("created_at", precision...).Nullable()
	bp.TimestampTz("updated_at", precision...).Nullable()
}

// Datetimes adds created_at and updated_at nullable datetime columns.
func (bp *Blueprint) Datetimes(precision ...int) {
	bp.DateTime("created_at", precision...).Nullable()
	bp.DateTime("updated_at", precision...).Nullable()
}

// SoftDeletes adds a nullable deleted_at timestamp column.
func (bp *Blueprint) SoftDeletes(column ...string) *ColumnDefinition {
	col := "deleted_at"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.Timestamp(col).Nullable()
}

// SoftDeletesTz adds a nullable deleted_at timestampTz column.
func (bp *Blueprint) SoftDeletesTz(column ...string) *ColumnDefinition {
	col := "deleted_at"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.TimestampTz(col).Nullable()
}

// SoftDeletesDatetime adds a nullable deleted_at datetime column.
func (bp *Blueprint) SoftDeletesDatetime(column ...string) *ColumnDefinition {
	col := "deleted_at"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.DateTime(col).Nullable()
}

// Year adds a year column.
func (bp *Blueprint) Year(column string) *ColumnDefinition {
	return bp.addColumn(column, "year", 0, 0, 0, false, false)
}

// ---- Special Types ----

// Binary adds a binary column.
func (bp *Blueprint) Binary(column string, length ...int) *ColumnDefinition {
	l := 0

	if len(length) > 0 {
		l = length[0]
	}

	return bp.addColumn(column, "binary", l, 0, 0, false, false)
}

// JSON adds a json column.
func (bp *Blueprint) JSON(column string) *ColumnDefinition {
	return bp.addColumn(column, "json", 0, 0, 0, false, false)
}

// JSONB adds a jsonb column.
func (bp *Blueprint) JSONB(column string) *ColumnDefinition {
	return bp.addColumn(column, "jsonb", 0, 0, 0, false, false)
}

// Enum adds an enum column.
func (bp *Blueprint) Enum(column string, allowed []string) *ColumnDefinition {
	c := bp.addColumn(column, "enum", 0, 0, 0, false, false)
	c.Allowed = allowed

	return c
}

// Set adds a set column.
func (bp *Blueprint) Set(column string, allowed []string) *ColumnDefinition {
	c := bp.addColumn(column, "set", 0, 0, 0, false, false)
	c.Allowed = allowed

	return c
}

// UUID adds a uuid column.
func (bp *Blueprint) UUID(column ...string) *ColumnDefinition {
	col := "uuid"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.addColumn(col, "uuid", 0, 0, 0, false, false)
}

// ULID adds a ulid column (char(26)).
func (bp *Blueprint) ULID(column ...string) *ColumnDefinition {
	col := "ulid"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.Char(col, 26)
}

// IPAddress adds an IP address column.
func (bp *Blueprint) IPAddress(column ...string) *ColumnDefinition {
	col := "ip_address"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.addColumn(col, "ipAddress", 0, 0, 0, false, false)
}

// MacAddress adds a MAC address column.
func (bp *Blueprint) MacAddress(column ...string) *ColumnDefinition {
	col := "mac_address"

	if len(column) > 0 {
		col = column[0]
	}

	return bp.addColumn(col, "macAddress", 0, 0, 0, false, false)
}

// Geometry adds a geometry column.
func (bp *Blueprint) Geometry(column string, subtype ...string) *ColumnDefinition {
	return bp.addColumn(column, "geometry", 0, 0, 0, false, false)
}

// Geography adds a geography column.
func (bp *Blueprint) Geography(column string, subtype ...string) *ColumnDefinition {
	return bp.addColumn(column, "geography", 0, 0, 0, false, false)
}

// Vector adds a vector column.
func (bp *Blueprint) Vector(column string, dimensions ...int) *ColumnDefinition {
	d := 0

	if len(dimensions) > 0 {
		d = dimensions[0]
	}

	c := bp.addColumn(column, "vector", 0, 0, 0, false, false)
	c.Total = d

	return c
}

// Tsvector adds a tsvector column (PostgreSQL full-text search).
func (bp *Blueprint) Tsvector(column string) *ColumnDefinition {
	return bp.addColumn(column, "tsvector", 0, 0, 0, false, false)
}

// RememberToken adds a nullable remember_token varchar(100) column.
func (bp *Blueprint) RememberToken() *ColumnDefinition {
	return bp.String("remember_token", 100).Nullable()
}

// ---- Foreign Key Helpers ----

// ForeignID adds an unsigned big integer column for use as a foreign key.
func (bp *Blueprint) ForeignID(column string) *ColumnDefinition {
	return bp.UnsignedBigInteger(column)
}

// ForeignIDFor adds an unsigned big integer column named after the model.
func (bp *Blueprint) ForeignIDFor(column string) *ColumnDefinition {
	return bp.UnsignedBigInteger(column)
}

// ForeignUUID adds a UUID column for foreign key use.
func (bp *Blueprint) ForeignUUID(column string) *ColumnDefinition {
	return bp.UUID(column)
}

// ForeignULID adds a ULID column for foreign key use.
func (bp *Blueprint) ForeignULID(column string) *ColumnDefinition {
	return bp.ULID(column)
}

// ---- Polymorphic Columns ----

// Morphs adds the columns for a polymorphic relationship.
func (bp *Blueprint) Morphs(name string, indexName ...string) {
	bp.String(name + "_type")
	bp.UnsignedBigInteger(name + "_id")
	idx := name + "_type_" + name + "_id_index"

	if len(indexName) > 0 {
		idx = indexName[0]
	}

	bp.AddIndex([]string{name + "_type", name + "_id"}, idx)
}

// NullableMorphs adds nullable polymorphic columns.
func (bp *Blueprint) NullableMorphs(name string, indexName ...string) {
	bp.String(name + "_type").Nullable()
	bp.UnsignedBigInteger(name + "_id").Nullable()
	idx := name + "_type_" + name + "_id_index"

	if len(indexName) > 0 {
		idx = indexName[0]
	}

	bp.AddIndex([]string{name + "_type", name + "_id"}, idx)
}

// NumericMorphs adds polymorphic columns with unsigned big int ID.
func (bp *Blueprint) NumericMorphs(name string, indexName ...string) {
	bp.Morphs(name, indexName...)
}

// UUIDMorphs adds polymorphic columns with UUID ID.
func (bp *Blueprint) UUIDMorphs(name string, indexName ...string) {
	bp.String(name + "_type")
	bp.UUID(name + "_id")
	idx := name + "_type_" + name + "_id_index"

	if len(indexName) > 0 {
		idx = indexName[0]
	}

	bp.AddIndex([]string{name + "_type", name + "_id"}, idx)
}

// ULIDMorphs adds polymorphic columns with ULID ID.
func (bp *Blueprint) ULIDMorphs(name string, indexName ...string) {
	bp.String(name + "_type")
	bp.ULID(name + "_id")
	idx := name + "_type_" + name + "_id_index"

	if len(indexName) > 0 {
		idx = indexName[0]
	}

	bp.AddIndex([]string{name + "_type", name + "_id"}, idx)
}

// ---- Computed / Raw Columns ----

// Computed adds a computed/generated column.
func (bp *Blueprint) Computed(column, expression string) *ColumnDefinition {
	c := bp.addColumn(column, "computed", 0, 0, 0, false, false)
	c.VirtualAs = expression

	return c
}

// RawColumn adds a column with a raw type definition.
func (bp *Blueprint) RawColumn(column, definition string) *ColumnDefinition {
	return bp.addColumn(column, definition, 0, 0, 0, false, false)
}

// ---- Index Operations ----

// AddPrimary adds a primary key.
func (bp *Blueprint) AddPrimary(columns []string, name ...string) {
	n := bp.Table + "_" + columns[0] + "_primary"

	if len(name) > 0 {
		n = name[0]
	}

	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "primary", Columns: columns, Index: n})
}

// AddUnique adds a unique index.
func (bp *Blueprint) AddUnique(columns []string, name ...string) {
	n := bp.Table + "_" + columns[0] + "_unique"

	if len(name) > 0 {
		n = name[0]
	}

	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "unique", Columns: columns, Index: n})
}

// AddIndex adds an index.
func (bp *Blueprint) AddIndex(columns []string, name ...string) {
	n := bp.Table + "_" + columns[0] + "_index"

	if len(name) > 0 {
		n = name[0]
	}

	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "index", Columns: columns, Index: n})
}

// AddFulltext adds a fulltext index.
func (bp *Blueprint) AddFulltext(columns []string, name ...string) {
	n := bp.Table + "_" + columns[0] + "_fulltext"

	if len(name) > 0 {
		n = name[0]
	}

	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "fulltext", Columns: columns, Index: n})
}

// Foreign adds a foreign key constraint.
func (bp *Blueprint) Foreign(columns ...string) *ForeignKeyDefinition {
	fk := &ForeignKeyDefinition{Columns: columns}
	bp.ForeignKeys = append(bp.ForeignKeys, fk)

	return fk
}

// DropColumn adds a command to drop a column.
func (bp *Blueprint) DropColumn(columns ...string) {
	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "dropColumn", Columns: columns})
}

// RenameColumn adds a command to rename a column.
func (bp *Blueprint) RenameColumn(from, to string) {
	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "renameColumn", Columns: []string{from, to}})
}

// DropPrimary drops a primary key.
func (bp *Blueprint) DropPrimary(name ...string) {
	n := bp.Table + "_primary"

	if len(name) > 0 {
		n = name[0]
	}

	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "dropPrimary", Index: n})
}

// DropUnique drops a unique index.
func (bp *Blueprint) DropUnique(name string) {
	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "dropUnique", Index: name})
}

// DropIndex drops an index.
func (bp *Blueprint) DropIndex(name string) {
	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "dropIndex", Index: name})
}

// DropForeign drops a foreign key constraint.
func (bp *Blueprint) DropForeign(name string) {
	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "dropForeign", Index: name})
}

// DropFulltext drops a fulltext index.
func (bp *Blueprint) DropFulltext(name string) {
	bp.Commands = append(bp.Commands, BlueprintCommand{Name: "dropFulltext", Index: name})
}

// DropTimestamps drops the created_at and updated_at columns.
func (bp *Blueprint) DropTimestamps() {
	bp.DropColumn("created_at", "updated_at")
}

// DropSoftDeletes drops the deleted_at column.
func (bp *Blueprint) DropSoftDeletes(column ...string) {
	col := "deleted_at"

	if len(column) > 0 {
		col = column[0]
	}

	bp.DropColumn(col)
}

// DropRememberToken drops the remember_token column.
func (bp *Blueprint) DropRememberToken() {
	bp.DropColumn("remember_token")
}

// SetTemporary marks the table as temporary.
func (bp *Blueprint) SetTemporary() {
	bp.Temporary = true
}

// SetEngine sets the table storage engine (MySQL).
func (bp *Blueprint) SetEngine(engine string) {
	bp.Engine = engine
}

// SetCharset sets the table default character set.
func (bp *Blueprint) SetCharset(charset string) {
	bp.Charset = charset
}

// SetCollation sets the table default collation.
func (bp *Blueprint) SetCollation(collation string) {
	bp.Collation = collation
}

// SetOrderBy sets the table-level sort key columns. Used by ClickHouse
// MergeTree engines to populate the required ORDER BY clause.
func (bp *Blueprint) SetOrderBy(columns ...string) {
	bp.OrderBy = columns
}

// SetPartitionBy sets the table-level partition expression. ClickHouse-only;
// ignored by grammars that do not support partitioning.
func (bp *Blueprint) SetPartitionBy(expression string) {
	bp.PartitionBy = expression
}

// addColumn is the internal helper to add a column definition.
func (bp *Blueprint) addColumn(name, colType string, length, precision, scale int, unsigned, autoIncrement bool) *ColumnDefinition {
	col := &ColumnDefinition{
		Name:          name,
		Type:          colType,
		Length:        length,
		Precision:     precision,
		Scale:         scale,
		Unsigned:      unsigned,
		AutoIncrement: autoIncrement,
	}
	bp.Columns = append(bp.Columns, col)

	return col
}

// GetAddedColumns returns columns that are not being changed.
func (bp *Blueprint) GetAddedColumns() []*ColumnDefinition {
	var added []*ColumnDefinition

	for _, col := range bp.Columns {
		if !col.IsChange {
			added = append(added, col)
		}
	}

	return added
}

// GetChangedColumns returns columns that are being modified.
func (bp *Blueprint) GetChangedColumns() []*ColumnDefinition {
	var changed []*ColumnDefinition

	for _, col := range bp.Columns {
		if col.IsChange {
			changed = append(changed, col)
		}
	}

	return changed
}
