package passwords

import (
	"context"
	"database/sql"
	"time"

	auth "github.com/gollin/packages/auth"
)

// SQLTokenRepository stores reset tokens in SQL.
type SQLTokenRepository struct {
	db    *sql.DB
	table string
}

// NewSQLTokenRepository creates a SQL-backed token repository.
func NewSQLTokenRepository(db *sql.DB, table string) *SQLTokenRepository {
	if table == "" {
		table = "password_reset_tokens"
	}

	return &SQLTokenRepository{db: db, table: table}
}

// Save inserts or replaces a token.
func (r *SQLTokenRepository) Save(ctx context.Context, token *Token) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO `+r.table+` (user_id, token_hash, created_at, expires_at) VALUES (?, ?, ?, ?) ON CONFLICT(user_id) DO UPDATE SET token_hash = excluded.token_hash, created_at = excluded.created_at, expires_at = excluded.expires_at`,
		token.UserID,
		token.TokenHash,
		token.CreatedAt,
		token.ExpiresAt,
	)

	return err
}

// FindByTokenHash returns a token by hash.
func (r *SQLTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*Token, error) {
	row := r.db.QueryRowContext(ctx, `SELECT user_id, token_hash, created_at, expires_at FROM `+r.table+` WHERE token_hash = ?`, tokenHash)

	var token Token

	if err := row.Scan(&token.UserID, &token.TokenHash, &token.CreatedAt, &token.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrInvalidToken
		}

		return nil, err
	}

	return &token, nil
}

// DeleteByTokenHash removes a token.
func (r *SQLTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM `+r.table+` WHERE token_hash = ?`, tokenHash)

	return err
}

// DeleteByUserID removes a token.
func (r *SQLTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM `+r.table+` WHERE user_id = ?`, userID)

	return err
}

// RecentlyCreated reports whether a token exists since the provided time.
func (r *SQLTokenRepository) RecentlyCreated(ctx context.Context, userID string, since time.Time) (bool, error) {
	var createdAt time.Time
	err := r.db.QueryRowContext(ctx, `SELECT created_at FROM `+r.table+` WHERE user_id = ?`, userID).Scan(&createdAt)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return createdAt.After(since) || createdAt.Equal(since), nil
}
