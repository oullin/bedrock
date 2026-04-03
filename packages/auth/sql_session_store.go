package auth

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLSessionStore persists sessions in SQL.
type SQLSessionStore struct {
	db    *sql.DB
	table string
}

// NewSQLSessionStore creates a new SQL session store.
func NewSQLSessionStore(db *sql.DB, table string) *SQLSessionStore {
	if table == "" {
		table = "auth_sessions"
	}

	return &SQLSessionStore{db: db, table: table}
}

// Create inserts a session.
func (s *SQLSessionStore) Create(ctx context.Context, session *Session) error {
	if session == nil {
		return fmt.Errorf("auth: session is required")
	}

	_, err := s.db.ExecContext(ctx, `INSERT INTO `+s.table+` (id, user_id, pending_two_factor, pending_remember, password_confirmed_at, authenticated_at, last_seen_at, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		session.ID,
		session.UserID,
		boolToInt(session.PendingTwoFactor),
		boolToInt(session.PendingRemember),
		nil,
		session.AuthenticatedAt,
		session.LastSeenAt,
		session.CreatedAt,
		session.ExpiresAt,
	)

	return err
}

// FindByID retrieves a session by id.
func (s *SQLSessionStore) FindByID(ctx context.Context, id string) (*Session, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, user_id, pending_two_factor, pending_remember, authenticated_at, last_seen_at, created_at, expires_at FROM `+s.table+` WHERE id = ?`, id)

	var session Session

	var pendingTwoFactor int

	var pendingRemember int

	var authenticatedAt sql.NullTime

	if err := row.Scan(&session.ID, &session.UserID, &pendingTwoFactor, &pendingRemember, &authenticatedAt, &session.LastSeenAt, &session.CreatedAt, &session.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUnauthorized
		}

		return nil, err
	}

	session.PendingTwoFactor = pendingTwoFactor == 1
	session.PendingRemember = pendingRemember == 1

	if authenticatedAt.Valid {
		value := authenticatedAt.Time
		session.AuthenticatedAt = &value
	}

	return &session, nil
}

// Update persists a changed session.
func (s *SQLSessionStore) Update(ctx context.Context, session *Session) error {
	if session == nil {
		return fmt.Errorf("auth: session is required")
	}

	result, err := s.db.ExecContext(ctx, `UPDATE `+s.table+` SET user_id = ?, pending_two_factor = ?, pending_remember = ?, authenticated_at = ?, last_seen_at = ?, created_at = ?, expires_at = ? WHERE id = ?`,
		session.UserID,
		boolToInt(session.PendingTwoFactor),
		boolToInt(session.PendingRemember),
		session.AuthenticatedAt,
		session.LastSeenAt,
		session.CreatedAt,
		session.ExpiresAt,
		session.ID,
	)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrUnauthorized
	}

	return nil
}

// Delete removes a session.
func (s *SQLSessionStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM `+s.table+` WHERE id = ?`, id)

	return err
}

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}
