package grammars

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/database/schema"
)

// PostgresGrammar compiles schema blueprints into PostgreSQL DDL.
type PostgresGrammar struct {
	tablePrefix string
}

var _ schema.Grammar = (*PostgresGrammar)(nil)

// NewPostgresGrammar creates a new PostgreSQL schema grammar.
func NewPostgresGrammar() *PostgresGrammar { return &PostgresGrammar{} }

func (g *PostgresGrammar) CompileCreate(bp *schema.Blueprint) []string {
	columns := g.getColumns(bp)
	temporary := ""

	if bp.Temporary {
		temporary = "temporary "
	}

	sql := fmt.Sprintf("create %stable %s (%s)", temporary, g.wrapTable(bp.Table), strings.Join(columns, ", "))

	return []string{sql}
}

func (g *PostgresGrammar) CompileAdd(bp *schema.Blueprint) []string {
	var statements []string

	for _, col := range bp.GetAddedColumns() {
		statements = append(statements, fmt.Sprintf("alter table %s add column %s", g.wrapTable(bp.Table), g.compileColumn(col)))
	}

	return statements
}

func (g *PostgresGrammar) CompileChange(bp *schema.Blueprint) []string {
	var statements []string

	for _, col := range bp.GetChangedColumns() {
		base := "alter table " + g.wrapTable(bp.Table)
		statements = append(statements, fmt.Sprintf("%s alter column %s type %s", base, g.wrap(col.Name), g.getType(col)))

		if col.IsNullable {
			statements = append(statements, fmt.Sprintf("%s alter column %s drop not null", base, g.wrap(col.Name)))
		} else {
			statements = append(statements, fmt.Sprintf("%s alter column %s set not null", base, g.wrap(col.Name)))
		}

		if col.HasDefault {
			statements = append(statements, fmt.Sprintf("%s alter column %s set default %s", base, g.wrap(col.Name), g.getDefaultValue(col.DefaultValue)))
		}
	}

	return statements
}

func (g *PostgresGrammar) CompileDrop(table string) string {
	return "drop table " + g.wrapTable(table)
}

func (g *PostgresGrammar) CompileDropIfExists(table string) string {
	return "drop table if exists " + g.wrapTable(table)
}

func (g *PostgresGrammar) CompileRename(from, to string) string {
	return "alter table " + g.wrapTable(from) + " rename to " + g.wrap(to)
}

func (g *PostgresGrammar) CompileDropColumn(bp *schema.Blueprint, columns []string) string {
	cols := make([]string, len(columns))

	for i, c := range columns {
		cols[i] = "drop column " + g.wrap(c)
	}

	return "alter table " + g.wrapTable(bp.Table) + " " + strings.Join(cols, ", ")
}

func (g *PostgresGrammar) CompileRenameColumn(bp *schema.Blueprint, from, to string) string {
	return "alter table " + g.wrapTable(bp.Table) + " rename column " + g.wrap(from) + " to " + g.wrap(to)
}

func (g *PostgresGrammar) CompileCreateIndex(bp *schema.Blueprint, cmd schema.BlueprintCommand) string {
	cols := make([]string, len(cmd.Columns))

	for i, c := range cmd.Columns {
		cols[i] = g.wrap(c)
	}

	switch cmd.Name {
	case "primary":
		return fmt.Sprintf("alter table %s add primary key (%s)", g.wrapTable(bp.Table), strings.Join(cols, ", "))
	case "unique":
		return fmt.Sprintf("create unique index %s on %s (%s)", g.wrap(cmd.Index), g.wrapTable(bp.Table), strings.Join(cols, ", "))
	case "index":
		return fmt.Sprintf("create index %s on %s (%s)", g.wrap(cmd.Index), g.wrapTable(bp.Table), strings.Join(cols, ", "))
	case "fulltext":
		return fmt.Sprintf("create index %s on %s using gin(to_tsvector('english', %s))", g.wrap(cmd.Index), g.wrapTable(bp.Table), strings.Join(cols, " || ' ' || "))
	}

	return ""
}

func (g *PostgresGrammar) CompileDropIndex(_ *schema.Blueprint, name string) string {
	return "drop index if exists " + g.wrap(name)
}

func (g *PostgresGrammar) CompileCreateForeignKey(bp *schema.Blueprint, fk *schema.ForeignKeyDefinition) (string, error) {
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

	return sql, nil
}

func (g *PostgresGrammar) CompileDropForeignKey(bp *schema.Blueprint, name string) string {
	return "alter table " + g.wrapTable(bp.Table) + " drop constraint " + g.wrap(name)
}

