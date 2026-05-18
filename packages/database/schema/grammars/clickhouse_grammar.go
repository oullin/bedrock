package grammars

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/database/schema"
)

// ClickHouseGrammar compiles schema blueprints into ClickHouse DDL.
//
// Every CREATE TABLE statement must declare a table engine. The grammar
// defaults to MergeTree with the primary-key column (or, if absent, the first
// declared column) as the ORDER BY key. Apps override these with
// Blueprint.SetEngine, SetOrderBy, and SetPartitionBy.
//
// Foreign keys, per-row mutations without WHERE, and a few other row-store
// constructs are not supported and produce empty DDL (callers should consult
// the driver's sentinel errors before relying on them).
type ClickHouseGrammar struct {
	tablePrefix string
}

var _ schema.Grammar = (*ClickHouseGrammar)(nil)

// NewClickHouseGrammar creates a new ClickHouse schema grammar.
func NewClickHouseGrammar() *ClickHouseGrammar { return &ClickHouseGrammar{} }

func (g *ClickHouseGrammar) CompileCreate(bp *schema.Blueprint) []string {
	columns := g.getColumns(bp)
	temporary := ""

	if bp.Temporary {
		temporary = "temporary "
	}

	sql := fmt.Sprintf("create %stable %s (%s)", temporary, g.wrapTable(bp.Table), strings.Join(columns, ", "))

	engine := bp.Engine

	if engine == "" {
		engine = "MergeTree"
	}

	sql += " engine = " + engine + "()"

	if pb := strings.TrimSpace(bp.PartitionBy); pb != "" {
		sql += " partition by " + pb
	}

	orderBy := bp.OrderBy

	if len(orderBy) == 0 {
		orderBy = []string{g.defaultOrderColumn(bp)}
	}

	wrappedOrderBy := make([]string, len(orderBy))

	for i, c := range orderBy {
		wrappedOrderBy[i] = g.wrap(c)
	}

	sql += " order by (" + strings.Join(wrappedOrderBy, ", ") + ")"

	return []string{sql}
}

func (g *ClickHouseGrammar) CompileAdd(bp *schema.Blueprint) []string {
	var statements []string

	for _, col := range bp.GetAddedColumns() {
		statements = append(statements,
			fmt.Sprintf("alter table %s add column %s", g.wrapTable(bp.Table), g.compileColumn(col)))
	}

	return statements
}

func (g *ClickHouseGrammar) CompileChange(bp *schema.Blueprint) []string {
	var statements []string

	for _, col := range bp.GetChangedColumns() {
		statements = append(statements,
			fmt.Sprintf("alter table %s modify column %s %s", g.wrapTable(bp.Table), g.wrap(col.Name), g.columnType(col)))
	}

	return statements
}

func (g *ClickHouseGrammar) CompileDrop(table string) string {
	return "drop table " + g.wrapTable(table)
}

func (g *ClickHouseGrammar) CompileDropIfExists(table string) string {
	return "drop table if exists " + g.wrapTable(table)
}

func (g *ClickHouseGrammar) CompileRename(from, to string) string {
	return "rename table " + g.wrapTable(from) + " to " + g.wrapTable(to)
}

func (g *ClickHouseGrammar) CompileDropColumn(bp *schema.Blueprint, columns []string) string {
	cols := make([]string, len(columns))

	for i, c := range columns {
		cols[i] = "drop column " + g.wrap(c)
	}

	return "alter table " + g.wrapTable(bp.Table) + " " + strings.Join(cols, ", ")
}

func (g *ClickHouseGrammar) CompileRenameColumn(bp *schema.Blueprint, from, to string) string {
	return "alter table " + g.wrapTable(bp.Table) + " rename column " + g.wrap(from) + " to " + g.wrap(to)
}

func (g *ClickHouseGrammar) CompileCreateIndex(bp *schema.Blueprint, cmd schema.BlueprintCommand) string {
	// ClickHouse primary keys are declared via the MergeTree ORDER BY clause,
	// not via ALTER TABLE ... ADD PRIMARY KEY. Skip primary commands at the
	// command stage; they were already wired into the engine clause.
	if cmd.Name == "primary" {
		return ""
	}

	cols := make([]string, len(cmd.Columns))

	for i, c := range cmd.Columns {
		cols[i] = g.wrap(c)
	}

	switch cmd.Name {
	case "unique":
		// ClickHouse has no UNIQUE constraint; closest is a data-skipping
		// index. Emit a minmax skip index as a useful approximation.
		return fmt.Sprintf("alter table %s add index %s (%s) type minmax granularity 1",
			g.wrapTable(bp.Table), g.wrap(cmd.Index), strings.Join(cols, ", "))
	case "index":
		return fmt.Sprintf("alter table %s add index %s (%s) type minmax granularity 1",
			g.wrapTable(bp.Table), g.wrap(cmd.Index), strings.Join(cols, ", "))
	case "fulltext":
		return fmt.Sprintf("alter table %s add index %s (%s) type tokenbf_v1(8192, 3, 0) granularity 1",
			g.wrapTable(bp.Table), g.wrap(cmd.Index), strings.Join(cols, ", "))
	}

	return ""
}

func (g *ClickHouseGrammar) CompileDropIndex(bp *schema.Blueprint, name string) string {
	return "alter table " + g.wrapTable(bp.Table) + " drop index " + g.wrap(name)
}

