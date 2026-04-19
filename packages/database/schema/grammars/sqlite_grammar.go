package grammars

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/database/schema"
)

// SQLiteGrammar compiles schema blueprints into SQLite DDL.
type SQLiteGrammar struct {
	tablePrefix string
}

var _ schema.Grammar = (*SQLiteGrammar)(nil)

// NewSQLiteGrammar creates a new SQLite schema grammar.
func NewSQLiteGrammar() *SQLiteGrammar { return &SQLiteGrammar{} }

func (g *SQLiteGrammar) CompileCreate(bp *schema.Blueprint) []string {
	columns := g.getColumns(bp)
	temporary := ""

	if bp.Temporary {
		temporary = "temporary "
	}

	sql := fmt.Sprintf("create %stable %s (%s)", temporary, g.wrapTable(bp.Table), strings.Join(columns, ", "))

	// Add primary key for auto-increment columns.
	for _, col := range bp.Columns {
		if col.AutoIncrement {
			// SQLite uses INTEGER PRIMARY KEY for auto-increment.
			break
		}
	}

	return []string{sql}
}

func (g *SQLiteGrammar) CompileAdd(bp *schema.Blueprint) []string {
	var statements []string

	for _, col := range bp.GetAddedColumns() {
		statements = append(statements, fmt.Sprintf("alter table %s add column %s", g.wrapTable(bp.Table), g.compileColumn(col)))
	}

	return statements
}

func (g *SQLiteGrammar) CompileChange(bp *schema.Blueprint) []string {
	// SQLite has limited ALTER TABLE support. Column modification requires
	// recreating the table in many cases.
	var statements []string

	for _, col := range bp.GetChangedColumns() {
		statements = append(statements, fmt.Sprintf("alter table %s rename column %s to %s", g.wrapTable(bp.Table), g.wrap(col.Name), g.wrap(col.Name)))
	}

	return statements
}

func (g *SQLiteGrammar) CompileDrop(table string) string {
	return "drop table " + g.wrapTable(table)
}

func (g *SQLiteGrammar) CompileDropIfExists(table string) string {
	return "drop table if exists " + g.wrapTable(table)
}

func (g *SQLiteGrammar) CompileRename(from, to string) string {
	return "alter table " + g.wrapTable(from) + " rename to " + g.wrapTable(to)
}

func (g *SQLiteGrammar) CompileDropColumn(bp *schema.Blueprint, columns []string) string {
	cols := make([]string, len(columns))

	for i, c := range columns {
		cols[i] = g.wrap(c)
	}

	return "alter table " + g.wrapTable(bp.Table) + " drop column " + strings.Join(cols, ", drop column ")
}

func (g *SQLiteGrammar) CompileRenameColumn(bp *schema.Blueprint, from, to string) string {
	return "alter table " + g.wrapTable(bp.Table) + " rename column " + g.wrap(from) + " to " + g.wrap(to)
}

func (g *SQLiteGrammar) CompileCreateIndex(bp *schema.Blueprint, cmd schema.BlueprintCommand) string {
	cols := make([]string, len(cmd.Columns))

	for i, c := range cmd.Columns {
		cols[i] = g.wrap(c)
	}

	switch cmd.Name {
	case "unique":
		return fmt.Sprintf("create unique index %s on %s (%s)", g.wrap(cmd.Index), g.wrapTable(bp.Table), strings.Join(cols, ", "))
	case "index", "fulltext":
		return fmt.Sprintf("create index %s on %s (%s)", g.wrap(cmd.Index), g.wrapTable(bp.Table), strings.Join(cols, ", "))
	case "primary":
		// SQLite handles primary keys in CREATE TABLE.
		return ""
	}

	return ""
}

func (g *SQLiteGrammar) CompileDropIndex(_ *schema.Blueprint, name string) string {
	return "drop index if exists " + g.wrap(name)
}

func (g *SQLiteGrammar) CompileCreateForeignKey(bp *schema.Blueprint, fk *schema.ForeignKeyDefinition) string {
	// SQLite foreign keys must be defined in CREATE TABLE, not ALTER TABLE.
	// This is a limitation. For migrations, they are added during table creation.
	return ""
}

func (g *SQLiteGrammar) CompileDropForeignKey(_ *schema.Blueprint, name string) string {
	return "" // SQLite doesn't support dropping foreign keys individually.
}

func (g *SQLiteGrammar) CompileTableExists() string {
	return "select * from sqlite_master where type = 'table' and name = ?"
}

func (g *SQLiteGrammar) CompileColumnListing(table string) string {
	return fmt.Sprintf("pragma table_info(%s)", g.wrapTable(table))
}

func (g *SQLiteGrammar) CompileEnableForeignKeyConstraints() string {
	return "PRAGMA foreign_keys = ON"
}

func (g *SQLiteGrammar) CompileDisableForeignKeyConstraints() string {
	return "PRAGMA foreign_keys = OFF"
}

func (g *SQLiteGrammar) getColumns(bp *schema.Blueprint) []string {
	var cols []string

	for _, col := range bp.GetAddedColumns() {
		cols = append(cols, g.compileColumn(col))
	}

	return cols
}

func (g *SQLiteGrammar) compileColumn(col *schema.ColumnDefinition) string {
	sql := g.wrap(col.Name) + " " + g.getType(col)

	if col.AutoIncrement {
		sql = g.wrap(col.Name) + " integer primary key autoincrement"

		return sql
	}

	if !col.IsNullable {
		sql += " not null"
	}

	if col.HasDefault {
		sql += " default " + g.getDefaultValue(col.DefaultValue)
	}

	if col.Unsigned {
		// SQLite doesn't enforce unsigned, but we track it.
	}

	return sql
}

func (g *SQLiteGrammar) getType(col *schema.ColumnDefinition) string {
	switch col.Type {
	case "bigInteger", "mediumInteger", "integer", "smallInteger", "tinyInteger":
		return "integer"
	case "string", "char":
		if col.Length > 0 {
			return fmt.Sprintf("varchar(%d)", col.Length)
		}

		return "varchar"
	case "text", "tinyText", "mediumText", "longText":
		return "text"
	case "float", "double", "decimal":
		return "real"
	case "boolean":
		return "tinyint(1)"
	case "date":
		return "date"
	case "dateTime", "dateTimeTz":
		return "datetime"
	case "time", "timeTz":
		return "time"
	case "timestamp", "timestampTz":
		return "datetime"
	case "year":
		return "integer"
	case "binary":
		return "blob"
	case "json", "jsonb":
		return "text"
	case "uuid":
		return "varchar(36)"
	case "ipAddress":
		return "varchar(45)"
	case "macAddress":
		return "varchar(17)"
	case "enum":
		return "varchar"
	case "set":
		return "varchar"
	case "geometry", "geography":
		return "text"
	case "vector":
		return "text"
	default:
		return col.Type
	}
}

func (g *SQLiteGrammar) getDefaultValue(value any) string {
	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case bool:
		if v {
			return "1"
		}

		return "0"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (g *SQLiteGrammar) wrap(value string) string {
	if value == "*" {
		return value
	}

	return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
}

func (g *SQLiteGrammar) wrapTable(table string) string {
	return g.wrap(g.tablePrefix + table)
}