func (g *PostgresGrammar) CompileTableExists() string {
	return "select * from information_schema.tables where table_catalog = current_database() and table_schema = current_schema() and table_name = ? and table_type = 'BASE TABLE'"
}

func (g *PostgresGrammar) CompileColumnListing(table string) string {
	return "select column_name from information_schema.columns where table_catalog = current_database() and table_schema = current_schema() and table_name = '" + table + "'"
}

func (g *PostgresGrammar) CompileEnableForeignKeyConstraints() string {
	return "SET CONSTRAINTS ALL IMMEDIATE"
}

func (g *PostgresGrammar) CompileDisableForeignKeyConstraints() string {
	return "SET CONSTRAINTS ALL DEFERRED"
}

func (g *PostgresGrammar) getColumns(bp *schema.Blueprint) []string {
	var cols []string

	for _, col := range bp.GetAddedColumns() {
		cols = append(cols, g.compileColumn(col))
	}

	return cols
}

func (g *PostgresGrammar) compileColumn(col *schema.ColumnDefinition) string {
	sql := g.wrap(col.Name) + " " + g.getType(col)

	if col.AutoIncrement {
		// Use serial types for auto-increment in PostgreSQL.
		switch col.Type {
		case "bigInteger":
			sql = g.wrap(col.Name) + " bigserial primary key"
		case "integer":
			sql = g.wrap(col.Name) + " serial primary key"
		case "smallInteger":
			sql = g.wrap(col.Name) + " smallserial primary key"
		default:
			sql = g.wrap(col.Name) + " bigserial primary key"
		}

		return sql
	}

	if !col.IsNullable {
		sql += " not null"
	}

	if col.HasDefault {
		sql += " default " + g.getDefaultValue(col.DefaultValue)
	}

	return sql
}

func (g *PostgresGrammar) getType(col *schema.ColumnDefinition) string {
	switch col.Type {
	case "bigInteger":
		return "bigint"
	case "mediumInteger", "integer":
		return "integer"
	case "smallInteger":
		return "smallint"
	case "tinyInteger":
		return "smallint"
	case "string":
		return fmt.Sprintf("varchar(%d)", col.Length)
	case "char":
		return fmt.Sprintf("char(%d)", col.Length)
	case "text", "tinyText", "mediumText", "longText":
		return "text"
	case "float":
		return "double precision"
	case "double":
		return "double precision"
	case "decimal":
		return fmt.Sprintf("decimal(%d, %d)", col.Precision, col.Scale)
	case "boolean":
		return "boolean"
	case "date":
		return "date"
	case "dateTime":
		if col.Precision > 0 {
			return fmt.Sprintf("timestamp(%d) without time zone", col.Precision)
		}

		return "timestamp without time zone"
	case "dateTimeTz":
		if col.Precision > 0 {
			return fmt.Sprintf("timestamp(%d) with time zone", col.Precision)
		}

		return "timestamp with time zone"
	case "time":
		if col.Precision > 0 {
			return fmt.Sprintf("time(%d) without time zone", col.Precision)
		}

		return "time without time zone"
	case "timeTz":
		if col.Precision > 0 {
			return fmt.Sprintf("time(%d) with time zone", col.Precision)
		}

		return "time with time zone"
	case "timestamp":
		if col.Precision > 0 {
			return fmt.Sprintf("timestamp(%d) without time zone", col.Precision)
		}

		return "timestamp without time zone"
	case "timestampTz":
		if col.Precision > 0 {
			return fmt.Sprintf("timestamp(%d) with time zone", col.Precision)
		}

		return "timestamp with time zone"
	case "year":
		return "integer"
	case "binary":
		return "bytea"
	case "json":
		return "json"
	case "jsonb":
		return "jsonb"
	case "uuid":
		return "uuid"
	case "ipAddress":
		return "inet"
	case "macAddress":
		return "macaddr"
	case "enum":
		return "varchar(255)"
	case "set":
		return "varchar(255)"
	case "geometry":
		return "geometry"
	case "geography":
		return "geography"
	case "vector":
		if col.Total > 0 {
			return fmt.Sprintf("vector(%d)", col.Total)
		}

		return "vector"
	default:
		return col.Type
	}
}

func (g *PostgresGrammar) getDefaultValue(value any) string {
	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case bool:
		if v {
			return "'true'"
		}

		return "'false'"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func (g *PostgresGrammar) wrap(value string) string {
	return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
}

func (g *PostgresGrammar) wrapTable(table string) string {
	return g.wrap(g.tablePrefix + table)
}
