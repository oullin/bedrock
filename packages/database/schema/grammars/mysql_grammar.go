package grammars

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/database/schema"
)

// MySQLGrammar compiles schema blueprints into MySQL DDL.
type MySQLGrammar struct {
	tablePrefix string
}

var _ schema.Grammar = (*MySQLGrammar)(nil)

// NewMySQLGrammar creates a new MySQL schema grammar.
func NewMySQLGrammar() *MySQLGrammar { return &MySQLGrammar{} }

func (g *MySQLGrammar) CompileCreate(bp *schema.Blueprint) []string {
	columns := g.getColumns(bp)
	temporary := ""

	if bp.Temporary {
		temporary = "temporary "
	}

	sql := fmt.Sprintf("create %stable %s (%s)", temporary, g.wrapTable(bp.Table), strings.Join(columns, ", "))

	if bp.Engine != "" {
		sql += " engine = " + bp.Engine
	}

	if bp.Charset != "" {
		sql += " default character set " + bp.Charset
	}

	if bp.Collation != "" {
		sql += " collate " + bp.Collation
	}

	return []string{sql}
}

func (g *MySQLGrammar) CompileAdd(bp *schema.Blueprint) []string {
	var parts []string

	for _, col := range bp.GetAddedColumns() {
		parts = append(parts, "add "+g.compileColumn(col))
	}

	if len(parts) == 0 {
		return nil
	}

	return []string{fmt.Sprintf("alter table %s %s", g.wrapTable(bp.Table), strings.Join(parts, ", "))}
}

func (g *MySQLGrammar) CompileChange(bp *schema.Blueprint) []string {
	var parts []string

	for _, col := range bp.GetChangedColumns() {
		parts = append(parts, fmt.Sprintf("modify %s", g.compileColumn(col)))
	}

	if len(parts) == 0 {
		return nil
	}

	return []string{fmt.Sprintf("alter table %s %s", g.wrapTable(bp.Table), strings.Join(parts, ", "))}
}

func (g *MySQLGrammar) CompileDrop(table string) string {
	return "drop table " + g.wrapTable(table)
}

func (g *MySQLGrammar) CompileDropIfExists(table string) string {
	return "drop table if exists " + g.wrapTable(table)
}

func (g *MySQLGrammar) CompileRename(from, to string) string {
	return "rename table " + g.wrapTable(from) + " to " + g.wrapTable(to)
}

func (g *MySQLGrammar) CompileDropColumn(bp *schema.Blueprint, columns []string) string {
	cols := make([]string, len(columns))

	for i, c := range columns {
		cols[i] = "drop " + g.wrap(c)
	}

	return "alter table " + g.wrapTable(bp.Table) + " " + strings.Join(cols, ", ")
}

func (g *MySQLGrammar) CompileRenameColumn(bp *schema.Blueprint, from, to string) string {
	return "alter table " + g.wrapTable(bp.Table) + " rename column " + g.wrap(from) + " to " + g.wrap(to)
}

func (g *MySQLGrammar) CompileCreateIndex(bp *schema.Blueprint, cmd schema.BlueprintCommand) string {
	cols := make([]string, len(cmd.Columns))

	for i, c := range cmd.Columns {
		cols[i] = g.wrap(c)
	}

	switch cmd.Name {
	case "primary":
		return fmt.Sprintf("alter table %s add primary key (%s)", g.wrapTable(bp.Table), strings.Join(cols, ", "))
	case "unique":
		return fmt.Sprintf("alter table %s add unique %s(%s)", g.wrapTable(bp.Table), g.wrap(cmd.Index), strings.Join(cols, ", "))
	case "index":
		return fmt.Sprintf("alter table %s add index %s(%s)", g.wrapTable(bp.Table), g.wrap(cmd.Index), strings.Join(cols, ", "))
	case "fulltext":
		return fmt.Sprintf("alter table %s add fulltext %s(%s)", g.wrapTable(bp.Table), g.wrap(cmd.Index), strings.Join(cols, ", "))
	}

	return ""
}

func (g *MySQLGrammar) CompileDropIndex(bp *schema.Blueprint, name string) string {
	return "alter table " + g.wrapTable(bp.Table) + " drop index " + g.wrap(name)
}

func (g *MySQLGrammar) CompileCreateForeignKey(bp *schema.Blueprint, fk *schema.ForeignKeyDefinition) string {
	cols := make([]string, len(fk.Columns))

	for i, c := range fk.Columns {
		cols[i] = g.wrap(c)
	}

	refCols := make([]string, len(fk.RefColumns))

	for i, c := range fk.RefColumns {
		refCols[i] = g.wrap(c)
	}

	sql := fmt.Sprintf("alter table %s add constraint %s foreign key (%s) references %s (%s)",
		g.wrapTable(bp.Table), g.wrap(fk.IndexName),
		strings.Join(cols, ", "), g.wrapTable(fk.RefTable), strings.Join(refCols, ", "))

	if fk.OnDeleteAction != "" {
		sql += " on delete " + fk.OnDeleteAction
	}

	if fk.OnUpdateAction != "" {
		sql += " on update " + fk.OnUpdateAction
	}

	return sql
}

func (g *MySQLGrammar) CompileDropForeignKey(bp *schema.Blueprint, name string) string {
	return "alter table " + g.wrapTable(bp.Table) + " drop foreign key " + g.wrap(name)
}

