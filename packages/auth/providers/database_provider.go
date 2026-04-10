package providers

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/auth"
)

// DBQuerier is the minimal raw-SQL interface for DatabaseUserProvider.
type DBQuerier interface {
	// QueryRow executes a query returning at most one row.
	QueryRow(ctx context.Context, query string, args ...any) DBRow
	// Exec executes a statement returning no rows.
	Exec(ctx context.Context, query string, args ...any) error
}

// DBRow is a single result row.
type DBRow interface {
	Scan(dest ...any) error
}

// RowMapper converts a scanned row map into an Authenticatable.
type RowMapper func(row map[string]any) auth.Authenticatable

// DatabaseUserProvider retrieves users from a raw SQL table.
type DatabaseUserProvider struct {
	db         DBQuerier
	table      string
	hasher     auth.PasswordHasher
	rowMapper  RowMapper
}

// NewDatabaseUserProvider creates a DatabaseUserProvider.
// table is the users table name. rowMapper converts a row map to Authenticatable.
func NewDatabaseUserProvider(db DBQuerier, table string, hasher auth.PasswordHasher, rowMapper RowMapper) *DatabaseUserProvider {
	return &DatabaseUserProvider{
		db:        db,
		table:     table,
		hasher:    hasher,
		rowMapper: rowMapper,
	}
}

func (p *DatabaseUserProvider) RetrieveByID(ctx context.Context, id any) (auth.Authenticatable, error) {
	row := p.db.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id = $1 LIMIT 1", p.table), id)

	return p.mapRow(row)
}

func (p *DatabaseUserProvider) RetrieveByToken(ctx context.Context, id any, token string) (auth.Authenticatable, error) {
	row := p.db.QueryRow(ctx,
		fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND remember_token = $2 LIMIT 1", p.table),
		id, token,
	)

	return p.mapRow(row)
}

func (p *DatabaseUserProvider) UpdateRememberToken(ctx context.Context, user auth.Authenticatable, token string) error {
	return p.db.Exec(ctx,
		fmt.Sprintf("UPDATE %s SET remember_token = $1 WHERE id = $2", p.table),
		token, user.GetAuthIdentifier(),
	)
}

func (p *DatabaseUserProvider) RetrieveByCredentials(ctx context.Context, credentials map[string]any) (auth.Authenticatable, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ", p.table)

	args := make([]any, 0, len(credentials))
	i := 1

	for k, v := range credentials {
		if k == "password" {
			continue
		}

		if i > 1 {
			query += " AND "
		}

		query += fmt.Sprintf("%s = $%d", k, i)
		args = append(args, v)
		i++
	}

	query += " LIMIT 1"
	row := p.db.QueryRow(ctx, query, args...)

	return p.mapRow(row)
}

func (p *DatabaseUserProvider) ValidateCredentials(_ context.Context, user auth.Authenticatable, credentials map[string]any) bool {
	plain, ok := credentials["password"].(string)
	if !ok {
		return false
	}

	return p.hasher.Check(plain, user.GetAuthPassword())
}

func (p *DatabaseUserProvider) RehashPasswordIfRequired(ctx context.Context, user auth.Authenticatable, credentials map[string]any, force bool) error {
	if !force && !p.hasher.NeedsRehash(user.GetAuthPassword()) {
		return nil
	}

	plain, ok := credentials["password"].(string)
	if !ok {
		return nil
	}

	hash, err := p.hasher.Hash(plain)
	if err != nil {
		return err
	}

	return p.db.Exec(ctx,
		fmt.Sprintf("UPDATE %s SET password = $1 WHERE id = $2", p.table),
		hash, user.GetAuthIdentifier(),
	)
}

func (p *DatabaseUserProvider) mapRow(row DBRow) (auth.Authenticatable, error) {
	// Scan into a map via column names is not directly supported by the DBRow
	// interface; callers must provide a RowMapper that matches their DB driver's
	// row type. Here we delegate to the injected mapper.
	//
	// When the row mapper cannot be called directly (scan requires concrete types),
	// callers should embed the concrete row type and cast appropriately. This
	// interface-based approach matches the driver-agnostic design of this package.
	//
	// For a concrete example using database/sql, wrap *sql.Row so Scan returns
	// column values into a []any and then build the map before calling rowMapper.
	var dest map[string]any
	if err := row.Scan(&dest); err != nil {
		return nil, nil // Row not found.
	}

	return p.rowMapper(dest), nil
}
