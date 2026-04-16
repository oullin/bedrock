package schema

import (
	"context"

	dbcontract "github.com/bedrock/packages/contracts/database"
)

// Grammar compiles Blueprint instances into DDL statements.
type Grammar interface {
	CompileCreate(bp *Blueprint) []string
	CompileAdd(bp *Blueprint) []string
	CompileChange(bp *Blueprint) []string
	CompileDrop(table string) string
	CompileDropIfExists(table string) string
	CompileRename(from, to string) string
	CompileDropColumn(bp *Blueprint, columns []string) string
	CompileRenameColumn(bp *Blueprint, from, to string) string
	CompileCreateIndex(bp *Blueprint, cmd BlueprintCommand) string
	CompileDropIndex(bp *Blueprint, name string) string
	CompileCreateForeignKey(bp *Blueprint, fk *ForeignKeyDefinition) string
	CompileDropForeignKey(bp *Blueprint, name string) string
	CompileTableExists() string
	CompileColumnListing(table string) string
	CompileEnableForeignKeyConstraints() string
	CompileDisableForeignKeyConstraints() string
}

// Builder provides methods for creating and modifying database tables.
// It mirrors Framework\Database\Schema\Builder.
type Builder struct {
	connection dbcontract.Connection
	grammar    Grammar
}

// NewBuilder creates a new schema Builder.
func NewBuilder(connection dbcontract.Connection, grammar Grammar) *Builder {
	return &Builder{connection: connection, grammar: grammar}
}

// Create creates a new table on the schema.
func (b *Builder) Create(ctx context.Context, table string, callback func(*Blueprint)) error {
	bp := NewBlueprint(table)
	callback(bp)

	statements := b.grammar.CompileCreate(bp)
	for _, sql := range statements {
		_, err := b.connection.Statement(ctx, sql)
		if err != nil {
			return err
		}
	}

	// Create indexes.
	for _, cmd := range bp.Commands {
		sql := b.compileCommand(bp, cmd)
		if sql != "" {
			_, err := b.connection.Statement(ctx, sql)
			if err != nil {
				return err
			}
		}
	}

	// Create foreign keys.
	for _, fk := range bp.ForeignKeys {
		sql := b.grammar.CompileCreateForeignKey(bp, fk)
		if sql != "" {
			_, err := b.connection.Statement(ctx, sql)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Table modifies an existing table.
func (b *Builder) Table(ctx context.Context, table string, callback func(*Blueprint)) error {
	bp := NewBlueprint(table)
	callback(bp)

	// Add new columns.
	if added := bp.GetAddedColumns(); len(added) > 0 {
		for _, sql := range b.grammar.CompileAdd(bp) {
			_, err := b.connection.Statement(ctx, sql)
			if err != nil {
				return err
			}
		}
	}

	// Change existing columns.
	if changed := bp.GetChangedColumns(); len(changed) > 0 {
		for _, sql := range b.grammar.CompileChange(bp) {
			_, err := b.connection.Statement(ctx, sql)
			if err != nil {
				return err
			}
		}
	}

	// Execute commands (indexes, drops, renames, etc.).
	for _, cmd := range bp.Commands {
		sql := b.compileCommand(bp, cmd)
		if sql != "" {
			_, err := b.connection.Statement(ctx, sql)
			if err != nil {
				return err
			}
		}
	}

	// Create foreign keys.
	for _, fk := range bp.ForeignKeys {
		sql := b.grammar.CompileCreateForeignKey(bp, fk)
		if sql != "" {
			_, err := b.connection.Statement(ctx, sql)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Drop drops a table.
func (b *Builder) Drop(ctx context.Context, table string) error {
	_, err := b.connection.Statement(ctx, b.grammar.CompileDrop(table))
	return err
}

// DropIfExists drops a table if it exists.
func (b *Builder) DropIfExists(ctx context.Context, table string) error {
	_, err := b.connection.Statement(ctx, b.grammar.CompileDropIfExists(table))
	return err
}

// Rename renames a table.
func (b *Builder) Rename(ctx context.Context, from, to string) error {
	_, err := b.connection.Statement(ctx, b.grammar.CompileRename(from, to))
	return err
}

// HasTable checks if a table exists.
func (b *Builder) HasTable(ctx context.Context, table string) (bool, error) {
	sql := b.grammar.CompileTableExists()
	rows, err := b.connection.Select(ctx, sql, table)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

// HasColumn checks if a column exists on a table.
func (b *Builder) HasColumn(ctx context.Context, table, column string) (bool, error) {
	columns, err := b.GetColumnListing(ctx, table)
	if err != nil {
		return false, err
	}
	for _, col := range columns {
		if col == column {
			return true, nil
		}
	}
	return false, nil
}

// GetColumnListing returns the column names for a table.
func (b *Builder) GetColumnListing(ctx context.Context, table string) ([]string, error) {
	sql := b.grammar.CompileColumnListing(table)
	rows, err := b.connection.Select(ctx, sql)
	if err != nil {
		return nil, err
	}
	var columns []string
	for _, row := range rows {
		if name, ok := row["column_name"]; ok {
			if s, ok := name.(string); ok {
				columns = append(columns, s)
			}
		}
	}
	return columns, nil
}

// EnableForeignKeyConstraints enables foreign key constraints.
func (b *Builder) EnableForeignKeyConstraints(ctx context.Context) error {
	_, err := b.connection.Statement(ctx, b.grammar.CompileEnableForeignKeyConstraints())
	return err
}

// DisableForeignKeyConstraints disables foreign key constraints.
func (b *Builder) DisableForeignKeyConstraints(ctx context.Context) error {
	_, err := b.connection.Statement(ctx, b.grammar.CompileDisableForeignKeyConstraints())
	return err
}

// GetConnection returns the underlying connection.
func (b *Builder) GetConnection() dbcontract.Connection { return b.connection }

// GetGrammar returns the schema grammar.
func (b *Builder) GetGrammar() Grammar { return b.grammar }

func (b *Builder) compileCommand(bp *Blueprint, cmd BlueprintCommand) string {
	switch cmd.Name {
	case "primary", "unique", "index", "fulltext":
		return b.grammar.CompileCreateIndex(bp, cmd)
	case "dropColumn":
		return b.grammar.CompileDropColumn(bp, cmd.Columns)
	case "renameColumn":
		if len(cmd.Columns) >= 2 {
			return b.grammar.CompileRenameColumn(bp, cmd.Columns[0], cmd.Columns[1])
		}
	case "dropPrimary", "dropUnique", "dropIndex", "dropFulltext":
		return b.grammar.CompileDropIndex(bp, cmd.Index)
	case "dropForeign":
		return b.grammar.CompileDropForeignKey(bp, cmd.Index)
	}
	return ""
}