func (g *MySQLGrammar) CompileTableExists() string {
	return "select * from information_schema.tables where table_schema = database() and table_name = ? and table_type = 'BASE TABLE'"
}

func (g *MySQLGrammar) CompileColumnListing(table string) string {
	return "select column_name as `column_name` from information_schema.columns where table_schema = database() and table_name = '" + table + "'"
}

func (g *MySQLGrammar) CompileEnableForeignKeyConstraints() string {
	return "SET FOREIGN_KEY_CHECKS=1"
}

func (g *MySQLGrammar) CompileDisableForeignKeyConstraints() string {
	return "SET FOREIGN_KEY_CHECKS=0"
}

func (g *MySQLGrammar) getColumns(bp *schema.Blueprint) []string {
	var cols []string

	for _, col := range bp.GetAddedColumns() {
		cols = append(cols, g.compileColumn(col))
	}

	return cols
}

func (g *MySQLGrammar) compileColumn(col *schema.ColumnDefinition) string {
	sql := g.wrap(col.Name) + " " + g.getType(col)

	if col.Unsigned {
		sql += " unsigned"
	}

	if col.AutoIncrement {
		sql += " auto_increment primary key"
	}

	if col.CharsetName != "" {
		sql += " character set " + col.CharsetName
	}

	if col.Collation != "" {
		sql += " collate " + col.Collation
	}

	if !col.IsNullable && !col.AutoIncrement {
		sql += " not null"
	}

	if col.IsNullable {
		sql += " null"
	}

	if col.HasDefault {
		sql += " default " + g.getDefaultValue(col.DefaultValue)
	}

	if col.CommentText != "" {
		sql += " comment '" + strings.ReplaceAll(col.CommentText, "'", "\\'") + "'"
	}

	if col.AfterColumn != "" {
		sql += " after " + g.wrap(col.AfterColumn)
	}

	if col.IsFirst {
		sql += " first"
	}

	if col.VirtualAs != "" {
		sql += " as (" + col.VirtualAs + ")"
	}

	if col.StoredAs != "" {
		sql += " as (" + col.StoredAs + ") stored"
	}

	return sql
}

func (g *MySQLGrammar) getType(col *schema.ColumnDefinition) string {
	switch col.Type {
	case "bigInteger":
		return "bigint"
	case "mediumInteger":
		return "mediumint"
	case "integer":
		return "int"
	case "smallInteger":
		return "smallint"
	case "tinyInteger":
		return "tinyint"
	case "string":
		return fmt.Sprintf("varchar(%d)", col.Length)
	case "char":
		return fmt.Sprintf("char(%d)", col.Length)
	case "text":
		return "text"
	case "tinyText":
		return "tinytext"
	case "mediumText":
		return "mediumtext"
	case "longText":
		return "longtext"
	case "float":
		return "double"
	case "double":
		return "double"
	case "decimal":
		return fmt.Sprintf("decimal(%d, %d)", col.Precision, col.Scale)
	case "boolean":
		return "tinyint(1)"
	case "date":
		return "date"
	case "dateTime":
		if col.Precision > 0 {
			return fmt.Sprintf("datetime(%d)", col.Precision)
		}

		return "datetime"
	case "dateTimeTz":
		if col.Precision > 0 {
			return fmt.Sprintf("datetime(%d)", col.Precision)
		}

		return "datetime"
	case "time":
		if col.Precision > 0 {
			return fmt.Sprintf("time(%d)", col.Precision)
		}

		return "time"
	case "timeTz":
		if col.Precision > 0 {
			return fmt.Sprintf("time(%d)", col.Precision)
		}

		return "time"
	case "timestamp":
		if col.Precision > 0 {
			return fmt.Sprintf("timestamp(%d)", col.Precision)
		}

		return "timestamp"
	case "timestampTz":
		if col.Precision > 0 {
			return fmt.Sprintf("timestamp(%d)", col.Precision)
		}

		return "timestamp"
	case "year":
		return "year"
	case "binary":
		if col.Length > 0 {
			return fmt.Sprintf("binary(%d)", col.Length)
		}

		return "blob"
	case "json":
		return "json"
	case "jsonb":
		return "json"
	case "uuid":
		return "char(36)"
	case "ipAddress":
		return "varchar(45)"
	case "macAddress":
		return "varchar(17)"
	case "enum":
		return "enum('" + strings.Join(col.Allowed, "', '") + "')"
	case "set":
		return "set('" + strings.Join(col.Allowed, "', '") + "')"
	case "geometry":
		return "geometry"
	case "geography":
		return "geometry"
	case "vector":
		if col.Total > 0 {
			return fmt.Sprintf("vector(%d)", col.Total)
		}

		return "vector"
	default:
		return col.Type
	}
}

func (g *MySQLGrammar) getDefaultValue(value any) string {
	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "\\'") + "'"
	case bool:
		if v {
			return "'1'"
		}

		return "'0'"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func (g *MySQLGrammar) wrap(value string) string {
	return "`" + strings.ReplaceAll(value, "`", "``") + "`"
}

func (g *MySQLGrammar) wrapTable(table string) string {
	return g.wrap(g.tablePrefix + table)
}