// CompileCreateForeignKey returns an empty string: ClickHouse does not enforce
// foreign keys. The driver's ErrForeignKeysNotSupported sentinel documents the
// limitation; callers should check it before relying on referential integrity.
func (g *ClickHouseGrammar) CompileCreateForeignKey(_ *schema.Blueprint, _ *schema.ForeignKeyDefinition) string {
	return ""
}

func (g *ClickHouseGrammar) CompileDropForeignKey(_ *schema.Blueprint, _ string) string {
	return ""
}

func (g *ClickHouseGrammar) CompileTableExists() string {
	return "select name from system.tables where database = currentDatabase() and name = ?"
}

func (g *ClickHouseGrammar) CompileColumnListing(table string) string {
	return "select name from system.columns where database = currentDatabase() and table = '" + table + "'"
}

func (g *ClickHouseGrammar) CompileEnableForeignKeyConstraints() string  { return "" }
func (g *ClickHouseGrammar) CompileDisableForeignKeyConstraints() string { return "" }

func (g *ClickHouseGrammar) getColumns(bp *schema.Blueprint) []string {
	var cols []string

	for _, col := range bp.GetAddedColumns() {
		cols = append(cols, g.compileColumn(col))
	}

	return cols
}

func (g *ClickHouseGrammar) compileColumn(col *schema.ColumnDefinition) string {
	sql := g.wrap(col.Name) + " " + g.columnType(col)

	if col.HasDefault {
		sql += " default " + g.getDefaultValue(col.DefaultValue)
	}

	if col.CommentText != "" {
		sql += " comment '" + strings.ReplaceAll(col.CommentText, "'", "\\'") + "'"
	}

	return sql
}

// columnType returns the ClickHouse type for a column, wrapping nullable
// columns in Nullable(...). ClickHouse columns are NOT NULL by default — the
// inverse of standard SQL — so we wrap rather than appending "not null".
func (g *ClickHouseGrammar) columnType(col *schema.ColumnDefinition) string {
	base := g.getType(col)

	if col.IsNullable {
		return "Nullable(" + base + ")"
	}

	return base
}

func (g *ClickHouseGrammar) getType(col *schema.ColumnDefinition) string {
	switch col.Type {
	case "bigInteger":
		if col.Unsigned {
			return "UInt64"
		}

		return "Int64"
	case "integer", "mediumInteger":
		if col.Unsigned {
			return "UInt32"
		}

		return "Int32"
	case "smallInteger":
		if col.Unsigned {
			return "UInt16"
		}

		return "Int16"
	case "tinyInteger":
		if col.Unsigned {
			return "UInt8"
		}

		return "Int8"
	case "string", "char":
		if col.Length > 0 {
			return fmt.Sprintf("FixedString(%d)", col.Length)
		}

		return "String"
	case "text", "tinyText", "mediumText", "longText":
		return "String"
	case "float":
		return "Float32"
	case "double":
		return "Float64"
	case "decimal":
		precision := col.Precision

		if precision <= 0 {
			precision = 18
		}

		scale := col.Scale

		return fmt.Sprintf("Decimal(%d, %d)", precision, scale)
	case "boolean":
		return "UInt8"
	case "date":
		return "Date"
	case "dateTime", "dateTimeTz", "timestamp", "timestampTz":
		if col.Precision > 0 {
			return fmt.Sprintf("DateTime64(%d)", col.Precision)
		}

		return "DateTime64(3)"
	case "time", "timeTz":
		return "String"
	case "year":
		return "UInt16"
	case "binary":
		return "String"
	case "json", "jsonb":
		// ClickHouse's JSON type is experimental; map to String for stability.
		// Apps needing structured access should model with Nested or Tuple.
		return "String"
	case "uuid":
		return "UUID"
	case "ipAddress":
		return "IPv6"
	case "macAddress":
		return "String"
	case "enum":
		// Build Enum8('a' = 1, 'b' = 2, ...).
		parts := make([]string, len(col.Allowed))

		for i, v := range col.Allowed {
			escaped := strings.ReplaceAll(v, "'", "\\'")
			parts[i] = fmt.Sprintf("'%s' = %d", escaped, i+1)
		}

		return "Enum8(" + strings.Join(parts, ", ") + ")"
	case "vector":
		if col.Total > 0 {
			return fmt.Sprintf("Array(Float32) /* dimensions=%d */", col.Total)
		}

		return "Array(Float32)"
	default:
		return col.Type
	}
}

func (g *ClickHouseGrammar) getDefaultValue(value any) string {
	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "\\'") + "'"
	case bool:
		if v {
			return "1"
		}

		return "0"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func (g *ClickHouseGrammar) wrap(value string) string {
	return "`" + strings.ReplaceAll(value, "`", "``") + "`"
}

func (g *ClickHouseGrammar) wrapTable(table string) string {
	return g.wrap(g.tablePrefix + table)
}

// defaultOrderColumn picks a sensible ORDER BY column when the blueprint does
// not declare one explicitly: the first auto-increment column if any, else
// the first declared column.
func (g *ClickHouseGrammar) defaultOrderColumn(bp *schema.Blueprint) string {
	for _, col := range bp.Columns {
		if col.AutoIncrement {
			return col.Name
		}
	}

	if len(bp.Columns) > 0 {
		return bp.Columns[0].Name
	}

	return "tuple()"
}
